package record

import (
	"constants"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	. "logger"
	"rao"
	. "types"
)

func UpsertPlayerTheaterRank(rankKey, username, nickname string, score int) {
	conn := rao.GetConn()
	defer conn.Close()

	su := SimpleUserInfo{
		Username: username,
		Nickname: nickname,
	}
	sub, _ := json.Marshal(su)

	conn.Send("ZADD", rankKey, score, "'"+string(sub)+"'")
	conn.Send("EXPIRE", rankKey, 30*24*3600)

	if err := conn.Flush(); err != nil {
		Err(err)
	}
}

func GetTheaterRank(rankKey, un, nn string) (res []PlayerTheaterInfo, rank int, mscore int, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	su := SimpleUserInfo{
		Username: un,
		Nickname: nn,
	}
	sub, _ := json.Marshal(su)
	sus := "'" + string(sub) + "'"

	conn.Send("ZREVRANGE", rankKey, 0, 99, "WITHSCORES")
	conn.Send("ZREVRANK", rankKey, sus)
	conn.Send("ZSCORE", rankKey, sus)

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	res = make([]PlayerTheaterInfo, 0)
	var values []interface{}
	if values, err = redis.Values(conn.Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no theater rank:", rankKey)
			err = nil
			return
		} else {
			Err(err)
			return
		}
	} else {
		for len(values) > 0 {
			var us string
			if values, err = redis.Scan(values, &us); err != nil {
				Err(err)
				return
			} else if len(us) > 2 {
				us = us[1 : len(us)-1]

				var sui SimpleUserInfo
				if err = json.Unmarshal([]byte(us), &sui); err != nil {
					Err(err)
					return
				}

				var score int
				if values, err = redis.Scan(values, &score); err != nil {
					Err(err)
					return
				}

				usi := PlayerTheaterInfo{
					SimpleUserInfo: sui,
					Score:          score,
				}
				res = append(res, usi)
			}
		}
	}

	if rank, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no theater rank:", rankKey)
			err = nil
			return
		} else {
			Err(err)
			return
		}
	} else {
		rank += 1
	}

	if mscore, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no theater rank:", rankKey)
			err = nil
			return
		} else {
			Err(err)
			return
		}
	}

	return
}
