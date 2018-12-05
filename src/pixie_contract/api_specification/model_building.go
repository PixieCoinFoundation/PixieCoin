package api_specification

//默认服装类建筑
type Building struct {
	ID           int64
	Name         string
	Star         int
	BuildingType int //0.小建筑 1.大建筑

	UnLockPrice     int64
	UnLockPriceType int   //解锁币种
	BuildPrice      int64 //装修费用
	BuildPriceType  int   //装修币种
	BuildDuration   int64 //装修时常(秒)

	IsOfficial bool   //是否官方
	BuildDesc  string //建筑描述

	Profit float64 //盈利能力
}

//娱乐、餐饮类建筑 特点：无视土地大小都能装修、不需要解锁、可以升级
type DemandBuilding struct {
	ID   int64
	Name string
	Star int

	NextID int64 //下一等级建筑的ID

	BuildPrice     int64 //装修费用
	BuildPriceType int   //装修币种
	BuildDuration  int64 //装修时常(秒)

	Offer     float64 //接待能力参数
	Profit    float64 //盈利能力参数
	CleanNeed float64 //卫生需求参数
}

//清洁设施
type CleanObject struct {
	ID             string
	Name           string
	Power          float64 //清洁能力
	SalaryInMinute float64 //每分钟工资
	EmployFee      float64 //雇佣费用
}

type Skin struct {
	ID         string
	StreetID   string
	Name       string
	Price      int64
	PriceType  int
	Order      string
	PicID      string
	SourceName string
}
