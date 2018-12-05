package types

import (
	. "pixie_contract/api_specification"
)

type StatusDetailName string

type GMStatus struct {
	//base info
	Username            string `bson:"username" json:"-"`
	UId                 string `bson:"uid" json:"uid"`
	Nickname            string `bson:"nickname" json:"nickname"`
	Money               int    `bson:"money" json:"money"`
	Diamond             int    `bson:"diamond" json:"diamond"`
	Tili                int
	Head                string
	HeadFile            string
	MaxTili             int `json:"mac_tili"`
	TiliPri             string
	BuyDiamond          int `bson:"buyDiamond" json:"buyDiamond"`
	VIP                 int
	CheatCount          int     `bson:"cheatCount" json:"cheatCount"`
	Ban                 bool    `bson:"ban" json:"ban"`
	BanInfo             BanData `bson:"banInfo" json:"banInfo"`
	DeviceToken         string  `bson:"deviceToken" json:"deviceToken"` // 推送的deviceToken
	Badge               int     `bson:"badge" json:"badge"`             // 推送时，iOS显示在角落的数字
	OpTime              string  `bson:"opTime" json:"opTime"`           // 最近一次的操作时间
	CurrentMissionLevel int
	CurrentMissionNo    int
	PkPoint             int
	CosPoint            int
	DesignerPoint       int
	PKWinCount          int

	MonthCardStatus string

	Backgrounds []GFBackground
	//付费背景列表
	BGlist map[string]*GFBackground
	//item info
	Items   []Item
	PKCount int

	//records
	Records []Record

	//tasks
	Tasks []Task

	//hopes
	Hopes []HopeView

	HopePP []HopeExtra

	//board msgs
	BoardMsgs []BoardMsgView

	//clothes
	Clothes []GMClothesInfo

	//clothes piece
	ClothesPieces []GMClothesPieceInfo

	//mails
	Mails []SMail

	//helps
	Helps []GMHelp

	// PartyInfo PartyStatus

	Info string

	ClothCnt  int
	ButtonCnt int

	Platform     string
	ExtraInfo1   Extra1
	DayExtraInfo Extra

	TheaterStatus string
	EndingStatus  string

	RegTimes        string
	BlockStartTimes string
	BlockEndTimes   string
	DelTimes        string

	Deleting bool
	Deleted  bool

	ClothesDesc string

	GuildID       int64
	GuildName     string
	GuildActivity int
	GuildShareCnt int

	PayDiamond  int
	FreeDiamond int

	BUsername         string
	LastSendHopeTimes string

	IsDesigner      bool
	DesignerDiamond int64
	DesignerGold    int64
	DesignerDesc    string

	Magic             int
	UnlockedModelDesc string
}

const (
	STATUS_DETAIL_STREET_OPERATE StatusDetailName = "STREET_OPERATE"

	STATUS_DETAIL_DAY_EXTRA_TILI   StatusDetailName = "DAY_EXTRA_TILI"
	STATUS_DETAIL_DAY_EXTRA        StatusDetailName = "DAY_EXTRA"
	STATUS_DETAIL_DAY_EXTRA_CLO    StatusDetailName = "DAY_EXTRA_CLO"
	STATUS_DETAIL_DAY_EXTRA_CLOP   StatusDetailName = "DAY_EXTRA_CLP"
	STATUS_DETAIL_DAY_EXTRA_MNYDMD StatusDetailName = "DAY_EXTRA_MNYDMD"
	STATUS_DETAIL_HELP_FOLLOW      StatusDetailName = "HELP_FOLLOW"
	STATUS_DETAIL_ITEMS            StatusDetailName = "ITEMS"
	STATUS_DETAIL_ITEMS_DAY_EXTRA  StatusDetailName = "ITEMS_DAY_EXTRA"
	STATUS_DETAIL_TASKS            StatusDetailName = "TASKS"
	STATUS_DETAIL_CLO_MNY_DMD      StatusDetailName = "CLO_MNY_DMD"
	STATUS_DETAIL_CLO_TASK_MNY_DMD StatusDetailName = "CLO_TASK_MNY_DMD"
	// STATUS_DETAIL_CLO_DMD               StatusDetailName = "CLO_DMD"
	STATUS_DETAIL_CLO_MNY               StatusDetailName = "CLO_MNY"
	STATUS_DETAIL_CLO_TASK_DMD          StatusDetailName = "CLO_TASK_DMD"
	STATUS_DETAIL_CLO_TASK_MNY          StatusDetailName = "CLO_TASK_MNY"
	STATUS_DETAIL_MNY_DMD               StatusDetailName = "MNY_DMD"
	STATUS_DETAIL_ENDING                StatusDetailName = "ENDING"
	STATUS_DETAIL_SCRIPT                StatusDetailName = "SCRIPT"
	STATUS_DETAIL_TILI                  StatusDetailName = "TILI"
	STATUS_DETAIL_PK_WINCNT_PNT         StatusDetailName = "PK_WINCNT_PNT"
	STATUS_DETAIL_HEART                 StatusDetailName = "HEART"
	STATUS_DETAIL_HEART_TILI            StatusDetailName = "HEART_TILI"
	STATUS_DETAIL_HEAD                  StatusDetailName = "HEAD"
	STATUS_DETAIL_VIP                   StatusDetailName = "VIP"
	STATUS_DETAIL_EVERYDAY              StatusDetailName = "EVERYDAY"
	STATUS_DETAIL_CLO                   StatusDetailName = "CLO"
	STATUS_DETAIL_RECORD                StatusDetailName = "RECORD"
	STATUS_DETAIL_RECORD_TASK           StatusDetailName = "RECORD_TASK"
	STATUS_DETAIL_RECORD_RANK           StatusDetailName = "RECORD_RANK"
	STATUS_DETAIL_RECORD_RANK_DAY_EXTRA StatusDetailName = "RECORD_RANK_DAY_EXTRA"
	STATUS_DETAIL_PLANET                StatusDetailName = "PLANET"
	STATUS_DETAIL_SUITS                 StatusDetailName = "SUITS"
	STATUS_DETAIL_S_SUITS               StatusDetailName = "S_SUITS"
	STATUS_DETAIL_S_SUITS_DAY_EXTRA     StatusDetailName = "S_SUITS_DAY_EXTRA"
	STATUS_DETAIL_WEIBO                 StatusDetailName = "WEIBO"
	STATUS_DETAIL_BAN                   StatusDetailName = "BAN"
	STATUS_DETAIL_CONTACT               StatusDetailName = "CONTACT"
	STATUS_DETAIL_TAG                   StatusDetailName = "TAG"
	STATUS_DETAIL_CLOP                  StatusDetailName = "CLOP"
	STATUS_DETAIL_PKCNT                 StatusDetailName = "PK_CNT"
	STATUS_DETAIL_MAIL_GIFT             StatusDetailName = "MAIL_GIFT"
	STATUS_DETAIL_PARTY_INFO            StatusDetailName = "PARTY_INFO"
	STATUS_DETAIL_PARTY_INFO_DAY_EXTRA  StatusDetailName = "PARTY_INFO_DAY_EXTRA"
	STATUS_DETAIL_DAY_EXTRA_ALL         StatusDetailName = "DAY_EXTRA_ALL"
	STATUS_DETAIL_PKCNT_DAY_EXTRA       StatusDetailName = "PK_CNT_DAY_EXTRA"
	STATUS_DETAIL_PK_WINCNT_PNT_EXTRA   StatusDetailName = "PK_WINCNT_PNT_EXTRA"
	STATUS_DETAIL_E1_MNY_DMD            StatusDetailName = "E1_MNY_DMD"
	STATUS_DETAIL_HEART_TASK            StatusDetailName = "HEART_TASK"
	STATUS_DETAIL_HOPE                  StatusDetailName = "HOPE"
	STATUS_DETAIL_E1_CLO                StatusDetailName = "E1_CLO"
	STATUS_DETAIL_TASK_REWARD           StatusDetailName = "TASK_REWARD"
	STATUS_DETAIL_LOTTERY               StatusDetailName = "LOTTERY"
	STATUS_DETAIL_E1_MNY_DMD_DAY_EXTRA  StatusDetailName = "E1_MNY_DMD_DAY_EXTRA"
)
