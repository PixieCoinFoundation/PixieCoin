package score

import (
	"db/cosplay/items"
	"errors"
	"math"
	"reflect"
)

import (
	. "logger"
	"service/clothes"
	. "types"
)

var AlreadyScoreErr = errors.New("player already scored")

type Weight struct {
	Type int
	W    float64
}

func GetCosSysScore(cosClothes map[string]string, cosParams CosParams) (result int64) {
	weights := map[string]*Weight{
		"hair":     &Weight{Type: 100, W: 0.2},
		"coat":     &Weight{Type: 200, W: 0.2},
		"shirt":    &Weight{Type: 300, W: 0.15},
		"dress":    &Weight{Type: 800, W: 0.35},
		"trousers": &Weight{Type: 400, W: 0.15},
		"socks":    &Weight{Type: 500, W: 0.05},
		"shoes":    &Weight{Type: 600, W: 0.15},
		"hat":      &Weight{Type: 701, W: 0.0125},
		"glass":    &Weight{Type: 702, W: 0.0125},
		"ear":      &Weight{Type: 703, W: 0.0125},
		"necklace": &Weight{Type: 704, W: 0.0125},
		"scarf":    &Weight{Type: 705, W: 0.0125},
		"belt":     &Weight{Type: 706, W: 0.0125},
		"bag":      &Weight{Type: 707, W: 0.0125},
		"other":    &Weight{Type: 708, W: 0.0125},
	}

	// 判断是否有穿外套
	_, ok := cosClothes["coat"]
	if !ok {
		weights["shirt"].W = weights["shirt"].W + weights["coat"].W
		weights["coat"].W = 0
	}
	// 判断是否有穿袜子
	_, okSocks := cosClothes["socks"]
	if !okSocks {
		weights["shoes"].W = weights["shoes"].W + weights["socks"].W
		weights["socks"].W = 0
	}
	// 饰品
	decoCount := 0
	for k, _ := range cosClothes {
		if k == "hat" || k == "glass" || k == "ear" || k == "necklace" || k == "scarf" || k == "belt" || k == "bag" || k == "other" {
			decoCount++
		}
	}
	// 根据饰品，修改饰品的权重
	for _, v := range weights {
		if v.Type >= 700 && v.Type <= 799 {
			if decoCount != 0 {
				v.W = 0.15 / float64(decoCount)
			} else {
				v.W = 0
			}
		}
	}

	attr := [...]string{
		"warm",
		"formal",
		"tight",
		"bright",
		"dark",
		"cute",
		"man",
		"tough",
		"noble",
		"strange",
		"sexy",
		"sport",
	}

	clothScore := 0.0
	for _, v := range cosClothes {
		cloth := clothes.GetClothesById(v)
		if cloth != nil {
			for _, vv := range weights {
				if cloth.SmallType == vv.Type {
					for _, vvv := range attr {
						clothValue := getValueByTag(cloth, vvv).(float64)
						levelValue := getValueByTag(&cosParams, vvv).(float64)
						clothScore = clothScore + clothValue*levelValue*vv.W
					}
				}
			}
		}
	}

	result = int64(math.Floor(clothScore) + cosParams.Adjust)

	return
}

func GetCosPlayerScore(cosItem CosItem, cosplay Cosplay, username string, score int) (result int64, err error) {
	if scores, e := items.GetUserMarks(cosplay.CosplayID, cosItem.ItemID); e != nil {
		Err(err)
		err = e
		return
	} else {
		result = 0
		for _, v := range scores {
			if v.Username == username {
				err = AlreadyScoreErr
				return
			}
			result += int64(v.Score * 50)
		}
		result += int64(score * 50)
	}

	return
}

func GetTotalCosScore(sysScore int64, playerScore float64, playerCount int) (totalScore int64) {
	sysW := 1.0
	pW := 0.0
	if playerCount <= 5 && playerCount > 0 {
		pW = 0.1
	} else if playerCount <= 10 && playerCount > 5 {
		pW = 0.2
	} else if playerCount <= 15 && playerCount > 10 {
		pW = 0.3
	} else if playerCount > 15 {
		pW = 0.4
	}

	sysW = 1 - pW
	totalScore = int64(float64(sysScore)*sysW + playerScore*pW)
	return
}

func getValueByTag(i interface{}, tagName string) (ret interface{}) {
	val := reflect.ValueOf(i).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		if tag.Get("bson") == tagName {
			ret = valueField.Interface()
			return
		}
	}

	return
}
