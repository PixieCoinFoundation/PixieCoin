package board

import (
	// "constants"
	// "encoding/json"
	"github.com/garyburd/redigo/redis"
	. "logger"
	"rao"
	// . "types"
)

func SetBoardNew(username string, fresh bool) {
	key := rao.GetBoardNewSignKey(username)

	v := 1
	if !fresh {
		v = 0
	}

	conn := rao.GetConn()
	defer conn.Close()

	if fresh {
		conn.Send("SET", key, v)
		conn.Send("EXPIRE", key, 24*3600)
	} else {
		conn.Send("DEL", key)
	}

	if err := conn.Flush(); err != nil {
		Err(err)
	}
}

func GetBoardNew(username string) bool {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetBoardNewSignKey(username)
	if v, err := redis.Int(conn.Do("GET", key)); err != nil {
		// Err(err)
		return false
	} else if v > 0 {
		return true
	}
	return false
}

// func GetBoardFirstPageR(username string) (res []BoardMsg, err error) {
// 	key := rao.GetBoardFirstPageKey(username)

// 	conn := rao.GetConn()
// 	defer conn.Close()

// 	var exist bool
// 	if exist, err = rao.ExistKey(key, &conn, constants.DEFAULT_EXPIRE_TIME); err != nil {
// 		return
// 	} else if exist {
// 		var values []interface{}
// 		if values, err = redis.Values(conn.Do("ZREVRANGE", key, 0, -1)); err != nil {
// 			Err(err)
// 			return
// 		}
// 		res = make([]BoardMsg, 0)
// 		for len(values) > 0 {
// 			var bs string
// 			if values, err = redis.Scan(values, &bs); err != nil {
// 				Err(err)
// 				return
// 			}
// 			bs = bs[1 : len(bs)-1]

// 			var m BoardMsg
// 			if err = json.Unmarshal([]byte(bs), &m); err != nil {
// 				Err(err)
// 				return
// 			}

// 			res = append(res, m)
// 		}
// 	} else {
// 		if res, err = getBoardMsgsDB(username, 1, constants.BOARD_PAGE_SIZE); err != nil {
// 			Err(err)
// 			return
// 		}
// 		key := rao.GetBoardFirstPageKey(username)
// 		for _, m := range res {
// 			mb, _ := json.Marshal(m)
// 			conn.Send("ZADD", key, m.ID, "'"+string(mb)+"'")
// 		}
// 		conn.Send("ZREMRANGEBYRANK", key, 0, -(constants.BOARD_PAGE_SIZE + 1))

// 		if err = conn.Flush(); err != nil {
// 			Err(err)
// 			return
// 		}
// 	}
// 	return
// }

// func AddBoardMsgR(m *BoardMsg) (err error) {
// 	//only add when exist
// 	conn := rao.GetConn()
// 	defer conn.Close()

// 	key := rao.GetBoardFirstPageKey(m.Owner)
// 	nkey := rao.GetBoardNewSignKey(m.Owner)
// 	var exist bool
// 	if exist, err = rao.ExistKey(key, &conn, constants.DEFAULT_EXPIRE_TIME); err != nil {
// 		Err(err)
// 	} else if exist {
// 		mb, _ := json.Marshal(m)

// 		conn.Send("ZADD", key, m.ID, "'"+string(mb)+"'")
// 		conn.Send("ZREMRANGEBYRANK", key, 0, -(constants.BOARD_PAGE_SIZE + 1))

// 		if m.Owner != m.Author {
// 			conn.Send("SET", nkey, 1)
// 			conn.Send("EXPIRE", nkey, 24*3600)
// 		}

// 		if err = conn.Flush(); err != nil {
// 			Err(err)
// 		}
// 	}

// 	return
// }

// func DelBoardMsgR(id int64, owner string) (err error) {
// 	conn := rao.GetConn()
// 	defer conn.Close()

// 	key := rao.GetBoardFirstPageKey(owner)

// 	if _, err = conn.Do("ZREMRANGEBYSCORE", key, id, id); err != nil {
// 		Err(err)
// 	}
// 	return
// }
