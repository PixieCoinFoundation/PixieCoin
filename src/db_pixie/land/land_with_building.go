package land

import (
	"constants"
	"dao_pixie"
	"database/sql"
	"db_pixie/paper"
	"encoding/json"
	. "logger"
	. "pixie_contract/api_specification"
	"tools"
)

func LandBuildingLevelUp(opUsername string, landID, buildingID, nt int64, buildDuration int64, fromLevel int) (endTime int64, err error) {
	var res sql.Result
	var effect int64

	endTime = nt + buildDuration
	if res, err = dao_pixie.LevelupLandBuildingStmt.Exec(nt, endTime, opUsername, landID, buildingID, fromLevel); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect <= 0 {
		err = constants.LandSqlEffectZero
		return
	}

	return
}

func CancelLandBuildingLevelUp(opUsername string, landID, buildingID int64, fromLevel int) (err error) {
	var res sql.Result
	var effect int64

	if res, err = dao_pixie.CancelLevelupLandBuildingStmt.Exec(opUsername, landID, buildingID, fromLevel); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect <= 0 {
		err = constants.LandSqlEffectZero
		return
	}

	return
}

func RemoveBuilding(landID, buildingID int64, opUsername string) (success bool, err error) {
	var l Land

	if l, err = GetOneLand(landID); err != nil {
		return
	}

	var res sql.Result
	var effect int64

	if opUsername == l.OwnerUsername {
		//owner remove building
		res, err = dao_pixie.RemoveBuildingByOwnerStmt.Exec(LAND_STATUS_NB_EMPTY, landID, opUsername, LAND_STATUS_WB_NORMAL, buildingID)
	} else if opUsername == l.RenterUsername {
		//renter remove building
		if l.BuildUsername == opUsername {
			res, err = dao_pixie.RemoveBuildingByRenterStmt.Exec(LAND_STATUS_NB_RENTED_BUILD, landID, opUsername, opUsername, LAND_STATUS_WB_NORMAL, buildingID)
		} else {
			Err("build user not match", l.OwnerUsername, l.RenterUsername, l.BuildUsername, opUsername)
			err = constants.BuildUserNotMatch
			return
		}
	} else {
		Err("land not occupy", l.OwnerUsername, l.RenterUsername, l.BuildUsername, opUsername)
		err = constants.NotOccupyLand
		return
	}

	if err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		success = true
	}

	return
}

func GetLandSale(landID int64, queryUsername string, nt int64) (salePapers []*OwnPaper, md HomeShowDetail, err error) {
	var l Land

	if l, err = GetOneLand(landID); err != nil {
		return
	}

	if l.ShopModelDetail != "" {
		if err = json.Unmarshal([]byte(l.ShopModelDetail), &md); err != nil {
			Err(err)
			return
		}
	}

	occupier := tools.GetLandOccupier(&l, nt)

	if occupier == queryUsername {
		if l.Status != int(LAND_STATUS_WB_IN_BUSINESS) && l.Status != int(LAND_STATUS_WB_NORMAL) && l.Status != int(LAND_STATUS_WB_RENTED_BUSINESS) {
			err = constants.LandStatusWrong
			return
		}
	} else {
		if l.Status != int(LAND_STATUS_WB_IN_BUSINESS) {
			err = constants.LandStatusWrong
			return
		}
	}

	if l.SaleInfo != "" {
		tmpList := make([]int64, 0)
		if err = json.Unmarshal([]byte(l.SaleInfo), &tmpList); err != nil {
			Err(err)
			return
		}

		err, salePapers = paper.ListPaperOwnByOwnerIDList(tmpList, occupier)
	}

	return
}

func SetShopModel(landID, buildingID int64, username string, model HomeShowDetail) (err error) {
	mb, _ := json.Marshal(model)

	if _, err = dao_pixie.SetShopModelStmt.Exec(string(mb), landID, buildingID, username, username); err != nil {
		Err(err)
	}
	return
}

func StartBusiness(opUsername, shopName string, landID int64, paperIDList []int64, typee int) (err error) {
	var lb []byte
	if lb, err = json.Marshal(paperIDList); err != nil {
		Err(err)
		return
	}

	saleInfo := string(lb)

	var l Land
	if l, err = GetOneLand(landID); err != nil {
		return
	}

	if l.Status == int(LAND_STATUS_WB_IN_BUSINESS) && l.SaleInfo == saleInfo && l.ShopName == shopName {
		//no change just return
		return
	}

	var landFromStatus LandStatus

	if typee == 1 {
		//from owner
		landFromStatus = LAND_STATUS_WB_NORMAL
	} else if typee == 2 {
		//from renter
		if l.BuildUsername == opUsername && l.BuildStartTime >= l.RentStartTime {
			//租地人在租期内建设的建筑
			landFromStatus = LAND_STATUS_WB_NORMAL
		} else {
			landFromStatus = LAND_STATUS_WB_RENTED_BUSINESS
		}
	} else {
		err = constants.LandBusinessTypeWrong
		return
	}
	return updateBusiness(LAND_STATUS_WB_IN_BUSINESS, landFromStatus, LAND_STATUS_WB_IN_BUSINESS, shopName, opUsername, saleInfo, landID, typee, false)
}

func StopBusiness(opUsername string, landID int64, typee int) (err error) {
	var l Land
	if l, err = GetOneLand(landID); err != nil {
		return
	}

	var landToStatus LandStatus
	if typee == 1 {
		//from owner
		landToStatus = LAND_STATUS_WB_NORMAL
	} else if typee == 2 {
		//from renter
		if l.BuildUsername == opUsername && l.BuildStartTime >= l.RentStartTime {
			//租地人在租期内建设的建筑
			landToStatus = LAND_STATUS_WB_NORMAL
		} else {
			landToStatus = LAND_STATUS_WB_RENTED_BUSINESS
		}
	} else {
		err = constants.LandBusinessTypeWrong
		return
	}

	return updateBusiness(landToStatus, LAND_STATUS_WB_IN_BUSINESS, LAND_STATUS_WB_IN_BUSINESS, "", opUsername, "", landID, typee, true)
}

func updateBusiness(landToStatus, landFromStatus, landFromStatus2 LandStatus, shopName, opUsername, saleInfo string, landID int64, originType int, excludeSaleInfo bool) (err error) {
	var res sql.Result

	if originType == 1 {
		//from owner
		if excludeSaleInfo {
			res, err = dao_pixie.UpdateLandBusinessStatusOwnerStmt.Exec(landToStatus, opUsername, landID, landFromStatus, landFromStatus2)
		} else {
			res, err = dao_pixie.UpdateLandBusinessOwnerStmt.Exec(landToStatus, saleInfo, shopName, opUsername, landID, landFromStatus, landFromStatus2)
		}
	} else if originType == 2 {
		//from renter
		if excludeSaleInfo {
			res, err = dao_pixie.UpdateLandBusinessStatusRenterStmt.Exec(landToStatus, opUsername, landID, landFromStatus, landFromStatus2)
		} else {
			res, err = dao_pixie.UpdateLandBusinessRenterStmt.Exec(landToStatus, saleInfo, shopName, opUsername, landID, landFromStatus, landFromStatus2)
		}
	} else {
		err = constants.LandBusinessTypeWrong
		return
	}

	var effect int64
	if err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect <= 0 {
		Err("update business effect 0", landToStatus, landFromStatus, landFromStatus2, opUsername, saleInfo, landID, originType)
		err = constants.LandSqlEffectZero
		return
	}

	return
}

func LandFinishLevelUp(landID, buildingID, st, et int64) (ok bool, err error) {
	var res sql.Result
	var effect int64

	if res, err = dao_pixie.FinishLevelUpStmt.Exec(landID, buildingID, st, et); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		ok = true
	}
	return
}

func LandDemandUpdateStatus(fromStatus, toStatus LandStatus, ownerUsername string, landID, buildingID int64) (ok bool, err error) {
	var res sql.Result
	var effect int64

	if res, err = dao_pixie.UpdateLandDemandStatusStmt.Exec(toStatus, ownerUsername, landID, buildingID, fromStatus); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		ok = true
	}
	return
}

func LandRentForBusiness(startTime int64, auctionDuration int64, duration int64, maxPrice int64, minPrice int64, priceType int, ownerUsername string, landID int64) (err error) {
	endTime := startTime + auctionDuration
	return updateLandRenting(LAND_STATUS_WB_RENTING_FOR_BUSINESS, LAND_STATUS_WB_NORMAL, startTime, endTime, duration, maxPrice, minPrice, priceType, ownerUsername, landID)
}

func LandStopRentForBusiness(ownerUsername string, landID int64) (err error) {
	return updateLandRenting(LAND_STATUS_WB_NORMAL, LAND_STATUS_WB_RENTING_FOR_BUSINESS, 0, 0, 0, 0, 0, 0, ownerUsername, landID)
}

func RentBuildingForBusiness(renterUsername, renterNickname, renterHead string, renterSex int, auctionStartTime int64, auctionEndTime int64, buildingID int64, landID int64) (err error) {
	return rentLand(LAND_STATUS_WB_RENTED_BUSINESS, LAND_STATUS_WB_RENTING_FOR_BUSINESS, renterUsername, renterNickname, renterHead, renterSex, auctionStartTime, auctionEndTime, buildingID, landID)
}

//go store
func GetLandByRenterUsername(renterUsername string, paperID int64, tNow int64) (error, []Land) {
	var rows *sql.Rows
	var err error
	if rows, err = dao_pixie.GetLandsByRenterUsernameStmt.Query(renterUsername, tNow); err != nil {
		Err(err)
		return err, nil
	} else {
		return getLandsFromRows(rows)
		// if err, lands := getLandsFromRows(rows); err != nil {
		// 	Err(err)
		// 	return
		// } else {
		// 	for _, v := range lands {
		// 		if v.SaleInfo != "" {
		// 			var paperIDList []int64
		// 			json.Unmarshal([]byte(v.SaleInfo), &paperIDList)
		// 		}
		// 		// if
		// 	}
		// }
	}

}

// func Get
