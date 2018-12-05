package lock

import (
	"appcfg"
	"github.com/garyburd/redigo/redis"
	"rao"
	"time"

	. "logger"
)

const (
	ALREADY_LOCKED_BY_OTHERS_ERR_MSG = "redigo: nil returned"
)

var expireMillisecond, retryCount, retryGapMillisecond int

func init() {
	expireMillisecond = appcfg.GetInt("redis_lock_expire_millisecond", 10000)   //default 10s
	retryCount = appcfg.GetInt("redis_lock_retry_count", 5)                     //default retry 5 times
	retryGapMillisecond = appcfg.GetInt("redis_lock_retry_gap_millisecond", 50) //default retry gap 50 ms
}

func Lock(key string) (ok bool, err error) {
	c := rao.GetConn()
	defer c.Close()

	var res string

	for i := 0; i < retryCount; i++ {
		if i != 0 {
			time.Sleep(time.Duration(retryGapMillisecond) * time.Millisecond)
		}

		if res, err = redis.String(c.Do("SET", key, 1, "PX", expireMillisecond, "NX")); err != nil {
			if err.Error() == ALREADY_LOCKED_BY_OTHERS_ERR_MSG {
				//already locked by others
				//do nothing.retry.
				err = nil
			} else {
				Err("lock redis key error", err)
				return
			}
		} else if res == "OK" {
			ok = true
			return
		}
	}

	return
}

func Unlock(key string) {
	c := rao.GetConn()
	defer c.Close()

	if _, err := redis.Int(c.Do("DEL", key)); err != nil {
		Err("unlock redis key error", err)
		return
	}
}
