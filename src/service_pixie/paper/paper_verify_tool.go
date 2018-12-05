package paper

import (
	"appcfg"
	"constants"
	"db_pixie/paper"
	"encoding/json"
	. "logger"
	. "pixie_contract/api_specification"
	"reflect"
	"service_pixie/paperInfo"
	"sort"
	"strings"
	. "types"
)

func checkCopySuccess(paperID int64) (err error, success bool, rewardUsernameList []string) {
	target := appcfg.GetInt("fake_verify_target", constants.VERIFY_TARGET)

	var list []PaperCopy
	if err, list = PaperCopyQuery(paperID); err != nil {
		Err(err)
		return
	} else {
		supportCount := 0
		maxSupport := 0
		once := true
		rewardUsernameList = make([]string, 0)

		for _, v := range list {
			if once {
				maxSupport = v.SupportNum
				once = false
			}

			if v.SupportNum == maxSupport {
				rewardUsernameList = append(rewardUsernameList, v.Username)
			}

			supportCount += v.SupportNum
		}
		supportCount += len(list)

		if float64(supportCount) >= float64(target)*constants.COPY_REPORT_SUCCESS_LIMIT {
			success = true
			return
		}
	}
	return
}

func computePaperVerify(paperID int64) (err error, score float64, tag1, tag2, style1, style2, styleJson string, resMap map[int64]float64) {
	var tagMap, styleMap map[string]int
	var list []PaperVerify
	var scoreCount, verifySize int

	if err, list, scoreCount, verifySize, tagMap, styleMap = paper.ListPaperVerifyWithDetail(paperID); err != nil {
		return
	}

	if verifySize > 0 {
		resMap = make(map[int64]float64, 0)

		tag1, tag2, _ = getTop3ID(tagMap)
		var style3 string
		style1, style2, style3 = getTop3ID(styleMap)

		if styleOpposite(style1, style2) {
			style2 = style3
		}

		score = float64(scoreCount) / float64(verifySize)
		styleJson = genPaperStyle(score, style1, style2, styleMap[style1], styleMap[style2], verifySize)

		//记录每个评选的匹配度
		for _, v := range list {
			//记录每个verify id的匹配度
			resMap[v.ID] = paperInfo.ComputeMatchDegree(float64(v.Score), score, v.Tag1, v.Tag2, tag1, tag2, v.Style1, v.Style2, style1, style2)
		}
	} else {
		Err("paper verify size 0", paperID)
		err = constants.PaperVerifySizeZero
		return
	}

	return
}

func genPaperStyle(score float64, style1, style2 string, style1Cnt, style2Cnt, verifySize int) (style string) {
	style1Key, style1Value := genStyleKeyAndValue(style1, style1Cnt, verifySize, score)
	style2Key, style2Value := genStyleKeyAndValue(style2, style2Cnt, verifySize, score)

	ps := PaperStyle{}

	vps := reflect.ValueOf(&ps).Elem()
	tps := reflect.TypeOf(ps)
	for i := 0; i < tps.NumField(); i++ {
		f := vps.Field(i)

		sn := tps.Field(i).Name

		if strings.EqualFold(sn, style1Key) {
			f.SetInt(int64(style1Value))
		} else if strings.EqualFold(sn, style2Key) {
			f.SetInt(int64(style2Value))
		}
	}

	tb, _ := json.Marshal(ps)
	style = string(tb)
	return
}

func genStyleKeyAndValue(style string, styleCnt, verifySize int, paperScore float64) (key string, value float64) {
	if style != "" {
		ss := strings.Split(style, "_")

		if len(ss) == 2 {
			key = ss[0]
			neg := 1.0
			if ss[1] == "-1" {
				neg = -1.0
			}
			value = float64(constants.STYLE_MAX) * (float64(styleCnt) / float64(verifySize)) * (paperScore / 100.0) * neg
		}
	}

	return
}

func styleOpposite(style1, style2 string) bool {
	if style1 != "" && style2 != "" {
		s1 := strings.Split(style1, "_")
		s2 := strings.Split(style2, "_")
		if len(s1) == 2 && len(s2) == 2 {
			if s1[0] == s2[0] {
				return true
			}
		}
	}

	return false
}

func getTop3ID(pm map[string]int) (id1, id2, id3 string) {
	idList := VerifyResultList{}
	for k, v := range pm {
		var temp VerifyResult
		temp.ID = k
		temp.Num = v
		idList = append(idList, temp)
	}
	sort.Stable(idList)

	if len(idList) >= 1 {
		id1 = idList[0].ID
	}

	if len(idList) >= 2 {
		id2 = idList[1].ID
	}

	if len(idList) >= 3 {
		id3 = idList[2].ID
	}

	return
}
