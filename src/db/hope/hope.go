package hope

import (
	"constants"
	"dao_friend"
	"database/sql"
	. "logger"
	"time"
	. "types"
)

func GetOneHope(id int64) (h Hope, err error) {
	if err = dao_friend.GetOneHopeStmt.QueryRow(id).Scan(&h.ID, &h.Sender, &h.SenderNickname, &h.OfferClothes, &h.OfferPart, &h.OfferNum, &h.NeedClothes, &h.NeedPart, &h.NeedNum, &h.SendTime, &h.Status, &h.Helper, &h.HelperNickname); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("no hope with id:", id)
		} else {
			Err(err)
		}
	}
	return
}

func SendHope(sender string, senderNickname string, offerClothes string, offerPart, offerNum int, needClothes string, needPart, needNum int) (id int64, err error) {
	var res sql.Result
	if res, err = dao_friend.SendHopeStmt.Exec(sender, senderNickname, offerClothes, offerPart, offerNum, needClothes, needPart, needNum, time.Now().Unix(), constants.HOPE_STATUS_SENDED, "", ""); err != nil {
		Err(err)
		return
	} else if id, err = res.LastInsertId(); err != nil {
		Err(err)
		return
	}
	return
}

func CancelHope(hopeID int64, sender string) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_friend.CancelHopeStmt.Exec(hopeID, sender, constants.HOPE_STATUS_SENDED); err != nil {
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

func DelHope(id int64) {
	if _, err := dao_friend.DelHopeStmt.Exec(id); err != nil {
		Err(err)
	}
	return
}

func GetPlayerHopes(sender string, lessThanStatus int) (res []Hope, err error) {
	var rows *sql.Rows
	if rows, err = dao_friend.QueryPlayerHopeStmt.Query(sender, lessThanStatus); err != nil {
		Err(err)
	} else {
		defer rows.Close()
		res = make([]Hope, 0)
		for rows.Next() {
			var r Hope
			if err = rows.Scan(&r.ID, &r.Sender, &r.SenderNickname, &r.OfferClothes, &r.OfferPart, &r.OfferNum, &r.NeedClothes, &r.NeedPart, &r.NeedNum, &r.SendTime, &r.Status, &r.Helper, &r.HelperNickname); err != nil {
				Err(err)
				return
			}

			res = append(res, r)
		}
	}
	return
}

func SelectHope(size int) (res []Hope, err error) {
	var rows *sql.Rows
	if rows, err = dao_friend.SelectHopeStmt.Query(size); err != nil {
		Err(err)
	} else {
		defer rows.Close()
		res = make([]Hope, 0)
		for rows.Next() {
			var r Hope
			if err = rows.Scan(&r.ID, &r.Sender, &r.SenderNickname, &r.OfferClothes, &r.OfferPart, &r.OfferNum, &r.NeedClothes, &r.NeedPart, &r.NeedNum, &r.SendTime, &r.Status, &r.Helper, &r.HelperNickname); err != nil {
				Err(err)
				return
			}

			res = append(res, r)
		}
	}
	return
}

func UpdateHopeStatus(id int64, status int, oriStatus int) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_friend.UpdateHopeStatusStmt.Exec(status, id, oriStatus); err != nil {
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

func UpdateHope(id int64, status int, oriStatus int, helper, helperNickname string) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_friend.UpdateHopeStmt.Exec(status, helper, helperNickname, id, oriStatus); err != nil {
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
