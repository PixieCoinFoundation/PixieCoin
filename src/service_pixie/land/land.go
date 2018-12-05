package land

import (
	"constants"
	"db_pixie/land"
	"db_pixie/paper"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	"service_pixie/building"
	"time"
)

func CheckPlayerBuildingLevelUp(username string, buildingID, landID int64, fromLevel, toLevel int) (ok bool, resp PixieRespInfo, buildPrice int64, buildPriceType int, buildDuration int64) {
	if _, existFrom := building.GetEntertainBuildingByIDLevel(buildingID, fromLevel); existFrom {
		if bTo, existTo := building.GetEntertainBuildingByIDLevel(buildingID, toLevel); existTo {
			if has, err := playerHasBuildingWithDetail(username, buildingID, fromLevel); err != nil {
				resp = PIXIE_ERR_LIST_LAND
				return
			} else if has {
				buildPrice = bTo.BuildPrice
				buildPriceType = bTo.BuildPriceType
				buildDuration = bTo.BuildDuration
				goto OK
			} else {
				resp = PIXIE_ERR_LAND_BUILDING_NOT_EXIST
				return
			}
		} else {
			resp = PIXIE_ERR_LAND_BUILDING_ILLEGAL
			return
		}
	} else if _, existFrom := building.GetRestaurantBuildingByIDLevel(buildingID, fromLevel); existFrom {
		if bTo, existTo := building.GetRestaurantBuildingByIDLevel(buildingID, toLevel); existTo {
			if has, err := playerHasBuildingWithDetail(username, buildingID, fromLevel); err != nil {
				resp = PIXIE_ERR_LIST_LAND
				return
			} else if has {
				buildPrice = bTo.BuildPrice
				buildPriceType = bTo.BuildPriceType
				buildDuration = bTo.BuildDuration
				goto OK
			} else {
				resp = PIXIE_ERR_LAND_BUILDING_NOT_EXIST
				return
			}
		} else {
			resp = PIXIE_ERR_LAND_BUILDING_ILLEGAL
			return
		}
	} else {
		resp = PIXIE_ERR_LAND_BUILDING_ILLEGAL
		return
	}

OK:
	ok = true
	return
}

func CheckPlayerBuilding(username string, buildingID, landID int64) (ok bool, resp PixieRespInfo, buildPrice int64, buildPriceType int, buildDuration int64) {
	if land, err := land.GetOneLand(landID); err != nil {
		resp = PIXIE_ERR_GET_ONE_LAND
		return
	} else {
		if b, exist := building.GetEntertainBuildingByIDLevel(buildingID, 1); exist {
			if land.OwnerUsername != username {
				resp = PIXIE_ERR_DEMAND_BUILD_NOT_SELF
				return
			}

			if has, err := playerHasBuilding(username, buildingID); err != nil {
				resp = PIXIE_ERR_LIST_LAND
				return
			} else if has {
				resp = PIXIE_ERR_LAND_BUILDING_ALREADY_EXIST
				return
			} else {
				buildPrice = b.BuildPrice
				buildPriceType = b.BuildPriceType
				buildDuration = b.BuildDuration
				goto OK
			}
		}

		if b, exist := building.GetRestaurantBuildingByIDLevel(buildingID, 1); exist {
			if land.OwnerUsername != username {
				resp = PIXIE_ERR_DEMAND_BUILD_NOT_SELF
				return
			}

			if has, err := playerHasBuilding(username, buildingID); err != nil {
				resp = PIXIE_ERR_LIST_LAND
				return
			} else if has {
				resp = PIXIE_ERR_LAND_BUILDING_ALREADY_EXIST
				return
			} else {
				buildPrice = b.BuildPrice
				buildPriceType = b.BuildPriceType
				buildDuration = b.BuildDuration
				goto OK
			}
		}

		if b, exist := building.GetBuildingByID(buildingID); exist {
			if land.Type != PIXIE_LAND_TYPE_MY_SHOP && b.BuildingType != int(land.Type) {
				resp = PIXIE_ERR_LAND_BUILDING_NOT_MATCH
				return
			} else {
				buildPrice = b.BuildPrice
				buildPriceType = b.BuildPriceType
				buildDuration = b.BuildDuration
			}
		}
	}

OK:
	ok = true
	return
}

func playerHasBuilding(username string, buildingID int64) (has bool, err error) {
	var lands []Land

	if lands, err = land.ListPlayerLand(username); err != nil {
		return
	} else {
		for _, l := range lands {
			if l.BuildingID == buildingID {
				has = true
				return
			}
		}
	}

	return
}

func playerHasBuildingWithDetail(username string, buildingID int64, level int) (has bool, err error) {
	var lands []Land

	if lands, err = land.ListPlayerLand(username); err != nil {
		return
	} else {
		for _, l := range lands {
			if l.BuildingID == buildingID && l.BuildingLevel == level {
				has = true
				return
			}
		}
	}

	return
}

func GetSaleInfo(username, yesterday string, page int, pageSize int, tNow time.Time) (err error, inSaleCount int64, inSaleInventory int64, inSaleShop int, ownShop int, buildingShop int, waitBusiness int, waitBuild int, todayPxcSale, yesterdayPxcSale, allPxcSale float64, todayGoldSale, yesterdayGoldSale, allGoldSale int64, todaySaleClothes, yesterdaySaleClothes, allSaleClothes int64, list []LandStat) {
	var allClothesIDList []interface{}
	//获取当前营业商店,拥有商店,装修中的商店,空闲中的商店,空闲中的地块
	err, inSaleCount, inSaleShop, ownShop, buildingShop, waitBusiness, waitBuild, allClothesIDList = land.GetLandByStatusMasternameList(username)
	if err != nil {
		return
	}

	//获取在售库存
	for len(allClothesIDList) < constants.PIXIE_OCCUPY_PAPERID_IN_SIZE {
		allClothesIDList = append(allClothesIDList, 0)
	}

	if err, inSaleInventory = paper.GetPaperInventoryByIDList(allClothesIDList[:constants.PIXIE_OCCUPY_PAPERID_IN_SIZE]); err != nil {
		return
	}

	//获取分页land数据
	if pageSize < 1 {
		pageSize = constants.PIXIE_DEFAULT_PAGE_SIZE
	}
	if err, list = land.ListLandExceptRentingByMastername(username, yesterday, page, pageSize); err != nil {
		return
	}
	today := tNow.Format("20060102")
	tYesterday := tNow.AddDate(0, 0, -1).Format("20060102")
	all := constants.GET_SALE_INFO_KEY_ALL
	//获取销售额
	if err, todayPxcSale, yesterdayPxcSale, allPxcSale, todayGoldSale, yesterdayGoldSale, allGoldSale, todaySaleClothes, yesterdaySaleClothes, allSaleClothes = paper.GetPaperSaleByUsernameDate(username, all, today, tYesterday); err != nil {
		return
	}
	return
}

func GetLandByStatusMastername(status LandStatus, username, yesterday string) (error, []LandStat) {
	return land.GetLandByStatusMastername(status, username, yesterday)
}

func GetLandNBWaitBuildByMastername(username, yesterday string) (error, []LandStat) {
	return land.GetLandNBWaitBuildByMastername(username, yesterday)
}

func GetLandWBWaitBusinessByMastername(username, yesterday string) (error, []LandStat) {
	return land.GetLandWBWaitBusinessByMastername(username, yesterday)
}
