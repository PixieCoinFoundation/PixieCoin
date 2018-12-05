package api_specification

import (
	"fmt"
)

var PAPER_CAN_DELETE_STATUS_IN_SQL = fmt.Sprintf("(%d,%d,%d,%d,%d)", PAPER_STATUS_ADMIN_QUEUE, PAPER_STATUS_QUEUE, PAPER_STATUS_WAIT_PUBLISH, PAPER_STATUS_FAIL, PAPER_STATUS_FREEZE)

//paper状态
type PaperStatus int

const (
	//名词定义
	//paper:玩家上传的作品
	//paper author:paper的创作者
	//paper owner:paper当前的所有者
	//注意：一旦paper有过一次owner，paper的owner字段之后就必然一直不为空

	//UI场景视角：比如在我的设计UI场景下，是【设计师视角】。在我的图纸UI场景下，是【所有者视角】

	//服务端状态：【待管理员审核】
	PAPER_STATUS_ADMIN_QUEUE PaperStatus = -3

	//服务端状态：【冻结】（玩家举报判定生效后）
	PAPER_STATUS_FREEZE PaperStatus = -2

	//服务端状态：【已删除】（玩家主动删除，客户端不会遇到）
	PAPER_STATUS_DELETED PaperStatus = -1

	//服务端状态：【待玩家审核】
	//客户端文案只有一种
	PAPER_STATUS_QUEUE PaperStatus = 0

	//服务端状态：【待发布】
	//客户端文案只有一种
	PAPER_STATUS_WAIT_PUBLISH PaperStatus = 1

	//服务端状态：【未通过】（管理员初审未通过）
	//客户端文案只有一种
	PAPER_STATUS_FAIL PaperStatus = 2

	//服务端状态：【竞拍中】
	//该状态在客户端有两套文案
	//显示为"发布中"的条件:查询者是author 且 paper当前无owner 且 UI场景是设计师视角
	//显示为"竞拍中"的条件:查询者不是author
	PAPER_STATUS_ON_SALE PaperStatus = 3

	//服务端状态：【持有中】
	//该状态在客户端有两套文案
	//显示为"已发布"的条件:查询者是author 且 paper有owner 且 UI场景是设计师视角
	//显示为"正常"的条件:查询者不是author
	PAPER_STATUS_OCCUPY PaperStatus = 4

	//备注："竞拍失败提醒","生产中"这一类服务端定义为：独立于paper状态之外的"属性"
	//当(paper status为PAPER_STATUS_ON_SALE或PAPER_STATUS_OCCUPY)且is_in_production属性为1时 额外显示"生产中"
	//当(paper status为PAPER_STATUS_OCCUPY)且auction_fail_unread属性为1时 额外显示感叹号(代表"竞拍失败提醒")
)

type LandStatus int

var LAND_CAN_LEVEL_UP_STATUS = fmt.Sprintf("(%d,%d)", LAND_STATUS_WB_NORMAL, LAND_STATUS_WB_IN_BUSINESS)
var LAND_CHECK_RENT_STATUS = fmt.Sprintf("(%d,%d,%d)", LAND_STATUS_WB_NORMAL, LAND_STATUS_WB_RENTED_BUSINESS, LAND_STATUS_WB_IN_BUSINESS)
var LAND_NOT_RENTING_STATUS = fmt.Sprintf("(%d,%d,%d,%d,%d,%d)", LAND_STATUS_WB_IN_BUSINESS, LAND_STATUS_WB_RENTED_BUSINESS, LAND_STATUS_WB_NORMAL, LAND_STATUS_NB_RENTED_BUILD, LAND_STATUS_NB_BUILDING, LAND_STATUS_NB_EMPTY)
var LAND_WB_WAIT_BUSINESS_STATUS = fmt.Sprintf("(%d,%d)", LAND_STATUS_WB_NORMAL, LAND_STATUS_WB_RENTED_BUSINESS)
var LAND_NB_WAIT_BUILD_STATUS = fmt.Sprintf("(%d,%d)", LAND_STATUS_NB_EMPTY, LAND_STATUS_NB_RENTED_BUILD)

const (
	//地块无建筑物时的状态
	LAND_STATUS_NB_EMPTY             LandStatus = 0
	LAND_STATUS_NB_RENTING_FOR_BUILD LandStatus = 1
	LAND_STATUS_NB_RENTED_BUILD      LandStatus = 2
	LAND_STATUS_NB_BUILDING          LandStatus = 3

	//地块有建筑物时的状态
	LAND_STATUS_WB_NORMAL               LandStatus = 10
	LAND_STATUS_WB_RENTING_FOR_BUSINESS LandStatus = 11
	LAND_STATUS_WB_RENTED_BUSINESS      LandStatus = 12
	LAND_STATUS_WB_IN_BUSINESS          LandStatus = 13

	// LAND_STATUS_WB_LEVEL_UP_ING         LandStatus = 14 //地块建筑升级中
)

type RecordRank string

const (
	RECORD_RANK_S = "S"
	RECORD_RANK_A = "A"
	RECORD_RANK_B = "B"
)

type ClothesOrigin string

const (
	CLOTHES_ORIGIN_OFFICIAL_PAPER ClothesOrigin = "O"
	CLOTHES_ORIGIN_PLAYER_PAPER   ClothesOrigin = "P"
)

type LandType int

const (
	PIXIE_LAND_TYPE_DEFAULT LandType = 0
	PIXIE_LAND_TYPE_BIG     LandType = 1
	PIXIE_LAND_TYPE_MY_SHOP LandType = 2
)

//0.待结算,1.已经结算,2.已删除
type PaperVerifyStatus int

const (
	PIXIE_PAPER_VERIFY_READY     PaperVerifyStatus = 0
	PIXIE_PAPER_VERIFY_PRE_TRIAL PaperVerifyStatus = 2 //预审状态
	PIXIE_PAPER_VERIFY_FINISH    PaperVerifyStatus = 1
)

//1.是举报奖励，2.是评审奖励 3.无奖励
type PaperNotifyType int

const (
	PIXIE_PAPER_NOTIFY_COPY_REWARD   PaperNotifyType = 1
	PIXIE_PAPER_NOTIFY_VERIFY_REWARD PaperNotifyType = 2
	PIXIE_PAPER_NOTIFY_NO_REWARD     PaperNotifyType = 3
)
