package api_specification

const DEPRECATED_HTTP_PLAYER_LAND_MANAGE_HANDLE = "api/shopManage"
const HTTP_PLAYER_LAND_MANAGE_HANDLE = "api/landManage"

type LandManageReq struct {
	Page     int
	PageSize int //可以不传,默认50
}

type LandManageResp struct {
	TodayPxcSale      float64 //今日pxc销量
	TodayGoldSale     int64   //今日gold销量
	YesterdayPxcSale  float64 //昨日pxc销量
	YesterdayGoldSale int64   //昨日gold销量
	AllPxcSale        float64 //总pxc销售额
	AllGoldSale       int64   //总金币销售额

	TodaySaleClothes     int64 //今日卖出服装
	YesterdaySaleClothes int64 //昨日卖出服装
	AllSaleClothes       int64 //累计卖出服装

	//inSaleShop int, ownShop int, buildingShop int, rentedBusiness int, nbRentedBuild int,
	InSaleShopCount   int //在售商店数量
	OwnShopCount      int //拥有的商店数量
	BuildingLandCount int //建造中数量

	//错误的两个参数 未来会去掉
	WBRentedBusinessCount int //空闲中建筑数量
	NBRentedBuildCount    int //空闲地块数量
	//正确的两个参数
	WaitBusinessCount int //空闲中建筑数量
	WaitBuildCount    int //空闲中地块数量

	InSaleClothesNum       int64 //在售服装种类
	InSaleClothesInventory int64 //在售服装库存

	List []LandStat
}

const HTTP_PLAYER_GET_STATUS_LAND_HANDLE = "api/getStatusLand"

type GetStatusLandReq struct {
	//1.营业中
	//2.装修中
	//3.地块有建筑 但没有开业
	//4.地块空闲中 还没有装修
	LandType int
}
type GetStatusLandResp struct {
	List []LandStat
}

const HTTP_PLAYER_GO_STORE string = "api/goStore"

type GoStoreReq struct {
	ClothesID string //图纸ID
}

type GoStoreResp struct {
	// 玩家商店
	Land
	Result []*OwnPaper
	//展示的模特信息
	ShopModel HomeShowDetail

	//官方商店
	Res []*OfficialPaper
}

const HTTP_PLAYER_GET_SALE_CLOTHES string = "api/getSaleClothes"

type GetSaleClothesReq struct {
	Page     int
	PageSize int

	//可选参数-模特类型
	ClothesType string

	//可选参数-部位类型
	PartType string
}

type GetSaleClothesResp struct {
	List []*OwnPaper
}
