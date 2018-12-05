package pixie_contract

var PIXIE_ERR_UPLOAD_PAPER = PixieRespInfo{
	RespCodeStr: "1004",
	RespDesc:    "upload paper failed",
	RespDescCN:  "上传图纸失败",
}

var PIXIE_ERR_GET_PAPER_BY_AUTHORNAME = PixieRespInfo{
	RespCodeStr: "1005",
	RespDesc:    "get paper by author wrong",
	RespDescCN:  "获取设计失败",
}

var PIXIE_ERR_GET_PAPER_BY_ID = PixieRespInfo{
	RespCodeStr: "1006",
	RespDesc:    "get paper by id wrong",
	RespDescCN:  "获取图纸失败",
}

var PIXIE_ERR_GET_ALL_PAPER = PixieRespInfo{
	RespCodeStr: "1007",
	RespDesc:    "get all paper wrong",
	RespDescCN:  "查询图纸失败",
}

var PIXIE_ERR_GET_BY_STATUS = PixieRespInfo{
	RespCodeStr: "1007",
	RespDesc:    "get paper by status wrong",
	RespDescCN:  "过滤图纸失败",
}

var PIXIE_ERR_GET_TYPE_ERR = PixieRespInfo{
	RespCodeStr: "1008",
	RespDesc:    "get paper type wrong",
	RespDescCN:  "查询参数错误",
}

var PIXIE_ERR_DELETE_PAPER_ERR = PixieRespInfo{
	RespCodeStr: "1009",
	RespDesc:    "delete paper db failed",
	RespDescCN:  "删除图纸失败",
}

var PIXIE_ERR_CANCEL_PAPER_TRADE_ERR = PixieRespInfo{
	RespCodeStr: "1010",
	RespDesc:    "cancel paper trade db failed",
	RespDescCN:  "撤销作品失败",
}

var PIXIE_ERR_PAPER_VERIFY_INSERT = PixieRespInfo{
	RespCodeStr: "1011",
	RespDesc:    "paper verify insert failed",
	RespDescCN:  "审核失败",
}

var PIXIE_ERR_PAPER_VERIFY_QUERY_DETAIL = PixieRespInfo{
	RespCodeStr: "1012",
	RespDesc:    "paper verify query detail type wrong",
	RespDescCN:  "查询审核详情失败",
}

var PIXIE_ERR_LOCKPAPER_FAIL = PixieRespInfo{
	RespCodeStr: "1013",
	RespDesc:    "paper lock file",
	RespDescCN:  "锁定作品失败",
}

var PIXIE_ERR_TRADE_FAIL = PixieRespInfo{
	RespCodeStr: "1014",
	RespDesc:    "paper trade fail",
	RespDescCN:  "作品交易失败",
}

var PIXIE_ERR_BUY_PAPER = PixieRespInfo{
	RespCodeStr: "1015",
	RespDesc:    "buy paper fail",
	RespDescCN:  "购买图纸失败",
}

var PIXIE_ERR_BUY_PAPER_TIME_OUT = PixieRespInfo{
	RespCodeStr: "1016",
	RespDesc:    "buy paper auction time out",
	RespDescCN:  "图纸竞拍超时",
}

var PIXIE_ERR_PRODUCT_PAPER_FAIL = PixieRespInfo{
	RespCodeStr: "1017",
	RespDesc:    "product paper fail",
	RespDescCN:  "生产失败",
}

var PIXIE_ERR_CANCEL_PRODUCT_PAPER_FAIL = PixieRespInfo{
	RespCodeStr: "1018",
	RespDesc:    "cancel product paper fail",
	RespDescCN:  "取消生产失败",
}

var PIXIE_ERR_AUCTION_STATUS_ERR = PixieRespInfo{
	RespCodeStr: "1019",
	RespDesc:    "auction status change err",
	RespDescCN:  "竞拍状态错误",
}

var PIXIE_ERR_LIST_DESIGN_PAPER = PixieRespInfo{
	RespCodeStr: "1020",
	RespDesc:    "list design paper error",
	RespDescCN:  "获取设计列表失败",
}

var PIXIE_ERR_LIST_OWN_PAPER = PixieRespInfo{
	RespCodeStr: "1021",
	RespDesc:    "list own paper error",
	RespDescCN:  "获取已拥有图纸列表失败",
}

var PIXIE_ERR_LIST_SALE_PAPER = PixieRespInfo{
	RespCodeStr: "1022",
	RespDesc:    "list sale paper error",
	RespDescCN:  "获取在售图纸失败",
}

var PIXIE_ERR_TRADE_NOT_EXIST = PixieRespInfo{
	RespCodeStr: "1023",
	RespDesc:    "trade not exist",
	RespDescCN:  "竞拍不存在",
}

var PIXIE_ERR_QUERY_PAPER_TRADE_HISTORY = PixieRespInfo{
	RespCodeStr: "1024",
	RespDesc:    "query paper trade history",
	RespDescCN:  "查询作品交易历史失败",
}

var PIXIE_ERR_PAPER_PRODUCT_LOCK_FAIL = PixieRespInfo{
	RespCodeStr: "1025",
	RespDesc:    "lock paper product fail",
	RespDescCN:  "锁定作品生产失败",
}

var PIXIE_ERR_CANCEL_PAPER_PRODUCT_LOCK_FAIL = PixieRespInfo{
	RespCodeStr: "1026",
	RespDesc:    "lock paper cancel product fail",
	RespDescCN:  "锁定作品生产失败",
}

var PIXIE_ERR_QUERY_PAPER_SALE = PixieRespInfo{
	RespCodeStr: "1028",
	RespDesc:    "error when query paper sale",
	RespDescCN:  "查询作品销售失败",
}

var PIXIE_ERR_BUY_CLOTHES = PixieRespInfo{
	RespCodeStr: "1029",
	RespDesc:    "buy clothes err",
	RespDescCN:  "购买衣服失败",
}

var PIXIE_ERR_BUY_SELF = PixieRespInfo{
	RespCodeStr: "1030",
	RespDesc:    "buy self's clothes err",
	RespDescCN:  "不能购买自己的衣服",
}

var PIXIE_ERR_PAPER_CLOTHES_PRICING_TYPE_WRONG = PixieRespInfo{
	RespCodeStr: "1031",
	RespDesc:    "clothes pricing type err",
	RespDescCN:  "服装定价类型参数错误",
}

var PIXIE_ERR_PAPER_CLOTHES_PRICING_SET = PixieRespInfo{
	RespCodeStr: "1032",
	RespDesc:    "clothes pricing add err",
	RespDescCN:  "服装定价设置错误",
}

var PIXIE_ERR_PAPER_CLOTHES_PRICING_CHANGE = PixieRespInfo{
	RespCodeStr: "1033",
	RespDesc:    "clothes pricing change err",
	RespDescCN:  "服装定价更改价格错误",
}

var PIXIE_ERR_PAPER_CLOTHES_PRICING_PRICE_IILEGAL = PixieRespInfo{
	RespCodeStr: "1034",
	RespDesc:    "clothes pricing price illegal",
	RespDescCN:  "服装定价价格有误",
}

var PIXIE_ERR_PAPER_VERIFY_QUERY = PixieRespInfo{
	RespCodeStr: "1035",
	RespDesc:    "paper verify query failed",
	RespDescCN:  "查询审核失败",
}

var PIXIE_ERR_PAPER_VERIFY_LIMIT = PixieRespInfo{
	RespCodeStr: "1036",
	RespDesc:    "paper verify limit num",
	RespDescCN:  "评审次数用尽",
}

var PIXIE_ERR_RANDOM_PAPER_NO_DATA = PixieRespInfo{
	RespCodeStr: "1037",
	RespDesc:    "no paper for verify",
	RespDescCN:  "当前没有可以评审的图纸哦,稍后再来吧~",
}

var PIXIE_ERR_PAPER_VERIFY_GET_WRONG = PixieRespInfo{
	RespCodeStr: "1038",
	RespDesc:    "paper verify get wrong",
	RespDescCN:  "获取指定图纸出错",
}

var PIXIE_ERR_PAPER_VERIFY_COUNT_LIMIT = PixieRespInfo{
	RespCodeStr: "1039",
	RespDesc:    "paper verify count limit",
	RespDescCN:  "评审图纸达到上限",
}

var PIXIE_ERR_REPORT_WRONG = PixieRespInfo{
	RespCodeStr: "1040",
	RespDesc:    "paper copy report wrong",
	RespDescCN:  "图纸举报出错",
}

var PIXIE_ERR_COPY_QUERY_WRONG = PixieRespInfo{
	RespCodeStr: "1041",
	RespDesc:    "paper copy query wrong",
	RespDescCN:  "图纸举报查询出错",
}

var PIXIE_ERR_RANDOM_PAPER_FOR_VERIFY = PixieRespInfo{
	RespCodeStr: "1042",
	RespDesc:    "get random paper for verify wrong",
	RespDescCN:  "获取审核作品出错",
}

var PIXIE_PAPER_COPY_SUPPORT_WRONG = PixieRespInfo{
	RespCodeStr: "1043",
	RespDesc:    "paper copy support wrong",
	RespDescCN:  "支持抄袭举报失败",
}

var PIXIE_ERR_PAPER_GET_VERIFIED_DETAIL_WRONG = PixieRespInfo{
	RespCodeStr: "1044",
	RespDesc:    "paper get verified detail wrong",
	RespDescCN:  "获取审核详情失败",
}

var PIXIE_ERR_GET_PAPER_VERIFY_REWARD_WRONG = PixieRespInfo{
	RespCodeStr: "1045",
	RespDesc:    "get paper verified reward wrong",
	RespDescCN:  "获取审核奖励失败",
}

var PIXIE_ERR_UPDATE_PAPER_NOTIFY_STATUS = PixieRespInfo{
	RespCodeStr: "1046",
	RespDesc:    "update paper notify status wrong",
	RespDescCN:  "更新通知状态失败",
}

var PIXIE_ERR_UPLOAD_PAPER_FORMAT_WRONG = PixieRespInfo{
	RespCodeStr: "1047",
	RespDesc:    "upload paper format wrong",
	RespDescCN:  "上传图纸文件格式错误",
}

var PIXIE_ERR_BATCH_GET_PAPER = PixieRespInfo{
	RespCodeStr: "1048",
	RespDesc:    "batch get paper failed",
	RespDescCN:  "批量获取图纸信息失败",
}

var PIXIE_ERR_LIST_OFFICIAL_CLOTHES = PixieRespInfo{
	RespCodeStr: "LIST_OFFICIAL_CLOTHES_ERROR",
	RespDesc:    "get official clothes list error",
	RespDescCN:  "获取官方服装信息出错",
}

var PIXIE_ERR_OFFICIAL_CLOTHES_PRICE_TYPE_WRONG = PixieRespInfo{
	RespCodeStr: "OFFICIAL_CLOTHES_PRICE_TYPE_WRONG",
	RespDesc:    "official clothes price type wrong",
	RespDescCN:  "官方服装价格类型错误",
}

var PIXIE_ERR_OFFICIAL_CLOTHES_NOT_EXIST = PixieRespInfo{
	RespCodeStr: "OFFICIAL_CLOTHES_NOT_EXIST",
	RespDesc:    "this official clothes not exist",
	RespDescCN:  "此官方服装不存在",
}

var PIXIE_ERR_LIST_SPECIFIED_OWN_PAPER = PixieRespInfo{
	RespCodeStr: "ERR_LIST_SPECIFIED_OWN_PAPER",
	RespDesc:    "query specified own paper error",
	RespDescCN:  "查询指定服装列表失败",
}

var PIXIE_ERR_GET_ADMIN_CHECK_CNT = PixieRespInfo{
	RespCodeStr: "ERR_GET_ADMIN_CHECK_CNT",
	RespDesc:    "query admin check size failed",
	RespDescCN:  "获取服装排队情况失败",
}

var PIXIE_ERR_GET_SALE_CLOTHES = PixieRespInfo{
	RespCodeStr: "GET_SALE_CLOTHES_WRONG",
	RespDesc:    "get sale clothes wrong",
	RespDescCN:  "获取在售服装",
}
