package land

import (
	"constants"
	"dao_pixie"
	"database/sql"

	. "logger"
	. "pixie_contract/api_specification"
)

func LandStartRentingForBuild(nt, auctionDuration, rentDuration, maxPrice, minPrice int64, priceType int, ownerUsername string, landID int64) (err error) {
	auctionEndTime := nt + auctionDuration
	return updateLandRenting(LAND_STATUS_NB_RENTING_FOR_BUILD, LAND_STATUS_NB_EMPTY, nt, auctionEndTime, rentDuration, maxPrice, minPrice, priceType, ownerUsername, landID)
}

func LandStopRentingForBuild(ownerUsername string, landID int64) (err error) {
	return updateLandRenting(LAND_STATUS_NB_EMPTY, LAND_STATUS_NB_RENTING_FOR_BUILD, 0, 0, 0, 0, 0, 0, ownerUsername, landID)
}

func RentLandForBuild(renterUsername, renterNickname, renterHead string, renterSex int, renterStartTime int64, renterEndTime int64, landID int64) (err error) {
	return rentLand(LAND_STATUS_NB_RENTED_BUILD, LAND_STATUS_NB_RENTING_FOR_BUILD, renterUsername, renterNickname, renterHead, renterSex, renterStartTime, renterEndTime, 0, landID)
}

func LandStartBuilding(opUsername string, landID, buildingID, nt int64, buildDuration int64, buildType int) (endTime int64, err error) {
	var landFromStatus LandStatus
	var res sql.Result
	var buildStmt *sql.Stmt
	var effect int64

	if buildType == 1 {
		landFromStatus = LAND_STATUS_NB_EMPTY
		buildStmt = dao_pixie.UpdateLandBuildingOwnerStmt
	} else if buildType == 2 {
		landFromStatus = LAND_STATUS_NB_RENTED_BUILD
		buildStmt = dao_pixie.UpdateLandBuildingRenterStmt
	} else {
		err = constants.LandBuildTypeWrong
		return
	}

	endTime = nt + buildDuration
	if res, err = buildStmt.Exec(buildingID, 1, opUsername, nt, endTime, LAND_STATUS_NB_BUILDING, opUsername, landID, landFromStatus, 0); err != nil {
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

func StopLandBuilding(opUsername string, landID, buildingID int64, buildType int) (err error) {
	var landToStatus LandStatus
	var res sql.Result
	var buildStmt *sql.Stmt
	var effect int64

	if buildType == 1 {
		landToStatus = LAND_STATUS_NB_EMPTY
		buildStmt = dao_pixie.UpdateLandBuildingOwnerStmt
	} else if buildType == 2 {
		landToStatus = LAND_STATUS_NB_RENTED_BUILD
		buildStmt = dao_pixie.UpdateLandBuildingRenterStmt
	} else {
		err = constants.LandBuildTypeWrong
		return
	}

	if res, err = buildStmt.Exec(0, 0, "", 0, 0, landToStatus, opUsername, landID, LAND_STATUS_NB_BUILDING, buildingID); err != nil {
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
