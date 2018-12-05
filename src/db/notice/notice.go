package notice

import (
	"appcfg"
	"common"
	"constants"
	"dao"
	"database/sql"
	"db/custom"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	. "language"
	. "logger"
	"rao"
	"service/mails"
	"shutdown"
	"strings"
	"time"
	. "types"
)

const (
	NOTICE_EXPIRE            = 24 * 3600 //SECONDS 24hours
	NOTICE_CHECK_TIME        = 600       //second
	NEW_FRIEND_NOTICE_EXPIRE = 3600      //1 hour
	NOTICE_LIST_SIZE         = 100
)

func init() {
	//check expired notices in db
	if appcfg.GetBool("main_server", false) && appcfg.GetServerType() == "" {
		go checkExpiredNotice()
	}
}

func checkExpiredNotice() {
	for {
		rejected := checkExpiredNotice1()

		if rejected > 0 {
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(NOTICE_CHECK_TIME * time.Second)
		}

	}
}

func checkExpiredNotice1() (rejected int) {
	shutdown.AddOneRequest("checkExpiredNotice")
	defer shutdown.DoneOneRequest("checkExpiredNotice")

	Info("checking expired notices..")
	now := time.Now().Unix()

	if rows, err := dao.QueryNoticeListStmt.Query(); err == nil {
		for rows.Next() {
			var n Notice
			var toUsername string
			var toUID int64
			var sendTime int64
			var token string
			if err = rows.Scan(&n.NoticeID, &n.FromUsername, &n.FromNickname, &n.FromUID, &toUsername, &toUID, &n.Type, &n.Content, &token, &sendTime); err != nil {
				Err(err)
			} else if now-sendTime >= NOTICE_EXPIRE {
				//expired notice
				rejected++
				RejectNotice(&n, toUsername, toUID, constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_TIMEOUT)
			}
		}
	} else {
		Err(err)
	}

	return
}

func getNoticeListFromDB(to string, start, limit int) (res []*Notice, err error) {
	var rows *sql.Rows
	// if !ignoreFriendResp {
	rows, err = dao.QueryNoticeListByUserStmt.Query(to, start, limit)
	// } else {
	// 	rows, err = dao.QueryNoticeListByUserAndNotTypeStmt.Query(to, constants.NOTICE_TYPE_ADD_FRIEND_RESPONSE, start, limit)
	// }

	if err != nil {
		Err(err)
		return
	}
	defer rows.Close()

	now := time.Now().Unix()
	res = make([]*Notice, 0)
	for rows.Next() {
		var n Notice
		var toUsername string
		var toUID int64
		var token string
		if err = rows.Scan(&n.NoticeID, &n.FromUsername, &n.FromNickname, &n.FromUID, &toUsername, &toUID, &n.Type, &n.Content, &token, &n.SendTime); err != nil {
			Err(err)
		} else if now-n.SendTime < NOTICE_EXPIRE {
			//expired notice
			res = append(res, &n)
		}
	}

	return
}

func getNoticeSizeFromDB(to string) (size int, err error) {
	// var rows *sql.Rows
	// if !ignoreFriendResp {
	err = dao.CountNoticeByUserStmt.QueryRow(to).Scan(&size)
	// } else {
	// 	err = dao.CountNoticeByUserStmt.QueryRow(to, constants.NOTICE_TYPE_ADD_FRIEND_RESPONSE).Scan(&size)
	// }

	if err != nil {
		Err(err)
	}

	return
}

func getNoticeFromDB(nid int64) (n Notice, err error) {
	if err = dao.GetNoticeStmt.QueryRow(nid).Scan(&n.NoticeID, &n.FromUsername, &n.FromUID, &n.Type, &n.Content); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("notice:", nid, "not exist in db")
			// rao.SetNoticeExist(nid, false, nil)
		} else {
			Err(err)
		}
	}

	return
}

func GetNotice(nid int64) (n Notice, err error) {
	//try get notice from redis first
	// conn := rao.GetConn()
	// defer conn.Close()

	// key := rao.GetNoticeSelfKey(nid)

	// var ns string
	// if ns, err = redis.String(conn.Do("GET", key)); err != nil {
	// 	if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
	// 		Info("can't get notice from redis:", nid)
	// 	} else {
	// 		Err(err)
	// 	}

	// if rao.NoticeExist(nid, &conn) {
	return getNoticeFromDB(nid)
	// }
	// } else {
	// 	ns = ns[1 : len(ns)-1]
	// 	if err = json.Unmarshal([]byte(ns), &n); err != nil {
	// 		Err(err)
	// 		return
	// 	}
	// }

	// return
}

//timeout or player reject
func RejectNotice(n *Notice, rUsername string, rUID int64, ntype int) (gotPieces map[string]int, err error) {
	if (n.Type == constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_TIMEOUT || n.Type == constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_RESP) && ntype == constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_TIMEOUT {
		var sc SwapCloPieceContent
		json.Unmarshal([]byte(n.Content), &sc)
		//do nothing
		if n.Type == constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_TIMEOUT {
			mc := MailContent{
				Content:   L("cp1"),
				CloPieces: map[string]int{fmt.Sprintf("%s:%d", sc.FromCloID, sc.FromCloPart): sc.FromCloPieceCnt},
			}
			mcb, _ := json.Marshal(mc)
			mails.SendToOne("", n.FromUsername, L("cp2"), string(mcb), 0, 0, "", true, time.Now().Unix()+constants.DEFAULT_MAIL_EXPIRE_TIME)
		} else if n.Type == constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_RESP {
			mc := MailContent{
				Content:   L("cp3"),
				CloPieces: map[string]int{fmt.Sprintf("%s:%d", sc.ToCloID, sc.ToCloPart): sc.ToCloPieceCnt},
			}
			mcb, _ := json.Marshal(mc)
			mails.SendToOne("", n.FromUsername, L("cp4"), string(mcb), 0, 0, "", true, time.Now().Unix()+constants.DEFAULT_MAIL_EXPIRE_TIME)
		}

		Info("delete notice:", n.NoticeID)
		if err = DeleteNotice(rUsername, n.NoticeID); err != nil {
			return
		}
	} else {
		if n.Type == constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE {
			Info("send delete_notice notice:", n.NoticeID)

			no := Notice{
				// NoticeID:     nid,
				FromUsername: rUsername,
				FromUID:      rUID,
				Type:         ntype,
				Content:      n.Content,
			}

			if err = flushNotice(n.FromUsername, n.FromUID, &no); err != nil {
				if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
					Info("duplicate notice:", rUsername, ntype, n.Content)
				} else {
					Err(err)
				}
				return
			}
		} else if n.Type == constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_TIMEOUT || n.Type == constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_RESP || n.Type == constants.NOTICE_TYPE_REJECT_SWAP_CLO_PIECE {
			//swap fail.add from piece
			gotPieces = make(map[string]int)
			var sc SwapCloPieceContent
			if err = json.Unmarshal([]byte(n.Content), &sc); err != nil {
				Err(err)
				return
			}

			if (n.Type == constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_TIMEOUT || n.Type == constants.NOTICE_TYPE_REJECT_SWAP_CLO_PIECE) && sc.FromCloPieceCnt > 0 {
				gotPieces[fmt.Sprintf("%s:%d", sc.FromCloID, sc.FromCloPart)] = sc.FromCloPieceCnt
			} else if sc.ToCloPieceCnt > 0 {
				gotPieces[fmt.Sprintf("%s:%d", sc.ToCloID, sc.ToCloPart)] = sc.ToCloPieceCnt
			}
		}

		Info("delete notice:", n.NoticeID)
		if err = DeleteNotice(rUsername, n.NoticeID); err != nil {
			return
		}
	}

	return
}

func GetUnexpireNoticeList(toUsername string) (res []Notice, cnt int, err error) {
	var rows *sql.Rows
	if rows, err = dao.QueryUnexpireNoticeStmt.Query(toUsername, constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_RESP, constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_TIMEOUT, time.Now().Unix()-NOTICE_EXPIRE); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]Notice, 0)
		for rows.Next() {
			var n Notice
			if err = rows.Scan(&n.NoticeID, &n.FromUsername, &n.FromUID, &n.Type, &n.Content, &n.SendTime); err != nil {
				Err(err)
				return
			}
			res = append(res, n)
		}
	}

	return
}

func GetNoticeList(toUsername string, page int, pageCnt int, ignoreNotice bool) (res []*Notice, cnt int, err error) {
	if ignoreNotice {
		return
	}
	// if page*pageCnt > NOTICE_LIST_SIZE {
	// 	return
	// }

	// conn := rao.GetConn()
	// defer conn.Close()

	// setKey := rao.GetUserNoticeListKey(toUsername)
	// var result []interface{}
	// now := time.Now().Unix()
	start := (page - 1) * pageCnt
	// end := page*pageCnt - 1

	// var exist bool
	// if exist, err = rao.ExistKey(setKey, &conn, 0); err != nil {
	// 	return
	// } else if exist {
	// 	conn.Send("ZREMRANGEBYSCORE", setKey, 0, now-NOTICE_EXPIRE)
	// 	conn.Send("ZCOUNT", setKey, now-NOTICE_EXPIRE+1, now)
	// 	conn.Send("ZREVRANGE", setKey, (page-1)*pageCnt, page*pageCnt-1)
	// 	conn.Flush()
	// 	conn.Receive()

	// 	if cnt, err = redis.Int(conn.Receive()); err != nil {
	// 		Err(err)
	// 		return
	// 	}

	// 	//get ids
	// 	ids := make([]interface{}, 0)
	// 	if result, err = redis.Values(conn.Receive()); err != nil {
	// 		Err(err)
	// 		return
	// 	}

	// 	for len(result) > 0 {
	// 		var nid int64
	// 		if result, err = redis.Scan(result, &nid); err != nil {
	// 			Err(err)
	// 			return
	// 		}

	// 		ids = append(ids, rao.GetNoticeSelfKey(nid))
	// 	}

	// 	if len(ids) > 0 {
	// 		//get notices
	// 		if result, err = redis.Values(conn.Do("MGET", ids...)); err != nil {
	// 			Err(err)
	// 			return
	// 		}

	// 		for len(result) > 0 && len(res) < pageCnt {
	// 			var ns string
	// 			if result, err = redis.Scan(result, &ns); err != nil {
	// 				Err(err)
	// 			}

	// 			if len(ns) > 3 {
	// 				ns = ns[1 : len(ns)-1]

	// 				var n Notice
	// 				if err = json.Unmarshal([]byte(ns), &n); err != nil {
	// 					Err(err)
	// 					return
	// 				}

	// 				res = append(res, &n)
	// 			}
	// 		}
	// 	}

	// 	if cnt > 0 && len(res) <= 0 {
	// 		Info("notice cache error.all size:", cnt, "actual page size:", len(res), "player:", toUsername, "key:", setKey)
	// 		cnt = len(res)
	// 	}
	// } else {
	// var res1 []*Notice
	if cnt, err = getNoticeSizeFromDB(toUsername); err != nil {
		return
	}

	if res, err = getNoticeListFromDB(toUsername, start, pageCnt); err != nil {
		return
	}
	// else if len(res1) > 0 {
	// 	for _, n := range res1 {
	// 		flushNotice(toUsername, toUID, n, &conn, true)
	// 	}

	// 	if end > len(res1)-1 {
	// 		end = len(res1) - 1
	// 	}

	// 	if start > len(res1)-1 || start > end {
	// 		return
	// 	} else {
	// 		res = res1[start : end+1]
	// 		cnt = len(res1)
	// 	}

	// 	return
	// }
	// }

	return
}

func AddSwapCloPieceNotice(fromUsername, fromNickname string, fromUID int64, toUsername, toNickname string, toUID int64, fromCloID string, fromPart int, fromCnt int, toCloID string, toPart int, toCnt int, isRequest bool) (err error) {
	info := SwapCloPieceContent{
		FromCloID:       fromCloID,
		FromCloPart:     fromPart,
		FromCloPieceCnt: fromCnt,
		ToCloID:         toCloID,
		ToCloPart:       toPart,
		ToCloPieceCnt:   toCnt,
	}

	var content []byte
	if content, err = json.Marshal(info); err != nil {
		Err(err)
		return
	} else {
		nType := constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE
		fu := fromUsername
		fn := fromNickname
		fuid := fromUID
		tu := toUsername
		tuid := toUID

		if !isRequest {
			nType = constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_RESP
			fu = toUsername
			fn = toNickname
			fuid = toUID
			tu = fromUsername
			tuid = fromUID
		}

		notice := Notice{
			// NoticeID:     nid,
			FromUsername: fu,
			FromNickname: fn,
			FromUID:      fuid,
			Type:         nType,
			Content:      string(content),
		}

		return flushNotice(tu, tuid, &notice)
	}
}

func getUnreadUnexpireNoticeCnt(toUsername string) (cnt int) {
	if err := dao.QueryUnexpireNoticeCntStmt.QueryRow(toUsername, constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_RESP, constants.NOTICE_TYPE_SWAP_CLOTHES_PIECE_TIMEOUT, time.Now().Unix()-NOTICE_EXPIRE).Scan(&cnt); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("can't find notice in db:", toUsername)
		} else {
			Err(err)
		}
	}
	return
}

func GetUnreadNoticeCount(toUsername string, ignoreNotice bool) int {
	if ignoreNotice {
		return 0
	}
	// conn := rao.GetConn()
	// defer conn.Close()

	// key := rao.GetUserNoticeListKey(toUsername)
	// now := time.Now().Unix()

	// if exist, err := rao.ExistKey(key, &conn, 0); err != nil {
	// 	return 0
	// } else if exist {
	// 	if cnt, err := redis.Int(conn.Do("ZCOUNT", key, now-NOTICE_EXPIRE+1, now)); err != nil {
	// 		Err(err)
	// 		return 0
	// 	} else {
	// 		return cnt
	// 	}
	// } else {
	// 	if res, err := getNoticeListFromDB(toUsername); err != nil {
	// 		return 0
	// 	} else if len(res) > 0 {
	// 		for _, n := range res {
	// 			flushNotice(toUsername, toUID, n, &conn, true)
	// 		}

	// 		return len(res)
	// 	} else {
	// 		return 0
	// 	}
	// }
	cnt, _ := getNoticeSizeFromDB(toUsername)
	return cnt
}

func addNoticeToDB(toUsername string, toUserUID int64, n *Notice) (id int64, err error) {
	token := fmt.Sprintf("%d_%d_%d", n.FromUID, toUserUID, time.Now().Unix())
	if n.Type == constants.NOTICE_TYPE_ADD_FRIEND {
		token = fmt.Sprintf("%s%d_%d", constants.NOTICE_ID_ADD_FRIEND_PREFIX, n.FromUID, toUserUID)
	}

	var res sql.Result
	if res, err = dao.AddNoticeStmt.Exec(n.FromUsername, n.FromNickname, n.FromUID, toUsername, toUserUID, n.Type, n.Content, token, n.SendTime); err != nil {
		if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
			Info("already send friend notice:", n.FromUsername, toUsername)
		} else {
			Err(err)
		}

		return
	} else if id, err = res.LastInsertId(); err != nil {
		Err(err)
		return
	}
	return
}

func deleteNoticeInDB(nid int64) (err error) {
	if _, err = dao.DeleteNoticeStmt.Exec(nid); err != nil {
		Err(err)
		return
	}

	return
}

func flushNotice(toUsername string, toUserUID int64, n *Notice) (err error) {
	n.SendTime = time.Now().Unix()
	//add to db first
	if n.NoticeID, err = addNoticeToDB(toUsername, toUserUID, n); err != nil {
		return
	}

	// Info("flush notice:", n.NoticeID)

	// if conn == nil {
	// 	c := rao.GetConn()
	// 	defer c.Close()

	// 	conn = &c
	// }

	// var data []byte
	// if data, err = json.Marshal(n); err != nil {
	// 	Err(err)
	// 	return
	// }

	// setKey := rao.GetUserNoticeListKey(toUsername)
	// key := rao.GetNoticeSelfKey(n.NoticeID)

	// (*conn).Send("SET", key, "'"+string(data)+"'")

	// (*conn).Send("ZADD", setKey, n.SendTime, n.NoticeID)
	// (*conn).Send("ZREMRANGEBYSCORE", setKey, 0, time.Now().Unix()-NOTICE_EXPIRE)

	// (*conn).Send("EXPIRE", setKey, constants.DEFAULT_EXPIRE_TIME)
	// (*conn).Send("EXPIRE", key, constants.DEFAULT_EXPIRE_TIME)

	// if err = (*conn).Flush(); err != nil {
	// 	Err(err)
	// 	return
	// }
	return
}

func DeleteNotice(toUsername string, nid int64) (err error) {
	//delete notice in db
	if err = deleteNoticeInDB(nid); err != nil {
		Err(err)
		return
	}

	// conn := rao.GetConn()
	// defer conn.Close()

	// setKey := rao.GetUserNoticeListKey(toUsername)
	// key := rao.GetNoticeSelfKey(nid)
	// conn.Send("ZREM", setKey, nid)
	// conn.Send("DEL", key)

	// if err = conn.Flush(); err != nil {
	// 	Err(err)
	// 	return
	// }

	return
}

func AddFriendNotice(fromUsername string, fromNickname string, fromLevel, fromExp int, fromUID int64, fromHead string, fromVIP int, fromBCCnt int, fromGCCnt int, fromPKPoint int, toUsername string, toUID int64, isRequest bool) (err error) {
	nType := constants.NOTICE_TYPE_ADD_FRIEND
	var content string

	if !isRequest {
		nType = constants.NOTICE_TYPE_ADD_FRIEND_RESPONSE

		content = fromNickname
	} else {
		var dLevel int
		var d GFDesigner
		if d, err = custom.GetDesignerInfo(fromUsername, true, nil); err != nil {
			return
		} else {
			dLevel = common.GenDesignerLevel(d.Points)
		}

		f := FriendInfo{
			UID:            fromUID,
			Username:       fromUsername,
			Nickname:       fromNickname,
			Head:           fromHead,
			VIP:            fromVIP,
			BoyClothesCnt:  fromBCCnt,
			GirlClothesCnt: fromGCCnt,
			DesignerLevel:  dLevel,
			PKPoint:        fromPKPoint,
			Level:          fromLevel,
			Exp:            fromExp,
		}

		var cb []byte
		if cb, err = json.Marshal(f); err != nil {
			Err(err)
			return
		} else {
			content = string(cb)
		}
	}

	n := Notice{
		FromUsername: fromUsername,
		FromNickname: fromNickname,
		FromUID:      fromUID,
		Type:         nType,
		Content:      content,
	}

	//friend notice only add to redis
	return flushNotice(toUsername, toUID, &n)
}

func GetAllNewFriendNotice(toUsername string) (adds []Friend, rems []Friend, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetUserNewFriendNoticeListKey(toUsername)

	var values []interface{}

	conn.Send("ZRANGEBYSCORE", key, 0, 1, "WITHSCORES")
	conn.Send("DEL", key)

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	if values, err = redis.Values(conn.Receive()); err != nil {
		if err.Error() == "redigo: nil returned" {
			Info("no new friend for:", toUsername)
			err = nil
		} else {
			Err(err)
			return
		}
	}

	adds = make([]Friend, 0)
	rems = make([]Friend, 0)
	for len(values) > 0 {
		var fs string
		if values, err = redis.Scan(values, &fs); err != nil {
			Err(err)
			return
		}
		var addType int
		if values, err = redis.Scan(values, &addType); err != nil {
			Err(err)
			return
		}

		fs = fs[1 : len(fs)-1]

		var f Friend
		if err = json.Unmarshal([]byte(fs), &f); err != nil {
			Err(err)
			return
		}

		if addType == 0 {
			adds = append(adds, f)
		} else {
			rems = append(rems, f)
		}

	}

	return

}

type Friend struct {
	Username string
	UID      int64
}

// func AddNewFriendNotice(toUsername string, fUsername string, fUID int64, addFriend bool) {
// 	f := Friend{
// 		Username: fUsername,
// 		UID:      fUID,
// 	}

// 	addType := 0 //add
// 	if !addFriend {
// 		addType = 1
// 	}

// 	if fb, err := json.Marshal(f); err != nil {
// 		Err(err)
// 		return
// 	} else {
// 		conn := rao.GetConn()
// 		defer conn.Close()

// 		key := rao.GetUserNewFriendNoticeListKey(toUsername)

// 		conn.Send("ZADD", key, addType, "'"+string(fb)+"'")
// 		conn.Send("EXPIRE", key, NEW_FRIEND_NOTICE_EXPIRE)
// 		if err = conn.Flush(); err != nil {
// 			Err(err)
// 		}
// 	}
// }
