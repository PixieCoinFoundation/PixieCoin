package types

type EthereumTransaction struct {
	ID           int64
	Hash         string
	SubmitTime   int64
	FromUsername string
	FromNickname string
	FromAccount  string
	ToAccount    string

	//用pxc购买游戏内衣服时以下字段才有值
	ToUsername  string
	ToNickname  string
	ClothesID   string
	ClothesName string
	Price       int
	BuyCount    int

	//总金额16进制
	Amount string

	//1:eth
	//2:pxc
	AmountType int

	//0:提交
	//1:成功
	//2:失败
	//3:未知
	Status int

	//事务原始数据
	Data string
	//回执
	Receipt string

	ClothesList []SimpleClothes1
}

type SimpleClothes1 struct {
	ClothesID    string
	ClothesName  string
	ClothesCount int
	Price        int
}

type BuyDesignerInfo struct {
	Username    string
	Nickname    string
	TotalAmount int64
}
