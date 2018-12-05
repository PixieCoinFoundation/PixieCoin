package types

type CustomStatus int

const (
	//以下常量没有特殊情况不允许变动
	QUEUE               CustomStatus = 0  // 等待审核
	PROCESSING          CustomStatus = 1  // 正在审核
	FAIL                CustomStatus = 2  // 审核未通过
	PASSED              CustomStatus = 3  // 审核通过且上架
	PASSED_NO_INVENTORY CustomStatus = 4  //没有库存限制
	CUSTOM_CANCEL       CustomStatus = 5  // 玩家取消设计
	CUSTOM_READY        CustomStatus = 6  // 待上架状态
	CUSTOM_BAN          CustomStatus = -1 // 举报冻结
	CUSTOM_OUTOFPRINT   CustomStatus = -3 //绝版 无法再上下架
	CUSTOM_SALE_BAN     CustomStatus = -2 //禁止上架销售
)

type GetHotCustomResp struct {
	Designers []HotDesigner `json:"hot1"`
	Customs   []HotClothes  `json:"hot2"`
}

type HotDesigner struct {
	Name     string       `json:"name"`
	Head     string       `json:"head"`
	FanCount int          `json:"fansCnt"`
	HotList  []HotClothes `json:"clothes"`
}

type HotClothes struct {
	Icon string `json:"icon"`
	Cnt  int    `json:"cnt"`
}

type CustomFile struct {
	Front  string //背景前景
	Back   string //背景
	Icon   string
	Main   string
	Bottom string
	Collar string
	Shadow string
}

type RecentCustom struct {
	CID      int64
	Hearts   int
	IconPath string
}

//ADCustom 广告位信息
type ADCustom struct {
	Custom
	AuctionPos   int   //0无广告位，1上，2下，
	AuctionState int8  // 0普通状态， 竞拍中1，竞拍成功2，即在竞拍中又在热卖中，竞拍失败3
	EndTime      int64 //竞拍结束时间
}
type RewardCustom struct {
	Custom
	IsCanSetReward bool
}

//服装pool
type CustomPool struct {
	Custom
	Topic      int
	Poolstatus int
	IsStick    bool
	EntryTime  int64
}

// 服装类型
type Custom struct {
	CID             int64        `bson:"cid" json:"cid"`               // 服装序号
	ModelNo         int          `bson:"modelNo" json:"modelNo"`       // 模特编号
	Type            int          `bson:"type" json:"type"`             // 衣服类型
	CName           string       `bson:"cname" json:"cname"`           // 衣服的名字
	Desc            string       `bson:"desc" json:"desc"`             // 衣服的描述
	CloSet          CustomFile   `bson:"clothesSet" json:"clothesSet"` // 服装文件ID
	Status          CustomStatus `bson:"status" json:"status"`         // 当前的状态
	UploadTime      int64        `bson:"uploadTime" json:"uploadTime"` // 上传时间
	CheckTime       int64        `bson:"checkTime" json:"checkTime"`   // 审核时间
	MoneyType       int          `bson:"mType" json:"mType"`           // 价格类型，1金币，2钻石
	MoneyAmount     int          `bson:"mAmount" json:"mAmount"`       // 价格数量
	Hearts          int          `bson:"hearts" json:"hearts"`         // 心数
	Username        string       `bson:"username" json:"username"`     // 设计师用户名
	Nickname        string       `bson:"nickname" json:"nickname"`     // 设计师昵称
	BuyCount        int          `bson:"buyCount" json:"buyCount"`     // 售出数量
	Inventory       int          `bson:"inventory" json:"inventory"`   // 库存
	Extra           CustomExtra
	Tag1            int
	Tag2            int
	ClothesIDInGame string `bson:"clothesIDInGame" json:"clothesIDInGame"`

	AdminUsername string `bson:"adminUsername" json:"adminUsername"` // 审核人员用户名
	Info          string `bson:"info" json:"info"`                   // 拒绝时的信息
	NewSold       bool   `bson:"newSold" json:"newSold"`             // 新售出标识
	Platform      string
	Remark        string

	GoldProfit       int
	DiamondProfit    int
	Exp              int
	DesignCoinProfit int //设计师币收益
	PlayerConfig     int //玩家设置是否为非卖品

	StickTime int64 //上次置顶时间
}

type ElasticCustom struct {
	CID             int64        `bson:"cid" json:"cid"`               // 服装序号
	ModelNo         int          `bson:"modelNo" json:"modelNo"`       // 模特编号
	Type            int          `bson:"type" json:"type"`             // 衣服类型
	CName           string       `bson:"cname" json:"cname"`           // 衣服的名字
	Desc            string       `bson:"desc" json:"desc"`             // 衣服的描述
	CloSet          CustomFile   `bson:"clothesSet" json:"clothesSet"` // 服装文件ID
	Status          CustomStatus `bson:"status" json:"status"`         // 当前的状态
	UploadTime      int64        `bson:"uploadTime" json:"uploadTime"` // 上传时间
	CheckTime       int64        `bson:"checkTime" json:"checkTime"`   // 审核时间
	MoneyType       int          `bson:"mType" json:"mType"`           // 价格类型，1金币，2钻石
	MoneyAmount     int          `bson:"mAmount" json:"mAmount"`       // 价格数量
	Hearts          int          `bson:"hearts" json:"hearts"`         // 心数
	Username        string       `bson:"username" json:"username"`     // 设计师用户名
	Nickname        string       `bson:"nickname" json:"nickname"`     // 设计师昵称
	BuyCount        int          `bson:"buyCount" json:"buyCount"`     // 售出数量
	Inventory       int          `bson:"inventory" json:"inventory"`   // 库存
	Extra           CustomExtra
	Tag1            int
	Tag2            int
	ClothesIDInGame string `bson:"clothesIDInGame" json:"clothesIDInGame"`

	UpTime int64 `json:"upTime"`
}

type GVCustom struct {
	CID             int64  `bson:"cid" json:"cid"`         // 服装序号
	ModelNo         int    `bson:"modelNo" json:"modelNo"` // 模特编号
	ModelName       string `json:"modelName"`
	Type            int    `bson:"type" json:"type"` // 衣服类型
	TypeName        string `json:"typeName"`
	CName           string `bson:"cname" json:"cname"`                 // 衣服的名字
	UploadTime      string `bson:"uploadTime" json:"uploadTime"`       // 上传时间
	AdminUsername   string `bson:"adminUsername" json:"adminUsername"` // 审核人员用户名
	MoneyType       string `bson:"mType" json:"mType"`                 // 价格类型，1金币，2钻石
	MoneyAmount     int    `bson:"mAmount" json:"mAmount"`             // 价格数量
	Hearts          int    `bson:"hearts" json:"hearts"`               // 心数
	Username        string `bson:"username" json:"username"`           // 设计师用户名
	Nickname        string `bson:"nickname" json:"nickname"`           // 设计师昵称
	ClothesIDInGame string `bson:"clothesIDInGame" json:"clothesIDInGame"`
	BuyCount        int    `bson:"buyCount" json:"buyCount"`   // 售出数量
	Inventory       int    `bson:"inventory" json:"inventory"` // 库存
	Status          string `json:"status"`                     //衣服状态
	Info            string `bson:"info" json:"info"`           // 拒绝时的信息
	Params          string `bson:"params" json:"params"`
	PlayerConfig    int    `json:"playerConfig"` //玩家设置的衣服状态 1为非卖品
	ZValue          int    `json:"zValue"`
	TagDesc         string `json:"tagDesc"`
}

type ExtraInfo struct {
	Name string
	X    float64
	Y    float64
}

type CustomExtra struct {
	ZValue int

	IconG   string
	MainG   string
	BottomG string
	CollarG string
	ShadowG string

	IconX   float64
	MainX   float64
	BottomX float64
	CollarX float64
	ShadowX float64

	IconY   float64
	MainY   float64
	BottomY float64
	CollarY float64
	ShadowY float64

	Warm    float64
	Formal  float64
	Tight   float64
	Bright  float64
	Dark    float64
	Cute    float64
	Man     float64
	Tough   float64
	Noble   float64
	Strange float64
	Sexy    float64
	Sport   float64
}

type CopyReport struct {
	ID       int64
	Username string
	CID      int64
	Reason   string
	Evidence string
	Contact  string
	Status   int
	Type     int
}

type DesignMsg struct {
	//老版会有ID，改后均为0
	ID int64

	DesignerHead     string
	DesignerUsername string
	DesignerNickname string

	Type int //1:发布 2:上架 3:下架 4:购物车上架

	//type为4时读
	ClothesID   string
	ClothesName string
	IconFileID  string
	MainFileID  string

	//type为1-3时读
	Details []DesignMsgDetail

	Time int64
}

type OldDesignMsg struct {
	//老版会有ID，改后均为0
	ID int64

	DesignerHead     string
	DesignerUsername string
	DesignerNickname string

	Type int //1:发布 2:上架 3:下架

	ClothesID   string
	ClothesName string
	IconFileID  string

	Time int64
}

type DesignMsgDetail struct {
	ClothesID   string
	ClothesName string
	IconFileID  string
	MainFileID  string
}

type DesignMsgInfo struct {
	Username string
	Nickname string
	Head     string
	Type     int
	Time     int64
}

type SDesignMsg struct {
	ID               int64
	DesignerUsername string
	Type             int
	Time             int64
}

type SimpleCustom struct {
	ClothesIDInGame string
	Name            string
}
