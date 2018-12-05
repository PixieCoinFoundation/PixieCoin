package land

import (
	"constants"
	"dao_pixie"
	"database/sql"
	"encoding/json"
	. "logger"
	. "pixie_contract/api_specification"
)

func AddLand(typee LandType, location int, ou, on, oh string, os int, tx *sql.Tx) (err error) {
	if _, err = tx.Stmt(dao_pixie.AddLandStmt).Exec(ou, on, oh, os, location, typee); err != nil {
		Err(err)
	}
	return
}

func ListPlayerLand(username string) (list []Land, err error) {
	var rows *sql.Rows

	if rows, err = dao_pixie.ListPlayerLandStmt.Query(username); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		err, list = getLandsFromRows(rows)
	}

	return
}

func GetLandByStatus(status LandStatus, offset int, pageSize int) (err error, list []Land) {
	var rows *sql.Rows
	if rows, err = dao_pixie.GetLandByStatusStmt.Query(status, offset, pageSize); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		err, list = getLandsFromRows(rows)
	}
	return
}

func ListLevelUpingLand(offset int, pageSize int) (err error, list []Land) {
	var rows *sql.Rows
	if rows, err = dao_pixie.ListLevelUpingLandStmt.Query(offset, pageSize); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		err, list = getLandsFromRows(rows)
	}
	return
}

func GetLandByStatusMastername(status LandStatus, username, yesterday string) (err error, list []LandStat) {
	var rows *sql.Rows
	if rows, err = dao_pixie.ListLandWithStatStmt.Query(username, yesterday, username, constants.GET_SALE_INFO_KEY_ALL, status, status, username, username); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		err, list = getLandStatsFromRows(rows)
	}
	return
}

func GetLandNBWaitBuildByMastername(username, yesterday string) (err error, list []LandStat) {
	var rows *sql.Rows
	if rows, err = dao_pixie.ListLandWithStatStmt.Query(username, yesterday, username, constants.GET_SALE_INFO_KEY_ALL, LAND_STATUS_NB_EMPTY, LAND_STATUS_NB_RENTED_BUILD, username, username); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		err, list = getLandStatsFromRows(rows)
	}
	return
}

func GetLandWBWaitBusinessByMastername(username, yesterday string) (err error, list []LandStat) {
	var rows *sql.Rows
	if rows, err = dao_pixie.ListLandWithStatStmt.Query(username, yesterday, username, constants.GET_SALE_INFO_KEY_ALL, LAND_STATUS_WB_NORMAL, LAND_STATUS_WB_RENTED_BUSINESS, username, username); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		err, list = getLandStatsFromRows(rows)
	}
	return
}

//我的街的所有land(不显示招租和招商中)
func ListLandExceptRentingByMastername(username, yesterday string, page int, pageSize int) (err error, list []LandStat) {
	var offset int
	if page < 0 {
		offset = 0
	} else {
		offset = (page - 1) * pageSize
	}

	var rows *sql.Rows
	if rows, err = dao_pixie.ListLandWithStatNotRentingStmt.Query(username, yesterday, username, constants.GET_SALE_INFO_KEY_ALL, username, username, offset, pageSize); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		err, list = getLandStatsFromRows(rows)
	}
	return
}

func GetLandByStatusMasternameList(username string) (err error, inSaleCount int64, inSaleShop int, ownShop int, buildingShop int, waitBusiness int, waitBuild int, allClothesIDList []interface{}) {
	var rows *sql.Rows
	allClothesIDMap := make(map[int64]bool)

	if rows, err = dao_pixie.GetLandSaleInfoByMasternameStmt.Query(username, username); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		var status LandStatus
		var saleInfo string

		for rows.Next() {
			if err = rows.Scan(&status, &saleInfo); err != nil {
				Err(err)
				return
			}
			//当前拥有商店
			ownShop++

			if saleInfo != "" {
				var paperIDList []int64
				json.Unmarshal([]byte(saleInfo), &paperIDList)

				for _, id := range paperIDList {
					if !allClothesIDMap[id] {
						allClothesIDMap[id] = true
						allClothesIDList = append(allClothesIDList, id)
					}
				}
				inSaleCount += int64(len(paperIDList))
			}

			if status == LAND_STATUS_WB_IN_BUSINESS {
				//当前营业商店
				inSaleShop++
			} else if status == LAND_STATUS_NB_BUILDING {
				//装修中的商店
				buildingShop++
			} else if status == LAND_STATUS_WB_RENTED_BUSINESS || status == LAND_STATUS_WB_NORMAL {
				//空闲中的商店
				waitBusiness++
			} else if status == LAND_STATUS_NB_RENTED_BUILD || status == LAND_STATUS_NB_EMPTY {
				//空闲中的地块
				waitBuild++
			}
		}
	}
	return
}

func GetLandCheckRentEnd(offset int, pageSize int) (err error, list []Land) {
	var rows *sql.Rows
	if rows, err = dao_pixie.GetLandCheckRentEndStmt.Query(offset, pageSize); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		err, list = getLandsFromRows(rows)
	}
	return
}

func UpdateLandStatus(toStatus LandStatus, landID int64, originStatus LandStatus) (err error) {
	var res sql.Result
	var effect int64

	if res, err = dao_pixie.UpdateLandStatusStmt.Exec(toStatus, landID, originStatus); err != nil {
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

func LandRentingTimeOut(toStatus LandStatus, landID int64, originStatus LandStatus) (err error) {
	var res sql.Result
	var effect int64

	if res, err = dao_pixie.LandRentingTimeOutStmt.Exec(toStatus, landID, originStatus); err != nil {
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

func EndLandRent(toStatus LandStatus, buildingID int64, landID int64) (err error) {
	var res sql.Result
	var effect int64

	if res, err = dao_pixie.EndLandRentStmt.Exec(toStatus, buildingID, landID); err != nil {
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
