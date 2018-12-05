package api_specification

//保存套装
const HTTP_PLAYER_ADD_SUIT_HANDLE = "api/addSuit"

type AddSuitReq struct {
	ModelType string //模特类型
	SuitModel Suit   //套装内容
}

type AddSuitResp struct {
	SuitMD5 string //套装MD5
}

//更新套装
const HTTP_PLAYER_UPDATE_SUIT_HANDLE = "api/updateSuit"

type UpdateSuitReq struct {
	SuitMD5   string //套装MD5
	ModelType string //模特类型
	SuitModel Suit   //套装内容
}

type UpdateSuitResp struct {
	SuitMD5 string //更新后的md5
}

//更新套装名称
const HTTP_PLAYER_UPDATE_SUIT_NAME_HANDLE = "api/updateSuitName"

type UpdateSuitNameReq struct {
	SuitMD5   string //套装md5
	ModelType string //模特类型
	SuitName  string //套装名字
}

type UpdateSuitNameResp struct {
}

//删除套装
const HTTP_PLAYER_DELETE_SUIT_HANDLE = "api/deleteSuit"

type DeleteSuitReq struct {
	SuitMD5   string
	ModelType string //模特类型
}

type DeleteSuitResp struct {
}

//更新我的家
const HTTP_PLAYER_UPDATE_HOME_SHOW_HANDLE = "api/updateHomeShow"

type UpdateHomeShowReq struct {
	HomeShow HomeShowDetail
}

type UpdateHomeShowResp struct {
}
