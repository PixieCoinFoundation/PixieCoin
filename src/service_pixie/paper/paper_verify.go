package paper

import (
	"constants"
	dp "db_pixie/paper"
	"encoding/json"
	. "pixie_contract/api_specification"
	"time"
	"tools"
)

func checkPaperVerify(sourceList []DesignPaper, tNow time.Time) {
	date := tNow.Format("20060102")

	for _, paper := range sourceList {
		if paper.CopyMark {
			if err, success, rewardList := checkCopySuccess(paper.PaperID); success {
				if err := dp.UpdateDesignPaperAfterVerify(PAPER_STATUS_FREEZE, 0, "", "", "", tNow, paper.PaperID, PAPER_STATUS_QUEUE); err == nil {
					for _, name := range rewardList {
						dp.InsertNotify(constants.PIXIE_GOLD_TYPE, constants.VERIFY_COPY_REWARD, 0, 0, name, paper.Cname, tNow, PIXIE_PAPER_NOTIFY_COPY_REWARD)
					}

					//更新paper redis
					dp.UpdatePaperStatusRedis(paper.PaperID, PAPER_STATUS_FREEZE)
				}

				continue
			} else if err != nil {
				continue
			}
		}

		//not copy or copy report failed
		if err, score, tag1, tag2, style1, style2, styleJson, resMap := computePaperVerify(paper.PaperID); err == nil {
			star := tools.ComputeStarFromScore(int(score))
			if err := dp.UpdateDesignPaperAfterVerify(PAPER_STATUS_WAIT_PUBLISH, star, tag1, tag2, styleJson, tNow, paper.PaperID, PAPER_STATUS_QUEUE); err == nil {
				result := PaperVerifyResult{
					Score:  score,
					Tag1:   tag1,
					Tag2:   tag2,
					Style1: style1,
					Style2: style2,
				}

				data, _ := json.Marshal(result)

				//更新匹配度
				for id, matchDergee := range resMap {
					dp.UpdatePaperVerifyAfterSettle(id, date, string(data), matchDergee, tNow.Unix())
				}

				//更新paper redis
				dp.UpdatePaperVerifyRedis(paper.PaperID, star, tag1, tag2, styleJson)
			}
		}
	}
}
