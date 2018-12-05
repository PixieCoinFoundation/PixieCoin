package api_specification

//图纸评审
const HTTP_PLAYER_PAPER_VERIFY = "api/paperVerify"

type PaperVerifyReq struct {
	PaperID int64
	Score   int

	//标签
	Tag1 string
	Tag2 string

	//风格
	Style1 string
	Style2 string
}

type PaperVerifyResp struct {
	VerifiedNum      int //参与评审人数
	VerifyCountLimit int //参与评审人数限制
	VerifiedCount    int //已经评审次数
	VerifyLimit      int //每日评审限制次数

	//下一件待审核作品的数据
	IsReport bool        //下一件作品我是否举报过
	Paper    DesignPaper //下一件待审核作品数据
}

//获取一个要评审的图纸
const HTTP_PLAYER_GET_PAPER_VERIFY string = "api/getPaperVerify"

type GetVerifyPaperReq struct {
}
type GetVerifyPaperResp struct {
	VerifiedNum      int //参与评审人数
	VerifyCountLimit int //参与评审人数限制
	VerifiedCount    int //已经评审次数
	VerifyLimit      int //每日评审限制次数

	IsReport bool //我是否举报过
	Paper    DesignPaper
}

//评审结果
const HTTP_PLAYER_PAPER_NOTIFY_QUERY = "api/paperNotifyQuery"

type PaperNotifyQueryReq struct {
	Page     int
	PageSize int
}

type PaperNotifyQueryResp struct {
	Res           []PaperNotify
	YesBenefit    float64 //昨日pxc收益
	YesVerify     int     //昨日评审数量
	BenefitCount  float64 //累计pxc收益
	VerifyCount   int     //累计评审数量
	VerifiedCount int     //每日已经评审次数
	VerifyLimit   int     //每日评审限制次数
}

//查询评审详情
const HTTP_PLAYER_PAPER_VERIFIED_DETAIL_QUERY = "api/paperVerifiedDetailQuery"

type PaperVerifiedDetailReq struct {
	NotifyID  int64 //通知ID
	CheckTime int64 //评审检查日期
}

type PaperVerifiedDetailResp struct {
	Res         []PaperVerify
	RewardCount int64
}

//举报图纸
const HTTP_PLAYER_PAPER_REPORT_COPY = "api/paperReportCopy"

type PaperReportCopyReq struct {
	PaperID int64
	Reason  string
	Pic     string
	Contact string
}

type PaperReportCopyResp struct {
}

//举报查询
const HTTP_PLAYER_COPY_QUERY = "api/paperCopyQuery"

type PaperCopyQueryReq struct {
	PaperID int64
}

type PaperCopyQueryResp struct {
	List []PaperCopy
}

//支持举报
const HTTP_PLAYER_PAPER_COPY_SUPPORT = "api/paperCopySupport"

type PaperCopySupportReq struct {
	CopyID  int64
	PaperID int64
}

type PaperCopySupportResp struct {
}

//领取奖励
const HTTP_PLAYER_GET_PAPER_COPY_REWARD = "api/getPaperCopyReward"

type GetPaperCopyRewardReq struct {
	CommonReq
	NotifyID int64 //奖励ID
}

type GetPaperCopyRewardResp struct {
	CurrentGold int64
	CurrentPxc  float64
}
