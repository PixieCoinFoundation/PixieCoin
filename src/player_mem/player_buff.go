package player_mem

import (
	"constants"
	. "logger"
	. "pixie_contract/api_specification"
	"service_pixie/npc"
	"strconv"
	"strings"
)

func (self *GFPlayer) getCleanObjectBuff(coid string) (buffPower, buffSalary float64) {
	if coid == "101" {
		buffPower = (100.0 + self.status.NpcBuff.BuffCleanObject101Power + self.status.NpcBuff.BuffAllCleanPower) / 100.0
		buffSalary = (100.0 + self.status.NpcBuff.BuffCleanObject101Salary + self.status.NpcBuff.BuffAllCleanSalary) / 100.0
	} else if coid == "201" {
		buffPower = (100.0 + self.status.NpcBuff.BuffCleanObject201Power + self.status.NpcBuff.BuffAllCleanPower) / 100.0
		buffSalary = (100.0 + self.status.NpcBuff.BuffCleanObject201Salary + self.status.NpcBuff.BuffAllCleanSalary) / 100.0
	} else {
		buffPower = 1.0
		buffSalary = 1.0
	}
	return
}

func (self *GFPlayer) getMyShopProfitBuff() (buffProfit float64) {
	return (100.0 + self.status.NpcBuff.BuffMyShopProfit + self.status.NpcBuff.BuffAllProfit) / 100
}

func (self *GFPlayer) getRestaurantBuff(bid int64) (buffOffer, buffCleanNeed, buffProfit float64) {
	br := self.status.NpcBuff

	bo := br.BuffRestaurantOffer
	if br.BuffSpecifiedRestaurantOffer != nil {
		bo += br.BuffSpecifiedRestaurantOffer[bid]
	}

	bcn := br.BuffRestaurantCleanNeed + br.BuffAllCleanNeed
	if br.BuffSpecifiedRestaurantCleanNeed != nil {
		bcn += br.BuffSpecifiedRestaurantCleanNeed[bid]
	}

	bp := br.BuffRestaurantProfit + br.BuffAllProfit
	if br.BuffSpecifiedRestaurantProfit != nil {
		bp += br.BuffSpecifiedRestaurantProfit[bid]
	}

	buffOffer = (100.0 + bo) / 100.0
	buffCleanNeed = (100.0 + bcn) / 100.0
	buffProfit = (100.0 + bp) / 100.0
	return
}

func (self *GFPlayer) getEntertainBuff(bid int64) (buffOffer, buffCleanNeed, buffProfit float64) {
	br := self.status.NpcBuff

	bo := br.BuffEntertainOffer
	if br.BuffSpecifiedEntertainOffer != nil {
		bo += br.BuffSpecifiedEntertainOffer[bid]
	}

	bcn := br.BuffEntertainCleanNeed + br.BuffAllCleanNeed
	if br.BuffSpecifiedEntertainCleanNeed != nil {
		bcn += br.BuffSpecifiedEntertainCleanNeed[bid]
	}

	bp := br.BuffEntertainProfit + br.BuffAllProfit
	if br.BuffSpecifiedEntertainProfit != nil {
		bp += br.BuffSpecifiedEntertainProfit[bid]
	}

	buffOffer = (100.0 + bo) / 100.0
	buffCleanNeed = (100.0 + bcn) / 100.0
	buffProfit = (100.0 + bp) / 100.0
	return
}

func (self *GFPlayer) refreshNpcBuffResult() {
	br := NpcBuffResult{
		NpcBuffList:                      self.status.NpcBuff.NpcBuffList,
		BuffSpecifiedRestaurantOffer:     make(map[int64]float64),
		BuffSpecifiedEntertainOffer:      make(map[int64]float64),
		BuffSpecifiedRestaurantProfit:    make(map[int64]float64),
		BuffSpecifiedEntertainProfit:     make(map[int64]float64),
		BuffSpecifiedRestaurantCleanNeed: make(map[int64]float64),
		BuffSpecifiedEntertainCleanNeed:  make(map[int64]float64),
		BuffUnlockBuilding:               make(map[int64]bool),
	}

	for npcID, fd := range self.status.NpcFriendship {
		tmpLevel := fd.Level
		for tmpLevel > 1 {
			b1, b1p1, b1p2, b2, b2p1, b2p2 := npc.GetNpcBuff(npcID, tmpLevel)

			refreshBuff1(b1, b1p1, b1p2, &br)
			refreshBuff1(b2, b2p1, b2p2, &br)

			tmpLevel--
		}
	}

	self.status.NpcBuff = br
}

func refreshBuff1(bt, bp1, bp2 string, br *NpcBuffResult) {
	if bt == "" {
		return
	}

	switch bt {
	case constants.BUFF_STREET_VISIT_CAP:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_STREET_VISIT_CAP 0")
		} else {
			br.BuffStreetVisitCap += value
		}
	case constants.BUFF_RESTAURANT_OFFER:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_RESTAURANT_OFFER 0")
		} else {
			br.BuffRestaurantOffer += value
		}
	case constants.BUFF_ENTERTAIN_OFFER:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_ENTERTAIN_OFFER 0")
		} else {
			br.BuffEntertainOffer += value
		}
	case constants.BUFF_CLEAN_OBJECT_101_POWER:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_CLEAN_OBJECT_101_POWER 0")
		} else {
			br.BuffCleanObject101Power += value
		}
	case constants.BUFF_CLEAN_OBJECT_201_POWER:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_CLEAN_OBJECT_201_POWER 0")
		} else {
			br.BuffCleanObject201Power += value
		}
	case constants.BUFF_ALL_CLEAN_POWER:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_ALL_CLEAN_POWER 0")
		} else {
			br.BuffAllCleanPower += value
		}
	case constants.BUFF_CLEAN_OBJECT_101_SALARY:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_CLEAN_OBJECT_101_SALARY 0")
		} else {
			br.BuffCleanObject101Salary += value
		}
	case constants.BUFF_CLEAN_OBJECT_201_SALARY:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_CLEAN_OBJECT_201_SALARY 0")
		} else {
			br.BuffCleanObject201Salary += value
		}
	case constants.BUFF_ALL_CLEAN_SALARY:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_ALL_CLEAN_SALARY 0")
		} else {
			br.BuffAllCleanSalary += value
		}
	case constants.BUFF_RESTAURANT_PROFIT:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_RESTAURANT_PROFIT 0")
		} else {
			br.BuffRestaurantProfit += value
		}
	case constants.BUFF_ENTERTAIN_PROFIT:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_ENTERTAIN_PROFIT 0")
		} else {
			br.BuffEntertainProfit += value
		}
	case constants.BUFF_MY_SHOP_PROFIT:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_MY_SHOP_PROFIT 0")
		} else {
			br.BuffMyShopProfit += value
		}
	case constants.BUFF_ALL_PROFIT:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_ALL_PROFIT 0")
		} else {
			br.BuffAllProfit += value
		}
	case constants.BUFF_RESTAURANT_CLEAN_NEED:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_RESTAURANT_CLEAN_NEED 0")
		} else {
			br.BuffRestaurantCleanNeed += value
		}
	case constants.BUFF_ENTERTAIN_CLEAN_NEED:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_ENTERTAIN_CLEAN_NEED 0")
		} else {
			br.BuffEntertainCleanNeed += value
		}
	case constants.BUFF_ALL_CLEAN_NEED:
		if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
			Err(e, "BUFF_ALL_CLEAN_NEED 0")
		} else {
			br.BuffAllCleanNeed += value
		}
	case constants.BUFF_SPECIFIED_RESTAURANT_OFFER:
		bidList := strings.Split(bp2, ",")
		if len(bidList) > 0 {
			for _, bids := range bidList {
				if bid, e := strconv.ParseInt(bids, 10, 64); e != nil || bid <= 0 {
					Err(e, "BUFF_SPECIFIED_RESTAURANT_OFFER bid 0")
				} else {
					if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
						Err(e, "BUFF_SPECIFIED_RESTAURANT_OFFER 0")
					} else {
						br.BuffSpecifiedRestaurantOffer[bid] += value
					}
				}
			}
		} else {
			Err("BUFF_SPECIFIED_RESTAURANT_OFFER size 0")
		}
	case constants.BUFF_SPECIFIED_ENTERTAIN_OFFER:
		bidList := strings.Split(bp2, ",")
		if len(bidList) > 0 {
			for _, bids := range bidList {
				if bid, e := strconv.ParseInt(bids, 10, 64); e != nil || bid <= 0 {
					Err(e, "BUFF_SPECIFIED_ENTERTAIN_OFFER bid 0")
				} else {
					if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
						Err(e, "BUFF_SPECIFIED_ENTERTAIN_OFFER 0")
					} else {
						br.BuffSpecifiedEntertainOffer[bid] += value
					}
				}
			}
		} else {
			Err("BUFF_SPECIFIED_ENTERTAIN_OFFER size 0")
		}
	case constants.BUFF_SPECIFIED_RESTAURANT_PROFIT:
		bidList := strings.Split(bp2, ",")
		if len(bidList) > 0 {
			for _, bids := range bidList {
				if bid, e := strconv.ParseInt(bids, 10, 64); e != nil || bid <= 0 {
					Err(e, "BUFF_SPECIFIED_RESTAURANT_PROFIT bid 0")
				} else {
					if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
						Err(e, "BUFF_SPECIFIED_RESTAURANT_PROFIT 0")
					} else {
						br.BuffSpecifiedRestaurantProfit[bid] += value
					}
				}
			}
		} else {
			Err("BUFF_SPECIFIED_RESTAURANT_PROFIT size 0")
		}
	case constants.BUFF_SPECIFIED_ENTERTAIN_PROFIT:
		bidList := strings.Split(bp2, ",")
		if len(bidList) > 0 {
			for _, bids := range bidList {
				if bid, e := strconv.ParseInt(bids, 10, 64); e != nil || bid <= 0 {
					Err(e, "BUFF_SPECIFIED_ENTERTAIN_PROFIT bid 0")
				} else {
					if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
						Err(e, "BUFF_SPECIFIED_ENTERTAIN_PROFIT 0")
					} else {
						br.BuffSpecifiedEntertainProfit[bid] += value
					}
				}
			}
		} else {
			Err("BUFF_SPECIFIED_ENTERTAIN_PROFIT size 0")
		}
	case constants.BUFF_SPECIFIED_RESTAURANT_CLEAN_NEED:
		bidList := strings.Split(bp2, ",")
		if len(bidList) > 0 {
			for _, bids := range bidList {
				if bid, e := strconv.ParseInt(bids, 10, 64); e != nil || bid <= 0 {
					Err(e, "BUFF_SPECIFIED_RESTAURANT_CLEAN_NEED bid 0")
				} else {
					if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
						Err(e, "BUFF_SPECIFIED_RESTAURANT_CLEAN_NEED 0")
					} else {
						br.BuffSpecifiedRestaurantCleanNeed[bid] += value
					}
				}
			}
		} else {
			Err("BUFF_SPECIFIED_RESTAURANT_CLEAN_NEED size 0")
		}
	case constants.BUFF_SPECIFIED_ENTERTAIN_CLEAN_NEED:
		bidList := strings.Split(bp2, ",")
		if len(bidList) > 0 {
			for _, bids := range bidList {
				if bid, e := strconv.ParseInt(bids, 10, 64); e != nil || bid <= 0 {
					Err(e, "BUFF_SPECIFIED_ENTERTAIN_CLEAN_NEED bid 0")
				} else {
					if value, e := strconv.ParseFloat(bp1, 64); e != nil || value <= 0 {
						Err(e, "BUFF_SPECIFIED_ENTERTAIN_CLEAN_NEED 0")
					} else {
						br.BuffSpecifiedEntertainCleanNeed[bid] += value
					}
				}
			}
		} else {
			Err("BUFF_SPECIFIED_ENTERTAIN_CLEAN_NEED size 0")
		}
	case constants.BUFF_UNLOCK_BUILDING:
		bidList := strings.Split(bp1, ",")
		if len(bidList) > 0 {
			for _, bids := range bidList {
				if bid, e := strconv.ParseInt(bids, 10, 64); e != nil || bid <= 0 {
					Err(e, "BUFF_RESTAURANT_OFFER 0")
				} else {
					br.BuffUnlockBuilding[bid] = true
				}
			}
		} else {
			Err("BUFF_UNLOCK_BUILDING size 0")
		}
	default:
		Err("unknown buff type", bt, bp1, bp2)
	}
}
