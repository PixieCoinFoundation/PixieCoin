package event

import (
	"appcfg"
	"constants"
	"dao"
	"database/sql"
	"encoding/json"
	"fmt"
	. "logger"
	"math/rand"
	"service/params"
	"strconv"
	"strings"
	"sync"
	"time"
	. "types"
)

type EventInfoDetail struct {
	//pk multi
	PKMulti           float64
	PKMultiStartVIP   int
	PKMultiStartLevel int
	PKMultiEndVIP     int
	PKMultiEndLevel   int

	//cos multi
	CosCmtAlwaysReward bool
	CosMultiStartVIP   int
	CosMultiStartLevel int
	CosMultiEndVIP     int
	CosMultiEndLevel   int

	//normal mission drop
	EventDrop         bool
	EventDropID       int64
	EventDropRate     int
	EventDropMaxNum   int
	EventDropDayLimit int

	//sum game content
	LogGame        bool
	LogGameEventID int64
	LogGameEndTime int64

	//sum diamond consume
	LogDiamondConsume        bool
	LogDiamondConsumeEventID int64

	//mission beilv
	MissionBeilvConfig          map[string]float64
	MissionBeilvMultiStartVIP   int
	MissionBeilvMultiStartLevel int
	MissionBeilvMultiEndVIP     int
	MissionBeilvMultiEndLevel   int
}

var lock sync.RWMutex

var OpenEvents []Event

var CosCmtAlwaysReward bool
var CosMultiStartVIP int
var CosMultiStartLevel int
var CosMultiEndVIP int
var CosMultiEndLevel int

var PKMulti float64
var PKMultiStartVIP int
var PKMultiStartLevel int
var PKMultiEndVIP int
var PKMultiEndLevel int

//event drop
var EventDrop bool
var EventDropID int64
var EventDropRate int
var EventDropMaxNum int
var EventDropDayLimit int
var EventDropType int

//sum game content
var LogGame bool
var LogGameEventID int64
var LogGameEndTime int64

//sum diamond consume
var LogDiamondConsume bool
var LogDiamondConsumeEventID int64

var MissionBeilvConfig map[string]float64
var MissionBeilvMultiStartVIP int
var MissionBeilvMultiStartLevel int
var MissionBeilvMultiEndVIP int
var MissionBeilvMultiEndLevel int

var specialMissionEventID int64
var specialMissionStartTime int64
var specialMissionEndTime int64
var specialMissionMap map[string]int

var pkDropEventID int64
var pkDropStartTime int64
var pkDropEndTime int64

// var pkDropRate int
// 	var		pkDropMaxNum int
// 	var 		pkDropDayLimit int

func init() {
	if appcfg.GetServerType() == "" {
		go refreshEvents()
	}
}

func refreshEvents() {
	for {
		if events, dt, info, err := getOpenEventsFromDB(); err != nil {
			Err(err)
		} else {
			var smst, smet, smeid, pkdst, pkdet, pkeid int64
			smm := make(map[string]int)
			for _, e := range events {
				if e.Info.Type == 1002 {
					smst = e.StartTime
					smet = e.EndTime

					for _, l := range e.Info.SpecialLevels {
						smm[l.LevelID] = 1
					}
					smeid = e.ID
					break
				} else if e.Info.Type == 1001 {
					pkeid = e.ID
					pkdst = e.StartTime
					pkdet = e.EndTime
					// pkDropRate = e.Info.DropRate
					// pkDropMaxNum = e.Info.DropMaxNum
					// pkDropDayLimit = e.Info.DropDayLimit
				}
			}

			lock.Lock()
			OpenEvents = events

			CosCmtAlwaysReward = info.CosCmtAlwaysReward
			CosMultiStartVIP = info.CosMultiStartVIP
			CosMultiStartLevel = info.CosMultiStartLevel
			CosMultiEndVIP = info.CosMultiEndVIP
			CosMultiEndLevel = info.CosMultiEndLevel

			//event drop
			EventDrop = info.EventDrop
			EventDropID = info.EventDropID
			EventDropRate = info.EventDropRate
			EventDropMaxNum = info.EventDropMaxNum
			EventDropDayLimit = info.EventDropDayLimit
			EventDropType = dt

			//sum game content
			LogGame = info.LogGame
			LogGameEventID = info.LogGameEventID
			LogGameEndTime = info.LogGameEndTime

			//sum diamond consume
			LogDiamondConsume = info.LogDiamondConsume
			LogDiamondConsumeEventID = info.LogDiamondConsumeEventID

			PKMulti = info.PKMulti
			PKMultiStartVIP = info.PKMultiStartVIP
			PKMultiStartLevel = info.PKMultiStartLevel
			PKMultiEndVIP = info.PKMultiEndVIP
			PKMultiEndLevel = info.PKMultiEndLevel

			specialMissionStartTime = smst
			specialMissionEndTime = smet
			specialMissionMap = smm
			specialMissionEventID = smeid
			Info("special mission:", specialMissionStartTime, specialMissionEndTime, specialMissionMap)

			pkDropEventID = pkeid
			pkDropStartTime = pkdst
			pkDropEndTime = pkdet

			MissionBeilvConfig = info.MissionBeilvConfig
			MissionBeilvMultiStartVIP = info.MissionBeilvMultiStartVIP
			MissionBeilvMultiStartLevel = info.MissionBeilvMultiStartLevel
			MissionBeilvMultiEndVIP = info.MissionBeilvMultiEndVIP
			MissionBeilvMultiEndLevel = info.MissionBeilvMultiEndLevel

			Info("multi config:", MissionBeilvConfig, "log diamond consume:", LogDiamondConsume, "log diamond consume event id:", LogDiamondConsumeEventID)
			lock.Unlock()
		}

		//every 1 minute
		time.Sleep(60 * time.Second)
	}
}

func PKDrop() (bool, int64) {
	lock.RLock()
	defer lock.RUnlock()

	nt := time.Now().Unix()
	if pkDropStartTime <= nt && pkDropEndTime >= nt && pkDropEventID > 0 {
		return true, pkDropEventID
	}
	return false, 0
}

func SpecialMission(lid string, t time.Time) (bool, int64) {
	lock.RLock()
	defer lock.RUnlock()

	nt := t.Unix()
	if specialMissionStartTime <= nt && nt <= specialMissionEndTime {
		return specialMissionMap[lid] > 0, specialMissionEventID
	}

	return false, 0
}

func GetTinyPackInfo() (int64, int, int64, int, int, int, string) {
	lock.RLock()
	defer lock.RUnlock()

	for _, e := range OpenEvents {
		if e.Info.Type == 10 && e.EndTime-time.Now().Unix() >= 10 {
			return e.ID, e.Info.TinyReward.RMBCost, e.Info.TinyReward.Gold, e.Info.TinyReward.Diamond, e.Info.TinyReward.Tili, e.Info.TinyReward.MaxBuyCnt, e.Info.TinyReward.Clothes
		}
	}
	return 0, 0, 0, 0, 0, 0, ""
}

func GetDefaultTinyPackInfo() (int64, int, int64, int, int, int, string) {
	lock.RLock()
	defer lock.RUnlock()

	for _, e := range OpenEvents {
		if e.Info.Type == 16 && e.EndTime-time.Now().Unix() >= 10 {
			return e.ID, e.Info.TinyReward.RMBCost, e.Info.TinyReward.Gold, e.Info.TinyReward.Diamond, e.Info.TinyReward.Tili, e.Info.TinyReward.MaxBuyCnt, e.Info.TinyReward.Clothes
		}
	}
	return 0, 0, 0, 0, 0, 0, ""
}

func GetFirstIAPInfo() int64 {
	lock.RLock()
	defer lock.RUnlock()

	for _, e := range OpenEvents {
		if e.Info.Type == 12 {
			return e.ID
		}
	}
	return 0
}

func GetSumIAPInfo() int64 {
	lock.RLock()
	defer lock.RUnlock()

	for _, e := range OpenEvents {
		if e.Info.Type == 17 {
			return e.ID
		}
	}
	return 0
}

func getOpenEventsFromDB() (res []Event, dt int, info EventInfoDetail, err error) {
	info.PKMulti = 1
	now := time.Now().Unix()
	info.MissionBeilvConfig = make(map[string]float64)

	var rows *sql.Rows
	if rows, err = dao.QueryEventsStmt.Query(now, now); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]Event, 0)

		for rows.Next() {
			e := Event{}
			var infos string
			if err = rows.Scan(&e.ID, &e.Type, &e.Title, &e.Content, &e.StartTime, &e.EndTime, &e.Banner, &e.RedirectType, &e.AutoPop, &e.RefreshTime, &e.ContentType, &e.IsLong, &infos); err != nil {
				Err(err)
				return
			}

			if err = json.Unmarshal([]byte(infos), &e.Info); err != nil {
				Err(err)
				return
			}

			if e.Info.Type == 4 || e.Info.Type == 15 {
				happen := false
				nhm := time.Now().Format("1504")
				for _, tg := range e.Info.TimeGap {
					tgs := strings.Split(tg, "-")
					tg1 := tgs[0]
					tg2 := tgs[1]
					if nhm >= tg1 && nhm <= tg2 {
						happen = true
						break
					}
				}

				if happen {
					if e.Info.MultiType == 2 {
						info.PKMulti = e.Info.Multi
						info.PKMultiStartVIP = e.Info.VIPStart
						info.PKMultiStartLevel = e.Info.LevelStart
						info.PKMultiEndVIP = e.Info.VIPEnd
						info.PKMultiEndLevel = e.Info.LevelEnd
					} else if e.Info.MultiType == 1 {
						info.CosCmtAlwaysReward = true
						info.CosMultiStartVIP = e.Info.VIPStart
						info.CosMultiStartLevel = e.Info.LevelStart
						info.CosMultiEndVIP = e.Info.VIPStart
						info.CosMultiEndLevel = e.Info.LevelEnd
					} else {
						info.MissionBeilvConfig[fmt.Sprintf("%d:%d", e.Info.MultiType, e.Info.MultiRes)] = e.Info.Multi
						info.MissionBeilvMultiStartVIP = e.Info.VIPStart
						info.MissionBeilvMultiStartLevel = e.Info.LevelStart
						info.MissionBeilvMultiEndVIP = e.Info.VIPEnd
						info.MissionBeilvMultiEndLevel = e.Info.LevelEnd
					}
				}
			} else if e.Info.Type == 7 || e.Info.Type == 1000 {
				info.EventDrop = true
				info.EventDropID = e.ID
				info.EventDropRate = e.Info.DropReward.DropRate
				info.EventDropMaxNum = e.Info.DropReward.DropMaxNum
				info.EventDropDayLimit = e.Info.DropReward.DropDayLimit
				dt = e.Info.Type
			} else if e.Info.Type == 8 {
				info.LogDiamondConsume = true
				info.LogDiamondConsumeEventID = e.ID
			} else if e.Info.Type == 9 {
				info.LogGame = true
				info.LogGameEventID = e.ID
				info.LogGameEndTime = e.EndTime
			}

			res = append(res, e)
		}
	}

	return
}

func GetLogGameInfo() (bool, int64, int64) {
	lock.RLock()
	defer lock.RUnlock()

	return LogGame, LogGameEventID, LogGameEndTime
}

func GetLogDiamondConsumeInfo() (bool, int64) {
	lock.RLock()
	defer lock.RUnlock()
	return LogDiamondConsume, LogDiamondConsumeEventID
}

func GetEventDropDayLimit() int {
	lock.RLock()
	defer lock.RUnlock()

	return EventDropDayLimit
}

func GetEventDropNum() (int64, int, int) {
	lock.RLock()
	defer lock.RUnlock()

	if EventDrop {
		v := rand.Intn(100)
		if v <= EventDropRate-1 {
			//drop
			return EventDropID, EventDropType, rand.Intn(EventDropMaxNum) + 1
		}
	}
	return EventDropID, EventDropType, 0
}

func GetOpenEvents(channel string) (res []Event) {
	lock.RLock()
	defer lock.RUnlock()

	if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
		res = OpenEvents
	} else {
		if channel == "" {
			res = OpenEvents
		} else {
			res = make([]Event, 0)
			for _, e := range OpenEvents {
				if e.Info.Type != 16 {
					res = append(res, e)
				}
			}
		}
	}

	now := time.Now().Unix()
	if tp := params.GetTimeDoor(); tp.Kor_shikongzhimen_start_time <= now && tp.Kor_shikongzhimen_end_time > now {
		tpEvent := Event{
			ID: constants.TIME_DOOR_EVENT_ID,

			StartTime: tp.Kor_shikongzhimen_start_time,
			EndTime:   tp.Kor_shikongzhimen_end_time,

			Banner:       tp.BannerLocation,
			RedirectType: -1,
			AutoPop:      true,

			Info: EventInfo{
				Type:          18,
				TimeDoorParam: tp,
				ShowInList:    true,
			},
		}

		res = append(res, tpEvent)
	}

	return
}

func GetOpenEventByID(eid int64) (e Event) {
	lock.RLock()
	defer lock.RUnlock()

	now := time.Now().Unix()
	for _, et := range OpenEvents {
		if et.ID == eid && et.StartTime <= now && et.EndTime >= now {
			return et
		}
	}
	return
}

func EventOpen(eid int64) bool {
	lock.RLock()
	defer lock.RUnlock()

	now := time.Now().Unix()
	for _, et := range OpenEvents {
		if et.ID == eid && et.StartTime <= now && et.EndTime >= now {
			return true
		}
	}

	return false
}

func AddEventPartyLog(key string, username string) {
	if log, id, et := GetLogGameInfo(); log {
		if _, err := dao.AddEventPartyLogStmt.Exec(id, key, username, et, 1); err != nil {
			Err(err)
		}
	}
}

func QueryEventPartyUserLegal(eid int64, key string, username string) bool {
	keys := strings.Split(key, ":")
	if rows, err := dao.QueryEventPartyUserLogStmt.Query(username, time.Now().Unix()); err != nil {
		Err(err)
		return false
	} else {
		defer rows.Close()
		for rows.Next() {
			var eventID int64
			var k string
			var cnt int
			if err = rows.Scan(&eventID, &k, &cnt); err != nil {
				Err(err)
				return false
			}

			if eventID == eid {
				ks := strings.Split(k, ":")
				if len(keys) == 6 && len(ks) == 5 && keys[0] == ks[0] && keys[1] == ks[1] && keys[2] == ks[2] {
					if !compareLessOrEqual(keys[3], ks[3]) {
						//player_num
						return false
					} else if !compareLessOrEqual(ks[4], keys[4]) {
						//rank
						return false
					} else if !compareLessOrEqual(keys[5], strconv.Itoa(cnt)) {
						return false
					}

					return true
				}
			}
		}
		return false
	}
}

func QueryEventPartyLog(username string) (res map[string]int) {
	res = make(map[string]int)
	if rows, err := dao.QueryEventPartyUserLogStmt.Query(username, time.Now().Unix()); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var eid int64
			var k string
			var cnt int
			if err = rows.Scan(&eid, &k, &cnt); err != nil {
				Err(err)
				return
			}

			res[k] = cnt
		}
		return
	}
}

func compareLessOrEqual(less string, more string) bool {
	li, _ := strconv.Atoi(less)
	mi, _ := strconv.Atoi(more)

	if li <= mi {
		return true
	}

	return false
}

func DeleteEventPartyLog() (err error) {
	if _, err = dao.DeleteEventPartyLogStmt.Exec(time.Now().Unix() - 1800); err != nil {
		Err(err)
	}
	return
}

func GetPKMulti(vip, level int) float64 {
	lock.RLock()
	defer lock.RUnlock()

	if PKMultiStartVIP > 0 && vip < PKMultiStartVIP {
		return 1
	}

	if PKMultiStartLevel > 0 && level < PKMultiStartLevel {
		return 1
	}

	if PKMultiEndVIP > 0 && vip > PKMultiEndVIP {
		return 1
	}

	if PKMultiEndLevel > 0 && level > PKMultiEndLevel {
		return 1
	}

	return PKMulti
}

func IsCosCmtAlwaysReward(vip, level int) bool {
	lock.RLock()
	defer lock.RUnlock()

	if CosMultiStartVIP > 0 && vip < CosMultiStartVIP {
		return false
	}

	if CosMultiStartLevel > 0 && level < CosMultiStartLevel {
		return false
	}

	if CosMultiEndVIP > 0 && vip > CosMultiEndVIP {
		return false
	}

	if CosMultiEndLevel > 0 && level > CosMultiEndLevel {
		return false
	}

	return CosCmtAlwaysReward
}

func GetRecordMulti(rtypee int, vip, level int) (g float64, t float64, cp float64, exp float64) {
	lock.RLock()
	defer lock.RUnlock()

	g = 1
	t = 1
	cp = 1
	exp = 1
	if MissionBeilvMultiStartVIP > 0 && vip < MissionBeilvMultiStartVIP {
		return
	}

	if MissionBeilvMultiStartLevel > 0 && level < MissionBeilvMultiStartLevel {
		return
	}

	if MissionBeilvMultiEndVIP > 0 && vip > MissionBeilvMultiEndVIP {
		return
	}

	if MissionBeilvMultiEndLevel > 0 && level > MissionBeilvMultiEndLevel {
		return
	}

	g = MissionBeilvConfig[fmt.Sprintf("%d:1", rtypee)]
	if g == 0 {
		g = 1
	}

	t = MissionBeilvConfig[fmt.Sprintf("%d:2", rtypee)]
	if t == 0 {
		t = 1
	}

	cp = MissionBeilvConfig[fmt.Sprintf("%d:3", rtypee)]
	if cp == 0 {
		cp = 1
	}

	exp = MissionBeilvConfig[fmt.Sprintf("%d:4", rtypee)]
	if exp == 0 {
		exp = 1
	}

	return
}
