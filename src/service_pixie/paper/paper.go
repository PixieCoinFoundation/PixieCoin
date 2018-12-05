package paper

import (
	"appcfg"
	"constants"
	"db_pixie/paper"
	"encoding/json"
	. "logger"
	"math/rand"
	. "pixie_contract/api_specification"
	"time"
	"tools"
)

func BatchGetPaperFromClothesIDs(cloids []string) (res []*BasePaper, err error) {
	pids := make([]int64, 0)
	oids := make([]int64, 0)
	for _, cid := range cloids {
		src, id := tools.GetPixieClothesDetail(cid)

		if src == string(CLOTHES_ORIGIN_OFFICIAL_PAPER) {
			oids = append(oids, id)
		} else {
			pids = append(pids, id)
		}
	}

	res = make([]*BasePaper, 0)
	var list []*BasePaper
	if err, list = paper.BatchGetPaperFromRedis(pids); err != nil {
		return
	} else {
		res = append(res, list...)
	}

	if list, err = BatchGetOfficialBasePaper(oids); err != nil {
		return
	} else {
		res = append(res, list...)
	}

	return
}

func BatchGetOfficialBasePaper(cids []int64) (res []*BasePaper, err error) {
	res = make([]*BasePaper, 0)
	for _, id := range cids {
		if p, exist := GetOneOfficialPaperByID(id); exist {
			res = append(res, &p.BasePaper)
		}
	}

	return
}

func AddPaper(username, nickname, head string, sex int, clothesType, partType, cname string, desc string, uploadTime int64, icon, main, bottom, collar, shadow string, beijingFileType int, front, back string, zValue int) (err error, paperID int64) {
	extra := PaperExtra{
		ZValue: zValue,
	}

	if partType != PAPER_TYPE_BEIJING {
		if err = ProcessFile(username, &extra, icon, main, bottom, collar, shadow, uploadTime); err != nil {
			err = constants.ProcessFileError
			return
		}
	}

	cloSet := PaperFile{
		Icon:   icon,
		Main:   main,
		Bottom: bottom,
		Collar: collar,
		Shadow: shadow,

		BeijingFileType: beijingFileType,
		Front:           front,
		Back:            back,
	}

	paperFile, _ := json.Marshal(cloSet)
	data, _ := json.Marshal(extra)

	return paper.AddPaper(username, nickname, head, sex, clothesType, partType, cname, desc, uploadTime, string(paperFile), string(data))
}

func ListPaperByAuthor(authorUsername string, page, pageSize int) (error, []DesignPaper) {
	return paper.ListPaperByAuthor(authorUsername, page, pageSize)
}

func ListPaperByOwner(username, clothesType, partType string, page, pageSize int, nt int64) (error, []*OwnPaper) {
	return paper.ListPaperByOwner(username, clothesType, partType, page, pageSize, nt)
}

func ListSalePaper(page, pageSize, sortType int, qUsername string) (error, []SalePaper) {
	if page <= 1 {
		page = 0
	} else {
		page -= 1
	}
	if pageSize <= 0 {
		pageSize = constants.GET_PAPER_DEFAULT_PAGE_SIZE
	}
	offset := page * pageSize
	return paper.ListSalePaper(offset, pageSize, sortType, qUsername)
}

func DeletePaper(paperID int64, username string) (bool, error) {
	return paper.DeletePaper(paperID, username)
}

func PaperVerifyInsert(paperID int64, cname, paperExtra, username, nickname string, tag1, tag2, style1, style2 string, score int, submitTime int64) error {
	return paper.PaperVerifyInsert(paperID, cname, paperExtra, username, nickname, tag1, tag2, style1, style2, score, submitTime)
}

func PaperNotifyQuery(username string, page int, pageSize int, tNow time.Time) (error, float64, float64, []PaperNotify) {
	if page <= 1 {
		page = 0
	} else {
		page -= 1
	}
	if pageSize <= 0 {
		pageSize = constants.GET_PAPER_DEFAULT_PAGE_SIZE
	}
	offset := page * pageSize
	return paper.PaperNotifyQuery(username, offset, pageSize, tNow)
}

func PaperVerifiedDetailQuery(username string, date string) (error, []PaperVerify) {
	return paper.PaperVerifiedDetailQuery(username, date)
}

func SubmitPaperTrade(username, nickname, head string, sex, maxPrice, minPrice int, startTime, duration int64, priceType int, paperID int64, circulationMax, sequence int) (int64, error) {
	return paper.SubmitPaperTrade(username, nickname, head, sex, maxPrice, minPrice, startTime, duration, priceType, paperID, circulationMax, sequence)
}

func CancelPaperTrade(paperID int64, username string, nt int64, tradeID int64) error {
	return paper.CancelPaperTrade(paperID, username, nt, tradeID)
}

func BuyPaper(tradeID, paperID int64, sequence int, tNow int64, username, nickname, head string, sex int, pirce float64, priceType int) error {
	return paper.BuyPaper(tradeID, paperID, sequence, tNow, username, nickname, head, sex, pirce, priceType)
}

func UpdateAuctionFailUnreadStatus(paperID int64, status bool, username string, typee int) (err error) {
	if err = paper.UpdateAuctionFailUnreadStatus(username, paperID, status, typee); err != nil {
		Err(err)
		return err
	}

	return
}

func GetPaperOwnerUsername(paperID int64) (error, string) {
	return paper.GetPaperOwnerUsername(paperID)
}

func SetPaperProduct(username string, paperID int64, startTime int64, produceNum int) error {
	return paper.SetPaperProduct(username, paperID, startTime, produceNum)
}

func CancelPaperProduct(paperID int64, username string, tNow int64) error {
	return paper.CancelPaperProduct(paperID, username, tNow)
}

func QueryPaperSaleInfo(paperID int64, username string, sequence int) (tradeID int64, priceType int, startTime int64, duration int64, maxPrice int64, minPrice int64, err error) {
	return paper.QueryPaperSaleInfo(paperID, username, sequence)
}

func QueryPaperTradeHistory(paperID int64, sequence, page int) (error, []PaperTradeHistory, int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * constants.PAGE_TRADE_LOG_PAGE_SIZE
	return paper.QueryPaperTradeHistory(paperID, sequence, offset, constants.PAGE_TRADE_LOG_PAGE_SIZE)
}

// func UpdateOccupyAfterSale(paperID int64, logID int64, saleNum int, username string) error {
// 	return paper.UpdateOccupyAfterSale(paperID, logID, saleNum, username)
// }

func ClothesPricing(price int64, priceType int, paperID int64, ownerUsername string) error {
	return paper.ClothesPricing(price, priceType, paperID, ownerUsername)
}

func ClothesPricingChangePrice(price int64, paperID int64, ownerUsername string, priceType int) error {
	return paper.ClothesPricingChangePrice(price, paperID, ownerUsername, priceType)
}

func GetDesignPaperByID(paperID int64) (err error, p DesignPaper) {
	return paper.GetDesignPaperByID(paperID)
}

func GetRandomPaperForVerify(username string) (err error, isReportCopy bool, p DesignPaper) {
	//获取今天已经评审过的图纸
	idMap := make(map[int64]bool, 0)
	var paperList []DesignPaper
	var exist bool

	if !appcfg.GetBool("verify_temp", false) {
		if err, idMap = paper.GetVerifiedPaperByUsername(username); err != nil {
			return
		}
	}

	//获取评审中的图纸
	if err, paperList = paper.GetDesignPaperByStatus(); err != nil {
		return
	}

	if p, err = getRandomPaper(idMap, username, paperList); err != nil {
		return
	} else {
		if err, exist = paper.GetCopyByIDUsername(username, p.PaperID); err != nil {
			return
		} else if exist {
			isReportCopy = true
		}
	}

	return
}

func getRandomPaper(idMap map[int64]bool, username string, paperList []DesignPaper) (p DesignPaper, err error) {
	//没有作品可供审核
	if len(paperList) <= 0 {
		err = constants.PixieNoPaperForVerify
		return
	}

	tempList := []DesignPaper{}
	for _, v := range paperList {
		if !idMap[v.PaperID] && v.VerifiedCount < appcfg.GetInt("fake_verify_target", constants.VERIFY_TARGET) {
			if appcfg.GetBool("verify_temp", false) || v.AuthorUsername != username {
				tempList = append(tempList, v)
			}
		}
	}

	//所有待审核作品都已经被审核过了
	if len(tempList) <= 0 {
		err = constants.PixieNoPaperForVerify
		return
	}

	p = tempList[rand.Intn(len(tempList))]
	return
}

func PaperReportCopy(username string, paperID int64, reason string, pic string, contact string) error {
	return paper.PaperReportCopy(username, paperID, reason, pic, contact)
}

func PaperCopyQuery(paperID int64) (error, []PaperCopy) {
	return paper.PaperCopyQuery(paperID)
}

func PaperCopySupprot(copyID int64, paperID int64) error {
	return paper.PaperCopySupprot(copyID, paperID)
}

func GetPaperNotifyReward(id int64) (err error, rewardType int, rewardVal float64) {
	var notify PaperNotify
	var reward RewardContent
	if err, notify = paper.GetPaperNotifyByID(id); err != nil {
		Err(err)
		return
	}
	if notify.Content != "" {
		if err = json.Unmarshal([]byte(notify.Content), &reward); err != nil {
			Err(err)
			return
		} else {
			rewardType = reward.CurrencyType
			rewardVal = reward.CurrencyVal
		}
	} else {
		err = constants.PaperVerifyNoReward
		return
	}
	return
}

func UpdatePaperNotifyRewarded(id int64) (err error) {
	return paper.UpdatePaperNotifyRewarded(id)
}
