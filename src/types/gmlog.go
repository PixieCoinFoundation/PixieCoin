package types

import (
	. "pixie_contract/api_specification"
)

type CustomReclaim struct {
	BeforeGold    int
	BeforeDiamond int
	AfterGold     int
	AfterDiamond  int
	Gold          int
	Diamond       int

	BeforeDesignCoin int
	DesignCoin       int
	AfterDesignCoin  int

	Exp   int
	CID   int64
	CloID string
}

type TailorExtra struct {
	GoldBefore int
	GoldChange int
	GoldAfter  int

	ClothBefore int
	ClothChange int
	ClothAfter  int

	ButtonBefore int
	ButtonChange int
	ButtonAfter  int

	DropClothes string
	GotClothes  string
	ClothesStar string
}

//DecoClothes 分解衣服
type DecoClothes struct {
	ClothesList []string
	ClothesStar []string

	GotCloth     int
	CurrentCloth int
}
type VIPChangeExtra struct {
	BuyDiamondBefore int
	BuyDiamondChange int
	BuyDiamiondAfter int
	VIPBefore        int
	VIPAfter         int
}

//UploadClothesList 衣服上架信息
type UploadClothesList struct {
	ShelvesTime string //衣服上架时间
	ClothesID   string //衣服ID
	ModType     string //模特类型
	ClothesType string //衣服类型
	ShelvesNum  int    //上架数量
	SellPrice   int    //销售价格
	SellType    int    //销售类型
}

//Flower 花
type Flower struct {
	PartyId            int64
	Username           string
	SendFlowerUsername string //
	PartyType          int
	PartyCloseTime     int64
}

type PartyItemCmtExtra struct {
	PartyID       int64
	Username      string
	Content       string
	ServerAddress string
}

type MarkPartyItemExtra struct {
	PartyID       int64
	Username      string
	Type          int //0:support 1:nonsupport 2:flower
	ServerAddress string
}

type BeMarkPartyItemExtra struct {
	PartyID       int64
	MarkUsername  string
	MarkNickname  string
	Type          int
	ServerAddress string
}

type HostPartyJoinExtra struct {
	JoinUsername  string
	JoinNickname  string
	ServerAddress string
}

type PKDestroyCloExtra struct {
	ClothesID     string
	ServerAddress string
}

type CosItemScoreExtra struct {
	DestUsername    string
	DestCosplayID   int64
	DestCosplayName string
	DestCosItemID   int64
	DestCosItemName string
	Score           int
}

type SystemDelCloExtra struct {
	ClothesID        string
	Cnt              int
	ReturnGold       int64
	ReturnDiamond    int64
	ReturnDesignCoin int64
}

type LoginDataSnapshotExtra struct {
	Nickname               string
	Level                  int
	VIP                    int
	Money                  int64
	Diamond                int
	Hearts                 int
	ThirdChannel           string
	UTCOffset              int
	TheatreStatus          string
	EndingStatus           string
	DesignerPoints         int
	CosPoints              int
	CurrentMissionLevel    int
	CurrentMissionNo       int
	OfficialBoyClothesCnt  int
	OfficialGirlClothesCnt int
	DesignerBoyClothesCnt  int
	DesignerGirlClothesCnt int
	PKGWeekPoints          int
	PKWeekPoints           int

	Cloth  int
	Button int
}

type NewPlayerExtra struct {
	ThirdChannel  string
	ServerAddress string
}

type AbnormalExtra struct {
	ServerAddress string
	Info          string
}

type OrderExtra struct {
	VIP             int
	TotalBuyDiamond int
	OS              string
	ThirdUsername   string
	Nickname        string
	OrderID         int64
	OrderUID        string
	DiamondID       int
	RMB             float64 //money cost:usd,cny,krw
	FirstIAP        int
	ServerAddress   string
	ThirdChannel    string
	OrderName       string
	CurrencyType    string

	KORThirdCode  int
	KOROSCode     int
	KORStoreCode  int
	KORGameCode   int
	KORPayCountry string
	KORLogType    int
}

type QinmiExtra struct {
	TargetUsername string
	Delta          int
	Reason         string
}

type DesignerIncomeExtra struct {
	ClothesID    string
	DiamondPrice int
	GoldPrice    int

	DiamondProfit int
	GoldProfit    int
	Rate          float64
}

type Row struct {
	C1       string
	C2       string
	C3       string
	Username string
	Time     string
	Extra    string
}
type BG struct {
	BID         string
	MoneyAgo    int
	MoneyChange int
	DiamAgo     int
	DiamChange  int
	DeviceId    string
}
type NItemLogExtra struct {
	Nickname      string
	Level         int
	VIP           int
	Before        int
	ItemName      string
	ItemDelta     int
	ItemLeft      int
	DeviceId      string
	ServerAddress string
}

//MatchPKExtra 匹配的用户信息
type MatchPKExtra struct {
	UID      int64
	KakaoID  string
	Nickname string
}

type BuyBGlog struct {
	BID       string
	MoneyAgo  int
	MoneyCost int
	MoneyLeft int
	DiamAgo   int
	DiamCost  int
	DiamLeft  int
}

//DesignerSellExtra 设计师衣服信息 和卖衣服细节
type DesignerSellExtra struct {
	ClothesID      string
	ClothesPrice   int
	MoneyType      int
	SellCnt        int
	DesignerProfit int
	Customer       string
	DeviceId       string
	UID            int64  //客户的UID
	Kakaoid        string //客户的kakaoid
	ClothesName    string //衣服名称
	ClothesPart    string //服装部位
	ClothesTag1    int    //服装tag1
	ClothesTag2    int    //服装tag2
	Price          int    //价格
	MoneyAgo       int    //购买前持有的金额
	MoneyAfter     int    //购买后持有的金额
	ServerAddress  string

	BuyTransactionHash string
}
type PartyHolder struct {
	Username string
	Nickname string
	PartyID  int64
}
type AddPartyItemExtra struct {
	Username       string
	Nickname       string
	PartyID        int64
	PartyName      string
	PartyHost      string
	PartyStartTime string
	PartyCloseTime string
	Img            string
	DyeCnt         int
	ServerAddress  string
	JoinType       int //0:normal 1:join as partner
	Rank           int
	RewardGold     int64
	RewardDiamond  int64

	PartnerType   int //0:no partner 1:all 2:chosed friend
	AllowPartners []UserInfo
}

type ExpChangeExtra struct {
	GotExp        int
	BeforeLevel   int
	AfterLevel    int
	BeforeExp     int
	AfterExp      int
	Reason        string
	ServerAddress string
}

type OrderUpdateExtra struct {
	OrderUID string
	PayTime  int64
	Status   int
	Channel  string
	Info     string
}

type OnlineExtra struct {
	OnlineNum     int
	ServerAddress string
}

type CurrencyChangeExtra struct {
	Nickname      string
	Level         int
	VIP           int
	Amount        float64
	Before        float64
	Left          float64
	Extra         string
	DeviceId      string `json:"deviceId, omitempty"`
	ServerAddress string
	Info          string
	NormalStore   bool
	DesignerStroe bool
}

type CloPieceChangeExtra struct {
	Key           string
	Amount        int
	Left          int
	DeviceId      string `json:"deviceId, omitempty"`
	ServerAddress string
}

type CurrencyExtra struct {
	Amount   int
	DeviceId string
}

type CurrencySnapshotExtra struct {
	Before        int
	After         int
	DeviceId      string `json:"deviceId, omitempty"`
	ServerAddress string
}

type LevelExtra struct {
	LevelID       string
	DeviceId      string `json:"deviceId, omitempty"`
	ServerAddress string
}

type UnlockTagExtra struct {
	UnlockIndex int
	DeviceId    string `json:"deviceId, omitempty"`
}

type DeltaClothesExtra struct {
	ClothesID     string
	Amount        int
	Left          int
	Role          string
	StoreType     string
	Part          string
	DeviceId      string `json:"deviceId, omitempty"`
	ServerAddress string
	TheaterID     string
	SellType      string
}

type DeltaFriendExtra struct {
	Username      string
	DeviceId      string `json:"deviceId, omitempty"`
	ServerAddress string
}

type LotteryExtra struct {
	MagicBefore   int
	MagicChange   int
	MagicLeft     int
	Count         int
	MoneyBefore   int
	MoneyType     string
	MoneyCount    int
	MoneyLeft     int
	GotClothes    []*ClothesDetailInfo
	DeviceId      string `json:"deviceId, omitempty"`
	ServerAddress string
	PoolID        int
}

//DesignerInfoExtra 设计师信息变更
type DesignerInfoExtra struct {
	CID             int64 `json:"cid"`
	Reason          string
	DesignerUID     int64
	DesignerContent string
	DeviceId        string `json:"deviceId, omitempty"`
	Status          string
}

//FollowLog 关注日志
type FollowLog struct {
	PlayerUID        int64
	PlayerNickname   string
	Type             int //1:follow 2:unfollow
	DesignerUsername string
	DesignerNickname string
}

//HopeExtra 愿望模型
type HopeExtra struct {
	HopeTime           string
	HopePlayerUsername string //许愿用户ID
	HopePlayerNickname string //许愿用户昵称
	HopeFinishTime     string
	HopeID             int64
	ClothesID          string
	Part               int //1-4
	PieceCnt           int

	HelpPlayerUID      int64  //帮助完成愿望用户的ID
	HelpPlayerNickname string //昵称
	ReqClothesID       string
	ReqPart            int //1-4
	ReqPieceCnt        int
	ServerAddress      string
}

type HelpHopeExtra struct {
	HopeID int64

	InClothesID string
	InPart      int //1-4
	InPieceCnt  int

	OutClothesID  string
	OutPart       int //1-4
	OutPieceCnt   int
	ServerAddress string
}

type LotteryExtra1 struct {
	MagicLeft  int
	Count      int
	MoneyType  string
	MoneyCount int
	GotClothes []ClothesDetailInfo
	PoolID     int
	DeviceId   string `json:"deviceId, omitempty"`
}

type MailExtra struct {
	ID            int64
	From          string
	Title         string
	Content       string
	Gold          int64
	ReadTime      string `json:"readtime"`
	Diamond       int64
	Clothes       []ClothesInfo
	ServerAddress string
}

type SwapCloPieceExtra struct {
	FromUsername    string
	ToUsername      string
	FromCloID       string
	FromCloPart     int
	FromCloPieceCnt int
	ToCloID         string
	ToCloPart       int
	ToCloPieceCnt   int
	ServerAddress   string
}
