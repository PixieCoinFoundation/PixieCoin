package types

import (
	. "model"
)

type Event struct {
	ID   int64
	Type int //0:event 1:announcement

	Title   string
	Content string

	StartTime int64
	EndTime   int64

	Banner string

	//0:order page 1:content page 2:pk page 3:cosplay page 4:tili page 5:clothes page 6:mission reward page
	//101:theater 1
	//102:theater 2
	//103:theater 3
	//104:theater 4
	//105:theater 5
	RedirectType int
	AutoPop      bool
	RefreshTime  int64

	ContentType int //0:text 1:img location
	IsLong      bool

	Info EventInfo
}

type EventInfo struct {
	//1:mission reward 2:day&time tili 3:day clothes 4:game multi 5.month signin 6.new player 7day 7.mission drop
	//8:sum consume reward 9:sum game content reward 10:tiny pack 11:friend cnt 12:first iap 13:month card 14:first iap double
	//15:beilv 16:default tiny pack 17:sum diamond iap 18:time door bonus

	//1000:normal magic beed clothes
	//1001:pk magic beed clothes
	//1002:special mission magic beed clothes
	Type int

	TimeGap            []string          //["1200-1400","1800-2000"]
	MultiType          int               //1:cosplay 2:pk 3:chuzhong 4:gaozhong 5:daxue 6:shehui 7:juchang
	Multi              float64           //*1.5 *2 *3
	MultiRes           int               //1:gold 2:tili 3:clo piece 4:exp 5:pk point
	ClothesReward      map[string]string //progress_id:clothes_id,type 1:mission_id:clothes_id,type 3:progress:clothes_id
	TiliReward         int               //type 2.["DEFAULT":30]
	GoldReward         map[string]int    //type 1.level_id:gold
	DiamondReward      map[string]int    //type 1.level_id:diamond
	TiliRewards        map[string]int
	ClothesPieceReward map[string]string
	CosReward          map[string]string
	PkReward           map[string]string
	CombineReward      map[string]string
	LevelReward        map[string]string
	SReward            map[string]string
	DoubleConfig       map[string]int //月签到专用 天数:双倍所需最小VIP等级

	SpecialLevels []GFLevelS //1002 config

	//party key: J/C:X/J:hour:player_num:rank
	//J:X:0:100:0 host a party with player num>=100
	//J:X:3:0:0 host a 3-hour party
	//C:X:0:0:20 join a casual party and rank in top 20
	//C:J:3:0:0 join a 3-hour prize party
	PartyReward map[string]string
	DropReward  DropInfo
	ShowInList  bool
	LinkAddress string
	TinyReward  TinyRewardInfo
	LevelStart  int
	LevelEnd    int
	VIPStart    int
	VIPEnd      int

	TimeDoorParam TimeDoor
}

type TinyRewardInfo struct {
	RMBCost   int
	Gold      int64
	Diamond   int
	Tili      int
	Clothes   string
	MaxBuyCnt int
}

type DropInfo struct {
	DropRate     int
	DropMaxNum   int
	DropDayLimit int
	Reward       map[string]int
}

type RawEvent struct {
	ID    int64
	Type  int //0:event 1:announcement
	EType int

	Title   string
	Content string

	ST        int64
	StartTime string
	ET        int64
	EndTime   string

	Banner       string
	RedirectType int //0:order page 1:content page
	AutoPop      bool
	RefreshTime  int64

	ContentType int //0:text 1:img location
	IsLong      bool

	Info string

	Platform string
}

type RollingAnnouncement struct {
	ID         int64
	Content    string
	StartTime  int64
	EndTime    int64
	StartTimes string
	EndTimes   string
}

type RollingAnnoun struct {
	ID        int64
	Content   string
	StartTime int64
	EndTime   int64
}
