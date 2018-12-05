package kor

import (
	"constants"
	"dao"
	// "dao_korlog"
	"database/sql"
	. "logger"
	"strings"
	"time"
	. "types"
)

func AddUnregister(username string, thirdUsername string, delTime int64) (err error) {
	if _, err = dao.AddUnregisterStmt.Exec(username, thirdUsername, delTime, delTime); err != nil {
		Err(err)
		return
	}
	return
}

func DelUnregister(id int64) (err error) {
	if _, err = dao.DelUnregisterStmt.Exec(id); err != nil {
		Err(err)
		return
	}
	return
}

func DelUserUnregister(username string) (err error) {
	if _, err = dao.DelUserUnregisterStmt.Exec(username); err != nil {
		Err(err)
		return
	}
	return
}

//GMKorCode 游戏码日志
func GMKorCode(codeType, code, username string) {
	t := time.Now().Format("2006-01-02 15:04:05")
	logToDBKorCode(codeType, code, username, t)
}

//logToDBKorCode 韩国游戏码日志
func logToDBKorCode(codeType string, code string, username string, t string) {
	if codeType == constants.C3_KOR_CODE1 {
		stmt := dao.AddKorGameCode1LogStmt
		if _, err := stmt.Exec(username, t, code); err != nil {
			Err(err)
			return
		}
	} else if codeType == constants.C3_KOR_CODE2 {
		stmt := dao.AddKorGameCode2LogStmt
		if _, err := stmt.Exec(username, t, code); err != nil {
			Err(err)
			return
		}
	} else if codeType == constants.C3_KOR_CODE3 {
		stmt := dao.AddKorGameCode3LogStmt
		if _, err := stmt.Exec(username, t, code); err != nil {
			Err(err)
			return
		}
	}
}

// //CleanKorLoginLog 清除登录log
// func CleanKorLoginLog(ds string) (bool, int64, error) {
// 	if len(ds) < 1 {
// 		Info("ds can not empty:", ds)
// 		return false, 0, constants.ErrorNilData
// 	}
// 	rows, err := dao.CleanKorLoginlogStmt.Exec(ds)
// 	if err != nil {
// 		Info("can't find data in db:", ds)
// 		return false, 0, err
// 	}
// 	if num, err := rows.RowsAffected(); err != nil {
// 		return false, 0, err
// 	} else if num > 0 {
// 		return false, num, nil
// 	}
// 	return false, 0, constants.ErrorDelect

// }

// //QueryKorLoginlog 查询韩国某小时的日志
// func QueryKorLoginlog(lim int, et int64) (ds string) {
// 	rows, err := dao.QueryKorLoginLogStmt.Query(lim)
// 	if err != nil {
// 		Err(err)
// 		return
// 	} else {
// 		for rows.Next() {
// 			var tid int64
// 			var fbt string
// 			if err := rows.Scan(&tid, &fbt); err != nil {
// 				Err(err)
// 				return
// 			}
// 			if tid < et {
// 				ds = fbt
// 			}
// 		}
// 	}
// 	return
// }

// //CleanKorLogin1Log 清除login log
// func CleanKorLogin1Log(ds string) (bool, int64, error) {
// 	if len(ds) < 1 {
// 		Info("ds can not empty:", ds)
// 		return false, 0, constants.ErrorNilData
// 	}
// 	rows, err := dao.CleanKorLogin1logStmt.Exec(ds)
// 	if err != nil {
// 		Info("can't find data in db:", ds)
// 		return false, 0, err
// 	}
// 	if num, err := rows.RowsAffected(); err != nil {
// 		return false, 0, err
// 	} else if num > 0 {
// 		return false, num, nil
// 	}
// 	return false, 0, constants.ErrorDelect
// }

// //QueryKorLogin1log 查询韩国某小时的日志
// func QueryKorLogin1log(lim int, et int64) (ds string) {
// 	rows, err := dao.QueryKorLogin1LogStmt.Query(lim)
// 	if err != nil {
// 		Err(err)
// 		return
// 	} else {
// 		for rows.Next() {
// 			var tid int64
// 			var fbt string
// 			if err := rows.Scan(&tid, &fbt); err != nil {
// 				Err(err)
// 				return
// 			}
// 			if tid < et {
// 				ds = fbt
// 			}
// 		}
// 	}
// 	return
// }

func ListUnregister() (res []KorUnregisterInfo, err error) {
	var rows *sql.Rows

	if rows, err = dao.ListUngisterStmt.Query(); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]KorUnregisterInfo, 0)
		for rows.Next() {
			var r KorUnregisterInfo
			if err = rows.Scan(&r.ID, &r.Username, &r.ThirdUsername, &r.DelTime); err != nil {
				Err(err)
				return
			}

			res = append(res, r)
		}
	}
	return
}

func AddInvite(username string, invite string) (err error) {
	if _, err = dao.AddKorInviteStmt.Exec(username, invite); err != nil {
		if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
			Info("already invite kor:", username, invite)
		} else {
			Err(err)

		}
	}

	return
}

func CountInvite(username string) (l []string, err error) {
	var rows *sql.Rows
	if rows, err = dao.ListKorInviteStmt.Query(username); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		l = make([]string, 0)
		for rows.Next() {
			var i string
			if err = rows.Scan(&i); err != nil {
				Err(err)
				return
			}

			l = append(l, i)
		}
	}

	return
}

// func UpdateInvite(invite string) {
// 	if _, err := dao.UpdateInviteStmt.Exec(invite); err != nil {
// 		Err(err)
// 	}
// }
