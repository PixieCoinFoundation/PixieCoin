package types

import (
// "gopkg.in/mgo.v2/bson"
)

type CosplayStatus int

const (
	COS_OPEN       CosplayStatus = 0 //opening
	COS_PROCESSING CosplayStatus = 1 //processing
	COS_CLOSED     CosplayStatus = 2 //closed
	COS_ALIVE      CosplayStatus = 3 //alive cosplay
)

// type CosplayType int

// const (
// 	NORMAL  CosplayType = 0 // 初级
// 	ADVANCE CosplayType = 1 // 进阶
// )

type CosParams struct {
	Warm      float64 `bson:"warm" json:"warm"`
	Formal    float64 `bson:"formal" json:"formal"`
	Tight     float64 `bson:"tight" json:"tight"`
	Bright    float64 `bson:"bright" json:"bright"`
	Dark      float64 `bson:"dark" json:"dark"`
	Cute      float64 `bson:"cute" json:"cute"`
	Man       float64 `bson:"man" json:"man"`
	Tough     float64 `bson:"tough" json:"tough"`
	Noble     float64 `bson:"noble" json:"noble"`
	Strange   float64 `bson:"strange" json:"strange"`
	Sexy      float64 `bson:"sexy" json:"sexy"`
	Sport     float64 `bson:"sport" json:"sport"`
	Adjust    float64 `bson:"adjust" json:"adjust"`
	StarScore int64   `bson:"starScore" json:"starScore"`
}

type Cosplay struct {
	CosplayID int64         `bson:"cosplayID" json:"cosplayID"` // 此次活动的编号
	Title     string        `bson:"title" json:"title"`         // 此次活动的标题
	Keyword   string        `bson:"keyword" json:"keyword"`     // 此次活动的关键字
	Type      int           `bson:"type" json:"type"`           // 此次活动的等级上限
	Status    CosplayStatus `bson:"status" json:"status"`       // 此次活动的状态

	OpenTime  int64 `bson:"openTime" json:"openTime"`   // 此次活动的开始时间
	CloseTime int64 `bson:"closeTime" json:"closeTime"` // 此次活动的发布时间

	AdminUsername string `bson:"adminUsername" json:"adminUsername"` // 发布人员用户名
	AdminNickname string `bson:"adminNickname" json:"adminNickname"` // 发布人员昵称

	Reward1 string `bson:"reward1" json:"reward1"` // 此次活动的冠军奖品 已废弃
	Reward2 string `bson:"reward2" json:"reward2"` // 此次活动的亚军奖品 已废弃
	Reward3 string `bson:"reward3" json:"reward3"` // 此次活动的季军奖品 已废弃

	Icon   string `bson:"icon" json:"icon"`     // 此次活动的奖品图片
	CosBg  string `bson:"cosBg" json:"cosBg"`   // 此次活动的背景图片
	ListBg string `bson:"listBg" json:"listBg"` // 此次活动的列表界面图片

	Params CosParams `bson:"cosParams" json:"cosParams"` // 当前活动的算分参数
}

type FCosplay struct {
	CosplayID     int64
	Title         string
	Keyword       string
	Type          int
	Status        CosplayStatus
	OpenTime      int64
	CloseTime     int64
	AdminUsername string
	AdminNickname string
	Icon          string
	CosBg         string
	ListBg        string
	Params        CosParams
}

type RawCosplay struct {
	CosplayID int64         `bson:"cosplayID" json:"cosplayID"` // 此次活动的编号
	Title     string        `bson:"title" json:"title"`         // 此次活动的标题
	Keyword   string        `bson:"keyword" json:"keyword"`     // 此次活动的关键字
	Type      int           `bson:"type" json:"type"`           // 此次活动的等级上限
	Status    CosplayStatus `bson:"status" json:"status"`       // 此次活动的状态

	OpenTime  string `bson:"openTime" json:"openTime"`   // 此次活动的开始时间
	CloseTime string `bson:"closeTime" json:"closeTime"` // 此次活动的发布时间

	AdminUsername string `bson:"adminUsername" json:"adminUsername"` // 发布人员用户名
	AdminNickname string `bson:"adminNickname" json:"adminNickname"` // 发布人员昵称

	Params   string `bson:"cosParams" json:"cosParams"` // 当前活动的算分参数
	Platform string
}

type UploadCosExtra struct {
	CosplayID     int64  `json:"cosplayID, omitempty"`
	Title         string `json:"title, omitempty"`
	CosTitle      string `json:"cosTitle, omitempty"`
	Desc          string
	ModelNo       int
	ServerAddress string
	DyeCnt        int

	DiamondBefore int
	DiamondChange int
	DiamondAfter  int

	Item101Before int
	Item101Change int
	Item101After  int

	PlayerSize int
}
