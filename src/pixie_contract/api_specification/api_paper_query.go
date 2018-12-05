package api_specification

//获取某一件设计的详情
const HTTP_PLAYER_GET_ONE_DESIGN_PAPER = "api/getOneDesignPaper"

type GetOneDesignPaperReq struct {
	PaperID int64
}

type GetOneDesignPaperResp struct {
	Paper DesignPaper

	//当paper status为PAPER_STATUS_ADMIN_QUEUE 时 才读此值
	AdminCheckBefore int
}

//获取我的设计列表
const HTTP_PLAYER_LIST_DESIGN_PAPER = "api/listDesignPaper"

type ListDesignPaperReq struct {
	Page     int
	PageSize int
	//若空则代表查自己的设计列表
	//若不为空则代表查询别人的设计列表
	TargetUsername string
}

type ListDesignPaperResp struct {
	List []DesignPaper
}

//获取我的图纸列表
const HTTP_PLAYER_LIST_OWN_PAPER = "api/listOwnPaper"

type ListOwnPaperReq struct {
	Page     int
	PageSize int

	//可选参数-模特类型
	ClothesType string

	//可选参数-部位类型
	PartType string

	//若空则代表查自己的设计列表
	//若不为空则代表查询别人的设计列表
	TargetUsername string
}

type ListOwnPaperResp struct {
	List []*OwnPaper
}

//获取指定拥有图纸列表 只能查非官方图纸信息
const HTTP_PLAYER_LIST_SPECIFIED_OWN_PAPER = "api/listSpecifiedOwnPaper"

type ListSpecifiedOwnPaperReq struct {
	ClothesIDList []string

	//若空则代表查自己的设计列表
	//若不为空则代表查询别人的设计列表
	TargetUsername string
}

type ListSpecifiedOwnPaperResp struct {
	List []*OwnPaper
}

//获取在售图纸列表
const HTTP_PLAYER_LIST_SALE_PAPER = "api/listSalePaper"

type ListSalePaperReq struct {
	//从1开始
	Page     int
	PageSize int

	//0：默认是发布时间
	//1：竞拍时间
	SortType int
}

type ListSalePaperResp struct {
	List []SalePaper
}

//获取图纸的销售信息
const HTTP_PLAYER_QUERY_PAPER_SALE_INFO = "api/queryPaperSaleInfo"

type QueryPaperSaleInfoReq struct {
	PaperID int64

	Sequence int
}

type QueryPaperSaleInfoResp struct {
	TradeID   int64
	PriceType int
	StartTime int64
	Duration  int64
	MaxPrice  int64
	MinPrice  int64
}

const HTTP_PLAYER_LIST_TRADE_HISTORY = "api/listPaperTradeHistory"

type ListPaperTradeHistoryReq struct {
	PaperID  int64
	Sequence int
	Page     int
}

type ListPaperTradeHistoryResp struct {
	List      []PaperTradeHistory
	TotalSale int64
}
