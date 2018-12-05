package background

import (
	"appcfg"
	"constants"
	. "logger"
	. "types"
	"xlsx"
)

var bgMap map[string]*GFBackground

func init() {
	if appcfg.GetServerType() != "" && appcfg.GetServerType() != constants.SERVER_TYPE_SIMPLE_TEST && appcfg.GetServerType() != constants.SERVER_TYPE_GM {
		return
	}

	bgMap = make(map[string]*GFBackground)

	fn := "scripts/data/data_scene.xlsx"
	if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
		fn = "scripts/data/data_scene_kor.xlsx"
	}

	xlFile, err := xlsx.OpenFile(fn)
	if err != nil {
		panic(err)
	}

	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i > 0 {
				var id, name string
				var pt, pr, from int

				for j, cell := range row.Cells {
					if j == 0 {
						id, _ = cell.String()
					} else if j == 5 {
						name, _ = cell.String()
					} else if j == 2 {
						pt, _ = cell.Int()
					} else if j == 3 {
						pr, _ = cell.Int()
					} else if j == 4 {
						from, _ = cell.Int()
					}
				}

				bgMap[id] = &GFBackground{
					ID:        id,
					PriceType: pt,
					Price:     pr,
					From:      from,
					Name:      name,
				}
			}
		}
	}

	Info("bg size:", len(bgMap))
}

func BackgroundInShop(bid string) (bool, int, int) {
	bg := bgMap[bid]

	if bg != nil && bg.From == 1 {
		return true, bg.PriceType, bg.Price
	}

	return false, 0, 0
}

//GetBackgroundList 获取增加的列表
func GetBackgroundList() (costBGList map[string]*GFBackground) {
	costBGList = map[string]*GFBackground{}
	for k, v := range bgMap {
		if v != nil && v.From == 1 {
			costBGList[k] = v
		}
	}
	return costBGList
}
func GetGFBackground(bid string) (b GFBackground) {
	if v, ok := bgMap[bid]; ok {
		b = *v
	}
	return
}

func BackgroundFree(bgid string) bool {
	bg := bgMap[bgid]

	if bg != nil && bg.From == 0 {
		return true
	}

	return false
}
