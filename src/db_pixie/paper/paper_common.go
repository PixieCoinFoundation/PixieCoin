package paper

import (
	"database/sql"
	"fmt"
	. "logger"
	. "pixie_contract/api_specification"
)

func getOwnPapers(rows *sql.Rows) (err error, list []*OwnPaper) {
	list = make([]*OwnPaper, 0)

	for rows.Next() {
		var temp OwnPaper
		var id int64
		if err = rows.Scan(&id, &temp.PaperID, &temp.AuthorUsername, &temp.AuthorNickname, &temp.AuthorHead, &temp.AuthorSex, &temp.OwnerUsername, &temp.OwnerNickname, &temp.OwnerHead, &temp.OwnerSex, &temp.ClothesType, &temp.PartType, &temp.Cname, &temp.Desc, &temp.Extra, &temp.File, &temp.Style, &temp.Star, &temp.Tag1, &temp.Tag2, &temp.STag, &temp.SaleCount, &temp.Inventory, &temp.ClothPrice, &temp.ClothPriceType, &temp.Status, &temp.AuctionFailUnread, &temp.IsInProduction, &temp.ProductStartTime, &temp.ProductLeftCount, &temp.ProductDoneCount, &temp.OccupyLogID, &temp.PriceType, &temp.LastDuration, &temp.LastMaxPrice, &temp.LastMinPrice, &temp.Sequence, &temp.CirculationMax); err != nil {
			Err(err)
			return
		}
		list = append(list, &temp)
	}

	return
}

func getDesignPapers(rows *sql.Rows, exclude bool, excludeStatus int) (err error, list []DesignPaper) {
	list = make([]DesignPaper, 0)

	for rows.Next() {
		var temp DesignPaper
		if err = rows.Scan(&temp.PaperID, &temp.AuthorUsername, &temp.AuthorNickname, &temp.AuthorHead, &temp.AuthorSex, &temp.OwnerUsername, &temp.OwnerNickname, &temp.OwnerHead, &temp.OwnerSex, &temp.ClothesType, &temp.PartType, &temp.Cname, &temp.Desc, &temp.UploadTime, &temp.Extra, &temp.File, &temp.Style, &temp.Status, &temp.AuctionFailUnread, &temp.Star, &temp.Tag1, &temp.Tag2, &temp.STag, &temp.PassTime, &temp.VerifiedCount, &temp.PriceType, &temp.LastDuration, &temp.LastMaxPrice, &temp.LastMinPrice, &temp.CopyMark, &temp.Reason, &temp.CirculationMax, &temp.Circulation); err != nil {
			Err(err)
			return
		}

		if exclude && temp.Status == excludeStatus {
			continue
		}

		list = append(list, temp)
	}

	return
}

func getRecommendPapers(rows *sql.Rows) (err error, list []RecommendClothes) {
	list = make([]RecommendClothes, 0)
	for rows.Next() {
		var temp RecommendClothes
		if err = rows.Scan(&temp.PaperID, &temp.AuthorUsername, &temp.AuthorNickname, &temp.AuthorHead, &temp.AuthorSex, &temp.OwnerUsername, &temp.OwnerNickname, &temp.OwnerHead, &temp.OwnerSex, &temp.ClothesType, &temp.PartType, &temp.Cname, &temp.Desc, &temp.Extra, &temp.File, &temp.Style, &temp.Star, &temp.Tag1, &temp.Tag2, &temp.STag, &temp.Price, &temp.PriceType); err != nil {
			Err(err)
			return
		}
		temp.ClothesID = fmt.Sprintf("P-%d", temp.PaperID)

		list = append(list, temp)
	}
	return
}
