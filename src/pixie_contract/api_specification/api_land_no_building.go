package api_specification

//土地招租
const HTTP_PLAYER_LAND_RENT_FOR_BUILD string = "api/landRentForBuild"

type LandRentForBuildReq struct {
	LandID   int64
	MaxPrice int64
	MinPrice int64

	//租赁时长 客户端传小时
	Duration int

	//竞拍时长 客户端传小时
	AuctionDuration int

	MoneyType int //货币类型
}

type LandRentForBuildResp struct {
	AuctionStartTime int64
	AuctionEndTime   int64
}

//取消土地招租
const HTTP_PLAYER_STOP_LAND_RENT_FOR_BUILD string = "api/stopLandRentForBuild"

type StopLandRentForBuildReq struct {
	LandID int64
}

type StopLandRentForBuildResp struct {
}

//土地装修
const HTTP_PLAYER_LAND_BUILDING string = "api/landBuilding"

type LandBuildingReq struct {
	//1:owner装修
	//2:renter装修
	BuildType int

	LandID     int64
	BuildingID int64
}

type LandBuildingResp struct {
	BuildStartTime int64
	BuildEndTime   int64
	CurrentGold    int64
	CurrentPxc     float64
}

//停止土地装修
const HTTP_PLAYER_STOP_LAND_BUILDING string = "api/stopLandBuilding"

type StopLandBuildingReq struct {
	//1:owner停止装修
	//2:renter停止装修
	BuildType int

	LandID     int64
	BuildingID int64
}

type StopLandBuildingResp struct {
}

//租下地块
const HTTP_PLAYER_RENT_LAND_FOR_BUILD string = "api/rentLandForBuild"

type RentLandForBuildReq struct {
	Price     int64
	PriceType int
	LandID    int64
}

type RentLandForBuildResp struct {
	CurrentGold   int64
	RentStartTime int64
	RentEndTime   int64
}
