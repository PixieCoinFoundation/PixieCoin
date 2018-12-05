package types

// 游戏中的微博
type Tweet struct {
	ID        int  `bson:"ID", json:"ID"`               // tweet id
	Retweeted bool `bson:"Retweeted", json:"Retweeted"` // 是否已经转发过
	Show      bool `bson:"Show", json:"Show"`           // 是否需要显示
}
