package types

type JobProcess struct {
	ID         int64
	Username   string
	PaperID    int64
	StartTime  int64
	EndTime    int64
	MoneyType  int   //货币类型
	CostAmount int64 //花费金额
	ProduceNum int64 //生产数量
	FinshNum   int   //生产数量
	Status     int8  //是否完成
}
