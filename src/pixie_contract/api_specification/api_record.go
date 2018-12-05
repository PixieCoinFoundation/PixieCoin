package api_specification

const HTTP_PLAYER_SAVE_RECORD_HANDLE = "api/saveRecord"

type SaveRecordReq struct {
	CommonReq

	RecordID string  //关卡ID
	Score    float64 //分数
	Rank     string
	Clothes  string
}

type SaveRecordResp struct {
	//玩家个人属性
	Tili    int
	MaxTili int
	Level   int
	Exp     int

	Res PixieRecord

	//新的友好度
	NpcFriendship map[string]*FriendshipDetail

	//新的buff
	NpcBuff NpcBuffResult

	//我的街目前状态
	StreetDetail *StreetOperateResult
}

const HTTP_PLAYER_VIEW_RECORD_HANDLE = "api/viewRecord"

type ViewRecordReq struct {
	RecordID string //关卡ID
}

type ViewRecordResp struct {
}
