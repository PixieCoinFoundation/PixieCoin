package paperInfo

import (
	"appcfg"
	"constants"
	"fmt"
	"math"
	. "pixie_contract/api_specification"
	"reflect"
	"strings"
	"tools"
)

var tagMap map[string]string   //key是id value是tag名称
var styleMap map[string]string //key是style id value是style名字
var sTagMap map[string]string  //key是stag id value是stag名字

var styleStructNameMap map[string]bool

func init() {
	if appcfg.GetServerType() != "" && appcfg.GetServerType() != constants.SERVER_TYPE_GM {
		return
	}

	loadStyleStructNameMap()
	loadPaperTag()
	loadPaperStyle()
	loadPaperSTag()
}

func ComputeMatchDegree(playerScore, resultScore float64, tag1, tag2, resultTag1, resultTag2, style1, style2, resultStyle1, resultStyle2 string) float64 {
	var matchTag, matchStyle int
	var resultTagCount, resultStyleCount float64
	if TagLegal(resultTag1) {
		resultTagCount++
	}
	if TagLegal(resultTag2) {
		resultTagCount++
	}
	if (tag1 == resultTag1 || tag1 == resultTag2) && tag1 != "" {
		matchTag++
	}
	if (tag2 == resultTag1 || tag2 == resultTag2) && tag2 != "" {
		matchTag++
	}

	if StyleLegal(resultStyle1) {
		resultStyleCount++
	}
	if StyleLegal(resultStyle2) {
		resultStyleCount++
	}
	if (style1 == resultStyle1 || style1 == resultStyle2) && style1 != "" {
		matchStyle++
	}
	if (style2 == resultStyle1 || style2 == resultStyle2) && style2 != "" {
		matchStyle++
	}

	if resultTagCount <= 0 {
		return (1.0 - (math.Abs(resultScore-playerScore) / 100.0) + float64(matchStyle)/resultStyleCount) / 2.0
	} else {
		return (1.0 - (math.Abs(resultScore-playerScore) / 100.0) + float64(matchStyle)/resultStyleCount + float64(matchTag)/resultTagCount) / 3.0
	}
}

func GetSTagMap() map[string]string {
	return sTagMap
}

func GetStyleMap() map[string]string {
	return styleMap
}

func GetTagMap() map[string]string {
	return tagMap
}

func TagLegal(tag string) bool {
	if tagMap[tag] != "" {
		return true
	}

	return false
}

func StyleLegal(id string) bool {
	if styleMap[id] != "" {
		return true
	}

	return false
}

func loadStyleStructNameMap() {
	styleStructNameMap = make(map[string]bool)

	t := PaperStyle{}

	vt := reflect.ValueOf(t)
	tt := reflect.TypeOf(t)
	for i := 0; i < tt.NumField(); i++ {
		fname := tt.Field(i).Name
		fkind := vt.Field(i).Kind()

		if fkind != reflect.Int {
			panic("struct " + tt.Name() + " field not int " + fname)
		}

		styleStructNameMap[strings.ToLower(fname)] = true
	}
}

func loadPaperTag() {
	tagMap = make(map[string]string)

	path := "scripts/data_pixie/tag.csv"

	records := tools.LoadCSV(path)

	for k, v := range records {
		if k > 0 {
			var id, name string
			for k1, v1 := range v {
				if k1 == 0 {
					id = strings.TrimSpace(v1)
				} else if k1 == 1 {
					name = strings.TrimSpace(v1)
				}
			}

			if id == "" || name == "" {
				continue
			}

			if !strings.HasPrefix(id, "tag_") {
				panic("unknown tag prefix " + id)
			}

			id = id[4:]

			tagMap[id] = name
		}
	}
}

func loadPaperSTag() {
	sTagMap = make(map[string]string)

	path := "scripts/data_pixie/hidden_tag.csv"

	records := tools.LoadCSV(path)

	for k, v := range records {
		if k > 0 {
			var id, part, name string
			for k1, v1 := range v {
				if k1 == 0 {
					id = strings.TrimSpace(v1)
				} else if k1 == 1 {
					part = strings.TrimSpace(v1)
				} else if k1 == 2 {
					name = strings.TrimSpace(v1)
				}
			}
			if id == "" || part == "" || name == "" {
				continue
			}

			shouldPrefix := "hidden_tag_" + part + "_"
			if !strings.HasPrefix(id, shouldPrefix) {
				panic("illegal stag " + id)
			}

			id = id[len("hidden_tag_"):]
			sTagMap[id] = name
		}
	}
}

func loadPaperStyle() {
	styleMap = make(map[string]string)

	path := "scripts/data_pixie/style.csv"

	records := tools.LoadCSV(path)

	for k, v := range records {
		if k > 0 {
			var key, neg, name string
			for k1, v1 := range v {
				if k1 == 0 {
					key = strings.TrimSpace(strings.ToLower(v1))
				} else if k1 == 1 {
					neg = strings.TrimSpace(v1)
				} else if k1 == 2 {
					name = strings.TrimSpace(v1)
				}
			}

			if key == "" || neg == "" || name == "" {
				continue
			}

			if !styleStructNameMap[key] {
				panic("style key not exist in struct " + key)
			}

			neg = strings.TrimSpace(neg)
			if neg != "1" && neg != "-1" {
				panic("style neg wrong " + neg)
			}

			styleMap[fmt.Sprintf("%s_%s", key, neg)] = name
		}
	}
}
