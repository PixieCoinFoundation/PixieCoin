package pixie_contract

var PIXIE_ERR_SUIT_LIMIT_WRONG = PixieRespInfo{
	RespCodeStr: "4001",
	RespDesc:    "suit count limit wrong",
	RespDescCN:  "套装件数超过限制",
}

var PIXIE_ERR_SUIT_ALREADY_EXIST = PixieRespInfo{
	RespCodeStr: "4002",
	RespDesc:    "suit already exist",
	RespDescCN:  "套装已存在",
}

var PIXIE_ERR_SUIT_HANDLE_UNKNOW_WRONG = PixieRespInfo{
	RespCodeStr: "4003",
	RespDesc:    "suit handle unknow wrong",
	RespDescCN:  "套装操作未知错误",
}

var PIXIE_ERR_SUIT_NOT_EXIST = PixieRespInfo{
	RespCodeStr: "4004",
	RespDesc:    "suit not exist",
	RespDescCN:  "套装不存在",
}

var PIXIE_UPDATE_SUIT_TYPE_WRONG = PixieRespInfo{
	RespCodeStr: "4005",
	RespDesc:    "update suit type wrong",
	RespDescCN:  "更新套装type错误",
}

var PIXIE_ERR_NOT_HAVE_CLOTHES = PixieRespInfo{
	RespCodeStr: "4006",
	RespDesc:    "player do not have clothes",
	RespDescCN:  "未拥有该服装",
}

var PIXIE_ERR_SUIT_NOT_LEGAL = PixieRespInfo{
	RespCodeStr: "4007",
	RespDesc:    "suit content not legal",
	RespDescCN:  "套装部件不合法",
}
