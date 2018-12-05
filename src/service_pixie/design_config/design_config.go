package design_config

import (
	"appcfg"
	"fmt"
	. "pixie_contract/api_specification"
	"reflect"
	"strings"
	"tools"
)

var designConfig DesignConfigInfo

type DesignConfigInfo struct {
	//评级对应经验
	SRankExp int
	ARankExp int
	BRankExp int

	//评级对应npc友好度
	SNpcFriendshipExpAdd int
	ANpcFriendshipExpAdd int
	BNpcFriendshipExpAdd int

	//回复一点体力需要的时间
	MinutePerEnergy int64 //比如 5

	//人流容量浮动上下限
	VisitCapFloor float64 //比如 0.9
	VisitCapUp    float64 //比如 1.1

	//初始需求能力
	RestaurantOfferInit float64
	EntertainOfferInit  float64
	CleanOfferInit      float64

	//街区收益最终浮动参数
	StreetIncomeFinalAdjust float64 //比如 0.0001
}

func init() {
	if appcfg.GetServerType() != "" {
		return
	}

	loadDesignConfig()
}

func GetSecondsPerTili() int64 {
	return designConfig.MinutePerEnergy * 60
}

func GetExpAdd(rank string) int {
	if rank == string(RECORD_RANK_S) {
		return designConfig.SRankExp
	} else if rank == string(RECORD_RANK_A) {
		return designConfig.ARankExp
	} else if rank == string(RECORD_RANK_B) {
		return designConfig.BRankExp
	}

	return 0
}

func GetNpcFriendshipExpAdd(rank string) int {
	if rank == string(RECORD_RANK_S) {
		return designConfig.SNpcFriendshipExpAdd
	} else if rank == string(RECORD_RANK_A) {
		return designConfig.ANpcFriendshipExpAdd
	} else if rank == string(RECORD_RANK_B) {
		return designConfig.BNpcFriendshipExpAdd
	}

	return 0
}

func GetRestaurantOfferInit() float64 {
	return designConfig.RestaurantOfferInit
}

func GetEntertainOfferInit() float64 {
	return designConfig.EntertainOfferInit
}

func GetCleanOfferInit() float64 {
	return designConfig.CleanOfferInit
}

func GetStreetIncomeFinalAdjust() float64 {
	return designConfig.StreetIncomeFinalAdjust
}

func loadDesignConfig() {
	path := "scripts/data_pixie/parameter.csv"
	records := tools.LoadCSV(path)

	vt := reflect.ValueOf(designConfig)
	fieldCnt := vt.NumField()

	configCnt := 0
	for k, v := range records {
		if v[1] == "0" || v[1] == "" {
			panic(fmt.Sprintf("design config wrong %d %s", k, v[1]))
		}

		configCnt++
		key := strings.TrimSpace(v[0])
		kb := tools.IgnoreUTF8BOM([]byte(key))
		key = string(kb)

		if key == "S-exp" {
			designConfig.SRankExp = tools.LoadGetInt(v[1], "s exp")
		} else if key == "A-exp" {
			designConfig.ARankExp = tools.LoadGetInt(v[1], "a exp")
		} else if key == "B-exp" {
			designConfig.BRankExp = tools.LoadGetInt(v[1], "b exp")
		} else if key == "S-friendship" {
			designConfig.SNpcFriendshipExpAdd = tools.LoadGetInt(v[1], "s friend add")
		} else if key == "A-friendship" {
			designConfig.ANpcFriendshipExpAdd = tools.LoadGetInt(v[1], "a friend add")
		} else if key == "B-friendship" {
			designConfig.BNpcFriendshipExpAdd = tools.LoadGetInt(v[1], "b friend add")
		} else if key == "minutes-per-energy" {
			designConfig.MinutePerEnergy = tools.LoadGetInt64(v[1], "minutes per energy")
		} else if key == "passager-random-floor" {
			designConfig.VisitCapFloor = tools.LoadGetFloat64(v[1], "visit random floor")
		} else if key == "passager-random-upper" {
			designConfig.VisitCapUp = tools.LoadGetFloat64(v[1], "visit random up")
		} else if key == "restaurant-init" {
			designConfig.RestaurantOfferInit = tools.LoadGetFloat64(v[1], "restaurant offer init")
		} else if key == "entertainment-init" {
			designConfig.EntertainOfferInit = tools.LoadGetFloat64(v[1], "entertain offer init")
		} else if key == "clean-init" {
			designConfig.CleanOfferInit = tools.LoadGetFloat64(v[1], "clean offer init")
		} else if key == "earning-correction" {
			designConfig.StreetIncomeFinalAdjust = tools.LoadGetFloat64(v[1], "street income final adjust")
		} else {
			panic(fmt.Sprintf("unknown design config #%d# #%s# #%s#", k, key, v[0]))
		}
	}

	if configCnt != fieldCnt {
		panic(fmt.Sprintf("design config field and file not match %d %d", fieldCnt, configCnt))
	}
}
