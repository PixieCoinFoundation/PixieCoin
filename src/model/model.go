package model

type BufferLog struct {
	C1       string
	C2       string
	C3       string
	Username string
	Extra    string
	T        string
	T64      int64
}

//TimeDoor 时空之门
type TimeDoor struct {
	//开始时间
	Kor_shikongzhimen_start_time int64
	//结束时间
	Kor_shikongzhimen_end_time int64
	//衣服概率
	// Kor_shikongzhimen_clothes_chance string `json:"kor_shikongzhimen_clothes_chance"`
	Chance13 int
	Chance4  int
	//抽奖衣服
	LotteryClothesShikongLevel13 []string // 抽奖服装level
	LotteryClothesShikongLevel4  []string // 抽奖服装level
	LotteryClothesShikongLevel5  []string // 抽奖服装leve

	BannerLocation string
}

// 在线玩家的信息保存
type GFClothes struct {
	CId   string `json:"id"` // 衣服id
	Index int    `json:"index"`
	Name  string `json:"name"`
	// Desc      string `json:"string"`
	Type      int `json:"type"`       // 衣服类型，上衣，外套等
	SmallType int `json:"small_type"` // 衣服小类型
	Sex       int `json:"sex"`        // 衣服性别，1男，2女
	PriceType int `json:"price_type"` // 1金币，2钻石，其它不卖
	Price     int `json:"price"`      // 价格
	SellPrice int `json:"sell_price"`
	Star      int `json:"star"` // 星级
	// Level     int    `json:"level"` // 打到哪一层解锁

	// DropMission int `json:"mission_drop"` // 关卡掉落与否
	DropCos     int `json:"cos_drop"`     // Cos掉落与否
	DropPK      int `json:"pk_drop"`      // PK掉落与否
	DropVIP     int `json:"vip_drop"`     // VIP掉落与否
	DropLottery int `json:"lottery_drop"` // 抽奖掉落

	Warm    float64 `bson:"warm" json:"warm"`
	Formal  float64 `bson:"formal" json:"formal"`
	Tight   float64 `bson:"tight" json:"tight"`
	Bright  float64 `bson:"bright" json:"bright"`
	Dark    float64 `bson:"dark" json:"dark"`
	Cute    float64 `bson:"cute" json:"cute"`
	Man     float64 `bson:"man" json:"man"`
	Tough   float64 `bson:"tough" json:"tough"`
	Noble   float64 `bson:"noble" json:"noble"`
	Strange float64 `bson:"strange" json:"strange"`
	Sexy    float64 `bson:"sexy" json:"sexy"`
	Sport   float64 `bson:"sport" json:"sport"`

	Char int `bson:"char" json:"char"`

	// 设计师
	// Username string `bson:"username" json:"username"`
	// Nickname string `bson:"nickname" json:"nickname"`
}
