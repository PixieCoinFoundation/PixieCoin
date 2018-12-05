package board

import (
	"dao_friend"
	"database/sql"
	"fmt"
	. "logger"
	"time"
	. "types"
)

func AddBoardMsg(owner, author, authorNickname, head, reply_to, content string) (err error) {
	n := time.Now().Unix()

	// var id int64
	// var res sql.Result
	if _, err = dao_friend.AddBoardMsgStmt.Exec(owner, author, authorNickname, head, reply_to, content, n); err != nil {
		Err(err)
	} else {
		// AddBoardMsgR(&BoardMsg{ID: id, Owner: owner, Author: author, AuthorNickname: authorNickname, Head: head, ReplyTo: reply_to, Content: content, Time: n})
		if owner != author {
			SetBoardNew(owner, true)
		}
	}
	return
}

func DeleteBoardMsgByPlayer(id int64, player string) (err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_friend.DeleteBoardMsgByPlayerStmt.Exec(id, player, player); err != nil {
		Err(err)
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
	} else if effect <= 0 {
		Info("del board effect 0", id, player)
		err = fmt.Errorf("del board effect 0 %d %s", id, player)
	} else {
		//del in redis
		// DelBoardMsgR(id, player)
	}
	return
}

func GetBoardMsgs(owner string, pageNo, pageCnt int) (res []BoardMsg, err error) {
	// if pageNo == 1 {
	// 	if res, err = GetBoardFirstPageR(owner); err == nil {
	// 		return
	// 	}
	// }

	return getBoardMsgsDB(owner, pageNo, pageCnt)
}
