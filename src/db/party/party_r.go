package party

import (
	// "appcfg"
	"constants"
	"encoding/json"
	// "errors"
	cp "common_db/party"
	"fmt"
	"github.com/garyburd/redigo/redis"
	. "logger"
	"math/rand"
	"rao"
	"strconv"
	"strings"
	"time"
	"tools"
	. "types"
)

func refreshPartyItemR(partyID int64, username, partner, partnerNickname, img, head string, partyCloseTime int64, clos []string, dyeMap map[string][7][4]float64, modelNo, qm int) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	selfKey := rao.GetPartyItemSelfKey(partyID, username)
	partySelfKey := rao.GetPartySelfKey(partyID)

	closb, _ := json.Marshal(clos)
	dmb, _ := json.Marshal(dyeMap)

	conn.Send("HSET", selfKey, "Partner", partner)
	conn.Send("HSET", selfKey, "PartnerNickname", partnerNickname)
	conn.Send("HSET", selfKey, "Img", img)
	conn.Send("HSET", selfKey, "PartnerHead", head)
	conn.Send("HSET", selfKey, "PartnerClothes", string(closb))
	conn.Send("HSET", selfKey, "PartnerModelNo", modelNo)
	conn.Send("HSET", selfKey, "PartnerDyeMap", string(dmb))
	conn.Send("HSET", selfKey, "PartnerQinmi", qm)

	//add player cnt
	conn.Send("HINCRBY", partySelfKey, "PlayerCnt", 1)

	//add party invite list
	conn.Send("LPUSH", rao.GetPartyInviteListKey(partyID), fmt.Sprintf("%s_%s", username, partner))

	//global party invite list
	gpisKey := rao.GetGlobalPartyInviteSetKey()
	conn.Send("ZREMRANGEBYSCORE", gpisKey, 0, time.Now().Unix())
	conn.Send("ZADD", gpisKey, partyCloseTime, fmt.Sprintf("%d_%s_%s", partyID, username, partner))

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}
	return
}

func GetPartyJoinerByPartners(partyID int64, us map[string]*FriendInfo) (res map[string]*FriendInfo, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	pilKey := rao.GetPartyInviteListKey(partyID)

	var values []interface{}
	if values, err = redis.Values(conn.Do("LRANGE", pilKey, 0, -1)); err != nil {
		Err(err)
		return
	}

	res = make(map[string]*FriendInfo)
	for len(values) > 0 {
		var content string
		if values, err = redis.Scan(values, &content); err != nil {
			Err(err)
			return
		}

		cs := strings.Split(content, "_")
		if len(cs) == 2 && us[cs[1]] != nil && us[cs[1]].UID > 0 {
			joiner := cs[0]
			res[joiner] = &FriendInfo{Username: joiner, UID: tools.GetUIDFromUsername(joiner)}
		}
	}

	return
}

func AddInviteFolwer(partyID int64, username string) (err error) {
	key := rao.GetPartyInviteFlowerKey(partyID)

	conn := rao.GetConn()
	defer conn.Close()

	var effect int
	if effect, err = redis.Int(conn.Do("SADD", key, username)); err != nil {
		Err(err)
		return
	} else if effect <= 0 {
		Info("already invite flower:", partyID, username)
		err = fmt.Errorf("already invite flower %d %s", partyID, username)
		return
	}
	return
}

func AlreadyInviteFlower(partyID int64, username string) bool {
	key := rao.GetPartyInviteFlowerKey(partyID)

	conn := rao.GetConn()
	defer conn.Close()

	if effect, err := redis.Int(conn.Do("SISMEMBER", key, username)); err != nil {
		Err(err)
		return true
	} else if effect <= 0 {
		return false
	}
	return true
}

func AddInviteAttend(partyID int64, username string) (err error) {
	key := rao.GetPartyInviteAttendKey(partyID)

	conn := rao.GetConn()
	defer conn.Close()

	var effect int
	if effect, err = redis.Int(conn.Do("SADD", key, username)); err != nil {
		Err(err)
		return
	} else if effect <= 0 {
		Info("already invite attend:", partyID, username)
		err = fmt.Errorf("already invite attend %d %s", partyID, username)
		return
	}
	return
}

func AlreadyInviteAttend(partyID int64, username string) bool {
	key := rao.GetPartyInviteAttendKey(partyID)

	conn := rao.GetConn()
	defer conn.Close()

	if effect, err := redis.Int(conn.Do("SISMEMBER", key, username)); err != nil {
		Err(err)
		return true
	} else if effect <= 0 {
		return false
	}
	return true
}

func ClearUserPartyCache(username string, username2 string) {
	conn := rao.GetConn()
	defer conn.Close()

	conn.Send("DEL", rao.GetPartyFriendHostCacheKey(username))
	conn.Send("DEL", rao.GetPartyFriendJoinCacheKey(username))
	conn.Send("DEL", rao.GetPartyFriendHostCacheKey(username2))
	conn.Send("DEL", rao.GetPartyFriendJoinCacheKey(username2))
	conn.Flush()
}

func getPartyFromRedis(partyID int64) (p Party, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	conn.Send("HGETALL", rao.GetPartySelfKey(partyID))
	// conn.Send("LLEN", rao.GetPartyItemListKey(partyID))

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	var values []interface{}
	if values, err = redis.Values(conn.Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no party in redis:", partyID)
		} else {
			Err(err)
		}
		return
	} else if err = redis.ScanStruct(values, &p); err != nil {
		Err(err)
		return
	}

	// if p.PlayerCnt, err = redis.Int(conn.Receive()); err != nil {
	// 	if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
	// 		Info("no party in redis:", partyID)
	// 	} else {
	// 		Err(err)
	// 	}
	// 	return
	// }
	return
}

// func getPartys1(length int, conn *redis.Conn) (res []Party, err error) {
// 	if err = (*conn).Flush(); err != nil {
// 		Err(err)
// 		return
// 	}

// 	now := time.Now().Unix()
// 	res = make([]Party, 0)
// 	for i := 0; i < length; i++ {
// 		var p Party
// 		if p, err = getParty(conn); err != nil {
// 			return
// 		}

// 		if p.CloseTime > now {
// 			res = append(res, p)
// 		}
// 	}
// 	return
// }

func getPartyListByIdMapFromRedis(ids map[int64]int, conn *redis.Conn, nowt int64, checkCloseTime, excludeCos bool) (res []Party, err error) {
	if conn == nil {
		c := rao.GetConn()
		defer c.Close()
		conn = &c
	}

	for id, _ := range ids {
		(*conn).Send("HGETALL", rao.GetPartySelfKey(id))
		(*conn).Send("LLEN", rao.GetPartyItemListKey(id))
	}

	return cp.GetPartys1(len(ids), conn, nowt, checkCloseTime, excludeCos)
}

// func getPartyListByIdsFromRedis(ids map[int64]int, conn *redis.Conn) (res []Party, err error) {
// 	if conn == nil {
// 		c := rao.GetConn()
// 		defer c.Close()
// 		conn = &c
// 	}

// 	for id, _ := range ids {
// 		(*conn).Send("HGETALL", rao.GetPartySelfKey(id))
// 		(*conn).Send("LLEN", rao.GetPartyItemListKey(id))
// 	}

// 	return getPartys1(len(ids), conn)
// }

func getPartyListByIdsFromRedis1(ids []int64, conn *redis.Conn, nowt int64, checkCloseTime, excludeCos bool) (res []Party, err error) {
	if conn == nil {
		c := rao.GetConn()
		defer c.Close()
		conn = &c
	}

	for _, id := range ids {
		(*conn).Send("HGETALL", rao.GetPartySelfKey(id))
		(*conn).Send("LLEN", rao.GetPartyItemListKey(id))
	}

	return cp.GetPartys1(len(ids), conn, nowt, checkCloseTime, excludeCos)
}

// func getPartyListFromRedis(typee int, pageNo int, pageSize int) (res []Party, err error) {
// 	conn := rao.GetConn()
// 	defer conn.Close()

// 	start := (pageNo - 1) * pageSize
// 	end := pageNo*pageSize - 1

// 	if typee == constants.PARTY_CASUAL_TYPE {
// 		if pageNo > 1 {
// 			return
// 		}
// 		conn.Send("ZREVRANGE", rao.GetGlobalCasualPartyListKey(), 0, -1)
// 	} else if typee == constants.PARTY_PRIZE_TYPE {
// 		conn.Send("ZREVRANGE", rao.GetGlobalPrizePartyListKey(), start, end)
// 	} else if typee == constants.PARTY_COSPLAY_TYPE {
// 		conn.Send("ZREVRANGE", rao.GetGlobalCosplayPartyListKey(), 0, 6)
// 	}

// 	if err = conn.Flush(); err != nil {
// 		Err(err)
// 		return
// 	}

// 	var ids map[int64]int
// 	if ids, err = getPartyIDs(typee, &conn, nil); err != nil || len(ids) <= 0 {
// 		return
// 	}

// 	return getPartyListByIdsFromRedis(ids, &conn)
// }

func getPartyItemRankListFromRedis(partyID int64, size int, qu string) (res []PartyItem, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	rankListKey := rao.GetPartyItemRankListKey(partyID)

	var values []interface{}
	if values, err = redis.Values(conn.Do("ZREVRANGE", rankListKey, 0, size-1)); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no rank list:", rankListKey)
			err = nil
		} else {
			Err(err)
		}
		return
	} else {
		return cp.GetPartyItemsFromValues(values, &conn, partyID, qu, false)
	}

	return
}

func getPartyUsersItemsFromRedis(partyID int64, us map[string]*FriendInfo) (r []PartyItem, err error) {
	if len(us) <= 0 {
		return
	}

	conn := rao.GetConn()
	defer conn.Close()

	for u, _ := range us {
		conn.Send("HGETALL", rao.GetPartyItemSelfKey(partyID, u))
		conn.Send("ZREVRANK", rao.GetPartyItemRankListKey(partyID), u)
	}

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	r = make([]PartyItem, 0)
	var values []interface{}
	var rank int
	for i := 0; i < len(us); i++ {
		var p PartyItem
		if values, err = redis.Values(conn.Receive()); err != nil {
			if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
				err = nil
				conn.Receive()
				continue
			} else {
				Err(err)
				return
			}
		}

		if rank, err = redis.Int(conn.Receive()); err != nil {
			if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
				err = nil
				continue
			} else {
				Err(err)
				return
			}
		}

		var pr PartyItemRaw
		if err = redis.ScanStruct(values, &pr); err != nil {
			Err(err)
			return
		} else {
			p = cp.GetFromPartyItemRaw(&pr)
			p.Rank = rank + 1
			r = append(r, p)
		}
	}

	return
}

func getPartyItemListFromRedis(partyID int64, pageNo int, pageSize int, queryUsername string) (res []PartyItem, err error) {
	conn := rao.GetConn()
	defer conn.Close()
	return cp.GetPartyItemListFromRedis(partyID, pageNo, pageSize, queryUsername, &conn)
}

func getRandomPartyItemsFromRedis(username string) (ps []PartyItem, err error) {
	ps = make([]PartyItem, 0)

	conn := rao.GetConn()
	defer conn.Close()

	conn.Send("ZREMRANGEBYSCORE", constants.PARTY_ITEM_GLOBAL_SET, 0, time.Now().Unix())
	conn.Send("ZCARD", constants.PARTY_ITEM_GLOBAL_SET)

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	conn.Receive()

	var size int
	if size, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no random item")
		} else {
			Err(err)
		}
		return
	} else if size > 0 {
		var values []interface{}
		m := make(map[int]int)
		for i := 0; i < 10; i++ {
			index := rand.Intn(size)
			if m[index] == 0 {
				m[index] = 1
				if values, err = redis.Values(conn.Do("ZRANGE", constants.PARTY_ITEM_GLOBAL_SET, index, index)); err != nil {
					if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
						Info("no random item,size:", size, "random index:", index)
						continue
					} else {
						Err(err)
						return
					}
				} else {
					var token string
					if values, err = redis.Scan(values, &token); err != nil {
						Err(err)
						return
					} else {
						ts := strings.Split(token, "_")
						pid, _ := strconv.ParseInt(ts[0], 10, 64)
						u := ts[1]

						if u != username {
							var p PartyItem
							if p, err = getPartyItemForPartyUserFromRedis(pid, u, &conn, false, ""); err != nil {
								return
							} else if p.Partner != username && p.ID > 0 {
								ps = append(ps, p)
							}
						}

					}
				}
			}

		}
	}

	return
}

func getPartyItemForPartyUserFromRedis(partyID int64, username string, conn *redis.Conn, needRank bool, qUsername string) (p PartyItem, err error) {
	if conn == nil {
		c := rao.GetConn()
		defer c.Close()
		conn = &c
	}

	selfFlowerKey := rao.GetPartyItemFlowerSetKey(partyID, username)

	(*conn).Send("HGETALL", rao.GetPartyItemSelfKey(partyID, username))
	if needRank {
		(*conn).Send("ZREVRANK", rao.GetPartyItemRankListKey(partyID), username)
	}
	if qUsername != "" {
		(*conn).Send("SISMEMBER", selfFlowerKey, qUsername)
	}

	if err = (*conn).Flush(); err != nil {
		Err(err)
		return
	}

	var pr PartyItemRaw
	var values []interface{}
	if values, err = redis.Values((*conn).Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no party item for:", partyID, username)
		} else {
			Err(err)
		}
		return
	} else if err = redis.ScanStruct(values, &pr); err != nil {
		Err(err)
		return
	} else {
		p = cp.GetFromPartyItemRaw(&pr)
	}

	if needRank {
		var rank int
		if rank, err = redis.Int((*conn).Receive()); err != nil {
			if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
				Info("no party item rank for:", partyID, username)
			} else {
				Err(err)
			}
			return
		} else {
			p.Rank = rank
		}
	}

	if qUsername != "" {
		var send int
		if send, err = redis.Int((*conn).Receive()); err != nil {
			if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
				Info("no flower info for getPartyItem")
			} else {
				Err(err)
			}
			return
		}

		if send == 1 {
			p.HasSendFlower = true
		}
	}

	return
}

// func getPartyItemCmtsFromRedis(partyID int64, username string, pageNo int, pageSize int) (res []PartyComment, err error) {
// 	conn := rao.GetConn()
// 	defer conn.Close()

// 	start := (pageNo - 1) * pageSize
// 	end := pageNo*pageSize - 1
// 	var values []interface{}
// 	if values, err = redis.Values(conn.Do("LRANGE", rao.GetPartyItemCmtListKey(partyID, username), start, end)); err != nil {
// 		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
// 			Info("no party item cmts for:", partyID, username)
// 		} else {
// 			Err(err)
// 		}
// 		return
// 	}

// 	res = make([]PartyComment, 0)
// 	for len(values) > 0 {
// 		var ps string
// 		if values, err = redis.Scan(values, &ps); err != nil {
// 			Err(err)
// 			return
// 		}

// 		ps = ps[1 : len(ps)-1]

// 		var p PartyComment
// 		if err = json.Unmarshal([]byte(ps), &p); err != nil {
// 			Err(err)
// 			return
// 		}

// 		res = append(res, p)
// 	}
// 	return
// }

// func addPartyItemCmtToRedis(id int64, partyID int64, username string, sendUsername string, sendNickname string, content string, t int64, partyCloseTime int64) (pc PartyComment, err error) {
// 	conn := rao.GetConn()
// 	defer conn.Close()

// 	pc = PartyComment{
// 		ID:           id,
// 		SendUsername: sendUsername,
// 		SendNickname: sendNickname,
// 		Content:      content,
// 		Time:         t,
// 	}

// 	var pb []byte
// 	if pb, err = json.Marshal(pc); err != nil {
// 		Err(err)
// 		return
// 	}

// 	listKey := rao.GetPartyItemCmtListKey(partyID, username)
// 	conn.Send("LPUSH", listKey, "'"+string(pb)+"'")
// 	conn.Send("LTRIM", listKey, 0, constants.COS_ITEM_CMT_LIST_SIZE-1)

// 	if err = conn.Flush(); err != nil {
// 		Err(err)
// 		return
// 	}

// 	return
// }

func resetPartyItemPopuToRedis(partyID int64, username string, likeCnt int, unlikeCnt int, partyType int) (popu int, lc int, ulc int, fc int, rank int, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	selfKey := rao.GetPartyItemSelfKey(partyID, username)

	conn.Send("HINCRBY", selfKey, "LikeCnt", likeCnt)
	conn.Send("HINCRBY", selfKey, "UnlikeCnt", unlikeCnt)
	conn.Send("HGET", selfKey, "FlowerCnt")
	conn.Send("HGET", selfKey, "Partner")
	conn.Send("HGET", selfKey, "PartnerQinmi")

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	if popu, lc, ulc, fc, err = getPartyItemPopuInRedis(partyID, username, &conn, partyType); err != nil {
		return
	}

	if rank, err = setPartyItemPopu(selfKey, rao.GetPartyItemRankListKey(partyID), username, int(popu), &conn); err != nil {
		return
	}

	return
}

func addPartyItemFlowerToRedis(partyID int64, username string, sendUsername string, partyCloseTime int64, partyType int) (popu int, lc int, ulc int, fc int, rank int, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	selfFlowerKey := rao.GetPartyItemFlowerSetKey(partyID, username)
	selfKey := rao.GetPartyItemSelfKey(partyID, username)

	conn.Send("HGET", selfKey, "LikeCnt")
	conn.Send("HGET", selfKey, "UnlikeCnt")
	conn.Send("HINCRBY", selfKey, "FlowerCnt", 1)
	conn.Send("HGET", selfKey, "Partner")
	conn.Send("HGET", selfKey, "PartnerQinmi")
	conn.Send("SADD", selfFlowerKey, sendUsername)

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	if popu, lc, ulc, fc, err = getPartyItemPopuInRedis(partyID, username, &conn, partyType); err != nil {
		return
	}

	conn.Receive()

	if rank, err = setPartyItemPopu(selfKey, rao.GetPartyItemRankListKey(partyID), username, int(popu), &conn); err != nil {
		return
	}

	return
}

func getPartysByHostUsers(us map[string]*FriendInfo, conn *redis.Conn, nowt int64, checkCloseTime, excludeCos bool) (res []Party, err error) {
	// for i := 0; i < 24; i++ {
	// (*conn).Send("ZRANGE", rao.GetCasultPartyListKey(-i), 0, -1)
	(*conn).Send("ZRANGE", rao.GetGlobalCasualPartyListKey(), 0, -1)
	// }

	// lkey := rao.GetPartyListKey(constants.PARTY_PRIZE_TYPE)
	(*conn).Send("ZRANGE", rao.GetGlobalPrizePartyListKey(), 0, -1)

	if err = (*conn).Flush(); err != nil {
		Err(err)
		return
	}

	var idMap map[int64]int
	var idList []int64
	if idMap, idList, err = cp.GetPartyIDs(-1, conn, us); err != nil || len(idMap) <= 0 {
		return
	}

	return cp.GetPartyListByIdsFromRedis(idList, conn, nowt, checkCloseTime, excludeCos)
}

func getPartysByJoinUsers(us map[string]*FriendInfo, conn *redis.Conn, nowt int64, checkCloseTime, excludeCos bool) (res []Party, err error) {
	(*conn).Send("ZREVRANGE", constants.PARTY_ITEM_GLOBAL_SET, 0, -1)
	(*conn).Send("ZREVRANGE", rao.GetGlobalPartyInviteSetKey(), 0, -1)

	if err = (*conn).Flush(); err != nil {
		Err(err)
		return
	}

	idMap := make(map[int64]int)
	var values []interface{}

	if values, err = redis.Values((*conn).Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no info for:", constants.PARTY_ITEM_GLOBAL_SET)
			err = nil
		} else {
			Err(err)
			return
		}
	} else {
		for len(values) > 0 {
			var token string
			if values, err = redis.Scan(values, &token); err != nil {
				Err(err)
				return
			}

			ts := strings.Split(token, "_")
			if len(ts) == 2 {
				id, _ := strconv.ParseInt(ts[0], 10, 64)
				username := ts[1]

				if us != nil && us[username] != nil && us[username].UID > 0 {
					//legal idnp
					idMap[id] = 1
				}
			}
		}
	}

	if values, err = redis.Values((*conn).Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no info for:", constants.PARTY_ITEM_GLOBAL_SET)
			err = nil
		} else {
			Err(err)
			return
		}
	} else {
		for len(values) > 0 {
			var token string
			if values, err = redis.Scan(values, &token); err != nil {
				Err(err)
				return
			}

			ts := strings.Split(token, "_")
			if len(ts) == 3 {
				id, _ := strconv.ParseInt(ts[0], 10, 64)
				// joiner := ts[1]
				partner := ts[2]

				if us != nil && us[partner] != nil && us[partner].UID > 0 {
					//legal idnp
					idMap[id] = 1
				}
			}

		}
	}

	if len(idMap) > 0 {
		return getPartyListByIdMapFromRedis(idMap, conn, nowt, checkCloseTime, excludeCos)
	}
	return
}

func addPartyItemToRedis(
	partyID int64,
	partyItemID int64,
	username string,
	uid int64,
	nickname string,
	head string,
	img string,
	uploadTime int64,
	partyCloseTime int64,
	subject string,
	desc string,
	bgBannerType int,
	partyUsername string,
	partyNickname string,
	partnerType int,
	allowPartners []UserInfo,
	clothes []string,
	modelNo int,
	dyeMap map[string][7][4]float64,
	bannerFile, bgFile string) (p PartyItem, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	cb, _ := json.Marshal(clothes)

	as := make([]string, 0)
	for _, u := range allowPartners {
		as = append(as, u.Username)
	}
	ab, _ := json.Marshal(allowPartners)

	db, _ := json.Marshal(dyeMap)

	pr := PartyItemRaw{
		ID:             partyItemID,
		PartyID:        partyID,
		Username:       username,
		UID:            uid,
		Nickname:       nickname,
		Head:           head,
		Img:            img,
		PartyUsername:  partyUsername,
		PartyNickname:  partyNickname,
		Subject:        subject,
		Desc:           desc,
		BgBannerType:   bgBannerType,
		BannerFile:     bannerFile,
		BackgroundFile: bgFile,
		UploadTime:     uploadTime,
		PartnerType:    partnerType,
		AllowPartners:  string(ab),
		Clothes:        string(cb),
		ModelNo:        modelNo,
		Popularity:     100,
		DyeMap:         string(db),
	}

	p = PartyItem{
		ID:             partyItemID,
		PartyID:        partyID,
		Username:       username,
		UID:            uid,
		Nickname:       nickname,
		Head:           head,
		Img:            img,
		PartyUsername:  partyUsername,
		PartyNickname:  partyNickname,
		Subject:        subject,
		Desc:           desc,
		BgBannerType:   bgBannerType,
		BannerFile:     bannerFile,
		BackgroundFile: bgFile,
		UploadTime:     uploadTime,
		PartnerType:    partnerType,
		AllowPartners:  as,
		Clothes:        clothes,
		ModelNo:        modelNo,
		Popularity:     100,
		DyeMap:         dyeMap,
	}

	listKey := rao.GetPartyItemListKey(partyID)
	rankListKey := rao.GetPartyItemRankListKey(partyID)
	selfKey := rao.GetPartyItemSelfKey(partyID, username)
	partySelfKey := rao.GetPartySelfKey(partyID)

	//global party item
	conn.Send("ZADD", constants.PARTY_ITEM_GLOBAL_SET, partyCloseTime, fmt.Sprintf("%d_%s", partyID, username))

	//party item list
	conn.Send("LPUSH", listKey, username)

	//party item rank list
	conn.Send("ZADD", rankListKey, constants.DEFAULT_PARTY_ITEM_POPU, username)

	//party item itself
	conn.Send("HMSET", redis.Args{}.Add(selfKey).AddFlat(&pr)...)

	//add player cnt
	conn.Send("HINCRBY", partySelfKey, "PlayerCnt", 1)

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	return
}

// func addPartyToRedis(selfKey string, p *Party, conn *redis.Conn, typee int) (err error) {
// 	token := fmt.Sprintf("%d_%s", p.ID, p.Username)

// 	var lk string
// 	if typee == constants.PARTY_CASUAL_TYPE {
// 		lk = rao.GetGlobalCasualPartyListKey()
// 	} else {
// 		lk = rao.GetGlobalPrizePartyListKey()
// 	}
// 	(*conn).Send("ZREMRANGEBYSCORE", lk, 0, time.Now().Unix())
// 	(*conn).Send("ZADD", lk, p.CloseTime, token)
// 	(*conn).Send("HMSET", redis.Args{}.Add(selfKey).AddFlat(p)...)

// 	if err = (*conn).Flush(); err != nil {
// 		Err(err)
// 		return
// 	}

// 	return
// }

// func tryAddParty(typee int, conn *redis.Conn, cntKey string) bool {
// 	if typee == constants.PARTY_CASUAL_TYPE {
// 		(*conn).Send("ZREMRANGEBYSCORE", rao.GetGlobalCasualPartyListKey(), 0, time.Now().Unix())
// 		(*conn).Send("INCR", cntKey)
// 		(*conn).Send("EXPIRE", cntKey, 3601)
// 	} else {
// 		plk := rao.GetGlobalPrizePartyListKey()
// 		(*conn).Send("ZREMRANGEBYSCORE", plk, 0, time.Now().Unix())
// 		(*conn).Send("ZCARD", plk)
// 	}

// 	if err := (*conn).Flush(); err != nil {
// 		Err(err)
// 		return false
// 	}

// 	(*conn).Receive()

// 	if size, err := redis.Int((*conn).Receive()); err != nil {
// 		Err(err)
// 		return false
// 	} else {
// 		if typee == constants.PARTY_CASUAL_TYPE {
// 			(*conn).Receive()

// 			if size <= constants.CASUAL_PARTY_HOUR_SIZE {
// 				return true
// 			}
// 		} else if typee == constants.PARTY_PRIZE_TYPE && size < appcfg.GetInt("prize_party_max_size", 1000) {
// 			return true
// 		}
// 	}

// 	return false
// }

//support function
// func getParty(conn *redis.Conn) (p Party, err error) {
// 	var values []interface{}

// 	if values, err = redis.Values((*conn).Receive()); err != nil {
// 		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
// 			Info("no info for getParty")
// 		} else {
// 			Err(err)
// 		}
// 		return
// 	} else if err = redis.ScanStruct(values, &p); err != nil {
// 		Err(err)
// 		return
// 	}

// 	if p.JoinerCnt, err = redis.Int((*conn).Receive()); err != nil {
// 		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
// 			Info("no joiner cnt for getParty")
// 		} else {
// 			Err(err)
// 		}
// 		return
// 	}

// 	return
// }

func getPartyItem(conn *redis.Conn, queryFlowerSend bool, i int, needRank bool) (p PartyItem, err error) {
	var values []interface{}

	var pr PartyItemRaw
	if values, err = redis.Values((*conn).Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no self info for getPartyItem")
		} else {
			Err(err)
		}
		return
	} else if err = redis.ScanStruct(values, &pr); err != nil {
		Err(err)
		return
	} else {
		p = getFromPartyItemRaw(&pr)
	}

	if queryFlowerSend {
		var send int
		if send, err = redis.Int((*conn).Receive()); err != nil {
			if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
				Info("no flower info for getPartyItem")
			} else {
				Err(err)
			}
			return
		}

		if send == 1 {
			p.HasSendFlower = true
		}
	}

	if needRank {
		var rank int
		if rank, err = redis.Int((*conn).Receive()); err != nil {
			if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
				Info("no rank info for getPartyItem")
			} else {
				Err(err)
			}
			return
		}
		p.Rank = rank + 1
	} else {
		p.Rank = i + 1
	}

	return
}

// func getPartyIDs(typee int, conn *redis.Conn, us map[string]*FriendInfo) (ids map[int64]int, err error) {
// 	ids = make(map[int64]int)
// 	var values []interface{}

// 	cnt := 1
// 	if typee == -1 {
// 		cnt = 2
// 	}

// 	for i := 0; i < cnt; i++ {
// 		if values, err = redis.Values((*conn).Receive()); err != nil {
// 			if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
// 				Info("no party type:", typee)
// 			} else {
// 				Err(err)
// 			}
// 			return
// 		} else {
// 			for len(values) > 0 {
// 				var token string
// 				if values, err = redis.Scan(values, &token); err != nil {
// 					Err(err)
// 					return
// 				}
// 				ts := strings.Split(token, "_")

// 				username := ts[1]
// 				id, _ := strconv.ParseInt(ts[0], 10, 64)

// 				if us != nil {
// 					//must be friends
// 					f := false
// 					if us[username] != nil && us[username].UID > 0 {
// 						f = true
// 					}

// 					if !f {
// 						continue
// 					}
// 				}

// 				ids[id] = 1
// 			}
// 		}
// 	}
// 	return
// }

func getPartyItemsFromValues(values []interface{}, conn *redis.Conn, partyID int64, queryUsername string, needRank bool) (res []PartyItem, err error) {
	size := len(values)

	if size <= 0 {
		return
	}

	for len(values) > 0 {
		var username string
		if values, err = redis.Scan(values, &username); err != nil {
			Err(err)
			return
		} else {
			(*conn).Send("HGETALL", rao.GetPartyItemSelfKey(partyID, username))
			if queryUsername != "" {
				(*conn).Send("SISMEMBER", rao.GetPartyItemFlowerSetKey(partyID, username), queryUsername)
			}
			if needRank {
				(*conn).Send("ZREVRANK", rao.GetPartyItemRankListKey(partyID), username)
			}
		}
	}

	if err = (*conn).Flush(); err != nil {
		Err(err)
		return
	}

	res = make([]PartyItem, 0)
	for i := 0; i < size; i++ {
		var p PartyItem
		if p, err = getPartyItem(conn, queryUsername != "", i, needRank); err != nil {
			return
		} else {
			res = append(res, p)
		}
	}

	return
}

func getPartyItemPopuInRedis(partyID int64, username string, conn *redis.Conn, partyType int) (popu int, lc int, ulc int, fc int, err error) {
	if lc, err = redis.Int((*conn).Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no like cnt:", partyID, username)
		} else {
			Err(err)
		}
		return
	}

	if ulc, err = redis.Int((*conn).Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no unlike cnt:", partyID, username)
		} else {
			Err(err)
		}
		return
	}

	if fc, err = redis.Int((*conn).Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no flower cnt:", partyID, username)
		} else {
			Err(err)
		}
		return
	}

	// if partyType == constants.PARTY_CASUAL_TYPE {
	// 	if fc > constants.PARTY_FLOWER_MAX_CNT {
	// 		fc = constants.PARTY_FLOWER_MAX_CNT
	// 	}
	// } else {
	// 	if fc > constants.PRIZE_PARTY_FLOWER_MAX_CNT {
	// 		fc = constants.PRIZE_PARTY_FLOWER_MAX_CNT
	// 	}
	// }

	if fc > constants.PARTY_FLOWER_MAX {
		fc = constants.PARTY_FLOWER_MAX
	}

	var partner string
	var pqm int
	if partner, err = redis.String((*conn).Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		}
	}

	if pqm, err = redis.Int((*conn).Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		}
	}

	popu = constants.DEFAULT_PARTY_ITEM_POPU + lc*5 + fc*5 - ulc*5
	if popu < 0 {
		popu = 0
	}

	if partner != "" {
		if pqm < 60 {
			popu = int(popu * 12 / 10)
		} else {
			popu = int(popu * 14 / 10)
		}
	}

	return
}

func setPartyItemPopu(selfKey string, rankListKey string, username string, popu int, conn *redis.Conn) (rank int, err error) {
	(*conn).Send("HSET", selfKey, "Popularity", popu)
	(*conn).Send("ZADD", rankListKey, popu, username)
	(*conn).Send("ZREVRANK", rankListKey, username)

	if err = (*conn).Flush(); err != nil {
		Err(err)
		return
	}

	(*conn).Receive()
	(*conn).Receive()

	if rank, err = redis.Int((*conn).Receive()); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no key:", rankListKey)
		} else {
			Err(err)
		}
		return
	}
	rank += 1
	return
}

func getFromPartyItemRaw(pr *PartyItemRaw) (p PartyItem) {
	p.ID = pr.ID
	p.PartyID = pr.PartyID
	p.Username = pr.Username
	p.UID = pr.UID
	p.Nickname = pr.Nickname
	p.Head = pr.Head
	p.Img = pr.Img
	p.PartyUsername = pr.PartyUsername
	p.PartyNickname = pr.PartyNickname
	p.Subject = pr.Subject
	p.Desc = pr.Desc
	p.BgBannerType = pr.BgBannerType
	p.BannerFile = pr.BannerFile
	p.BackgroundFile = pr.BackgroundFile

	json.Unmarshal([]byte(pr.Clothes), &p.Clothes)
	p.ModelNo = pr.ModelNo
	p.PartnerType = pr.PartnerType //0:no 1:all 2:friend

	p.AllowPartners = make([]string, 0)
	tl := make([]UserInfo, 0)
	json.Unmarshal([]byte(pr.AllowPartners), &tl)
	for _, u := range tl {
		p.AllowPartners = append(p.AllowPartners, u.Username)
	}

	p.Partner = pr.Partner
	p.PartnerQinmi = pr.PartnerQinmi
	p.PartnerNickname = pr.PartnerNickname
	p.PartnerHead = pr.PartnerHead
	json.Unmarshal([]byte(pr.DyeMap), &p.DyeMap)

	p.PartnerModelNo = pr.PartnerModelNo
	json.Unmarshal([]byte(pr.PartnerClothes), &p.PartnerClothes)
	json.Unmarshal([]byte(pr.PartnerDyeMap), &p.PartnerDyeMap)

	p.Popularity = pr.Popularity
	p.LikeCnt = pr.LikeCnt
	p.UnlikeCnt = pr.UnlikeCnt
	p.FlowerCnt = pr.FlowerCnt

	// if p.FlowerCnt > constants.PARTY_FLOWER_MAX_CNT {
	// 	p.FlowerCnt = constants.PARTY_FLOWER_MAX_CNT
	// }

	p.UploadTime = pr.UploadTime
	p.Rank = pr.Rank

	p.HasSendFlower = pr.HasSendFlower

	// Info(p.Clothes, "^^^", p.AllowPartners)
	return
}
