package rt_party

import (
	// "appcfg"
	"constants"
	"db/rt_party"
	. "logger"
	"time"
	. "types"
	. "zk_manager"
)

func init() {
	// if appcfg.GetServerType() == "" {
	// 	if appcfg.GetBool("main_server", false) {
	// 		go cleanRTPartyLoop()
	// 	}

	// 	if appcfg.GetBool("rt_match_test", false) {
	// 		testRTPartyMatch()
	// 	}
	// }
}

func SendChat(username, content string, hostID int64, typee int, t int64) ([]RTPartyChat, error) {
	return rt_party.SendChat(username, content, hostID, typee, t)
}

func GetStage(hostID int64, now time.Time) (status RTPartyHostStatus, err error) {
	var h RTPartyHost
	if h, err = Sync(hostID, false, "", now); err != nil {
		return
	} else {
		status = GetRTPartyStage(&h, now)
		return
	}
}

func Sync(hostID int64, needExtra bool, syncUser string, now time.Time) (h RTPartyHost, err error) {
	return rt_party.Sync(hostID, needExtra, syncUser, now)
}

func SyncWithChat(hostID int64, needExtra bool, syncUser string, now time.Time) (h RTPartyHost, chats []RTPartyChat, err error) {
	return rt_party.SyncWithChat(hostID, needExtra, syncUser, now)
}

func CancelMatch(username string, hostID int64, now time.Time) (err error) {
	if ok, path := LockRTPartyHost(hostID, constants.RT_PARTY_DEFAULT_TIMEOUT); ok {
		defer Unlock(path)

		return LeaveHostBeforeBegin(username, hostID, now)
	}

	return
}

func MatchHost(username, nickname, head string, level, vip int, now time.Time) (hostID int64, status ErrorCode) {
	tnow := now
	for i := 0; i < 600; i++ {
		if createHostID, foundHostID, err := findHost(username, nickname, head, level, vip, tnow); err != nil {
			status = ERR_MATCH_RT_PARTY
			return
		} else if foundHostID > 0 {
			if err = joinHost(foundHostID, username, nickname, head, level, vip, tnow); err == nil {
				hostID = foundHostID
				return
			} else if err != constants.RTPartyHostFullErr && err != constants.RTPartyHostDeadErr {
				status = ERR_JOIN_RT_PARTY
				return
			}
			//continue
		} else if createHostID > 0 {
			hostID = createHostID
			return
		} else {
			status = ERR_MATCH_RT_PARTY_HOSTID
			return
		}

		time.Sleep(10 * time.Millisecond)
		tnow = time.Now()
	}

	status = ERR_MATCH_RT_PARTY_TIMEOUT

	return
}

func findHost(username, nickname, head string, level, vip int, now time.Time) (createHostID int64, foundHostID int64, err error) {
	var createHost bool
	if createHost, foundHostID, err = rt_party.CheckEmptyHost(now); err != nil {
		return
	} else if foundHostID > 0 {
		// Info("username", username, "find host", foundHostID)
		return
	} else if createHost {
		if ok, path := LockRTPartyHostList(constants.RT_PARTY_MATCH_TIMEOUT); ok {
			defer Unlock(path)

			var createHost bool
			if createHost, foundHostID, err = rt_party.CheckEmptyHost(now); err != nil {
				return
			} else if foundHostID > 0 {
				// Info("username", username, "find host", foundHostID)
				return
			} else if createHost {
				rdSubjects := getRandomRTPartySubject(now)
				createHostID, err = rt_party.CreateHost(username, nickname, head, level, vip, now, rdSubjects)
				// Info("username", username, "create host", createHostID)
			}
		} else {
			Err("lock party list", username, nickname)
			err = constants.LockRTPartyListErr
		}
	}

	return
}

func joinHost(hostID int64, username, nickname, head string, level, vip int, now time.Time) (err error) {
	if ok, path := LockRTPartyHost(hostID, constants.RT_PARTY_DEFAULT_TIMEOUT); ok {
		defer Unlock(path)

		return rt_party.JoinHost(username, nickname, head, level, vip, hostID, now)
	} else {
		Err("lock party host", username, nickname, hostID)
		err = constants.LockRTPartyHostErr
		return
	}
}

func LeaveHostBeforeBegin(username string, hostID int64, now time.Time) (err error) {
	if ok, path := LockRTPartyHost(hostID, constants.RT_PARTY_DEFAULT_TIMEOUT); ok {
		defer Unlock(path)

		return rt_party.LeaveHostBeforeBegin(username, hostID, now)
	} else {
		Err("lock party host", username, hostID)
		err = constants.LockRTPartyHostErr
		return
	}
}

func CancelAndLeaveHost(username string, hostID int64, now time.Time) (err error) {
	if ok, path := LockRTPartyHost(hostID, constants.RT_PARTY_DEFAULT_TIMEOUT); ok {
		defer Unlock(path)

		return rt_party.CancelAndLeaveHost(username, hostID, now)
	} else {
		Err("lock party host", username, hostID)
		err = constants.LockRTPartyHostErr
		return
	}
}

func OfferSubject(username, subject string, hostID int64, now time.Time) (err error) {
	if ok, path := LockRTPartyHost(hostID, constants.RT_PARTY_DEFAULT_TIMEOUT); ok {
		defer Unlock(path)

		return rt_party.OfferSubject(username, subject, hostID, now)
	} else {
		Err("lock party host", username, subject, hostID)
		err = constants.LockRTPartyHostErr
		return
	}
}

func VoteSubject(self, u1, s1, u2, s2 string, hostID int64, now time.Time) (err error) {
	if ok, path := LockRTPartyHost(hostID, constants.RT_PARTY_DEFAULT_TIMEOUT); ok {
		defer Unlock(path)

		return rt_party.VoteSubject(self, u1, s1, u2, s2, hostID, now)
	} else {
		Err("lock party host", self, u1, s1, u2, s2, hostID)
		err = constants.LockRTPartyHostErr
		return
	}
}

func UploadDress(username string, dyeCnt int, dyeMap map[string][7][4]float64, clothes []string, modelNo int, backgroundID, img, dressWord string, dressWordType int, hostID int64, now time.Time) (err error) {
	return rt_party.UploadDress(username, dyeCnt, dyeMap, clothes, modelNo, backgroundID, img, dressWord, dressWordType, hostID, now)
}

func VoteDress(self, u1, u2 string, hostID int64, now time.Time) (err error) {
	if ok, path := LockRTPartyHost(hostID, constants.RT_PARTY_DEFAULT_TIMEOUT); ok {
		defer Unlock(path)

		return rt_party.VoteDress(self, u1, u2, hostID, now)
	} else {
		Err("lock party host", self, u1, u2, hostID)
		err = constants.LockRTPartyHostErr
		return
	}
}
