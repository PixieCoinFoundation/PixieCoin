package api_specification

//获取杂志列表 只返回ID banner信息 【推荐套装】信息
const HTTP_PLAYER_LIST_SUBJECT = "api/listSubject"

type GetSubjectReq struct {
	sType int //0是最新的一个，1.是所有杂志
}

type GetSubjectResp struct {
	SubjectInfo []Subject
}

//获取杂志-【推荐单品】页
const HTTP_PLAYER_GET_SUBJECT_DETAIL = "api/getSubjectDetail"

type GetSubjectDetailReq struct {
	SubjectID int64
	IDList    []string //已经存在的id数组
	Page      int      //页码，页码打于1只返回第三个列表
}
type GetSubjectDetailResp struct {
	HasHistory  bool
	ItemsTop    []*TopPapers
	ItemsMedium []*ItemsPaper //18个推荐单品
	ItemsDown   []*ItemsPaper
}

//获取杂志-【推荐套装】信息
const HTTP_PLAYER_GET_SUBJECT_SUIT = "api/getSubjectSuit"

type GetSubjectSuitReq struct {
	SubjectID int64
}

type GetSubjectSuitResp struct {
	SuitList []RecommendClothes //套装列表
}

//获取专题列表
const HTTP_PLAYER_GET_SUBJECT_TOPIC = "api/getSubjectTopic"

type GetSubjectTopicReq struct {
	SubjectID int64
	Page      int64
}

type GetSubjectTopicResp struct {
	Topic []SubjectTopic //专题列表
}

//历史推荐单品(带link)

const HTTP_PLAYER_GET_HISTORY_SUBJECT_ITEM = "api/getHistorySubjectItem"

type GetHistorySubjectItemReq struct {
	Page int64
}

type GetHistorySubjectItemResp struct {
	ItemsTop []*TopPapers
}
