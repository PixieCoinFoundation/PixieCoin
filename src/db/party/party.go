package party

import (
	"constants"
	"encoding/json"
	// "errors"
	cp "common_db/party"
	"dao_party"
	"github.com/garyburd/redigo/redis"
	. "logger"
	"rao"
	"time"
	. "types"
)

func AddParty(typee int, singleType int, startTime int64, closeTime int64, username string, nickname string, subject string, desc string, moneyType int, startPool int, ticket int, bgBannerType int) (id int64, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	// cntKey := rao.GetPartyCasualCntKey()
	// if tryAddParty(typee, &conn, cntKey) {
	// 	if id, err = addPartyToDB(typee, singleType, startTime, closeTime, username, nickname, subject, desc, moneyType, startPool, ticket, bgBannerType); err != nil {
	// 		return
	// 	}

	// 	p := Party{
	// 		ID:           id,
	// 		Type:         typee,
	// 		SingleType:   singleType,
	// 		StartTime:    startTime,
	// 		CloseTime:    closeTime,
	// 		Username:     username,
	// 		Nickname:     nickname,
	// 		Subject:      subject,
	// 		Desc:         desc,
	// 		MoneyType:    moneyType,
	// 		StartPool:    startPool,
	// 		Ticket:       ticket,
	// 		PlayerCnt:    0,
	// 		BgBannerType: bgBannerType,
	// 	}

	// 	if err = addPartyToRedis(rao.GetPartySelfKey(id), &p, &conn, typee); err != nil {
	// 		go DeletePartyInDB(id, true)
	// 		return
	// 	}
	// } else {
	// 	return 0, constants.PartySizeErr
	// }

	return cp.AddParty(typee, singleType, startTime, closeTime, username, nickname, subject, desc, moneyType, startPool, ticket, bgBannerType, "", "", &conn, dao_party.AddPartyStmt, dao_party.DeletePartyStmt, dao_party.DeletePartyItemStmt, dao_party.DeletePartyFlowerStmt, dao_party.DeletePartyCmtStmt)
}

func AddPartyItem(
	partyID int64,
	username string,
	uid int64,
	nickname, head, img string,
	partyCloseTime int64,
	subject, desc string,
	bgBannerType int,
	partyUsername, partyNickname string,
	partnerType int,
	allowPartners []UserInfo,
	clothes []string,
	modelNo int,
	dyeMap map[string][7][4]float64,
	bannerFile, bgFile string) (p PartyItem, err error) {

	t := time.Now().Unix()
	var id int64
	if id, err = addPartyItemToDB(partyID, username, nickname, head, img, t); err != nil {
		return
	}

	return addPartyItemToRedis(partyID, id, username, uid, nickname, head, img, t, partyCloseTime, subject, desc, bgBannerType, partyUsername, partyNickname, partnerType, allowPartners, clothes, modelNo, dyeMap, bannerFile, bgFile)
}

func AddPartyItemCmt(partyID int64, username string, content string, sendUsername string, sendNickname string, partyCloseTime int64) (pc PartyComment, err error) {
	t := time.Now().Unix()
	var id int64
	if id, err = addPartyCmtToDB(partyID, username, sendUsername, sendNickname, content, t); err != nil {
		return
	}

	pc = PartyComment{
		ID:           id,
		SendUsername: sendUsername,
		SendNickname: sendNickname,
		Content:      content,
		Time:         t,
	}

	return
}

func GetPartyItemCmts(partyID int64, username string, pageNo int, pageSize int, asc bool) ([]PartyComment, error) {
	return getPartyCmtFromDB(partyID, username, (pageNo-1)*pageSize, pageSize, asc)
}

func GetPartyItems(partyID int64, pageNo int, pageSize int, queryUsername string) ([]PartyItem, error) {
	return getPartyItemListFromRedis(partyID, pageNo, pageSize, queryUsername)
}

func GetPartyUsersItems(partyID int64, usernames map[string]*FriendInfo) ([]PartyItem, error) {
	return getPartyUsersItemsFromRedis(partyID, usernames)
}

func GetPartyRankItems(partyID int64, size int, qu string) ([]PartyItem, error) {
	return getPartyItemRankListFromRedis(partyID, size, qu)
}

func GetPartys(typee int, pageNo int, pageSize int, nowt int64) (res []Party, err error) {
	conn := rao.GetConn()
	defer conn.Close()
	return cp.GetPartyListFromRedis(typee, pageNo, pageSize, &conn, false, nowt, typee != constants.PARTY_COSPLAY_TYPE)
}

func GetPartysByIds(ids []int64, nowt int64, checkCloseTime, excludeCos bool) (res []Party, err error) {
	if len(ids) > 0 {
		return getPartyListByIdsFromRedis1(ids, nil, nowt, checkCloseTime, excludeCos)
	}
	return
}

func GetPartysByHostUsers(fs map[string]*FriendInfo, username string, pageNo, pageSize int, nowt int64, checkCloseTime, excludeCos bool) (res []Party, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	start := (pageNo - 1) * pageSize
	end := pageNo*pageSize - 1

	var exist bool
	var res1 []Party
	cacheKey := rao.GetPartyFriendHostCacheKey(username)
	if exist, err = rao.ExistKey(cacheKey, &conn, 0); err != nil {
		Err(err)
		return
	} else if exist {
		var ps string
		if ps, err = redis.String(conn.Do("GET", cacheKey)); err != nil {
			if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
				Info("no cache:", cacheKey)
			} else {
				Err(err)
			}
			return
		}
		ps = ps[1 : len(ps)-1]

		var res2 []*Party
		if err = json.Unmarshal([]byte(ps), &res2); err != nil {
			Err(err)
			return
		}

		for _, p := range res2 {
			if fs != nil && fs[p.Username] != nil && fs[p.Username].UID > 0 {
				res1 = append(res1, *p)
			}
		}
	} else {
		if res1, err = getPartysByHostUsers(fs, &conn, nowt, checkCloseTime, excludeCos); err != nil {
			return
		}

		var psb []byte
		if psb, err = json.Marshal(res1); err != nil {
			Err(err)
			return
		}

		conn.Send("SET", cacheKey, "'"+string(psb)+"'")
		conn.Send("EXPIRE", cacheKey, constants.FRIEND_PARTY_EXPIRE)
		if err = conn.Flush(); err != nil {
			Err(err)
			return
		}
	}

	if end > len(res1) {
		end = len(res1)
	}

	if start > len(res1)-1 || start >= end {
		return
	} else {
		res = res1[start:end]
		return
	}
}

func GetPartysByJoinUsers(fs map[string]*FriendInfo, username string, pageNo, pageSize int, nowt int64, checkCloseTime, excludeCos bool) (res []Party, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	start := (pageNo - 1) * pageSize
	end := pageNo*pageSize - 1

	var exist bool
	var res1 []Party
	cacheKey := rao.GetPartyFriendJoinCacheKey(username)
	if exist, err = rao.ExistKey(cacheKey, &conn, 0); err != nil {
		Err(err)
		return
	} else if exist {
		var ps string
		if ps, err = redis.String(conn.Do("GET", cacheKey)); err != nil {
			if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
				Info("no cache:", cacheKey)
			} else {
				Err(err)
			}
			return
		}
		ps = ps[1 : len(ps)-1]

		if err = json.Unmarshal([]byte(ps), &res1); err != nil {
			Err(err)
			return
		}
	} else {
		if res1, err = getPartysByJoinUsers(fs, &conn, nowt, checkCloseTime, excludeCos); err != nil {
			return
		}

		var psb []byte
		if psb, err = json.Marshal(res1); err != nil {
			Err(err)
			return
		}

		conn.Send("SET", cacheKey, "'"+string(psb)+"'")
		conn.Send("EXPIRE", cacheKey, constants.FRIEND_PARTY_EXPIRE)
		if err = conn.Flush(); err != nil {
			Err(err)
			return
		}
	}

	if end > len(res1) {
		end = len(res1)
	}

	if start > len(res1)-1 || start >= end {
		return
	} else {
		res = res1[start:end]
		return
	}
}

func GetPartyByID(partyID int64) (p Party, err error) {
	return getPartyFromRedis(partyID)
}

func GetPartyItem(partyID int64, username string, needRank bool, qUsername string) (PartyItem, error) {
	return getPartyItemForPartyUserFromRedis(partyID, username, nil, needRank, qUsername)
}

func GetRandomPartyItems(username string) ([]PartyItem, error) {
	return getRandomPartyItemsFromRedis(username)
}

func GetRandomPartyItemPairs(username string, nowt int64) ([]PartyItemPair, error) {
	return getRandomPartyItemPairsFromRedis(username, nowt)
}

func RefreshPartyItem(partyID int64, username, partner, partnerNickname, img, head string, pct int64, clos []string, dyeMap map[string][7][4]float64, modelNo int, qm int, partyType int) (success bool, err error) {
	if success, _ = refreshPartyItem(partyID, username, partner, partnerNickname, qm); success {
		refreshPartyItemR(partyID, username, partner, partnerNickname, img, head, pct, clos, dyeMap, modelNo, qm)
		resetPartyItemPopuToRedis(partyID, username, 0, 0, partyType)
	}
	return
}

func MarkPartyItem(partyID int64, username string, supportUsername, nonsupportUsername string, partyType int) (popu int, lc int, ulc int, fc int, rank int, err error) {
	if err = addPartyItemLikeCnt(partyID, supportUsername, 1, 0); err != nil {
		return
	}

	if err = addPartyItemLikeCnt(partyID, nonsupportUsername, 0, 1); err != nil {
		return
	}

	resetPartyItemPopuToRedis(partyID, nonsupportUsername, 0, 1, partyType)

	return resetPartyItemPopuToRedis(partyID, supportUsername, 1, 0, partyType)
}

//SendFlower 送花
func SendFlower(partyID int64, username string, sendFlowerUsername string, partyCloseTime int64, partyType int) (popu int, lc int, ulc int, fc int, rank int, err error) {
	if sendFlowerUsername != "" {
		if err = addPartyItemFlowerToDB(partyID, username, sendFlowerUsername); err != nil {
			return
		}
		fdata := Flower{
			PartyId:            partyID,
			Username:           username,
			SendFlowerUsername: sendFlowerUsername,
			PartyType:          partyType,
			PartyCloseTime:     partyCloseTime,
		}
		fbyte, _ := json.Marshal(fdata)
		GMLog(constants.C1_PLAYER, constants.C2_PARTY, constants.C3_SEND_FLOWER, sendFlowerUsername, string(fbyte))
		GMLog(constants.C1_PLAYER, constants.C2_PARTY, constants.C3_SEND_FLOWER, username, string(fbyte))
		return addPartyItemFlowerToRedis(partyID, username, sendFlowerUsername, partyCloseTime, partyType)
	}

	return
}

func JoinedParty(partyID int64, username string) (string, error) {
	return joinedParty(partyID, username)
}
