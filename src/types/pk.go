package types

type GFPK struct {
	ID               int64
	Username         string `bson:"username", json:"-"`
	LevelID          string `bson:"lid", json:"lid"`               // 当前pk操作对应的关卡id
	Nickname         string `bson:"nickname", json:"nickname"`     // 玩家上传时的昵称
	Head             string `bson:"head", json:"head"`             // 玩家上传时的头像
	PKPoint          int    `bson:"pkPoint", json:"pkPoint"`       // 玩家积分
	PKLevel          int    `bson:"pkLevel", json:"pkLevel"`       // pk等级
	GirlClothesCount int    `bson:"gcCount", json:"gcCount"`       // 女装数
	BoyClothesCount  int    `bson:"bcCount", json:"bcCount"`       // 男装数
	Operations       string `bson:"ops", json:"ops"`               // 操作记录串
	UploadTime       int64  `bson:"uploadTime", json:"uploadTime"` // 更新时间
}

type PKOP struct {
	Type int `json:"type"`
	Time int `json:"time"`
}

type PKRank struct {
	Username string
	Nickname string `json:"name"`
	Point    int    `json:"score"`
}
