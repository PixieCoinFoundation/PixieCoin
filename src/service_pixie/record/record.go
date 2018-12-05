package record

import (
	"appcfg"
	"constants"
	. "logger"
	. "model"
	"strings"
	"tools"
)

var (
	allRecordMap map[string]*RawRecord
	levelIDList  []string
)

func init() {
	if appcfg.GetServerType() != "" && appcfg.GetServerType() != constants.SERVER_TYPE_GM {
		return
	}

	loadRecordCSV()
}

func GetRecordDetail(levelID string) (friendNpcID string, costTili int) {
	r := allRecordMap[levelID]
	if r != nil {
		friendNpcID = r.FriendNpcID
		costTili = r.CostTili
	}

	return
}

func LevelIDLegal(levelID string) bool {
	if allRecordMap[levelID] != nil {
		return true
	}

	return false
}

func GetLevelIDList() []string {
	return levelIDList
}

func loadRecordCSV() {
	path := "scripts/data_pixie/story.csv"
	records := tools.LoadCSV(path)

	allRecordMap = make(map[string]*RawRecord)
	levelIDList = make([]string, 0)

	for k, v := range records {
		if k > 0 {
			var record RawRecord
			for k1, v1 := range v {
				content := strings.TrimSpace(v1)
				if k1 == 0 {
					record.RecordID = content
				} else if k1 == 7 {
					record.CostTili = tools.LoadGetInt(content, "cost tili")
				} else if k1 == 1 {
					record.FriendNpcID = content
				}
			}
			allRecordMap[record.RecordID] = &record
			levelIDList = append(levelIDList, record.RecordID)
		}
	}

	size := len(allRecordMap)
	if size < 1 {
		panic("load record csv size 0")
	} else {
		Info("all level size", size)
	}
}
