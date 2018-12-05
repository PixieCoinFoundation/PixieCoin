package types

type FriendInfo struct {
	UID            int64
	Username       string
	Nickname       string
	Head           string
	VIP            int
	BoyClothesCnt  int
	GirlClothesCnt int
	DesignerLevel  int
	PKPoint        int
	Level          int
	Exp            int
	// Type           int //0 add friend 1 remove friend

	Qinmi              int
	QinmiLastUpdateDay string
	// LastInviteTime     int64

	CanInvite bool
	QinmiUp   bool

	PassedCount int

	ThirdUsername string
}

type UserInfo struct {
	Username string
	Nickname string
	Head     string
}

type QinmiInfo struct {
	Nickname string
	Qinmi    int
}
