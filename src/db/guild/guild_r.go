package guild

import (
	cg "common_db/guild"
	"constants"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	. "logger"
	"rao"
	. "types"
)

func SetGuildBoardNewestIDR(gid int64, id int64) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGuildBoardNewestKey(gid)
	conn.Send("SET", key, id)
	conn.Send("EXPIRE", key, 72*3600)

	if err := conn.Flush(); err != nil {
		Err(err)
	}
	return
}

func GetGuildBoardNewestIDR(gid int64) (id int64) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGuildBoardNewestKey(gid)

	var err error
	if id, err = redis.Int64(conn.Do("GET", key)); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no guild newest board msg:", gid)
		} else {
			Err(err)
		}
		return
	}

	return
}

func SetGuildMedalRankCache(gwt string, v *string) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGuildMedalRankKey(gwt)
	conn.Send("SET", key, "'"+*v+"'")
	conn.Send("EXPIRE", key, 24*3600)

	if err = conn.Flush(); err != nil {
		Err(err)
	}
	return
}

func GetGuildMedalRankCache(gwt string) (res []Guild, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGuildMedalRankKey(gwt)
	var v string
	res = make([]Guild, 0)
	if v, err = redis.String(conn.Do("GET", key)); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no guild medal rank cache:", key)
		} else {
			Err(err)
		}

		return
	} else if len(v) > 2 {
		v = v[1 : len(v)-1]
		if err = json.Unmarshal([]byte(v), &res); err != nil {
			Err(err)
			return
		}
	}
	return
}

func SetGuildActivityRankCache(wt string, v *string) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGuildActivityRankKey(wt)
	conn.Send("SET", key, "'"+*v+"'")
	conn.Send("EXPIRE", key, 24*3600)

	if err = conn.Flush(); err != nil {
		Err(err)
	}
	return
}

func GetGuildActivityRankCache(wt string) (res []Guild, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGuildActivityRankKey(wt)
	var v string
	res = make([]Guild, 0)
	if v, err = redis.String(conn.Do("GET", key)); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no guild activity rank cache:", wt)
		} else {
			Err(err)
		}
		return
	} else if len(v) > 2 {
		v = v[1 : len(v)-1]
		if err = json.Unmarshal([]byte(v), &res); err != nil {
			Err(err)
			return
		}
	}
	return
}

func AddNewApplyCnt(gid int64, cnt int) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGuildNewApplyCntKey(gid)
	conn.Send("INCRBY", key, cnt)
	conn.Send("EXPIRE", key, 24*3600)

	if err = conn.Flush(); err != nil {
		Err(err)
	}
	return
}

func ReadNewApply(gid int64) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	if _, err = conn.Do("DEL", rao.GetGuildNewApplyCntKey(gid)); err != nil {
		Err(err)
	}
	return
}

func GetNewApplyCnt(gid int64) (cnt int, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	if cnt, err = redis.Int(conn.Do("GET", rao.GetGuildNewApplyCntKey(gid))); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			// Info("no new apply for guild id:", gid)
			err = nil
		} else {
			Err(err)
		}
	}
	return
}

func DelGuildRank(gid int64, wt string, gwt string) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	akey := rao.GetGuildActivityWeekKey(wt)
	mkey := rao.GetGuildMedalGWeekKey(gwt)
	Info("del guild rank:", gid, akey, mkey)

	conn.Send("ZREM", akey, gid)
	conn.Send("ZREM", mkey, gid)
	// conn.Do("ZREMRANGEBYRANK", key, 0, -constants.GUILD_RANK_SIZE-1)

	if err = conn.Flush(); err != nil {
		Err(err)
	}

	return
}

func RefreshGuildMedalRank(gid int64, medal int, gwt string) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGuildMedalGWeekKey(gwt)

	conn.Do("ZADD", key, medal, gid)
	conn.Do("ZREMRANGEBYRANK", key, 0, -constants.GUILD_RANK_SIZE-1)

	if err = conn.Flush(); err != nil {
		Err(err)
	}

	return
}

func AddGuildMedalR(gid int64, medal int, gwt string) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGuildMedalGWeekKey(gwt)
	conn.Do("ZINCRBY", key, medal, gid)
	conn.Do("ZREMRANGEBYRANK", key, 0, -constants.GUILD_RANK_SIZE-1)

	if err = conn.Flush(); err != nil {
		Err(err)
	}

	return
}

func GetGuildMedalRank(gid int64, gwt string) (r int, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	return cg.GetGuildMedalRank(gid, gwt, &conn)
}

func ListGuildMedalRank(gwt string) (res []int64, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGuildMedalGWeekKey(gwt)

	var values []interface{}
	if values, err = redis.Values(conn.Do("ZREVRANGE", key, 0, 99)); err != nil {
		Err(err)
		return
	}

	res = make([]int64, 0)
	for len(values) > 0 {
		var id int64
		if values, err = redis.Scan(values, &id); err != nil {
			Err(err)
			return
		}

		res = append(res, id)
	}

	return
}

func RefreshGuildActivityRank(gid int64, activity int, wt string) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGuildActivityWeekKey(wt)

	conn.Do("ZADD", key, activity, gid)
	conn.Do("ZREMRANGEBYRANK", key, 0, -constants.GUILD_RANK_SIZE-1)

	if err = conn.Flush(); err != nil {
		Err(err)
	}

	return
}

func AddGuildActivityR(gid int64, activity int, wt string) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGuildActivityWeekKey(wt)
	conn.Do("ZINCRBY", key, activity, gid)
	conn.Do("ZREMRANGEBYRANK", key, 0, -constants.GUILD_RANK_SIZE-1)

	if err = conn.Flush(); err != nil {
		Err(err)
	}

	return
}

func GetGuildActivityRank(gid int64, wt string) (r int, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	return cg.GetGuildActivityRank(gid, wt, &conn)
}

func ListGuildActivityRank(wt string) (res []int64, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetGuildActivityWeekKey(wt)

	var values []interface{}
	if values, err = redis.Values(conn.Do("ZREVRANGE", key, 0, 99)); err != nil {
		Err(err)
		return
	}

	res = make([]int64, 0)
	for len(values) > 0 {
		var id int64
		if values, err = redis.Scan(values, &id); err != nil {
			Err(err)
			return
		}

		res = append(res, id)
	}

	return
}
