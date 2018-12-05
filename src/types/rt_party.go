package types

type RTPartyHostStatus int

const (
	//房间非正常状态
	// RT_PARTY_HOST_STATUS_TIMEOUT RTPartyHostStatus = -1 //房间在规定时间内没有足够玩家进入 超时

	//房价正常状态
	RT_PARTY_HOST_STATUS_START          RTPartyHostStatus = 1 //房间开启 等待玩家进入
	RT_PARTY_HOST_STATUS_SUBJECT_SELECT RTPartyHostStatus = 2 //房间开启 主题选择阶段
	RT_PARTY_HOST_STATUS_SUBJECT_VOTE   RTPartyHostStatus = 3 //房间开启 主题投票阶段
	RT_PARTY_HOST_STATUS_DRESS          RTPartyHostStatus = 4 //房间开启 搭配阶段
	RT_PARTY_HOST_STATUS_DRESS_VOTE     RTPartyHostStatus = 5 //房间开启 搭配投票阶段
	RT_PARTY_HOST_STATUS_DRESS_VOTE_END RTPartyHostStatus = 6 //房间开启 搭配投票结束 投票公示阶段
	// RT_PARTY_HOST_STATUS_CLOSE          RTPartyHostStatus = 7 //房间关闭
)

type RTPartySubject struct {
	Name      string
	OfferTime int64
	GotVote   int
}

type RTPartyHost struct {
	ID                  int64
	StartTime           int64 //创建时间
	BeginTime           int64 //人满后的正式开始时间
	AllOfferSubjectTime int64 //全部成员提供完主题的时间
	AllVoteSubjectTime  int64 //全部成员主题投票结束的时间
	AllFinishDressTime  int64 //全部成员完成搭配的时间
	AllVoteDressTime    int64 //全部成员搭配投票完的时间

	DefaultSubjects []RTPartySubject //默认主题列表

	FinalSubject string             //最终选出的主题
	FinalRank    RTPartyHostJoiners //最终排名

	Joiners []*RTPartyHostJoiner //舞会参与者
}

type RTPartyHostJoiner struct {
	Username string
	Nickname string
	Head     string
	Level    int
	VIP      int

	Left bool //是否已主动离开

	OfferSubject   RTPartySubject
	SubjectVoteCnt int

	LastSyncTime int64

	DyeCnt        int
	DyeMap        map[string][7][4]float64
	Clothes       []string
	ModelNo       int
	BackgroundID  string
	Img           string
	DressWord     string
	DressWordType int //0:普通 1:语音

	//extra attribute
	OfferSubjectTime int64
	VoteSubjectTime  int64
	VoteDressTime    int64

	DressGotVote  int   //搭配得到投票数
	DressVoteCnt  int   //搭配投票次数
	DressDoneTime int64 //搭配结束时间

	DressVoteUsername1 string
	DressVoteUsername2 string
}

type RTPartyHostJoinerExtra struct {
	DyeCnt        int    `redis:"dyeCnt"`
	DyeMap        string `redis:"dyeMap"`
	Clothes       string `redis:"clothes"`
	ModelNo       int    `redis:"modelNo"`
	BackgroundID  string `redis:"backgroundID"`
	Img           string `redis:"img"`
	DressWord     string `redis:"dressWord"`
	DressWordType int    `redis:"dressWordType"`
	DressDoneTime int64  `redis:"dressDoneTime"`
	LastSyncTime  int64  `redis:"lastSyncTime"`
}

type RTPartyChat struct {
	U  string //用户名
	C  string //内容
	T  int64  //时间
	Te int    //0:普通聊天 1:语音
}

type RTPartyHostJoiners []*RTPartyHostJoiner

func (a RTPartyHostJoiners) Len() int {
	return len(a)
}

func (a RTPartyHostJoiners) Less(i, j int) bool {
	if a[i].DressGotVote < a[j].DressGotVote {
		return false
	} else if a[i].DressGotVote == a[j].DressGotVote {
		if a[i].Left {
			return false
		}

		if a[j].Left {
			return true
		}

		if a[i].DressDoneTime == 0 {
			return false
		}

		if a[j].DressDoneTime == 0 {
			return true
		}

		return a[i].DressDoneTime < a[j].DressDoneTime
	}
	return true
}

func (a RTPartyHostJoiners) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
