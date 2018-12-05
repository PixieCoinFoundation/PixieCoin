package paper

import (
	"constants"
	"dao_pixie"
	"database/sql"
	"encoding/json"
	. "logger"
	. "pixie_contract/api_specification"
	"service_pixie/paperInfo"
	"time"
)

func ListPaperVerifyWithDetail(paperID int64) (err error, list []PaperVerify, scoreCount, verifySize int, tagMap map[string]int, styleMap map[string]int) {
	var rows *sql.Rows

	if rows, err = dao_pixie.GetPaperVerifyByPaperIDStmt.Query(paperID, false); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		list = make([]PaperVerify, 0)
		tagMap = make(map[string]int)
		styleMap = make(map[string]int)
		for rows.Next() {
			var paperVerify PaperVerify
			if err = rows.Scan(&paperVerify.ID, &paperVerify.PaperID, &paperVerify.Cname, &paperVerify.Extra, &paperVerify.Username, &paperVerify.Nickname, &paperVerify.CheckDate, &paperVerify.Tag1, &paperVerify.Tag2, &paperVerify.Style1, &paperVerify.Style2, &paperVerify.Score, &paperVerify.MatchDegree, &paperVerify.RewardPxc, &paperVerify.SubmitTime, &paperVerify.CheckoutTime, &paperVerify.Status, &paperVerify.Result); err != nil {
				Err(err)
				return
			}

			list = append(list, paperVerify)
			scoreCount += paperVerify.Score
			verifySize++

			if paperInfo.TagLegal(paperVerify.Tag1) {
				tagMap[paperVerify.Tag1] += 1
			}

			if paperInfo.TagLegal(paperVerify.Tag2) {
				tagMap[paperVerify.Tag2] += 1
			}

			if paperInfo.StyleLegal(paperVerify.Style1) {
				styleMap[paperVerify.Style1] += 1
			}

			if paperInfo.StyleLegal(paperVerify.Style2) {
				styleMap[paperVerify.Style2] += 1
			}
		}
	}

	return
}

func PaperVerifyInsert(paperID int64, cname, paper_extra, username, nickname, tag1, tag2, style1, style2 string, score int, submitTime int64) (err error) {
	var tx *sql.Tx
	if tx, err = dao_pixie.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	if _, err = tx.Stmt(dao_pixie.VerifyPaperInsertStmt).Exec(paperID, cname, paper_extra, username, nickname, tag1, tag2, style1, style2, score, submitTime); err != nil {
		Err(err)
		return
	}

	if _, err = tx.Stmt(dao_pixie.UpdateVerifyCountStmt).Exec(1, paperID, PAPER_STATUS_QUEUE); err != nil {
		Err(err)
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func PaperNotifyQuery(username string, offset int, pageSize int, tNow time.Time) (err error, yesterdayProfit float64, allProfit float64, list []PaperNotify) {
	var rows *sql.Rows
	if rows, err = dao_pixie.GetPaperNotifyStmt.Query(username, offset, pageSize); err != nil {
		Err(err)
		return
	} else {
		once := true
		defer rows.Close()
		for rows.Next() {
			var temp PaperNotify
			if err = rows.Scan(&temp.ID, &temp.Username, &temp.Content, &temp.DayPxcProfit, &temp.Cname, &temp.Time, &temp.Status, &temp.Type); err != nil {
				Err(err)
				return
			}
			if once {
				if time.Unix(temp.Time, 0).Format("20060102") == tNow.Format("20060102") {
					yesterdayProfit = temp.DayPxcProfit
					once = false
				}
			}
			list = append(list, temp)
		}
	}
	if err, allProfit = getAllProfitByUsername(username); err != nil {
		Err(err)
		return
	}
	return
}

func getAllProfitByUsername(username string) (err error, allProfit float64) {
	if err = dao_pixie.GetAllProfitByUsernameStmt.QueryRow(username).Scan(&allProfit); err != nil {
		Err(err)
		return
	}
	return
}

func PaperVerifiedDetailQuery(username string, date string) (err error, list []PaperVerify) {
	var rows *sql.Rows
	if rows, err = dao_pixie.PaperVerifiedDetailStmt.Query(username, date); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var paperVerify PaperVerify
			if err = rows.Scan(&paperVerify.ID, &paperVerify.PaperID, &paperVerify.Cname, &paperVerify.Extra, &paperVerify.Username, &paperVerify.Nickname, &paperVerify.CheckDate, &paperVerify.Tag1, &paperVerify.Tag2, &paperVerify.Style1, &paperVerify.Style2, &paperVerify.Score, &paperVerify.MatchDegree, &paperVerify.RewardPxc, &paperVerify.SubmitTime, &paperVerify.CheckoutTime, &paperVerify.Status, &paperVerify.Result); err != nil {
				Err(err)
				return
			}
			list = append(list, paperVerify)
		}
	}
	return
}

func InsertNotify(currencyType int, currencyVal float64, verifyResultNum, dayVerifyCount int, username string, cname string, tn time.Time, notifyType PaperNotifyType) (err error) {
	var content RewardContent
	content.CurrencyType = currencyType
	content.CurrencyVal = currencyVal
	content.VerifyResultNum = verifyResultNum
	content.DayVerifyCount = dayVerifyCount

	data, _ := json.Marshal(content)
	if _, err = dao_pixie.NotifyInsertStmt.Exec(username, string(data), currencyVal, cname, tn.Unix(), int(notifyType)); err != nil {
		Err(err, currencyType, currencyVal, verifyResultNum, dayVerifyCount, username, cname, tn.Unix(), notifyType)
		return
	}
	return
}

func UpdateDesignPaperAfterVerify(toStatus PaperStatus, star int, tag1, tag2, style string, tNow time.Time, paperID int64, originStatus PaperStatus) (err error) {
	var num int64
	var res sql.Result
	if res, err = dao_pixie.PassPaperStmt.Exec(toStatus, star, tag1, tag2, style, tNow.Unix(), paperID, originStatus); err != nil {
		Err(err)
		return
	} else if num, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if num < 1 {
		Err("updatePaperVerifiedInfo effect 0", paperID, star)
		err = constants.UpdatePaperVerifiedInfoEffectZero
		return
	}
	return
}

func GetVerifiedPaperByUsername(username string) (err error, idMap map[int64]bool) {
	var rows *sql.Rows

	if rows, err = dao_pixie.GetVerifiedPaperByUsernameStmt.Query(username); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		idMap = make(map[int64]bool)
		for rows.Next() {
			var paperID int64
			if err = rows.Scan(&paperID); err != nil {
				Err(err)
				return
			}
			idMap[paperID] = true
		}
	}
	return
}

func PaperReportCopy(username string, paperID int64, reason string, pic string, contact string) (err error) {
	var tx *sql.Tx
	if tx, err = dao_pixie.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()
	if _, err = tx.Stmt(dao_pixie.PaperReportCopyStmt).Exec(username, paperID, reason, pic, contact); err != nil {
		Err(err)
		return
	}

	if _, err = tx.Stmt(dao_pixie.UpdatePaperDesignCopyInfoStmt).Exec(1, true, paperID); err != nil {
		Err(err)
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func GetCopyByIDUsername(username string, paperID int64) (err error, exist bool) {
	var count int64
	if err = dao_pixie.GetCopyByIDUsernameStmt.QueryRow(username, paperID).Scan(&count); err != nil {
		Err(err)
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

func PaperCopyQuery(paperID int64) (err error, list []PaperCopy) {
	var rows *sql.Rows
	if rows, err = dao_pixie.GetPaperCopyByPaperIDStmt.Query(paperID); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var temp PaperCopy
			if err = rows.Scan(&temp.CopyID, &temp.Username, &temp.PaperID, &temp.Reason, &temp.Pic, &temp.Contact, &temp.SupportNum, &temp.RejectNum, &temp.Time, &temp.CheckTime, &temp.Status); err != nil {
				Err(err)
				return
			}
			list = append(list, temp)
		}
	}
	return
}

func UpdatePaperVerifyAfterSettle(verifyID int64, date, verifyResult string, matchDegree float64, nt int64) (err error) {
	if _, err = dao_pixie.VerifiedUpdateInfoStmt.Exec(date, matchDegree, nt, verifyResult, verifyID, false); err != nil {
		Err(err, verifyID, date, verifyResult, matchDegree, nt)
	}
	return
}

func PaperCopySupprot(copyID int64, paperID int64) (err error) {
	var tx *sql.Tx
	if tx, err = dao_pixie.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()
	var num int64
	var res sql.Result
	if res, err = tx.Stmt(dao_pixie.UpdatePaperSupportStmt).Exec(1, copyID); err != nil {
		Err(err)
		return
	} else if num, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if num < 1 {
		Err("UpdatePaperSupportStmt effect 0", copyID, paperID)
		err = constants.UpdatePaperSupportEffectZero
		return
	}

	if res, err = tx.Stmt(dao_pixie.UpdateVerifyCountStmt).Exec(1, paperID, PAPER_STATUS_QUEUE); err != nil {
		Err(err)
		return
	} else if num, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if num < 1 {
		Err("UpdateVerifyCountStmt effect 0", copyID, paperID)
		err = constants.UpdateVerifyCountEffectZero
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func GetPaperNotifyByID(id int64) (err error, temp PaperNotify) {
	if err = dao_pixie.GetPaperNotifyByIDStmt.QueryRow(id).Scan(&temp.ID, &temp.Username, &temp.Content, &temp.DayPxcProfit, &temp.Cname, &temp.Time, &temp.Status, &temp.Type); err != nil {
		Err(err)
		return
	}
	return
}

func UpdatePaperNotifyRewarded(id int64) (err error) {
	if _, err = dao_pixie.UpdatePaperNotifyRewardStatusStmt.Exec(true, id); err != nil {
		Err(err)
		return
	}
	return
}
