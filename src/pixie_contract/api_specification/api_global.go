package api_specification

//心跳接口
const HTTP_PLAYER_HEART_BEAT string = "api/heartBeat"

type HeartBeatResp struct {
	Time      int64
	UTCOffset int

	Tili      int   //剧情体力
	MaxTili   int   //最大剧情体力
	TiliRTime int64 //上次恢复体力的时间

	//我的街目前状态
	StreetDetail *StreetOperateResult
}

//游戏登录接口
const HTTP_PLAYER_LOGIN string = "api/loginGFServer"

type PixieLoginReq struct {
	ThirdUsername string
	Username      string
	AccessToken   string
	DeviceToken   string
	UID           int64
	// GFToken             string
	Platform            string
	DeviceId            string
	Channel             string
	ThirdChannel        string
	GuestLogin          bool   //游客登录
	DeviceModelType     string //设备名称
	KorRevertUnregister bool
}

//游戏登录接口

//获取玩家数据接口
const HTTP_PLAYER_GET_ALL_USER_DATA string = "api/getAllUserData"

type GetAllUserDataResp struct {
	UserAttributes     Status
	IsContractDesigner bool
	CurrentTime        int64
	UTCOffset          int

	ShouldUploadLog bool

	//玩家拥有的地块信息
	Lands []Land

	ShopCart map[string]*ShopCartDetail

	RTPartyTime []string

	PrizePartyClose bool

	MailCount   int
	NoticeCount int
	BoardNew    bool

	SpecialAccount  bool
	OrderOpen       bool
	GiftPackOpen    bool
	MonthCardCanBuy bool

	GuanZhuDesigners map[string]int

	LevelMap map[string]int

	ChannelShareURL string

	GuildBoardNewMsg bool

	SkinMap map[string]bool
}

// type TheaterInfo struct {
// 	TheaterStartTime int64
// 	TheaterCloseTime int64
// 	TheaterID        int
// 	TheaterName      string
// 	UnlockModel      int
// 	SuitID           string
// 	RankKey          string
// }

type ShopCartDetail struct {
	Count   int
	AddTime int64
}
