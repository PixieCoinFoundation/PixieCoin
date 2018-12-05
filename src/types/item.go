package types

// 玩家拥有的道具
//id:101 入场券 上传一次作品
//id:102 发言券 2次发言机会
//id:103 加速券 直接回满发言机会
//id:104 发言权 发言槽永久增加1
//id:105 加速器
//id:106 插播券
type Item struct {
	ItemID    int `bson:"itemID" json:"itemID"`
	Count     int `bson:"count" json:"count"`
	MaxCount  int `bson:"maxCount" json:"maxCount"`
	PriceType int `bson:"priceType" json:"priceType"`
	Price     int `bson:"price" json:"price"`
}
