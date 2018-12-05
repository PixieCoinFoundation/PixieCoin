package types

type Trade struct {
	ID              int64
	OwnerUsername   string //拥有者名称
	OwnerNickname   string // 拥有者昵称
	MaxPrice        int64  //最高出价
	MinPrice        int64  //最低出价
	StartTime       int64  //开始时间
	EndTime         int64  //结束时间
	PaperID         int64  //图纸ID
	AuctionTime     int64  //竞拍时间
	AuctionUsername string //竞拍者名称
	AuctionNickname string //竞拍者昵称
	PriceType       int    //价格类型
	Price           int    //价格
	Status          int8   //交易状态
}
