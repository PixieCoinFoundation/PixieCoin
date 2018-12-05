package pixie_contract

type PixieRespInfo struct {
	//错误码 后面都用这个 不一定是数字了 正常返回空字符串 非空字符串代表出错
	RespCodeStr string

	//错误描述 英文
	RespDesc string

	//错误描述 中文
	RespDescCN string
}

var PIXIE_ERR_UNKNOWN = PixieRespInfo{
	RespCodeStr: "-1",
	RespDesc:    "unknown error happen",
	RespDescCN:  "服务器未知错误",
}

var PIXIE_ERR_SHUTDOWN = PixieRespInfo{
	RespCodeStr: "-2",
	RespDesc:    "server wrong status",
	RespDescCN:  "服务器状态异常",
}

var PIXIE_ERR_WRONG_PARAMS = PixieRespInfo{
	RespCodeStr: "-3",
	RespDesc:    "client data format wrong",
	RespDescCN:  "客户端数据格式错误",
}

var PIXIE_ERR_DECRYPT = PixieRespInfo{
	RespCodeStr: "-4",
	RespDesc:    "client data decrypt wrong",
	RespDescCN:  "客户端数据解密失败",
}

var PIXIE_ERR_WRONG_JSON = PixieRespInfo{
	RespCodeStr: "-5",
	RespDesc:    "client data json wrong",
	RespDescCN:  "客户端数据格式异常",
}

var PIXIE_ERR_MAINTAIN = PixieRespInfo{
	RespCodeStr: "-6",
	RespDesc:    "server maintaining",
	RespDescCN:  "服务器维护中",
}

var PIXIE_ERR_ACCESS_TOKEN = PixieRespInfo{
	RespCodeStr: "-7",
	RespDesc:    "client access token wrong",
	RespDescCN:  "已离线，请重新登录",
}

var PIXIE_ERR_PLAYER_BAN = PixieRespInfo{
	RespCodeStr: "-8",
	RespDesc:    "player baned.contact admin.",
	RespDescCN:  "账号被锁定。请联系管理员。",
}

var PIXIE_ERR_PLAYER_BLOCK = PixieRespInfo{
	RespCodeStr: "-9",
	RespDesc:    "player blocked.contact admin.",
	RespDescCN:  "账户被禁止访问",
}

var PIXIE_ERR_CHANNEL_FORBID = PixieRespInfo{
	RespCodeStr: "-10",
	RespDesc:    "special maintain",
	RespDescCN:  "维护中请稍后",
}

var PIXIE_ERR_API_UNMARSHAL = PixieRespInfo{
	RespCodeStr: "-11",
	RespDesc:    "api data format wrong",
	RespDescCN:  "接口数据格式错误",
}

var PIXIE_ERR_API_NOT_SUPPORT = PixieRespInfo{
	RespCodeStr: "-12",
	RespDesc:    "api not support",
	RespDescCN:  "该接口尚未被支持",
}

var PIXIE_ERR_LOGIN_PARAM_EMPTY = PixieRespInfo{
	RespCodeStr: "-13",
	RespDesc:    "login param empty",
	RespDescCN:  "登录参数缺失",
}

var PIXIE_ERR_CHANNEL_LIMIT = PixieRespInfo{
	RespCodeStr: "-14",
	RespDesc:    "server limit",
	RespDescCN:  "服务器限制",
}

var PIXIE_ERR_PLAYER_DELETING = PixieRespInfo{
	RespCodeStr: "-15",
	RespDesc:    "player data deleting",
	RespDescCN:  "玩家数据删除中",
}

var PIXIE_ERR_PLAYER_DELETED = PixieRespInfo{
	RespCodeStr: "-16",
	RespDesc:    "player data deleted",
	RespDescCN:  "玩家数据已删除",
}

var PIXIE_ERR_PLAYER_LOGIN_ERR = PixieRespInfo{
	RespCodeStr: "-17",
	RespDesc:    "player login error",
	RespDescCN:  "玩家登录错误",
}

var PIXIE_ERR_DB_TOKEN_WRONG = PixieRespInfo{
	RespCodeStr: "-18",
	RespDesc:    "player login token wrong",
	RespDescCN:  "登录参数异常",
}

var PIXIE_ERR_SAVE_DB_WORNG = PixieRespInfo{
	RespCodeStr: "-19",
	RespDesc:    "player save db error",
	RespDescCN:  "玩家数据存储失败",
}

var PIXIE_ERR_PARAM_ILLEGAL = PixieRespInfo{
	RespCodeStr: "-20",
	RespDesc:    "param illegal",
	RespDescCN:  "参数非法",
}

var PIXIE_ERR_FLUSH_PLAYER = PixieRespInfo{
	RespCodeStr: "-21",
	RespDesc:    "player data flush err",
	RespDescCN:  "玩家数据存储异常",
}
