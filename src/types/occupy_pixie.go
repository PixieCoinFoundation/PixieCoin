package types

type Occupy struct {
	ID             int64
	Username       string
	Nickname       string
	PaperID        int
	MoneyAmount    int64 //金币销售量
	PxcAmount      int64 //pxc销售量
	MoneyInventory int   // 金币投资额
	PxcInvertory   int   //pxc投资额
	SaleCount      int   //售卖数量
	StockCount     int   //库存数量
	IsInProduction bool  //是否是生产中
}
