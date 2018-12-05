package dao

import (
	"constants"
	. "pixie_contract/api_specification"
)

func GetOnePaperByID(id int64) (err error, exist bool, paper PaperInfo) {
	var ownPaper OwnPaper
	var designPaper DesignPaper
	var salePaper SalePaper
	var have bool

	designPaper.PaperID = id
	salePaper.PaperID = id

	if exist, err = GetEngine().Table(constants.GM_TABLE_NAME_PAPER).ID(id).Get(&ownPaper); err != nil {
		// Err(err)
		return
	} else if !exist {
		return
	}
	paper.BasePaper = ownPaper.BasePaper
	paper.SaleCount = ownPaper.SaleCount                 //销售数量
	paper.ClothPrice = ownPaper.ClothPrice               //衣服价格
	paper.ClothPriceType = ownPaper.ClothPriceType       //衣服价格类型
	paper.Inventory = ownPaper.Inventory                 //库存
	paper.IsInProduction = ownPaper.IsInProduction       // 是否生产中
	paper.ProductStartTime = ownPaper.ProductStartTime   //开始生产时间
	paper.ProductLeftCount = ownPaper.ProductLeftCount   //剩余生产数量
	paper.ProductDoneCount = ownPaper.ProductDoneCount   //已经生产的
	paper.OccupyLogID = ownPaper.OccupyLogID             //销售记录ID
	paper.AuctionFailUnread = ownPaper.AuctionFailUnread //竞拍失败-是否未读
	paper.PriceType = ownPaper.PriceType

	if have, err = GetEngine().Table(constants.GM_TABLE_NAME_PAPER_VERIFY).Get(&designPaper); err != nil {
		// Err(err)
		return
	} else if have {
		paper.UploadTime = designPaper.UploadTime
		paper.PassTime = designPaper.PassTime
		// paper.VerifyCount = designPaper.VerifyCount

	}

	if have, err = GetEngine().Table(constants.GM_TABLE_NAME_OCCUPY).Get(&salePaper); err != nil {
		// Err(err)
		return
	} else if have {
		paper.PriceType = salePaper.PriceType //图纸价格类型
		paper.StartTime = salePaper.StartTime
		paper.Duration = salePaper.Duration
		paper.MaxPrice = salePaper.MaxPrice
		paper.MinPrice = salePaper.MinPrice
		paper.AuctionTime = salePaper.AuctionTime
		paper.AuctionUsername = salePaper.AuctionUsername
		paper.AuctionNickname = salePaper.AuctionNickname
		paper.AuctionPrice = salePaper.AuctionPrice
	}
	return
}

func FindPaperByIDList(idList []interface{}) (err error, list []OwnPaper) {
	if err = GetEngine().Table(constants.GM_TABLE_NAME_PAPER).In("id", idList).Find(&list); err != nil {
		return
	}
	return
}
