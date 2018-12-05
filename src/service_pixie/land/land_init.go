package land

import (
	"appcfg"
	"constants"
	"db_pixie/land"
	. "logger"
	. "pixie_contract/api_specification"
	"service_pixie/building"
	"time"
)

func init() {
	if appcfg.GetServerType() == "" && appcfg.GetBool("main_server", false) {
		go checkLandRentAndBuild()
	}
}

func checkLandRentAndBuild() {
	for range time.Tick(10 * time.Second) {
		//检查租期结束
		checkRent()
		//检查装修结束
		checkBuild()
		//检查升级结束
		checkLevelUp()
		//检查招租超时
		checkRentingForBuild()
		//检查招商超时
		checkRentingForBusiness()
	}
}

func checkRentingForBuild() {
	Info("checkRentingForBuild process...")
	i := 0
	for {
		offset := i * constants.PIXIE_CHECK_PAGE_SIZE
		tNow := time.Now().Unix()
		_, List := land.GetLandByStatus(LAND_STATUS_NB_RENTING_FOR_BUILD, offset, constants.PIXIE_CHECK_PAGE_SIZE)
		for _, v := range List {
			if tNow >= v.AuctionEndTime {
				Info("land", v.ID, "renting for build timeout")
				land.LandRentingTimeOut(LAND_STATUS_NB_EMPTY, v.ID, LAND_STATUS_NB_RENTING_FOR_BUILD)
			}
		}
		i++

		if len(List) < constants.PIXIE_CHECK_PAGE_SIZE {
			break
		}
	}
}

func checkRentingForBusiness() {
	Info("checkRentingForBusiness process...")
	i := 0
	for {
		offset := i * constants.PIXIE_CHECK_PAGE_SIZE
		tNow := time.Now().Unix()
		_, List := land.GetLandByStatus(LAND_STATUS_WB_RENTING_FOR_BUSINESS, offset, constants.PIXIE_CHECK_PAGE_SIZE)

		for _, v := range List {
			if tNow >= v.AuctionEndTime {
				Info("land", v.ID, "renting for business timeout")
				land.LandRentingTimeOut(LAND_STATUS_WB_NORMAL, v.ID, LAND_STATUS_WB_RENTING_FOR_BUSINESS)
			}
		}
		i++

		if len(List) < constants.PIXIE_CHECK_PAGE_SIZE {
			break
		}
	}
}

func checkLevelUp() {
	Info("check level up process..")
	i := 0
	for {
		offset := i * constants.PIXIE_CHECK_PAGE_SIZE
		tNow := time.Now().Unix()
		_, list := land.ListLevelUpingLand(offset, constants.PIXIE_CHECK_PAGE_SIZE)
		Info("check level up list size", len(list))

		for _, v := range list {
			if v.LevelUpStartTime > 0 && v.LevelUpEndTime > 0 && tNow >= v.LevelUpEndTime {
				if ok, _ := land.LandFinishLevelUp(v.ID, v.BuildingID, v.LevelUpStartTime, v.LevelUpEndTime); ok {
					Info("land", v.ID, "level up done")
				} else {
					Info("land", v.ID, "may level up from other place")
				}
			}
		}
		i++

		if len(list) < constants.PIXIE_CHECK_PAGE_SIZE {
			break
		}
	}
}

func checkBuild() {
	Info("check build process..")
	i := 0
	for {
		offset := i * constants.PIXIE_CHECK_PAGE_SIZE
		tNow := time.Now().Unix()
		_, list := land.GetLandByStatus(LAND_STATUS_NB_BUILDING, offset, constants.PIXIE_CHECK_PAGE_SIZE)
		Info("check build list size", len(list))

		for _, v := range list {
			if v.BuildEndTime > 0 && tNow >= v.BuildEndTime {
				Info("land", v.ID, "build done from LAND_STATUS_NB_BUILDING")
				if building.BuildingIsDemand(v.BuildingID) {
					land.UpdateLandStatus(LAND_STATUS_WB_IN_BUSINESS, v.ID, LAND_STATUS_NB_BUILDING)
				} else {
					land.UpdateLandStatus(LAND_STATUS_WB_NORMAL, v.ID, LAND_STATUS_NB_BUILDING)
				}

			}
		}
		i++

		if len(list) < constants.PIXIE_CHECK_PAGE_SIZE {
			break
		}
	}
}

func checkRent() {
	Info("check rent process..")
	i := 0
	for {
		offset := i * constants.PIXIE_CHECK_PAGE_SIZE
		tNow := time.Now().Unix()
		_, list := land.GetLandCheckRentEnd(offset, constants.PIXIE_CHECK_PAGE_SIZE)
		Info("check rent list size", len(list))

		for _, v := range list {
			if v.RentEndTime > 0 && tNow >= v.RentEndTime {
				Info("land", v.ID, "rent end")
				if v.BuildingID <= 0 {
					land.EndLandRent(LAND_STATUS_NB_EMPTY, 0, v.ID)
				} else {
					if v.BuildUsername == v.OwnerUsername {
						land.EndLandRent(LAND_STATUS_WB_NORMAL, v.BuildingID, v.ID)
					} else {
						land.EndLandRent(LAND_STATUS_NB_EMPTY, 0, v.ID)
					}
				}
			}
		}
		i++

		if len(list) < constants.PIXIE_CHECK_PAGE_SIZE {
			break
		}
	}
}
