package cronjob

import (
	"dao"
	"database/sql"
	. "logger"
)

func LogCronJob(token string) (ok bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = dao.AddCronJobLogStmt.Exec(token); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		ok = true
	}
	return
}

func QueryCronJobProgress(token string) int {
	var v int
	if err := dao.QueryCronJobLogStmt.QueryRow(token).Scan(&v); err != nil {
		Err(err)
	}
	return v
}
