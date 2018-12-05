package params

import (
	"appcfg"
	"constants"
	"dao"
	"database/sql"
	"db/cart"
	"db/global"
	dp "db/push"
	"encoding/json"
	"fmt"
	. "language"
	. "logger"
	"math/rand"
	. "model"
	p "push"
	"service/clothes"
	sp "service/push"
	"service/pushlist"
	"strconv"
	"strings"
	"sync"
	"time"
	"tools"
	. "types"
	"xlsx"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
)

var channelLimitLock sync.RWMutex
var channelLimit bool
var channelNum map[string]int

var maintainLock sync.RWMutex
var maintainInfo string
var maintainStartTime int64
var maintainEndTime int64
var maintainURL string

var whiteListLock sync.RWMutex
var whiteList []string

var pkDiscountLock sync.RWMutex
var pkDiscount int

var orderOpenLock sync.RWMutex
var orderOpen bool

var monthCardStatusLock sync.RWMutex
var monthCardStatus bool

var giftPackStatusLock sync.RWMutex
var giftPackStatus bool

var copyCustomsLock sync.RWMutex
var copyCustoms map[string]int

var rollingsLock sync.RWMutex
var rollings []RollingAnnoun

var delUserDurationLock sync.RWMutex
var delUserDuration int

var errMailTosLock sync.RWMutex
var errMailTos string

var channelShareURLMapLock sync.RWMutex
var channelShareURLMap map[string]string

var channelForbidMapLock sync.RWMutex
var channelForbidMap map[string]int64

var hotDesignerListLock sync.RWMutex
var hotDesignerList []string

var timeDoorParamLock sync.RWMutex
var timeDoorParam TimeDoor

var blockMapLock sync.RWMutex
var blockMap map[string]BlockInfo

var designerCoinListLock sync.RWMutex
var designerCoinList map[string]int

var banDesignerMapLock sync.RWMutex
var banDesignerMap map[string]int64

var tempSaleEventLock sync.RWMutex
var tempSaleEventStartTime int64
var tempSaleEventEndTime int64
var tempSaleEventUploadBeginTime int64

var saleRankChangeDesignCoinLock sync.RWMutex
var saleRankChangeDesignCoin bool

var monthCardRewardMapLock sync.RWMutex
var monthCardRewardMap map[string]MonthCardReward

var rtPartyDefaultSubjectLock sync.RWMutex
var rtPartyDefaultSubjectList *sll.List

var prizePartyClose bool
var prizePartyCloseLock sync.RWMutex

var logUploadPlayerMap map[string]bool
var logUploadPlayerMapLock sync.RWMutex

var defaultRTSubjects map[string]int

func init() {
	if appcfg.GetServerType() == "" {
		loadRTPartySubjects()
		refreshParams1()
		go paramsLoop()

		if appcfg.GetBool("main_server", false) {
			go pushLoop()
		}
	}
}

func paramsLoop() {
	for {
		time.Sleep(60 * time.Second)

		refreshParams1()
	}
}

func pushLoop() {
	for {
		time.Sleep(60 * time.Second)

		sendPushs()
	}
}

func loadRTPartySubjects() {
	defaultRTSubjects = make(map[string]int)
	//读取默认主题池
	xlFile, err := xlsx.OpenFile("scripts/data/data_salon.xlsx")
	if err != nil {
		panic(err)
	}

	for _, sheet := range xlFile.Sheets {
		if sheet.Name == "Sheet1" {
			for i, row := range sheet.Rows {
				if i > 0 {
					for j, cell := range row.Cells {
						if j == 1 {
							v, _ := cell.String()
							defaultRTSubjects[v] = 1
						}
					}
				}
			}
		}
	}
	Info("default rt subject size", len(defaultRTSubjects))
}

func refreshParams1() {
	var pd, srcdc, ppck int
	var tsest, tseet, tseubt int64
	var emls string
	var gp GameParam
	var res ChannelLimit

	mcrm := make(map[string]MonthCardReward)
	wl := make([]string, 0)
	eml := make([]string, 0)
	limits := make(map[string]int)
	ccm := make(map[string]int)
	csum := make(map[string]string)
	cfbl := make(map[string]int64)
	dbm := make(map[string]int64)
	hdl := make([]string, 0)
	bml := make([]BlockInfo, 0)
	bm := make(map[string]BlockInfo)
	dcl := make(map[string]int)
	rpdsm := make(map[string]int)
	rpdsl := sll.New()
	lupm := make(map[string]bool, 0)

	mst, met, mi, murl, _ := checkMaintain()
	rls, _ := getRollingAnnouncements()

	var tp TimeDoor

	if rows, err := dao.ListConfigStmt.Query(); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var k, v string
			if err := rows.Scan(&k, &v); err != nil {
				Err(err)
				return
			}
			if strings.HasPrefix(k, constants.CHANNEL_SHARE_URL_KEY_PREFIX) {
				csum[k] = v
			} else {
				switch k {
				case constants.KOR_SHIKONGZHIMEN_PARAMS:
					if err := json.Unmarshal([]byte(v), &tp); err != nil {
						Err(err)
					}
				case constants.MAINTAIN_WHITE_LIST_KEY:
					if err = json.Unmarshal([]byte(v), &wl); err != nil {
						Err(err)
					}
				case constants.ERR_MAIL_RECEIVER_LIST_KEY:
					if err = json.Unmarshal([]byte(v), &eml); err != nil {
						Err(err)
					} else {
						for _, s := range eml {
							emls += s + ";"
						}
					}
				case constants.GAME_PARAM_KEY:
					if err = json.Unmarshal([]byte(v), &gp); err != nil {
						Err(err)
					} else {
						pd = gp.PKDiscount
						if pd != 100 {
							Info("pk discount:", pd)
						}
						if pd == 0 {
							pd = 100
						}
					}
				case constants.CHANNEL_LIMIT_KEY:
					if err = json.Unmarshal([]byte(v), &res); err != nil {
						Err(err)
					} else {
						for k, v := range res.LimitInfo {
							limits[k] = v.LimitNum
						}
					}
				case constants.DESIGN_COIN_KEY:
					//获取白名单设计师列表
					if err = json.Unmarshal([]byte(v), &dcl); err != nil {
						Err(err)
					}
				case constants.COPY_CUSTOM_KEY:
					if err = json.Unmarshal([]byte(v), &ccm); err != nil {
						Err(err)
					}
				case constants.CHANNEL_FORBID_LIST_KEY:
					if err = json.Unmarshal([]byte(v), &cfbl); err != nil {
						Err(err)
					}
				case constants.DESIGNER_BAN_MAP_KEY:
					if err = json.Unmarshal([]byte(v), &dbm); err != nil {
						Err(err)
					}
				case constants.HOT_DESIGNER_LIST_KEY:
					if err = json.Unmarshal([]byte(v), &hdl); err != nil {
						Err(err)
					}
				case constants.BLOCK_LIST_KEY:
					if err = json.Unmarshal([]byte(v), &bml); err != nil {
						Err(err)
					} else {
						for _, b := range bml {
							bm[b.Username] = b
						}
					}
				case constants.TEMP_SALE_EVENT_START_TIME_KEY:
					tsest, _ = strconv.ParseInt(v, 10, 64)
				case constants.TEMP_SALE_EVENT_END_TIME_KEY:
					tseet, _ = strconv.ParseInt(v, 10, 64)
				case constants.TEMP_SALE_EVENT_UPLOAD_BEGIN_TIME_KEY:
					tseubt, _ = strconv.ParseInt(v, 10, 64)
				case constants.SALE_RANK_CHANGE_DESIGN_COIN_KEY:
					srcdc, _ = strconv.Atoi(v)
				case constants.MONTH_CARD_REWARD_KEY:
					if err = json.Unmarshal([]byte(v), &mcrm); err != nil {
						Err(err)
					}
				case constants.RT_PARTY_DEFAULT_SUBJECT_KEY:
					if err = json.Unmarshal([]byte(v), &rpdsm); err != nil {
						Err(err)
					}
				case constants.PRIZE_PARTY_CLOSE_KEY:
					ppck, _ = strconv.Atoi(v)
				case constants.LOG_UPLOAD_PLAYER_LIST_KEY:
					if err = json.Unmarshal([]byte(v), &lupm); err != nil {
						Err(err)
					}
				}
			}
		}
	}

	for k, _ := range defaultRTSubjects {
		rpdsm[k] = 1
	}
	for k, _ := range rpdsm {
		rpdsl.Add(k)
	}

	channelLimitLock.Lock()
	channelLimit = res.Limit
	channelNum = limits
	channelLimitLock.Unlock()

	maintainLock.Lock()
	maintainStartTime = mst
	maintainEndTime = met
	maintainInfo = mi
	maintainURL = murl
	maintainLock.Unlock()

	whiteListLock.Lock()
	whiteList = wl
	whiteListLock.Unlock()

	pkDiscountLock.Lock()
	pkDiscount = pd
	pkDiscountLock.Unlock()

	orderOpenLock.Lock()
	orderOpen = gp.OrderStatus
	orderOpenLock.Unlock()

	monthCardStatusLock.Lock()
	monthCardStatus = gp.MonthCardStatus
	monthCardStatusLock.Unlock()

	giftPackStatusLock.Lock()
	giftPackStatus = gp.GiftPackStatus
	giftPackStatusLock.Unlock()

	copyCustomsLock.Lock()
	copyCustoms = ccm
	if len(copyCustoms) > 0 {
		Info("copy customs size:", len(copyCustoms))
	}
	copyCustomsLock.Unlock()

	rollingsLock.Lock()
	rollings = rls
	if len(rollings) > 0 {
		Info("rolling size:", len(rollings))
	}
	rollingsLock.Unlock()

	delUserDurationLock.Lock()
	delUserDuration = gp.DelUserDuration
	if delUserDuration == 0 {
		delUserDuration = constants.DEFAULT_DEL_USER_DURATION
	}
	delUserDurationLock.Unlock()

	errMailTosLock.Lock()
	errMailTos = emls
	Info("err mail send to:", emls)
	errMailTosLock.Unlock()

	channelShareURLMapLock.Lock()
	channelShareURLMap = csum
	channelShareURLMapLock.Unlock()

	channelForbidMapLock.Lock()
	channelForbidMap = cfbl
	if len(channelForbidMap) > 0 {
		Info("channel forbid:", channelForbidMap)
	}
	channelForbidMapLock.Unlock()

	banDesignerMapLock.Lock()
	banDesignerMap = dbm
	if len(banDesignerMap) > 0 {
		Info("designer ban:", banDesignerMap)
	}
	banDesignerMapLock.Unlock()

	blockMapLock.Lock()
	blockMap = bm
	blockMapLock.Unlock()
	if len(bm) > 0 {
		Info("block map:", bm)
	}

	hotDesignerListLock.Lock()
	hotDesignerList = hdl
	hotDesignerListLock.Unlock()

	timeDoorParamLock.Lock()
	timeDoorParam = tp
	timeDoorParamLock.Unlock()
	Info("time door:", tp.Kor_shikongzhimen_start_time, tp.Kor_shikongzhimen_end_time, len(tp.LotteryClothesShikongLevel13), len(tp.LotteryClothesShikongLevel4), len(tp.LotteryClothesShikongLevel5))

	tempSaleEventLock.Lock()
	tempSaleEventStartTime = tsest
	tempSaleEventEndTime = tseet
	tempSaleEventUploadBeginTime = tseubt
	tempSaleEventLock.Unlock()

	designerCoinListLock.Lock()
	designerCoinList = dcl
	designerCoinListLock.Unlock()

	saleRankChangeDesignCoinLock.Lock()
	saleRankChangeDesignCoin = srcdc == 1
	saleRankChangeDesignCoinLock.Unlock()

	monthCardRewardMapLock.Lock()
	monthCardRewardMap = mcrm
	monthCardRewardMapLock.Unlock()

	rtPartyDefaultSubjectLock.Lock()
	rtPartyDefaultSubjectList = rpdsl
	rtPartyDefaultSubjectLock.Unlock()

	prizePartyCloseLock.Lock()
	prizePartyClose = ppck == 1
	prizePartyCloseLock.Unlock()

	logUploadPlayerMapLock.Lock()
	logUploadPlayerMap = lupm
	logUploadPlayerMapLock.Unlock()
}

func PlayerShouldUploadLog(username string) bool {
	logUploadPlayerMapLock.RLock()
	defer logUploadPlayerMapLock.RUnlock()

	return logUploadPlayerMap[username]
}

func PrizePartyClose() bool {
	prizePartyCloseLock.RLock()
	defer prizePartyCloseLock.RUnlock()

	return prizePartyClose
}

func GetRTPartyDefaultSubjects() *sll.List {
	rtPartyDefaultSubjectLock.RLock()
	defer rtPartyDefaultSubjectLock.RUnlock()

	nl := sll.New()
	for _, v := range rtPartyDefaultSubjectList.Values() {
		nl.Add(v)
	}
	return nl
}

func GetMonthCardReward(month string) (r MonthCardReward) {
	monthCardRewardMapLock.RLock()
	defer monthCardRewardMapLock.RUnlock()

	if v, ok := monthCardRewardMap[month]; ok {
		r = v
	}
	return
}

func SaleRankChangeDesignCoin() bool {
	saleRankChangeDesignCoinLock.RLock()
	defer saleRankChangeDesignCoinLock.RUnlock()

	return saleRankChangeDesignCoin
}

func GetTempSaleEventTime() (int64, int64, int64) {
	tempSaleEventLock.RLock()
	defer tempSaleEventLock.RUnlock()

	return tempSaleEventStartTime, tempSaleEventEndTime, tempSaleEventUploadBeginTime
}

func GetHotDesignerList() []string {
	hotDesignerListLock.RLock()
	defer hotDesignerListLock.RUnlock()

	return hotDesignerList
}

//getRandomLevelShikong 获取动态概率
func getRandomLevelShikong(force5 bool, chance1 int, chance2 int) int {
	if force5 {
		return 5
	}
	randValue := rand.Intn(100)
	if randValue <= chance1 {
		return rand.Intn(3) + 1
	} else if randValue <= chance2 {
		return 4
	} else {
		return 5
	}
}

func GetTimeDoor() TimeDoor {
	timeDoorParamLock.RLock()
	defer timeDoorParamLock.RUnlock()

	return timeDoorParam
}

//GetRandomLotteryShikongClothes 获取时空之门的衣服
func GetRandomLotteryTimeDoorClothes(force5 bool, tp *TimeDoor) (string, int, string) {
	level := getRandomLevelShikong(force5, tp.Chance13, tp.Chance4)
	var tmpClothes []string
	if level <= 3 && len(tp.LotteryClothesShikongLevel13) > 0 {
		tmpClothes = tp.LotteryClothesShikongLevel13
	} else if level == 4 && len(tp.LotteryClothesShikongLevel4) > 0 {
		tmpClothes = tp.LotteryClothesShikongLevel4
	} else if level == 4 && len(tp.LotteryClothesShikongLevel4) <= 0 {
		tmpClothes = tp.LotteryClothesShikongLevel13
	} else if level == 5 && len(tp.LotteryClothesShikongLevel5) > 0 {
		tmpClothes = tp.LotteryClothesShikongLevel5
	}
	lotteryClothesCount := len(tmpClothes)

	if lotteryClothesCount != 0 {
		index := rand.Intn(lotteryClothesCount)
		return clothes.LotteryPoolActivity[tmpClothes[index]].CId, clothes.LotteryPoolActivity[tmpClothes[index]].Star, clothes.LotteryPoolActivity[tmpClothes[index]].Name
	} else {
		return "", 0, ""
	}
}

func GetDesignerBanEndTime(du string) int64 {
	banDesignerMapLock.RLock()
	defer banDesignerMapLock.RUnlock()

	return banDesignerMap[du]
}

func Blocked(username string) (bool, string) {
	blockMapLock.RLock()
	defer blockMapLock.RUnlock()

	if v, ok := blockMap[username]; ok {
		return true, v.Reason
	}
	return false, ""
}

//GetDesignCornList 获取白名单设计师列表
func InDesignCornList(username string) bool {
	designerCoinListLock.RLock()
	defer designerCoinListLock.RUnlock()

	if _, exist := designerCoinList[username]; exist {
		return true
	}
	return false
}

func GetDesignCoinDesignerList() []string {
	designerCoinListLock.RLock()
	defer designerCoinListLock.RUnlock()

	res := make([]string, 0)
	for k, _ := range designerCoinList {
		res = append(res, k)
	}

	return res
}

func ChannelForbided(channel string, uid int64) bool {
	channelForbidMapLock.RLock()
	defer channelForbidMapLock.RUnlock()

	if v, ok := channelForbidMap[channel]; ok {
		return uid > v
	} else {
		return false
	}
}

func GetChannelShareUrl(channel string) string {
	if channel == "" {
		return ""
	} else {
		channelShareURLMapLock.RLock()
		defer channelShareURLMapLock.RUnlock()
		return channelShareURLMap[constants.CHANNEL_SHARE_URL_KEY_PREFIX+channel]
	}
}

func GetErrMailTos() string {
	errMailTosLock.RLock()
	defer errMailTosLock.RUnlock()

	return errMailTos
}

func GetDelUserDuration() int {
	delUserDurationLock.RLock()
	defer delUserDurationLock.RUnlock()

	return delUserDuration
}

func GetRollings() []RollingAnnoun {
	rollingsLock.RLock()
	defer rollingsLock.RUnlock()

	return rollings
}

func GetConfig(key string) string {
	var value string
	if err := dao.GetConfigStmt.QueryRow(key).Scan(&value); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("no config key:", key)
		} else {
			Err(err)
		}
	}
	return value
}

func SetConfig(key, value string) {
	if _, err := dao.SetConfigStmt.Exec(key, value, value); err != nil {
		Err(err)
	}
	return
}
func CustomCopy(cid string) (ok bool, rd int) {
	copyCustomsLock.RLock()
	defer copyCustomsLock.RUnlock()

	rd, ok = copyCustoms[cid]
	return
}

func ChannelLimited() bool {
	channelLimitLock.RLock()
	defer channelLimitLock.RUnlock()

	return channelLimit
}

func ChannelLimitNum(channel string) int {
	channelLimitLock.RLock()
	defer channelLimitLock.RUnlock()

	return channelNum[channel]
}

func OrderOpen() bool {
	orderOpenLock.RLock()
	defer orderOpenLock.RUnlock()

	return orderOpen
}

func MonthCardCanBuy() bool {
	monthCardStatusLock.RLock()
	defer monthCardStatusLock.RUnlock()

	return monthCardStatus
}

func GiftPackOpen() bool {
	giftPackStatusLock.RLock()
	defer giftPackStatusLock.RUnlock()

	return giftPackStatus
}

func GetPKDiscount() int {
	pkDiscountLock.RLock()
	defer pkDiscountLock.RUnlock()

	return pkDiscount
}

func IsMaintainStatus() (bool, string, string) {
	maintainLock.RLock()
	defer maintainLock.RUnlock()

	now := time.Now().Unix()
	return maintainStartTime <= now && maintainEndTime >= now, maintainInfo, maintainURL
}

func IsInWhiteList(username string) bool {
	whiteListLock.RLock()
	defer whiteListLock.RUnlock()

	for _, v := range whiteList {
		if v == username {
			return true
		}
	}

	return false
}

func checkMaintain() (st int64, et int64, info string, url string, err error) {
	now := time.Now().Unix()

	if err = dao.GetMaintainJobStmt.QueryRow(now, now).Scan(&info, &st, &et, &url); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("no maintain job")
		} else {
			Err(err)
		}
	}

	return
}

func getRollingAnnouncements() (res []RollingAnnoun, err error) {
	now := time.Now().Unix()

	var rows *sql.Rows
	if rows, err = dao.GetRollingAnnouncementStmt.Query(now, now); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]RollingAnnoun, 0)
		for rows.Next() {
			var r RollingAnnoun
			if err = rows.Scan(&r.ID, &r.Content, &r.StartTime, &r.EndTime); err != nil {
				Err(err)
				return
			}

			res = append(res, r)
		}
	}
	return
}

func sendPushs() {
	nowt := time.Now()

	//gm配置推送
	loopSendGMPush(nowt)

	if appcfg.GetLanguage() == "" {
		//新品发布推送
		// loopSendNewCustomPush(nowt)

		//购物车推送
		loopSendShopCartInventoryPush(nowt)
	}

	//设计师库存紧张推送
	// loopSendDesignerInventoryEmptyPush(nowt)
}

func loopSendGMPush(nowt time.Time) {
	now := nowt.Unix()
	if ps, err := dp.GetLegalPushJobs(); err == nil {
		for _, pj := range ps {
			if pj.PushTime > now {
				continue
			}
			if success, _ := dp.DonePush(pj.ID); success {
				if pj.To == "ALL" {
					p.TPushToAll(pj.Title, pj.Content)
				} else if pj.To == "ANDROID" {
					p.TPushToAllAndroid(pj.Title, pj.Content)
				} else if pj.To == "IOS" {
					p.TPushToAllIOS(pj.Title, pj.Content)
				} else {
					tos := strings.Split(pj.To, ",")
					for _, t := range tos {
						if appcfg.GetLanguage() == "" {
							p.TPushToOne(t, pj.Title, pj.Content)
						} else if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
							sp.PushKorTC(t, pj.Title, pj.Content)
						}
					}
				}
			}
		}
	}
}

func loopSendShopCartInventoryPush(nowt time.Time) {
	nowHour := nowt.Hour()

	jobToken := fmt.Sprintf("loopSendShopCartInventoryPush_%s", nowt.Format("20060102"))

	if !appcfg.GetBool("test_shop_cart_push", false) {
		if nowHour != 18 {
			Info("not time for loopSendShopCartInventoryPush")
			return
		}

		if global.CheckOnceJob(jobToken) {
			Info("loopSendShopCartInventoryPush done")
			return
		} else {
			if !global.DoOnceJob(jobToken) {
				Info("loopSendShopCartInventoryPush done fail")
				return
			}
		}
	}

	cloPage := 1
	userJobMap := make(map[string]int)

	for {
		if clos, err := cart.ListSoldClothes(cloPage, constants.SHOP_CART_PUSH_BATCH_SIZE); err == nil && len(clos) > 0 {
			for _, c := range clos {
				if !tools.ClothesIDLegal(c.ClothesIDInGame) {
					Err("cloid illegal", c.ClothesIDInGame)
					continue
				}

				userPage := 1
				for {
					if us, err := cart.ListCartUserByPage(c.ClothesIDInGame, userPage, constants.SHOP_CART_PUSH_BATCH_SIZE); err == nil && len(us) > 0 {
						for _, u := range us {
							userJobToken := fmt.Sprintf("sci_push_%s_%d", u, nowHour)
							if userJobMap[userJobToken] <= 0 {
								userJobMap[userJobToken] = 1
								pushlist.AddSinglePushJob(u, L("o15"), fmt.Sprintf(L("o16"), c.Name), constants.PUSH_TYPE_SHOP_CART)
							}
						}

						userPage++
					} else {
						break
					}

				}
			}
		} else {
			break
		}

		cloPage++
	}
}

func UpCustomPushShopCart(designerUsername string, clos []*SimpleCustom) {
	designerJobToken := fmt.Sprintf("duc_push_%s", designerUsername)
	if global.CheckOnceJob(designerJobToken) && !appcfg.GetBool("test_shop_cart_push", false) {
		Info("designer up custom shop cart push done", designerUsername)
		return
	}

	if !global.DoOnceJob(designerJobToken) {
		Info("designer up custom shop cart push done 1", designerUsername)
		return
	}

	if len(clos) > 0 {
		for _, c := range clos {
			if !tools.ClothesIDLegal(c.ClothesIDInGame) {
				Err("cloid illegal", c.ClothesIDInGame)
				continue
			}

			userPage := 1
			for {
				if us, err := cart.ListCartUserByPage(c.ClothesIDInGame, userPage, constants.SHOP_CART_PUSH_BATCH_SIZE); err == nil && len(us) > 0 {
					for _, u := range us {
						if u == designerUsername {
							continue
						}

						jobToken := fmt.Sprintf("uc_push_%s", u)
						if global.DoOnceJob(jobToken) || appcfg.GetBool("test_shop_cart_push", false) {
							pushlist.AddSinglePushJob(u, L("o17"), fmt.Sprintf(L("o18"), c.Name), constants.PUSH_TYPE_SHOP_CART)
						}
					}

					userPage++
				} else {
					break
				}
			}
		}
	}
}
