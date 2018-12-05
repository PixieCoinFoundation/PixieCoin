package designerZone

import (
	// "appcfg"
	"constants"
	"db/designerzone"
	. "logger"
	"math/rand"
	"strconv"
	"sync"
	"time"
	. "types"
)

func init() {
	// if appcfg.GetServerType() != "" {
	// 	return
	// }
	// bannerList = make([]DesignerZoneBanner, 0)
	// topicList = make([]Topic, 0)

	// refreshLoop()
	// go loop()
}

func loop() {
	for {
		time.Sleep(60 * time.Second)

		refreshLoop()
	}
}

var (
	recommendDesiger     GFDesignerPool //推荐设计师
	recommendDesigerLock sync.RWMutex

	stickDesigerList     []*GFDesignerPool // 置顶设计师列表
	stickDesigerListLock sync.RWMutex

	stickNewDesigerList     []*GFDesignerPool // 置顶新人设计师列表
	stickNewDesigerListLock sync.RWMutex

	recommendNewDesiger     GFDesignerPool //推荐新人设计师
	recommendNewDesigerLock sync.RWMutex

	stickNewClothesList     []*CustomPool
	stickNewClothesListLock sync.RWMutex

	stickPrefectClothesList     []*CustomPool
	stickPrefectClothesListLock sync.RWMutex

	publicBuyList     []*CustomPool
	publicBuyListLock sync.RWMutex

	bannerList     []DesignerZoneBanner
	bannerListLock sync.RWMutex

	topicList     []Topic
	topicListLock sync.RWMutex

	stickTopicList     []Topic
	stickTopicListLock sync.RWMutex

	recommendCustomsLock sync.RWMutex
	recommendCustoms     []*CustomPool
)

func refreshLoop() {
	tNow := time.Now()

	//加载设计师相关推荐
	designerList, newDesigerList, recommendDesingerList, newRecommendDesignerList, err := designerzone.GetDesignerIDListFromDesignerPool()
	if err != nil {
		Err(err)
		Info(err)
	} else {
		designer := getRandomDesigner(designerList)
		recommendDesigerLock.Lock()
		Info("recommendDesiger", designer.Username)
		recommendDesiger = designer
		recommendDesigerLock.Unlock()

		designerNew := getRandomDesigner(newDesigerList)
		recommendNewDesigerLock.Lock()
		Info("recommendNewDesiger", designerNew.Username)
		recommendNewDesiger = designerNew
		recommendNewDesigerLock.Unlock()

		stickDesigerListLock.Lock()
		Info("stickDesigerList size", len(recommendDesingerList))
		stickDesigerList = recommendDesingerList
		stickDesigerListLock.Unlock()

		stickNewDesigerListLock.Lock()
		Info("stickNewDesigerList size", len(newRecommendDesignerList))
		stickNewDesigerList = newRecommendDesignerList
		stickNewDesigerListLock.Unlock()
	}

	//加载作品相关推荐
	clothesList, newClothesList, recommendClothesList, recommendNewClothesList, publicBuyListTemp, err := designerzone.GetCustomIDListFromCustomPool()
	if err != nil {
		Err(err)
	} else {
		trl := make([]*CustomPool, 0)
		p1, p2, p3, p4, p5, p6 := splitCustomPoolByDate(clothesList, tNow)
		l1 := getRandomClothes(p1, p2, p3, p4, p5, p6, constants.DESIGNER_ZONE_CUSTOM_SIZE)
		trl = append(trl, l1...)

		newp1, newp2, newp3, newp4, newp5, newp6 := splitCustomPoolByDate(newClothesList, tNow)
		l2 := getRandomClothes(newp1, newp2, newp3, newp4, newp5, newp6, constants.DESIGNER_ZONE_CUSTOM_SIZE)
		trl = append(trl, l2...)

		recommendCustomsLock.Lock()
		Info("recommendCustoms size", len(trl))
		recommendCustoms = trl
		recommendCustomsLock.Unlock()

		stickNewClothesListLock.Lock()
		Info("stickNewClothesList size", len(recommendNewClothesList))
		stickNewClothesList = recommendNewClothesList
		stickNewClothesListLock.Unlock()

		stickPrefectClothesListLock.Lock()
		Info("stickPrefectClothesList size", len(recommendClothesList))
		stickPrefectClothesList = recommendClothesList
		stickPrefectClothesListLock.Unlock()

		publicBuyListLock.Lock()
		Info("publicBuyList size", len(publicBuyListTemp))
		publicBuyList = publicBuyListTemp
		publicBuyListLock.Unlock()

	}

	bannerListTemp, err := designerzone.GetBannerListFromCBannerPool(tNow)
	if err != nil {
		Err(err)
	} else {
		bannerListLock.Lock()
		Info("bannerList size", len(bannerListTemp))
		bannerList = bannerListTemp
		bannerListLock.Unlock()
	}

	topicListtemp, stickTopic, err := designerzone.GetTopicListFromTopicPool(tNow)
	if err != nil {
		Err(err)
	} else {
		topicListLock.Lock()
		Info("topicList size", len(topicListtemp))
		topicList = topicListtemp
		topicListLock.Unlock()

		stickTopicListLock.Lock()
		Info("stickTopicList size", len(stickTopic))
		stickTopicList = stickTopic
		stickTopicListLock.Unlock()
	}
}

func GetRecommendCustoms() []*CustomPool {
	recommendCustomsLock.RLock()
	defer recommendCustomsLock.RUnlock()

	return recommendCustoms
}

func GetRandDesigner() GFDesignerPool {
	recommendDesigerLock.RLock()
	defer recommendDesigerLock.RUnlock()

	return recommendDesiger
}

func GetRandNewDesinger() GFDesignerPool {
	recommendNewDesigerLock.RLock()
	defer recommendNewDesigerLock.RUnlock()

	return recommendNewDesiger
}

//获取新人置顶设计师
func GetRecommendRandDesigner() (designer GFDesignerPool) {
	stickDesigerListLock.RLock()
	defer stickDesigerListLock.RUnlock()

	length := len(stickDesigerList)
	if length <= 0 {
		return
	}

	index := rand.Intn(length)
	return *stickDesigerList[index]
}

//获取置顶优秀设计师
func GetRecommendRandNewDesigner() (designer GFDesignerPool) {
	stickNewDesigerListLock.RLock()
	defer stickNewDesigerListLock.RUnlock()

	length := len(stickNewDesigerList)
	if length <= 0 {
		return
	}

	index := rand.Intn(length)
	return *stickNewDesigerList[index]
}

//获取置顶衣服
func GetRecommendClothes() (list []Custom) {
	stickPrefectClothesListLock.RLock()
	defer stickPrefectClothesListLock.RUnlock()

	val := len(stickPrefectClothesList)
	if val <= 0 {
		return
	} else if val <= 2 {
		for _, v := range stickPrefectClothesList {
			list = append(list, v.Custom)
		}
		return
	} else {
		index1 := rand.Intn(val)
		tempOne := *stickPrefectClothesList[index1]
		list = append(list, tempOne.Custom)

		for {
			index := rand.Intn(val)
			if index != index1 {
				tempTwo := *stickPrefectClothesList[index]
				list = append(list, tempTwo.Custom)
				break
			}
		}
	}
	return
}
func PublicBuy() (list []*CustomPool) {
	publicBuyListLock.RLock()
	defer publicBuyListLock.RUnlock()
	return publicBuyList
}

func GetStickTopic() (topic TopicInfo) {
	stickTopicListLock.RLock()
	defer stickTopicListLock.RUnlock()

	val := len(stickTopicList)
	if val <= 0 {
		return
	}

	index := rand.Intn(val)

	topic.TopicID = stickTopicList[index].ID
	topic.MainIcon = stickTopicList[index].MainIcon
	topic.SliderIcon = stickTopicList[index].SliderIcon
	topic.TitleIcon = stickTopicList[index].TitleIcon
	return
}

//获取新人置顶衣服
func GetNewRecommendClothes() (list []Custom) {
	stickNewClothesListLock.RLock()
	defer stickNewClothesListLock.RUnlock()

	val := len(stickNewClothesList)
	if val <= 0 {
		return
	} else if val <= 2 {
		for _, v := range stickNewClothesList {
			list = append(list, v.Custom)
		}
		return
	} else {
		index1 := rand.Intn(val)
		tempOne := *stickNewClothesList[index1]
		list = append(list, tempOne.Custom)

		for {
			index := rand.Intn(val)

			if index != index1 {
				tempTwo := *stickNewClothesList[index]
				list = append(list, tempTwo.Custom)
				break
			}
		}
	}

	return
}

// func GetRecommendClothesList() map[string]*CustomPool {
// 	recommendClothesLock.RLock()
// 	defer recommendClothesLock.RUnlock()

// 	return recommendClothesList
// }
// func GetRecommendNewDesigerList() []*GFDesignerPool {
// 	recommendNewDesigerLock.RLock()
// 	defer recommendNewDesigerLock.RUnlock()

// 	return recommendNewDesigerList
// }

// func GetRandLothesList(queryType int) (list []CustomPool) {
// 	var Pool1, Pool2, Pool3, Pool4, Pool5, Pool6 []CustomPool
// 	if queryType == 1 {
// 		customPoolLock.RLock()
// 		defer customPoolLock.RUnlock()

// 		Pool1 = customPool1
// 		Pool2 = customPool2
// 		Pool3 = customPool3
// 		Pool4 = customPool4
// 		Pool5 = customPool5
// 		Pool6 = customPool6
// 	} else if queryType == 2 {
// 		newCustomPoolLock.RLock()
// 		defer newCustomPoolLock.RUnlock()

// 		Pool1 = newCustomPool1
// 		Pool2 = newCustomPool2
// 		Pool3 = newCustomPool3
// 		Pool4 = newCustomPool4
// 		Pool5 = newCustomPool5
// 		Pool6 = newCustomPool6
// 	}

// 	if len(Pool1)+len(Pool2)+len(Pool3)+len(Pool4)+len(Pool5)+len(Pool6) > constants.DESIGNER_ZONE_CUSTOM_SIZE {
// 		list = getRandomClothes(Pool1, Pool2, Pool3, Pool4, Pool5, Pool6, constants.DESIGNER_ZONE_CUSTOM_SIZE)
// 	} else {
// 		tempPool := []CustomPool{}
// 		tempPool = append(tempPool, Pool1...)
// 		tempPool = append(tempPool, Pool2...)
// 		tempPool = append(tempPool, Pool3...)
// 		tempPool = append(tempPool, Pool4...)
// 		tempPool = append(tempPool, Pool5...)
// 		tempPool = append(tempPool, Pool6...)

// 		for _, v := range tempPool {
// 			list = append(list, v)
// 		}
// 	}
// 	return list
// }

func GetRandomTopic() (output []TopicInfo) {
	input := GetTopicList()
	size := len(input)
	if size <= 0 {
		return output
	}
	index := rand.Intn(size)
	model := input[index]
	output = append(output, model)
	return output
}
func GetTopicList() []TopicInfo {
	topicListLock.RLock()
	defer topicListLock.RUnlock()
	list := []TopicInfo{}
	for _, v := range topicList {
		var temp TopicInfo
		temp.TopicID = v.ID
		temp.Title = v.OwerTitle
		temp.TitleIcon = v.TitleIcon
		temp.MainIcon = v.MainIcon
		temp.SliderIcon = v.SliderIcon
		list = append(list, temp)
	}
	return list
}

func GetRecommendBannerList() []BannerInfo {
	bannerListLock.RLock()
	defer bannerListLock.RUnlock()
	list := []BannerInfo{}
	// stickMap := make(map[i])
	for _, v := range bannerList {
		recommendInt, _ := strconv.Atoi(v.IndexOrRecom)
		if recommendInt == constants.RECOMMEND {
			var temp BannerInfo
			temp.BannerType = v.BannerType
			temp.URL = v.IMGURL
			temp.Params = v.URL
			temp.Order = v.Shunxu
			list = append(list, temp)
		}
	}
	return list
}

func GetIndexBannerList() []IndexBannerInfo {
	bannerListLock.RLock()
	defer bannerListLock.RUnlock()
	list := []IndexBannerInfo{}
	for _, v := range bannerList {
		indexInt, _ := strconv.Atoi(v.IndexOrRecom)
		if indexInt == constants.INDEX {
			var temp IndexBannerInfo
			temp.BanenrID = v.ID
			temp.IconPath = v.IMGURL
			temp.Order = v.Shunxu
			temp.Config = v.URL
			list = append(list, temp)
		}
	}
	return list
}
