package building

import (
	"appcfg"
	. "pixie_contract/api_specification"
	"strings"
	"tools"
)

var allCleanObjectMap map[string]*CleanObject
var cleanMaxEmployMap map[string]int

func init() {
	if appcfg.GetServerType() != "" {
		return
	}

	loadCleanObjects()
}

func GetCleanObjectDetail(oid string) (power float64, salaryInMinute float64, employFee float64) {
	v := allCleanObjectMap[oid]
	if v != nil {
		return v.Power, v.SalaryInMinute, v.EmployFee
	}

	return
}

func GetCleanObjectMaxEmployCount(oid string) int {
	return cleanMaxEmployMap[oid]
}

func loadCleanObjects() {
	allCleanObjectMap = make(map[string]*CleanObject)
	cleanMaxEmployMap = make(map[string]int)
	path := "scripts/data_pixie/clean.csv"

	records := tools.LoadCSV(path)

	for k, v := range records {
		if k > 0 {
			var c CleanObject
			var maxEmployCnt int
			for k1, v1 := range v {
				if k1 == 0 {
					c.ID = strings.TrimSpace(v1)
				} else if k1 == 1 {
					c.Name = strings.TrimSpace(v1)
				} else if k1 == 2 {
					c.Power = tools.LoadGetFloat64(v1, "power")
				} else if k1 == 3 {
					c.SalaryInMinute = tools.LoadGetFloat64(v1, "salary in minute")
				} else if k1 == 4 {
					maxEmployCnt = tools.LoadGetInt(v1, "max employ")
				} else if k1 == 5 {
					c.EmployFee = tools.LoadGetFloat64(v1, "employ fee")
				}
			}

			allCleanObjectMap[c.ID] = &c
			cleanMaxEmployMap[c.ID] = maxEmployCnt
		}
	}

	if len(allCleanObjectMap) <= 0 {
		panic("clean size 0")
	}
}
