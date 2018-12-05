package board

import (
	"dao_friend"
	"database/sql"
	. "logger"
	. "types"
)

func getBoardMsgsDB(owner string, pageNo, pageCnt int) (res []BoardMsg, err error) {
	start := (pageNo - 1) * pageCnt

	var rows *sql.Rows
	if rows, err = dao_friend.GetBoardMsgStmt.Query(owner, start, pageCnt); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]BoardMsg, 0)
		for rows.Next() {
			m := BoardMsg{}
			if err = rows.Scan(&m.ID, &m.Owner, &m.Author, &m.AuthorNickname, &m.Head, &m.ReplyTo, &m.Content, &m.Time); err != nil {
				Err(err)
				return
			}
			res = append(res, m)
		}
	}
	return
}

func BatchQueryBoardMsgs(start, size int) (res []BoardMsg, err error) {
	var rows *sql.Rows
	if rows, err = dao_friend.BatchQueryBoardStmt.Query(start, size); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]BoardMsg, 0)
		for rows.Next() {
			m := BoardMsg{}
			if err = rows.Scan(&m.ID, &m.Time); err != nil {
				Err(err)
				return
			}
			res = append(res, m)
		}
	}
	return
}

func BatchDelBoardMsgs(id int64) (err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_friend.BatchDeleteBoardStmt.Exec(id); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else {
		Info("delete board msgs before id:", id, "effect:", effect)
	}
	return
}
