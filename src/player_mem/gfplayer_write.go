package player_mem

import (
	"appcfg"
	"common"
	"constants"
	"dao_player"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

import (
	. "language"
	. "logger"
	. "pixie_contract/api_specification"
	"service/clothes"
	"service/guild"
	"service/mails"
	"service_pixie/design_config"
	"service_pixie/level"
	"service_pixie/npc"
	. "types"
)

func (self *GFPlayer) AddFriendshipExp(npcID string, exp int) {
	if self.status.NpcFriendship == nil {
		self.status.NpcFriendship = make(map[string]*FriendshipDetail)
	}

	levelChange := false

	if self.status.NpcFriendship[npcID] == nil {
		levelChange = true
		self.status.NpcFriendship[npcID] = &FriendshipDetail{
			Level: 1,
			Exp:   exp,
		}
	} else {
		fd := self.status.NpcFriendship[npcID]
		fd.Exp += exp
		levelupExp := npc.GetNpcFriendshipLevelupExp(npcID, fd.Level)

		if levelupExp > 0 && fd.Exp >= levelupExp {
			levelChange = true
			fd.Level++
			fd.Exp -= levelupExp
		}
	}

	if levelChange {
		//refresh npc buff list
		if self.status.NpcBuff.NpcBuffList == nil {
			self.status.NpcBuff.NpcBuffList = make([]NpcBuffListDetail, 0)
		}

		self.status.NpcBuff.NpcBuffList = append(self.status.NpcBuff.NpcBuffList, NpcBuffListDetail{npcID, self.status.NpcFriendship[npcID].Level})

		self.refreshNpcBuffResult()
	}
}

func (self *GFPlayer) SetStreetName(name string) {
	self.status.StreetName = name
}

func (self *GFPlayer) SetMapID(skinID string) {
	self.status.MapID = skinID
}

func (self *GFPlayer) SetNickname(name string) (err error) {
	if err = dao_player.UpdatePlayerNickname(self.Username, name); err == nil {
		self.status.Nickname = name
		self.Nickname = name

		if !self.status.ExtraInfo1.NicknameSet {
			self.status.ExtraInfo1.NicknameSet = true
		}
	}

	return
}

func (self *GFPlayer) UseGold(moneyCount int64) {
	before := self.status.Money
	self.status.Money -= moneyCount
	after := self.status.Money

	c := CurrencyChangeExtra{
		Nickname:      self.Nickname,
		Level:         self.status.Level,
		Before:        float64(before),
		Amount:        float64(moneyCount),
		Left:          float64(after),
		ServerAddress: appcfg.GetAddress(),
	}
	data, _ := json.Marshal(c)
	GMLog(constants.C1_PLAYER, constants.C2_MONEY, "USE", self.Username, string(data))

	return
}

func (self *GFPlayer) UsePxc(moneyCount float64) {
	before := self.status.Pxc
	self.status.Pxc -= moneyCount
	after := self.status.Pxc

	c := CurrencyChangeExtra{
		Nickname:      self.Nickname,
		Level:         self.status.Level,
		Before:        before,
		Amount:        moneyCount,
		Left:          after,
		ServerAddress: appcfg.GetAddress(),
	}
	data, _ := json.Marshal(c)
	GMLog(constants.C1_PLAYER, constants.C2_PXC, "USE", self.Username, string(data))

	return
}

func (self *GFPlayer) AddGold(moneyCount int64) {
	self.status.Money += moneyCount
	return
}

func (self *GFPlayer) AddPxc(moneyCount float64) {
	self.status.Pxc += moneyCount
	return
}

func (self *GFPlayer) SetHomeShow(hs HomeShowDetail) {
	self.status.HomeShow = hs
}

func (self *GFPlayer) AddSkinMap(name string) {
	if self.status.SkinMap == nil {
		self.status.SkinMap = make(map[string]bool)
	}
	self.status.SkinMap[name] = true
}

func (self *GFPlayer) AddSuit(modelType string, sd Suit, md5Str string) {
	if self.status.SuitMap == nil {
		self.status.SuitMap = make(map[string][]Suit)
	}
	sd.SuitMD5 = md5Str
	self.status.SuitMap[modelType] = append(self.status.SuitMap[modelType], sd)
}

func (self *GFPlayer) UpdateSuit(modelType string, oldMd5 string, sd Suit, newMd5 string) {
	if self.status.SuitMap == nil {
		self.status.SuitMap = make(map[string][]Suit)
	}
	for i, v := range self.status.SuitMap[modelType] {
		if v.SuitMD5 == oldMd5 {
			sd.SuitMD5 = newMd5
			self.status.SuitMap[modelType][i] = sd
			return
		}
	}
}

func (self *GFPlayer) UpdateSuitName(modelType, oldMd5, suitName string) {
	for i, v := range self.status.SuitMap[modelType] {
		if v.SuitMD5 == oldMd5 {
			self.status.SuitMap[modelType][i].SuitName = suitName
			return
		}
	}
}

func (self *GFPlayer) DeleteSuit(modelType string, md5 string) {
	if self.status.SuitMap == nil {
		self.status.SuitMap = make(map[string][]Suit)
	}
	for i, v := range self.status.SuitMap[modelType] {
		if v.SuitMD5 == md5 {
			self.status.SuitMap[modelType] = append(self.status.SuitMap[modelType][:i], self.status.SuitMap[modelType][i+1:]...)
			return
		}
	}
}

func (self *GFPlayer) SetKorInviteRewarded(typee int, cnt int) {
	if self.status.ExtraInfo1.KorInviteRewardMap == nil {
		self.status.ExtraInfo1.KorInviteRewardMap = make(map[string]int)
	}

	self.status.ExtraInfo1.KorInviteRewardMap[fmt.Sprintf("%d_%d", typee, cnt)] = 1
}

func (self *GFPlayer) SetSuitRewarded(sid string) {
	if self.status.ExtraInfo1.RewardedSuitMap == nil {
		self.status.ExtraInfo1.RewardedSuitMap = make(map[string]int)
	}

	self.status.ExtraInfo1.RewardedSuitMap[sid] = 1
}

func (self *GFPlayer) UseSpecialMissionCnt(lid string, t time.Time) int {
	if self.status.DayExtraInfo.DaySpecialMissionCntMap == nil {
		self.status.DayExtraInfo.DaySpecialMissionCntMap = make(map[string]int)
	}
	token := fmt.Sprintf("%s_%s", t.Format("20060102"), lid)
	self.status.DayExtraInfo.DaySpecialMissionCntMap[token]++
	return self.status.DayExtraInfo.DaySpecialMissionCntMap[token]
}

func (self *GFPlayer) AddMonthOrder(mt string, currency string, cost float64) float64 {
	if self.status.ExtraInfo1.MonthIap == nil {
		self.status.ExtraInfo1.MonthIap = make(map[string]float64)
	}

	if self.status.ExtraInfo1.MonthToken == "" || self.status.ExtraInfo1.MonthToken != mt {
		self.status.ExtraInfo1.MonthToken = mt
		self.status.ExtraInfo1.MonthIap = make(map[string]float64)
	}
	self.status.ExtraInfo1.MonthIap[currency] += cost

	return self.status.ExtraInfo1.MonthIap[currency]
}

func (self *GFPlayer) SetWarDateWinOnceRewarded(wd string) {
	if self.status.DayExtraInfo.WarDateWinOnceRewardMap == nil {
		self.status.DayExtraInfo.WarDateWinOnceRewardMap = make(map[string]int)
	}
	self.status.DayExtraInfo.WarDateWinOnceRewardMap[wd] = 1
}

func (self *GFPlayer) AddGuildActivity(a int, t time.Time) error {
	if self.GuildInfo.GID > 0 {
		dt := t.Format("20060102")
		now := self.GetActivityGot(dt)
		if now < constants.DAY_ACTIVITY_GOT_LIMIT {
			delta := constants.DAY_ACTIVITY_GOT_LIMIT - now
			if a > delta {
				a = delta
			}

			wt, _ := common.GetGuildYWTokenByTime(t, 0, 0)
			if ok, err := guild.AddMemberActivity(self.GuildInfo.GID, self.Username, a, wt); err != nil {
				return err
			} else if ok {
				self.AddActivityGot(a, dt)
				self.GuildInfo.Activity += a
			}
		}
	}
	return nil
}

func (self *GFPlayer) CheckGuildActivityReward(g *Guild, week string, t time.Time) {
	if self.GuildInfo.GID > 0 {
		var cid string

		if !self.GuildInfo.L1Got && self.GuildInfo.Activity >= constants.GUILD_MEMBER_LOTTERY1_ACTIVITY && g.Lottery1.LotteryUsername != "" {
			//send reward
			cid = g.Lottery1.LotteryClothes
			if clo := clothes.GetClothesById(cid); clo != nil {
				if ok, _ := guild.TryLotteryUser(self.GuildInfo.GID, self.Username, 1, week); ok {
					//send mail
					mails.SendToOne("", self.Username, L("guild9"), fmt.Sprintf(L("guild10"), g.Lottery1.LotteryNickname, clo.Name), 0, 0, cid, true, t.Unix()+constants.DEFAULT_MAIL_EXPIRE_TIME)
				}
			}
		}

		if !self.GuildInfo.L2Got && self.GuildInfo.Activity >= constants.GUILD_MEMBER_LOTTERY2_ACTIVITY && g.Lottery2.LotteryUsername != "" {
			//send reward
			cid = g.Lottery2.LotteryClothes
			if clo := clothes.GetClothesById(cid); clo != nil {
				if ok, _ := guild.TryLotteryUser(self.GuildInfo.GID, self.Username, 2, week); ok {
					//send mail
					mails.SendToOne("", self.Username, L("guild9"), fmt.Sprintf(L("guild10"), g.Lottery2.LotteryNickname, clo.Name), 0, 0, cid, true, t.Unix()+constants.DEFAULT_MAIL_EXPIRE_TIME)
				}
			}
		}

		if !self.GuildInfo.L3Got && self.GuildInfo.Activity >= constants.GUILD_MEMBER_LOTTERY3_ACTIVITY && g.Lottery3.LotteryUsername != "" {
			//send reward
			cid = g.Lottery3.LotteryClothes
			if clo := clothes.GetClothesById(cid); clo != nil {
				if ok, _ := guild.TryLotteryUser(self.GuildInfo.GID, self.Username, 3, week); ok {
					//send mail
					mails.SendToOne("", self.Username, L("guild9"), fmt.Sprintf(L("guild10"), g.Lottery3.LotteryNickname, clo.Name), 0, 0, cid, true, t.Unix()+constants.DEFAULT_MAIL_EXPIRE_TIME)
				}
			}
		}

		if !self.GuildInfo.L4Got && self.GuildInfo.Activity >= constants.GUILD_MEMBER_LOTTERY4_ACTIVITY && g.Lottery4.LotteryUsername != "" {
			//send reward
			cid = g.Lottery4.LotteryClothes
			if clo := clothes.GetClothesById(cid); clo != nil {
				if ok, _ := guild.TryLotteryUser(self.GuildInfo.GID, self.Username, 4, week); ok {
					//send mail
					mails.SendToOne("", self.Username, L("guild9"), fmt.Sprintf(L("guild10"), g.Lottery4.LotteryNickname, clo.Name), 0, 0, cid, true, t.Unix()+constants.DEFAULT_MAIL_EXPIRE_TIME)
				}
			}
		}
	}
}

func (self *GFPlayer) RefreshQuitGuildTime() {
	self.status.ExtraInfo1.LastQuitGuildTime = time.Now().Unix()
}

func (self *GFPlayer) AddActivityGot(a int, dt string) {
	if self.status.DayExtraInfo.DayActivityGot == nil {
		self.status.DayExtraInfo.DayActivityGot = make(map[string]int)
	}
	self.status.DayExtraInfo.DayActivityGot[dt] += a
}

func (self *GFPlayer) AddWarClothesRewarded(wt string, index int) {
	if self.status.ExtraInfo1.WeekWarClothesReward == nil {
		self.status.ExtraInfo1.WeekWarClothesReward = make(map[string][]int)
	}

	if self.status.ExtraInfo1.WeekWarClothesReward[wt] == nil {
		self.status.ExtraInfo1.WeekWarClothesReward[wt] = make([]int, 0)
	}

	self.status.ExtraInfo1.WeekWarClothesReward[wt] = append(self.status.ExtraInfo1.WeekWarClothesReward[wt], index)
}

func (self *GFPlayer) SetGuildInfo(gid int64, owner bool) {
	self.GuildInfo.GID = gid
	self.GuildOwner = owner

	if gid == 0 {
		self.GuildInfo.IgnoreWar = false
		self.GuildInfo.IgnoreWarNotice = false
		self.GuildInfo.WarDate = ""
		self.GuildInfo.WarClothes = ""
		self.GuildInfo.LeftAttackCnt = 0
	}
}

func (self *GFPlayer) SetDelTime(t int64) {
	self.status.ExtraInfo1.DelTime = t
}

func (self *GFPlayer) RefreshLastSendHopeTime() {
	self.status.DayExtraInfo.LastSendHopeTime = time.Now().Unix()
}

func (self *GFPlayer) ResetLastSendHopeTime() {
	self.status.DayExtraInfo.LastSendHopeTime = 0
}

func (self *GFPlayer) AddDayBuQianCnt() int {
	self.status.DayExtraInfo.DayBuQianCnt++
	return self.status.DayExtraInfo.DayBuQianCnt
}

func (self *GFPlayer) UseDyeCnt(cnt int) int {
	self.status.DayExtraInfo.DayDyeCnt -= cnt
	return self.status.DayExtraInfo.DayDyeCnt
}

func (self *GFPlayer) AddExp(exp int, src string) (nlevel int, nexp int, maxTili int, gotExp int) {
	if self.status.Level == 0 {
		self.status.Level = 1
		self.status.Exp = 0
	} else if self.status.Level >= 50 {
		self.status.Level = 50
		self.status.Exp = 0
		return self.status.Level, self.status.Exp, self.status.MaxTili, 0
	}

	if exp == 0 {
		return self.status.Level, self.status.Exp, self.status.MaxTili, 0
	}

	bl := self.status.Level
	be := self.status.Exp

	self.status.Exp += exp

	needExp := level.GetLevelNeedExp(self.status.Level)

	for self.status.Exp >= needExp && needExp > 0 {
		self.status.Level++
		self.status.Exp -= needExp
		needExp = level.GetLevelNeedExp(self.status.Level)

		if needExp == 0 {
			Info("add exp err:", self.status.Level, self.status.Exp, needExp)
			break
		}
	}

	self.status.MaxTili = level.GetLevelMaxTili(self.status.Level)

	e := ExpChangeExtra{
		GotExp:        exp,
		BeforeLevel:   bl,
		BeforeExp:     be,
		AfterLevel:    self.status.Level,
		AfterExp:      self.status.Exp,
		Reason:        src,
		ServerAddress: appcfg.GetAddress(),
	}

	eb, _ := json.Marshal(e)
	GMLog(constants.C1_PLAYER, constants.C2_EXP, src, self.Username, string(eb))

	return self.status.Level, self.status.Exp, self.status.MaxTili, exp
}

func (self *GFPlayer) AddInvitePlayerCnt() int {
	self.status.ExtraInfo1.DayInviteFriendCnt++
	return self.status.ExtraInfo1.DayInviteFriendCnt
}

func (self *GFPlayer) GetMonthCardReward(deviceID string) (got bool, index int, rd int, cd int) {
	nowDI := time.Now()

	startDS, _ := time.ParseInLocation("2006-01-02", time.Unix(self.status.DayExtraInfo.MonthCardStart, 0).Format("2006-01-02"), time.Local)
	nowDS, _ := time.ParseInLocation("2006-01-02", nowDI.Format("2006-01-02"), time.Local)

	nowIndex := int(nowDS.Sub(startDS).Hours()/24 + 1)

	if self.status.DayExtraInfo.MonthCardStart <= nowDI.Unix() && nowDI.Unix() <= self.status.DayExtraInfo.MonthCardEnd && self.status.DayExtraInfo.MonthCardIndex < nowIndex {
		index = nowIndex
		got = true
		rd = 70

		self.status.DayExtraInfo.MonthCardIndex = nowIndex
	}

	return
}

func (self *GFPlayer) RefreshMonthCard() (int64, int64, int) {
	now := time.Now().Unix()

	if self.status.DayExtraInfo.MonthCardEnd <= now {
		//new month card
		self.status.DayExtraInfo.MonthCardStart = now
		self.status.DayExtraInfo.MonthCardEnd = now + 30*24*3600
		self.status.DayExtraInfo.MonthCardIndex = 0
	} else {
		//extend month card
		self.status.DayExtraInfo.MonthCardEnd += 30 * 24 * 3600
	}

	return self.status.DayExtraInfo.MonthCardStart, self.status.DayExtraInfo.MonthCardEnd, self.status.DayExtraInfo.MonthCardIndex
}

func (self *GFPlayer) UnlockModel(modelNo int) {
	if self.status.DayExtraInfo.UnlockModels == nil {
		self.status.DayExtraInfo.UnlockModels = make([]int, 0)
	}
	self.status.DayExtraInfo.UnlockModels = append(self.status.DayExtraInfo.UnlockModels, modelNo)
}

func (self *GFPlayer) SaveDefaultClothes(clos []string, left, right, show int, sceneID string) {
	self.status.DayExtraInfo.DefaultClothes = clos
	self.status.DayExtraInfo.DefaultLeftModel = left
	self.status.DayExtraInfo.DefaultRightModel = right
	self.status.DayExtraInfo.DefaultShowModel = show
	self.status.DayExtraInfo.DefaultSceneID = sceneID
}

func (self *GFPlayer) AddDayBuyTiliCnt() int {
	self.status.DayExtraInfo.DayBuyTiliCnt++
	return self.status.DayExtraInfo.DayBuyTiliCnt
}

func (self *GFPlayer) AddEventDropDayNum(cnt int) {
	self.status.DayExtraInfo.DayEventDropNum += cnt
}

func (self *GFPlayer) SetEventComplexProgress(id int64, progress string) {
	if self.status.DayExtraInfo.EventComplexProgress == nil {
		self.status.DayExtraInfo.EventComplexProgress = make(map[string]string)
	}

	ids := strconv.FormatInt(id, 10)
	self.status.DayExtraInfo.EventComplexProgress[ids] = progress
}

func (self *GFPlayer) SetEventLinearProgressWithKey(key string, progress int) {
	if self.status.DayExtraInfo.EventLinearProgress == nil {
		self.status.DayExtraInfo.EventLinearProgress = make(map[string]int)
	}

	self.status.DayExtraInfo.EventLinearProgress[key] = progress
}

func (self *GFPlayer) SetSpecialEventLinearProgress(key string, value int) {
	if self.status.DayExtraInfo.EventLinearProgress == nil {
		self.status.DayExtraInfo.EventLinearProgress = make(map[string]int)
	}
	self.status.DayExtraInfo.EventLinearProgress[constants.SPECIAL_EVENT_KEY_PREFIX+key] = value
}

func (self *GFPlayer) SetEventLinearProgress(eid int64, progress int) {
	if self.status.DayExtraInfo.EventLinearProgress == nil {
		self.status.DayExtraInfo.EventLinearProgress = make(map[string]int)
	}
	es := strconv.FormatInt(eid, 10)
	self.status.DayExtraInfo.EventLinearProgress[es] = progress
}

func (self *GFPlayer) AddEventDiscreteProgress(eid int64, progress string) {
	if self.status.DayExtraInfo.EventDiscreteProgress == nil {
		self.status.DayExtraInfo.EventDiscreteProgress = make(map[string][]string)
	}
	es := strconv.FormatInt(eid, 10)
	if _, ok := self.status.DayExtraInfo.EventDiscreteProgress[es]; ok {
		self.status.DayExtraInfo.EventDiscreteProgress[es] = append(self.status.DayExtraInfo.EventDiscreteProgress[es], progress)
	} else {
		s := make([]string, 0)
		s = append(s, progress)
		self.status.DayExtraInfo.EventDiscreteProgress[es] = s
	}
}

func (self *GFPlayer) IncreTask(taskID int, incre int, template *Task) (isOver bool, changed bool, currentProgress int, dayProgress int, err error) {
	if template == nil {
		Info("no such task:", taskID)
		err = errors.New("no such task")
		return
	}
	have := false
	for _, v := range self.status.Tasks {
		if v.TaskID == taskID {
			have = true
			if v.Progress < template.Target {
				v.Progress += incre
				if v.Progress > template.Target {
					v.Progress = template.Target
				}

				changed = true
				if v.Progress == template.Target {
					isOver = true
				} else {
					isOver = false
				}
			} else {
				isOver = true
				changed = false
			}
			currentProgress = v.Progress
			break
		}
	}

	if !have {
		changed = true
		currentProgress = incre
		if currentProgress > template.Target {
			currentProgress = template.Target
		}

		task := Task{
			TaskID:    taskID,
			Progress:  currentProgress,
			RewardGot: false,
		}

		self.status.Tasks = append(self.status.Tasks, &task)
	}

	dayProgress = self.CheckDayTask(false)

	return
}

func (self *GFPlayer) UpdateTask(taskID int, progress int, template *Task) (isOver bool, dayProgress int, err error) {
	if template == nil {
		err = errors.New("no such task")
		return
	}

	have := false
	for _, v := range self.status.Tasks {
		if v.TaskID == taskID {
			have = true
			v.Progress = progress

			if v.Progress == template.Target {
				isOver = true
			} else {
				isOver = false
			}

			break
		}
	}

	if !have {
		task := Task{
			TaskID:    taskID,
			Progress:  progress,
			RewardGot: false,
		}

		self.status.Tasks = append(self.status.Tasks, &task)
	}

	dayProgress = self.CheckDayTask(false)

	return
}

func (self *GFPlayer) RefreshTaskIfLessThan(taskID int, progress int, template *Task) (changed bool, nowProgress int, dayProgress int, err error) {
	if template == nil {
		err = errors.New("no such task")
		return
	}
	rp := progress
	if progress > template.Target {
		rp = template.Target
	}

	nowProgress = rp
	have := false
	for _, v := range self.status.Tasks {
		if v.TaskID == taskID {
			have = true
			if v.Progress < rp {
				changed = true
				v.Progress = rp
			}

			nowProgress = v.Progress

			break
		}
	}

	if !have {
		changed = true
		task := Task{
			TaskID:    taskID,
			Progress:  rp,
			RewardGot: false,
		}

		self.status.Tasks = append(self.status.Tasks, &task)
	}

	dayProgress = self.CheckDayTask(false)

	return
}

func (self *GFPlayer) AddEnding(endingID int) bool {
	self.status.Ending = append(self.status.Ending, endingID)

	return true
}

func (self *GFPlayer) AddScript(scriptID string) bool {
	self.status.Scripts = append(self.status.Scripts, scriptID)

	return true
}

func (self *GFPlayer) DeleteScript(scriptID string) bool {
	for i, v := range self.status.Scripts {
		if v == scriptID {
			self.status.Scripts = append(self.status.Scripts[:i], self.status.Scripts[i+1:]...)
			break
		}
	}

	return true
}

// 回复体力
func (self *GFPlayer) RecoverTili(deviceId string) (change bool) {
	now := time.Now().Unix()
	if self.status.Tili >= self.status.MaxTili {
		change = true
		self.status.Tili = self.status.MaxTili
		self.status.TiliRTime = now
	} else {
		gapSeconds := now - self.status.TiliRTime
		secondsPerTili := design_config.GetSecondsPerTili()

		// 计算一下回复的count，不信任客户端传来的count
		count := gapSeconds / secondsPerTili
		before := self.status.Tili
		if count > 0 {
			change = true
			if self.status.Tili+int(count) <= self.status.MaxTili {
				self.status.Tili = self.status.Tili + int(count)
			} else if self.status.Tili < self.status.MaxTili {
				self.status.Tili = self.status.MaxTili
			}
			self.status.TiliRTime = self.status.TiliRTime + count*secondsPerTili
			//gm log
			c := CurrencyChangeExtra{
				Nickname:      self.Nickname,
				Level:         self.status.Level,
				Before:        float64(before),
				Amount:        float64(count),
				Left:          float64(self.status.Tili),
				DeviceId:      deviceId,
				ServerAddress: appcfg.GetAddress(),
			}
			data, eerr := json.Marshal(c)
			if eerr != nil {
				Err(eerr)
			}
			GMLog(constants.C1_PLAYER, constants.C2_TILI, constants.C3_SYSTEM_RECOVER, self.Username, string(data))
		}
	}

	return
}

//add tili even when current tili is full
func (self *GFPlayer) AddTili(count int, diamond int, deviceId string, src string) (newTili int, newDiamond int, log bool) {
	before := self.status.Tili
	if count != 0 {
		self.status.Tili = self.status.Tili + count

		//gm log
		c := CurrencyChangeExtra{
			Nickname:      self.Nickname,
			Level:         self.status.Level,
			Before:        float64(before),
			Amount:        float64(count),
			Left:          float64(self.status.Tili),
			DeviceId:      deviceId,
			ServerAddress: appcfg.GetAddress(),
		}
		data, _ := json.Marshal(c)
		GMLog(constants.C1_PLAYER, constants.C2_TILI, src, self.Username, string(data))
	}

	newTili = self.status.Tili

	return
}

func (self *GFPlayer) UseTili(count int, deviceId string) {
	before := self.status.Tili

	self.status.Tili = self.status.Tili - count

	//gm log
	c := CurrencyChangeExtra{
		Nickname:      self.Nickname,
		Level:         self.status.Level,
		Before:        float64(before),
		Amount:        float64(-count),
		Left:          float64(self.status.Tili),
		DeviceId:      deviceId,
		ServerAddress: appcfg.GetAddress(),
	}
	data, eerr := json.Marshal(c)
	if eerr != nil {
		Err(eerr)
	}
	GMLog(constants.C1_PLAYER, constants.C2_TILI, constants.C3_DEFAULT, self.Username, string(data))
}

//MISSION gift
func (self *GFPlayer) PlusTili(count int, deviceId string) int {
	before := self.status.Tili
	if self.status.Tili < self.status.MaxTili {
		if self.status.Tili+count <= self.status.MaxTili {
			self.status.Tili = self.status.Tili + count
		} else {
			self.status.TiliRTime = time.Now().Unix()
			self.status.Tili = self.status.MaxTili
		}
	}

	//gm log
	c := CurrencyChangeExtra{
		Nickname:      self.Nickname,
		Level:         self.status.Level,
		Before:        float64(before),
		Amount:        float64(count),
		Left:          float64(self.status.Tili),
		DeviceId:      deviceId,
		ServerAddress: appcfg.GetAddress(),
	}
	data, eerr := json.Marshal(c)
	if eerr != nil {
		Err(eerr)
	}
	GMLog(constants.C1_PLAYER, constants.C2_TILI, constants.C3_MISSION, self.Username, string(data))

	return self.status.Tili
}

func (self *GFPlayer) SetHead(newHead string) (changed bool) {
	if self.status.Head == newHead && self.status.ExtraInfo1.KakaoHead == newHead {
		return false
	} else {
		changed = true
		self.status.Head = newHead
		self.status.ExtraInfo1.KakaoHead = newHead
	}

	return
}

func (self *GFPlayer) SetFirstIAPReward(diamondID int) {
	if self.status.DayExtraInfo.FirstRewardMap == nil {
		self.status.DayExtraInfo.FirstRewardMap = make(map[string]int)
	}
	self.status.DayExtraInfo.FirstRewardMap[strconv.Itoa(diamondID)] = 1

	return
}

func (self *GFPlayer) SetVerifyBenefit(num int) {
	self.status.ExtraInfo1.VerifyBenefit += num
}

func (self *GFPlayer) SetYesterdayBenefit(num int) {
	self.status.ExtraInfo1.YesterdayBenefit += num
}

func (self *GFPlayer) IncreVerify(key string) {
	if self.status.DayExtraInfo.DayVerifyNum == nil {
		self.status.DayExtraInfo.DayVerifyNum = make(map[string]int)
	}

	self.status.DayExtraInfo.DayVerifyNum[key] += 1
	self.status.DayExtraInfo.VerifyCount += 1
}

func (self *GFPlayer) AddPixieClothes(origin ClothesOrigin, paperID int64, count int, nt int64, reason string) (cc ClothesInfo) {
	cid := fmt.Sprintf("%s-%d", origin, paperID)

	return self.AddPixieClothesWithCID(cid, count, nt, reason)
}

func (self *GFPlayer) AddPixieClothesWithCID(cid string, count int, nt int64, reason string) (cc ClothesInfo) {
	cc.ClothesID = cid

	if v, ok := self.status.Clothes[cid]; ok {
		v.Count += count
		v.Time = nt

		cc.Count = v.Count
		cc.Time = nt
	} else {
		self.status.Clothes[cid] = &ClothesDetail{Count: count, Time: nt}

		cc.Count = count
		cc.Time = nt
	}

	return
}

func (self *GFPlayer) DeleteClothes(clothesID string, cnt int, source string, deviceId string) (clothesInfo *ClothesInfo, succeed bool) {
	changed := false
	for cid, v := range self.status.Clothes {
		if cid == clothesID && v.Count >= cnt {
			changed = true
			v.Count -= cnt

			clothesInfo = &ClothesInfo{ClothesID: cid, Count: v.Count, Time: v.Time}

			//gm log
			c := DeltaClothesExtra{
				ClothesID:     clothesID,
				Amount:        -cnt,
				Left:          v.Count,
				DeviceId:      deviceId,
				ServerAddress: appcfg.GetAddress(),
			}
			data, _ := json.Marshal(c)
			GMLog(constants.C1_ITEM, constants.C2_CLOTHES, source, self.Username, string(data))
			break
		}
	}

	if changed {
		succeed = true
	} else {
		succeed = false
	}
	return
}

func (self *GFPlayer) ViewRecord(levelID string) {
	if self.status.RecordMap == nil {
		self.status.RecordMap = make(map[string]*PixieRecord)
	}

	nr := self.status.RecordMap[levelID]

	if nr == nil {
		//new record
		self.status.RecordMap[levelID] = &PixieRecord{
			PlotRead:    true,
			HistoryRank: make(map[string]int),
		}
	}
}

func (self *GFPlayer) SaveRecord(recordID, rank, clothes string, score float64, fnpc, deviceId string, costTili int) {
	if self.status.RecordMap == nil {
		self.status.RecordMap = make(map[string]*PixieRecord)
	}

	//save record
	nr := self.status.RecordMap[recordID]

	if nr == nil {
		//new record
		self.status.RecordMap[recordID] = &PixieRecord{
			PlotRead:    true,
			Score:       score,
			BestScore:   score,
			Rank:        rank,
			Clothes:     clothes,
			BestClothes: clothes,
			HistoryRank: map[string]int{rank: 1},
		}
	} else {
		nr.PlotRead = true
		nr.Score = score
		if nr.BestScore < score {
			nr.BestScore = score
			nr.BestClothes = clothes
		}
		nr.Rank = rank
		nr.Clothes = clothes
		if nr.HistoryRank[rank] <= 0 {
			nr.HistoryRank[rank] = 1
		}
	}

	//add npc friendship exp
	if fnpc != "" {
		addExp := design_config.GetNpcFriendshipExpAdd(rank)
		self.AddFriendshipExp(fnpc, addExp)
	}

	//add player exp
	self.AddExp(design_config.GetExpAdd(rank), constants.SRC_MISSION)

	//use tili
	self.UseTili(costTili, deviceId)
}

func (self *GFPlayer) Ban(st, et int64, reason string) {
	self.status.BanInfo.BanStartTime = st
	self.status.BanInfo.BanEndTime = et
	self.status.BanInfo.BanReason = reason
}

func (self *GFPlayer) SetIAPReturned() {
	self.status.ExtraInfo1.IAPReturned = true
}

func (self *GFPlayer) SetYuyueRewarded() {
	self.status.ExtraInfo1.YuyueRewarded = true
}

func (self *GFPlayer) AddDayQinmiAddBoardCnt(fu string) {
	if self.status.DayExtraInfo.DayQinmiAddBoardMap == nil {
		self.status.DayExtraInfo.DayQinmiAddBoardMap = make(map[string]int)
	}
	self.status.DayExtraInfo.DayQinmiAddBoardMap[fu]++
}

func (self *GFPlayer) AddDayQinmiReplyBoardCnt(fu string) {
	if self.status.DayExtraInfo.DayQinmiReplyBoardMap == nil {
		self.status.DayExtraInfo.DayQinmiReplyBoardMap = make(map[string]int)
	}
	self.status.DayExtraInfo.DayQinmiReplyBoardMap[fu]++
}

func (self *GFPlayer) AddDayTheaterCnt(theaterID int, cnt int) int {
	if self.status.DayExtraInfo.DayTheaterCntMap == nil {
		self.status.DayExtraInfo.DayTheaterCntMap = make(map[string]int)
	}

	self.status.DayExtraInfo.DayTheaterCntMap[strconv.Itoa(theaterID)] += cnt
	return self.status.DayExtraInfo.DayTheaterCntMap[strconv.Itoa(theaterID)]
}

func (self *GFPlayer) AddDayCopyReportCnt(cnt int) int {
	self.status.DayExtraInfo.DayCopyReportCnt += cnt
	return self.status.DayExtraInfo.DayCopyReportCnt
}

func (self *GFPlayer) AddDayShensuCnt(cid string, cnt int) int {
	if self.status.DayExtraInfo.DayShensuMap == nil {
		self.status.DayExtraInfo.DayShensuMap = make(map[string]int)
	}
	self.status.DayExtraInfo.DayShensuMap[cid] += cnt
	return self.status.DayExtraInfo.DayShensuMap[cid]
}

func (self *GFPlayer) SetDayFreeLottery(poolID, typee int) {
	if self.status.DayExtraInfo.DayFreeLotteryMap == nil {
		self.status.DayExtraInfo.DayFreeLotteryMap = make(map[string]bool)
	}

	self.status.DayExtraInfo.DayFreeLotteryMap[fmt.Sprintf("%d:%d", poolID, typee)] = true
}

func (self *GFPlayer) AddGuanzhuDesignerCnt(cnt int) int {
	self.status.ExtraInfo1.GuanzhuDesignerCount += cnt
	return self.status.ExtraInfo1.GuanzhuDesignerCount
}

func (self *GFPlayer) AddMarkRewardGot(rt int) {
	if self.status.DayExtraInfo.MarkGotRewards == nil {
		self.status.DayExtraInfo.MarkGotRewards = make([]int, 0)
	}

	self.status.DayExtraInfo.MarkGotRewards = append(self.status.DayExtraInfo.MarkGotRewards, rt)
}

func (self *GFPlayer) SetGuildBoardReadID(id int64) {
	self.status.ExtraInfo1.LastReadGuildBoardID = id
}

func (self *GFPlayer) SetIgnoreNotice(value bool) {
	self.status.ExtraInfo1.IgnoreNotice = value
}

func (self *GFPlayer) SetPushConfig(value bool, ignoreGuildWarPush, ignoreGuildPush, acceptFriendRequestPush, ignoreHopeDonePush, acceptPartyInvitePush, ignoreNewCustomPush, ignoreCartPush bool) {
	self.status.ExtraInfo1.IgnorePush = value

	self.status.ExtraInfo1.IgnoreGuildWarPush = ignoreGuildWarPush
	self.status.ExtraInfo1.IgnoreGuildPush = ignoreGuildPush
	self.status.ExtraInfo1.AcceptFriendRequestPush = acceptFriendRequestPush
	self.status.ExtraInfo1.IgnoreHopeDonePush = ignoreHopeDonePush
	self.status.ExtraInfo1.AcceptPartyInvitePush = acceptPartyInvitePush
	self.status.ExtraInfo1.IgnoreNewCustomPush = ignoreNewCustomPush
	self.status.ExtraInfo1.IgnoreCartPush = ignoreCartPush
}

func (self *GFPlayer) SetIgnorePushKor(ignorePush, ignoreNightPush bool) {
	self.status.ExtraInfo1.IgnorePush = ignorePush
	self.status.ExtraInfo1.IgnoreNightPush = ignoreNightPush
}

func (self *GFPlayer) AddBackground(bid string) bool {
	if self.status.ExtraInfo1.Backgrounds == nil {
		self.status.ExtraInfo1.Backgrounds = make([]string, 0)
	}

	for _, i := range self.status.ExtraInfo1.Backgrounds {
		if i == bid {
			return false
		}
	}

	self.status.ExtraInfo1.Backgrounds = append(self.status.ExtraInfo1.Backgrounds, bid)
	return true
}

//CheckDesingCoin 检查设计师币
func (self *GFPlayer) CheckDesingCoin(Designcoin int) bool {
	if Designcoin >= 0 && self.status.ExtraInfo1.DesignCoin >= Designcoin {
		return true
	} else if Designcoin < 0 && self.status.ExtraInfo1.DesignCoin >= -Designcoin {
		return true
	} else {
		return false
	}
}

func (self *GFPlayer) SetKorInitPlayerRewarded(value bool) {
	self.status.ExtraInfo1.KorInitPlayerRewarded = value
}

func (self *GFPlayer) SetKorRecordCZRewarded(value bool) {
	self.status.ExtraInfo1.KorRecordCZRewarded = value
}

func (self *GFPlayer) BuyTheaterCnt(tid int) {
	if self.status.DayExtraInfo.DayTheaterCntMap == nil {
		self.status.DayExtraInfo.DayTheaterCntMap = make(map[string]int)
	}

	if self.status.DayExtraInfo.DayBuyTheaterCntMap == nil {
		self.status.DayExtraInfo.DayBuyTheaterCntMap = make(map[string]int)
	}

	self.status.DayExtraInfo.DayBuyTheaterCntMap[strconv.Itoa(tid)] += 1
	self.status.DayExtraInfo.DayTheaterCntMap[strconv.Itoa(tid)] = 0
}

func (self *GFPlayer) SetMonthCardRewarded(month string) {
	if self.status.ExtraInfo1.MonthCardRewardMap == nil {
		self.status.ExtraInfo1.MonthCardRewardMap = make(map[string]int)
	}

	self.status.ExtraInfo1.MonthCardRewardMap[month] = 1
}

func (self *GFPlayer) AddDesignCoinPriceToken(cnt int) int {
	self.status.ExtraInfo1.DesignCoinPriceToken += cnt
	return self.status.ExtraInfo1.DesignCoinPriceToken
}

func (self *GFPlayer) AddDesignCoinPriceTokenBuyCnt(month string, cnt int) int {
	if self.status.ExtraInfo1.DesignCoinPriceTokenBuyMap == nil {
		self.status.ExtraInfo1.DesignCoinPriceTokenBuyMap = make(map[string]int)
	}
	self.status.ExtraInfo1.DesignCoinPriceTokenBuyMap[month] += cnt
	return self.status.ExtraInfo1.DesignCoinPriceTokenBuyMap[month]
}

func (self *GFPlayer) AddAttendParty(hostID int64, now time.Time) (new bool) {
	if self.status.ExtraInfo1.AttendRTParty == nil {
		self.status.ExtraInfo1.AttendRTParty = make(map[string]int64)
	}

	token := fmt.Sprintf("%d", hostID)
	if self.status.ExtraInfo1.AttendRTParty[token] > 0 {
		return false
	} else {
		self.status.ExtraInfo1.AttendRTParty[token] = now.Unix()
		return true
	}
}

func (self *GFPlayer) RemoveAttendParty(hostID int64) (change bool) {
	if self.status.ExtraInfo1.AttendRTParty == nil {
		self.status.ExtraInfo1.AttendRTParty = make(map[string]int64)
		return false
	}

	token := fmt.Sprintf("%d", hostID)
	if _, ok := self.status.ExtraInfo1.AttendRTParty[token]; ok {
		change = true
		delete(self.status.ExtraInfo1.AttendRTParty, token)
	}
	return
}

func (self *GFPlayer) AddCurrentDayRTPartyReward(now time.Time, gold, diamond, wire, exp int) (int, int, int, int) {
	if self.status.ExtraInfo1.DayRTPartyRewardMap == nil {
		self.status.ExtraInfo1.DayRTPartyRewardMap = make(map[string]int)
	}
	ds := now.Format("20060102")
	self.status.ExtraInfo1.DayRTPartyRewardMap[ds+"_gold"] += gold
	self.status.ExtraInfo1.DayRTPartyRewardMap[ds+"_diamond"] += diamond
	self.status.ExtraInfo1.DayRTPartyRewardMap[ds+"_wire"] += wire
	self.status.ExtraInfo1.DayRTPartyRewardMap[ds+"_exp"] += exp

	return self.status.ExtraInfo1.DayRTPartyRewardMap[ds+"_gold"], self.status.ExtraInfo1.DayRTPartyRewardMap[ds+"_diamond"], self.status.ExtraInfo1.DayRTPartyRewardMap[ds+"_wire"], self.status.ExtraInfo1.DayRTPartyRewardMap[ds+"_exp"]
}

func (self *GFPlayer) SetETHPayPwdVerifyInfo(email string, code string, codeET int64) {
	self.status.ExtraInfo1.ETHPayPwdVerifyCodeEmailSend = email
	self.status.ExtraInfo1.ETHPayPwdVerifyCode = code
	self.status.ExtraInfo1.ETHPayPwdVerifyCodeExpireTime = codeET
}

func (self *GFPlayer) SetETHPayPwd(email string, pwdExcrypt string) {
	self.status.ExtraInfo1.ETHPayPwdEmail = email
	self.status.ExtraInfo1.ETHPayPwdEncrypt = pwdExcrypt
}
