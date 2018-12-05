package recommend

import (
	"appcfg"
	"constants"
	"db_pixie/recommend"
	. "logger"
	. "pixie_contract/api_specification"
	"sync"
	"time"
)

var (
	//期刊信息
	subject     Subject
	subjectLock sync.RWMutex

	subjectSuit     []RecommendClothes
	subjectSuitLock sync.RWMutex

	//推荐
	subjectDetailMedium     []*ItemsPaper
	subjectDetailMediumLock sync.RWMutex
	subjectDetailDown       []*ItemsPaper
	subjectDetailDownLock   sync.RWMutex
)

func init() {
	if appcfg.GetServerType() != "" {
		return
	}
	go loadSubjectInfo()
}

func loadSubjectInfo() {
	tNowFirst := time.Now()
	refresh(tNowFirst)
	for range time.Tick(60 * time.Second) {
		tNow := time.Now()
		refresh(tNow)
	}
}

func GetSubject() Subject {
	subjectLock.RLock()
	defer subjectLock.RUnlock()
	return subject
}

func GetSubjectSuit() []RecommendClothes {
	subjectSuitLock.RLock()
	defer subjectSuitLock.RUnlock()
	return subjectSuit
}

func GetMediumPaper() []*ItemsPaper {
	subjectDetailMediumLock.RLock()
	defer subjectDetailMediumLock.RUnlock()
	return subjectDetailMedium
}

func GetDownPaper(clothesIDList []string) (list []*ItemsPaper) {
	subjectDetailDownLock.RLock()
	defer subjectDetailDownLock.RUnlock()
	var clothesIDListTemp []string
	clothesIDMap := map[string]bool{}
	if len(clothesIDList) > 0 {
		for _, v := range clothesIDList {
			clothesIDListTemp = append(clothesIDListTemp, v)
			clothesIDMap[v] = true
		}
	}
	//按10条分页
	for _, v := range subjectDetailDown {
		if !clothesIDMap[v.ClothesID] && len(list) <= 10 {
			list = append(list, v)
		}
	}
	return list
}

func refresh(tNow time.Time) {
	if err, subje, suit := recommend.GetSubject(); err != nil {
		return
	} else {
		subjectLock.Lock()
		subject = subje
		subjectLock.Unlock()

		subjectSuitLock.Lock()
		subjectSuit = suit
		subjectSuitLock.Unlock()
	}

	if err, itemMedium := recommend.GetRecommendPool(constants.RECOMMEND_POOL); err != nil {
		Err(err)
		return
	} else {
		trl := make([]*ItemsPaper, 0)
		p1, p2, p3, p4, p5, p6, p7 := splitCustomPoolByDate(itemMedium, tNow)
		l1 := getRandomClothes(p1, p2, p3, p4, p5, p6, p7, constants.RECOMMEND_ZONE_SIZE)
		trl = append(trl, l1...)
		subjectDetailMediumLock.Lock()
		subjectDetailMedium = trl
		subjectDetailMediumLock.Unlock()

		subjectDetailDownLock.Lock()
		subjectDetailDown = itemMedium
		subjectDetailDownLock.Unlock()
		Info(len(itemMedium), "-----refresh GetRecommendPool-")
	}
}
