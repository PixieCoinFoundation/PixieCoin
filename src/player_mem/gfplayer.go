package player_mem

import (
	"appcfg"
	"common"
	"constants"
	"db/event"
	"db/global"
	"db/mail"
	"db_pixie/status"
	"encoding/json"
	"errors"
	"fmt"
	"service/kakao"
	"service/params"
	"strconv"
	"strings"
	"sync"
	"time"
	"tools"

	. "language"
	. "logger"
	"service/clothes"

	. "pixie_contract/api_specification"
	servCosItem "service/cosplay"
	"service/iap"
	"service/mails"
	"service/rank"
	"service/reward"
	"service/tasks"
	. "types"
)

// 在线玩家的信息保存
type GFPlayer struct {
	status Status // 玩家的状态

	ThirdChannel      string
	ThirdUsername     string //第三方 sdk用户名
	GameThirdUsername string

	OldChannel string //为了兼容

	Username      string // 用户名
	Nickname      string // 昵称
	AccessToken   string //
	UID           int64  // uid,用于sharding
	DeviceToken   string // 设备的push标识
	HeartbeatTime int64  // 上一次操作时间

	Platform string // 登陆的终端

	playerLock         sync.Mutex
	LastApplyGuildTime int64

	CustomCnt int

	//guild info
	GuildOwner   bool
	GuildInfo    GuildMember
	NotifiedUser map[string]int
	ShopCart     map[string]*ShopCartDetail //服装ID：最后加入购物车的时间&数量

	LastSendETHVerifyTime int64

	LastVistUsername string
}

var (
	startGold    int64
	startPxc     float64
	startDiamond int64
	startHead    string
	startItems   []byte
)

func init() {
	if appcfg.GetServerType() != "" {
		return
	}

	startHead = appcfg.GetString("start_head", "head-101.png")
	startItemss := appcfg.GetString("start_items", "[{\"itemID\":101,\"count\":1,\"maxCount\":99,\"priceType\":2,\"price\":30},{\"itemID\":102,\"count\":1,\"maxCount\":99,\"priceType\":2,\"price\":10},{\"itemID\":103,\"count\":1,\"maxCount\":99,\"priceType\":2,\"price\":20},{\"itemID\":104,\"count\":1,\"maxCount\":99,\"priceType\":2,\"price\":50},{\"itemID\":105,\"count\":1,\"maxCount\":99,\"priceType\":2,\"price\":25},{\"itemID\":106,\"count\":1,\"maxCount\":99,\"priceType\":2,\"price\":10}]")
	startItems = []byte(startItemss)

	startGold = appcfg.GetInt64("start_gold", 20000)
	startPxc = appcfg.GetFloat("start_pxc", 20000.0)

	startDiamond = appcfg.GetInt64("start_diamond", 100)
}

func (self *GFPlayer) Lock() {
	self.playerLock.Lock()
}

func (self *GFPlayer) Unlock() {
	self.playerLock.Unlock()
}

func (self *GFPlayer) InitGuild(setOwner bool, setData bool, wt string) (success bool) {
	// //player guild data
	// if mInfo, err := guild.GetUserGuild(self.Username, wt); err != nil {
	// 	success = false
	// 	return
	// } else {
	// 	self.GuildInfo = mInfo

	// 	if mInfo.GID > 0 {
	// 		if setOwner {
	// 			if self.GuildOwner, err = guild.IsOwner(mInfo.GID, self.Username); err != nil {
	// 				success = false
	// 				return
	// 			}
	// 		}

	// 		if setData {
	// 			if self.Nickname != mInfo.Nickname || self.status.VIP != mInfo.VIP || self.status.Head != mInfo.Head || self.status.Level != mInfo.Level {
	// 				//update
	// 				if err := guild.UpdateMember(mInfo.GID, self.Username, self.status.Head, self.status.VIP, self.status.Level); err != nil {
	// 					success = false
	// 					return
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// success = true
	return true
}

func (self *GFPlayer) GetNpcBuff() NpcBuffResult {
	return self.status.NpcBuff
}

func (self *GFPlayer) GetNpcFriendship() map[string]*FriendshipDetail {
	return self.status.NpcFriendship
}

func (self *GFPlayer) GetNickname() string {
	return self.status.Nickname
}

func (self *GFPlayer) GetStreetName() string {
	return self.status.StreetName
}

// func (self *GFPlayer) CurrencyEnough(moneyType int, moneyCount int64) bool {
// 	if moneyType == constants.PIXIE_GOLD_TYPE {
// 		return self.status.Money >= moneyCount
// 	} else if moneyType == constants.PIXIE_PXC_TYPE {
// 		return self.status.Pxc >= moneyCount
// 	}
// 	return false
// }

func (self *GFPlayer) GoldEnough(moneyCount int64) bool {
	return self.status.Money >= moneyCount
}

func (self *GFPlayer) PxcEnough(moneyCount float64) bool {
	return self.status.Pxc >= moneyCount
}

func (self *GFPlayer) GetHomeShow() HomeShowDetail {
	return self.status.HomeShow
}

func (self *GFPlayer) GetSuitMap() map[string][]Suit {
	return self.status.SuitMap
}

func (self *GFPlayer) GetSuitCountByModelType(modelType string) int {
	return len(self.status.SuitMap[modelType])
}

func (self *GFPlayer) ExistSuit(modelType string, suitMD5 string) (exist bool, suitName string) {
	for _, v := range self.status.SuitMap[modelType] {
		if v.SuitMD5 == suitMD5 {
			exist = true
			suitName = v.SuitName
			return
		}
	}
	return
}

func (self *GFPlayer) GetGuildConfig() GuildMember {
	return self.GuildInfo
}

func (self *GFPlayer) CheckFormula(f *Formula) bool {
	if f != nil {
		if !self.HaveClothes(f.OriginClothes) {
			return false
		}

		if self.status.ExtraInfo1.Cloth < f.ClothCnt {
			return false
		}

		if self.status.ExtraInfo1.Item1 < f.Item1 {
			return false
		}

		if self.status.ExtraInfo1.Needle < f.Item2 {
			return false
		}

		if self.status.ExtraInfo1.Wire < f.Item3 {
			return false
		}

		if !self.CheckMoney(f.Gold) {
			return false
		}

		return true
	}
	return false
}

func (self *GFPlayer) IsGuildOwner() bool {
	return self.GuildOwner
}

func (self *GFPlayer) GetHead() string {
	return self.status.Head
}

func (self *GFPlayer) GetSex() int {
	return self.status.Sex
}

func (self *GFPlayer) GetVerifyBenefit() int {
	return self.status.ExtraInfo1.VerifyBenefit
}

func (self *GFPlayer) GetDayVerifyNum(key string) int {
	if self.status.DayExtraInfo.DayVerifyNum == nil {
		self.status.DayExtraInfo.DayVerifyNum = make(map[string]int)
	}
	return self.status.DayExtraInfo.DayVerifyNum[key]
}

func (self *GFPlayer) GetVerifyCount() int {
	return self.status.DayExtraInfo.VerifyCount
}

func (self *GFPlayer) GetYesterdayBenefit() int {
	return self.status.ExtraInfo1.YesterdayBenefit
}

func (self *GFPlayer) GetOfficialBoyGirlClothesCnt() (bcnt int, gcnt int) {
	for cid, clo := range self.status.Clothes {
		if clo.Count <= 0 {
			continue
		}
		c := clothes.GetClothesById(cid)
		if c != nil {
			if c.Char > 200 {
				bcnt += 1
			} else if c.Char > 100 {
				gcnt += 1
			}
		}
	}
	return
}

func (self *GFPlayer) GetAllBoyGirlClothesCnt() (dbcnt int, dgcnt int, obcnt int, ogcnt int) {
	for cid, clo := range self.status.Clothes {
		if clo.Count <= 0 {
			continue
		}
		// c := clothes.GetClothesById(clo.ClothesID)
		_, m, o, _ := tools.GetInfoFromCustomIDInGame(cid)
		if o == constants.ORIGIN_DESIGNER {
			if m > 200 {
				dbcnt += 1
			} else if m > 100 {
				dgcnt += 1
			}
		} else {
			if m > 200 {
				obcnt += 1
			} else if m > 100 {
				ogcnt += 1
			}
		}

	}
	return
}

func (self *GFPlayer) GetEventComplexProgress(id int64) string {
	if self.status.DayExtraInfo.EventComplexProgress == nil {
		self.status.DayExtraInfo.EventComplexProgress = make(map[string]string)
	}

	ids := strconv.FormatInt(id, 10)
	return self.status.DayExtraInfo.EventComplexProgress[ids]
}

func (self *GFPlayer) GetEventProgress() (l map[string]int, d map[string][]string, c map[string]string) {
	return self.status.DayExtraInfo.EventLinearProgress, self.status.DayExtraInfo.EventDiscreteProgress, self.status.DayExtraInfo.EventComplexProgress
}

func isSpecialEventKey(key string) bool {
	return strings.HasPrefix(key, constants.SPECIAL_EVENT_KEY_PREFIX)
}

func (self *GFPlayer) checkEventProgress() {
	eps := self.status.DayExtraInfo.EventLinearProgress
	for ids, _ := range eps {
		if ids == "0:COS" {
			delete(self.status.DayExtraInfo.EventLinearProgress, ids)
			continue
		}

		if !isSpecialEventKey(ids) {
			var idi int64
			if strings.Contains(ids, ":") {
				if strings.Contains(ids, ":CMB:") || strings.HasSuffix(ids, ":COS") {
					idps := strings.Split(ids, ":")
					if len(idps) > 0 {
						idi, _ = strconv.ParseInt(idps[0], 10, 64)
					}
				} else if strings.HasSuffix(ids, constants.STAMP_EVENT_TOKEN) {
					idps := strings.Split(ids, "_")
					if len(idps) > 0 {
						idi, _ = strconv.ParseInt(idps[0], 10, 64)
					}
				}
			} else {
				idi, _ = strconv.ParseInt(ids, 10, 64)
			}

			if idi <= 0 {
				Err("event id wrong:", ids)
				continue
			}
			if !event.EventOpen(idi) {
				delete(self.status.DayExtraInfo.EventLinearProgress, ids)
			}
		}
	}

	eps1 := self.status.DayExtraInfo.EventDiscreteProgress
	for ids, _ := range eps1 {
		idi, _ := strconv.ParseInt(ids, 10, 64)
		if idi <= 0 {
			Err("event id wrong:", ids)
			continue
		}
		if !event.EventOpen(idi) {
			delete(self.status.DayExtraInfo.EventDiscreteProgress, ids)
		}
	}

	eps2 := self.status.DayExtraInfo.EventComplexProgress
	for ids, _ := range eps2 {
		idi, _ := strconv.ParseInt(ids, 10, 64)
		if idi <= 0 {
			Err("event id wrong:", ids)
			continue
		}
		if !event.EventOpen(idi) {
			delete(self.status.DayExtraInfo.EventComplexProgress, ids)
		}
	}
}

func (self *GFPlayer) GetPKLevel(weekToken string) int {
	if self.status.ExtraInfo1.PKWeekPoints == nil {
		self.status.ExtraInfo1.PKWeekPoints = make(map[string]int)
	}
	l, _ := common.GenPKLevel(self.status.ExtraInfo1.PKWeekPoints[weekToken])
	return l
}

func (self *GFPlayer) GetSpecialEventLinearProgress(key string) int {
	return self.status.DayExtraInfo.EventLinearProgress[constants.SPECIAL_EVENT_KEY_PREFIX+key]
}

func (self *GFPlayer) GetEventLinearProgress(eid int64) int {
	es := strconv.FormatInt(eid, 10)
	return self.status.DayExtraInfo.EventLinearProgress[es]
}

func (self *GFPlayer) GetEventLinearProgressByKey(key string) int {
	return self.status.DayExtraInfo.EventLinearProgress[key]
}

func (self *GFPlayer) GetEventDiscreteProgress() map[string][]string {
	return self.status.DayExtraInfo.EventDiscreteProgress
}

func (self *GFPlayer) EventDiscreteProgressPassed(eid int64, progress string) bool {
	es := strconv.FormatInt(eid, 10)
	if v, ok := self.status.DayExtraInfo.EventDiscreteProgress[es]; ok {
		for _, ps := range v {
			if ps == progress {
				return true
			}
		}
	}

	return false
}

func (self *GFPlayer) EventDiscreteProgressLen(eid int64) int {
	es := strconv.FormatInt(eid, 10)
	if v, ok := self.status.DayExtraInfo.EventDiscreteProgress[es]; ok {
		return len(v)
	}

	return 0
}

func (self *GFPlayer) GetTaskReward(template *Task, deviceId string) (rewardDiamond int, rewardGold int64, rewardItem1 int, curDiamond int, curGold int64, curItem1 int, err error) {
	have := false
	canReward := true
	for _, v := range self.status.Tasks {
		if v.TaskID == template.TaskID {
			have = true
			if v.Progress >= template.Target && v.RewardGot == false {
				rewardDiamond, rewardGold, rewardItem1 = template.Diamond, template.Gold, template.Fastener
				// curGold, _ = self.AddMoney(rewardGold, constants.SRC_TASK, deviceId, fmt.Sprintf("%d", v.TaskID))
				// curItem1 = self.AddItem1(rewardItem1, constants.SRC_TASK, deviceId)

				v.RewardGot = true
			} else {
				canReward = false
			}

			break
		}
	}

	if !have {
		err = errors.New("task not found")
	} else if !canReward {
		err = errors.New("can not get task reward")
	}

	return
}

func (self *GFPlayer) GetTasks() []*Task {
	return self.status.Tasks
}

func (self *GFPlayer) GetEnding() []int {
	// return self.ending
	return self.status.Ending
}

func (self *GFPlayer) GetEndingSize() int {
	return len(self.status.Ending)
}

func (self *GFPlayer) HaveEnding(endingID int) bool {
	for _, v := range self.status.Ending {
		if v == endingID {
			return true
		}
	}
	return false
}

func (self *GFPlayer) GetScripts() []string {
	// return self.scripts
	return self.status.Scripts
}

func (self *GFPlayer) HaveScript(scriptID string) bool {
	for _, v := range self.status.Scripts {
		if v == scriptID {
			return true
		}
	}
	return false
}

func (self *GFPlayer) getIaped() int {
	return self.status.DayExtraInfo.IAPed
}

func (self *GFPlayer) setIaped(c int) int {
	self.status.DayExtraInfo.IAPed = c
	return self.status.DayExtraInfo.IAPed
}

func (p *GFPlayer) ConfirmOrder1(o Order, deviceId, platform, code, platformOrderID string) (rewardDiamond int, rewardGold int64, diamond int, gold int64, buyDiamond, vip, rewardTili, tili, maxTili, tinyBuyCnt int, isFirst, firstIAPEventReward bool, clot ClothesInfo, status ErrorCode, t17p, t18p, t19p, t20p, rewardDesignCoin, currentDesignCoin int, monthIAP float64) {
	if o.Status != ORDER_SUCCEED || platform == constants.ORDER_PLATFORM_PINGPP {
		success := false
		or := OrderReward{
			DiamondID:       o.DiamondId,
			PlatformOrderID: platformOrderID,
		}
		oe := OrderExtra{
			OS:            p.Platform,
			ThirdUsername: p.ThirdUsername,
			Nickname:      p.Nickname,
			OrderID:       o.OrderId,
			OrderUID:      o.OrderUID,
			DiamondID:     o.DiamondId,
			RMB:           o.RmbCost,
			ServerAddress: appcfg.GetAddress(),
			ThirdChannel:  p.ThirdChannel,
			OrderName:     iap.GetOrderName(o.DiamondId),
			CurrencyType:  constants.CURRENCY_CN,
		}
		if iap.IsTinyPack(o.DiamondId) {
			//tiny pack
			var id int64
			var rmb int
			var clo string
			if o.DiamondId != 11 && o.DiamondId != 1009 {
				id, rmb, rewardGold, rewardDiamond, rewardTili, _, clo = event.GetTinyPackInfo()
			} else {
				id, rmb, rewardGold, rewardDiamond, rewardTili, _, clo = event.GetDefaultTinyPackInfo()
			}

			if iap.CheckTinyPack(rmb, o.DiamondId) || appcfg.GetBool("order_test", false) {
				if appcfg.GetBool("order_test", false) {
					if o.DiamondId == 8 || o.DiamondId == 11 {
						rewardDiamond = 10
						rewardGold = 1000
					} else if o.DiamondId == 9 {
						rewardDiamond = 30
						rewardGold = 3000
					} else if o.DiamondId == 10 {
						rewardDiamond = 60
						rewardGold = 6000
					}
				}

				or.Tili = rewardTili
				or.Gold = rewardGold
				or.Diamond = rewardDiamond
				or.Clothes = clo

				// p.AddMoney(rewardGold, constants.SRC_TINY_PACK, deviceId, fmt.Sprintf("%d", o.OrderId))
				p.AddTili(rewardTili, 0, deviceId, constants.SRC_TINY_PACK)
				if tools.ClothesIDLegal(clo) {
					clot = p.AddPixieClothesWithCID(clo, 1, time.Now().Unix(), constants.SRC_TINY_PACK)
				}

				//tiny pack progress
				if id > 0 {
					p.SetEventLinearProgress(id, p.GetEventLinearProgress(id)+1)
					tinyBuyCnt = p.GetEventLinearProgress(id)
				}

				//first iap progress
				if fid := event.GetFirstIAPInfo(); fid > 0 && !p.EventDiscreteProgressPassed(fid, "1") {
					if (appcfg.GetLanguage() == "" && rmb >= 6) || (appcfg.GetLanguage() == constants.KOR_LANGUAGE && o.DiamondId != 1007) {
						firstIAPEventReward = true
						p.SetEventLinearProgress(fid, 1)
					}
				}
				// //vip progress
				// if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
				// 	if o.DiamondId == 1007 {
				// 		p.AddBuyDiamondOnly(6*20, constants.SRC_TINY_PACK, deviceId)
				// 	} else if o.DiamondId == 1009 {
				// 		p.AddBuyDiamondOnly(60*20, constants.SRC_TINY_PACK, deviceId)
				// 	}
				// } else {
				// 	p.AddBuyDiamondOnly(rmb*20, constants.SRC_TINY_PACK, deviceId)
				// }

				success = true
			} else {
				status = ERR_TINY_PACK_NOT_EXIST
			}
		} else {
			di := iap.GetDiamondItemById(o.DiamondId)
			if di != nil {
				rewardDesignCoin = di.DesignCoin
				if rewardDesignCoin > 0 {
					or.DesignCoin = rewardDesignCoin
					// currentDesignCoin, _ = p.AddDesignCoin(rewardDesignCoin, false, constants.SRC_PAY, deviceId, true)
					// p.AddBuyDiamondOnly(rewardDesignCoin, constants.SRC_PAY, deviceId)
				}
			}

			rewardDiamond = o.Diamond
			if appcfg.GetBool("order_test", false) && o.DiamondId == 7 {
				rewardDiamond = 300
			}

			or.Diamond = rewardDiamond

			if !p.GetFirstIAPReward(o.DiamondId) && ((o.DiamondId >= 1 && o.DiamondId <= 6) || (o.DiamondId >= 1001 && o.DiamondId <= 1006)) {
				reward.SendRewardByTag("first_iap_"+strconv.Itoa(o.DiamondId), p.Username)
				p.SetFirstIAPReward(o.DiamondId)
				isFirst = true

				//buy diamond first double
				or.Diamond += rewardDiamond
			}

			if iap.IsMonthCard(o.DiamondId) {
				p.RefreshMonthCard()

				ot := time.Unix(o.OrderTime, 0).Format("200601")
				if mr := params.GetMonthCardReward(ot); mr.Type > 0 && mr.RewardInfo != "" {
					if !p.MonthCardRewarded(ot) {
						//send reward
						if mr.Type == 1 {
							clos := clothes.GetSuitClothesList(mr.RewardInfo)
							if len(clos) > 0 {
								mails.SendToOneClos(p.Username, L("o13"), L("o14"), clos)
								p.SetMonthCardRewarded(ot)
							}
						} else {
							ErrMail("unknown month card reward:", mr.Type, mr.RewardInfo)
						}
					}
				}

			}

			//first iap progress
			if fid := event.GetFirstIAPInfo(); fid > 0 && !p.EventDiscreteProgressPassed(fid, "1") {
				firstIAPEventReward = true
				p.SetEventLinearProgress(fid, 1)
			}

			//sum diamond iap
			if fid := event.GetSumIAPInfo(); fid > 0 {
				p.SetEventLinearProgress(fid, p.GetEventLinearProgress(fid)+rewardDiamond)
			}
			success = true
		}
		orb, _ := json.Marshal(or)

		if success && o.Status != ORDER_SUCCEED {
			iap.SetOrderStatus(o.OrderUID, ORDER_SUCCEED, p.Username, platform, code, string(orb))
			//kakaogame sendlog 发送请求
			if err := kakao.KorKakaoSendLog(p.ThirdUsername, o.DiamondId, p.Platform); err != nil {
				Err(err)
			}
			if p.getIaped() <= 0 {
				oe.FirstIAP = 1
			}
			p.setIaped(1)

			//订单log插入
			if appcfg.GetLanguage() == constants.KOR_LANGUAGE && appcfg.GetBool("upload_kor_log", false) {
				var okl OrderKorLog
				oe.KORLogType = 90101
				oe.KORThirdCode = tools.GetKorThirdCode(p.ThirdChannel)
				oe.KOROSCode, oe.KORStoreCode, oe.CurrencyType = tools.GetKOROrderInfo(p.Platform)
				oe.KORGameCode = constants.KOR_GAME_CODE
				oe.KORPayCountry = "KR"
				//韩国版本订单日志
				ordertime := time.Now().Format("2006-01-02 15:04:05")
				okl = OrderKorLog{
					Paylogid:       o.OrderId,
					Eventdate:      ordertime,
					Logtype:        90101,
					Guid:           p.UID,
					Platformuserid: p.ThirdUsername,
					Platformcode:   7,
					Oscode:         oe.KOROSCode,
					Storetype:      oe.KORStoreCode,
					Paydate:        ordertime,
					Paytool:        0,
					Gamecode:       constants.KOR_GAME_CODE,
					Gamename:       "yjdyc",
					Itemname:       o.DiamondId,
					Currentcycode:  oe.CurrencyType,
					Paycountry:     "KR",
				}

				if len(platformOrderID) <= 50 {
					okl.Orderid = platformOrderID
				} else {
					Err("platform order id too long:", o.OrderId, platformOrderID, p.Username, p.ThirdUsername)
				}

				kc, uc := iap.GetKorOrderCost(o.DiamondId)
				if oe.CurrencyType == constants.CURRENCY_US {
					okl.Amount = uc
				} else if oe.CurrencyType == constants.CURRENCY_KOR {
					okl.Amount = kc
				}

				okldata, _ := json.Marshal(okl)
				GMLogKor(constants.C1_SYSTEM, constants.C2_ORDER, constants.C3_ORDER_LOG_KOR, p.Username, string(okldata))
			}

			oeb, _ := json.Marshal(oe)
			GMLog(constants.C1_SYSTEM, constants.C2_ORDER, constants.C3_DEFAULT, p.Username, string(oeb))
			global.AddOrderCnt(oe.RMB, oe.CurrencyType, p.ThirdChannel, time.Now())

			p.AddMonthOrder(time.Now().Format("200601"), oe.CurrencyType, oe.RMB)
		} else {
			iap.SetOrderStatus(o.OrderUID, ORDER_FAILED, p.Username, platform, code, string(orb))
		}
	}

	gold = p.GetMoney()
	tili = p.GetTili()
	maxTili = p.GetMaxTili()
	monthIAP = p.GetMonthIAP()

	return
}

func (self *GFPlayer) checkDay(ds string) {
	if self.status.DayExtraInfo.DayVerifyNum == nil {
		self.status.DayExtraInfo.DayVerifyNum = make(map[string]int)
		return
	}

	for k, _ := range self.status.DayExtraInfo.DayVerifyNum {
		if k != ds {
			delete(self.status.DayExtraInfo.DayVerifyNum, k)
		}
	}
}

func (self *GFPlayer) SuitRewarded(sid string) bool {
	if self.status.ExtraInfo1.RewardedSuitMap == nil {
		self.status.ExtraInfo1.RewardedSuitMap = make(map[string]int)
	}

	return self.status.ExtraInfo1.RewardedSuitMap[sid] > 0
}

func (self *GFPlayer) RefreshHeartbeatTime() {
	self.HeartbeatTime = time.Now().Unix()
}

func (self *GFPlayer) GetHeartBeatTime() int64 {
	return self.HeartbeatTime
}

func (self *GFPlayer) GetHeartbeatTime() int64 {
	return self.HeartbeatTime
}

func (self *GFPlayer) GetTili() int {
	return self.status.Tili
}

func (self *GFPlayer) GetMaxTili() int {
	return self.status.MaxTili
}

func (self *GFPlayer) GetTiliRTime() int64 {
	return self.status.TiliRTime
}

func (self *GFPlayer) SaveStatusToDB(name StatusDetailName) error {
	if err := status.UpdateStatusDetail(&self.status, name); err != nil {
		Info("flush error for player:", self.Username, name)
		logOutUser1(self)
		return err
	} else if appcfg.GetBool("persist_test", false) {
		ns, _ := status.Find(self.Username)

		if e, i := Equal(&self.status, ns); !e {
			tools.ReportErr(nil, self.Username+":"+string(name)+":"+i, params.GetErrMailTos())
		}
	}

	return nil
}

func (self *GFPlayer) GetOneClothesCnt(cloID string) int {
	for cid, c := range self.status.Clothes {
		if cid == cloID {
			return c.Count
		}
	}
	return 0
}

func (self *GFPlayer) GetStatus() Status {
	return self.status
}

func (self *GFPlayer) GetFirstIAPReward(diamondID int) bool {
	if self.status.DayExtraInfo.FirstRewardMap == nil {
		self.status.DayExtraInfo.FirstRewardMap = make(map[string]int)
	}
	return self.status.DayExtraInfo.FirstRewardMap[strconv.Itoa(diamondID)] > 0
}

func (self *GFPlayer) CheckMoney(money int64) bool {
	if money >= 0 && self.status.Money >= money {
		return true
	} else if money < 0 && self.status.Money >= -money {
		return true
	} else {
		return false
	}
}

func (self *GFPlayer) GetMoney() int64 {
	return self.status.Money
}

func (self *GFPlayer) GetPxc() float64 {
	return self.status.Pxc
}

func (self *GFPlayer) CheckSkinExist(skinID string) bool {
	if self.status.SkinMap == nil {
		self.status.SkinMap = make(map[string]bool)
	}
	return self.status.SkinMap[skinID]
}

func (self *GFPlayer) WarClothesRewarded(wt string, index int) bool {
	if self.status.ExtraInfo1.WeekWarClothesReward == nil {
		self.status.ExtraInfo1.WeekWarClothesReward = make(map[string][]int)
	}

	// if self.status.ExtraInfo1.WeekWarClothesReward[wt] == nil {
	// 	self.status.ExtraInfo1.WeekWarClothesReward[wt] = make([]int, 0)
	// }

	for _, i := range self.status.ExtraInfo1.WeekWarClothesReward[wt] {
		if i == index {
			return true
		}
	}

	return false
}

func (self *GFPlayer) GetWarClothesRewardInfo(wt string) []int {
	if self.status.ExtraInfo1.WeekWarClothesReward == nil {
		self.status.ExtraInfo1.WeekWarClothesReward = make(map[string][]int)
	}

	// if self.status.ExtraInfo1.WeekWarClothesReward[wt] == nil {
	// 	self.status.ExtraInfo1.WeekWarClothesReward[wt] = make([]int, 0)
	// }

	return self.status.ExtraInfo1.WeekWarClothesReward[wt]
}

func (self *GFPlayer) GetNewPlayerLoginDayCnt() int {
	return self.status.DayExtraInfo.NewPlayerLoginDayCnt
}

func (self *GFPlayer) GetDayBuyTiliCnt() int {
	return self.status.DayExtraInfo.DayBuyTiliCnt
}

// func (self *GFPlayer) GetLeftBuyTiliCnt() int {
// 	if self.status.DayExtraInfo.DayBuyTiliCnt <= constants.DAY_BUY_TILI_LIMIT {
// 		return constants.DAY_BUY_TILI_LIMIT - self.status.DayExtraInfo.DayBuyTiliCnt
// 	} else {
// 		return 0
// 	}
// }

func (self *GFPlayer) GetEventDropDayNum() int {
	return self.status.DayExtraInfo.DayEventDropNum
}

func (self *GFPlayer) GetRegTime() int64 {
	return self.status.DayExtraInfo.RegTime
}

func (self *GFPlayer) GetUploadHelpDayCnt() int {
	return self.status.DayExtraInfo.UploadHelpDayCnt
}

func (self *GFPlayer) AddUploadHelpDayCnt() {
	self.status.DayExtraInfo.UploadHelpDayCnt++

}

func (self *GFPlayer) GetUploadCustomDayCnt(t string) int {
	// return self.status.DayExtraInfo.UploadCustomDayCnt
	if self.status.DayExtraInfo.UploadCustomCntMap == nil {
		self.status.DayExtraInfo.UploadCustomCntMap = make(map[string]int)
	}
	return self.status.DayExtraInfo.UploadCustomCntMap[t]
}

func (self *GFPlayer) AddUploadCustomDayCnt(t string) {
	// self.status.DayExtraInfo.UploadCustomDayCnt++
	if self.status.DayExtraInfo.UploadCustomCntMap == nil {
		self.status.DayExtraInfo.UploadCustomCntMap = make(map[string]int)
	}
	self.status.DayExtraInfo.UploadCustomCntMap[t] += 1
}

func (self *GFPlayer) GetRecordHistoryClothes(lid string) (string, string) {
	if self.status.RecordMap == nil {
		self.status.RecordMap = make(map[string]*PixieRecord)
	}
	if val, exist := self.status.RecordMap[lid]; exist {
		return val.Clothes, val.BestClothes
	}
	return "", ""
}

func (self *GFPlayer) GetFreeMagicSpeedCnt() int {
	return self.status.DayExtraInfo.MagicSpeedLeftCnt
}

// 检查每日登陆任务
func (self *GFPlayer) checkEverydayTask() {
	// taskID:1是每日登陆, taskID:2是7天连续登陆
	if len(self.status.Tasks) > 0 {
		for _, v := range self.status.Tasks {
			// if v.TaskID == 1 {
			// 	v.Progress = 1
			// 	v.RewardGot = false
			// } else if v.TaskID == 2 {
			// 	if v.Target != 7 {
			// 		v.Target = 7
			// 	}
			// 	if v.Progress == 0 {
			// 		v.Progress = 1
			// 		v.RewardGot = false
			// 	} else if v.Progress < 7 {
			// 		today := time.Now().Format("2006-01-02")
			// 		todayInDB := time.Unix(int64(v.Info+24*3600), 0).Format("2006-01-02")

			// 		if today == todayInDB {
			// 			v.Progress = v.Progress + 1
			// 		} else {
			// 			v.Progress = 1
			// 		}
			// 		v.RewardGot = false
			// 	} else if v.Progress >= 7 {
			// 		if v.RewardGot {
			// 			v.Progress = 1
			// 			v.RewardGot = false
			// 		}
			// 	}
			// 	v.Info = int(time.Now().Unix())
			// } else {
			// 每日任务
			if tasks.IsDayTask(v.TaskID) {
				if v.Progress != 0 || v.RewardGot {
					v.Progress = 0
					v.RewardGot = false
				}
			}
			// }
		}
	}
	// else {
	// 	self.status.Tasks = append(self.status.Tasks, &Task{TaskID: 1, Progress: 1, Target: 1, RewardGot: false, Type: 1, RewardType: 1})
	// 	self.status.Tasks = append(self.status.Tasks, &Task{TaskID: 2, Progress: 1, Target: 7, RewardGot: false, Info: int(time.Now().Unix()), Type: 1, RewardType: 1})
	// }
}

func (self *GFPlayer) GetDefaultWardrobe() ([]string, int, int, int, string) {
	return self.status.DayExtraInfo.DefaultClothes, self.status.DayExtraInfo.DefaultLeftModel, self.status.DayExtraInfo.DefaultRightModel, self.status.DayExtraInfo.DefaultShowModel, self.status.DayExtraInfo.DefaultSceneID
}

func (self *GFPlayer) GetClothes() map[string]*ClothesDetail {
	return self.status.Clothes
}

func (self *GFPlayer) GetUnreadMailCount() int {
	cnt, _ := mail.GetUnreadCount(self.Username, self.GuildInfo.IgnoreWarNotice)
	return cnt
}

func (self *GFPlayer) HaveClothes(clothesID string) bool {
	if v, ok := self.status.Clothes[clothesID]; ok {
		if v.Count > 0 {
			return true
		}
	}

	return false
}

func (self *GFPlayer) GetRecords() map[string]*PixieRecord {
	return self.status.RecordMap
}

func (self *GFPlayer) GetRecordBylevelNo(levelNo string) PixieRecord {
	if self.status.RecordMap == nil {
		self.status.RecordMap = make(map[string]*PixieRecord)
	}
	if val, exist := self.status.RecordMap[levelNo]; exist {
		return *val
	}
	return PixieRecord{}
}

func (self *GFPlayer) LevelNoPassed(levelNo string) bool {
	for k, _ := range self.status.RecordMap {
		if fmt.Sprintf("%d", rank.GetLevelNo(k)) == levelNo {
			return true
		}
	}
	return false
}

func (self *GFPlayer) GetClothesCnt() int {
	cnt := 0
	for cid, v := range self.status.Clothes {
		if cid == "" || v.Count == 0 {
			continue
		}
		cnt += 1
	}
	return cnt
}

func (self *GFPlayer) GetMyCosItemCntByCosplayID(cosplayID int64) (cnt int, err error) {
	return servCosItem.GetMyCosItemCntByCosplayID(self.Username, cosplayID)
}

func (self *GFPlayer) Clear() {
	self.status = Status{}
	self.ThirdUsername = ""
	self.GameThirdUsername = ""
	self.Username = ""
	self.Nickname = ""
	self.AccessToken = ""
	self.UID = -1
	self.DeviceToken = ""
	self.HeartbeatTime = 0
	self.GuildOwner = false

	self.ThirdChannel = ""
	self.OldChannel = ""

	self.LastApplyGuildTime = 0
	self.CustomCnt = 0
	self.GuildInfo = GuildMember{}
	self.Platform = ""
	self.NotifiedUser = nil
	self.ShopCart = nil
	self.LastSendETHVerifyTime = 0
	self.LastVistUsername = ""
}

func (self *GFPlayer) IsBanned() (banned bool, st int64, et int64, reason string) {
	now := time.Now().Unix()
	banned = self.status.BanInfo.BanStartTime >= now && self.status.BanInfo.BanEndTime <= now
	st = self.status.BanInfo.BanStartTime
	et = self.status.BanInfo.BanEndTime
	reason = self.status.BanInfo.BanReason
	return
}

func (self *GFPlayer) GetOneMail(mailID int64) (Mail, error) {
	return mail.FindOne(self.Username, mailID)
}

func (self *GFPlayer) DeleteMail(mailID int64) (err error) {
	return mails.DeleteMail(self.Username, mailID)
}

func (self *GFPlayer) SetMailRead(mailID int64) (err error) {
	return mails.SetMailRead(self.Username, mailID)
}

func (self *GFPlayer) GetMails(ignoreWarNotice bool) []*Mail {
	ret, _ := mail.Find(self.Username, ignoreWarNotice)
	return ret
}

func (self *GFPlayer) ModelUnlocked(modelNo int) bool {
	if self.status.DayExtraInfo.UnlockModels == nil {
		self.status.DayExtraInfo.UnlockModels = make([]int, 0)
	}
	for _, v := range self.status.DayExtraInfo.UnlockModels {
		if v == modelNo {
			return true
		}
	}
	return false
}

func (self *GFPlayer) ModelUnlockNum() int {
	if self.status.DayExtraInfo.UnlockModels == nil {
		self.status.DayExtraInfo.UnlockModels = make([]int, 0)
	}
	return len(self.status.DayExtraInfo.UnlockModels)
}

func (self *GFPlayer) GetMonthCardInfo() (int64, int64, int) {
	return self.status.DayExtraInfo.MonthCardStart, self.status.DayExtraInfo.MonthCardEnd, self.status.DayExtraInfo.MonthCardIndex
}

func (self *GFPlayer) GetLevel() int {
	return self.status.Level
}

func (self *GFPlayer) GetLevelAndExp() (int, int) {
	return self.status.Level, self.status.Exp
}

func (self *GFPlayer) GetDayDyeCnt() int {
	return self.status.DayExtraInfo.DayDyeCnt
}

func (self *GFPlayer) GetDayBuQianCnt() int {
	return self.status.DayExtraInfo.DayBuQianCnt
}

func (self *GFPlayer) GetLastSendHopeTime() int64 {
	if appcfg.GetBool("mac", false) {
		return 0
	}
	return self.status.DayExtraInfo.LastSendHopeTime
}

func (self *GFPlayer) GetPKWeekPoints(token string) int {
	if self.status.ExtraInfo1.PKWeekPoints == nil {
		self.status.ExtraInfo1.PKWeekPoints = make(map[string]int)
	}

	return self.status.ExtraInfo1.PKWeekPoints[token]
}

func (self *GFPlayer) GetJuhuaCnt() int {
	return self.status.ExtraInfo1.JuhuaCnt
}

func (self *GFPlayer) GetJuhuaChaCnt() int {
	return self.status.ExtraInfo1.JuhuaChaCnt
}

func (self *GFPlayer) GetIAPReturned() bool {
	return self.status.ExtraInfo1.IAPReturned
}

func (self *GFPlayer) GetYuyueRewarded() bool {
	return self.status.ExtraInfo1.YuyueRewarded
}

func (self *GFPlayer) GetDayQinmiAddBoardCnt(fu string) int {
	if self.status.DayExtraInfo.DayQinmiAddBoardMap == nil {
		self.status.DayExtraInfo.DayQinmiAddBoardMap = make(map[string]int)
	}
	return self.status.DayExtraInfo.DayQinmiAddBoardMap[fu]
}

func (self *GFPlayer) GetDayQinmiReplyBoardCnt(fu string) int {
	if self.status.DayExtraInfo.DayQinmiReplyBoardMap == nil {
		self.status.DayExtraInfo.DayQinmiReplyBoardMap = make(map[string]int)
	}
	return self.status.DayExtraInfo.DayQinmiReplyBoardMap[fu]
}

func (self *GFPlayer) GetDayTheaterCnt(theaterID int) int {
	if self.status.DayExtraInfo.DayTheaterCntMap == nil {
		self.status.DayExtraInfo.DayTheaterCntMap = make(map[string]int)
	}
	return self.status.DayExtraInfo.DayTheaterCntMap[strconv.Itoa(theaterID)]
}

func (self *GFPlayer) CanSaveRecord(levelID string, t time.Time) (bool, bool, int64, bool) {
	// if tid, st, et, _, _, _ := params.GetTheaterInfo(rank.GetTheaterID(levelID)); tid > 0 {
	// 	if st > t.Unix() || et <= t.Unix() {
	// 		return false, false, 0, false
	// 	}
	// 	if self.GetDayTheaterCnt(tid) < constants.DAY_THEATER_PLAY_CNT {
	// 		return true, false, 0, false
	// 	}
	// } else {
	return true, false, 0, false
	// }
	// return false, false, 0, false
}

func (self *GFPlayer) GetDayCopyReportCnt() int {
	return self.status.DayExtraInfo.DayCopyReportCnt
}

func (self *GFPlayer) GetDayShensuCnt(cid string) int {
	if self.status.DayExtraInfo.DayShensuMap == nil {
		self.status.DayExtraInfo.DayShensuMap = make(map[string]int)
	}
	return self.status.DayExtraInfo.DayShensuMap[cid]
}

func (self *GFPlayer) GetDayFreeLottery(poolID, typee int) bool {
	if self.status.DayExtraInfo.DayFreeLotteryMap == nil {
		self.status.DayExtraInfo.DayFreeLotteryMap = make(map[string]bool)
	}

	return self.status.DayExtraInfo.DayFreeLotteryMap[fmt.Sprintf("%d:%d", poolID, typee)]
}

func (self *GFPlayer) GetStageTaskProgress() (t12p, t13p, t14p int) {
	for _, t := range self.status.Tasks {
		if t.TaskID == 12 {
			t12p = t.Progress
		} else if t.TaskID == 13 {
			t13p = t.Progress
		} else if t.TaskID == 14 {
			t14p = t.Progress
		}
	}
	return
}

func (self *GFPlayer) GetGuanzhuDesignerCnt() int {
	return self.status.ExtraInfo1.GuanzhuDesignerCount
}

func (self *GFPlayer) MarkRewardGot(rt int) bool {
	if self.status.DayExtraInfo.MarkGotRewards == nil {
		self.status.DayExtraInfo.MarkGotRewards = make([]int, 0)
		return false
	}

	for _, r := range self.status.DayExtraInfo.MarkGotRewards {
		if r == rt {
			return true
		}
	}

	return false
}

func (self *GFPlayer) GetActivityGot(dt string) int {
	if self.status.DayExtraInfo.DayActivityGot == nil {
		self.status.DayExtraInfo.DayActivityGot = make(map[string]int)
	}
	return self.status.DayExtraInfo.DayActivityGot[dt]
}

func (self *GFPlayer) GetQuitGuildTime() int64 {
	return self.status.ExtraInfo1.LastQuitGuildTime
}

func (self *GFPlayer) WarDateWinOnceRewarded(wd string) bool {
	if self.status.DayExtraInfo.WarDateWinOnceRewardMap == nil {
		self.status.DayExtraInfo.WarDateWinOnceRewardMap = make(map[string]int)
	}
	return self.status.DayExtraInfo.WarDateWinOnceRewardMap[wd] > 0
}

func (self *GFPlayer) CanBuyCustomCnt() bool {
	var cnt int
	for cid, _ := range self.status.Clothes {
		if _, _, ori, _ := tools.GetInfoFromCustomIDInGame(cid); ori == constants.ORIGIN_DESIGNER {
			cnt++
		}
	}

	return true
}

func (self *GFPlayer) GetMonthIAP() float64 {
	return self.status.ExtraInfo1.MonthIap[""]
}

func (self *GFPlayer) HasSuit(sid string) bool {
	sm := clothes.GetSuitClothesMap(sid)
	if len(sm) > 0 {
		var mcnt int

		for cid, _ := range self.status.Clothes {
			if sm[cid] > 0 {
				mcnt++
			}
		}

		if mcnt >= len(sm) {
			return true
		}
	}
	return false
}

func (self *GFPlayer) KorInviteRewarded(typee int, cnt int) bool {
	if self.status.ExtraInfo1.KorInviteRewardMap == nil {
		self.status.ExtraInfo1.KorInviteRewardMap = make(map[string]int)
	}

	return self.status.ExtraInfo1.KorInviteRewardMap[fmt.Sprintf("%d_%d", typee, cnt)] > 0
}

func (self *GFPlayer) GetCloth() int {
	return self.status.ExtraInfo1.Cloth
}

// func (self *GFPlayer) GetItem1() int {
// 	return self.status.ExtraInfo1.Item1
// }

func (self *GFPlayer) GetTheatreStatus() string {
	passedTheaterCnt := make(map[int]int)
	// for _, val := range self.status.RecordMap {
	// 	if !val.NewUnlock {
	// 		if tid := rank.GetTheaterID(r.LevelID); tid > 0 {
	// 			//theater level
	// 			// Info(r.LevelID, r.Rank, tid)
	// 			passedTheaterCnt[tid] += 1
	// 		}
	// 	}
	// }

	var res string
	if passedTheaterCnt[1] > 0 {
		res += "1:xz" + fmt.Sprintf("%d", 1000+passedTheaterCnt[1])
	}

	if passedTheaterCnt[2] > 0 {
		res += "#2:xdrl" + fmt.Sprintf("%d", 1000+passedTheaterCnt[2])
	}

	if passedTheaterCnt[3] > 0 {
		res += "#3:cafe" + fmt.Sprintf("%d", 1000+passedTheaterCnt[3])
	}

	return res
}

func (self *GFPlayer) GetEndingStatus() string {
	var res string
	length := len(self.status.Ending)
	if length > 0 {
		for i := 0; i < length; i++ {
			res += fmt.Sprintf("%d", self.status.Ending[i])
			if i < length-1 {
				res += "#"
			}
		}
	}
	return res
}

func (self *GFPlayer) GetButton() int {
	return self.status.ExtraInfo1.Item1
}

func (self *GFPlayer) GetWire() int {
	return self.status.ExtraInfo1.Wire
}

func (self *GFPlayer) GetNeedle() int {
	return self.status.ExtraInfo1.Needle
}

func (self *GFPlayer) GetGuildBoardReadID() int64 {
	return self.status.ExtraInfo1.LastReadGuildBoardID
}

func (self *GFPlayer) IgnoreNotice() bool {
	return self.status.ExtraInfo1.IgnoreNotice
}

func (self *GFPlayer) IgnorePush() bool {
	return self.status.ExtraInfo1.IgnorePush
}

func (self *GFPlayer) IgnoreNewCustomPush() bool {
	return self.status.ExtraInfo1.IgnoreNewCustomPush
}

func (self *GFPlayer) HasBackground(bid string) bool {
	if self.status.ExtraInfo1.Backgrounds == nil {
		self.status.ExtraInfo1.Backgrounds = make([]string, 0)
		return false
	}

	for _, i := range self.status.ExtraInfo1.Backgrounds {
		if i == bid {
			return true
		}
	}

	return false
}

func (self *GFPlayer) GetDesignCoin() int {
	return self.status.ExtraInfo1.DesignCoin
}

func (self *GFPlayer) KorInitPlayerRewarded() bool {
	return self.status.ExtraInfo1.KorInitPlayerRewarded
}

func (self *GFPlayer) KorRecordCZRewarded() bool {
	return self.status.ExtraInfo1.KorRecordCZRewarded
}

func (self *GFPlayer) MonthCardRewarded(month string) bool {
	if self.status.ExtraInfo1.MonthCardRewardMap == nil {
		self.status.ExtraInfo1.MonthCardRewardMap = make(map[string]int)
		return false
	}

	return self.status.ExtraInfo1.MonthCardRewardMap[month] > 0
}

func (self *GFPlayer) GetDesignCoinPriceToken() int {
	return self.status.ExtraInfo1.DesignCoinPriceToken
}

func (self *GFPlayer) GetDesignCoinPriceTokenBuyCnt(month string) int {
	if self.status.ExtraInfo1.DesignCoinPriceTokenBuyMap == nil {
		self.status.ExtraInfo1.DesignCoinPriceTokenBuyMap = make(map[string]int)
	}
	return self.status.ExtraInfo1.DesignCoinPriceTokenBuyMap[month]
}

func (self *GFPlayer) GetDesignCoinPriceTokenBuyMap() map[string]int {
	if self.status.ExtraInfo1.DesignCoinPriceTokenBuyMap == nil {
		self.status.ExtraInfo1.DesignCoinPriceTokenBuyMap = make(map[string]int)
	}
	return self.status.ExtraInfo1.DesignCoinPriceTokenBuyMap
}

func (self *GFPlayer) GetAttendRTParty(now time.Time) int64 {
	if self.status.ExtraInfo1.AttendRTParty == nil {
		self.status.ExtraInfo1.AttendRTParty = make(map[string]int64)
		return 0
	}

	nt := now.Unix()
	for k, v := range self.status.ExtraInfo1.AttendRTParty {
		if nt < v+constants.RT_PARTY_RECONNECT_TIMEOUT {
			hostID, _ := strconv.ParseInt(k, 10, 64)
			return hostID
		}
	}

	return 0
}

func (self *GFPlayer) GetCurrentDayRTPartyReward(now time.Time) (gold int, diamond int, wire int, gotExp int) {
	if self.status.ExtraInfo1.DayRTPartyRewardMap == nil {
		self.status.ExtraInfo1.DayRTPartyRewardMap = make(map[string]int)
	}

	ds := now.Format("20060102")
	return self.status.ExtraInfo1.DayRTPartyRewardMap[ds+"_gold"], self.status.ExtraInfo1.DayRTPartyRewardMap[ds+"_diamond"], self.status.ExtraInfo1.DayRTPartyRewardMap[ds+"_wire"], self.status.ExtraInfo1.DayRTPartyRewardMap[ds+"_exp"]
}

func (self *GFPlayer) GetETHAccount() string {
	return self.status.ExtraInfo1.ETHAccount
}

func (self *GFPlayer) GetETHPassword() string {
	return tools.GenETHPassword(self.Username, self.status.DayExtraInfo.RegTime)
}

func (self *GFPlayer) GetETHPayPwdEmail() string {
	return self.status.ExtraInfo1.ETHPayPwdEmail
}

func (self *GFPlayer) GetETHPayPasswordEncrypt() string {
	return self.status.ExtraInfo1.ETHPayPwdEncrypt
}

func (self *GFPlayer) GetETHPayPasswordVerifyInfo() (mailSend string, codeSend string, codeExpire int64) {
	mailSend = self.status.ExtraInfo1.ETHPayPwdVerifyCodeEmailSend
	codeSend = self.status.ExtraInfo1.ETHPayPwdVerifyCode
	codeExpire = self.status.ExtraInfo1.ETHPayPwdVerifyCodeExpireTime

	return
}
