package blog

import (
	"appcfg"
	"constants"
	. "logger"
	"strconv"
	"xlsx"
)

var blogRewardMap map[int]int

func init() {
	if appcfg.GetServerType() != "" {
		return
	}
	blogRewardMap = make(map[int]int)

	f := "scripts/data/data_blog.xlsx"
	if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
		f = "scripts/data/data_blog_kor.xlsx"
	}

	xlFile, err := xlsx.OpenFile(f)
	if err != nil {
		panic(err)
		return
	}

	cnt := 0
	for _, sheet := range xlFile.Sheets {
		if sheet.Name == "Sheet1" {
			for i, row := range sheet.Rows {
				if i > 0 {
					var blogID, diamond int
					for j, cell := range row.Cells {
						str, _ := cell.String()
						if j == 1 {
							blogID, _ = strconv.Atoi(str)
						} else if j == 10 {
							diamond, _ = strconv.Atoi(str)
						}
					}
					cnt++
					blogRewardMap[blogID] = diamond
				}
			}
		}
	}

	if cnt != len(blogRewardMap) {
		Info(cnt, len(blogRewardMap))
		panic("blog cnt not match")
	}
}

func GetDiamondByBlog(blogID int) int {
	return blogRewardMap[blogID]
}
