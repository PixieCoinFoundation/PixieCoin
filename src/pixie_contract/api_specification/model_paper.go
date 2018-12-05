package api_specification

//paper基础信息
type BasePaper struct {
	PaperID        int64
	AuthorUsername string
	AuthorNickname string
	AuthorHead     string
	AuthorSex      int
	OwnerUsername  string
	OwnerNickname  string
	OwnerHead      string
	OwnerSex       int

	Circulation    int //已流通数量
	CirculationMax int //最大流通数量

	ClothesType string //衣服模特
	PartType    string //衣服类别
	Cname       string //图纸名称
	Desc        string //描述
	File        string //图纸源文件 PaperFile json
	Extra       string //图纸文件处理后 PaperExtra json
	Status      int    //图纸状态
	Star        int    //图纸星级

	//普通标签
	Tag1 string
	Tag2 string

	//风格 内容是PaperStyle的json字符串
	Style string

	//暗标签
	STag string
}

type OfficialPaper struct {
	BasePaper

	AdminName   string
	UploadTime  int64
	SaleTime    int64
	UnlockLevel string
	Price       int64
	PriceType   int
}

type PaperStyle struct {
	Warm     int `json:"Warm,omitempty"`
	Formal   int `json:"Formal,omitempty"`
	Tight    int `json:"Tight,omitempty"`
	Bright   int `json:"Bright,omitempty"`
	Cute     int `json:"Cute,omitempty"`
	Man      int `json:"Man,omitempty"`
	Tough    int `json:"Tough,omitempty"`
	Sexy     int `json:"Sexy,omitempty"`
	Gorgeous int `json:"Gorgeous,omitempty"`
	Noble    int `json:"Noble,omitempty"`
	Strange  int `json:"Strange,omitempty"`
	Sport    int `json:"Sport,omitempty"`
}

type PaperExtra struct {
	ZValue int `json:"ZValue,omitempty"`

	IconG   string `json:"IconG,omitempty"`
	MainG   string `json:"MainG,omitempty"`
	BottomG string `json:"BottomG,omitempty"`
	CollarG string `json:"CollarG,omitempty"`
	ShadowG string `json:"ShadowG,omitempty"`

	IconX   float64 `json:"IconX,omitempty"`
	MainX   float64 `json:"MainX,omitempty"`
	BottomX float64 `json:"BottomX,omitempty"`
	CollarX float64 `json:"CollarX,omitempty"`
	ShadowX float64 `json:"ShadowX,omitempty"`

	IconY   float64 `json:"IconY,omitempty"`
	MainY   float64 `json:"MainY,omitempty"`
	BottomY float64 `json:"BottomY,omitempty"`
	CollarY float64 `json:"CollarY,omitempty"`
	ShadowY float64 `json:"ShadowY,omitempty"`

	IconW   float64 `json:"IconW,omitempty"`
	MainW   float64 `json:"MainW,omitempty"`
	BottomW float64 `json:"BottomW,omitempty"`
	CollarW float64 `json:"CollarW,omitempty"`
	ShadowW float64 `json:"ShadowW,omitempty"`

	IconH   float64 `json:"IconH,omitempty"`
	MainH   float64 `json:"MainH,omitempty"`
	BottomH float64 `json:"BottomH,omitempty"`
	CollarH float64 `json:"CollarH,omitempty"`
	ShadowH float64 `json:"ShadowH,omitempty"`
}

type PaperFile struct {
	//背景专用
	BeijingFileType int    `json:"BeijingFileType,omitempty"` //文件类型 0图片 1动画
	Front           string `json:"Front,omitempty"`           //背景前景
	Back            string `json:"Back,omitempty"`            //背景

	//普通服装属性
	Icon   string `json:"Icon,omitempty"`
	Main   string `json:"Main,omitempty"`
	Bottom string `json:"Bottom,omitempty"`
	Collar string `json:"Collar,omitempty"`
	Shadow string `json:"Shadow,omitempty"`
}

type SimplePaper struct {
	PaperID   int64
	PriceType int
	Price     float64
	BuyCount  int
}

type SimpleClothes struct {
	ClothesID string
	PriceType int
	Price     float64
	BuyCount  int
}

type OwnPaper struct {
	BasePaper

	Sequence int //图纸流通序号

	SaleCount        int     //销售数量
	ClothPrice       float64 //衣服价格
	ClothPriceType   int     //衣服价格类型
	Inventory        int     //库存
	IsInProduction   bool    // 是否生产中
	ProductStartTime int64   //开始生产时间
	ProductLeftCount int     //剩余生产数量
	ProductDoneCount int     //已经生产的
	OccupyLogID      int64   //销售记录ID

	//附加属性
	AuctionFailUnread bool //竞拍失败-是否未读

	//拍卖价格类型
	PriceType    int   //status为PAPER_STATUS_ONSALE时才读这个值
	LastDuration int64 //上一次时常 hour *3600
	LastMaxPrice int64 //上一次最高价
	LastMinPrice int64 //上一次最低价
}

type DesignPaper struct {
	BasePaper

	UploadTime    int64 //上传时间
	PassTime      int64 //通过时间
	VerifiedCount int   //审核人数

	//附加属性
	AuctionFailUnread bool //发布失败-是否未读

	PriceType    int   //status为PAPER_STATUS_ONSALE时才读这个值
	LastDuration int64 //上一次时常 hour *3600
	LastMaxPrice int64 //上一次最高价
	LastMinPrice int64 //上一次最低价

	CopyMark bool   //是否被举报
	Reason   string //拒绝理由
}

type RecommendClothes struct {
	BasePaper
	ClothesID string //服装ID
	Price     float64
	PriceType int
}

type SalePaper struct {
	BasePaper

	TradeID  int64
	Sequence int

	PriceType       int //图纸价格类型
	StartTime       int64
	Duration        int64
	MaxPrice        int64
	MinPrice        int64
	AuctionTime     int64
	AuctionUsername string
	AuctionNickname string
	AuctionPrice    int64
	InSale          bool //是否在售卖中
}

type PaperVerify struct {
	ID           int64
	PaperID      int64
	Cname        string //图纸名称
	Extra        string //图纸内容
	Username     string
	Nickname     string
	CheckDate    string //"20060102"
	Tag1         string
	Tag2         string
	Style1       string
	Style2       string
	Score        int
	MatchDegree  float64 //匹配度
	RewardPxc    float64 //奖励数量
	SubmitTime   int64   //审核时间
	CheckoutTime int64   //更新结果时间
	Status       int     //0.待结算,1.已经结算,2.已删除
	Result       string  //计算结果star,tag1,tag2
}

type PaperVerifyResult struct {
	Score  float64
	Tag1   string
	Tag2   string
	Style1 string
	Style2 string
}

type PaperVerifyReward struct {
	PaperID     int64
	Cname       string
	Content     string  //奖励内容
	MatchDegree float64 //匹配度
}

type PaperNotify struct {
	ID           int64
	Username     string
	Content      string  //奖励内容 RewardContent
	DayPxcProfit float64 //每日pxc奖励
	Cname        string
	Time         int64 //发送通知
	Status       int   //1是领取奖励，0是未领取
	Type         int   //1.是举报奖励，2.是评审奖励 3.无奖励
}

type RewardContent struct {
	VerifyResultNum int     //审核结果次数
	CurrencyType    int     //货币类型
	CurrencyVal     float64 //货币数量
	DayVerifyCount  int     //每日的评审人数
}

type PaperCopy struct {
	CopyID     int64
	PaperID    int64
	Username   string
	Reason     string
	Pic        string
	Contact    string
	SupportNum int
	RejectNum  int
	Time       int64
	CheckTime  int64
	Status     bool //是否处理
}

type PaperTradeHistory struct {
	PaperID        int64
	Username       string
	Nickname       string
	TimeBuy        int64 //交易时间
	TradePrice     int64 //交易价格
	TradePriceType int
	SaleCount      int //销售量
}

type PaperInfo struct {
	BasePaper

	SaleCount        int     //销售数量
	ClothPrice       float64 //衣服价格
	ClothPriceType   int     //衣服价格类型
	Inventory        int     //库存
	IsInProduction   bool    // 是否生产中
	ProductStartTime int64   //开始生产时间
	ProductLeftCount int     //剩余生产数量
	ProductDoneCount int     //已经生产的
	OccupyLogID      int64   //销售记录ID

	//附加属性
	AuctionFailUnread bool //竞拍失败-是否未读

	//拍卖价格类型
	PriceType int //status为PAPER_STATUS_ONSALE时才读这个值

	UploadTime  int64 //上传时间
	PassTime    int64 //通过时间
	VerifyCount int64 //审核人数

	StartTime       int64
	Duration        int64
	MaxPrice        int64
	MinPrice        int64
	AuctionTime     int64
	AuctionUsername string
	AuctionNickname string
	AuctionPrice    int64
}
