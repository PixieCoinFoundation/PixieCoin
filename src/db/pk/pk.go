package pk

import (
	"common"
	"constants"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	. "logger"
	// "math/rand"
	"rao"
	"strings"
	. "types"
)

const (
	gap = "^_^"
)

//only return GF_PK_MAX_SIZE instances
func Find(username, levelID string, point, pkType int) (ret []*GFPK, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetPKKey(levelID, pkType)
	start, end := common.ComputePKPointRange(point)

	cnt := 0
	var exist bool
	for cnt < 3 {
		um := make(map[string]int)

		if exist, err = rao.ExistKey(key, &conn, 0); err != nil {
			return
		} else if exist {
			ret = make([]*GFPK, 0)
			var values []interface{}
			if values, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, start, end)); err != nil {
				Err(err)
				return
			}

			for len(values) > 0 {
				var pj string
				if values, err = redis.Scan(values, &pj); err != nil {
					Err(err)
					return
				}

				pj = pj[1 : len(pj)-1]

				pk := GFPK{}
				if err = json.Unmarshal([]byte(pj), &pk); err != nil {
					Err(err)
					return
				}
				if pk.Username != username {
					ret = append(ret, &pk)
					um[pk.Username] = 1
				}
			}
		} else {
			break
		}

		if len(um) >= 10 {
			return
		}

		cnt++
		start -= cnt * 100
		end += cnt * 100
	}
	return
}

//conn is not nil means recover mode
func Save(pkInfo *GFPK, rconn *redis.Conn, pkType int) (err error) {
	var pjb []byte
	pjb, err = json.Marshal(pkInfo)
	if err != nil {
		Err(err)
		return
	}

	key := rao.GetPKKey(pkInfo.LevelID, pkType)

	if rconn == nil {
		//normal mode.async add to db for recover
		// go pk_batch.AddPK(pkInfo)
		//add to db first
		// if err = addPKInDB(pkInfo); err != nil {
		// 	return
		// }

		//add to redis
		conn := rao.GetConn()
		defer conn.Close()

		// var size int
		if _, err = conn.Do("ZADD", key, pkInfo.PKPoint, "'"+string(pjb)+"'"); err != nil {
			Err(err)
			return
		}

		// if size > constants.GF_PK_MAX_SIZE*2 {
		// 	go common.TrimRedisList(key, 0, constants.GF_PK_MAX_SIZE-1)
		// }

		// rao.SetGlobalHasPK(pkInfo.LevelID, true, &conn)
		return nil
	} else {
		//recover mode should be rpush---FIFO
		(*rconn).Send("RPUSH", key, "'"+string(pjb)+"'")
		return nil
	}

}

// func recoverPKCache(lid string, ps []*GFPK, conn *redis.Conn) {
// 	for _, p := range ps {
// 		Save(p, conn)
// 	}

// 	if err := (*conn).Flush(); err != nil {
// 		Err(err)
// 	}
// }
func DelWeekRank(weekToken string) {
	delRank(rao.GetPKWeekRankKey(weekToken))
}

func DelGWeekRank(gweekToken string) {
	delRank(rao.GetPKGWeekRankKey(gweekToken))
}

func AddWeekRank(username, nickname string, point int, weekToken string) (rank int, err error) {
	return addRank(username, nickname, rao.GetPKWeekRankKey(weekToken), point)
}

func AddGWeekRank(username, nickname string, point int, gweekToken string) (rank int, err error) {
	return addRank(username, nickname, rao.GetPKGWeekRankKey(gweekToken), point)
}

func GetGWeekRank(gweekToken string) (res []PKRank, err error) {
	return getRank(rao.GetPKGWeekRankKey(gweekToken))
}

func GetWeekRank(weekToken string) (res []PKRank, err error) {
	return getRank(rao.GetPKWeekRankKey(weekToken))
}

func delRank(key string) {
	conn := rao.GetConn()
	defer conn.Close()

	if _, err := conn.Do("DEL", key); err != nil {
		Err(err)
	}
}

func addRank(username, nickname, key string, point int) (rank int, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	element := fmt.Sprintf("%s%s%s", username, gap, nickname)

	conn.Send("ZADD", key, point, element)
	conn.Send("ZREMRANGEBYRANK", key, 0, -constants.PK_GRANK_SIZE-1)
	conn.Send("ZREVRANK", key, element)

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	conn.Receive()
	conn.Receive()

	if rank, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			err = nil
			return
		} else {
			Err(err)
			return
		}
	} else {
		rank += 1
	}

	return
}

func GetPlayerRank(username, nickname, wt, gwt string) (wr int, gwr int) {
	wk := rao.GetPKWeekRankKey(wt)
	gwk := rao.GetPKGWeekRankKey(gwt)

	conn := rao.GetConn()
	defer conn.Close()

	element := fmt.Sprintf("%s%s%s", username, gap, nickname)

	conn.Send("ZREVRANK", wk, element)
	conn.Send("ZREVRANK", gwk, element)

	if err := conn.Flush(); err != nil {
		Err(err)
		return
	}

	if rank, err := redis.Int(conn.Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			// Info("no player week rank:", username)
		} else {
			Err(err)
		}
	} else {
		wr = rank + 1
	}

	if rank, err := redis.Int(conn.Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			// Info("no player gweek rank:", username)
		} else {
			Err(err)
		}
	} else {
		gwr = rank + 1
	}

	return
}

func getRank(key string) (res []PKRank, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	var values []interface{}
	if values, err = redis.Values(conn.Do("ZREVRANGE", key, 0, 99, "WITHSCORES")); err != nil {
		Err(err)
		return
	} else {
		for len(values) > 0 {
			var ele string
			var p int
			if values, err = redis.Scan(values, &ele); err != nil {
				Err(err)
				return
			}

			if values, err = redis.Scan(values, &p); err != nil {
				Err(err)
				return
			}

			eles := strings.Split(ele, gap)

			if len(eles) >= 2 {
				nn := eles[1]

				for i := 2; i < len(eles); i++ {
					nn += gap + eles[i]
				}

				pr := PKRank{
					Username: eles[0],
					Nickname: nn,
					Point:    p,
				}

				res = append(res, pr)
			} else {
				Info("data format error:", ele)
			}
		}
	}

	return
}
