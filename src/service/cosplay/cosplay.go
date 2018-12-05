package cosplay

import (
	// "encoding/json"
	"fmt"
	. "logger"
	"shutdown"
	"sync"
	"time"
)

import (
	"appcfg"
	"db/cosplay"
	"db/cosplay/items"
	"service/score"
	// "timer"
	. "types"
)

var (
	openCosplay     []Cosplay
	openCosplayLock sync.RWMutex
	// openCosReset    bool

	closedCosplay     []Cosplay
	closedCosplayLock sync.RWMutex
	// closedCosReset    bool
)

func init() {
	if appcfg.GetServerType() != "" {
		return
	}

	// openCosReset = true

	if appcfg.GetBool("main_server", false) {
		// 检查Cosplay是否结束
		go checkCosplay()
		// 每隔一分钟，移除打头的置顶
		// go removeTopItem() 目前的逻辑，暂时不需要咯
		// 每隔一分钟，刷新一遍topItem
		// go refreshTopItem()
	}
	// 每隔一分钟，刷新一遍cosplay活动列表
	// go refreshCosplay()
	// 每隔一分钟，刷新一遍结束的cosplay活动列表
	// go refreshClosedCosplay()
	// // 每隔一分钟，刷新一遍CosplayItem列表
	// go refreshCosItemsMap()
	// // 每隔一分钟，刷新一遍Cosplay活动的rank
	// go refreshCosItemRank()
	// // 每个一分钟，刷新一次cosComment
	// go refreshCosComment()
}

func checkCosplay() {
	for {
		checkCosplay1()
		time.Sleep(600 * time.Second)
	}
}

func checkCosplay1() {
	shutdown.AddOneRequest("checkCosplay")
	defer shutdown.DoneOneRequest("checkCosplay")

	if coses, err := cosplay.FindCosplayByStatus(COS_PROCESSING); err == nil {
		for _, cos := range coses {
			if err := items.ProcessCos(cos.CosplayID, cos.Title, cos.Type); err == nil {
				cosplay.CloseCosplay(cos.Type, cos.CosplayID)
			}
		}
	}
}

// func refreshCosplay() {
// 	for {
// 		if cosp, err := cosplay.FindCosplayByStatus(COS_OPEN); err == nil {
// 			openCosplayLock.Lock()
// 			openCosplay = cosp
// 			openCosplayLock.Unlock()
// 		}

// 		if cosp, err := cosplay.FindCosplayByStatus(COS_CLOSED); err == nil {
// 			closedCosplayLock.Lock()
// 			closedCosplay = cosp
// 			closedCosplayLock.Unlock()
// 		}

// 		time.Sleep(60 * time.Second)

// openCosReset = true
// closedCosReset = true
// }

// t := make(chan int32, 1)
// timer.Add(int32(TIMER_REFRESH_COSPLAY), time.Now().Unix()+60, t)
// for {
// 	select {
// 	case <-t:
// 		openCosReset = true
// 		timer.Add(int32(TIMER_REFRESH_COSPLAY), time.Now().Unix()+60, t)
// 	}
// }
// }

// func refreshClosedCosplay() {
// for{
// 	time.Sleep(65 * time.Second)

// }

// t := make(chan int32, 1)
// timer.Add(int32(TIMER_REFRESH_FINISH_COSPLAY), time.Now().Unix()+65, t)
// for {
// 	select {
// 	case <-t:
// 		closedCosReset = true
// 		timer.Add(int32(TIMER_REFRESH_FINISH_COSPLAY), time.Now().Unix()+60, t)
// 	}
// }
// }

// func refreshTopItem() {
// 	i := 0
// 	for {
// 		openCosplayLock.RLock()
// 		for _, v := range openCosplay {
// 			if i%12 == 0 {
// 				items.PopTopCosItem(v.CosplayID)
// 			} else {
// 				items.CheckTopFrequently(v.CosplayID)
// 			}

// 		}
// 		openCosplayLock.RUnlock()

// 		if i%12 == 0 {
// 			i = 0
// 		}

// 		i++
// 		time.Sleep(5 * time.Second)
// 	}
// }

func GetOpenCosplays(t time.Time) (coses []Cosplay, err error) {
	coses = make([]Cosplay, 0)

	openCosplayLock.RLock()
	defer openCosplayLock.RUnlock()

	for _, c := range openCosplay {
		if c.OpenTime <= t.Unix() && c.CloseTime > t.Unix() {
			coses = append(coses, c)
		}
	}

	return
}

func GetFinishedCosplays() (coses []Cosplay, err error) {
	closedCosplayLock.RLock()
	defer closedCosplayLock.RUnlock()

	coses = closedCosplay

	return
}

func SaveCosplayItem(username string, nickname string, cosplayID int64, title string, desc string, modelNo int, icon string, head string, cosClothes map[string]string, cosBg string, img string, dyeMap map[string][7][4]float64) (int, error) {
	cos, have := getCosplayByID(cosplayID)
	if !have {
		return 0, fmt.Errorf("cannot find cosplay with ID:%d", cosplayID)
	}

	var sysScore int64
	sysScore = 0
	// 计算系统分数
	sysScore = score.GetCosSysScore(cosClothes, cos.Params)

	item := CosItem{
		Username:   username,
		Nickname:   nickname,
		CosplayID:  cosplayID,
		Title:      title,
		Desc:       desc,
		ModelNo:    modelNo,
		Icon:       icon,
		Head:       head,
		Clothes:    cosClothes,
		CosBg:      cosBg,
		SysScore:   sysScore,
		Score:      sysScore,
		Scores:     make([]CosItemScore, 0),
		UploadTime: time.Now().Unix(),
		Img:        img,
		DyeMap:     dyeMap,
	}
	return items.Add(&item)
}

func GetCosItems(cosplayID int64, pageID int, status bool, pageCount int) (retItems []CosItem, err error) {
	// var ret []CosItemR
	return items.FindByPageAndCountByCosplayID(cosplayID, pageID, pageCount)

	// var tmp CosItem
	// for _, v := range ret {
	// 	tmp = CosItem{
	// 		Username:  v.Username,
	// 		UID:       v.UID,
	// 		Nickname:  v.Nickname,
	// 		CosplayID: v.CosplayID,
	// 		ItemID:    v.ItemID,
	// 		Title:     v.Title,
	// 		Desc:      v.Desc,
	// 		ModelNo:   v.ModelNo,
	// 		Score:     int64(v.Score),
	// 		SysScore:  v.SysScore,
	// 		Icon:      v.Icon,
	// 		Head:      v.Head,
	// 		CosBg:     v.CosBg,
	// 		CosPoints: v.CosPoints,

	// 		UploadTime: v.UploadTime,

	// 		Rank: v.Rank,
	// 		Img:  v.Img,
	// 	}
	// 	json.Unmarshal([]byte(v.Clothes[1:len(v.Clothes)-1]), &tmp.Clothes)
	// 	json.Unmarshal([]byte(v.Clothes[1:len(v.DyeMap)-1]), &tmp.DyeMap)

	// 	retItems = append(retItems, tmp)
	// }

	// return
}

func GetLeaderCos(cosplayID int64, pageID int, pageCount int) (result []CosItem) {
	items, err := items.FindByPageAndCountOrderByScore(cosplayID, pageID, pageCount)
	if err == nil {
		// var tmp CosItem
		// for _, v := range itemsR {
		// 	tmp = CosItem{
		// 		Username:   v.Username,
		// 		UID:        v.UID,
		// 		Nickname:   v.Nickname,
		// 		CosplayID:  v.CosplayID,
		// 		ItemID:     v.ItemID,
		// 		Title:      v.Title,
		// 		Desc:       v.Desc,
		// 		ModelNo:    v.ModelNo,
		// 		SysScore:   v.SysScore,
		// 		Icon:       v.Icon,
		// 		Head:       v.Head,
		// 		CosBg:      v.CosBg,
		// 		CosPoints:  v.CosPoints,
		// 		UploadTime: v.UploadTime,
		// 		Score:      int64(v.Score),
		// 		Rank:       v.Rank,
		// 		Img:        v.Img,
		// 	}
		// 	json.Unmarshal([]byte(v.Clothes[1:len(v.Clothes)-1]), &tmp.Clothes)
		// 	json.Unmarshal([]byte(v.Clothes[1:len(v.DyeMap)-1]), &tmp.DyeMap)
		// 	result = append(result, tmp)
		// }
		return items
	}

	return
}

func AddTopCosItem(cosItem CosItem) (int, error) {
	return items.AddTopCosItem(cosItem, nil)
}

func GetTopCosItem(cosplayID int64) (ret CosItem) {
	ret, _ = items.GetTopCosItem(cosplayID)
	return
}

func GetRankOneItem(cosplayID int64) (cosItem CosItem, have bool) {
	if list, err := items.FindByPageAndCountOrderByScore(cosplayID, 1, 1); err != nil {
		have = false
		return
	} else {
		if len(list) > 0 {
			// have = true
			// cosItem = CosItem{
			// 	Username:  list[0].Username,
			// 	UID:       list[0].UID,
			// 	Nickname:  list[0].Nickname,
			// 	CosplayID: list[0].CosplayID,
			// 	ItemID:    list[0].ItemID,
			// 	Title:     list[0].Title,
			// 	Desc:      list[0].Desc,
			// 	ModelNo:   list[0].ModelNo,
			// 	Score:     int64(list[0].Score),
			// 	SysScore:  list[0].SysScore,
			// 	Icon:      list[0].Icon,
			// 	Head:      list[0].Head,
			// 	CosBg:     list[0].CosBg,
			// 	CosPoints: list[0].CosPoints,

			// 	UploadTime: list[0].UploadTime,

			// 	Rank: list[0].Rank,
			// 	Img:  list[0].Img,
			// }
			// json.Unmarshal([]byte(list[0].Clothes[1:len(list[0].Clothes)-1]), &cosItem.Clothes)
			// json.Unmarshal([]byte(list[0].Clothes[1:len(list[0].DyeMap)-1]), &cosItem.DyeMap)
			return list[0], true
		} else {
			have = false
		}

		return
	}
}

func GetMyCosItemCntByCosplayID(username string, cosplayID int64) (cnt int, err error) {
	var itemlist []CosItem
	if itemlist, err = items.FindMyCosItemsByCosplayID(username, cosplayID); err != nil {
		return
	} else {
		cnt = len(itemlist)
	}
	return
}

func GetMyCosItemsByCosplayID(username string, cosplayID int64) (retItems []CosItem, err error) {
	// var itemsR []CosItemR
	return items.FindMyCosItemsByCosplayID(username, cosplayID)

	// for _, v := range itemsR {
	// 	if v.Username == username {
	// 		var tmp CosItem
	// 		tmp = CosItem{
	// 			Username:  v.Username,
	// 			UID:       v.UID,
	// 			Nickname:  v.Nickname,
	// 			CosplayID: v.CosplayID,
	// 			ItemID:    v.ItemID,
	// 			Title:     v.Title,
	// 			Desc:      v.Desc,
	// 			ModelNo:   v.ModelNo,
	// 			Score:     int64(v.Score),
	// 			SysScore:  v.SysScore,
	// 			Icon:      v.Icon,
	// 			Head:      v.Head,
	// 			CosBg:     v.CosBg,
	// 			CosPoints: v.CosPoints,

	// 			UploadTime: v.UploadTime,

	// 			Rank: v.Rank,
	// 			Img:  v.Img,
	// 		}
	// 		json.Unmarshal([]byte(v.Clothes[1:len(v.Clothes)-1]), &tmp.Clothes)
	// 		json.Unmarshal([]byte(v.Clothes[1:len(v.DyeMap)-1]), &tmp.DyeMap)

	// 		retItems = append(retItems, &tmp)
	// 	}
	// }

	// return
}

func AddScore(cosplayID int64, item CosItem, username string, addedScore int) (ret CosItem, err error) {
	cosplay, have := getCosplayByID(cosplayID)
	if !have {
		Err("cosplay id:", cosplayID, "not found")
		err = fmt.Errorf("cannot find cosplay with ID:%d", cosplayID)
		return
	}

	var totalScore int64
	var ps int64

	if ps, err = score.GetCosPlayerScore(item, cosplay, username, addedScore); err != nil {
		return
	}

	totalScore = item.SysScore + ps

	if err = items.AddScore(cosplayID, item.ItemID, username, addedScore, totalScore); err != nil {
		return
	}

	ret = item
	ret.Score = totalScore

	return
}

func IsCosplayOpen(cosplayID int64) (bool, string) {
	openCosplayLock.RLock()
	defer openCosplayLock.RUnlock()

	now := time.Now().Unix()
	for _, v := range openCosplay {
		if v.CosplayID == cosplayID && v.CloseTime > now {
			return true, v.Title
		}
	}
	return false, ""
}

func AddCosComment(cosplayID int64, username string, nickname string, cosItemID int64, content string) (retItem CosComment, err error) {
	retItem = CosComment{
		CosplayID: cosplayID,
		CosItemID: cosItemID,
		Username:  username,
		Nickname:  nickname,
		Content:   content,
		ReplyTime: time.Now().Unix(),
	}

	err = cosplay.AddComment(retItem, nil, true)

	return
}

func GetCosComment(cosplayID int64, itemID int64, page int, pageCount int) (comments []*CosComment, err error) {
	return cosplay.GetCosComment(cosplayID, itemID, page, pageCount)
}

func getCosplayByID(cosplayID int64) (cos Cosplay, have bool) {
	have = false
	openCosplayLock.RLock()
	for _, v := range openCosplay {
		if v.CosplayID == cosplayID {
			cos = v
			have = true
			break
		}
	}
	openCosplayLock.RUnlock()

	if !have {
		closedCosplayLock.RLock()
		for _, v := range closedCosplay {
			if v.CosplayID == cosplayID {
				cos = v
				have = true
				break
			}
		}
		closedCosplayLock.RUnlock()
	}
	return
}
