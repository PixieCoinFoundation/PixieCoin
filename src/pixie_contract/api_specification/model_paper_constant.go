package api_specification

//模特
const (
	MODEL_MALE     = "1"
	MODEL_FEMALE   = "2"
	MODEL_KID      = "3"
	MODEL_NO_LIMIT = "4"
)

//服装标签
const (
	PAPER_TAG_ZHONGGUOFENG = "1"
	PAPER_TAG_HEFENG       = "2"
	PAPER_TAG_OUSHIGUDIAN  = "3"
	PAPER_TAG_MINZUFENG    = "4"
	PAPER_TAG_LUOLITA      = "5"
	PAPER_TAG_ANHEIXI      = "6"
	PAPER_TAG_KEHUAN       = "7"
	PAPER_TAG_GUZHUANG     = "8"
	PAPER_TAG_NVPU         = "9"
	PAPER_TAG_XIAOYUAN     = "10"
	PAPER_TAG_ZHIFUXI      = "11"
	PAPER_TAG_HAITAN       = "12"
	PAPER_TAG_PIGEXI       = "13"
	PAPER_TAG_JUJIA        = "14"
	PAPER_TAG_XIAOQINGXIN  = "15"
	PAPER_TAG_YONGZHUANG   = "16"
	PAPER_TAG_SHUIYI       = "17"
	PAPER_TAG_YAOGUN       = "18"
	PAPER_TAG_HUNLI        = "19"
	PAPER_TAG_SENNV        = "20"
	PAPER_TAG_WANOUFU      = "21"
)

var PAPER_TYPE_ZVALUE_MAP = map[string][]int{
	PAPER_TYPE_FAXING:    []int{325},
	PAPER_TYPE_WAITAO:    []int{320, 328},
	PAPER_TYPE_SHANGYI:   []int{300},
	PAPER_TYPE_XIAZHUANG: []int{295, 280, 305},
	PAPER_TYPE_WAZI:      []int{260},
	PAPER_TYPE_XIEZI:     []int{270, 285},
	PAPER_TYPE_LIANTIFU:  []int{300, 280},
	PAPER_TYPE_MAKEUP:    []int{202},
	PAPER_TYPE_BEIJING:   []int{50, 350},
	PAPER_TYPE_HEAD:      []int{326},
	PAPER_TYPE_FACE:      []int{240, 327},
	PAPER_TYPE_EAR:       []int{323},
	PAPER_TYPE_NECK:      []int{210, 310, 315, 340},
	PAPER_TYPE_HAND:      []int{290, 329},
	PAPER_TYPE_INHAND:    []int{294, 330},
	PAPER_TYPE_PIFU:      []int{201},
	PAPER_TYPE_OTHER:     []int{322, 301, 345, 60},
}

//服装类别和暗标签
const (
	//发型
	PAPER_TYPE_FAXING = "100"
	// PAPER_TYPE_FAXING_STAG_CHANGFA = "100_1"
	// PAPER_TYPE_FAXING_STAG_DUANFA  = "100_2"
	// PAPER_TYPE_FAXING_STAG_FAJI    = "100_3"
	// PAPER_TYPE_FAXING_STAG_BIANZI  = "100_4"

	//外套
	PAPER_TYPE_WAITAO = "200"
	// PAPER_TYPE_WAITAO_STAG_FENGYI   = "200_1"
	// PAPER_TYPE_WAITAO_STAG_JIAKE    = "200_2"
	// PAPER_TYPE_WAITAO_STAG_XIFU     = "200_3"
	// PAPER_TYPE_WAITAO_STAG_CHANGPAO = "200_4"
	// PAPER_TYPE_WAITAO_STAG_MAJIA    = "200_5"

	//上衣
	PAPER_TYPE_SHANGYI = "300"
	// PAPER_TYPE_SHANGYI_STAG_CHENSHAN = "300_1"
	// PAPER_TYPE_SHANGYI_STAG_TXU      = "300_2"
	// PAPER_TYPE_SHANGYI_STAG_BEIXIN   = "300_3"

	//下装
	PAPER_TYPE_XIAZHUANG = "400"
	// PAPER_TYPE_XIAZHUANG_STAG_DUANKU   = "400_1"
	// PAPER_TYPE_XIAZHUANG_STAG_CHANGKU  = "400_2"
	// PAPER_TYPE_XIAZHUANG_STAG_DUANQUN  = "400_3"
	// PAPER_TYPE_XIAZHUANG_STAG_CHANGQUN = "400_4"

	//袜子
	PAPER_TYPE_WAZI = "500"

	//鞋子
	PAPER_TYPE_XIEZI = "600"
	// PAPER_TYPE_XIEZI_STAG_XUEZI     = "600_1"
	// PAPER_TYPE_XIEZI_STAG_TUOXIE    = "600_2"
	// PAPER_TYPE_XIEZI_STAG_CHIJIAO   = "600_3"
	// PAPER_TYPE_XIEZI_STAG_GAOGENXIE = "600_4"

	//连体服
	PAPER_TYPE_LIANTIFU = "700"
	// PAPER_TYPE_LIANTIFU_STAG_LIANTIKU = "700_1"
	// PAPER_TYPE_LIANTIFU_STAG_DUANQUN  = "700_2"
	// PAPER_TYPE_LIANTIFU_STAG_CHANGQUN = "700_3"

	//妆容
	PAPER_TYPE_MAKEUP = "800"

	//背景
	PAPER_TYPE_BEIJING = "900"

	//以下为饰品类别

	//头部
	PAPER_TYPE_HEAD = "1001"
	// PAPER_TYPE_HEAD_STAG_MAOZI  = "1001_1"
	// PAPER_TYPE_HEAD_STAG_FAZAN  = "1001_2"
	// PAPER_TYPE_HEAD_STAG_FAQIA  = "1001_3"
	// PAPER_TYPE_HEAD_STAG_TOUTAO = "1001_4"

	//脸部
	PAPER_TYPE_FACE = "1002"
	// PAPER_TYPE_FACE_STAG_YANJING = "1002_1"
	// PAPER_TYPE_FACE_STAG_MIANJU  = "1002_2"
	// PAPER_TYPE_FACE_STAG_MIANSHA = "1002_3"

	//耳部
	PAPER_TYPE_EAR = "1003"

	//颈部
	PAPER_TYPE_NECK = "1004"
	// PAPER_TYPE_NECK_STAG_XIANGLIAN = "1004_1"
	// PAPER_TYPE_NECK_STAG_LINGDAI   = "1004_2"
	// PAPER_TYPE_NECK_STAG_LINGJIE   = "1004_3"
	// PAPER_TYPE_NECK_STAG_WEIJIN    = "1004_4"

	//手部
	PAPER_TYPE_HAND = "1005"

	//手持
	PAPER_TYPE_INHAND = "1006"
	// PAPER_TYPE_INHAND_STAG_BAO  = "1006_1"
	// PAPER_TYPE_INHAND_STAG_ITEM = "1006_2"

	//皮肤
	PAPER_TYPE_PIFU = "1007"

	//其他
	PAPER_TYPE_OTHER = "1008"
)

var PixiePaperClothesTypeMap = map[string]int{
	MODEL_MALE:     1,
	MODEL_FEMALE:   1,
	MODEL_KID:      1,
	MODEL_NO_LIMIT: 1,
}

var PixiePaperPartTypeMap = map[string]string{
	PAPER_TYPE_FAXING:    "发型",
	PAPER_TYPE_WAITAO:    "外套",
	PAPER_TYPE_SHANGYI:   "上衣",
	PAPER_TYPE_XIAZHUANG: "下装",
	PAPER_TYPE_WAZI:      "袜子",
	PAPER_TYPE_XIEZI:     "鞋子",
	PAPER_TYPE_LIANTIFU:  "连体服",
	PAPER_TYPE_MAKEUP:    "妆容",
	PAPER_TYPE_BEIJING:   "背景",
	PAPER_TYPE_HEAD:      "头部",
	PAPER_TYPE_FACE:      "面部",
	PAPER_TYPE_EAR:       "耳部",
	PAPER_TYPE_NECK:      "颈部",
	PAPER_TYPE_HAND:      "手部",
	PAPER_TYPE_INHAND:    "手持",
	PAPER_TYPE_PIFU:      "皮肤",
	PAPER_TYPE_OTHER:     "其他",
}

type PartTypeInfo struct {
	Part string
	Name string
}

type PartOfZvalue struct {
	Part       string
	PartofList []ZvalueWithName
}

type ZvalueWithName struct {
	Name   string //层级名称
	Zvalue string //层级
}

var ZValueMap = map[int]string{
	50:  "背景",
	60:  "身后",
	200: "身体",
	201: "皮肤",
	202: "妆容",
	210: "小项链",
	240: "眼镜",
	260: "袜子",
	270: "低帮鞋",
	280: "小脚裤裤子",
	285: "高帮鞋",
	290: "手饰",
	294: "男左手、女右手包",
	295: "大腿裤裤子、裙子",
	300: "连衣裙、大裤腿连体裤",
	301: "系着",
	305: "背带裤",
	310: "领带、大项链",
	313: "领子",
	315: "领结",
	320: "外套",
	322: "披着",
	323: "耳饰",
	325: "头发",
	326: "帽子",
	327: "架在头上的眼镜",
	328: "带帽子的外套",
	329: "大手套",
	330: "男右手、女左手包",
	340: "围巾",
	345: "身前",
	350: "前景",
}
