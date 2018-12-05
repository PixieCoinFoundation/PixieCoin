package types

import (
// "gopkg.in/mgo.v2/bson"
)

// type CosItemSex int

// const (
// 	COS_ITEM_BOY  CosItemSex = 1
// 	COS_ITEM_GIRL CosItemSex = 2
// )

type CosItemScore struct {
	Username  string `bson:"username" json:"username"`
	Score     int    `bson:"score" json:"score"`
	CosItemID int64  `bson:"cosItemID" json:"cosItemID"`
}

type CosItem struct {
	Username  string            `bson:"username" json:"username, omitempty"`
	UID       string            `bson:"uid" json:"uid, omitempty"`
	Nickname  string            `bson:"nickname" json:"nickname, omitempty"`
	CosplayID int64             `bson:"cosplayID" json:"cosplayID, omitempty"`
	ItemID    int64             `bson:"itemID" json:"itemID, omitempty"`       // 系统分配给该套服装的编号
	Title     string            `bson:"title" json:"title, omitempty"`         // 玩家取的搭配名字
	Desc      string            `bson:"desc" json:"desc, omitempty"`           // 玩家给搭配的描述
	ModelNo   int               `bson:"modelNo" json:"modelNo, omitempty"`     // 搭配的性别
	Clothes   map[string]string `bson:"clothes" json:"clothes, omitempty"`     // 搭配的服装列表
	Scores    []CosItemScore    `bson:"scores" json:"-"`                       // 搭配的服装分数列表
	Score     int64             `bson:"score" json:"score, omitempty"`         // 搭配的服装分数
	SysScore  int64             `bson:"sysScore" json:"sysScore, omitempty"`   // 搭配的服装分数
	Icon      string            `bson:"icon" json:"icon, omitempty"`           // 搭配的服装预览小图地址
	Head      string            `bson:"head" json:"head, omitempty"`           // 玩家上传这个服装时的head
	CosBg     string            `bson:"cosBg" json:"cosBg, omitempty"`         // 当前cos的背景图
	CosPoints int               `bson:"cosPoints" json:"cosPoints, omitempty"` // cosPoints

	UploadTime int64 `bson:"uploadTime" json:"uploadTime, omitempty"` // 服装上传的时间
	SetTopTime int64 `bson:"setTopTime" json:"setTopTime, omitempty"` // 服装置顶的时间

	// 临时变量，不做存储
	HasVoted bool `bson:"-" json:"hasVoted"` // 是否已经投票过
	MyScore  int  `bson:"-" json:"myScore"`  // 我投的分数

	Rank   int                      `bson:"rank" json:"rank"` // 排名
	Img    string                   `bson:"img" json:"img"`
	DyeMap map[string][7][4]float64 `bson:"dyeMap" json:"dyeMap"`
}

type FCosItem struct {
	ItemID   int64
	Username string
	Icon     string
	Img      string
}

type CosItemR struct {
	Username   string `bson:"username" json:"username, omitempty"`
	UID        string `bson:"uid" json:"uid, omitempty"`
	Nickname   string `bson:"nickname" json:"nickname, omitempty"`
	CosplayID  int64  `bson:"cosplayID" json:"cosplayID, omitempty"`
	ItemID     int64  `bson:"itemID" json:"itemID, omitempty"`         // 系统分配给该套服装的编号
	Title      string `bson:"title" json:"title, omitempty"`           // 玩家取的搭配名字
	Desc       string `bson:"desc" json:"desc, omitempty"`             // 玩家给搭配的描述
	ModelNo    int    `bson:"modelNo" json:"modelNo, omitempty"`       // 搭配的性别
	Clothes    string `bson:"clothes" json:"clothes, omitempty"`       // 搭配的服装列表
	SysScore   int64  `bson:"sysScore" json:"sysScore, omitempty"`     // 搭配的服装分数
	Icon       string `bson:"icon" json:"icon, omitempty"`             // 搭配的服装预览小图地址
	Head       string `bson:"head" json:"head, omitempty"`             // 玩家上传这个服装时的head
	CosBg      string `bson:"cosBg" json:"cosBg, omitempty"`           // 当前cos的背景图
	CosPoints  int    `bson:"cosPoints" json:"cosPoints, omitempty"`   // cosPoints
	UploadTime int64  `bson:"uploadTime" json:"uploadTime, omitempty"` // 服装上传的时间
	Score      int
	Rank       int
	Img        string `bson:"img" json:"img, omitempty"`
	DyeMap     string `bson:"dyeMap" json:"dyeMap"`
}

type CosComment struct {
	CosplayID int64
	CosItemID int64  `bson:"cosItemID" json:"cosItemID, omitempty"`
	Username  string `bson:"username" json:"username, omitempty"`
	Nickname  string `bson:"nickname" json:"nickname, omitempty"`
	Content   string `bson:"content" json:"content, omitempty"`
	ReplyTime int64  `bson:"replytime" json:"replytime, omitempty"`
}
