package types

type Party struct {
	ID int64

	Type       int //0:normal 1:prize 2:cosplay
	SingleType int //0:single 1:multi
	StartTime  int64
	CloseTime  int64

	Username string
	Nickname string

	Subject string
	Desc    string

	PlayerCnt int
	MoneyType int //0:gold 1:diamond
	StartPool int
	Ticket    int

	//bg & banner pair type
	//大于0是官方的默认背景和banner
	//-1:管理员指定的图片
	BgBannerType int

	//BgBannerType为-1时这两个值才有意义
	BannerFile     string
	BackgroundFile string

	JoinerCnt int
}

type GMParty struct {
	Party
	StartTimes string
	CloseTimes string
}

type PartyItemPair struct {
	Item1 PartyItem
	Item2 PartyItem
}

type PartyItem struct {
	ID            int64
	PartyID       int64
	Username      string
	UID           int64
	Nickname      string
	Head          string
	Img           string
	PartyUsername string
	PartyNickname string
	Subject       string
	Desc          string

	//bg & banner pair type
	//大于0是官方的默认背景和banner
	//-1:管理员指定的图片
	BgBannerType int

	//BgBannerType为-1时这两个值才有意义
	BannerFile     string
	BackgroundFile string

	Clothes         []string
	ModelNo         int
	PartnerType     int //0:no 1:all 2:friend
	AllowPartners   []string
	Partner         string
	PartnerQinmi    int
	PartnerNickname string
	PartnerHead     string
	DyeMap          map[string][7][4]float64

	//partner
	PartnerClothes []string
	PartnerModelNo int
	PartnerDyeMap  map[string][7][4]float64

	Popularity int
	LikeCnt    int
	UnlikeCnt  int
	FlowerCnt  int
	UploadTime int64
	Rank       int

	HasSendFlower bool
}

type PartyItemRaw struct {
	ID            int64
	PartyID       int64
	Username      string
	UID           int64
	Nickname      string
	Head          string
	Img           string
	PartyUsername string
	PartyNickname string
	Subject       string
	Desc          string
	BgBannerType  int

	BannerFile     string
	BackgroundFile string

	Clothes         string
	ModelNo         int
	PartnerType     int //0:no 1:all 2:friend
	AllowPartners   string
	Partner         string
	PartnerQinmi    int
	PartnerNickname string
	PartnerHead     string
	DyeMap          string

	PartnerClothes string
	PartnerModelNo int
	PartnerDyeMap  string

	Popularity int
	LikeCnt    int
	UnlikeCnt  int
	FlowerCnt  int
	UploadTime int64
	Rank       int

	HasSendFlower bool
}

type PartyComment struct {
	ID           int64
	SendUsername string
	SendNickname string
	Content      string
	Time         int64
}
