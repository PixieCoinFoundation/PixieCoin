package api_specification

const HTTP_PLAYER_SET_SHOP_MODEL string = "api/setShopModel"

type SetShopModelReq struct {
	LandID     int64
	BuildingID int64

	ShopModel HomeShowDetail
}

type SetShopModelResp struct{}

//土地建筑升级
const HTTP_PLAYER_LAND_BUILDING_LEVEL_UP string = "api/landBuildingLevelup"

type LandBuildingLevelUpReq struct {
	CommonReq

	LandID     int64
	BuildingID int64

	FromLevel int
}

type LandBuildingLevelUpResp struct {
	LevelUpStartTime int64
	LevelUpEndTime   int64
	CurrentGold      int64
}

//土地建筑取消升级
const HTTP_PLAYER_CANCEL_LAND_BUILDING_LEVEL_UP string = "api/cancelLandBuildingLevelup"

type CancelLandBuildingLevelUpReq struct {
	CommonReq

	LandID     int64
	BuildingID int64

	FromLevel int
}

type CancelLandBuildingLevelUpResp struct{}

//建筑招商
const HTTP_PLAYER_LAND_RENT_FOR_BUSINESS string = "api/landRentForBusiness"

type LandRentForBusinessReq struct {
	LandID     int64
	BuildingID int64

	MaxPrice int64
	MinPrice int64

	//租赁时长 客户端传小时
	Duration int

	//竞拍时长 客户端传小时
	AuctionDuration int

	MoneyType int //货币类型
}

type LandRentForBusinessResp struct {
	AuctionStartTime int64
	AuctionEndTime   int64
}

//取消建筑招商
const HTTP_PLAYER_STOP_LAND_RENT_FOR_BUSINESS string = "api/stopLandRentForBusiness"

type StopLandRentForBusinessReq struct {
	LandID     int64
	BuildingID int64
}

type StopLandRentForBusinessResp struct {
}

//建筑拆除
const HTTP_PLAYER_REMOVE_BUILDING string = "api/removeBuilding"

type RemoveBuildingReq struct {
	LandID     int64
	BuildingID int64
}

type RemoveBuildingResp struct {
}

//租下建筑
const HTTP_PLAYER_RENT_BUILDING_FOR_BUSINESS string = "api/rentBuildingForBusiness"

type RentBuildingForBusinessReq struct {
	LandID     int64
	BuildingID int64
	Price      float64
	PriceType  int
}

type RentBuildingForBusinessResp struct {
	CurrentGold   int64
	RentStartTime int64
	RentEndTime   int64
}

//餐饮或娱乐建筑开始营业
const HTTP_PLAYER_DEMAND_START_BUSINESS string = "api/demandStartBusiness"

type DemandStartBusinessReq struct {
	LandID     int64
	BuildingID int64
}

type DemandStartBusinessResp struct {
	//我的街目前状态
	StreetDetail *StreetOperateResult
}

//餐饮或娱乐建筑停止营业
const HTTP_PLAYER_DEMAND_STOP_BUSINESS string = "api/demandStopBusiness"

type DemandStopBusinessReq struct {
	LandID     int64
	BuildingID int64
}

type DemandStopBusinessResp struct {
	//我的街目前状态
	StreetDetail *StreetOperateResult
}

//普通商店开始营业
const HTTP_PLAYER_LAND_START_BUSINESS string = "api/startBusiness"

type LandStartBusinessReq struct {
	//1:owner
	//2:renter
	Type int

	LandID      int64
	PaperIDList []int64
	ShopName    string
}

type LandStartBusinessResp struct {
	//我的街目前状态
	StreetDetail *StreetOperateResult
}

//普通商店停止营业
const HTTP_PLAYER_LAND_STOP_BUSINESS string = "api/stopBusiness"

type LandStopBusinessReq struct {
	//1:owner
	//2:renter
	Type int

	LandID int64
}

type LandStopBusinessResp struct {
	//我的街目前状态
	StreetDetail *StreetOperateResult
}

//获取商店商品列表
const HTTP_PLAYER_LIST_LAND_SALE = "api/listLandSale"

type ListLandSaleReq struct {
	LandID int64
}

type ListLandSaleResp struct {
	Result []*OwnPaper

	//展示的模特信息
	ShopModel HomeShowDetail
}

//购买衣服
const HTTP_PLAYER_BUY_CLOTHES = "api/buyClothes"

type BuyClothesReq struct {
	SaleUsername string
	LandID       int64
	PaperID      int64
	BuyCount     int

	PriceType int
	Price     float64
}

type BuyClothesResp struct {
	CurrentGold    int64
	CurrentPxc     float64
	CurrentClothes ClothesInfo
}

//批量购买衣服-在指定的地块上购买
const HTTP_PLAYER_BATCH_BUY_CLOTHES = "api/batchBuyClothes"

type BatchBuyClothesReq struct {
	SaleUsername string
	LandID       int64

	PaperList []*SimplePaper
	CommonReq
}

type BatchBuyClothesResp struct {
	CurrentGold int64
	CurrentPxc  float64

	SuccessClothesList []ClothesInfo
}

//批量购买衣服-只有服装信息，没有地块和卖家信息
const HTTP_PLAYER_BUY_CLOTHES_LIST = "api/buyClothesList"

type BuyClothesListReq struct {
	ClothesList []*SimpleClothes
	CommonReq
}

type BuyClothesListResp struct {
	CurrentGold int64
	CurrentPxc  float64

	SuccessClothesList []ClothesInfo
}
