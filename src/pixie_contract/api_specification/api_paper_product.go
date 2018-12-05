package api_specification

//生产衣服
const HTTP_PLAYER_PAPER_PRODUCT string = "api/paperProduct"

type PaperProductReq struct {
	PaperID    int64
	ProductNum int //生产数量
}
type PaperProductResp struct {
	ProductStartTime int64
}

//取消生产
const HTTP_PLAYER_CANCEL_PAPER_PRODUCT string = "api/cancelPaperProduct"

type CancelPaperProductReq struct {
	PaperID int64
}
type CancelPaperProductResp struct {
}

//服装定价
const HTTP_PLAYER_CLOTHES_PRICING string = "api/clothesPricing"

type PaperClothesPricingReq struct {
	PaperID   int64
	Price     int64
	PriceType int
	Type      int //1.服装定价，2.修改服装定价
}

type PaperClothesPricingResp struct {
}
