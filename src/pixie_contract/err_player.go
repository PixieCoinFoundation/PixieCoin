package pixie_contract

var PIXIE_ERR_PLAYER_CURRENCY_LIMIT = PixieRespInfo{
	RespCodeStr: "3000",
	RespDesc:    "player currency not enough",
	RespDescCN:  "货币不足",
}

var PIXIE_ERR_PLAYER_NOT_FOUND = PixieRespInfo{
	RespCodeStr: "3001",
	RespDesc:    "player not found",
	RespDescCN:  "玩家不存在",
}

var PIXIE_ERR_NICKNAME_ALREADY_EXIST = PixieRespInfo{
	RespCodeStr: "NICKNAME_ALREADY_EXIST",
	RespDesc:    "nickname already exist",
	RespDescCN:  "昵称已存在",
}

var PIXIE_ERR_SET_NICKNAME = PixieRespInfo{
	RespCodeStr: "ERR_SET_NICKNAME",
	RespDesc:    "set nickname error",
	RespDescCN:  "设置昵称时出错",
}

var PIXIE_ERR_NICKNAME_NOT_CHANGE = PixieRespInfo{
	RespCodeStr: "NICKNAME_NOT_CHANGE",
	RespDesc:    "nickname not change",
	RespDescCN:  "昵称无变化",
}

var PIXIE_ERR_SET_STREET_NAME = PixieRespInfo{
	RespCodeStr: "ERR_SET_STREET_NAME",
	RespDesc:    "set street name error",
	RespDescCN:  "设置街区命名出错",
}

var PIXIE_ERR_STREET_NAME_NOT_CHANGE = PixieRespInfo{
	RespCodeStr: "STREET_NAME_NOT_CHANGE",
	RespDesc:    "street name not change",
	RespDescCN:  "街区命名无变化",
}

var PIXIE_ERR_NICKNAME_ILLEGAL = PixieRespInfo{
	RespCodeStr: "NICKNAME_ILLEGAL",
	RespDesc:    "nickname illegal",
	RespDescCN:  "昵称非法",
}

var PIXIE_ERR_TILI_NOT_ENOUGH = PixieRespInfo{
	RespCodeStr: "TILI_NOT_ENOUGH",
	RespDesc:    "energy not enough",
	RespDescCN:  "体力不足",
}
