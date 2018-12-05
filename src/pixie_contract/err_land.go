package pixie_contract

var PIXIE_ERR_LAND_RENT_PRICE_WRONG = PixieRespInfo{
	RespCodeStr: "2001",
	RespDesc:    "price illegal",
	RespDescCN:  "价格非法",
}

var PIXIE_ERR_LAND_RENT = PixieRespInfo{
	RespCodeStr: "2002",
	RespDesc:    "land rent error",
	RespDescCN:  "招租失败",
}

var PIXIE_ERR_LANG_BUILDING = PixieRespInfo{
	RespCodeStr: "2003",
	RespDesc:    "land building err",
	RespDescCN:  "装修失败",
}

var PIXIE_ERR_LAND_STATUS_WRONG = PixieRespInfo{
	RespCodeStr: "2004",
	RespDesc:    "land status wrong",
	RespDescCN:  "地块状态错误",
}

var PIXIE_ERR_LAND_START_BUSINESS = PixieRespInfo{
	RespCodeStr: "2005",
	RespDesc:    "land start business err",
	RespDescCN:  "营业失败",
}

var PIXIE_ERR_REMOVE_BUILDING = PixieRespInfo{
	RespCodeStr: "2006",
	RespDesc:    "remove building err",
	RespDescCN:  "拆除失败",
}

var PIXIE_ERR_LIST_LAND = PixieRespInfo{
	RespCodeStr: "2007",
	RespDesc:    "list land err",
	RespDescCN:  "查询地块列表失败",
}

var PIXIE_ERR_LIST_LAND_SALE = PixieRespInfo{
	RespCodeStr: "2008",
	RespDesc:    "list land sale err",
	RespDescCN:  "查询地块销售列表失败",
}

var PIXIE_ERR_START_BUSINESS = PixieRespInfo{
	RespCodeStr: "2009",
	RespDesc:    "land start business err",
	RespDescCN:  "营业失败",
}

var PIXIE_ERR_STOP_BUSINESS = PixieRespInfo{
	RespCodeStr: "2010",
	RespDesc:    "land stop business err",
	RespDescCN:  "停止营业失败",
}

var PIXIE_ERR_LOCK_LAND = PixieRespInfo{
	RespCodeStr: "2011",
	RespDesc:    "lock land err",
	RespDescCN:  "锁定地块失败",
}

var PIXIE_ERR_GET_ONE_LAND = PixieRespInfo{
	RespCodeStr: "2012",
	RespDesc:    "get one land err",
	RespDescCN:  "查询地块失败",
}

var PIXIE_ERR_LAND_NOT_IN_BUSINESS = PixieRespInfo{
	RespCodeStr: "2013",
	RespDesc:    "land not in business",
	RespDescCN:  "该地块不在营业中哦~",
}

var PIXIE_ERR_LAND_RENT_END = PixieRespInfo{
	RespCodeStr: "2014",
	RespDesc:    "land rent end",
	RespDescCN:  "地块租期结束",
}

var PIXIE_ERR_BUILDING_ID_WRONG = PixieRespInfo{
	RespCodeStr: "2015",
	RespDesc:    "building id not legal",
	RespDescCN:  "建筑ID非法",
}

var PIXIE_ERR_PRICE_ILLEGAL = PixieRespInfo{
	RespCodeStr: "2016",
	RespDesc:    "price illegal",
	RespDescCN:  "价格非法",
}

var PIXIE_ERR_LAND_RENT_FOR_BUILD = PixieRespInfo{
	RespCodeStr: "2017",
	RespDesc:    "update lend rent for build err",
	RespDescCN:  "招租地块出错",
}

var PIXIE_ERR_LAND_GET_ERR = PixieRespInfo{
	RespCodeStr: "2018",
	RespDesc:    "land get err",
	RespDescCN:  "获取信地块息错误",
}

var PIXIE_ERR_LAND_OWNER_ILLEGAL = PixieRespInfo{
	RespCodeStr: "2019",
	RespDesc:    "land owner illegal",
	RespDescCN:  "非地块主人，禁止操作",
}

var PIXIE_ERR_RENT_LAND_ERR = PixieRespInfo{
	RespCodeStr: "2020",
	RespDesc:    "rent land err",
	RespDescCN:  "租地发生错误",
}

var PIXIE_ERR_RENT_LAND_PRICE_ILLEGAL = PixieRespInfo{
	RespCodeStr: "2021",
	RespDesc:    "rent price illegal",
	RespDescCN:  "竞租地块价格非法",
}

var PIXIE_ERR_LOCKLAND_FAIL = PixieRespInfo{
	RespCodeStr: "2022",
	RespDesc:    "lock land err",
	RespDescCN:  "锁地块错误",
}

var PIXIE_ERR_GET_RANDOM_MAP = PixieRespInfo{
	RespCodeStr: "2023",
	RespDesc:    "get random player error",
	RespDescCN:  "随机访问失败",
}

var PIXIE_ERR_RENT_LAND_FOR_BUSINESS = PixieRespInfo{
	RespCodeStr: "2024",
	RespDesc:    "rent land for business fail",
	RespDescCN:  "发布招商失败",
}

var PIXIE_ERR_RENT_LAND_FOR_BUILD = PixieRespInfo{
	RespCodeStr: "2025",
	RespDesc:    "rent land for build",
	RespDescCN:  "招商失败",
}

var PIXIE_ERR_RENT_SELF = PixieRespInfo{
	RespCodeStr: "2026",
	RespDesc:    "can not rent self",
	RespDescCN:  "不能竞拍自己的地块或建筑",
}

var PIXIE_ERR_RENTING = PixieRespInfo{
	RespCodeStr: "2027",
	RespDesc:    "land or building in renting",
	RespDescCN:  "地块或建筑正在租期中",
}

var PIXIE_ERR_RENT_PRICE_TYPE_WRONG = PixieRespInfo{
	RespCodeStr: "2028",
	RespDesc:    "rent price type wrong",
	RespDescCN:  "竞拍价格类型错误",
}

var PIXIE_ERR_GET_STATUS_LAND_WRONG = PixieRespInfo{
	RespCodeStr: "2029",
	RespDesc:    "get status land wrong",
	RespDescCN:  "获取指定状态地块错误",
}

var PIXIE_ERR_GET_MY_LAND_WRONG = PixieRespInfo{
	RespCodeStr: "2030",
	RespDesc:    "get my land wrong",
	RespDescCN:  "获取我的地块错误",
}

var PIXIE_ERR_BUILD_USER_NOT_MATCH = PixieRespInfo{
	RespCodeStr: "2031",
	RespDesc:    "build username not match",
	RespDescCN:  "该建筑并非由您建造，无法拆除",
}

var PIXIE_ERR_NOT_OCCUPY_LAND = PixieRespInfo{
	RespCodeStr: "2032",
	RespDesc:    "not occupy land",
	RespDescCN:  "您没有对该地块的操作权限",
}

var PIXIE_ERR_LAND_BUILDING_NOT_MATCH = PixieRespInfo{
	RespCodeStr: "LAND_BUILDING_NOT_MATCH",
	RespDesc:    "land type and building type not match",
	RespDescCN:  "该地块不支持此类型的建筑",
}

var PIXIE_ERR_LAND_BUILDING_ALREADY_EXIST = PixieRespInfo{
	RespCodeStr: "LAND_BUILDING_ALREADY_EXIST",
	RespDesc:    "building of this type already exist",
	RespDescCN:  "该类型的建筑无法重复建造",
}

var PIXIE_ERR_LAND_BUILDING_NOT_EXIST = PixieRespInfo{
	RespCodeStr: "LAND_BUILDING_NOT_EXIST",
	RespDesc:    "building of this type not exist",
	RespDescCN:  "该类型的建筑不存在",
}

var PIXIE_ERR_LANG_BUILDING_LEVEL_UP = PixieRespInfo{
	RespCodeStr: "ERR_LANG_BUILDING_LEVEL_UP",
	RespDesc:    "building level up error",
	RespDescCN:  "建筑升级出错",
}

var PIXIE_ERR_LAND_BUILDING_ILLEGAL = PixieRespInfo{
	RespCodeStr: "LAND_BUILDING_ILLEGAL",
	RespDesc:    "building illegal",
	RespDescCN:  "该类型的建筑非法",
}

var PIXIE_ERR_BUILDING_NOT_DEMAND = PixieRespInfo{
	RespCodeStr: "BUILDING_NOT_DEMAND",
	RespDesc:    "building not demand",
	RespDescCN:  "该建筑类型错误",
}

var PIXIE_ERR_START_DEMAND_BUSINESS = PixieRespInfo{
	RespCodeStr: "ERR_START_DEMAND_BUSINESS",
	RespDesc:    "start demand business error",
	RespDescCN:  "需求类建筑开店失败",
}

var PIXIE_ERR_STOP_DEMAND_BUSINESS = PixieRespInfo{
	RespCodeStr: "ERR_STOP_DEMAND_BUSINESS",
	RespDesc:    "stop demand business error",
	RespDescCN:  "需求类建筑闭店失败",
}

var PIXIE_ERR_DEMAND_BUILD_NOT_SELF = PixieRespInfo{
	RespCodeStr: "DEMAND_BUILD_NOT_SELF",
	RespDesc:    "demand building can only be built on your land",
	RespDescCN:  "需求类建筑只能在自己的土地上进行建设",
}

var PIXIE_ERR_SET_SHOP_MODEL = PixieRespInfo{
	RespCodeStr: "ERR_SET_SHOP_MODEL",
	RespDesc:    "set shop model error",
	RespDescCN:  "设置店铺模特失败",
}

var PIXIE_ERR_CANCEL_LEVEL_UP_BUILDING = PixieRespInfo{
	RespCodeStr: "ERR_CANCEL_LEVEL_UP_BUILDING",
	RespDesc:    "cancel building level up error",
	RespDescCN:  "取消升级失败",
}

var PIXIE_ERR_LAND_SALE_FORMAT = PixieRespInfo{
	RespCodeStr: "ERR_LAND_SALE_FORMAT",
	RespDesc:    "buy paper clothes fail.land sale format.",
	RespDescCN:  "购买服装失败。土地销售格式有误。",
}

var PIXIE_ERR_CLOTHES_NOT_SALE = PixieRespInfo{
	RespCodeStr: "CLOTHES_NOT_SALE",
	RespDesc:    "paper clothes not for sale",
	RespDescCN:  "未上架的服装",
}
