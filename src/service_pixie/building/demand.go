package building

import (
	"appcfg"
	"constants"
	. "pixie_contract/api_specification"
	"strings"
	"tools"
)

var (
	//娱乐建筑 buildingID:level:building object
	allEntertainBuilding map[int64]map[int]*DemandBuilding

	//餐饮建筑 buildingID:level:building object
	allRestaurantBuilding map[int64]map[int]*DemandBuilding
)

func init() {
	if appcfg.GetServerType() != "" && appcfg.GetServerType() != constants.SERVER_TYPE_GM {
		return
	}

	loadDemandBuildings()
}

func loadDemandBuildings() {
	allEntertainBuilding = loadDemandBuildCSV("scripts/data_pixie/entertainment.csv")
	allRestaurantBuilding = loadDemandBuildCSV("scripts/data_pixie/restaurant.csv")
}

func loadDemandBuildCSV(path string) (res map[int64]map[int]*DemandBuilding) {
	records := tools.LoadCSV(path)

	res = make(map[int64]map[int]*DemandBuilding)
	for k, v := range records {
		if k > 0 {
			var building DemandBuilding
			var level int
			for k1, v1 := range v {
				if k1 == 0 {
					building.ID = tools.LoadGetInt64(v1, "id")
				} else if k1 == 1 {
					level = tools.LoadGetInt(v1, "level")
				} else if k1 == 2 {
					building.Name = strings.TrimSpace(v1)
				} else if k1 == 3 {
					building.Star = tools.LoadGetInt(v1, "star")
				} else if k1 == 4 {
					building.BuildPrice = tools.LoadGetInt64(v1, "build price")
				} else if k1 == 5 {
					building.BuildPriceType = tools.LoadGetInt(v1, "build price type")
				} else if k1 == 6 {
					building.BuildDuration = tools.LoadGetInt64(v1, "duration")
				} else if k1 == 10 {
					building.Offer = tools.LoadGetFloat64(v1, "demand")
				} else if k1 == 11 {
					building.Profit = tools.LoadGetFloat64(v1, "profit")
				} else if k1 == 12 {
					building.CleanNeed = tools.LoadGetFloat64(v1, "clean")
				}
			}

			if level <= 0 || building.ID <= 0 {
				panic("format illegal " + path)
			}

			if res[building.ID] == nil {
				res[building.ID] = make(map[int]*DemandBuilding)
			}
			res[building.ID][level] = &building
		}
	}
	if len(res) < 1 {
		panic("load building csv size 0 " + path)
	}

	return
}

func BuildingIsDemand(buildingID int64) bool {
	_, eb := allEntertainBuilding[buildingID][1]
	_, rb := allRestaurantBuilding[buildingID][1]

	return eb || rb
}

func GetEntertainBuildingByIDLevel(id int64, level int) (building *DemandBuilding, exist bool) {
	building, exist = allEntertainBuilding[id][level]
	return
}

func GetRestaurantBuildingByIDLevel(id int64, level int) (building *DemandBuilding, exist bool) {
	building, exist = allRestaurantBuilding[id][level]
	return
}
