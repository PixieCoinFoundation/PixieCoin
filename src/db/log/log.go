package log

import (
	"appcfg"
	"constants"
	"dao"
	"database/sql"
	"encoding/json"
	"fmt"
	. "model"
	"tools"
	. "types"
)

func DelKorLog(id int64) (err error) {
	if _, err = dao.DelKorLogStmt.Exec(id); err != nil {
		fmt.Println(err)
		return
	}
	return
}

func ListKorLog() (res []KorLogClean, err error) {
	var rows *sql.Rows
	if rows, err = dao.ListKorLogStmt.Query(); err != nil {
		fmt.Println(err)
		return
	} else {
		defer rows.Close()
		res = make([]KorLogClean, 0)
		for rows.Next() {
			var r KorLogClean

			if err = rows.Scan(&r.ID, &r.Time); err != nil {
				fmt.Println(err)
				return
			}

			res = append(res, r)
		}
	}
	return
}

func AddLog(C1 string, C2 string, C3 string, username string, t int64, extra string) {
	var stmt *sql.Stmt
	if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
		stmt = dao.AddKorLogStmt
	} else {
		if tools.GetDayTag() == 0 {
			stmt = dao.AddLog1Stmt
		} else {
			stmt = dao.AddLog2Stmt
		}
	}

	if _, err := stmt.Exec(C1, C2, C3, username, t, extra); err != nil {
		fmt.Println("add log error", err)
	}
}

func AddLog10(logs []BufferLog) {
	if len(logs) != 10 {
		// Err("add log not 10", logs)
		b, _ := json.Marshal(logs)
		tools.SendInternalMail("log not 10", "log not 10:"+string(b), nil)
		return
	}

	var stmt *sql.Stmt
	if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
		stmt = dao.AddKorLogStmt10
	} else {
		if tools.GetDayTag() == 0 {
			stmt = dao.AddLog1Stmt10
		} else {
			stmt = dao.AddLog2Stmt10
		}
	}

	if _, err := stmt.Exec(
		logs[0].C1, logs[0].C2, logs[0].C3, logs[0].Username, logs[0].T64, logs[0].Extra,
		logs[1].C1, logs[1].C2, logs[1].C3, logs[1].Username, logs[1].T64, logs[1].Extra,
		logs[2].C1, logs[2].C2, logs[2].C3, logs[2].Username, logs[2].T64, logs[2].Extra,
		logs[3].C1, logs[3].C2, logs[3].C3, logs[3].Username, logs[3].T64, logs[3].Extra,
		logs[4].C1, logs[4].C2, logs[4].C3, logs[4].Username, logs[4].T64, logs[4].Extra,
		logs[5].C1, logs[5].C2, logs[5].C3, logs[5].Username, logs[5].T64, logs[5].Extra,
		logs[6].C1, logs[6].C2, logs[6].C3, logs[6].Username, logs[6].T64, logs[6].Extra,
		logs[7].C1, logs[7].C2, logs[7].C3, logs[7].Username, logs[7].T64, logs[7].Extra,
		logs[8].C1, logs[8].C2, logs[8].C3, logs[8].Username, logs[8].T64, logs[8].Extra,
		logs[9].C1, logs[9].C2, logs[9].C3, logs[9].Username, logs[9].T64, logs[9].Extra); err != nil {
		fmt.Println("add log error", err)
	}
}
