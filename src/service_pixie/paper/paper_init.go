package paper

import (
	"appcfg"
	"constants"
	"dao_pixie"
	"database/sql"
	"db/global"
	"db_pixie/paper"
	"encoding/json"
	"fmt"

	. "logger"
	. "pixie_contract/api_specification"
	"shutdown"
	"sync"
	"time"
	"tools"
)

var officialBasePaperList []*OfficialPaper
var officialBasePaperMap map[string]*OfficialPaper
var officialBasePaperMapMD5 string
var officialBasePaperMapLock sync.RWMutex

func init() {
	if appcfg.GetServerType() == "" {
		if appcfg.GetBool("main_server", false) {
			//检查全局审核中图纸是否有审核完毕的
			go checkVerifyProcess()

			//检查在售图纸是否超时
			go checkTradeTimeout()

			//检查图纸生产进度
			go updateProductProcess()

			//每日定时统计玩家昨日的评审收益
			go settleVerifiedReward()
		}

		go refreshOfficialBasePaperLoop()
	}
}

func refreshOfficialBasePaperLoop() {
	for {
		refreshOfficialBasePaper()
		time.Sleep(60 * time.Second)
	}
}

func refreshOfficialBasePaper() {
	if nm, nl, md5, err := paper.ListAllOfficialPaperMap(); err != nil {
		Err(err)
		return
	} else {
		Info("official clothes size", len(nm), "new md5", md5, "old md5", officialBasePaperMapMD5)
		if md5 != officialBasePaperMapMD5 {
			officialBasePaperMapLock.Lock()
			defer officialBasePaperMapLock.Unlock()

			officialBasePaperMap = nm
			officialBasePaperList = nl
			officialBasePaperMapMD5 = md5
		}
	}
}

func GetAllOfficialPaper() []*OfficialPaper {
	officialBasePaperMapLock.RLock()
	defer officialBasePaperMapLock.RUnlock()

	return officialBasePaperList
}

func GetOneOfficialPaperByID(id int64) (op *OfficialPaper, exist bool) {
	officialBasePaperMapLock.RLock()
	defer officialBasePaperMapLock.RUnlock()

	memop := officialBasePaperMap[tools.GenPixieOfficialClothesID(id)]
	if memop != nil {
		op = memop
		exist = true
		return
	} else {
		if p, err := paper.GetOneOfficialPaper(id); err == nil && p.PaperID > 0 {
			op = &p
			exist = true
			return
		}
	}

	return
}

func GetOneOfficialPaperByCID(oid string) (op *OfficialPaper, exist bool) {
	src, id := tools.GetPixieClothesDetail(oid)
	if src != string(CLOTHES_ORIGIN_OFFICIAL_PAPER) {
		return
	}

	officialBasePaperMapLock.RLock()
	defer officialBasePaperMapLock.RUnlock()

	memop := officialBasePaperMap[oid]
	if memop != nil {
		op = memop
		exist = true
		return
	} else {
		if p, err := paper.GetOneOfficialPaper(id); err == nil && p.PaperID > 0 {
			op = &p
			exist = true
			return
		}
	}

	return
}

func updateProductProcess() {
	for range time.Tick(60 * time.Second) {
		updateProductProcess1()
	}
}

func checkTradeTimeout() {
	for range time.Tick(60 * time.Second) {
		checkTradeTimeout1()
	}
}

func checkVerifyProcess() {
	for range time.Tick(10 * time.Second) {
		checkVerifyPrcess1()
	}
}

func settleVerifiedReward() {
	if appcfg.GetBool("verify_temp", false) {
		for range time.Tick(2 * time.Minute) {
			settleVerifiedReward1()
		}
	} else {
		token := fmt.Sprintf("verified_reward_send_%s", time.Now().Format("20060102"))
		for range time.Tick(1 * time.Minute) {
			if time.Now().Hour() == 1 {
				if global.DoOnceJob(token) {
					settleVerifiedReward1()
				}
			}
		}
	}
}

func checkVerifyPrcess1() {
	shutdown.SimpleAddOneRequest()
	defer shutdown.SimpleDoneOneRequest()
	tNow := time.Now()
	target := appcfg.GetInt("fake_verify_target", constants.VERIFY_TARGET)

	var paperList []DesignPaper
	if rows, err := dao_pixie.GetPaperDesignCopyInfoByStatusStmt.Query(PAPER_STATUS_QUEUE); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		for rows.Next() {
			var temp DesignPaper
			if err := rows.Scan(&temp.PaperID, &temp.Cname, &temp.VerifiedCount, &temp.CopyMark); err != nil {
				Err(err)
				return
			}
			if temp.VerifiedCount >= target {
				paperList = append(paperList, temp)
			}
		}
	}

	checkPaperVerify(paperList, tNow)
}

func settleVerifiedReward1() {
	shutdown.SimpleAddOneRequest()
	defer shutdown.SimpleDoneOneRequest()
	tn := time.Now()
	yesterday := tn.AddDate(0, 0, -1).Format("20060102")

	var matchDegreeCount float64

	var paperVerifyList []PaperVerify
	personMap := make(map[string]*RewardContent)

	//查询verify表获取所有今日通过的审核信息
	if rows, err := dao_pixie.GetPaperTodayVerifiedStmt.Query(yesterday, false); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		for rows.Next() {
			var paperVerify PaperVerify
			var result PaperVerifyResult
			if err = rows.Scan(&paperVerify.ID, &paperVerify.PaperID, &paperVerify.Cname, &paperVerify.Extra, &paperVerify.Username, &paperVerify.Nickname, &paperVerify.CheckDate, &paperVerify.Tag1, &paperVerify.Tag2, &paperVerify.Style1, &paperVerify.Style2, &paperVerify.Score, &paperVerify.MatchDegree, &paperVerify.RewardPxc, &paperVerify.SubmitTime, &paperVerify.CheckoutTime, &paperVerify.Status, &paperVerify.Result); err != nil {
				Err(err)
				return
			}
			if err = json.Unmarshal([]byte(paperVerify.Result), &result); err != nil {
				Err(err)
				return
			}
			paperVerifyList = append(paperVerifyList, paperVerify)
			//匹配度
			matchDegreeCount += paperVerify.MatchDegree
		}
	}

	//插入notify信息
	dayVerifyCount := len(paperVerifyList)
	for _, v := range paperVerifyList {
		if _, exist := personMap[v.Username]; !exist {
			var content RewardContent
			personMap[v.Username] = &content
		}
		reward := float64(constants.VERIFY_REWARD_COUNT) / matchDegreeCount * v.MatchDegree
		personMap[v.Username].CurrencyVal += reward
		personMap[v.Username].VerifyResultNum++
		personMap[v.Username].DayVerifyCount = dayVerifyCount

		if _, err := dao_pixie.UpdateVerifyPaperRewardStmt.Exec(reward, true, v.ID, v.Username, v.PaperID, false); err != nil {
			Err(err)
			return
		}
	}
	for username, val := range personMap {
		paper.InsertNotify(constants.PIXIE_PXC_TYPE, val.CurrencyVal, val.VerifyResultNum, val.DayVerifyCount, username, "", tn, PIXIE_PAPER_NOTIFY_VERIFY_REWARD)
	}
}

func checkTradeTimeout1() {
	shutdown.SimpleAddOneRequest()
	defer shutdown.SimpleDoneOneRequest()

	tNow := time.Now().Unix()

	if rows, err := dao_pixie.GetAllPaperTradeStmt.Query(); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		for rows.Next() {
			var paperID, startTime, tradeID int64
			var duration, tradeStatus int
			var ownerUsername, authorUsername string

			if err := rows.Scan(&tradeID, &authorUsername, &ownerUsername, &paperID, &startTime, &duration, &tradeStatus); err != nil {
				Err(err)
				return
			}

			if tNow >= startTime+int64(duration*3600) {
				//trade timeout
				if ownerUsername != "" {
					paper.DeletePaperTradeWhenTimeout(tradeID, paperID, false, ownerUsername, tradeStatus)
				} else {
					paper.DeletePaperTradeWhenTimeout(tradeID, paperID, true, authorUsername, tradeStatus)
				}
			}
		}
	}
}

func updateProductProcess1() {
	shutdown.SimpleAddOneRequest()
	defer shutdown.SimpleDoneOneRequest()

	tNow := time.Now().Unix()
	var err error
	var rows *sql.Rows
	if rows, err = dao_pixie.GetPaperOccupyListInProductionStmt.Query(); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		for rows.Next() {
			var username string
			var paperID, productStartTime int64
			var productLeftCount, productDoneCount int

			if err = rows.Scan(&paperID, &username, &productStartTime, &productLeftCount, &productDoneCount); err != nil {
				Err(err)
				return
			}

			productCnt := tools.CheckProductNum(tNow, productStartTime, productLeftCount, productDoneCount)
			if productCnt > 0 {
				paper.DoProductPaper(username, productCnt, productStartTime, productLeftCount-productCnt, productDoneCount+productCnt, paperID, productLeftCount, productDoneCount)
			}

			Info("productLeftCount <= 0", paperID, productLeftCount, productDoneCount)
		}
	}
	return
}
