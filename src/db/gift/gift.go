package gift

import (
	"constants"
	"dao"
	"database/sql"
	"encoding/json"
	. "logger"
	"strings"
	"time"
	. "types"
)

func GetGiftPackIDByCode(code string) (id int64, typee int) {
	if err := dao.GetGiftPackIDStmt.QueryRow(code).Scan(&id); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("code pack id type single not exist:", code)
			if err := dao.GetGiftPackID2Stmt.QueryRow(code).Scan(&id); err != nil {
				if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
					Info("code pack id type multi not exist:", code)
				} else {
					Err(err)
					return
				}
			} else {
				typee = constants.MULTI_TYPE
			}
		} else {
			Err(err)
			return
		}
	} else {
		typee = constants.SINGLE_TYPE
	}
	return
}

func UseCode(code string, typee int, packID int64, username string) bool {
	var res sql.Result
	var err error

	if typee == constants.SINGLE_TYPE {
		res, err = dao.UpdateGiftCodeStmt.Exec(username, code)
	} else if typee == constants.MULTI_TYPE {
		res, err = dao.UpdateGiftCode2Stmt.Exec(username, code)
	} else {
		res, err = dao.UseOneCodeStmt.Exec(packID, username)
	}

	if err != nil {
		if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
			Info("already got reward:", username, packID)
		} else {
			Err(err)
		}

		return false
	} else {
		if effect, err := res.RowsAffected(); err != nil {
			Err(err)
			return false
		} else if effect > 0 {
			return true
		}
	}
	return false
}

func GetGictPacks() (res map[int64]*GiftPackContent, ns map[int64]string, err error) {
	var rows *sql.Rows
	if rows, err = dao.GetGiftPackStmt.Query(); err != nil {
		Err(err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*GiftPackContent)
	ns = make(map[int64]string)
	for rows.Next() {
		var id int64
		var name, content string
		if err = rows.Scan(&id, &name, &content); err != nil {
			Err(err)
			return
		}

		var gp GiftPackContent
		if err = json.Unmarshal([]byte(content), &gp); err != nil {
			Err(err)
			return
		}
		if gp.StartTime <= time.Now().Unix() && gp.EndTime >= time.Now().Unix() {
			res[id] = &gp
			ns[id] = name
		}
	}

	return
}
