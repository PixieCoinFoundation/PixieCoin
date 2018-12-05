package push

import (
	"constants"
	"dao"
	"dao_designer"
	"database/sql"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	. "logger"
	"rao"
	"strings"
	"time"
	. "types"
)

func PopDayNewCustom(ds string) (po NewCustomPushObject, err error) {
	//new custom push
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetDayNewCustomListKey(ds)
	var pos string
	if pos, err = redis.String(conn.Do("LPOP", key)); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		}
	} else if len(pos) > 2 {
		pos = pos[1 : len(pos)-1]
		if err = json.Unmarshal([]byte(pos), &po); err != nil {
			Err(err)
			return
		}
	}

	return
}

func AddDesignerPushHistory(designerUsername string) bool {
	token := time.Now().Format("20060102")
	if _, err := dao_designer.AddDesignPushHistoryStmt.Exec(designerUsername, token); err != nil {
		if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
			Info("push already send:", designerUsername, token)
		} else {
			Err(err)
		}
		return false
	}

	return true
}

func GetLegalPushJobs() (res []Push, err error) {
	var rows *sql.Rows
	if rows, err = dao.GetPushStmt.Query(constants.PUSH_STATUS_LINE); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]Push, 0)
		for rows.Next() {
			e := Push{}
			if err = rows.Scan(&e.ID, &e.Title, &e.Content, &e.PushTime, &e.To, &e.Status); err != nil {
				Err(err)
				return
			}

			// e.PushTimes = time.Unix(e.PushTime, 0).Format("2006-01-02 15:04:05")

			// if e.Status == constants.PUSH_STATUS_DONE {
			// 	closed = append(closed, e)
			// } else {
			// 	line = append(line, e)
			// }
			res = append(res, e)
		}
	}

	return
}

func DonePush(id int64) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = dao.UpdatePushStmt.Exec(constants.PUSH_STATUS_DONE, id, constants.PUSH_STATUS_LINE); err != nil {
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
