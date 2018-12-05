package level

import (
	"appcfg"
	"fmt"
	. "logger"
	"tools"
)

var levelMap map[int]*LevelDetail
var levelExpMap map[string]int
var MaxLevel int

type LevelDetail struct {
	CurrentLevel int
	NeedExp      int
	VisitCap     float64
	WalletMax    float64
	TiliMax      int
}

func init() {
	if appcfg.GetServerType() != "" {
		return
	}

	loadLevel()
}

func GetLevelMaxTili(level int) int {
	v := levelMap[level]
	if v != nil {
		return v.TiliMax
	}

	return 0
}

func GetLevelVisitCap(level int) float64 {
	v := levelMap[level]
	if v != nil {
		return v.VisitCap
	}

	return 0
}

func GetLevelWalletMax(level int) float64 {
	v := levelMap[level]
	if v != nil {
		return v.WalletMax
	}

	return 0
}

func GetLevelNeedExp(level int) int {
	v := levelMap[level]
	if v != nil {
		return v.NeedExp
	}

	return 0
}

func GetLevelExpMap() map[string]int {
	return levelExpMap
}

func loadLevel() {
	levelMap = make(map[int]*LevelDetail)
	levelExpMap = make(map[string]int)
	path := "scripts/data_pixie/player.csv"

	records := tools.LoadCSV(path)
	for k, v := range records {
		if k > 0 {
			var l LevelDetail
			for k1, v1 := range v {
				if k1 == 0 {
					l.CurrentLevel = tools.LoadGetInt(v1, "level")
				} else if k1 == 1 {
					l.VisitCap = tools.LoadGetFloat64(v1, "visitor gap")
				} else if k1 == 2 {
					l.WalletMax = tools.LoadGetFloat64(v1, "wallet max")
				} else if k1 == 3 {
					l.NeedExp = tools.LoadGetInt(v1, "need exp")
				} else if k1 == 4 {
					l.TiliMax = tools.LoadGetInt(v1, "tili max")
				}
			}

			if l.CurrentLevel <= 0 {
				continue
			}

			levelMap[l.CurrentLevel] = &l
			if MaxLevel < l.CurrentLevel {
				MaxLevel = l.CurrentLevel
			}
			levelExpMap[fmt.Sprintf("%d", l.CurrentLevel)] = l.NeedExp
		}
	}

	size := len(levelMap)
	if size < 1 {
		panic("load player level csv size 0")
	} else {
		Info("all player level size", size)
	}
}
