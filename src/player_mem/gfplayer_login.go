package player_mem

import (
	"appcfg"
	"common"
	"constants"
	"db/global"
	"db/mail"
	"db_pixie/status"
	"encoding/json"
	"math/rand"
	"service/kor"
	"service/params"
	"strconv"
	"strings"
	"time"
	"tools"

	"db/cart"
	. "logger"
	"rao"
	"service/clothes"

	. "pixie_contract/api_specification"
	"service/eth"
	"service/iap"
	"service/mails"
	"service/tasks"
	"service_pixie/level"
	. "types"
)

func (self *GFPlayer) initialize(username string, uid int64, deviceToken, platform, deviceId, channel, thirdUsername, accessToken, tc, oldChannel string, korRevertUnregister bool) (success bool, banned bool, limited bool, deleting bool, deled bool) {
	//player base info
	self.Username = username
	self.ThirdUsername = thirdUsername
	self.AccessToken = accessToken
	self.DeviceToken = deviceToken
	self.UID = uid
	self.Platform = platform
	self.NotifiedUser = make(map[string]int)
	self.ThirdChannel = tc
	self.OldChannel = oldChannel
	self.ShopCart = make(map[string]*ShopCartDetail)

	if tu, err := status.GetThirdUsername(self.UID); err != nil || tu == "" {
		success = false
		return
	} else {
		self.GameThirdUsername = tu

		tus := strings.Split(tu, constants.THIRD_USERNAME_GAP)
		if len(tus) == 2 {
			self.ThirdUsername = tus[1]
		} else {
			// Info("third username not have -:", tu)
			self.ThirdUsername = tu
		}
	}

	var newPlayer bool
	success, banned, limited, deleting, deled, newPlayer = self.initStatus(deviceId, channel, korRevertUnregister)
	if !success {
		return
	}

	if !newPlayer {
		wt, _ := common.GetGuildYearWeekToken()
		if success = self.InitGuild(true, true, wt); !success {
			return
		}
	}

	self.RefreshHeartbeatTime()

	return
}

func (self *GFPlayer) checkKorInitReward() {
	if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
		sendConfig := appcfg.GetBool("game_gift_is_available", false)

		if sendConfig && !self.KorInitPlayerRewarded() {
			self.SetKorInitPlayerRewarded(true)
			if err := self.SaveStatusToDB(STATUS_DETAIL_E1_MNY_DMD); err == nil {
				//给新玩家发送礼包
				reward1, reward2, reward3 := kor.GetKorInitPlayerReward(self.ThirdUsername)

				if reward3 {
					//发送邀请好友邮件
					if err := mails.SendToOneT("", self.Username, constants.Kor_NORMAL_INVITE_TITLE, constants.Kor_NORMAL_INVITE_CONTENT, 100, 0, "", true, 0); err != nil {
						Err(err)
					}
				}

				if reward2 {
					//发送开服预约邮件
					rmc := MailContent{
						Content: constants.Kor_NORMAL_PLAYER_CONTENT,
						TiliCnt: 30,
					}
					//发送邮件内容
					var mailTo Mail
					mailTo.To = self.Username
					mailTo.Diamond = 300
					mailTo.Gold = 20000
					mailTo.Title = constants.Kor_NORMAL_PLAYER_TITLE
					mailTo.Delete = true

					clos := make([]ClothesInfo, 0)
					closs := clothes.GetSuitClothesList("c900101000000000036")
					for _, c := range closs {
						ci := ClothesInfo{
							ClothesID: c,
							Count:     1,
						}
						clos = append(clos, ci)
					}
					mailTo.Time = time.Now().Unix()
					mailTo.Clothes = clos
					data, _ := json.Marshal(rmc)
					mailTo.Content = string(data)
					if err := mail.SendToOne(mailTo, nil); err != nil {
						Err(err)
					}
				}

				if reward1 {
					//发送cbt邮件
					if err := mails.SendToOneT("", self.Username, constants.Kor_GAMECBT_TITLE, constants.Kor_GAMECBT_CONTENT, 300, 0, "", true, 0); err != nil {
						Err(err)
					}
				}
			}

		}
	}
}

func (self *GFPlayer) checkUnlockModelTask(flush bool) {
	var has1016, has1017, change bool
	var target1016, target1017 int

	if template := tasks.GetTask(1016); template != nil {
		target1016 = template.Target
	}

	if template := tasks.GetTask(1017); template != nil {
		target1017 = template.Target
	}

	um := self.ModelUnlockNum()

	result1016 := um
	result1017 := um
	if result1016 > target1016 {
		result1016 = target1016
	}
	if result1017 > target1017 {
		result1017 = target1017
	}

	for _, v := range self.status.Tasks {
		if v.TaskID == 1016 {
			has1016 = true
			if v.Progress < result1016 {
				change = true
				v.Progress = result1016
			}
		} else if v.TaskID == 1017 {
			has1017 = true
			if v.Progress < result1017 {
				change = true
				v.Progress = result1017
			}
		}

	}

	if !has1016 {
		self.status.Tasks = append(self.status.Tasks, &Task{TaskID: 1016, Progress: result1016})
	}

	if !has1017 {
		self.status.Tasks = append(self.status.Tasks, &Task{TaskID: 1017, Progress: result1017})
	}

	if flush && (change || !has1016 || !has1017) {
		self.SaveStatusToDB(STATUS_DETAIL_TASKS)
	}
}

func (self *GFPlayer) checkE2() {
	if self.status.ExtraInfo2.TotalPxcRewardOnVerifyDueNotifyID == 0 {

	}
}

func (self *GFPlayer) checkE1(n time.Time, flush bool) {
	nt := n.Unix()
	ds := n.Format("20060102")
	var change bool
	if !self.status.ExtraInfo1.DesignCoinPriceTokenInited && self.status.ExtraInfo1.DesignCoinPriceToken == 0 {
		change = true
		self.status.ExtraInfo1.DesignCoinPriceToken = 1
		self.status.ExtraInfo1.DesignCoinPriceTokenInited = true
	}

	if len(self.status.ExtraInfo1.AttendRTParty) > 0 {
		for k, v := range self.status.ExtraInfo1.AttendRTParty {
			if v+constants.RT_PARTY_RECONNECT_TIMEOUT <= nt {
				change = true
				delete(self.status.ExtraInfo1.AttendRTParty, k)
			}
		}
	}

	if len(self.status.ExtraInfo1.DayRTPartyRewardMap) > 0 {
		for k, _ := range self.status.ExtraInfo1.DayRTPartyRewardMap {
			if !strings.HasPrefix(k, ds+"_") {
				change = true
				delete(self.status.ExtraInfo1.DayRTPartyRewardMap, k)
			}
		}
	}

	//check eth account
	if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
		pwd := tools.GenETHPassword(self.Username, self.status.DayExtraInfo.RegTime)

		if self.status.ExtraInfo1.ETHAccount == "" {
			if ethAccount, ethUnlock, err := eth.GenETHAccount(nil, pwd); err == nil {
				change = true
				self.status.ExtraInfo1.ETHAccount = ethAccount
				self.status.ExtraInfo1.ETHUnlock = ethUnlock
			}
		}

		if !self.status.ExtraInfo1.ETHUnlock {
			if ok, _ := eth.UnlockAccount(nil, self.status.ExtraInfo1.ETHAccount, pwd); ok {
				change = true
				self.status.ExtraInfo1.ETHUnlock = true
			}
		}
	}

	if change && flush {
		self.SaveStatusToDB(STATUS_DETAIL_E1_MNY_DMD)
	}
}

func (self *GFPlayer) CheckDayTask(flush bool) (dayProgress int) {
	for _, v := range self.status.Tasks {
		if template := tasks.GetTask(v.TaskID); template != nil {
			if template.Target > 0 && v.TaskID != 36 && v.Progress >= template.Target && template.Type == 1 {
				dayProgress++
			}
		}
	}

	have := false
	change := false
	for _, v := range self.status.Tasks {
		if v.TaskID == 36 {
			have = true

			if v.Progress != dayProgress {
				change = true
				v.Progress = dayProgress
			}

			break
		}
	}

	if !have {
		task := Task{
			TaskID:    36,
			Progress:  dayProgress,
			RewardGot: false,
		}

		self.status.Tasks = append(self.status.Tasks, &task)
	}

	if flush && (!have || change) {
		self.SaveStatusToDB(STATUS_DETAIL_TASKS)
	}

	return
}

func (self *GFPlayer) CheckOrder(deviceId string, flush bool) {
	//STATUS_DETAIL_VIP
	change := false
	if os, err := iap.GetPaidOrder(self.Username); err == nil {
		for _, o := range os {
			if success, _ := iap.SSetOrderStatus(o.OrderUID, ORDER_SUCCEED, ORDER_PAID); success {
				change = true
				var poid string
				if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
					poid = o.Info
				}

				self.ConfirmOrder1(o, deviceId, "", "", poid)
				rew := iap.GetDiamondItemById(o.DiamondId)
				mails.SendToOne("", self.Username, rew.Title, rew.Content, 0, 0, "", true, 0)
			}
		}
	}

	if flush && change {
		self.SaveStatusToDB(STATUS_DETAIL_VIP)
	}
}

func (self *GFPlayer) CheckCronMail() bool {
	if changed, nids := mails.CheckCronMail(self.Username, self.status.DayExtraInfo.RegTime, self.status.DayExtraInfo.ReceivedCronMails, self.ThirdChannel); changed {
		self.status.DayExtraInfo.ReceivedCronMails = nids
		return true
	}
	return false
}

func (self *GFPlayer) initStatus(deviceId string, channel string, korRevertUnregister bool) (success bool, banned bool, limited bool, deleting bool, deled bool, newPlayer bool) {
	now := time.Now()
	savedStatus, err := status.Find(self.Username)
	if err == nil {
		if savedStatus.BanInfo.BanStartTime <= now.Unix() && savedStatus.BanInfo.BanEndTime >= now.Unix() {
			Info("player banned:", self.UID, savedStatus.Username)
			banned = true
			return
		}

		revertDel := false
		if savedStatus.ExtraInfo1.DelTime > now.Unix()+120 {
			if korRevertUnregister && appcfg.GetLanguage() == constants.KOR_LANGUAGE {
				savedStatus.ExtraInfo1.DelTime = 0

				//delete from unregister job
				if err := kor.DelUserUnregister(self.Username); err != nil {
					success = false
					return
				}

				revertDel = true
			} else {
				deleting = true
				return
			}
		} else if savedStatus.ExtraInfo1.DelTime > 0 {
			deled = true
			// status.DelUser(self.Username, self.dbAddr)
			return
		}

		//player exist
		if self.ShopCart, err = cart.ListCart(self.Username); err != nil {
			return
		}

		self.status = *savedStatus
		if revertDel {
			self.SaveStatusToDB(STATUS_DETAIL_E1_MNY_DMD)
		}

		self.Nickname = savedStatus.Nickname
		if appcfg.GetBool("random_sex", false) {
			self.status.Sex = rand.Intn(2)
		}

		self.checkKorInitReward()
		self.CheckCronMail()

		self.checkUnlockModelTask(false)
		self.checkE1(now, false)
		self.CheckDayTask(false)
		self.CheckOrder(deviceId, false)
		self.RecoverTili(deviceId)
		self.RefreshStreet(false)

		self.SaveStatusToDB("")

		success = true
		return
	} else if err == constants.PlayerNotFoundErr {
		var pwd, ethAccount string
		ethUnlock := true

		if appcfg.SupportEthereum() {
			pwd = tools.GenETHPassword(self.Username, now.Unix())
			ethAccount, ethUnlock, err = eth.GenETHAccount(nil, pwd)
			if err != nil {
				return
			}
		}

		newPlayer = true
		rp := constants.CHANNEL_DEFAULT
		if channel != "" {
			rp = channel
		}
		if params.ChannelLimited() {
			limit := params.ChannelLimitNum(rp)
			if limit > 0 {

				//get now num
				nn := rao.GetChannelPlayerNum(rp)
				if nn >= limit {
					Info("player limited:", rp, self.Username, nn, limit)
					limited = true
					return
				}
			}
		}

		// 玩家不存在，新玩家
		self.status = Status{
			Username: self.Username,
			Nickname: constants.DEFAULT_NICKNAME_PREFIX + strconv.FormatInt(self.UID, 10),
			Uid:      self.UID,
			Head:     startHead,
			Updated:  time.Now(),
			MapID:    "0",

			Level: 1,
			Exp:   0,

			Tili:      constants.START_TILI,
			MaxTili:   level.GetLevelMaxTili(1),
			TiliRTime: now.Unix(),
		}
		self.Nickname = self.status.Nickname

		self.status.RecordMap = make(map[string]*PixieRecord)
		self.status.NpcFriendship = make(map[string]*FriendshipDetail)

		self.status.Scripts = make([]string, 0)
		self.status.Ending = make([]int, 0)
		self.status.Tasks = make([]*Task, 0)
		self.status.Clothes = make(map[string]*ClothesDetail)
		self.status.HomeShow = HomeShowDetail{}
		self.status.SuitMap = make(map[string][]Suit)
		self.status.SkinMap = map[string]bool{"0": true}
		self.refreshNpcBuffResult()
		if _, err = self.initStreet(1, now.Format(STREET_INCOME_DATE_FORMAT)); err != nil {
			return
		}

		self.status.DayExtraInfo = Extra{
			RegTime:                 now.Unix(),
			NewPlayerLoginDayCnt:    1,
			ReceivedCronMails:       make([]int64, 0),
			EventLinearProgress:     make(map[string]int),
			EventDiscreteProgress:   make(map[string][]string),
			DayQinmiAddBoardMap:     make(map[string]int),
			DayQinmiReplyBoardMap:   make(map[string]int),
			DayTheaterCntMap:        make(map[string]int),
			DayBuyTheaterCntMap:     make(map[string]int),
			DayShensuMap:            make(map[string]int),
			DayFreeLotteryMap:       make(map[string]bool),
			DayActivityGot:          make(map[string]int),
			WarDateWinOnceRewardMap: make(map[string]int),
			UploadCustomCntMap:      make(map[string]int),
			DaySpecialMissionCntMap: make(map[string]int),
			UnlockModels:            make([]int, 0),
			DayVerifyNum:            make(map[string]int),
		}

		self.status.ExtraInfo1 = Extra1{
			PKWeekPoints:               make(map[string]int),
			PKGWeekPoints:              make(map[string]int),
			WeekWarClothesReward:       make(map[string][]int),
			PartnerInviteMap:           make(map[string]int64),
			LoginDayCnt:                1,
			DayStampChangeCntMap:       make(map[string]int),
			RewardedSuitMap:            make(map[string]int),
			KorInviteRewardMap:         make(map[string]int),
			Backgrounds:                make([]string, 0),
			MonthCardRewardMap:         make(map[string]int),
			DesignCoinPriceToken:       1,
			DesignCoinPriceTokenInited: true,
			DesignCoinPriceTokenBuyMap: make(map[string]int),
			AttendRTParty:              make(map[string]int64),
			DayRTPartyRewardMap:        make(map[string]int),
			ETHAccount:                 ethAccount,
			ETHUnlock:                  ethUnlock,
		}

		// 每日登陆
		self.checkEverydayTask()
		// self.AddMoney(startGold, constants.SRC_FIRST_PLAY, deviceId, "")
		self.AddGold(startGold)
		self.AddPxc(startPxc)
		insertErr := status.AddNewPlayer(&self.status)
		if insertErr != nil {
			Err("player insert error:", self.Username, self.AccessToken, self.UID)
			return
		} else {
			global.AddNewPlayerCnt(self.ThirdChannel, now)

			e := NewPlayerExtra{
				ThirdChannel:  self.ThirdChannel,
				ServerAddress: appcfg.GetAddress(),
			}
			eb, _ := json.Marshal(e)
			GMLog(constants.C1_PLAYER, constants.C2_NEW, constants.C3_DEFAULT, self.Username, string(eb))
			if params.ChannelLimited() {
				rao.AddChannelPlayerNum(rp, 1)
			}

			self.checkKorInitReward()

			success = true
			return
		}
	} else {
		return
	}
}
