package types

// "gopkg.in/mgo.v2/bson"

//ADGFDesigner 广告设计师
// type ADGFDesigner struct {

// 	GFDesigner
// }

// 设计师相关信息
type GFDesigner struct {
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Diamond      int    `json:"diamond"`      // 当前未提出的钻石数量
	Gold         int    `json:"gold"`         // 当前未提出的金币数
	TotalDiamond int    `json:"totalDiamond"` // 总共售出的钻石数量
	TotalGold    int    `json:"totalGold"`    // 总共售出的金币数

	DesignCoin      int `json:"designCoin"`      //当前未提出的设计币
	TotalDesignCoin int `json:"totalDesignCoin"` //总共获得的设计币

	TotalSale int `json:"totalSale"` // 总共售出的服装件数

	Points      int `json:"points"`      // 设计师的分数
	PassedCount int `json:"passedCount"` // 总共审核通过的衣服
	TotalCount  int `json:"totalCount"`  // 总共审核的衣服（包括审核通过和审核没通过的）

	FensiCount int    `json:"fensiCount"` //粉丝数
	Desc       string `json:"desc"`       //个人描述
	Head       string `json:"head"`       //头像

	ETHAccount string `json:"ethAccount"` //eth账号地址

	Ad_DayReward              string
	Ad_MonthOneReward         string
	Ad_MonthTenReward         string
	Ad_MonthThirtyReward      string
	Ad_AutoConfig             bool
	DesignCoinPriceTokenCount int
}

//设计师池
type GFDesignerPool struct {
	GFDesigner
	ClothesList   []Custom
	Poolstatus    int
	IsNewDesigner bool
	IsStick       bool
	EntryTime     int64
}
type SimpleDesigner struct {
	Username string
	Nickname string
	Head     string
}

// 高产榜
type GFDesignerRankTopAmountItem struct {
	Username    string `bson:"username" json:"username"`
	Nickname    string `bson:"nickname" json:"nickname"`
	PassedCount int    `bson:"passedCount" json:"passedCount"` // 总共审核通过的衣服
}

// 最热榜
type GFDesignerRankTopSaleItem struct {
	Username  string `bson:"username" json:"username"`
	Nickname  string `bson:"nickname" json:"nickname"`
	TotalSale int    `bson:"totalSale" json:"totalSale"` // 总共审核通过的衣服
}

//妖精币榜
type GFDesignerRankTopDesignCoinItem struct {
	Username        string `json:"username"`
	Nickname        string `json:"nickname"`
	TotalDesignCoin int    `json:"totalDesignCoin"` // 总共审核通过的衣服
}

// 金币榜
type GFDesignerRankTopGoldItem struct {
	Username  string `bson:"username" json:"username"`
	Nickname  string `bson:"nickname" json:"nickname"`
	TotalGold int    `bson:"totalGold" json:"totalGold"` // 总共审核通过的衣服
}

// 钻石榜
type GFDesignerRankTopDiamondItem struct {
	Username     string `bson:"username" json:"username"`
	Nickname     string `bson:"nickname" json:"nickname"`
	TotalDiamond int    `bson:"totalDiamond" json:"totalDiamond"` // 总共审核通过的衣服
}
