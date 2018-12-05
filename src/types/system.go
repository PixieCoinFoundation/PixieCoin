package types

type KorUnregisterInfo struct {
	ID            int64
	Username      string
	ThirdUsername string
	DelTime       int64
}

type KorLogClean struct {
	ID   int64
	Time int64
}

type BlockInfo struct {
	Username string
	Reason   string
	AReason  string
}

type ClientLogType struct {
	ID               int
	Name             string
	LastRefreshTime  int64
	LastRefreshTimes string
}

type MonthCardReward struct {
	Type       int    //1:套装ID
	RewardInfo string //Type:1:套装ID

	OuterImg string
	InnerImg string
}

type SimpleUserInfo struct {
	Username string
	Nickname string
}

type UserSaleInfo struct {
	SimpleUserInfo
	Sale int
}

type PlayerTheaterInfo struct {
	SimpleUserInfo
	Score int
}

type SinglePushJob struct {
	To      string `json:"t"`
	Title   string `json:"tt"`
	Content string `json:"c"`
	Type    int    `json:"ty"`
}

type NewCustomPushObject struct {
	DesignerUsername string
	DesignerNickname string
	CName            string
}

type GFBackground struct {
	ID        string
	PriceType int
	Price     int
	From      int
	Name      string
}

type ResetResp struct {
	Info                string
	Platform            string
	HTTPType            string
	PicAddress          string
	MServer             string
	TotalCount          int
	DesignerBanHistorys []string
	List                []DesignerZoneBanner
	BanDesigners        map[string]string
}

type DownCustomData struct {
	CloID     string
	DownTime  int64
	Processed bool
}

type LotteryPubRespData struct {
	ID    string `json:"id"`
	SName string `json:"sname"`
	RName string `json:"rname"`
	Info  string `json:"info"`
	Star  string `json:"star"`
}

type ServerHB struct {
	ID            int64  `json:"id"`
	Address       string `json:"address"`
	Type          string `json:"type"`
	HeartBeatTime string `json:"heartbeat_time"`
}

type RTData struct {
	Time int64
	Data int
	Desc string
}

type ChannelLimit struct {
	Limit     bool
	LimitInfo map[string]ChannelNum
	Platform  string
}

type ChannelNum struct {
	LimitNum  int
	ActualNum int
}

type DBConfig struct {
	// Start   int64
	// End     int64
	IPPort string
	// IPPort2 string
	Weight int
}

type CResp struct {
	C1 string
	C2 string
}

type CommonReport struct {
	Name       string
	UserCnt    int
	AllCnt     int
	UserPerCnt float64
}

type PKReport struct {
	UserCnt    int
	AllCnt     int
	UserPerCnt float64
	WinCnt     int
	LoseCnt    int
	EscapeCnt  int
	WinRate    float64
	WinUser    int
	LoseUser   int
	EscapeUser int
}

type MaintainConfig struct {
	Maintain     bool
	MaintainInfo string
}

type GameParam struct {
	PKDiscount int
	// Theaters        []TheaterInfo
	OrderStatus     bool
	MonthCardStatus bool
	GiftPackStatus  bool
	DelUserDuration int
	GuestLogin      bool
}

type Push struct {
	ID        int64
	Title     string
	Content   string
	PushTime  int64
	PushTimes string
	To        string
	Status    int
}

type MaintainJob struct {
	ID         int64
	Content    string
	StartTime  int64
	EndTime    int64
	StartTimes string
	EndTimes   string
	ShowTime   int64
	ShowTimes  string
	URL        string
}
