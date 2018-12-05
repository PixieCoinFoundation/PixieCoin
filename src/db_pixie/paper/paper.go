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

func GetAdminCheckCountBefore(paperID int64) (cnt int, err error) {
	if err = dao_pixie.GetStatusCountBeforeIDStmt.QueryRow(PAPER_STATUS_ADMIN_QUEUE, paperID).Scan(&cnt); err != nil {
		Err(err)
	}
	return
}

func QueryPaperSaleInfo(paperID int64, username string, sequence int) (tradeID int64, priceType int, startTime int64, duration int64, maxPrice int64, minPrice int64, err error) {
	if err = dao_pixie.QueryPaperSaleInfoStmt.QueryRow(paperID, sequence).Scan(&tradeID, &priceType, &startTime, &duration, &maxPrice, &minPrice); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("not on sale", paperID)
			err = constants.TradeNotExist
		} else {
			Err(err)
		}
	}
	return
}

func AddPaper(username, nickname, head string, sex int, clothesType, partType, cname, desc string, uploadTime int64, paperFile, paperExtra string) (err error, paperID int64) {
	var res sql.Result
	if res, err = dao_pixie.UploadPaperStmt.Exec(username, nickname, head, sex, clothesType, partType, cname, desc, uploadTime, paperFile, paperExtra, PAPER_STATUS_ADMIN_QUEUE); err != nil {
		Err(err)
		return
	} else if paperID, err = res.LastInsertId(); err != nil {
		Err(err)
		return
	} else {
		AddPaperRedis(paperID, username, nickname, head, sex, clothesType, partType, cname, desc, paperExtra, paperFile, PAPER_STATUS_ADMIN_QUEUE)
	}

	return
}

func ListDesignPaper(start, size int) (err error, list []BasePaper) {
	var rows *sql.Rows
	if rows, err = dao_pixie.ListDesignPaperByPageStmt.Query(start, size); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		list = make([]BasePaper, 0)
		for rows.Next() {
			var temp BasePaper
			if err = rows.Scan(&temp.PaperID, &temp.AuthorUsername, &temp.AuthorNickname, &temp.AuthorHead, &temp.AuthorSex, &temp.OwnerUsername, &temp.OwnerNickname, &temp.OwnerHead, &temp.OwnerSex, &temp.ClothesType, &temp.PartType, &temp.Cname, &temp.Desc, &temp.Extra, &temp.Style, &temp.Status, &temp.Star, &temp.Tag1, &temp.Tag2, &temp.STag); err != nil {
				Err(err)
				return
			}
			list = append(list, temp)
		}
	}
	return
}

func DesignPaperInList(idList []interface{}) (err error, list []RecommendClothes) {
	var rows *sql.Rows
	if rows, err = dao_pixie.RecommendPaperInListStmt.Query(idList...); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		return getRecommendPapers(rows)
	}
	return
}

func ListPaperByAuthor(username string, page, pageSize int) (err error, list []DesignPaper) {
	var rows *sql.Rows
	if rows, err = dao_pixie.GetPaperByAuthorStmt.Query(username, (page-1)*pageSize, pageSize); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		return getDesignPapers(rows, appcfg.GetBool("exclude_freeze", false), int(PAPER_STATUS_FREEZE))
	}
	return
}

//获取已购图纸
func ListPaperByOwner(username, clothesType, partType string, page, pageSize int, tNow int64) (err error, list []*OwnPaper) {
	var rows *sql.Rows

	if partType != "" {
		if clothesType == "" {
			rows, err = dao_pixie.ListPaperByPartTypeOwnerStmt.Query(username, partType, (page-1)*pageSize, pageSize)
		} else {
			rows, err = dao_pixie.ListPaperByClothesTypePartTypeOwnerStmt.Query(username, clothesType, partType, (page-1)*pageSize, pageSize)
		}
	} else {
		if clothesType == "" {
			rows, err = dao_pixie.ListPaperByOwnerStmt.Query(username, (page-1)*pageSize, pageSize)
		} else {
			rows, err = dao_pixie.ListPaperByOwnerClothesTypeStmt.Query(username, clothesType, (page-1)*pageSize, pageSize)
		}
	}

	if err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		var tempList []*OwnPaper
		if err, tempList = getOwnPapers(rows); err != nil {
			return
		} else {
			for _, v := range tempList {
				if v.IsInProduction {
					productCnt := tools.CheckProductNum(tNow, v.ProductStartTime, v.ProductLeftCount, v.ProductDoneCount)
					if productCnt > 0 {
						if e := DoProductPaper(username, productCnt, v.ProductStartTime, v.ProductLeftCount-productCnt, v.ProductDoneCount+productCnt, v.PaperID, v.ProductLeftCount, v.ProductDoneCount); e == nil {
							v.ProductLeftCount -= productCnt
							v.ProductDoneCount += productCnt
							v.Inventory += productCnt

							if v.ProductLeftCount == 0 {
								v.IsInProduction = false
								v.ProductStartTime = 0
							}
						}
					}

				}
				list = append(list, v)
			}
		}
	}
	return
}

func ListPaperOwnByOwnerIDList(idList []int64, ownerUsername string) (err error, list []*OwnPaper) {
	argList := make([]interface{}, 0)

	argList = append(argList, ownerUsername)

	for _, v := range idList {
		argList = append(argList, v)
	}

	for len(argList) < constants.PIXIE_LAND_MAX_SALE_SIZE+1 {
		argList = append(argList, 0)
	}

	return getPaperOwnList(dao_pixie.ListPaperOwnByOwnerIdListStmt, argList)
}

func ListPaperOwnByIDMap(idMap map[int64]*SimpleClothes) (err error, list []*OwnPaper) {
	argList := make([]interface{}, 0)

	for v, _ := range idMap {
		argList = append(argList, v)
	}

	for len(argList) < constants.PIXIE_LAND_MAX_SALE_SIZE {
		argList = append(argList, 0)
	}

	return getPaperOwnList(dao_pixie.ListPaperOwnByIdListStmt, argList)
}

func getPaperOwnList(stmt *sql.Stmt, argList []interface{}) (err error, list []*OwnPaper) {
	var rows *sql.Rows
	if rows, err = stmt.Query(argList...); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		return getOwnPapers(rows)
	}
	return
}

func GetPaperSaleByUsernameDate(username, all, today, yesterday string) (err error, todayPxcSale, yesterdayPxcSale, allPxcSale float64, todayGoldSale, yesterdayGoldSale, allGoldSale int64, todaySaleClothes, yesterdaySaleClothes, allSaleClothes int64) {
	var rows *sql.Rows
	if rows, err = dao_pixie.GetPaperSaleByUsernameDateStmt.Query(username, all, today, yesterday); err != nil {
		Err(err)
		return
	} else {
		for rows.Next() {
			var temp ShopSale
			if err = rows.Scan(&temp.ID, &temp.Username, &temp.Pxc, &temp.Gold, &temp.Date, &temp.SaleCount); err != nil {
				Err(err)
				return
			}
			if temp.Date == all {
				allGoldSale = temp.Gold
				allPxcSale = temp.Pxc
				allSaleClothes = temp.SaleCount
			} else if temp.Date == today {
				todayGoldSale = temp.Gold
				todayPxcSale = temp.Pxc
				todaySaleClothes = temp.SaleCount
			} else if temp.Date == yesterday {
				yesterdayGoldSale = temp.Gold
				yesterdayPxcSale = temp.Pxc
				yesterdaySaleClothes = temp.SaleCount
			}
		}
	}
	return
}

//获取在售图纸
func ListSalePaper(offset, pageSize, sortType int, qUsername string) (err error, list []SalePaper) {
	var rows *sql.Rows

	if sortType == 0 {
		rows, err = dao_pixie.ListPaperTradeOrderByStartTimeStmt.Query(qUsername, qUsername, offset, pageSize)
	} else if sortType == 1 {
		rows, err = dao_pixie.ListPaperTradeOrderByLeftTimeStmt.Query(qUsername, qUsername, offset, pageSize)
	} else {
		err = constants.UnknownSortTypeErr
		return
	}

	if err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var temp SalePaper
			if err = rows.Scan(&temp.TradeID, &temp.AuthorUsername, &temp.AuthorNickname, &temp.AuthorHead, &temp.AuthorSex, &temp.OwnerUsername, &temp.OwnerNickname, &temp.OwnerHead, &temp.OwnerSex, &temp.PaperID, &temp.ClothesType, &temp.PartType, &temp.Cname, &temp.Desc, &temp.Extra, &temp.File, &temp.Style, &temp.Star, &temp.Tag1, &temp.Tag2, &temp.STag, &temp.MaxPrice, &temp.MinPrice, &temp.StartTime, &temp.Duration, &temp.PriceType, &temp.AuctionTime, &temp.AuctionUsername, &temp.AuctionNickname, &temp.AuctionPrice, &temp.Sequence, &temp.CirculationMax); err != nil {
				Err(err)
				return
			}
			temp.Status = int(PAPER_STATUS_ON_SALE)
			list = append(list, temp)
		}
	}
	return
}

func DeletePaper(paperID int64, au string) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_pixie.DeletePaperByIDStmt.Exec(PAPER_STATUS_DELETED, paperID, au); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		success = true
	}
	return
}

func UpdateAuctionFailUnreadStatus(ownerUsername string, paperID int64, status bool, typee int) (err error) {
	var res sql.Result
	var effect int64
	if typee == 1 {
		if res, err = dao_pixie.UpdatePaperAuctionFailStatusStmt.Exec(status, paperID); err != nil {
			Err(err)
			return
		} else if effect, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		} else if effect < 1 {
			err = constants.UpdatePaperAuctionFailUnreadZero
			return
		}
	} else {
		if res, err = dao_pixie.UpdatePaperOccupyAuctionFailStatusStmt.Exec(status, paperID, ownerUsername); err != nil {
			Err(err)
			return
		} else if effect, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		} else if effect < 1 {
			err = constants.UpdatePaperOccupyAuctionFailUnreadZero
			return
		}
	}
	return
}

func GetPaperOwnerUsername(paperID int64) (err error, ownername string) {
	if err = dao_pixie.GetPaperOwnerUsernameStmt.QueryRow(paperID).Scan(&ownername); err != nil {
		Err(err)
		return
	}
	return
}

func GetPaperAuthorUsername(paperID int64) (err error, username string) {
	if err = dao_pixie.GetPaperAuthorUsernameStmt.QueryRow(paperID).Scan(&username); err != nil {
		Err(err)
		return
	}
	return
}

func QueryPaperTradeHistory(paperID int64, sequence, offset, pageSize int) (err error, res []PaperTradeHistory, ts int64) {
	var rows *sql.Rows
	if rows, err = dao_pixie.GetTradeHistoryStmt.Query(paperID, sequence, offset, pageSize); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var temp PaperTradeHistory
			var id int64
			if err = rows.Scan(&id, &temp.PaperID, &temp.Username, &temp.Nickname, &temp.TimeBuy, &temp.TradePriceType, &temp.TradePrice, &temp.SaleCount); err != nil {
				Err(err)
				return
			}
			res = append(res, temp)
		}
	}

	if err = dao_pixie.GetPaperTotalSaleBySeqStmt.QueryRow(paperID, sequence).Scan(&ts); err != nil {
		Err(err)
		return
	}
	return
}

func UpdateOccupyAfterSale(paperID int64, logID int64, saleNum int, username string) (err error) {
	var tx *sql.Tx
	var res sql.Result
	var effect int64

	if tx, err = dao_pixie.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	//更新paperOccupy 加销量,减库存
	if res, err = tx.Stmt(dao_pixie.UpdateOccupySaleStmt).Exec(saleNum, saleNum, username, paperID); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect < 1 {
		err = constants.PaperTradeLogSaleCountEffectZero
		return
	}

	//更新tradelog表的库存
	if res, err = tx.Stmt(dao_pixie.UpdateOccupyLogSaleStmt).Exec(saleNum, logID); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect < 1 {
		err = constants.PaperTradeLogSaleCountEffectZero
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func ClothesPricing(price int64, priceType int, paperID int64, ownerUsername string) (err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_pixie.ClothesPricingStmt.Exec(price, priceType, paperID, ownerUsername, priceType); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect < 1 {
		err = constants.ClothesPricingEffectZero
		return
	}
	return
}

func ClothesPricingChangePrice(price int64, paperID int64, ownerUsername string, priceType int) (err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_pixie.ClothesPriceChangeStmt.Exec(price, paperID, ownerUsername, priceType); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect < 1 {
		err = constants.ClothesPriceChangeEffectZero
		return
	}
	return
}

func GetDesignPaperByID(paperID int64) (err error, temp DesignPaper) {
	if err = dao_pixie.GetOneDesignPaperStmt.QueryRow(paperID).Scan(&temp.PaperID, &temp.AuthorUsername, &temp.AuthorNickname, &temp.AuthorHead, &temp.AuthorSex, &temp.OwnerUsername, &temp.OwnerNickname, &temp.OwnerHead, &temp.OwnerSex, &temp.ClothesType, &temp.PartType, &temp.Cname, &temp.Desc, &temp.UploadTime, &temp.Extra, &temp.File, &temp.Style, &temp.Status, &temp.AuctionFailUnread, &temp.Star, &temp.Tag1, &temp.Tag2, &temp.STag, &temp.PassTime, &temp.VerifiedCount, &temp.PriceType, &temp.LastDuration, &temp.LastMaxPrice, &temp.LastMinPrice, &temp.CopyMark, &temp.Reason, &temp.CirculationMax, &temp.Circulation); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("can't find paper by id", paperID)
		} else {
			Err(err)
		}
	}
	return
}

func GetDesignPaperByStatus() (err error, designList []DesignPaper) {
	var rows *sql.Rows
	if rows, err = dao_pixie.GetDesignPaperByStatusStmt.Query(PAPER_STATUS_QUEUE); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		return getDesignPapers(rows, false, -999)
	}
	return
}

func GetDesignPaperInIDList(idList []interface{}) (err error, designList []DesignPaper) {
	for len(idList) < constants.PIXIE_PAPER_DESIGN_GET_SIZE {
		idList = append(idList, 0)
	}

	var rows *sql.Rows
	if rows, err = dao_pixie.GetDesignPaperInIDListStmt.Query(idList...); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		return getDesignPapers(rows, false, -999)

	}
}

func GetRecommendDesignPaperInIDList(idList []interface{}) (err error, recommendPaperList []RecommendClothes) {
	for len(idList) < constants.PIXIE_RECOMMEND_PAPER_DESIGN_GET_SIZE {
		idList = append(idList, 0)
	}
	var rows *sql.Rows
	if rows, err = dao_pixie.GetRecommendDesignPaperInIDListStmt.Query(idList...); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		return getRecommendPapers(rows)
	}
}
