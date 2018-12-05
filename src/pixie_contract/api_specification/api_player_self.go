package api_specification

//修改玩家昵称
const HTTP_PLAYER_CHANGE_NICKNAME = "api/changeNickname"

type ChangeNicknameReq struct {
	NewNickname string
}

type ChangeNicknameResp struct {
}

//修改街区名称
const HTTP_PLAYER_CHANGE_STREET_NAME = "api/changeStreetName"

type ChangeStreetNameReq struct {
	NewStreetName string
}

type ChangeStreetNameResp struct {
}
