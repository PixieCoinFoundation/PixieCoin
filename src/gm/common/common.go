package common

import (
	"constants"
	. "pixie_contract/api_specification"
)

var ClothesTypeStringMap = map[string]string{
	MODEL_MALE:   "男装",
	MODEL_FEMALE: "女装",
	MODEL_KID:    "童装",
}

var PartTypeStringMap = map[string]string{
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
	PAPER_TYPE_FACE:      "脸部",
	PAPER_TYPE_EAR:       "耳部",
	PAPER_TYPE_NECK:      "颈部",
	PAPER_TYPE_HAND:      "手部",
	PAPER_TYPE_INHAND:    "手持",
	PAPER_TYPE_PIFU:      "皮肤",
	PAPER_TYPE_OTHER:     "其他",
}

var CurrencyStringMap = map[int]string{
	constants.PIXIE_GOLD_TYPE: "gold",
	constants.PIXIE_PXC_TYPE:  "pxc",
}

var MenuList = `[{
    "title": "玩家信息",
    "icon": "fa-cubes",
    "id":100,
    "spread": true,
    "children": [{
    	"id":101,
    	"parent_id":100,
        "title": "基本信息",
        "icon": "&#xe641;",
        "href": "player.html"
    }]
},{
	"id":200,
    "title": "管理员管理",
    "icon": "fa-stop-circle",
    "href": "#",
    "spread": false,
    "children": [{
    	"id":201,
    	"parent_id":200,
        "title": "管理员列表",
        "icon": "fa-github",
        "href": "adminList.html"
    }, {
    	"id":202,
    	"parent_id":200,
        "title": "角色管理",
        "icon": "fa-github",
        "href": "roleManage.html"
    }]
},{
    "id":300,
    "title": "推荐服装内容",
    "icon": "fa-stop-circle",
    "href": "#",
    "spread": false,
    "children": [{
        "id":301,
        "parent_id":300,
        "title": "期刊列表",
        "icon": "fa-github",
        "href": "recommend/subjectList.html"
    }, {
        "id":302,
        "parent_id":300,
        "title": "更新套装",
        "icon": "fa-github",
        "href": "recommend/subjectSuit.html"
    },{
        "id":303,
        "parent_id":300,
        "title": "所有推荐专题",
        "icon": "fa-github",
        "href": "recommend/subjectTopicList.html"
    },{
        "id":304,
        "parent_id":300,
        "title": "h5单品列表",
        "icon": "fa-github",
        "href": "recommend/subjectItemList.html"
    },{
        "id":305,
        "parent_id":300,
        "title": "正在售卖中的服装列表",
        "icon": "fa-github",
        "href": "recommend/saleClothes.html"
    }, {
        "id":306,
        "parent_id":300,
        "title": "所有单品",
        "icon": "fa-github",
        "href": "recommend/subjectItemsList.html"
    },{
        "id":307,
        "parent_id":300,
        "title": "单品推荐预选池",
        "icon": "fa-github",
        "href": "recommend/poolClothes.html"
    },{
        "id":308,
        "parent_id":300,
        "title":"单品推荐池",
        "icon":"fa-github",
        "href":"recommend/poolClothes1.html"
    }]
}, {
    "id":400,
    "title": "官方商店上传图纸",
    "icon": "fa-address-book",
    "href": "",
    "spread": false,
    "children": [{
        "id":401,
        "parent_id":400,
        "title": "服装列表",
        "icon": "fa-github",
        "href": "clothesList.html"
    }, {
        "id":402,
        "parent_id":400,
        "title": "背景列表",
        "icon": "fa-qq",
        "href": "bgList.html"
    }, {
        "id":403,
        "parent_id":400,
        "title": "上传服装",
        "icon": "&#xe609;",
        "href": "uploadClothes.html"
    }, {
        "id":403,
        "parent_id":400,
        "title": "上传背景",
        "icon": "&#xe609;",
        "href": "uploadBackground.html"
    }]
}, {
	"id":500,
    "title": "官方推荐图纸",
    "icon": "fa-address-book",
    "href": "",
    "spread": false,
    "children": []
},{
    "id":600,
    "title": "玩家上传图纸",
    "icon": "fa-stop-circle",
    "href": "",
    "spread": false,
    "children": [{
        "id":601,
        "parent_id":600,
        "title": "待预审服装列表",
        "icon": "fa-github",
        "href": "verifyClothes.html"
    }, {
        "id":602,
        "parent_id":600,
        "title": "待预审背景列表",
        "icon": "fa-qq",
        "href": "verfiyBgList.html"
    }, {
        "id":603,
        "parent_id":600,
        "title": "通过预审服装列表",
        "icon": "fa-qq",
        "href": "clothes/auctionList.html"
    }, {
        "id":604,
        "parent_id":600,
        "title": "通过预审背景列表",
        "icon": "fa-qq",
        "href": "clothes/auctionBgList.html"
    }, {
        "id":605,
        "parent_id":600,
        "title": "拒绝预审服装列表",
        "icon": "fa-qq",
        "href": "clothes/rejectList.html"
    }, {
        "id":606,
        "parent_id":600,
        "title": "拒绝预审背景列表",
        "icon": "fa-qq",
        "href": "clothes/rejectBgList.html"
    }]
}, {
    "id":700,
    "title": "大众评审时图纸举报",
    "icon": "fa-stop-circle",
    "href": "",
    "spread": false,
    "children": [{
        "id":701,
        "parent_id":700,
        "title": "被举报作品列表",
        "icon": "fa-github",
        "href": "verifyClothes.html"
    }, {
        "id":702,
        "parent_id":700,
        "title": "被下架作品列表",
        "icon": "fa-qq",
        "href": "clothes/auctionList.html"
    }, {
        "id":703,
        "parent_id":700,
        "title": "待处理",
        "icon": "fa-qq",
        "href": "clothes/rejectList.html"
    }, {
        "id":704,
        "parent_id":700,
        "title": "申诉成功",
        "icon": "fa-qq",
        "href": "clothes/rejectList.html"
    }, {
        "id":705,
        "parent_id":700,
        "title": "申诉失败",
        "icon": "fa-qq",
        "href": "clothes/rejectList.html"
    }]
},{
    "id":800,
    "title": "发布后图纸举报",
    "icon": "fa-stop-circle",
    "href": "",
    "spread": false,
    "children": [{
        "id":801,
        "parent_id":800,
        "title": "被举报作品列表",
        "icon": "fa-github",
        "href": "verfiyBgList.html"
    }, {
        "id":802,
        "parent_id":800,
        "title": "被下架作品列表",
        "icon": "fa-qq",
        "href": "clothes/auctionBgList.html"
    }, {    
        "id":803,
        "parent_id":800,
        "title": "待处理",
        "icon": "fa-qq",
        "href": "clothes/rejectBgList.html"
    }, {    
        "id":804,
        "parent_id":800,
        "title": "申诉成功",
        "icon": "fa-qq",
        "href": "clothes/rejectBgList.html"
    }, {    
        "id":805,
        "parent_id":800,
        "title": "申诉失败",
        "icon": "fa-qq",
        "href": "clothes/rejectBgList.html"
    }]
},{
    "id":900,
    "title": "服务器管理",
    "icon": "fa-stop-circle",
    "href": "#",
    "spread": false,
    "children": [{
        "id":901,
        "parent_id":900,
        "title": "维护",
        "icon": "fa-github",
        "href": "maintenance.html"
    },{
        "id":902,
        "parent_id":900,
        "title": "运行管理",
        "icon": "fa-github",
        "href": "runtime.html"
    },{
        "id":903,
        "parent_id":900,
        "title":"上传热更",
        "icon":"fa-github",
        "href":"uploadPack.html"
    }]
}]`
