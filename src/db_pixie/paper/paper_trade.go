package paper

import (
	"appcfg"
	"constants"
	"dao_pixie"
	"database/sql"
	. "logger"
	. "pixie_contract/api_specification"
	"tools"
)

func BuyPaper(tradeID, paperID int64, oriSequence int, tNow int64, username, nickname, head string, sex int, price float64, priceType int) (err error) {
	var tx *sql.Tx
	var res sql.Result
	var idInsert, effect int64

	if tx, err = dao_pixie.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	var dp DesignPaper
	if err = tx.Stmt(dao_pixie.GetOneDesignPaperStmt).QueryRow(paperID).Scan(&dp.PaperID, &dp.AuthorUsername, &dp.AuthorNickname, &dp.AuthorHead, &dp.AuthorSex, &dp.OwnerUsername, &dp.OwnerNickname, &dp.OwnerHead, &dp.OwnerSex, &dp.ClothesType, &dp.PartType, &dp.Cname, &dp.Desc, &dp.UploadTime, &dp.Extra, &dp.File, &dp.Style, &dp.Status, &dp.AuctionFailUnread, &dp.Star, &dp.Tag1, &dp.Tag2, &dp.STag, &dp.PassTime, &dp.VerifiedCount, &dp.PriceType, &dp.LastDuration, &dp.LastMaxPrice, &dp.LastMinPrice, &dp.CopyMark, &dp.Reason, &dp.CirculationMax, &dp.Circulation); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("can't find paper for submit trade", paperID)
		} else {
			Err(err)
		}
		return
	}

	//从trade查询图纸信息
	var sp SalePaper
	if err = tx.Stmt(dao_pixie.GetPaperTradeByIDStmt).QueryRow(tradeID).Scan(&sp.TradeID, &sp.AuthorUsername, &sp.AuthorNickname, &sp.AuthorHead, &sp.AuthorSex, &sp.OwnerUsername, &sp.OwnerNickname, &sp.OwnerHead, &sp.OwnerSex, &sp.PaperID, &sp.ClothesType, &sp.PartType, &sp.Cname, &sp.Desc, &sp.Extra, &sp.File, &sp.Style, &sp.Star, &sp.Tag1, &sp.Tag2, &sp.STag, &sp.MaxPrice, &sp.MinPrice, &sp.StartTime, &sp.Duration, &sp.PriceType, &sp.AuctionTime, &sp.AuctionUsername, &sp.AuctionNickname, &sp.AuctionPrice, &sp.Sequence, &sp.CirculationMax); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Err("can't find trade for", paperID, tradeID)
			err = constants.TradeNotExist
		} else {
			Err(err)
		}

		return
	}

	if oriSequence != sp.Sequence {
		//校验客户端传来的sequence
		err = constants.SequenceNotMatchTrade
		Err("sequence not match")
		return
	}

	if tNow >= sp.StartTime+int64(sp.Duration*3600) {
		//竞拍已超时
		Info("trade end when buy paper", paperID)
		err = constants.TradeNotExist
		return
	}

	if ok, lp := tools.AuctionPriceLegal(sp.MaxPrice, sp.MinPrice, sp.StartTime, sp.Duration*3600, tNow, price); !ok {
		Err("auction price illegal", username, paperID, lp, price, sp.MaxPrice, sp.MinPrice, sp.Duration, sp.StartTime, tNow)
		err = constants.AuctionPriceIllegal
		return
	}

	var sequence int

	if sp.Sequence <= 0 {
		//设计师发布
		if dp.Circulation >= dp.CirculationMax {
			err = constants.CirculationLimit
			Err(err)
			return
		} else {
			sequence = dp.Circulation + 1
		}
	} else {
		sequence = sp.Sequence
	}

	if sp.Sequence > 0 || sequence >= dp.CirculationMax {
		//普通购买图纸 或 发布图纸数量到达上限 删除pixie_trade记录
		if res, err = tx.Stmt(dao_pixie.DeletePaperTradeByIDStmt).Exec(sp.TradeID); err != nil {
			Err(err)
			return
		} else if effect, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		} else if effect <= 0 {
			Err("delete or update trade for paper when buy paper fail", paperID, tradeID, "effect 0")
			err = constants.PaperSqlEffectZero
			return
		}
	} else {
		//更新trade状态
		if res, err = tx.Stmt(dao_pixie.UpdatePaperTradeStatusStmt).Exec(constants.PIXIE_TRADE_STATUS_SUCCESS, tradeID); err != nil {
			Err(err)
			return
		}
	}

	//插入pixie_trade_log
	if res, err = tx.Stmt(dao_pixie.InsertTradeLogStmt).Exec(paperID, sequence, username, nickname, tNow, price, sp.PriceType); err != nil {
		Err(err)
		return
	} else if idInsert, err = res.LastInsertId(); err != nil {
		Err(err)
		return
	} else if idInsert < 1 {
		Err("insert trade log get id wrong", idInsert, paperID, username)
		err = constants.PaperTradeEffectZero
		return
	}

	if sp.OwnerUsername != username {
		//插入
		res, err = tx.Stmt(dao_pixie.AddOccupyPaperStmt).Exec(paperID, sp.AuthorUsername, sp.AuthorNickname, sp.AuthorHead, sp.AuthorSex, username, nickname, head, sex, sp.ClothesType, sp.PartType, sp.Cname, sp.Desc, sp.Extra, sp.File, sp.Style, sp.Star, sp.Tag1, sp.Tag2, sp.STag, PAPER_STATUS_OCCUPY, idInsert, sequence, dp.CirculationMax)
	} else {
		//更新
		res, err = tx.Stmt(dao_pixie.UpdateOccupyPaperStmt).Exec(PAPER_STATUS_OCCUPY, idInsert, dp.CirculationMax, paperID, username, sequence)
	}

	//插入或更新paper occupy
	if err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect <= 0 {
		Err("add or update occupy for paper when buy paper fail", paperID, tradeID, "effect 0")
		err = constants.PaperSqlEffectZero
		return
	}

	if appcfg.GetBool("fake_buy_paper_with_product", false) {
		if res, err = tx.Stmt(dao_pixie.FakeProductPaperStmt).Exec(paperID, username); err != nil {
			Err(err)
			return
		} else if effect, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		} else if effect <= 0 {
			Err("update fake product for paper", paperID, "effect 0")
			err = constants.PaperSqlEffectZero
			return
		}
	}

	if sp.Sequence == 0 {
		//设计师发布图纸 只需更新pixie_paper
		toStatus := PAPER_STATUS_ON_SALE
		if sequence >= dp.CirculationMax {
			toStatus = PAPER_STATUS_OCCUPY
		}

		res, err = tx.Stmt(dao_pixie.UpdatePaperStatusStmt).Exec(sequence, toStatus, paperID, PAPER_STATUS_ON_SALE)
	} else if sp.OwnerUsername != username {
		//正常图纸流传 需删除原paper_occupy
		res, err = tx.Stmt(dao_pixie.UnoccupyPaperStmt).Exec(paperID, sp.OwnerUsername)
	}

	if err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect <= 0 {
		Err("update paper or paper occupy when buy paper", paperID, tradeID, "effect 0")
		err = constants.PaperSqlEffectZero
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}

	if sp.Sequence == 0 && sequence >= dp.CirculationMax {
		//更新redis
		UpdatePaperStatusRedis(paperID, PAPER_STATUS_OCCUPY)
	}

	return
}

func SubmitPaperTrade(submitUsername, submitNickname, submitHead string, submitSex int, maxPrice, minPrice int, startTime, duration int64, moneyType int, paperID int64, circulationMax, sequence int) (tid int64, err error) {
	var tx *sql.Tx
	var res sql.Result

	if tx, err = dao_pixie.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	var t DesignPaper
	if err = tx.Stmt(dao_pixie.GetOneDesignPaperStmt).QueryRow(paperID).Scan(&t.PaperID, &t.AuthorUsername, &t.AuthorNickname, &t.AuthorHead, &t.AuthorSex, &t.OwnerUsername, &t.OwnerNickname, &t.OwnerHead, &t.OwnerSex, &t.ClothesType, &t.PartType, &t.Cname, &t.Desc, &t.UploadTime, &t.Extra, &t.File, &t.Style, &t.Status, &t.AuctionFailUnread, &t.Star, &t.Tag1, &t.Tag2, &t.STag, &t.PassTime, &t.VerifiedCount, &t.PriceType, &t.LastDuration, &t.LastMaxPrice, &t.LastMinPrice, &t.CopyMark, &t.Reason, &t.CirculationMax, &t.Circulation); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("can't find paper for submit trade", paperID)
		} else {
			Err(err)
		}
		return
	}

	if sequence == 0 {
		//设计师发布图纸 需要检查最大流通量
		if t.CirculationMax <= 0 {
			//未设置过最大流通量
			if circulationMax <= 0 {
				//未提交最大流通量
				err = constants.CirculationNotSet
				Err(err)
				return
			}
		} else if t.Circulation > 0 && t.CirculationMax != circulationMax {
			//设置过流通量 流通过后 不能再更改最大流通量
			err = constants.CirculationChange
			Err(err)
			return
		}

		res, err = tx.Stmt(dao_pixie.PublishPaperStmt).Exec(t.AuthorUsername, t.AuthorNickname, t.AuthorHead, t.AuthorSex, "", "", "", 0, paperID, t.ClothesType, t.PartType, t.Cname, t.Desc, t.Extra, t.File, t.Style, t.Star, t.Tag1, t.Tag2, t.STag, maxPrice, minPrice, startTime, duration, moneyType, 0, "", "", 0, sequence, circulationMax)
	} else {
		//普通出售图纸
		res, err = tx.Stmt(dao_pixie.PublishPaperStmt).Exec(t.AuthorUsername, t.AuthorNickname, t.AuthorHead, t.AuthorSex, submitUsername, submitNickname, submitHead, submitSex, paperID, t.ClothesType, t.PartType, t.Cname, t.Desc, t.Extra, t.File, t.Style, t.Star, t.Tag1, t.Tag2, t.STag, maxPrice, minPrice, startTime, duration, moneyType, 0, "", "", 0, sequence, t.CirculationMax)
	}

	if err != nil {
		Err(err)
		return
	} else if tid, err = res.LastInsertId(); err != nil {
		Err(err)
		return
	}

	if sequence == 0 {
		//设计师发布
		res, err = tx.Stmt(dao_pixie.UpdatePaperWhenStartTradeStmt).Exec(PAPER_STATUS_ON_SALE, moneyType, duration, maxPrice, minPrice, circulationMax, paperID, PAPER_STATUS_WAIT_PUBLISH)
	} else {
		//正常竞拍
		res, err = tx.Stmt(dao_pixie.UpdateOccupyWhenStartTradeStmt).Exec(PAPER_STATUS_ON_SALE, moneyType, duration, maxPrice, minPrice, paperID, PAPER_STATUS_OCCUPY, submitUsername)
	}

	var effect int64
	if err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect <= 0 {
		err = constants.PaperSqlEffectZero
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}

	//更新redis
	UpdatePaperStatusRedis(paperID, PAPER_STATUS_ON_SALE)

	return
}

func CancelPaperTrade(paperID int64, username string, nt int64, tradeID int64) (err error) {
	var tx *sql.Tx
	var res sql.Result
	var num int64
	if tx, err = dao_pixie.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()
	var ownerUsername string
	var startTime int64
	var duration, sequence int

	//author
	if err = tx.Stmt(dao_pixie.GetSimplePaperTradeByIDStmt).QueryRow(paperID, tradeID).Scan(&ownerUsername, &startTime, &duration, &sequence); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("can't find trade for", paperID)
			err = constants.TradeNotExist
		} else {
			Err(err)
		}
		return
	}

	if sequence > 0 {
		if ownerUsername != username {
			err = constants.PaperAuthorOrOwnerWrong
			return
		}

		if res, err = tx.Stmt(dao_pixie.UpdateOccupyWhenCancelTradeStmt).Exec(PAPER_STATUS_OCCUPY, paperID, username, PAPER_STATUS_ON_SALE); err != nil {
			Err(err)
			return
		} else if num, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		}
	} else {
		if res, err = tx.Stmt(dao_pixie.UpdatePaperWhenCancelTradeStmt).Exec(PAPER_STATUS_WAIT_PUBLISH, paperID, username, PAPER_STATUS_ON_SALE); err != nil {
			Err(err)
			return
		} else if num, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		}
	}

	if num <= 0 {
		Err("paper status not right when cancel trade", paperID, username)
		err = constants.PaperStatusWrong
		return
	}

	//删除trade
	if res, err = tx.Stmt(dao_pixie.DeletePaperTradeByIDStmt).Exec(tradeID); err != nil {
		Err(err)
		return
	} else if num, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if num <= 0 {
		Err("delete paper trade effect 0", paperID, username)
		err = constants.PaperSqlEffectZero
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func DeletePaperTradeWhenTimeout(tradeID, paperID int64, firstTrade bool, originUsername string, tradeStatus int) {
	var tx *sql.Tx
	var res sql.Result
	var err error
	var num int64
	if tx, err = dao_pixie.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	//删除trade表
	if res, err = tx.Stmt(dao_pixie.DeletePaperTradeByIDStmt).Exec(tradeID); err != nil {
		Err(err)
		return
	} else if num, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if num < 1 {
		Err("DeletePaperTradeWhenTimeout effect 0", tradeID, paperID, originUsername)
		return
	}

	failUnread := 1

	if tradeStatus == constants.PIXIE_TRADE_STATUS_SUCCESS {
		failUnread = 0
	}

	if firstTrade {
		//更新Paper表
		if res, err = tx.Stmt(dao_pixie.UpdatePaperWhenTradeTimeoutStmt).Exec(PAPER_STATUS_WAIT_PUBLISH, failUnread, paperID, PAPER_STATUS_ON_SALE, originUsername); err != nil {
			Err(err)
			return
		} else if num, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		} else if num < 1 {
			Err("DeletePaperTradeWhenTimeout first trade effect 0", paperID, originUsername)
			return
		}
	} else {
		//更新Occupy
		if res, err = tx.Stmt(dao_pixie.UpdateOccupyWhenTradeTimeoutStmt).Exec(PAPER_STATUS_OCCUPY, failUnread, paperID, PAPER_STATUS_ON_SALE, originUsername); err != nil {
			Err(err)
			return
		} else if num, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		} else if num < 1 {
			Err("DeletePaperTradeWhenTimeout not first trade effect 0", paperID, originUsername)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
}
