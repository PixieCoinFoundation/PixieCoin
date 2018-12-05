package land

import (
	"constants"
	"dao_pixie"
	"database/sql"
	. "logger"
	. "pixie_contract/api_specification"
)

func GetOneLand(landID int64) (l Land, err error) {
	if err = dao_pixie.GetOneLandStmt.QueryRow(landID).Scan(&l.ID, &l.OwnerUsername, &l.OwnerNickname, &l.OwnerHead, &l.OwnerSex, &l.RenterUsername, &l.RenterNickname, &l.RenterHead, &l.RenterSex, &l.Location, &l.Type, &l.RentStartTime, &l.RentEndTime, &l.AuctionStartTime, &l.AuctionEndTime, &l.RentDuration, &l.RentMaxPrice, &l.RentMinPrice, &l.RentPriceType, &l.BuildingID, &l.BuildingLevel, &l.BuildUsername, &l.BuildStartTime, &l.BuildEndTime, &l.LevelUpStartTime, &l.LevelUpEndTime, &l.ShopName, &l.ShopModelDetail, &l.Status, &l.SaleInfo); err != nil {
		Err(err)
	}
	return
}

func getLandsFromRows(rows *sql.Rows) (err error, list []Land) {
	list = make([]Land, 0)
	for rows.Next() {
		var l Land
		if rows.Scan(&l.ID, &l.OwnerUsername, &l.OwnerNickname, &l.OwnerHead, &l.OwnerSex, &l.RenterUsername, &l.RenterNickname, &l.RenterHead, &l.RenterSex, &l.Location, &l.Type, &l.RentStartTime, &l.RentEndTime, &l.AuctionStartTime, &l.AuctionEndTime, &l.RentDuration, &l.RentMaxPrice, &l.RentMinPrice, &l.RentPriceType, &l.BuildingID, &l.BuildingLevel, &l.BuildUsername, &l.BuildStartTime, &l.BuildEndTime, &l.LevelUpStartTime, &l.LevelUpEndTime, &l.ShopName, &l.ShopModelDetail, &l.Status, &l.SaleInfo); err != nil {
			Err(err)
			return
		}
		list = append(list, l)
	}

	return
}

func getLandStatsFromRows(rows *sql.Rows) (err error, list []LandStat) {
	list = make([]LandStat, 0)
	for rows.Next() {
		var l LandStat
		if rows.Scan(&l.ID, &l.OwnerUsername, &l.OwnerNickname, &l.OwnerHead, &l.OwnerSex, &l.RenterUsername, &l.RenterNickname, &l.RenterHead, &l.RenterSex, &l.Location, &l.Type, &l.RentStartTime, &l.RentEndTime, &l.AuctionStartTime, &l.AuctionEndTime, &l.RentDuration, &l.RentMaxPrice, &l.RentMinPrice, &l.RentPriceType, &l.BuildingID, &l.BuildingLevel, &l.BuildStartTime, &l.BuildEndTime, &l.LevelUpStartTime, &l.LevelUpEndTime, &l.ShopName, &l.ShopModelDetail, &l.Status, &l.SaleInfo, &l.YesPxc, &l.YesGold, &l.AllPxc, &l.AllGold); err != nil {
			Err(err)
			return
		}
		list = append(list, l)
	}

	return
}

func updateLandRenting(toStatus, originStatus LandStatus, auctionStartTime, auctionEndTime, rentDuration, maxPrice, minPrice int64, priceType int, ownerUsername string, landID int64) (err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_pixie.UpdatLandRentingStmt.Exec(toStatus, auctionStartTime, auctionEndTime, rentDuration, maxPrice, minPrice, priceType, landID, ownerUsername, originStatus); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect < 1 {
		err = constants.LandSqlEffectZero
	}
	return
}

func rentLand(toStatus, originStatus LandStatus, renterUsername, renterNickname, renterHead string, renterSex int, renterStartTime int64, renterEndTime int64, buildingID int64, landID int64) (err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_pixie.RentLandStmt.Exec(toStatus, renterUsername, renterNickname, renterHead, renterSex, renterStartTime, renterEndTime, landID, buildingID, originStatus); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect < 1 {
		err = constants.LandSqlEffectZero
	}
	return
}
