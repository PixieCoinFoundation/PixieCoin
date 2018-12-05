package api_specification

//购买官方衣服
const HTTP_PLAYER_BUY_OFFICIAL_CLOTHES = "api/buyOfficialClothes"

type BuyOfficialClothesReq struct {
	PaperID  int64
	BuyCount int

	PriceType int
	Price     float64
}

type BuyOfficialClothesResp struct {
	CurrentGold    int64
	CurrentPxc     float64
	CurrentClothes ClothesInfo
}

//获取官方衣服列表
const HTTP_PLAYER_LIST_OFFICIAL_CLOTHES = "api/listOfficialClothes"

type ListOfficialClothesReq struct {
}

type ListOfficialClothesResp struct {
	//key是clothes id
	Res []*OfficialPaper
}
