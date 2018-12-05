package items

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	. "logger"
	"rao"
)

import (
	. "types"
)

const (
	COS_ITEM_LIST_SIZE = 1000
)

//player upload cos item
//rconn not nil means recover mode
func Add(ci *CosItem) (playerSize int, err error) {
	// if rconn == nil {
	if ci.ItemID, err = addCosItemToDB(ci); err != nil {
		return
	}

	if playerSize, err = countPlayerCos(ci.CosplayID); err != nil {
		return
	}

	// }

	scoreKey := rao.GetCosItemRankSetKey(ci.CosplayID)
	userKey := rao.GetCosItemUserListKey(ci.CosplayID, ci.Username)
	selfKey := rao.GetCosItemSelfKey(ci.ItemID)

	cir := getElement(ci)

	// if rconn == nil {
	//add to redis
	conn := rao.GetConn()
	defer conn.Close()

	//score set
	conn.Send("ZADD", scoreKey, ci.Score, ci.ItemID)

	//user list
	conn.Send("LPUSH", userKey, ci.ItemID)

	//item self
	conn.Send("HMSET", redis.Args{}.Add(selfKey).AddFlat(cir)...)

	err = conn.Flush()
	if err != nil {
		Err(err)
		return
	}
	// } else {
	// 	//recover mode
	// 	(*rconn).Send("ZADD", scoreKey, ci.Score, ci.ItemID)
	// 	(*rconn).Send("LPUSH", userKey, ci.ItemID)
	// }

	return
}

func FindByPageAndCountOrderByScore(cosplayID int64, page int, pageCount int) (result []CosItem, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	start := pageCount * (page - 1)
	end := page*pageCount - 1

	var values []interface{}
	if values, err = redis.Values(conn.Do("ZREVRANGE", rao.GetCosItemRankSetKey(cosplayID), start, end, "WITHSCORES")); err != nil {
		Err(err)
		return
	}

	scores := make([]int, 0)
	for len(values) > 0 {
		var itemID int64
		if values, err = redis.Scan(values, &itemID); err != nil {
			Err(err)
			return
		}

		var score int
		if values, err = redis.Scan(values, &score); err != nil {
			Err(err)
			return
		}

		scores = append(scores, score)

		conn.Send("HGETALL", rao.GetCosItemSelfKey(itemID))
	}

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	length := len(scores)
	result = make([]CosItem, 0)
	for i := 0; i < length; i++ {
		var cir CosItemR
		if values, err = redis.Values(conn.Receive()); err != nil {
			Err(err)
			return
		} else if err = redis.ScanStruct(values, &cir); err != nil {
			Err(err)
			return
		}

		cir.Rank = start + 1 + i
		cir.Score = scores[i]

		ci := getOrigin(&cir)

		result = append(result, ci)
	}

	return
}

func FindMyCosItemsByCosplayID(username string, cosplayID int64) (result []CosItem, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	var values []interface{}
	if values, err = redis.Values(conn.Do("LRANGE", rao.GetCosItemUserListKey(cosplayID, username), 0, 2)); err != nil {
		Err(err)
		return
	}

	length := len(values)
	rankKey := rao.GetCosItemRankSetKey(cosplayID)
	for len(values) > 0 {
		var itemID int64
		if values, err = redis.Scan(values, &itemID); err != nil {
			Err(err)
			return
		}

		conn.Send("HGETALL", rao.GetCosItemSelfKey(itemID))
		conn.Send("ZSCORE", rankKey, itemID)
		conn.Send("ZREVRANK", rankKey, itemID)
	}

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	result = make([]CosItem, 0)
	for i := 0; i < length; i++ {
		var cir CosItemR
		if values, err = redis.Values(conn.Receive()); err != nil {
			Err(err)
			return
		} else if err = redis.ScanStruct(values, &cir); err != nil {
			Err(err)
			return
		}

		if cir.Score, err = redis.Int(conn.Receive()); err != nil {
			Err(err)
			return
		}

		if cir.Rank, err = redis.Int(conn.Receive()); err != nil {
			Err(err)
			return
		}
		cir.Rank++

		ci := getOrigin(&cir)

		result = append(result, ci)
	}

	return
}

func FindByPageAndCountByCosplayID(cosplayID int64, page int, pageCount int) (result []CosItem, err error) {
	return FindByPageAndCountOrderByScore(cosplayID, page, pageCount)
}

func GetUserMark(username string, cosplayId int64, itemId int64) (voted bool, score int, err error) {
	conn := rao.GetConn()
	defer conn.Close()
	key := rao.GetCosItemUserMarkListKey(cosplayId, itemId)

	conn.Send("LRANGE", key, 0, -1)
	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	var values []interface{}
	if values, err = redis.Values(conn.Receive()); err != nil {
		Err(err)
		return
	} else {
		for len(values) > 0 {
			var cs string
			if values, err = redis.Scan(values, &cs); err != nil {
				Err(err)
				return
			}

			cs = cs[1 : len(cs)-1]
			c := CosItemScore{}
			if err = json.Unmarshal([]byte(cs), &c); err != nil {
				Err(err)
				return
			}

			if c.Username == username {
				voted = true
				score = c.Score
				return
			}
		}
	}

	return
}

func GetUserMarks(cosplayId int64, itemId int64) (result []CosItemScore, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetCosItemUserMarkListKey(cosplayId, itemId)

	var values []interface{}
	if values, err = redis.Values(conn.Do("LRANGE", key, 0, -1)); err != nil {
		Err(err)
		return
	}

	result = make([]CosItemScore, 0)

	for len(values) > 0 {
		var cis string
		values, err = redis.Scan(values, &cis)
		if err != nil {
			Err(err)
			return
		}

		cis = cis[1 : len(cis)-1]
		ci := CosItemScore{}
		err = json.Unmarshal([]byte(cis), &ci)

		if err != nil {
			Err(err)
			return
		}

		result = append(result, ci)
	}

	return
}

func AddScore(cosplayId int64, itemId int64, username string, score int, totalScore int64) (err error) {
	cmt := CosItemScore{
		Username:  username,
		Score:     score,
		CosItemID: itemId,
	}
	cmtb, err := json.Marshal(cmt)
	if err != nil {
		Err(err)
		return
	}

	//add to db first
	if err = addScoreToDB(itemId, username, score, totalScore); err != nil {
		return
	}

	conn := rao.GetConn()
	defer conn.Close()

	// Info("add score.username:", username, "score:", score, "cosplay_id:", cosplayId, "item_id:", itemId, "total_score:", totalScore)
	userMarkKey := rao.GetCosItemUserMarkListKey(cosplayId, itemId)
	scoreKey := rao.GetCosItemRankSetKey(cosplayId)
	//append user mark
	conn.Send("LPUSH", userMarkKey, "'"+string(cmtb)+"'")
	//new score
	conn.Send("ZADD", scoreKey, totalScore, itemId)
	err = conn.Flush()

	if err != nil {
		Err(err)
		return
	}

	return
}

func getElement(ci *CosItem) (cir CosItemR) {
	cb, _ := json.Marshal(ci.Clothes)
	dmb, _ := json.Marshal(ci.DyeMap)
	//do not set score & rank
	return CosItemR{
		Username:   ci.Username,
		UID:        ci.UID,
		Nickname:   ci.Nickname,
		CosplayID:  ci.CosplayID,
		ItemID:     ci.ItemID,
		Title:      ci.Title,
		Desc:       ci.Desc,
		ModelNo:    ci.ModelNo,
		Clothes:    "'" + string(cb) + "'",
		SysScore:   ci.SysScore,
		Icon:       ci.Icon,
		Head:       ci.Head,
		CosBg:      ci.CosBg,
		CosPoints:  ci.CosPoints,
		UploadTime: ci.UploadTime,
		Img:        ci.Img,
		DyeMap:     "'" + string(dmb) + "'",
	}
}

func getOrigin(v *CosItemR) (ci CosItem) {
	ci = CosItem{
		Username:   v.Username,
		UID:        v.UID,
		Nickname:   v.Nickname,
		CosplayID:  v.CosplayID,
		ItemID:     v.ItemID,
		Title:      v.Title,
		Desc:       v.Desc,
		ModelNo:    v.ModelNo,
		SysScore:   v.SysScore,
		Icon:       v.Icon,
		Head:       v.Head,
		CosBg:      v.CosBg,
		CosPoints:  v.CosPoints,
		UploadTime: v.UploadTime,
		Score:      int64(v.Score),
		Rank:       v.Rank,
		Img:        v.Img,
	}

	if len(v.Clothes) > 2 {
		json.Unmarshal([]byte(v.Clothes[1:len(v.Clothes)-1]), &ci.Clothes)
	}
	if len(v.DyeMap) > 2 {
		json.Unmarshal([]byte(v.DyeMap[1:len(v.DyeMap)-1]), &ci.DyeMap)
	}

	return ci
}
