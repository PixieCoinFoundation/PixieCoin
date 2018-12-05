package rt_party

import (
	"appcfg"
	"constants"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	. "logger"
	"rao"
	"time"
	. "types"
)

func listChatRedis(hostID int64, conn *redis.Conn) (res []RTPartyChat, err error) {
	var values []interface{}
	if values, err = redis.Values((*conn).Do("LRANGE", rao.GetHostChatListKey(hostID), 0, -1)); err != nil {
		Err(err)
		return
	} else {
		res = make([]RTPartyChat, 0)
		for len(values) > 0 {
			var cs string
			if values, err = redis.Scan(values, &cs); err != nil {
				Err(err)
				return
			}

			if len(cs) > 2 {
				cs = cs[1 : len(cs)-1]
				var c RTPartyChat

				if err = json.Unmarshal([]byte(cs), &c); err != nil {
					Err(err)
					continue
				}

				res = append(res, c)
			}
		}
	}

	return
}

func sendChatRedis(username, content string, hostID int64, typee int, t int64) (res []RTPartyChat, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetHostChatListKey(hostID)

	c := RTPartyChat{
		U:  username,
		C:  content,
		T:  t,
		Te: typee,
	}
	cb, _ := json.Marshal(c)

	conn.Send("LPUSH", key, "'"+string(cb)+"'")
	conn.Send("LTRIM", key, 0, constants.RT_PARTY_CHAT_SIZE-1)
	conn.Send("EXPIRE", key, constants.RT_PARTY_CHAT_EXPIRE_SECOND)
	conn.Send("LRANGE", key, 0, -1)

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	conn.Receive()
	conn.Receive()
	conn.Receive()

	var values []interface{}
	if values, err = redis.Values(conn.Receive()); err != nil {
		Err(err)
		return
	} else {
		res = make([]RTPartyChat, 0)
		for len(values) > 0 {
			var cs string
			if values, err = redis.Scan(values, &cs); err != nil {
				Err(err)
				return
			}

			if len(cs) > 2 {
				cs = cs[1 : len(cs)-1]
				var c RTPartyChat

				if err = json.Unmarshal([]byte(cs), &c); err != nil {
					Err(err)
					continue
				}

				res = append(res, c)
			}
		}
	}

	return
}

func voteDressRedis(self, u1, u2 string, hostID int64, now time.Time) (err error) {
	var vc int
	if u1 != "" {
		vc++
	}

	if u2 != "" {
		vc++
	}

	if vc <= 0 {
		return
	}

	c := rao.GetConn()
	defer c.Close()

	var h RTPartyHost
	if h, err = getHostRedis(hostID, &c, true, "", now); err != nil {
		return
	} else {
		allVoteDress := true
		for _, joiner := range h.Joiners {
			if joiner.Username == u1 || joiner.Username == u2 {
				if joiner.DressDoneTime <= 0 {
					err = constants.VoteToUndoneDressErr
					return
				}

				if joiner.Username == u1 {
					joiner.DressGotVote++
				}

				if joiner.Username == u2 {
					joiner.DressGotVote++
				}
			} else if joiner.Username == self {
				if joiner.Left {
					ErrMail("party player leave but vote dress", self, hostID)
					err = constants.PartyAlreadyLeftErr
					return
				}

				if joiner.DressVoteCnt > 0 {
					Err("already vote dress", self, u1, u2, hostID)
					err = constants.AlreadyVoteDressErr
					return
				} else {
					joiner.DressVoteCnt = vc
					joiner.DressVoteUsername1 = u1
					joiner.DressVoteUsername2 = u2
				}
			}

			if !joiner.Left && joiner.DressVoteCnt <= 0 {
				allVoteDress = false
			}
		}

		if allVoteDress && h.AllVoteDressTime <= 0 {
			h.AllVoteDressTime = now.Unix()
		}

		setHostRedis(&h, &c)
	}

	return
}

func offerSubjectRedis(username, subject string, hostID int64, now time.Time) (err error) {
	if subject == "" {
		return
	}

	c := rao.GetConn()
	defer c.Close()

	var h RTPartyHost
	if h, err = getHostRedis(hostID, &c, false, "", now); err != nil {
		return
	} else {
		allOfferSubject := true
		for _, joiner := range h.Joiners {
			if joiner.Username == username {
				if joiner.Left {
					ErrMail("party player leave but offer subject", username, subject, hostID)
					err = constants.PartyAlreadyLeftErr
					return
				}

				if joiner.OfferSubject.Name != "" {
					Err("already offer subject", username, subject, hostID, joiner.OfferSubject)
					err = constants.AlreadyOfferSubjectErr
					return
				} else if subject != "" {
					joiner.OfferSubjectTime = now.Unix()
					joiner.OfferSubject = RTPartySubject{Name: subject, OfferTime: now.Unix()}
				}
			}

			if !joiner.Left && joiner.OfferSubject.Name == "" {
				allOfferSubject = false
			}
		}

		if allOfferSubject && h.AllOfferSubjectTime <= 0 {
			h.AllOfferSubjectTime = now.Unix()
		}

		setHostRedis(&h, &c)
	}

	return
}

func voteSubjectRedis(self, u1, s1, u2, s2 string, hostID int64, now time.Time) (err error) {
	var vc int
	if s1 != "" {
		vc++
	}

	if s2 != "" {
		vc++
	}

	if vc <= 0 {
		return
	}

	c := rao.GetConn()
	defer c.Close()

	var h RTPartyHost
	if h, err = getHostRedis(hostID, &c, false, "", now); err != nil {
		return
	} else {
		var ignoreU1, ignoreU2 bool

		if u1 == "" {
			ignoreU1 = true
			for i, s := range h.DefaultSubjects {
				if s.Name == s1 {
					h.DefaultSubjects[i].GotVote++
					break
				}
			}
		}

		if u2 == "" {
			ignoreU2 = true
			for i, s := range h.DefaultSubjects {
				if s.Name == s2 {
					h.DefaultSubjects[i].GotVote++
					break
				}
			}
		}

		allVoteSubject := true
		for _, joiner := range h.Joiners {
			if !ignoreU1 && joiner.Username == u1 {
				if joiner.OfferSubject.Name == s1 {
					joiner.OfferSubject.GotVote++
				}
			} else if !ignoreU2 && joiner.Username == u2 {
				if joiner.OfferSubject.Name == s2 {
					joiner.OfferSubject.GotVote++
				}
			}

			if joiner.Username == self {
				if joiner.Left {
					ErrMail("party player leave but vote subject", self, hostID)
					err = constants.PartyAlreadyLeftErr
					return
				}

				if joiner.SubjectVoteCnt >= constants.RT_PARTY_SUBJECT_VOTE_CNT {
					Err("already vote subject", self, u1, s1, u2, s2, hostID)
					err = constants.AlreadyVoteSubjectErr
					return
				} else {
					joiner.SubjectVoteCnt = vc
				}
			}

			if !joiner.Left && joiner.SubjectVoteCnt <= 0 {
				allVoteSubject = false
			}
		}

		if allVoteSubject && h.AllVoteSubjectTime <= 0 {
			h.AllVoteSubjectTime = now.Unix()
		}

		setHostRedis(&h, &c)
	}

	return
}

func setJoinerDressRedis(hostID int64, username string, dyeCnt int, dyeMap map[string][7][4]float64, clothes []string, modelNo int, backgroundID, img, dressWord string, dressWordType int, now time.Time) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	var h RTPartyHost
	if h, err = getHostRedis(hostID, &conn, true, "", now); err != nil {
		return
	}

	for _, joiner := range h.Joiners {
		if joiner.Username == username && joiner.Left {
			Info("party player leave but upload dress", username, hostID)
			err = constants.PartyAlreadyLeftErr
			return
		}
	}

	var dyeMapStr string
	if dyeCnt > 0 {
		dmb, _ := json.Marshal(dyeMap)
		dyeMapStr = string(dmb)
	}

	clob, _ := json.Marshal(clothes)

	key := rao.GetHostJoinerHashKey(username, hostID)
	conn.Send("HMSET", key, "dyeCnt", dyeCnt, "dyeMap", "'"+dyeMapStr+"'", "clothes", "'"+string(clob)+"'", "modelNo", modelNo, "backgroundID", backgroundID, "dressWord", dressWord, "dressWordType", dressWordType, "img", img)

	if img != "" {
		conn.Send("HSET", key, "dressDoneTime", now.Unix())
	}
	conn.Send("EXPIRE", key, constants.RT_PARTY_DEFAULT_EXPIRE_SECOND)

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	rc := 2
	if img != "" {
		rc = 3
	}

	for i := 0; i < rc; i++ {
		conn.Receive()
	}

	if img != "" {
		for _, joiner := range h.Joiners {
			if joiner.Username == username {
				joiner.DressDoneTime = now.Unix()
			}
		}

		allFinishDress := true
		for _, joiner := range h.Joiners {
			if !joiner.Left && joiner.DressDoneTime <= 0 {
				allFinishDress = false
			}
		}

		if allFinishDress && h.AllFinishDressTime <= 0 {
			h.AllFinishDressTime = now.Unix()
		}

		setHostRedis(&h, &conn)
	}

	return
}

func getHostAndChatRedis(hostID int64, needExtra bool, syncUser string, now time.Time) (h RTPartyHost, chats []RTPartyChat, err error) {
	c := rao.GetConn()
	defer c.Close()

	if h, err = getHostRedis(hostID, &c, needExtra, syncUser, now); err != nil {
		return
	}

	chats, err = listChatRedis(hostID, &c)
	return
}

func getHostRedis(hostID int64, conn *redis.Conn, needExtra bool, syncUser string, now time.Time) (h RTPartyHost, err error) {
	var c redis.Conn
	if conn == nil {
		c = rao.GetConn()
		defer c.Close()
	} else {
		c = *conn
	}

	c.Send("GET", rao.GetHostSelfKey(hostID))
	if syncUser != "" {
		key := rao.GetHostJoinerHashKey(syncUser, hostID)
		c.Send("HSET", key, "lastSyncTime", now.Unix())
		c.Send("EXPIRE", key, constants.RT_PARTY_DEFAULT_EXPIRE_SECOND)
	}

	if err = c.Flush(); err != nil {
		Err(err, hostID)
		return
	}

	var vs string
	if vs, err = redis.String(c.Receive()); err != nil {
		Err(err, hostID)
		return
	} else if len(vs) > 2 {
		vs = vs[1 : len(vs)-1]

		if err = json.Unmarshal([]byte(vs), &h); err != nil {
			Err(err, vs, hostID)
			return
		}

		if syncUser != "" {
			c.Receive()
			c.Receive()
		}

		if needExtra {
			if err = setJoinerExtra(&h, &c); err != nil {
				return
			}
		}

		return
	} else {
		Err("rt party host format wrong:", hostID, vs)
		err = constants.PartyHostFormatErr
		return
	}
}

func setJoinerExtra(h *RTPartyHost, conn *redis.Conn) (err error) {
	if len(h.Joiners) > 0 {
		for _, joiner := range h.Joiners {
			(*conn).Send("HGETALL", rao.GetHostJoinerHashKey(joiner.Username, h.ID))
		}

		if err = (*conn).Flush(); err != nil {
			Err(err)
			return
		}

		var values []interface{}
		for _, joiner := range h.Joiners {
			var je RTPartyHostJoinerExtra

			if values, err = redis.Values((*conn).Receive()); err != nil {
				Err(err)
				return
			} else if err = redis.ScanStruct(values, &je); err != nil {
				Err(err)
				return
			}

			joiner.DyeCnt = je.DyeCnt
			if len(je.DyeMap) > 2 {
				json.Unmarshal([]byte(je.DyeMap[1:len(je.DyeMap)-1]), &joiner.DyeMap)
			}

			if len(je.Clothes) > 2 {
				json.Unmarshal([]byte(je.Clothes[1:len(je.Clothes)-1]), &joiner.Clothes)
			}

			joiner.Img = je.Img
			joiner.DressWord = je.DressWord
			joiner.DressWordType = je.DressWordType
			joiner.ModelNo = je.ModelNo
			joiner.BackgroundID = je.BackgroundID
			joiner.DressDoneTime = je.DressDoneTime
			joiner.LastSyncTime = je.LastSyncTime
		}
	}

	return
}

func setHostRedis(h *RTPartyHost, conn *redis.Conn) (err error) {
	var c redis.Conn
	if conn != nil {
		c = *conn
	} else {
		c = rao.GetConn()
		defer c.Close()
	}

	hb, _ := json.Marshal(*h)
	hk := rao.GetHostSelfKey(h.ID)
	c.Send("SET", hk, "'"+string(hb)+"'")
	c.Send("EXPIRE", hk, constants.RT_PARTY_DEFAULT_EXPIRE_SECOND)

	if conn == nil {
		if err = c.Flush(); err != nil {
			Err(err)
		}
	}

	return
}

func createHostRedis(h *RTPartyHost, now time.Time) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	//add in list
	lk := rao.GetHostListKey(now.Format("20060102"))
	conn.Send("ZADD", lk, len(h.Joiners), h.ID)
	conn.Send("EXPIRE", lk, constants.RT_PARTY_DEFAULT_EXPIRE_SECOND)

	//set host self
	setHostRedis(h, &conn)

	if err = conn.Flush(); err != nil {
		Err(err)
	}
	return
}

func joinHostRedis(joinUsername, joinNickname, joinHead string, level, vip int, hostID int64, now time.Time) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	var h RTPartyHost
	if h, err = getHostRedis(hostID, &conn, false, "", now); err != nil {
		return
	}

	if len(h.Joiners) >= appcfg.GetInt("rt_party_manual_full_size", constants.RT_PARTY_PLAYER_SIZE) {
		err = constants.RTPartyHostFullErr
		return
	} else if h.StartTime+constants.RT_PARTY_DEAD_TIME <= now.Unix() {
		Info("host:", hostID, "dead")

		//delete this party
		conn.Send("ZREM", rao.GetHostListKey(now.Format("20060102")), hostID)

		err = constants.RTPartyHostDeadErr
		return
	}

	for _, joiner := range h.Joiners {
		if joiner.Username == joinUsername {
			return
		}
	}

	joiner := RTPartyHostJoiner{
		Username: joinUsername,
		Nickname: joinNickname,
		Head:     joinHead,
		Level:    level,
		VIP:      vip,
	}

	if h.Joiners == nil {
		h.Joiners = make([]*RTPartyHostJoiner, 0)
	}
	h.Joiners = append(h.Joiners, &joiner)

	if len(h.Joiners) >= appcfg.GetInt("rt_party_manual_full_size", constants.RT_PARTY_PLAYER_SIZE) {
		h.BeginTime = now.Unix()
	}

	conn.Send("ZADD", rao.GetHostListKey(now.Format("20060102")), len(h.Joiners), hostID)
	setHostRedis(&h, &conn)

	if err = conn.Flush(); err != nil {
		Err(err)
	}
	return
}

func leaveHostBeforeBeginRedis(username string, hostID int64, now time.Time) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	var h RTPartyHost
	if h, err = getHostRedis(hostID, &conn, false, "", now); err != nil {
		return
	}

	if len(h.Joiners) >= appcfg.GetInt("rt_party_manual_full_size", constants.RT_PARTY_PLAYER_SIZE) {
		Info("party already started:", username, hostID)
		err = constants.PartyAlreadyStartErr
		return
	}

	have := false
	njs := make([]*RTPartyHostJoiner, 0)
	for _, joiner := range h.Joiners {
		if joiner.Username == username {
			have = true
		} else {
			njs = append(njs, joiner)
		}
	}

	if have {
		h.Joiners = njs
		setHostRedis(&h, &conn)
		conn.Send("ZADD", rao.GetHostListKey(now.Format("20060102")), len(h.Joiners), hostID)

		if err = conn.Flush(); err != nil {
			Err(err)
		}
	} else {
		Info("party not attend:", username, hostID)
		err = constants.PartyNotAttendErr
	}

	return
}

func cancelAndLeaveHostRedis(username string, hostID int64, now time.Time) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	var h RTPartyHost
	if h, err = getHostRedis(hostID, &conn, true, "", now); err != nil {
		return
	}

	if len(h.Joiners) < appcfg.GetInt("rt_party_manual_full_size", constants.RT_PARTY_PLAYER_SIZE) {
		//cancel and leave
		have := false
		njs := make([]*RTPartyHostJoiner, 0)
		for _, joiner := range h.Joiners {
			if joiner.Username == username {
				have = true
			} else {
				njs = append(njs, joiner)
			}
		}

		if have {
			h.Joiners = njs
			setHostRedis(&h, &conn)
			conn.Send("ZADD", rao.GetHostListKey(now.Format("20060102")), len(njs), hostID)

			if err = conn.Flush(); err != nil {
				Err(err)
			}
		} else {
			Info("party not attend:", username, hostID)
			err = constants.PartyNotAttendErr
		}
	} else {
		//mark leave
		change := false

		allOfferSubject := true
		allVoteSuject := true
		allDressDone := true
		allVoteDress := true

		for _, joiner := range h.Joiners {
			if joiner.Username == username {
				change = true
				joiner.Left = true
			} else if !joiner.Left {
				if joiner.OfferSubject.Name == "" {
					allOfferSubject = false
				}

				if joiner.SubjectVoteCnt <= 0 {
					allVoteSuject = false
				}

				if joiner.DressDoneTime <= 0 {
					allDressDone = false
				}

				if joiner.DressVoteCnt <= 0 {
					allVoteDress = false
				}
			}
		}

		if h.AllOfferSubjectTime <= 0 && allOfferSubject {
			change = true
			h.AllOfferSubjectTime = now.Unix()
		}

		if h.AllVoteSubjectTime <= 0 && allVoteSuject {
			change = true
			h.AllVoteSubjectTime = now.Unix()
		}

		if h.AllFinishDressTime <= 0 && allDressDone {
			change = true
			h.AllFinishDressTime = now.Unix()
		}

		if h.AllVoteDressTime <= 0 && allVoteDress {
			change = true
			h.AllVoteDressTime = now.Unix()
		}

		if change {
			setHostRedis(&h, &conn)

			if err = conn.Flush(); err != nil {
				Err(err)
			}
		} else {
			Info("party not attend:", username, hostID)
			err = constants.PartyNotAttendErr
		}
	}

	return
}

func checkEmptyHostRedis(now time.Time) (needCreate bool, hostID int64, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	lk := rao.GetHostListKey(now.Format("20060102"))

	var values []interface{}
	if values, err = redis.Values(conn.Do("ZRANGE", lk, 0, 0, "WITHSCORES")); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			err = nil
			needCreate = true
			Info("no rt party host opening...")
		} else {
			Err(err)
		}
		return
	} else {
		for len(values) > 0 {
			if values, err = redis.Scan(values, &hostID); err != nil {
				Err(err)
				return
			}

			var score int
			if values, err = redis.Scan(values, &score); err != nil {
				Err(err)
				return
			}

			if score >= appcfg.GetInt("rt_party_manual_full_size", constants.RT_PARTY_PLAYER_SIZE) {
				hostID = 0
				needCreate = true
			}

			return
		}

		needCreate = true
		return
	}

	return
}
