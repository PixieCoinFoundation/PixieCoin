package types

import (
	. "pixie_contract/api_specification"
)

// "gopkg.in/mgo.v2/bson"

type Mail struct {
	MailID     int64         `bson:"mailID" json:"mailID"`
	Type       int           `bson:"type" json:"type"`
	From       string        `bson:"from" json:"from"`
	To         string        `bson:"to" json:"to"`
	Title      string        `bson:"title" json:"title"`
	Content    string        `bson:"content" json:"content"`
	Gold       int64         `bson:"gold" json:"gold"`
	Diamond    int64         `bson:"diamond" json:"diamond"`
	Clothes    []ClothesInfo `bson:"clothes" json:"clothes"`
	Time       int64         `bson:"time" json:"time"`
	Read       bool          `bson:"read" json:"read"`
	Delete     bool          `bson:"Delete",json:"Delete"`
	ExpireTime int64         `bson:"expireTime",json:"expireTime"`
}

// MailExtraLog 邮件信息
type MailRqeExtra struct {
	Mail
	ReadTime string `json:"readtime"`
}

type DMail struct {
	MailID     int64
	To         string
	Time       int64
	ExpireTime int64
}

type MailContent struct {
	Content              string
	Item101Cnt           int
	Item102Cnt           int
	Item103Cnt           int
	Item104Cnt           int
	Item105Cnt           int
	Item106Cnt           int
	TiliCnt              int
	CloPieces            map[string]int
	JuhuaCnt             int
	JuhuaChaCnt          int
	Cloth                int
	Button               int
	Wire                 int
	Needle               int
	DesignCoin           int64
	DesignCoinPriceToken int

	//客户端忽略的值
	//竞拍广告位需要计入活动的钻石数
	ADLogDiamond int

	FailClothesList []SimpleClothes1

	//party related
	PartyInfo       Party
	PartyUsername   string
	PartyNickname   string
	InvitePartyType int //0:flower 1:attend 2:partner
}

type SMail struct {
	MailID  int64
	From    string
	Content string
	Time    string
	Gold    int
	Diamond int
	Clothes string
}

type MailInfo struct {
	From         string
	Title        string
	Content      string
	To           string
	Diamond      int64
	Gold         int64
	Clothes      []ClothesInfo
	ShouldDelete bool
}

type CronMail struct {
	Id          int64
	Title       string
	Content     string
	Gold        int64
	Diamond     int64
	Clothes     string
	CronTime    int64
	MinRegTime  int64
	MaxRegTime  int64
	Delete      bool
	Status      int
	Info        string
	CronTimeStr string
	Platform    string

	ChannelLimit bool
	ThirdChannel string
}
