package rt_party

import (
	"dao_party"
	"database/sql"
	. "logger"
	"time"
)

func createHostDB(username string, t int64) (hostID int64, err error) {
	ds := time.Unix(t, 0).Format("20060102")
	var res sql.Result
	if res, err = dao_party.AddRTPartyHostStmt.Exec(ds, username, 1, t); err != nil {
		Err(err)
		return
	} else if hostID, err = res.LastInsertId(); err != nil {
		Err(err)
		return
	}
	return
}

func ListHostDB(ds string) (res []int64, err error) {
	res = make([]int64, 0)

	var rows *sql.Rows
	if rows, err = dao_party.ListRTPartyHostStmt.Query(ds); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var id int64
			if err = rows.Scan(&id); err != nil {
				Err(err)
				return
			}

			res = append(res, id)
		}
	}
	return
}

func DelHostByDs(ds string) {
	if _, err := dao_party.BatchDelRTPartyHostStmt.Exec(ds); err != nil {
		Err(err)
		return
	}
}
