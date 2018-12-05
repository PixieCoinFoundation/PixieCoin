package rao

// import (
// 	"constants"
// 	"github.com/garyburd/redigo/redis"
// 	. "logger"
// )

// const (
// 	FRIEND_CHANGE_FIELD = "f_c"
// 	HAS_CUSTOM_FIELD    = "h_c"
// 	HAS_MAIL_FIELD      = "h_m"
// 	IS_DESIGNER_FIELD   = "i_d"
// 	EXIST_FIELD         = "e"
// 	HAS_PARTY_FIELD     = "h_pt"

// 	HAS_PK_FIELD = "h_pk"
// )

// //------------------------------------------------------------------------------

// func setVolitileProperty(key string, value int, conn *redis.Conn) {
// 	if conn == nil {
// 		c := GetConn()
// 		defer c.Close()

// 		conn = &c
// 	}

// 	(*conn).Send("SET", key, value)
// 	(*conn).Send("EXPIRE", key, constants.DEFAULT_EXPIRE_TIME)
// 	if err := (*conn).Flush(); err != nil {
// 		Err(err)
// 	}

// 	if _, e := (*conn).Receive(); e != nil {
// 		Err(e)
// 	}
// 	if _, e := (*conn).Receive(); e != nil {
// 		Err(e)
// 	}
// }

// func getVolitileProperty(key string, conn *redis.Conn, defaultValue int) int {
// 	if conn == nil {
// 		c := GetConn()
// 		defer c.Close()

// 		conn = &c
// 	}

// 	(*conn).Send("EXPIRE", key, constants.DEFAULT_EXPIRE_TIME)
// 	(*conn).Send("GET", key)
// 	(*conn).Flush()

// 	var exist int
// 	var err error
// 	var value int

// 	if exist, err = redis.Int((*conn).Receive()); err != nil {
// 		Err(err)
// 		exist = 0
// 	}

// 	if value, err = redis.Int((*conn).Receive()); err != nil {
// 		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
// 			Info("key not exist in redis:", key)
// 		} else {
// 			Err(err)
// 		}
// 		return defaultValue
// 	}

// 	if exist == 1 {
// 		return value
// 	} else {
// 		return defaultValue
// 	}
// }

// func setUserProperty(username string, field string, value int, conn *redis.Conn) {
// 	if conn == nil {
// 		c := GetConn()
// 		defer c.Close()

// 		conn = &c
// 	}

// 	key := constants.USER_KEY_PREFIX + username

// 	(*conn).Send("HSET", key, field, value)
// 	(*conn).Send("EXPIRE", key, constants.USER_EXPIRE_TIME)
// 	(*conn).Flush()
// 	if _, err := (*conn).Receive(); err != nil {
// 		Err(err)
// 	}
// 	if _, err := (*conn).Receive(); err != nil {
// 		Err(err)
// 	}
// }

// func getUserProperty(username string, field string, conn *redis.Conn) (res int, err error) {
// 	if conn == nil {
// 		c := GetConn()
// 		defer c.Close()

// 		conn = &c
// 	}

// 	key := constants.USER_KEY_PREFIX + username

// 	(*conn).Send("HGET", key, field)
// 	(*conn).Send("EXPIRE", key, constants.USER_EXPIRE_TIME)
// 	(*conn).Flush()
// 	if res, err = redis.Int((*conn).Receive()); err != nil {
// 		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
// 			// Info("user property not exist:", key, field)
// 		} else {
// 			Err(err)
// 		}
// 	}
// 	if _, e := (*conn).Receive(); e != nil {
// 		Err(e)
// 	}

// 	return
// }

// func setGlobalProperty(post string, field string, value int, conn *redis.Conn) (int, error) {
// 	if conn == nil {
// 		c := GetConn()
// 		defer c.Close()

// 		conn = &c
// 	}

// 	key := constants.GLOBAL_KEY

// 	return redis.Int((*conn).Do("HSET", key, field+"_"+post, value))
// }

// func getGlobalProperty(post string, field string, conn *redis.Conn) (int, error) {
// 	if conn == nil {
// 		c := GetConn()
// 		defer c.Close()

// 		conn = &c
// 	}

// 	key := constants.GLOBAL_KEY

// 	return redis.Int((*conn).Do("HGET", key, field+"_"+post))
// }

//------------------------------------------------------------------------------

// func SetUserIsDesigner(username string, has bool, conn *redis.Conn) {
// 	e := 1
// 	if !has {
// 		e = 0
// 	}

// 	setUserProperty(username, IS_DESIGNER_FIELD, e, conn)
// }

// func SetUserHasCustom(username string, has bool, conn *redis.Conn) {
// 	e := 1
// 	if !has {
// 		e = 0
// 	}

// 	setUserProperty(username, HAS_CUSTOM_FIELD, e, conn)
// }

// func SetUserHasMail(username string, has bool, conn *redis.Conn) {
// 	e := 1
// 	if !has {
// 		e = 0
// 	}

// 	setUserProperty(username, HAS_MAIL_FIELD, e, conn)
// }

// func SetUserExist(username string, exist bool, conn *redis.Conn) {
// 	e := 1
// 	if !exist {
// 		e = 0
// 	}

// 	setUserProperty(username, EXIST_FIELD, e, conn)
// }

// func SetUserFriendChange(username string, change bool, conn *redis.Conn) {
// 	e := 1
// 	if !change {
// 		e = 0
// 	}

// 	setUserProperty(username, FRIEND_CHANGE_FIELD, e, conn)
// }

// func SetGlobalHasPK(lid string, has bool, conn *redis.Conn) {
// 	e := 1
// 	if !has {
// 		e = 0
// 	}

// 	setGlobalProperty(lid, HAS_PK_FIELD, e, conn)
// }

// func SetHelpHasCmt(helpID int64, has bool, conn *redis.Conn) {
// 	e := 1
// 	if !has {
// 		e = 0
// 	}
// 	setVolitileProperty(GetHelpHasCmtKey(helpID), e, conn)
// }

// func SetMailExist(mailID int64, exist bool, conn *redis.Conn) {
// 	e := 1
// 	if !exist {
// 		e = 0
// 	}
// 	setVolitileProperty(GetMailExistKey(mailID), e, conn)
// }

// func SetCosItemHasCmt(itemID int64, has bool, conn *redis.Conn) {
// 	e := 1
// 	if !has {
// 		e = 0
// 	}
// 	setVolitileProperty(GetCosItemHasCmtKey(itemID), e, conn)
// }

// func SetNoticeExist(nid int64, exist bool, conn *redis.Conn) {
// 	e := 1
// 	if !exist {
// 		e = 0
// 	}
// 	setVolitileProperty(GetNoticeExistKey(nid), e, conn)
// }

//------------------------------------------------------------------------------

// func HelpHasCmt(helpID int64, conn *redis.Conn) bool {
// 	if getVolitileProperty(GetHelpHasCmtKey(helpID), conn, 1) == 1 {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func NoticeExist(nid int64, conn *redis.Conn) bool {
// 	if getVolitileProperty(GetNoticeExistKey(nid), conn, 1) == 1 {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func MailExist(mailID int64, conn *redis.Conn) bool {
// 	if getVolitileProperty(GetMailExistKey(mailID), conn, 1) == 1 {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func CosItemHasCmt(itemID int64, conn *redis.Conn) bool {
// 	if getVolitileProperty(GetCosItemHasCmtKey(itemID), conn, 1) == 1 {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func GlobalHasPK(lid string, conn *redis.Conn) bool {
// 	e := 1
// 	var err error
// 	if e, err = getGlobalProperty(lid, HAS_PK_FIELD, conn); err != nil {
// 		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
// 			Info("global has pk cache not exist")
// 		} else {
// 			Err(err)
// 		}
// 		e = 1
// 	}

// 	if e == 1 {
// 		return true
// 	}

// 	return false
// }

// func UserHasCustom(username string, conn *redis.Conn) bool {
// 	e := 1
// 	var err error
// 	if e, err = getUserProperty(username, HAS_CUSTOM_FIELD, conn); err != nil {
// 		e = 1
// 	}

// 	if e == 1 {
// 		return true
// 	}

// 	return false
// }

// func UserHasMail(username string, conn *redis.Conn) bool {
// 	e := 1
// 	var err error
// 	if e, err = getUserProperty(username, HAS_MAIL_FIELD, conn); err != nil {
// 		e = 1
// 	}

// 	if e == 1 {
// 		return true
// 	}

// 	return false
// }

// func UserIsDesigner(username string, conn *redis.Conn) bool {
// 	e := 1
// 	var err error
// 	if e, err = getUserProperty(username, IS_DESIGNER_FIELD, conn); err != nil {
// 		e = 1
// 	}

// 	if e == 1 {
// 		return true
// 	}

// 	return false
// }

// func UserExist(username string, conn *redis.Conn) bool {
// 	e := 1
// 	var err error
// 	if e, err = getUserProperty(username, EXIST_FIELD, conn); err != nil {
// 		e = 1
// 	}

// 	if e == 1 {
// 		return true
// 	}

// 	return false
// }

// func UserFriendChange(username string, conn *redis.Conn) bool {
// 	e := 1
// 	var err error
// 	if e, err = getUserProperty(username, FRIEND_CHANGE_FIELD, conn); err != nil {
// 		e = 1
// 	}

// 	if e == 1 {
// 		return true
// 	}

// 	return false
// }
