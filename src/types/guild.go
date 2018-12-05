package types

type GMGuild struct {
	Info      string
	Platform  string
	GuildInfo Guild
	Members   []GuildMember
}

type MatchGuildWarExtra struct {
	OpponentGID   int64
	ServerAddress string
}

type PassGuildApplyExtra struct {
	GID           int64
	ServerAddress string
	ApplyID       int64
	Username      string
	Nickname      string
}

type RejectGuildApplyExtra struct {
	ApplyID       int64
	ServerAddress string
}

type GuildExtra struct {
	GID           int64
	ServerAddress string
}

type TransferGuildReq struct {
	QuitGuild  bool
	ToUsername string
	ToNickname string
}

type TransferGuildExtra struct {
	TransferGuildReq
	GID              int64  //公会ID
	TransType        string //转让方式
	GuildName        string //公会名称
	OldGuildownerUID int64  //旧的公会长
	OldGuildNickname string //旧的公会长昵称
	ServerAddress    string
}

type KickGuildMemberExtra struct {
	Username      string
	GID           int64
	ServerAddress string
}

type SetDefendInfo struct {
	Username         string
	ShareClothesUsed []GuildClothes
}

type GuildWarSettleInfo struct {
	GID     int64
	Rank    int
	Medal   int
	Rewards []MemberReward
}

type MemberReward struct {
	Username      string
	RewardDiamond int
}

type WarLog struct {
	//from 进攻 to 防守
	FromGID      int64
	ToGID        int64
	FromUsername string
	FromNickname string
	ToUsername   string
	ToNickname   string
	WinMedal     int
	Time         int64
}

type SGuild struct {
	ID       int64
	MedalCnt int
	GWT      string
	Zombie   bool
}

type Guild struct {
	ID            int64
	Owner         string
	OwnerNickname string
	Name          string
	Desc          string
	LogoType      int

	MemberCnt    int
	MedalCnt     int //赛季钻冕
	Activity     int
	PKRank       int
	ActivityRank int
	NewApplyCnt  int
	WarMemberCnt int

	WarDate     string
	WarOpponent int64

	Lottery1s string
	Lottery2s string
	Lottery3s string
	Lottery4s string

	Lottery1 GuildLotteryInfo
	Lottery2 GuildLotteryInfo
	Lottery3 GuildLotteryInfo
	Lottery4 GuildLotteryInfo

	Zombie bool
}

type GuildLotteryInfo struct {
	LotteryUsername string `json:"lu"`
	LotteryNickname string `json:"ln"`
	LotteryClothes  string `json:"lc"`
	Week            string `json:"wk"`
}

type GuildApply struct {
	ID        int64
	GID       int64
	GName     string
	Username  string
	Nickname  string
	VIP       int
	Level     int
	Head      string
	ApplyTime int64
}

type SApply struct {
	ID        int64
	ApplyTime int64
}

type GuildMember struct {
	GID             int64
	Username        string
	Nickname        string
	Head            string
	VIP             int
	Level           int
	Activity        int
	IgnoreWar       bool
	IgnoreWarNotice bool
	ClothesShareCnt int

	WarDate        string
	WarClothes     string
	WarSubject     string
	WarClothesDone bool
	LeftAttackCnt  int
	LeftMedalCnt   int

	ActivityWeek string

	DefendCnt int
	WinCnt    int
	LoseCnt   int

	L1Got bool
	L2Got bool
	L3Got bool
	L4Got bool

	Zombie bool

	VicePresident bool
}

type GuildClothes struct {
	Username  string
	Nickname  string
	Head      string
	ClothesID string
	UseCnt    int
	Token     string
}

type GuildConfig struct {
	GID             int64
	IgnoreWar       bool
	IgnoreWarNotice bool

	L1Got bool
	L2Got bool
	L3Got bool
	L4Got bool
}
