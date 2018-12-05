package party

import (
	"appcfg"
	"constants"
	"github.com/garyburd/redigo/redis"
	. "logger"
	"math/rand"
	"rao"
	"strconv"
	"strings"
	"sync"
	"time"
	. "types"
)

type simpleParty struct {
	ID      int64
	EndTime int64
}

type partyMarkInfo struct {
	simpleParty
	Size    int
	SizeEnd int
}

type partyMark struct {
	TotalSize int
	Marks     []partyMarkInfo
}

var cosplayPartyMark, normalPartyMark partyMark
var cosLock, norLock sync.RWMutex

func init() {
	if appcfg.GetServerType() != "" {
		return
	}
	now := time.Now()

	if cosplayPartyRawInfo, size, err := getCosplayPartyMarkInfoRedis(now); err != nil {
		panic(err)
	} else {
		Info("cos mark info:", len(cosplayPartyRawInfo), size)
		cosplayPartyMark.Marks = cosplayPartyRawInfo
		cosplayPartyMark.TotalSize = size
	}

	if normalPartyRawInfo, size, err := getNormalPartyMarkInfoRedis(now); err != nil {
		panic(err)
	} else {
		Info("normal mark info:", len(normalPartyRawInfo), size)
		normalPartyMark.Marks = normalPartyRawInfo
		normalPartyMark.TotalSize = size
	}

	go loop()
}

func loop() {
	for {
		time.Sleep(1 * time.Minute)

		now := time.Now()
		loopCos(now)
		loopNormal(now)
	}
}

func loopCos(now time.Time) {
	if cosplayPartyRawInfo, size, err := getCosplayPartyMarkInfoRedis(now); err != nil {
		Err(err)
	} else {
		cosLock.Lock()
		defer cosLock.Unlock()
		Info("cos mark info:", len(cosplayPartyRawInfo), size)
		cosplayPartyMark.Marks = cosplayPartyRawInfo
		cosplayPartyMark.TotalSize = size
	}
}

func loopNormal(now time.Time) {
	if normalPartyRawInfo, size, err := getNormalPartyMarkInfoRedis(now); err != nil {
		Err(err)
	} else {
		norLock.Lock()
		defer norLock.Unlock()
		Info("normal mark info:", len(normalPartyRawInfo), size)
		normalPartyMark.Marks = normalPartyRawInfo
		normalPartyMark.TotalSize = size
	}
}

func getNormalPartyMarkInfoRedis(now time.Time) (pms []partyMarkInfo, totalSize int, err error) {
	conn := rao.GetConn()
	defer conn.Close()
	nt := now.Unix()

	conn.Send("ZREVRANGEBYSCORE", rao.GetGlobalCasualPartyListKey(), nt*10, nt, "WITHSCORES")
	conn.Send("ZREVRANGEBYSCORE", rao.GetGlobalPrizePartyListKey(), nt*10, nt, "WITHSCORES")

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	var values []interface{}
	pms = make([]partyMarkInfo, 0)
	pis := make([]simpleParty, 0)
	for i := 0; i < 2; i++ {
		if values, err = redis.Values(conn.Receive()); err != nil {
			if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
				Info("no info for getParty")
				err = nil
				continue
			} else {
				Err(err)
				return
			}
		} else {
			var token string
			var et int64
			for len(values) > 0 {
				if values, err = redis.Scan(values, &token); err != nil {
					Err(err)
					return
				}

				if values, err = redis.Scan(values, &et); err != nil {
					Err(err)
					return
				}

				pid, _ := strconv.ParseInt(strings.Split(token, "_")[0], 10, 64)

				if pid > 0 {
					pis = append(pis, simpleParty{ID: pid, EndTime: et})

					conn.Send("LLEN", rao.GetPartyItemListKey(pid))
				}
			}
		}
	}

	if len(pis) > 0 {
		if err = conn.Flush(); err != nil {
			Err(err)
			return
		}

		var size int
		for i := 0; i < len(pis); i++ {
			if size, err = redis.Int(conn.Receive()); err != nil {
				Err(err)
				return
			}

			if size >= 2 {
				totalSize += size
				pms = append(pms, partyMarkInfo{simpleParty: simpleParty{ID: pis[i].ID, EndTime: pis[i].EndTime}, SizeEnd: totalSize, Size: size})
			}
		}
	}

	return
}

func getCosplayPartyMarkInfoRedis(now time.Time) (pms []partyMarkInfo, totalSize int, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	//获取舞会列表
	lk := rao.GetGlobalCosplayPartyListKey()
	nt := now.Unix()

	var values []interface{}
	if values, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", lk, nt*10, nt, "WITHSCORES")); err != nil {
		if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
			Info("no info for getParty")
			err = nil
		} else {
			Err(err)
		}
		return
	} else {
		pms = make([]partyMarkInfo, 0)
		pis := make([]simpleParty, 0)
		var token string
		var et int64
		for len(values) > 0 {
			if values, err = redis.Scan(values, &token); err != nil {
				Err(err)
				return
			}

			if values, err = redis.Scan(values, &et); err != nil {
				Err(err)
				return
			}

			pid, _ := strconv.ParseInt(strings.Split(token, "_")[0], 10, 64)

			if pid > 0 {
				pis = append(pis, simpleParty{ID: pid, EndTime: et})

				conn.Send("LLEN", rao.GetPartyItemListKey(pid))
			}
		}

		if len(pis) > 0 {
			if err = conn.Flush(); err != nil {
				Err(err)
				return
			}

			var size int
			for i := 0; i < len(pis); i++ {
				if size, err = redis.Int(conn.Receive()); err != nil {
					Err(err)
					return
				}

				if size >= 2 {
					totalSize += size
					pms = append(pms, partyMarkInfo{simpleParty: simpleParty{ID: pis[i].ID, EndTime: pis[i].EndTime}, SizeEnd: totalSize, Size: size})
				}
			}
		}
	}
	return
}

func getRandomPartyForMark(nowt int64) (pid int64, size int) {
	r := rand.Intn(9)

	if r <= 2 {
		//cosplay
		pid, size = getRandomCosParty(nowt)
		if pid <= 0 {
			return getRandomNormalParty(nowt)
		}
	} else {
		pid, size = getRandomNormalParty(nowt)
		if pid <= 0 {
			return getRandomCosParty(nowt)
		}
	}

	return
}

func getRandomNormalParty(nowt int64) (pid int64, size int) {
	//normal
	norLock.RLock()
	defer norLock.RUnlock()
	if len(normalPartyMark.Marks) > 0 && normalPartyMark.TotalSize > 0 {
		for i := 0; i < 3; i++ {
			r := rand.Intn(normalPartyMark.TotalSize)

			for i := 0; i < len(normalPartyMark.Marks); i++ {
				if r <= normalPartyMark.Marks[i].SizeEnd {
					if normalPartyMark.Marks[i].EndTime <= nowt {
						break
					} else {
						return normalPartyMark.Marks[i].ID, normalPartyMark.Marks[i].Size
					}
				}
			}
		}

	}

	return
}

func getRandomCosParty(nowt int64) (pid int64, size int) {
	//cosplay
	cosLock.RLock()
	defer cosLock.RUnlock()

	if len(cosplayPartyMark.Marks) > 0 && cosplayPartyMark.TotalSize > 0 {
		for i := 0; i < 3; i++ {
			r := rand.Intn(cosplayPartyMark.TotalSize)

			for i := 0; i < len(cosplayPartyMark.Marks); i++ {
				if r <= cosplayPartyMark.Marks[i].SizeEnd {
					if cosplayPartyMark.Marks[i].EndTime <= nowt {
						break
					} else {
						return cosplayPartyMark.Marks[i].ID, cosplayPartyMark.Marks[i].Size
					}
				}
			}
		}

	}

	return
}

func getRandomPartyItemPairsFromRedis(username string, nowt int64) (ps []PartyItemPair, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	ps = make([]PartyItemPair, 0)

	for j := 0; j < 5; j++ {
		pid, size := getRandomPartyForMark(nowt)

		if pid == 0 || size == 0 {
			Info("no party at all..", username)
			err = constants.NoPartyOpenErr
			continue
		}

		m := make(map[int]int)
		lk := rao.GetPartyItemListKey(pid)

		pp := PartyItemPair{}
		for i := 0; i < 20; i++ {
			index := rand.Intn(size)
			if m[index] == 0 {
				m[index] = 1

				var joiner string
				if joiner, err = redis.String(conn.Do("LINDEX", lk, index)); err != nil {
					if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
						Info("no random item,size:", size, "random index:", index)
						continue
					} else {
						Err(err)
						return
					}
				} else if joiner != username {
					var p PartyItem
					if p, err = getPartyItemForPartyUserFromRedis(pid, joiner, &conn, false, ""); err != nil {
						return
					} else if p.Partner != username && p.ID > 0 {
						if pp.Item1.Username == "" {
							pp.Item1 = p
						} else if pp.Item2.Username == "" {
							pp.Item2 = p
						}

						if pp.Item1.Username != "" && pp.Item2.Username != "" {
							ps = append(ps, pp)
							break
						}
					}
				}
			}
		}
	}

	return
}
