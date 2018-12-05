package api_specification

//批量获取图纸基础信息接口
const HTTP_PLAYER_BATCH_GET_BASE_PAPER = "api/batchGetPaper"

type BatchGetPaperReq struct {
	//准备弃用 此方法只支持玩家图纸信息的获取
	PaperIDList []int64

	//请客户端迁移使用此参数 可以获取官方或非官方的所有图纸信息
	ClothesIDList []string
}

type BatchGetPaperResp struct {
	Papers []*BasePaper
}

//上传作品接口
const HTTP_PLAYER_ADD_PAPER string = "api/addPaper"

type UploadPaperReq struct {
	//通用参数
	CName    string
	PartType string
	Desc     string
	Icon     string

	//上传服装时
	ClothesType string
	Main        string
	Bottom      string
	Collar      string
	Shadow      string
	ZValue      int

	//上传背景
	BeijingFileType int //0:图片 1:动画
	Front           string
	Back            string
}

type UploadPaperResp struct {
}

//删除图纸
const HTTP_PLAYER_DELETE_PAPER string = "api/deletePaper"

type DeletePaperReq struct {
	PaperID int64
}

type DeletePaperResp struct {
}

//撤销发布图纸
const HTTP_PALYER_CANCEL_PAPER_TRADE = "api/cancelPaperTrade"

type CancelPaperTradeReq struct {
	TradeID int64
	PaperID int64
}

type CancelPaperTradeResp struct {
}

//发布竞拍(设计师发布及正常玩家发布竞拍)
const HTTP_PLAYER_PUBLISH_PAPER = "api/publishPaper"

type PublishPaperReq struct {
	PaperID   int64
	MaxPrice  int
	MinPrice  int
	Duration  int64
	MoneyType int //货币类型

	Sequence int
	//最大流通数量 只在设计师首次发布时可以设置
	CirculationMax int
}

type PublishPaperResp struct {
	TradeID   int64
	StartTime int64
}

//购买图纸
const HTTP_PLAYER_BUY_PAPER = "api/buyPaper"

type BuyPaperReq struct {
	TradeID          int64
	PaperID          int64
	Sequence         int
	AuctionPrice     float64 //竞拍价格
	AuctionPriceType int
}

type BuyPaperResp struct {
	CurrentGold    int64
	CurrentClothes ClothesInfo
}

//竞拍失败消息置为已读
const HTTP_PLAYER_SET_AUCTION_FAIL_READ = "api/setPaperStatusRead"

type MarkAuctionFailUnReadReq struct {
	PaperID int64

	Type int //1:发布图纸 2:普通拍卖
}

type MarkAuctionFailUnReadResp struct {
}
