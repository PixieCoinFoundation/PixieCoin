package cosplay

import (
	// "common"
	// "constants"
	"database/sql"
	// "encoding/json"
	"github.com/garyburd/redigo/redis"
)

import (
	"dao_party"
	. "logger"
	// "rao"
	. "types"
)

func DeleteCosComment(cosplayID int64) (err error) {
	if _, err = dao_party.DeleteCosCommentStmt.Exec(cosplayID); err != nil {
		Err(err)
	}
	return
}

func addCommentInDB(c *CosComment) (err error) {
	if _, err = dao_party.AddCosCommentStmt.Exec(c.CosplayID, c.CosItemID, c.Username, c.Nickname, c.Content, c.ReplyTime); err != nil {
		Err(err)
		return
	}

	return
}

//rconn not nil means recover mode
func AddComment(comment CosComment, rconn *redis.Conn, addToHead bool) (err error) {
	// var cj []byte
	// if cj, err = json.Marshal(comment); err != nil {
	// 	Err(err)
	// 	return
	// }

	// key := rao.GetCosItemCmtListKey(comment.CosItemID)
	// hasKey := rao.GetCosItemHasCmtKey(comment.CosItemID)

	// cmd := "LPUSH"
	// if !addToHead {
	// 	cmd = "RPUSH"
	// }

	// if rconn == nil {
	//ASYNC add to db
	// go cos_item_cmt_batch.AddCosItemCmt(&comment)
	//add to db first
	if err = addCommentInDB(&comment); err != nil {
		return
	}

	//add to redis
	// 	conn := rao.GetConn()
	// 	defer conn.Close()

	// 	conn.Send(cmd, key, "'"+string(cj)+"'")
	// 	conn.Send("SET", hasKey, 1)
	// 	if err = conn.Flush(); err != nil {
	// 		Err(err)
	// 		return
	// 	}
	// 	var size int
	// 	size, err = redis.Int(conn.Receive())

	// 	if err != nil {
	// 		Err(err)
	// 		return err
	// 	}

	// 	if size > constants.COS_ITEM_CMT_LIST_SIZE*2 {
	// 		go common.TrimRedisList(key, 0, constants.COS_ITEM_CMT_LIST_SIZE-1)
	// 	}
	// } else {
	// 	(*rconn).Send(cmd, key, "'"+string(cj)+"'")
	// }

	return nil
}

func getCosCommentFromDB(cosplayID int64, itemID int64, page int, pageCount int) (comments []*CosComment, err error) {
	comments = make([]*CosComment, 0, pageCount)

	var rows *sql.Rows
	if rows, err = dao_party.GetCosCommentStmt.Query(cosplayID, itemID, pageCount*(page-1), pageCount); err != nil {
		Err(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		comment := CosComment{}
		var id int64
		if err = rows.Scan(&id, &comment.CosplayID, &comment.CosItemID, &comment.Username, &comment.Nickname, &comment.Content, &comment.ReplyTime); err != nil {
			Err(err)
			return
		}

		comments = append(comments, &comment)
	}

	return
}

func GetCosComment(cosplayID int64, itemID int64, page int, pageCount int) (comments []*CosComment, err error) {
	// comments = make([]*CosComment, 0, pageCount)

	// if page*pageCount <= constants.COS_ITEM_CMT_LIST_SIZE {
	// 	//get from redis
	// 	conn := rao.GetConn()
	// 	defer conn.Close()

	// 	key := rao.GetCosItemCmtListKey(itemID)

	// 	var values []interface{}
	// 	values, err = redis.Values(conn.Do("LRANGE", key, pageCount*(page-1), page*pageCount-1))
	// 	if err != nil {
	// 		Err(err)
	// 		return
	// 	}
	// 	for len(values) > 0 {
	// 		var cj string
	// 		values, err = redis.Scan(values, &cj)
	// 		if err != nil {
	// 			Err(err)
	// 			return nil, err
	// 		}

	// 		cj = cj[1 : len(cj)-1]
	// 		com := CosComment{}
	// 		err = json.Unmarshal([]byte(cj), &com)
	// 		if err != nil {
	// 			Err(err)
	// 			return
	// 		}
	// 		comments = append(comments, &com)
	// 	}
	// } else {
	return getCosCommentFromDB(cosplayID, itemID, page, pageCount)
	// }

	// return
}
