package types

import (
	. "pixie_contract/api_specification"
)

// 登陆
type LoginReq struct {
	BiliUsername string `json:"biliUsername"`
	Username     string `json:"username"`
	AccessToken  string `json:"token"`
	DeviceToken  string `json:"deviceToken"`
	UID          int64  `json:"uid"`
	GFToken      string `json:"gfToken"`

	Platform string `json:"platform"`
	DeviceId string `json:"deviceId, omitempty"`

	Channel      string `json:"channel"`
	ThirdChannel string `json:"thirdChannel"`

	GuestLogin      bool   `json:"guestLogin"`      //游客登录
	DeviceModelType string `json:"deviceModelType"` //设备名称

	KorRevertUnregister bool `json:"korRevertUnregister"`
}

//LoginKorLog 登录韩国日志
type LoginKorLog struct {
	Guid            int64
	Uuid            string
	Cuid            int64
	Gamecode        int
	Nickname        string
	Sex             bool
	Age             int
	LevYn           string
	UserLv          int
	Authtoken       int
	Oscode          int
	Ostopverno      int
	Tempuserid      string
	TokenredisDate  string
	Registdate      string
	Mappingdate     string
	IP              string
	GuestLogin      bool
	DeviceModelType string
}

//LoginKorLog1 登录韩国日志1
type LoginKorLog1 struct {
	Logindate         string
	Logindateseq      int
	Cuid              int64
	Isloginsuccessdiv int
	Platformuserid    string
	Loginrootcode     int
	GameCode          int
	Deviceid          int
	Platformcode      int
	Oscode            int
	Ostopverno        int
}

//OrderKorLog 订单日志
type OrderKorLog struct {
	Paylogid       int64
	Eventdate      string
	Logtype        int64
	Guid           int64
	Platformuserid string
	Platformcode   int
	Oscode         int
	Storetype      int
	Orderid        string
	Paydate        string
	Paytool        int
	Gamecode       int
	Gamename       string
	Itemname       int
	Itemcode       int
	Currentcycode  string
	Amount         float64
	Paycountry     string
}

type LoginReqExtra struct {
	BiliUsername string
	Username     string
	AccessToken  string
	DeviceToken  string
	UID          int64
	GFToken      string

	Platform string
	DeviceId string

	Channel      string
	ThirdChannel string

	GuestLogin      bool
	DeviceModelType string

	KorRevertUnregister bool
	RemoteAddr          string
	ThirdUsername       string
	ServerAddress       string
	Nickname            string
	Level               int
	KOROSCode           int
	RegTime             string
	KORGameCode         int
	IP                  string
}

//过关
type SaveRecordRequest struct {
	//必填参数
	LevelID     string
	SkipPlot    bool //是否跳过剧情
	OptionScore int  //选项得分

	Exit bool //未完成即退出 此值为true时后面的参数会被无视

	CurRecord    Record `json:"record"`
	UnlockRecord Record `json:"unlockRecord"`
	MLevel       int    `json:"mLevel"`
	MNo          int    `json:"mNo"`
	NeedUnlock   bool   `json:"needUnlock"`
	Planet       int    `json:"planet"`
	DeviceId     string `json:"deviceId, omitempty"`
}

//GuildRecordExtra 公会信息
type GuildRecordExtra struct {
	GuildID    int64  //公会ID
	GuildName  string //公会名称
	GuildOwer  int64  //会长
	UID        int64  //会长UID
	Extra      string //预留字段
	Nickname   string //昵称
	CreateTime string //创建时间
	CloseTime  string //解散公会状态
	CloseType  string //解散公会类型 1会长解散，0为GMTool解散
	Reason     string
}

//GuildJoinExtra 申请者信息
type GuildJoinExtra struct {
	UID       int64 //申请者
	Nickname  string
	KaKaoid   string
	JoinTime  string //加入时间
	IsAllower string //是否允许
	GuildId   int64  //公会ID
	GuildName string //公会名称
}

type SaveRecordExtra struct {
	LevelID        string
	SkipPlot       bool //是否跳过剧情
	OptionScore    int  //选项得分
	Exit           bool //未完成即退出 此值为true时后面的参数会被无视
	Score          int
	Rank           string
	MLevel         int
	MNo            int
	GotGold        int
	GotExp         int
	ServerAddress  string
	Clothes        string
	Replay10       bool
	MagicSpeed     bool
	IsRewardPirces bool
	ReardPircesID  string
	DeviceId       string
}

type PKExtra struct {
	TimeCost int
	Win      string
	Bot      bool

	//赛季积分信息
	PrePKPoints int
	NowPKPoints int
	PrePKLevel  int
	NowPKLevel  int

	//周榜积分信息
	WeekPrePKPoints int
	WeekNowPKPoints int
	WeekPrePKLevel  int
	WeekNowPKLevel  int

	LevelID    string
	WeekToken  string
	GWeekToken string
	PKType     int //0:xiuxian 1:jingji

	WinClothes  string
	LoseClothes string

	DeviceId         string `json:"deviceId, omitempty"`
	ServerAddress    string
	OpponentUsername string //pk对手的用户名
	WinCnt           int
}

//UploadClothesRequestExtra 上传衣服cid
type UploadClothesRequestExtra struct {
	UploadClothesRequest
	CID int64 `json:"cid"`
}
type UploadClothesRequest struct {
	CName    string     `json:"cname, omitempty"`
	ModelNo  int        `json:"modelNo, omitempty"`
	Type     int        `json:"type, omitempty"`
	Desc     string     `json:"desc, omitempty"`
	CloSet   CustomFile `json:"clothesSet, omitempty"`
	DeviceId string     `json:"deviceId, omitempty"`

	Remark string `json:"remark, omitempty"`
	Tag1   int    `json:"tag1, omitempty"`
	Tag2   int    `json:"tag2, omitempty"`
	ZValue int    `json:"zValue, omitempty"`

	//0:正常 1:非卖品
	PlayerConfig int `json:"playerConfig"`
	/*
	   tag编号
	   1：保暖
	   2：正式
	   3：紧身
	   4：鲜艳
	   5：深色
	   6：可爱
	   7：帅气
	   8：粗犷
	   9：华丽
	   10：另类
	   11：性感
	   12：运动

	   13：温暖
	   14：严肃
	   15：修身
	   16：亮眼
	   17：黑暗
	   18：活泼
	   19：男人味
	   20：不羁
	   21：高贵
	   22：个性
	   23：迷人
	   24：活力

	   25：清凉
	   26：休闲
	   27：宽松
	   28：朴素
	   29：浅色
	   30：成熟
	   31：妩媚
	   32：斯文

	   33：凉快
	   34：随性
	   35：飘逸
	   36：低调
	   37：洁白
	   38：稳重
	   39：女人味
	   40：内敛
	*/
}
