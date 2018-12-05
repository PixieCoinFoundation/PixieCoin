package api_specification

type Land struct {
	ID int64

	//值>=1
	Location int

	//0:默认土地类型
	//1:大型土地
	//2:我的店
	Type LandType

	OwnerUsername string
	OwnerNickname string
	OwnerHead     string
	OwnerSex      int

	RenterUsername string
	RenterNickname string
	RenterHead     string
	RenterSex      int

	//租赁完成后属性
	RentStartTime int64
	RentEndTime   int64

	//发布租赁属性
	AuctionStartTime int64
	AuctionEndTime   int64
	RentDuration     int64 //租期:单位秒
	RentMaxPrice     int64
	RentMinPrice     int64
	RentPriceType    int

	//装修相关属性
	BuildingID     int64
	BuildingLevel  int
	BuildUsername  string //建筑者
	BuildStartTime int64
	BuildEndTime   int64

	//升级先关属性
	LevelUpStartTime int64
	LevelUpEndTime   int64

	//店铺信息
	ShopName        string
	ShopModelDetail string //ModelDetail json

	//地块状态
	Status int

	//卖的商品信息
	//目前是[]int64的json表达:paper id列表
	SaleInfo string
}

type LandStat struct {
	Land

	YesPxc  int64
	YesGold int64
	AllPxc  int64
	AllGold int64
}

type ShopSale struct {
	ID        int64
	Username  string
	Pxc       float64
	Gold      int64
	Date      string //20060102或者ALL(总)
	SaleCount int64  //每日销售数量,ALL(累计的销售数量)
}
