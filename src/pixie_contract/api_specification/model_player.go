package api_specification

import (
	"encoding/json"
	"time"
)

// 玩家的各个属性
type Status struct {
	//任何地方都不要使用ID
	Id int64 `xorm:"pk not null 'id'"`

	Username   string `xorm:"not null unique 'username'"`
	Uid        int64  `xorm:"not null 'uid'"`
	Nickname   string `xorm:"not null 'nickname'"`
	StreetName string `xorm:"not null 'street_name'"`
	Head       string `xorm:"not null 'head'"`
	Sex        int    `xorm:"not null 'sex'"`

	Level int `xorm:"not null 'level'"` //等级
	Exp   int `xorm:"not null 'exp'"`   //知名度

	//金币
	Money int64   `xorm:"not null 'money'"`
	Pxc   float64 `xorm:"not null 'pxc'"`

	Tili      int       `xorm:"not null 'tili'"`
	MaxTili   int       `xorm:"not null 'max_tili'"`
	TiliRTime int64     `xorm:"not null 'tili_r_time'"`
	MapID     string    `xorm:"not null 'map_id'"`
	Updated   time.Time `xorm:"not null updated"`
	//json fields
	SkinMap map[string]bool `xorm:"not null 'skin_map' json"`
	//衣服 服装ID:服装附加属性(数量、最近一次购买时间)
	Clothes map[string]*ClothesDetail `xorm:"not null 'clothes' json"`

	//我的家-模特信息
	HomeShow HomeShowDetail `xorm:"not null 'home_show' json"`

	//我的套装 key是model type:value是套装列表
	SuitMap map[string][]Suit `xorm:"not null 'suit_map' json"`

	//存储剧情完成情况 key是关卡id value是关卡的完成信息
	RecordMap map[string]*PixieRecord `xorm:"not null 'record_map' json"`

	//和npc的友好度情况
	NpcFriendship map[string]*FriendshipDetail `xorm:"not null 'npc_friendship' json"`

	//npc buff
	NpcBuff NpcBuffResult `xorm:"not null 'npc_buff' json"`

	Scripts []string `xorm:"not null 'scripts' json"`
	Ending  []int    `xorm:"not null 'ending' json"`
	Tasks   []*Task  `xorm:"not null 'tasks' json"`

	//我的街运营相关数据
	StreetOperate *StreetOperateResult `xorm:"not null 'street_operate' json"`

	DayExtraInfo Extra   `xorm:"not null 'day_extra_info' json"`
	ExtraInfo1   Extra1  `xorm:"not null 'extra1' json"`
	ExtraInfo2   Extra2  `xorm:"not null 'extra2' json"`
	BanInfo      BanData `xorm:"not null 'ban_info' json"`
}

type NpcBuffResult struct {
	NpcBuffList []NpcBuffListDetail

	//人流容量buff
	BuffStreetVisitCap float64

	//餐饮buff
	BuffRestaurantOffer              float64
	BuffRestaurantProfit             float64
	BuffRestaurantCleanNeed          float64
	BuffSpecifiedRestaurantOffer     map[int64]float64
	BuffSpecifiedRestaurantProfit    map[int64]float64
	BuffSpecifiedRestaurantCleanNeed map[int64]float64

	//娱乐buff
	BuffEntertainOffer              float64
	BuffEntertainProfit             float64
	BuffEntertainCleanNeed          float64
	BuffSpecifiedEntertainOffer     map[int64]float64
	BuffSpecifiedEntertainProfit    map[int64]float64
	BuffSpecifiedEntertainCleanNeed map[int64]float64

	//我的店buff
	BuffMyShopProfit float64

	//全局buff
	BuffAllProfit      float64
	BuffAllCleanNeed   float64
	BuffAllCleanPower  float64
	BuffAllCleanSalary float64

	//清洁设施101 buff
	BuffCleanObject101Power  float64
	BuffCleanObject101Salary float64

	//清洁设施201 buff
	BuffCleanObject201Power  float64
	BuffCleanObject201Salary float64

	//解锁buff
	BuffUnlockBuilding map[int64]bool
}

type NpcBuffListDetail struct {
	NpcID string
	Level int
}

type FriendshipDetail struct {
	Level int //友好等级
	Exp   int //好友度
}

type StreetOperateDemandDetail struct {
	//当前状态
	EntertainNeed   float64 //娱乐需求
	EntertainOffer  float64 //娱乐能力
	RestaurantNeed  float64 //餐饮需求
	RestaurantOffer float64 //餐饮能力

	//需求类设施开店状态
	EntertainParamOpen  float64
	RestaurantParamOpen float64

	//需求类设施开店下一级状态
	EntertainParamOpenNextLevel  float64
	RestaurantParamOpenNextLevel float64

	CleanNeed  float64 //卫生需求
	CleanOffer float64 //卫生能力

	EntertainLandID  int64
	RestaurantLandID int64
}

//我的街运营的数据结果
type StreetOperateResult struct {
	StreetOperateDemandDetail

	WalletLeft float64 //钱包剩余值
	WalletMax  float64 //钱包最大值

	ProfitInHour           float64 //每小时总利润(收入-支出)
	MyShopIncomeInHour     float64 //我的店每小时总收入
	EntertainIncomeInHour  float64 //娱乐设施每小时总收入
	RestaurantIncomeInHour float64 //餐饮设施每小时总收入
	CleanSalaryInHour      float64 //清洁设施每小时总支出

	VisitCap float64 //人流容量

	GeneralParam float64 //总体评价

	CleanPower map[string]*CleanObjectDetail //key是clean设备ID，value是雇佣的数量

	LastIncomeMinute string //上次收入计算时间。格式："20180806-1504"
}

type CleanObjectDetail struct {
	Count int
	Level int
}

//妖精购物街游戏剧情存储信息
type PixieRecord struct {
	PlotRead bool `json:"pr"`

	Score     float64 `json:"s"`
	BestScore float64 `json:"bs"`

	//rank为空字符串代表只是看过此关卡并没有通关
	Rank string `json:"r"`

	Clothes     string `json:"cs"`
	BestClothes string `json:"bcs"`

	HistoryRank map[string]int `json:"hr"`
}

type HomeShowDetail struct {
	Left       ModelDetail
	Right      ModelDetail
	Middle     ModelDetail
	Background string
	ShowOrder  []string
}

type ModelDetail struct {
	ModelType string
	ModelSuit SuitDetail
}

type SuitDetail struct {
	Faxing    string `json:"Faxing,omitempty"`
	Waitao    string `json:"Waitao,omitempty"`
	Shangyi   string `json:"Shangyi,omitempty"`
	Xiazhuang string `json:"Xiazhuang,omitempty"`
	Wazi      string `json:"Wazi,omitempty"`
	Xiezi     string `json:"Xiezi,omitempty"`
	Liantifu  string `json:"Liantifu,omitempty"`
	Makeup    string `json:"Makeup,omitempty"`

	//饰品 单类最多可穿两个
	Head   []string `json:"Head,omitempty"`
	Face   []string `json:"Face,omitempty"`
	Ear    []string `json:"Ear,omitempty"`
	Neck   []string `json:"Neck,omitempty"`
	Hand   []string `json:"Hand,omitempty"`
	Inhand []string `json:"Inhand,omitempty"`
	Pifu   []string `json:"Pifu,omitempty"`
	Other  []string `json:"Other,omitempty"`
}

type Suit struct {
	SuitDetail
	SuitMD5  string //套装MD5
	SuitName string
	SuitImg  string
}

func (s *Suit) GenIDPartMap() (res map[string]string) {
	res = make(map[string]string)

	if s.Faxing != "" {
		res[s.Faxing] = PAPER_TYPE_FAXING
	}

	if s.Waitao != "" {
		res[s.Waitao] = PAPER_TYPE_WAITAO
	}

	if s.Shangyi != "" {
		res[s.Shangyi] = PAPER_TYPE_SHANGYI
	}

	if s.Xiazhuang != "" {
		res[s.Xiazhuang] = PAPER_TYPE_XIAZHUANG
	}

	if s.Wazi != "" {
		res[s.Wazi] = PAPER_TYPE_WAZI
	}

	if s.Xiezi != "" {
		res[s.Xiezi] = PAPER_TYPE_XIEZI
	}

	if s.Liantifu != "" {
		res[s.Liantifu] = PAPER_TYPE_LIANTIFU
	}

	if s.Makeup != "" {
		res[s.Makeup] = PAPER_TYPE_MAKEUP
	}

	if len(s.Head) > 0 {
		for _, v := range s.Head {
			res[v] = PAPER_TYPE_HEAD
		}
	}

	if len(s.Face) > 0 {
		for _, v := range s.Face {
			res[v] = PAPER_TYPE_FACE
		}
	}

	if len(s.Ear) > 0 {
		for _, v := range s.Ear {
			res[v] = PAPER_TYPE_EAR
		}
	}

	if len(s.Neck) > 0 {
		for _, v := range s.Neck {
			res[v] = PAPER_TYPE_NECK
		}
	}

	if len(s.Hand) > 0 {
		for _, v := range s.Hand {
			res[v] = PAPER_TYPE_HAND
		}
	}

	if len(s.Inhand) > 0 {
		for _, v := range s.Inhand {
			res[v] = PAPER_TYPE_INHAND
		}
	}

	if len(s.Pifu) > 0 {
		for _, v := range s.Pifu {
			res[v] = PAPER_TYPE_PIFU
		}
	}

	if len(s.Other) > 0 {
		for _, v := range s.Other {
			res[v] = PAPER_TYPE_OTHER
		}
	}

	return
}

func (s Status) TableName() string {
	return "status"
}

type Extra struct {
	IAPed                   int            `json:"iaped"`
	DaySpecialMissionCntMap map[string]int `json:"dsmcm"`
	UploadCustomCntMap      map[string]int `json:"uccm"`   //上传设计师数量情况
	DayActivityGot          map[string]int `json:"dag"`    //日活跃度获取数
	TotalIAPCost            int            `json:"tic"`    //总充值数
	ClothesShareCount       int            `json:"csc"`    //服装分享次数
	WarDateWinOnceRewardMap map[string]int `json:"wdworm"` //
	MarkGotRewards          []int
	DayFreeLotteryMap       map[string]bool
	DayCopyReportCnt        int
	DayShensuMap            map[string]int

	DayTheaterCntMap    map[string]int
	DayBuyTheaterCntMap map[string]int

	DayQinmiAddBoardMap   map[string]int
	DayQinmiReplyBoardMap map[string]int
	MonthCardStart        int64
	MonthCardEnd          int64
	MonthCardIndex        int
	PKCnt                 int
	UnlockModels          []int
	PKLoseClothes         []string

	DefaultClothes    []string
	DefaultLeftModel  int
	DefaultRightModel int
	DefaultSceneID    string
	DefaultShowModel  int

	RegTime                 int64
	ReceivedCronMails       []int64
	ClothesPieceDayInCnt    int
	ClothesPieceDayOutCnt   int
	MagicSpeedLeftCnt       int
	UploadHelpDayCnt        int
	UploadCustomDayCnt      int
	UploadCosItemDayCnt     int
	DayEventDropNum         int
	DayBuyTiliCnt           int
	NewPlayerLoginDayCnt    int
	DayDyeCnt               int
	DayFreeDiamondBonusUsed bool
	DayBuQianCnt            int
	LastSendHopeTime        int64
	FirstRewardMap          map[string]int
	EventLinearProgress     map[string]int      //id:progress
	EventDiscreteProgress   map[string][]string //id:process accomplished
	EventComplexProgress    map[string]string

	DayVerifyNum map[string]int //审核 key日期 val次数
	VerifyCount  int            //累计审核

}

type Extra1 struct {
	NicknameSet        bool //是否首次设置了昵称
	DayInviteFriendCnt int  //每日邀请好友次数
	JuhuaCnt           int  //菊花数量
	JuhuaChaCnt        int  //菊花茶数量

	VerifyBenefit    int //累计评审收益
	YesterdayBenefit int //昨日审核收益

	MonthCardRewardMap map[string]int //月卡奖励是否已发放

	PKWeekPoints  map[string]int //pk周积分情况
	PKGWeekPoints map[string]int //pk赛季积分情况

	WeekWarClothesReward map[string][]int //公会周抽奖领取情况

	FirstPKReward bool //首次PK奖励是否已领取
	IAPReturned   bool //公测返利是否已领取
	YuyueRewarded bool //预约奖励是否已领取

	PartnerInviteMap     map[string]int64 //舞会邀请好友情况
	GuanzhuDesignerCount int              //关注设计师数量
	LastQuitGuildTime    int64            `json:"lqgt"` //上次退出公会时间
	LastReadGuildBoardID int64            `json:"lrgbid"`

	DelTime int64 //账号删除时间

	Item1  int //小菊的纽扣
	Cloth  int //布料
	Wire   int //五彩线
	Needle int //魔术针

	DayRTPartyRewardMap map[string]int

	MonthToken string             `json:"mt"`
	MonthIap   map[string]float64 `json:"miapm"`

	LoginDayCnt    int `json:"ldc"`
	IapDiamondLeft int `json:"idl"`

	DayStampChangeCntMap map[string]int

	RewardedSuitMap map[string]int `json:"rsm"`

	KorInviteRewardMap map[string]int `json:"kirm"`

	IgnoreNotice            bool `json:"in"`
	IgnorePush              bool `json:"iph"`
	IgnoreNightPush         bool `json:"inp"`
	IgnoreGuildWarPush      bool `json:"igwph"`
	IgnoreGuildPush         bool `json:"igph"`
	AcceptFriendRequestPush bool `json:"afrph"`
	IgnoreHopeDonePush      bool `json:"ihdph"`
	AcceptPartyInvitePush   bool `json:"apiph"`
	IgnoreNewCustomPush     bool `json:"incph"`
	IgnoreCartPush          bool `json:"icp"`

	Backgrounds                []string `json:"bgs"`
	DesignCoinPriceTokenInited bool     `json:"dcpti"`

	DesignCoin int //设计师币

	KakaoHead string

	KorInitPlayerRewarded bool `json:"kipr"`
	KorRecordCZRewarded   bool `json:"krczr"`

	DesignCoinPriceToken       int            //设计币定价券
	DesignCoinPriceTokenBuyMap map[string]int //设计比定价券购买次数

	AttendRTParty map[string]int64

	ETHAccount string
	ETHUnlock  bool

	//支付密码密文
	ETHPayPwdEncrypt string
	//支付密码绑定的邮箱
	ETHPayPwdEmail string

	//发送验证码的邮箱
	ETHPayPwdVerifyCodeEmailSend string
	//设置支付密码时的验证码
	ETHPayPwdVerifyCode string
	//设置支付密码时的验证码失效时间
	ETHPayPwdVerifyCodeExpireTime int64
}

type Extra2 struct {
	//累计通过评审获得的pxc收益
	TotalPxcRewardOnVerify float64 `json:"tprov"`
	//计算累计获得评审获得收益的截止notify id
	TotalPxcRewardOnVerifyDueNotifyID int64 `json:"tprovdvi"`
}

type ClothesInfo struct {
	ClothesID string `json:"I"`
	Count     int    `json:"C"`
	Time      int64  `json:"T, omitempty"` // 最近一次购买的时间
}

type ClothesDetail struct {
	Count int   `json:"C"`
	Time  int64 `json:"T"`
}

type Task struct {
	TaskID   int    `bson:"taskID" json:"taskID"`     // 任务id
	Progress int    `bson:"progress" json:"progress"` // 任务进度
	Desc     string `bson:"desc" json:"-"`            // 任务描述
	Target   int    `bson:"target" json:"Target"`     // 目标
	Type     int    `bson:"type" json:"type"`         // 任务类型
	Gold     int64  `bson:"gold" json:"gold"`
	Diamond  int    `bson:"diamond" json:"diamond"`
	Fastener int    `bson:"fastener" json:"fastener"`

	RewardGot bool `bson:"rewardGot" json:"rewardGot"` // 奖励是否已经被获得
	Info      int  `bson:"info" json:"-"`              // 任务的一些额外信息
}

type Record struct {
	LevelID     string         `bson:"lid" json:"lid"`
	Score       int            `bson:"score" json:"s"`
	Rank        string         `bson:"rank" json:"r"`
	Clothes     string         `bson:"rank" json:"clothes"`
	BestClothes string         `bson:"rank" json:"bestClothes"`
	BestScore   int            `json:"bestScore"`
	Lock        bool           `bson:"lock" json:"lock"`
	NewUnlock   bool           `bson:"newUnlock" json:"newUnlock"`
	GayCount    int            `bson:"gayCount" json:"gayCount"` // 过关搭配中，man属性小于0的服装件数
	HistoryRank map[string]int `bson:"historyRank" json:"historyRank"`
}

func (t Record) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["lid"] = t.LevelID

	if t.Score != 0 {
		m["s"] = t.Score
	}

	if t.Rank != "" {
		m["r"] = t.Rank
	}

	if t.Clothes != "" {
		m["clothes"] = t.Clothes
	}

	if t.BestClothes != "" {
		m["bestClothes"] = t.BestClothes
	}

	if t.Lock {
		m["locl"] = t.Lock
	}

	if t.NewUnlock {
		m["newUnlock"] = t.NewUnlock
	}

	if t.GayCount != 0 {
		m["gayCount"] = t.GayCount
	}

	if t.BestScore != 0 {
		m["bestScore"] = t.BestScore
	}

	if len(t.HistoryRank) > 0 {
		m["historyRank"] = t.HistoryRank
	}

	return json.Marshal(m)
}

func (t Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"taskID":    t.TaskID,
		"progress":  t.Progress,
		"rewardGot": t.RewardGot,
	})
}

type BanData struct {
	BanStartTime int64
	BanEndTime   int64
	BanReason    string
}
