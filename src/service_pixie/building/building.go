package building

import (
	"appcfg"
	"constants"
	. "pixie_contract/api_specification"
	"strings"
	"tools"
)

var (
	allBuilding map[int64]*Building
	allSkin     map[string]*Skin
) // 所有服装类建筑

func init() {
	if appcfg.GetServerType() != "" && appcfg.GetServerType() != constants.SERVER_TYPE_GM {
		return
	}

	loadClothesBuildCSV()
	loadSkinCSV()
}

func loadClothesBuildCSV() {
	allBuilding = make(map[int64]*Building, 0)
	path := "scripts/data_pixie/building.csv"

	records := tools.LoadCSV(path)
	for k, v := range records {
		if k > 0 {
			var building Building
			for k1, v1 := range v {
				if k1 == 0 {
					building.ID = tools.LoadGetInt64(v1, "id")
				} else if k1 == 1 {
					building.Name = strings.TrimSpace(v1)
				} else if k1 == 2 {
					building.Star = tools.LoadGetInt(v1, "star")
				} else if k1 == 3 {
					building.BuildingType = tools.LoadGetInt(v1, "building type")
				} else if k1 == 4 {
					building.UnLockPrice = tools.LoadGetInt64(v1, "unlock price")
				} else if k1 == 5 {
					building.UnLockPriceType = tools.LoadGetInt(v1, "unlock price type")
				} else if k1 == 6 {
					building.BuildPrice = tools.LoadGetInt64(v1, "build price")
				} else if k1 == 7 {
					building.BuildPriceType = tools.LoadGetInt(v1, "build price type")
				} else if k1 == 8 {
					building.BuildDuration = tools.LoadGetInt64(v1, "duration")
				} else if k1 == 9 {
					building.IsOfficial = tools.LoadGetBool(v1, "is official")
				} else if k1 == 10 {
					building.BuildDesc = strings.TrimSpace(v1)
				} else if k1 == 12 {
					building.Profit = tools.LoadGetFloat64(v1, "profit")
				}
			}
			allBuilding[building.ID] = &building
		}
	}
	if len(allBuilding) < 1 {
		panic("load building csv size 0")
	}
}

func loadSkinCSV() {
	allSkin = make(map[string]*Skin, 0)
	path := "scripts/data_pixie/skin.csv"

	records := tools.LoadCSV(path)
	for k, v := range records {
		if k > 0 {
			var skin Skin
			for k1, v1 := range v {
				if k1 == 0 {
					skin.ID = v1
				} else if k1 == 1 {
					skin.StreetID = v1
				} else if k1 == 2 {
					skin.Name = v1
				} else if k1 == 3 {
					skin.Price = tools.LoadGetInt64(v1, "price")
				} else if k1 == 4 {
					skin.PriceType = tools.LoadGetInt(v1, "skin price type")
				} else if k1 == 5 {
					skin.Order = v1
				} else if k1 == 6 {
					skin.PicID = v1
				} else if k1 == 7 {
					skin.SourceName = v1
				}
			}
			allSkin[skin.ID] = &skin
		}
	}
	if len(allBuilding) < 1 {
		panic("load building csv size 0")
	}
}

func GetBuildingByID(id int64) (building *Building, exist bool) {
	building, exist = allBuilding[id]
	return
}

func GetSkinByID(id string) (skin *Skin, exist bool) {
	skin, exist = allSkin[id]
	return
}
