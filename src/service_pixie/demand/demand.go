package demand

import (
	"appcfg"
	"tools"
)

type DemandDetail struct {
	From  float64 //>=
	To    float64 //<
	Value float64
}

var demandCleanList []*DemandDetail
var demandExtertainList []*DemandDetail
var demandRestaurantList []*DemandDetail

func init() {
	if appcfg.GetServerType() != "" {
		return
	}

	loadDemands()
}

func GetVisitCapCleanNeed(visitCap float64) float64 {
	var maxValue float64
	for _, d := range demandCleanList {
		if maxValue < d.Value {
			maxValue = d.Value
		}
		if visitCap >= d.From && visitCap < d.To {
			return d.Value
		}
	}

	return maxValue
}

func GetVisitCapEntertainNeed(visitCap float64) float64 {
	var maxValue float64
	for _, d := range demandExtertainList {
		if maxValue < d.Value {
			maxValue = d.Value
		}
		if visitCap >= d.From && visitCap < d.To {
			return d.Value
		}
	}

	return maxValue
}

func GetVisitCapRestaurantNeed(visitCap float64) float64 {
	var maxValue float64
	for _, d := range demandRestaurantList {
		if maxValue < d.Value {
			maxValue = d.Value
		}
		if visitCap >= d.From && visitCap < d.To {
			return d.Value
		}
	}

	return maxValue
}

func loadDemands() {
	demandCleanList = loadDemandFile("scripts/data_pixie/demand_clean.csv")
	demandExtertainList = loadDemandFile("scripts/data_pixie/demand_entertainment.csv")
	demandRestaurantList = loadDemandFile("scripts/data_pixie/demand_restaurant.csv")
}

func loadDemandFile(path string) (res []*DemandDetail) {
	records := tools.LoadCSV(path)

	res = make([]*DemandDetail, 0)
	for k, v := range records {
		if k > 0 {
			var l DemandDetail
			for k1, v1 := range v {
				if k1 == 0 {
					l.From = tools.LoadGetFloat64(v1, "from "+path)
				} else if k1 == 1 {
					l.To = tools.LoadGetFloat64(v1, "to "+path)
				} else if k1 == 2 {
					l.Value = tools.LoadGetFloat64(v1, "value "+path)
				}
			}
			res = append(res, &l)
		}
	}

	return
}
