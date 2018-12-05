package items

import (
	"constants"
	"db/mail"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	. "logger"
	"rao"
	// "service/push"
	"time"
)

import (
	. "language"
	. "types"
)

//push one cos item
//rconn not nil means recover mode
//only in redis
func AddTopCosItem(item CosItem, rconn *redis.Conn) (before int, err error) {
	cir := getElement(&item)

	var cirj []byte
	if cirj, err = json.Marshal(cir); err != nil {
		Err(err)
		return
	}

	key := rao.GetTopCosItemListKey(item.CosplayID)

	if rconn == nil {
		conn := rao.GetConn()
		defer conn.Close()

		conn.Send("RPUSH", key, "'"+string(cirj)+"'")
		conn.Send("EXPIRE", key, constants.TOP_COS_EXPIRE_TIME)
		if err = conn.Flush(); err != nil {
			Err(err)
			return
		}

		if before, err = redis.Int(conn.Receive()); err != nil {
			Err(err)
			return
		}
		if before > 0 {
			before--
		}
	} else {
		(*rconn).Send("RPUSH", key, "'"+string(cirj)+"'")
	}

	return
}

func GetTopCosItem(cosplayID int64) (item CosItem, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	var cj string
	if cj, err = redis.String(conn.Do("GET", rao.GetTopCosItemKey(cosplayID))); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
		}
		return
	} else if len(cj) > 2 {
		cj = cj[1 : len(cj)-1]
		var ct CosItemR
		if err = json.Unmarshal([]byte(cj), &ct); err != nil {
			Err(err)
			return
		}

		item.Username = ct.Username
		item.UID = ct.UID
		item.Nickname = ct.Nickname
		item.CosplayID = ct.CosplayID
		item.ItemID = ct.ItemID
		item.Title = ct.Title
		item.Desc = ct.Desc
		item.ModelNo = ct.ModelNo
		if len(ct.Clothes) > 2 {
			if err = json.Unmarshal([]byte(ct.Clothes[1:len(ct.Clothes)-1]), &item.Clothes); err != nil {
				Err(err)
				return
			}
		}

		item.SysScore = ct.SysScore
		item.Icon = ct.Icon
		item.Head = ct.Head
		item.CosBg = ct.CosBg
		item.CosPoints = ct.CosPoints
		item.UploadTime = ct.UploadTime
		item.Score, item.Rank = getScoreAndRank(cosplayID, item.ItemID)
		item.Img = ct.Img

		if len(ct.DyeMap) > 2 {
			if err = json.Unmarshal([]byte(ct.DyeMap[1:len(ct.DyeMap)-1]), &item.DyeMap); err != nil {
				Err(err)
				return
			}
		}
	}

	return
}

func getScoreAndRank(cosplayID int64, itemID int64) (score int64, rank int) {
	conn := rao.GetConn()
	defer conn.Close()

	rankKey := rao.GetCosItemRankSetKey(cosplayID)
	conn.Send("ZSCORE", rankKey, itemID)
	conn.Send("ZREVRANK", rankKey, itemID)

	var err error
	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	if score, err = redis.Int64(conn.Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("can't find score for item:", itemID)
		} else {
			Err(err)
		}

	}

	if rank, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("can't find rank for item:", itemID)
		} else {
			Err(err)
		}
	}

	return
}

func CheckTopFrequently(cosplayID int64) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	listKey := rao.GetTopCosItemListKey(cosplayID)
	topKey := rao.GetTopCosItemKey(cosplayID)

	conn.Send("GET", topKey)
	conn.Send("LINDEX", listKey, 0)
	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	var top string
	var first string
	if top, err = redis.String(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
		}
	}

	if first, err = redis.String(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
		}
	}

	if top == "" && first != "" {
		conn.Send("SET", topKey, first)
		conn.Send("EXPIRE", topKey, 60)
		if err = conn.Flush(); err != nil {
			Err(err)
			return
		}
	}

	return
}

//pop one cos item
//only in redis
//only main game server call this method
func PopTopCosItem(cosplayID int64) (err error) {
	srconn := rao.GetConn()
	defer srconn.Close()

	listKey := rao.GetTopCosItemListKey(cosplayID)

	var hj string

	srconn.Send("LPOP", listKey)
	srconn.Send("EXPIRE", listKey, constants.TOP_COS_EXPIRE_TIME)
	if err = srconn.Flush(); err != nil {
		Err(err)
		return
	}

	if hj, err = redis.String(srconn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
		} else {
			err = nil
		}

		return
	}
	srconn.Receive()

	//push to key
	srconn.Send("SET", rao.GetTopCosItemKey(cosplayID), hj)
	srconn.Send("EXPIRE", rao.GetTopCosItemKey(cosplayID), 61)
	if err = srconn.Flush(); err != nil {
		Err(err)
		return
	}

	return
}

//iterate all items.send cosplay reward
func ProcessCos(cosplayID int64, cosTitle string, levelLimit int) (err error) {
	Info("cosplay id:", cosplayID, "closed.begin process..")
	defer Info("cosplay process end:", cosplayID)
	conn := rao.GetConn()
	defer conn.Close()

	scoreKey := rao.GetCosItemRankSetKey(cosplayID)
	//get item cnt
	var cnt int
	if cnt, err = redis.Int(conn.Do("ZCARD", scoreKey)); err != nil {
		Err(err)
		return
	} else {
		Info("cosplay:", cosplayID, "total item count:", cnt)
		pos1 := int(cnt / 10) //upper 10%
		if pos1 <= 0 {
			pos1 = 1
		}
		pos2 := int(cnt * 3 / 10) //upper 30%
		if pos2 <= 0 {
			pos2 = 1
		}
		pos3 := int(cnt / 2) //upper 50%
		if pos3 <= 0 {
			pos3 = 1
		}

		i := 0
		//iterate all items
		var values []interface{}
		var values1 []interface{}
		if values, err = redis.Values(conn.Do("ZREVRANGE", scoreKey, 0, -1)); err != nil {
			Err(err)
			return
		} else {
			for len(values) > 0 {
				var itemID int64
				if values, err = redis.Scan(values, &itemID); err != nil {
					Err(err)
					return
				} else {
					var ci CosItemR
					if values1, err = redis.Values(conn.Do("HGETALL", rao.GetCosItemSelfKey(itemID))); err != nil {
						Err(err)
						return
					} else if err = redis.ScanStruct(values1, &ci); err != nil {
						Err(err)
						return
					} else if ci.Username != "" {
						i++
						var diamond int64
						var diamond2 int64

						level10 := false
						level30 := false
						level50 := false

						if i <= pos1 {
							level10 = true
							//level 1.80/50 diamond
							if levelLimit == 0 {
								//normal
								diamond = 80
							} else {
								diamond = 50
							}
						} else if i <= pos2 {
							level30 = true
							//level 2.40 diamond
							diamond = 40
						} else if i <= pos3 {
							level50 = true
							//level3.30 diamond
							diamond = 30
						} else {
							//10/5 diamond
							if levelLimit == 0 {
								diamond = 10
							} else {
								diamond = 5
							}
						}

						if i == 1 {
							if levelLimit == 0 {
								diamond2 = 1000
							} else {
								diamond2 = 500
							}
						} else if i == 2 {
							if levelLimit == 0 {
								diamond2 = 500
							} else {
								diamond2 = 250
							}
						} else if i == 3 {
							if levelLimit == 0 {
								diamond2 = 300
							} else {
								diamond2 = 150
							}
						} else if i <= 10 {
							if levelLimit == 0 {
								diamond2 = 200
							} else {
								diamond2 = 100
							}
						}

						var content string
						if !level10 && !level30 && !level50 {
							content = fmt.Sprintf(L("cos1"), cosTitle)
						} else if level10 {
							content = fmt.Sprintf(L("cos2"), cosTitle, 10)
						} else if level30 {
							content = fmt.Sprintf(L("cos2"), cosTitle, 30)
						} else if level50 {
							content = fmt.Sprintf(L("cos2"), cosTitle, 50)
						}

						now := time.Now().Unix()
						m := Mail{
							To:         ci.Username,
							Title:      L("cos3"),
							Content:    content,
							Diamond:    diamond,
							Delete:     true,
							Time:       now,
							ExpireTime: now + constants.DEFAULT_MAIL_EXPIRE_TIME,
						}
						mail.SendToOne(m, nil)

						if diamond2 != 0 {
							Info("cosplay:", cosplayID, "rank:", i, "player:", ci.Username)
							r := L("cos4")
							if i == 1 {
								r = L("cos5")
							} else if i == 2 {
								r = L("cos6")
							} else if i == 3 {
								r = L("cos7")
							}
							content := fmt.Sprintf(L("cos9"), cosTitle, r)
							m := Mail{
								To:         ci.Username,
								Title:      L("cos3"),
								Content:    content,
								Diamond:    diamond2,
								Delete:     true,
								Time:       now,
								ExpireTime: now + constants.DEFAULT_MAIL_EXPIRE_TIME,
							}
							mail.SendToOne(m, nil)
						}
					}
				}
			}
		}
	}

	return
}
