package mail

import (
	cm "common_db/mail"
	"constants"
	"dao"
	"database/sql"
	// "encoding/json"
	// "fmt"
	"github.com/garyburd/redigo/redis"
	. "logger"
	// "time"
)

import (
	"rao"
	. "types"
)

func List(start, limit int) (res []*DMail, err error) {
	var rows *sql.Rows
	if rows, err = dao.ListMailStmt.Query(start, limit); err != nil {
		Err(err)
		return
	}
	defer rows.Close()
	res = make([]*DMail, 0)
	// var cloStr string
	for rows.Next() {
		mail := DMail{}

		if err = rows.Scan(&mail.MailID, &mail.To, &mail.Time, &mail.ExpireTime); err != nil {
			Err(err)
			return
		}

		// if err = json.Unmarshal([]byte(cloStr), &mail.Clothes); err != nil {
		// 	Err(err)
		// 	return
		// }

		res = append(res, &mail)
	}

	return
}

func Find(to string, ignoreWarNotice bool) (result []*Mail, err error) {
	// key := rao.GetMailListKey(to)
	// now := time.Now().Unix()

	// conn := rao.GetConn()
	// defer conn.Close()

	// var exist bool
	// if exist, err = rao.ExistKey(key, &conn, constants.MAIL_EXPIRE_TIME); err != nil {
	// 	return
	// } else if exist {
	// 	var values []interface{}
	// 	if values, err = redis.Values(conn.Do("ZREVRANGE", key, (pageNumber-1)*pageCount, pageCount*pageNumber-1)); err != nil {
	// 		Err(err)
	// 		return
	// 	}

	// 	result = make([]*Mail, 0, pageCount)
	// 	for len(values) > 0 {
	// 		var mj string
	// 		values, err = redis.Scan(values, &mj)

	// 		if err != nil {
	// 			Err(err)
	// 			return
	// 		}

	// 		mj = mj[1 : len(mj)-1]
	// 		mail := Mail{}

	// 		err = json.Unmarshal([]byte(mj), &mail)
	// 		if err != nil {
	// 			Err(err)
	// 			return
	// 		}

	// 		if mail.Type == constants.MAIL_TYPE_NOTIFY_GUILD_WAR && ignoreWarNotice {
	// 			go DeleteMail(to, mail.MailID)
	// 			continue
	// 		}

	// 		if mail.ExpireTime != 0 && mail.ExpireTime <= now {
	// 			//delete mail
	// 			go DeleteMail(to, mail.MailID)
	// 		} else {
	// 			result = append(result, &mail)
	// 		}
	// 	}
	// } else {
	if result, err = cm.FindAllFromDB(to, dao.GetAllMailStmt, dao.DeleteMailStmt); err != nil {
		return
	}
	// 	if len(result) > 0 {
	// 		// rao.SetUserHasMail(to, true, &conn)
	// 		cm.RecoverMailCache(to, result, &conn)
	// 	}
	// }

	return
}

func FindOne(to string, mailID int64) (r Mail, err error) {
	// key := rao.GetMailListKey(to)

	// conn := rao.GetConn()
	// defer conn.Close()

	// values, err := redis.Values(conn.Do("ZRANGEBYSCORE", key, mailID, mailID))

	// if err != nil {
	// 	Err(err)
	// 	return
	// }

	// //just get the first one
	// r = Mail{}
	// if len(values) > 0 {
	// 	var mj string
	// 	_, err = redis.Scan(values, &mj)

	// 	if err != nil {
	// 		Err(err)
	// 		return
	// 	}

	// 	mj = mj[1 : len(mj)-1]

	// 	err = json.Unmarshal([]byte(mj), &r)
	// 	if err != nil {
	// 		Err(err)
	// 		return
	// 	}

	// 	return
	// } else {
	return cm.FindOneFromDB(mailID, dao.GetMailStmt)
	// }

	// return
}

func GetUnreadCount(to string, ignoreWarNotice bool) (count int, err error) {
	// conn := rao.GetConn()
	// defer conn.Close()

	// now := time.Now().Unix()

	// key := rao.GetMailListKey(to)

	// var exist bool
	// if exist, err = rao.ExistKey(key, &conn, constants.MAIL_EXPIRE_TIME); err != nil {
	// 	return
	// } else if exist {
	// 	var values []interface{}
	// 	values, err = redis.Values(conn.Do("ZREVRANGE", key, 0, constants.MAIL_MAX_SIZE-1))

	// 	if err != nil {
	// 		Err(err)
	// 		return
	// 	}

	// 	var mj string
	// 	var mail Mail
	// 	for len(values) > 0 {
	// 		values, err = redis.Scan(values, &mj)

	// 		if err != nil {
	// 			Err(err)
	// 			return
	// 		}

	// 		mj = mj[1 : len(mj)-1]

	// 		err = json.Unmarshal([]byte(mj), &mail)

	// 		if err != nil {
	// 			Err(err)
	// 			return
	// 		}

	// 		if mail.Type == constants.MAIL_TYPE_NOTIFY_GUILD_WAR && ignoreWarNotice {
	// 			go DeleteMail(to, mail.MailID)
	// 			continue
	// 		}

	// 		if !mail.Read {
	// 			if mail.ExpireTime != 0 && mail.ExpireTime <= now {
	// 				//delete mail
	// 				go DeleteMail(to, mail.MailID)
	// 			} else {
	// 				count++
	// 			}
	// 		}
	// 	}
	// } else {
	var result []*Mail
	if result, err = cm.FindAllFromDB(to, dao.GetAllMailStmt, dao.DeleteMailStmt); err != nil {
		return
	}
	if len(result) > 0 {
		// rao.SetUserHasMail(to, true, &conn)
		// cm.RecoverMailCache(to, result, &conn)
		for _, m := range result {
			if m.Type == constants.MAIL_TYPE_NOTIFY_GUILD_WAR && ignoreWarNotice {
				go DeleteMail(to, m.MailID)
				continue
			}

			if !m.Read {
				count++
			}
		}
	} else {
		count = 0
		// rao.SetUserHasMail(to, false, &conn)
	}
	// }

	return
}

func markReadInDB(mailID int64) (err error) {
	if _, err = dao.MarkMailReadStmt.Exec(mailID); err != nil {
		Err(err)
		return
	}

	return
}

// func markReadInRedis(to string, mailID int64) (err error) {
// 	key := rao.GetMailListKey(to)

// 	conn := rao.GetConn()
// 	defer conn.Close()

// 	var exist bool
// 	if exist, err = rao.ExistKey(key, &conn, constants.MAIL_EXPIRE_TIME); err != nil {
// 		return
// 	} else if exist {
// 		var values []interface{}
// 		if values, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, mailID, mailID)); err != nil {
// 			Err(err)
// 			return
// 		}

// 		//just get the first one
// 		if len(values) > 0 {
// 			var mj string
// 			if _, err = redis.Scan(values, &mj); err != nil {
// 				Err(err)
// 				return
// 			}

// 			mj = mj[1 : len(mj)-1]
// 			mail := Mail{}

// 			err = json.Unmarshal([]byte(mj), &mail)
// 			if err != nil {
// 				Err(err)
// 				return
// 			}
// 			mail.Read = true
// 			var nmb []byte
// 			nmb, err = json.Marshal(mail)
// 			if err != nil {
// 				Err(err)
// 				return err
// 			}

// 			// conn.Send("MULTI")
// 			conn.Send("ZREMRANGEBYSCORE", key, mailID, mailID)
// 			conn.Send("ZADD", key, mailID, "'"+string(nmb)+"'")
// 			if err = conn.Flush(); err != nil {
// 				Err(err)
// 			}

// 			return
// 		}
// 	} else {
// 		//do not recover when write.
// 	}
// 	return
// }

func GetNotSendCronMails() (res []CronMail) {
	res, _ = cm.GetNotSendCronMails(dao.GetCronMailListStmt)
	return
}

func SetRead(to string, mailID int64) (err error) {
	//update read in batch
	// go mail_batch.MarkReadMail(mailID)
	//update in db
	if err = markReadInDB(mailID); err != nil {
		return
	}

	//only udpate mails in redis
	// markReadInRedis(to, mailID)

	return
}

func DeleteMail(to string, mailID int64) (err error) {
	conn := rao.GetConn()
	defer conn.Close()
	return cm.DeleteMail(mailID, to, dao.DeleteMailStmt, &conn)
}

// func SendToList(infos []MailInfo) (err error) {
// 	return cm.SendToList(infos, dao.AddMailStmt, dao.GetAllMailStmt)
// }

//from,to,titie,content,diamond,gold,clothesid,should_delete
//rconn : recover connection only used in recover mode
func SendToOne(mail Mail, rconn *redis.Conn) (err error) {
	// conn := rao.GetConn()
	// defer conn.Close()
	return cm.SendToOne(mail, dao.AddMailStmt)
}
