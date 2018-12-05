package rt_party

import (
	"appcfg"
	"constants"
	"db/rt_party"
	"fmt"
	. "logger"
	"math/rand"
	"service/params"
	"sort"
	"sync"
	"time"
	. "types"
	. "zk_manager"
)

// func init() {
// 	if appcfg.GetServerType() == constants.SERVER_TYPE_SIMPLE_TEST {
// 		now := time.Now().Unix()
// 		s := RTPartyHostJoiners{}
// 		s = append(s, &RTPartyHostJoiner{Username: "1"})
// 		s = append(s, &RTPartyHostJoiner{Username: "2", Left: true})
// 		s = append(s, &RTPartyHostJoiner{Username: "5", DressDoneTime: now})
// 		s = append(s, &RTPartyHostJoiner{Username: "3"})
// 		s = append(s, &RTPartyHostJoiner{Username: "4"})

// 		sort.Sort(s)

// 		for i := 0; i < len(s); i++ {
// 			Info(s[i].Username)
// 		}

// 		panic("end")
// 	}
// }

func cleanRTPartyLoop() {
	for {
		yday := time.Now().Add(-24 * time.Hour).Format("20060102")
		if hosts, err := rt_party.ListHostDB(yday); err == nil {
			for _, h := range hosts {
				DelRTPartyHostLock(h)
			}
		}

		yday1 := time.Now().Add(-48 * time.Hour).Format("20060102")
		rt_party.DelHostByDs(yday1)

		time.Sleep(10 * time.Minute)
	}

}

func getRandomRTPartySubject(now time.Time) []RTPartySubject {
	src := params.GetRTPartyDefaultSubjects()
	res := make([]RTPartySubject, 0)

	if src.Size() > 0 {
		for i := 0; i < constants.RT_PARTY_DEFAULT_SUBJECT_SIZE; i++ {
			r := rand.Intn(src.Size())
			if name, ok := src.Get(r); ok {
				res = append(res, RTPartySubject{Name: name.(string), OfferTime: now.Unix()})
				src.Remove(r)

				if src.Size() <= 0 {
					break
				}
			}
		}

		var i int
		for len(res) < constants.RT_PARTY_DEFAULT_SUBJECT_SIZE {
			i++
			res = append(res, RTPartySubject{Name: fmt.Sprintf("测试主题%d", i), OfferTime: now.Unix()})
		}

		return res
	} else {
		return []RTPartySubject{RTPartySubject{Name: "春风得意马蹄疾", OfferTime: now.Unix()}, RTPartySubject{Name: "夏日荷花别样红", OfferTime: now.Unix()}, RTPartySubject{Name: "不辞冰雪为卿热", OfferTime: now.Unix()}}
	}
}

func getReward(rank int) (gold int, diamond int, wire int) {
	switch rank {
	case 1:
		return 2500, 5, 2
	case 2:
		return 1500, 0, 1
	case 3:
		return 1000, 0, 1
	case 4:
		return 700, 0, 1
	case 5:
		return 400, 0, 0
	}

	return 0, 0, 0
}

func GetPlayerRTPartyReward(h *RTPartyHost, username string, crg, crd, crw, cge int) (gold int, diamond int, wire int, exp int) {
	if len(h.FinalRank) <= 0 {
		return
	}

	var rank int
	var firstUsername, vu1, vu2 string
	for i, h := range h.FinalRank {
		if h.Username == username {
			if h.Left || h.DressDoneTime <= 0 {
				return
			}

			exp = constants.RT_PARTY_EXP
			rank = i + 1
			vu1 = h.DressVoteUsername1
			vu2 = h.DressVoteUsername2
		}

		if i == 0 {
			firstUsername = h.Username
		}
	}

	var voteGold int
	if vu1 == firstUsername {
		voteGold += 500
	}
	if vu2 == firstUsername {
		voteGold += 500
	}

	rg, rd, rw := getReward(rank)
	rg += voteGold

	lrg := constants.DAY_RT_PARTY_GOLD_LIMIT - crg
	lrd := constants.DAY_RT_PARTY_DIAMOND_LIMIT - crd
	lrw := constants.DAY_RT_PARTY_WIRE_LIMIT - crw
	lge := constants.DAY_RT_PARTY_EXP_LIMIT - cge

	if rg > lrg {
		rg = lrg
	}
	if rg > 0 {
		gold += rg
	}

	if rd > lrd {
		rd = lrd
	}
	if rd > 0 {
		diamond += rd
	}

	if rw > lrw {
		rw = lrw
	}
	if rw > 0 {
		wire += rw
	}

	if exp > lge {
		exp = lge
		if exp < 0 {
			exp = 0
		}
	}

	return
}

func GetRTPartyRewardToken(hostID int64, username string) string {
	return fmt.Sprintf("%s_%d_rtreward", username, hostID)
}

func CheckRTPartyTime(now time.Time) (timeOK bool) {
	if appcfg.GetBool("check_rt_party_time", true) {
		hour := now.Hour()
		if (hour >= 12 && hour <= 14) || (hour >= 19 && hour <= 23) {
			timeOK = true
		}
	} else {
		timeOK = true
	}

	return
}

func GetRTPartyStage(h *RTPartyHost, now time.Time) RTPartyHostStatus {
	beginTime := h.BeginTime

	if len(h.Joiners) < appcfg.GetInt("rt_party_manual_full_size", constants.RT_PARTY_PLAYER_SIZE) && beginTime <= 0 {
		return RT_PARTY_HOST_STATUS_START
	}

	nt := now.Unix()
	var allOfferSubjectTime, allVoteSujectTime, allFinishDressTime int64
	if nt < beginTime+constants.RT_PARTY_SUBJECT_OFFER_TIME && h.AllOfferSubjectTime <= 0 {
		return RT_PARTY_HOST_STATUS_SUBJECT_SELECT
	}

	allOfferSubjectTime = h.AllOfferSubjectTime
	if allOfferSubjectTime <= 0 {
		allOfferSubjectTime = beginTime + constants.RT_PARTY_SUBJECT_OFFER_TIME
	}

	if nt < allOfferSubjectTime+constants.RT_PARTY_SUBJECT_VOTE_TIME && h.AllVoteSubjectTime <= 0 {
		return RT_PARTY_HOST_STATUS_SUBJECT_VOTE
	}

	allVoteSujectTime = h.AllVoteSubjectTime
	if allVoteSujectTime <= 0 {
		allVoteSujectTime = allOfferSubjectTime + constants.RT_PARTY_SUBJECT_VOTE_TIME
	}

	if nt < allVoteSujectTime+constants.RT_PARTY_DRESS_TIME+appcfg.GetInt64("test_rt_party_dress_extend", 0) && h.AllFinishDressTime <= 0 {
		return RT_PARTY_HOST_STATUS_DRESS
	}

	allFinishDressTime = h.AllFinishDressTime
	if allFinishDressTime <= 0 {
		allFinishDressTime = allVoteSujectTime + constants.RT_PARTY_DRESS_TIME + appcfg.GetInt64("test_rt_party_dress_extend", 0)
	}

	if nt < allFinishDressTime+constants.RT_PARTY_DRESS_VOTE_TIME+appcfg.GetInt64("test_rt_party_dress_extend", 0) && h.AllVoteDressTime <= 0 {
		//check if leave 4 player.skip vote

		var leftCnt int
		for _, joiner := range h.Joiners {
			if joiner.Left {
				leftCnt++
			}
		}

		if leftCnt < appcfg.GetInt("rt_party_manual_full_size", constants.RT_PARTY_PLAYER_SIZE)-1 {
			return RT_PARTY_HOST_STATUS_DRESS_VOTE
		}
	}

	return RT_PARTY_HOST_STATUS_DRESS_VOTE_END
}

func CheckRTPartyStage(h *RTPartyHost, now time.Time) (stage RTPartyHostStatus, fs string, fr []*RTPartyHostJoiner, err error) {
	stage = GetRTPartyStage(h, now)
	fs = h.FinalSubject
	fr = h.FinalRank

	if h.FinalSubject == "" && stage >= RT_PARTY_HOST_STATUS_DRESS {
		//主题选择阶段结束后，选出最终主题
		if ok, path := LockRTPartyHost(h.ID, constants.RT_PARTY_DEFAULT_TIMEOUT); ok {
			defer Unlock(path)

			var nh RTPartyHost
			if nh, err = Sync(h.ID, true, "", now); err != nil {
				return
			} else {
				h = &nh
			}

			stage = GetRTPartyStage(h, now)
			if h.FinalSubject == "" && stage >= RT_PARTY_HOST_STATUS_DRESS {
				var vote int
				var tmpFs string
				//默认主题情况
				for _, s := range h.DefaultSubjects {
					if s.Name != "" && s.GotVote >= vote {
						tmpFs = s.Name
						vote = s.GotVote
					}
				}

				//玩家主题
				for _, joiner := range h.Joiners {
					if joiner.OfferSubject.Name != "" && joiner.OfferSubject.GotVote >= vote {
						tmpFs = joiner.OfferSubject.Name
						vote = joiner.OfferSubject.GotVote
					}
				}

				if tmpFs != "" {
					h.FinalSubject = tmpFs

					if h.AllVoteSubjectTime <= 0 {
						h.AllVoteSubjectTime = now.Unix()
					}
				} else {
					Err("final subject empty", h.ID)
					err = constants.SubjectEmptyErr
					return
				}

				if err = rt_party.Flush(h); err != nil {
					return
				}
			}

			fs = h.FinalSubject
		} else {
			err = constants.LockRTPartyHostErr
			return
		}
	}

	if (h.FinalRank == nil || len(h.FinalRank) == 0) && stage >= RT_PARTY_HOST_STATUS_DRESS_VOTE_END {
		//搭配投票结束后 检查是否有未投票的玩家 自动随机投票
		if ok, path := LockRTPartyHost(h.ID, constants.RT_PARTY_DEFAULT_TIMEOUT); ok {
			defer Unlock(path)

			var nh RTPartyHost
			if nh, err = Sync(h.ID, true, "", now); err != nil {
				return
			} else {
				h = &nh
			}
			stage = GetRTPartyStage(h, now)

			if (h.FinalRank == nil || len(h.FinalRank) == 0) && stage >= RT_PARTY_HOST_STATUS_DRESS_VOTE_END {
				if h.FinalRank == nil {
					h.FinalRank = make([]*RTPartyHostJoiner, 0)
				}
				for _, joiner := range h.Joiners {
					if joiner.DressVoteCnt == 0 && joiner.DressDoneTime > 0 {
						voteForPlayer(h, joiner.Username)
					}

					h.FinalRank = append(h.FinalRank, joiner)
				}

				if h.AllVoteDressTime <= 0 {
					h.AllVoteDressTime = now.Unix()
				}

				//得分排行
				sort.Sort(h.FinalRank)
				if err = rt_party.Flush(h); err != nil {
					return
				}
			}

			fr = h.FinalRank
		} else {
			err = constants.LockRTPartyHostErr
			return
		}
	}
	//TODO 机器人设定?

	return
}

func voteForPlayer(h *RTPartyHost, username string) {
	tmpJoiners := make([]*RTPartyHostJoiner, 0)

	for _, joiner := range h.Joiners {
		if joiner.DressDoneTime > 0 {
			if joiner.Username == username {
				joiner.DressVoteCnt = constants.RT_PARTY_DRESS_VOTE_CNT
			} else {
				tmpJoiners = append(tmpJoiners, joiner)
			}
		}
	}

	length := len(tmpJoiners)

	if length > 0 {
		for i := 0; i < constants.RT_PARTY_DRESS_VOTE_CNT; i++ {
			r := rand.Intn(length)

			tmpJoiners[r].DressGotVote++
		}
	} else {
		Err("auto vote dress joiner size 0.maybe too many player leave this host", length, h.ID, username)
	}
}

func testRTPartyMatch() {
	Info("----begin rt party match test---")
	testSize := 1000

	var okCntLock, errCntLock sync.Mutex
	var okCnt, errCnt int
	hostSizeMap := make(map[string]int)

	for i := 0; i < testSize; i++ {
		is := fmt.Sprintf("%d", i)

		go func() {
			hostID, status := MatchHost(is, is, is, i, i, time.Now())
			if status == NO_ERROR {
				okCntLock.Lock()
				defer okCntLock.Unlock()

				okCnt++
				hostSizeMap[fmt.Sprintf("%d", hostID)] += 1
				// Info("-----ok", i, hostID)
			} else {
				errCntLock.Lock()
				defer errCntLock.Unlock()

				errCnt++
				Info("-----error", i, status)
			}
		}()
	}

	time.Sleep(30 * time.Second)

	for k, v := range hostSizeMap {
		if v != appcfg.GetInt("rt_party_manual_full_size", constants.RT_PARTY_PLAYER_SIZE) {
			Info("host", k, "size", v)
		}
	}

	Info("ok cnt", okCnt)
	Info("err cnt", errCnt)
	Info("----end rt party match test---")
}
