package player_mem

import (
	"db_pixie/land"
	"math"
	"math/rand"
	. "pixie_contract/api_specification"
	"service_pixie/building"
	"service_pixie/demand"
	"service_pixie/design_config"
	"service_pixie/level"
	"time"
)

const (
	STREET_INCOME_DATE_FORMAT = "20060102-1504"
)

func (self *GFPlayer) SetStreetClean(cleanObjectID string, changeCnt int) (success bool) {
	if self.status.StreetOperate.CleanPower == nil {
		self.status.StreetOperate.CleanPower = make(map[string]*CleanObjectDetail)
	}

	if self.status.StreetOperate.CleanPower[cleanObjectID] == nil {
		self.status.StreetOperate.CleanPower[cleanObjectID] = &CleanObjectDetail{
			Level: 1,
		}
	}

	toCnt := self.status.StreetOperate.CleanPower[cleanObjectID].Count + changeCnt
	if toCnt >= 0 && toCnt <= building.GetCleanObjectMaxEmployCount(cleanObjectID) {
		success = true
		self.status.StreetOperate.CleanPower[cleanObjectID].Count += changeCnt

		//会影响街区的参数

	}

	return
}

func (self *GFPlayer) ClearStreetWallet() int64 {
	cnt := int64(self.status.StreetOperate.WalletLeft)
	self.status.StreetOperate.WalletLeft -= float64(cnt)

	return cnt
}

func (self *GFPlayer) GetStreetDetail() *StreetOperateResult {
	return self.status.StreetOperate
}

func (self *GFPlayer) initStreet(lvl int, ntInMinute string) (change bool, err error) {
	if self.status.StreetOperate == nil {
		visitCapParam := self.getVisitCapParam()
		if profitInMinute, generalDemandParam, myShopIncomeInMinute, entertainIncomeInMinute, restaurantInComeInMinute, cleanSalaryInMinute, dd, e := self.computeStreetIncome(visitCapParam, nil); e == nil {
			self.status.StreetOperate = &StreetOperateResult{
				LastIncomeMinute: ntInMinute,
				VisitCap:         visitCapParam,
				WalletMax:        level.GetLevelWalletMax(lvl),
				CleanPower:       make(map[string]*CleanObjectDetail),

				StreetOperateDemandDetail: dd,

				GeneralParam: generalDemandParam,

				ProfitInHour:           profitInMinute * 60,
				MyShopIncomeInHour:     myShopIncomeInMinute * 60,
				EntertainIncomeInHour:  entertainIncomeInMinute * 60,
				RestaurantIncomeInHour: restaurantInComeInMinute * 60,
				CleanSalaryInHour:      cleanSalaryInMinute * 60,
			}

			change = true
		} else {
			err = e
		}
	}

	return
}

//更新街区各参数、计算收入
func (self *GFPlayer) RefreshStreet(force bool) (change bool) {
	//当前时间信息
	nt := time.Now()
	ntInMinute := nt.Format(STREET_INCOME_DATE_FORMAT)

	//玩家当前等级
	lvl := self.GetLevel()

	if c, e := self.initStreet(lvl, ntInMinute); c || e != nil {
		return
	}

	//上次计算收入的时间，精确到分
	lastComputeTime, _ := time.ParseInLocation(STREET_INCOME_DATE_FORMAT, self.status.StreetOperate.LastIncomeMinute, time.Local)

	//人流容量
	var visitCapParam float64

	//检查人流容量是否需要更新
	if lastComputeTime.Before(nt) && ntInMinute != self.status.StreetOperate.LastIncomeMinute {
		visitCapParam = self.getVisitCapParam()
	} else {
		visitCapParam = self.status.StreetOperate.VisitCap
	}

	if force || (lastComputeTime.Before(nt) && ntInMinute != self.status.StreetOperate.LastIncomeMinute) {
		//街区需要更新
		if profitInMinute, generalDemandParam, myShopIncomeInMinute, entertainIncomeInMinute, restaurantInComeInMinute, cleanSalaryInMinute, dd, err := self.computeStreetIncome(visitCapParam, self.status.StreetOperate.CleanPower); err == nil {
			change = true

			for lastComputeTime.Before(nt) && ntInMinute != self.status.StreetOperate.LastIncomeMinute {
				//递增时间
				lastComputeTime = lastComputeTime.Add(time.Minute)

				//为钱包增加收入
				self.status.StreetOperate.WalletLeft += profitInMinute

				//退出计算条件
				if self.status.StreetOperate.WalletLeft >= self.status.StreetOperate.WalletMax {
					self.status.StreetOperate.WalletLeft = self.status.StreetOperate.WalletMax
					break
				}

				if self.status.StreetOperate.WalletLeft <= 0 {
					self.status.StreetOperate.WalletLeft = 0
					break
				}
			}

			//更新街区状态
			self.status.StreetOperate.VisitCap = visitCapParam
			self.status.StreetOperate.LastIncomeMinute = ntInMinute
			self.status.StreetOperate.GeneralParam = generalDemandParam
			self.status.StreetOperate.ProfitInHour = profitInMinute * 60
			self.status.StreetOperate.MyShopIncomeInHour = myShopIncomeInMinute * 60
			self.status.StreetOperate.EntertainIncomeInHour = entertainIncomeInMinute * 60
			self.status.StreetOperate.RestaurantIncomeInHour = restaurantInComeInMinute * 60
			self.status.StreetOperate.CleanSalaryInHour = cleanSalaryInMinute * 60
			self.status.StreetOperate.StreetOperateDemandDetail = dd
		}
	}

	return
}

//获取人流容量参数
func (self *GFPlayer) getVisitCapParam() float64 {
	//人流容量参数
	visitCapParamBase := level.GetLevelVisitCap(self.GetLevel())
	visitCapParamRand := float64(90+rand.Intn(21)) * 0.01

	return visitCapParamBase * visitCapParamRand * (100.0 + self.status.NpcBuff.BuffStreetVisitCap) / 100.0
}

//获得街的每分钟收入值
func (self *GFPlayer) computeStreetIncome(visitCapParam float64, cleanPowerMap map[string]*CleanObjectDetail) (profitInMinute, generalDemandParam, myShopIncomeInMinute, entertainIncomeInMinute, restaurantInComeInMinute, cleanSalaryInMinute float64, dd StreetOperateDemandDetail, err error) {
	//卫生需求 考虑buff
	cleanNeed := demand.GetVisitCapCleanNeed(visitCapParam) * getActualBuff(self.status.NpcBuff.BuffAllCleanNeed)

	//其他需求 无buff
	entertainNeed := demand.GetVisitCapEntertainNeed(visitCapParam)
	restaurantNeed := demand.GetVisitCapRestaurantNeed(visitCapParam)

	dd = StreetOperateDemandDetail{
		//初始能力
		EntertainOffer:  design_config.GetEntertainOfferInit(),
		RestaurantOffer: design_config.GetRestaurantOfferInit(),

		//需求
		EntertainNeed:  entertainNeed,
		RestaurantNeed: restaurantNeed,

		//卫生
		CleanNeed:  cleanNeed,
		CleanOffer: design_config.GetCleanOfferInit(),
	}

	//参与盈利的各个建筑
	var myShopProfit, entertainProfit, restaurantProfit float64
	var entertainBuildingID, restaurantBuildingID int64
	var lands []Land
	if lands, err = land.ListPlayerLand(self.Username); err == nil {
		for _, l := range lands {
			if bd, exist := building.GetEntertainBuildingByIDLevel(l.BuildingID, l.BuildingLevel); exist && entertainBuildingID == 0 {
				//娱乐建筑
				entertainBuildingID = bd.ID
				dd.EntertainLandID = l.ID

				//引入buff
				buffOffer, buffCleanNeed, buffProfit := self.getEntertainBuff(bd.ID)
				offer := bd.Offer * buffOffer
				cleanNeed := bd.CleanNeed * buffCleanNeed
				profit := bd.Profit * buffProfit

				//本建筑的固定效果
				dd.EntertainParamOpen = offer / dd.EntertainNeed
				if bdn, exist := building.GetEntertainBuildingByIDLevel(l.BuildingID, l.BuildingLevel+1); exist {
					noffer := bdn.Offer * buffOffer
					dd.EntertainParamOpenNextLevel += noffer / dd.EntertainNeed
				} else {
					dd.EntertainParamOpenNextLevel += offer / dd.EntertainNeed
				}

				//本建筑的实际效果
				if l.Status == int(LAND_STATUS_WB_IN_BUSINESS) {
					dd.EntertainOffer += offer
					dd.CleanNeed += cleanNeed
					entertainProfit = profit
				}
			} else if bd, exist := building.GetRestaurantBuildingByIDLevel(l.BuildingID, l.BuildingLevel); exist && restaurantBuildingID == 0 {
				//餐饮建筑
				restaurantBuildingID = bd.ID
				dd.RestaurantLandID = l.ID

				//引入buff
				buffOffer, buffCleanNeed, buffProfit := self.getRestaurantBuff(bd.ID)
				offer := bd.Offer * buffOffer
				cleanNeed := bd.CleanNeed * buffCleanNeed
				profit := bd.Profit * buffProfit

				//本建筑的固定效果
				dd.RestaurantParamOpen += offer / dd.RestaurantNeed
				if bdn, exist := building.GetRestaurantBuildingByIDLevel(l.BuildingID, l.BuildingLevel+1); exist {
					noffer := bdn.Offer * buffOffer
					dd.RestaurantParamOpenNextLevel += noffer / dd.RestaurantNeed
				} else {
					dd.RestaurantParamOpenNextLevel += offer / dd.RestaurantNeed
				}

				//本建筑的实际效果
				if l.Status == int(LAND_STATUS_WB_IN_BUSINESS) {
					dd.RestaurantOffer += offer
					dd.CleanNeed += cleanNeed
					restaurantProfit = profit
				}
			} else if l.Status == int(LAND_STATUS_WB_IN_BUSINESS) && l.Type == PIXIE_LAND_TYPE_MY_SHOP {
				//我的默认店铺 且在营业中
				if bd, exist := building.GetBuildingByID(l.BuildingID); exist {
					myShopProfit = bd.Profit * self.getMyShopProfitBuff()
				}
			}
		}
	}

	if len(cleanPowerMap) > 0 {
		for coid, v := range cleanPowerMap {
			//buff
			bp, bs := self.getCleanObjectBuff(coid)

			p, sim, _ := building.GetCleanObjectDetail(coid)

			dd.CleanOffer += p * bp * float64(v.Count)
			cleanSalaryInMinute += sim * bs * float64(v.Count)
		}
	}

	cleanParam := paramCheck(dd.CleanOffer / dd.CleanNeed)
	entertainParam := paramCheck(dd.EntertainOffer / dd.EntertainNeed)
	restaurantParam := paramCheck(dd.RestaurantOffer / dd.RestaurantNeed)
	demandParam := (entertainParam + restaurantParam) * 0.5
	generalDemandParam = math.Sqrt(demandParam * cleanParam)

	profitParam := visitCapParam * demandParam * cleanParam * design_config.GetStreetIncomeFinalAdjust()

	myShopIncomeInMinute = myShopProfit * profitParam
	entertainIncomeInMinute = entertainProfit * profitParam
	restaurantInComeInMinute = restaurantProfit * profitParam

	//盈利=三类建筑收入 - 卫生支出
	profitInMinute = myShopIncomeInMinute + entertainIncomeInMinute + restaurantInComeInMinute - cleanSalaryInMinute
	return
}

func getActualBuff(buff float64) float64 {
	return (100.0 + buff) / 100.0
}

func paramCheck(p float64) float64 {
	if p < 0.3 {
		return 0.3
	} else if p > 1.2 {
		return 1.2
	}
	return p
}
