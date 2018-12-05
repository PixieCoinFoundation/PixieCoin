package global

import (
	"constants"
	"github.com/garyburd/redigo/redis"
	. "logger"
	"math/rand"
	"rao"
)

func AddPlayerLoginForRandomVisit(username string) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGlobalUserSetKey()
	conn.Send("SADD", key, username)
	conn.Send("SCARD", key)

	if err := conn.Flush(); err != nil {
		Err(err)
		return
	}
	conn.Receive()

	if size, err := redis.Int(conn.Receive()); err != nil {
		Err(err)
		return
	} else if size > constants.PIXIE_RANDOM_VISIT_SIZE*2 {
		for i := 0; i < constants.PIXIE_REMOVE_RANDOM_VISIT_SIZE; i++ {
			conn.Send("SPOP", key)
		}
		conn.Flush()
	}
}

func GetRandomPlayerForVisit(queryUsername, lastVisitUsername string) (targetUsername string, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGlobalUserSetKey()
	var values []interface{}
	if values, err = redis.Values(conn.Do("SRANDMEMBER", key, 10)); err != nil {
		Err(err)
		return
	} else {
		tmpList := make([]string, 0)
		var tmpUsername string
		for len(values) > 0 {
			if values, err = redis.Scan(values, &tmpUsername); err != nil {
				Err(err)
				return
			}

			if tmpUsername == "" {
				continue
			}

			if queryUsername != "" && tmpUsername == queryUsername {
				continue
			}

			if lastVisitUsername != "" && tmpUsername == lastVisitUsername {
				continue
			}

			tmpList = append(tmpList, tmpUsername)
		}

		if len(tmpList) > 0 {
			targetUsername = tmpList[rand.Intn(len(tmpList))]
		}
	}
	return
}
