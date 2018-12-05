package rao

import (
	"appcfg"
	"constants"
	. "logger"
	"time"

	"github.com/garyburd/redigo/redis"
)

//id key names
// var ID_KEY_NAMES []string

// var RConn redis.Conn
var rpool *redis.Pool

func init() {
	if appcfg.GetServerType() == constants.SERVER_TYPE_GL_TS {
		return
	}

	if appcfg.GetServerType() == constants.SERVER_TYPE_SIMPLE_TEST {
		return
	}

	// ID_KEY_NAMES = append(ID_KEY_NAMES, "gf_id_custom_key", "gf_boot_count_key", "gf_id_help_key", "gf_id_order_key", "gf_id_cosplay_key", "gf_id_cos_item_key", "gf_id_mail_key", "gf_id_notice_key", "general_uid_key")
	dindex := appcfg.GetInt("redis_db_index", 0)
	rpool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		MaxActive:   appcfg.GetInt("redis_per_max_conn", 2500),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", appcfg.GetString("redis_ipport", "127.0.0.1:6379"), redis.DialConnectTimeout(30*time.Second), redis.DialReadTimeout(30*time.Second), redis.DialWriteTimeout(30*time.Second))
			if err != nil {
				return nil, err
			}
			pwd := appcfg.GetString("redis_password", "")
			if pwd != "" {
				if _, err = c.Do("AUTH", pwd); err != nil {
					c.Close()
					return nil, err
				}
			}

			if dindex > 0 {
				if _, err = c.Do("SELECT", dindex); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	//test redis connection
	c := GetConn()
	defer c.Close()

	if _, err := c.Do("PING"); err != nil {
		panic(err)
	}
}

// func GetUID() (uid int64, err error) {
// 	return getId(ID_KEY_NAMES[8])
// }

func GetConn() (c redis.Conn) {
	return rpool.Get()
}

func GetChannelPlayerNum(channel string) int {
	c := rpool.Get()
	defer c.Close()

	if nn, err := redis.Int(c.Do("GET", channel)); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			//can initialize
			Info("no redis key:", channel)
			return 0
		} else {
			Err(err)
			return 0
		}
	} else {
		return nn
	}
}

func AddChannelPlayerNum(channel string, cnt int) {
	c := rpool.Get()
	defer c.Close()

	if _, err := c.Do("INCRBY", channel, cnt); err != nil {
		Err(err)
	}
}

// func GetCustomID() (id int64, err error) {
// 	return getId(ID_KEY_NAMES[0])
// }

// func GetBootCount() (cnt int64, err error) {
// 	return getId(ID_KEY_NAMES[1])
// }

// func GetHelpID() (cnt int64, err error) {
// 	return getId(ID_KEY_NAMES[2])
// }

// func GetOrderID() (cnt int64, err error) {
// 	return getId(ID_KEY_NAMES[3])
// }

// func GetCosplayID() (cnt int64, err error) {
// 	return getId(ID_KEY_NAMES[4])
// }

// func GetCosItemID() (cnt int64, err error) {
// 	return getId(ID_KEY_NAMES[5])
// }

// func GetMailID() (cnt int64, err error) {
// 	return getId(ID_KEY_NAMES[6])
// }

// func GetNoticeID() (id string, err error) {
// 	var intID int64
// 	if intID, err = getId(ID_KEY_NAMES[7]); err != nil {
// 		return
// 	} else {
// 		id = strconv.FormatInt(intID, 10)
// 	}

// 	return
// }

func getId(key string) (id int64, err error) {
	conn := GetConn()
	defer conn.Close()

	id, err = redis.Int64(conn.Do("INCR", key))
	return
}

func ExistKey(key string, conn *redis.Conn, expireTime int) (exist bool, err error) {
	var rc redis.Conn
	if conn == nil {
		rc = GetConn()
		defer rc.Close()
	} else {
		rc = *conn
	}

	var res int
	if expireTime > 0 {
		if res, err = redis.Int(rc.Do("EXPIRE", key, expireTime)); err != nil {
			Err(err)
			return
		} else {
			return res == 1, nil
		}
	} else {
		if res, err = redis.Int(rc.Do("EXISTS", key)); err != nil {
			Err(err)
			return
		} else {
			return res == 1, nil
		}
	}

}

func ExpireKey(key string, ts int) {
	conn := GetConn()
	defer conn.Close()

	if _, err := conn.Do("EXPIRE", key, ts); err != nil {
		Err(err)
	}
}
