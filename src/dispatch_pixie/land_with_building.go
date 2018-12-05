package dispatch_pixie

import (
	"constants"
	"db_pixie/land"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service_pixie/building"
	sl "service_pixie/land"
	"time"
	. "zk_manager"
)

func SetShopModelHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input SetShopModelReq
	var output SetShopModelResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		nt := time.Now().Unix()
		if ok, own, r := sl.CheckLandOccupy(input.LandID, input.BuildingID, player.Username, nt); !ok {
			resp = r
		} else if own {
			if e := land.SetShopModel(input.LandID, input.BuildingID, player.Username, input.ShopModel); e != nil {
				resp = PIXIE_ERR_SET_SHOP_MODEL
			}
		} else {
			resp = PIXIE_ERR_LAND_OWNER_ILLEGAL
		}

		result, parseErr = PixieMarshal(output)
	}
	return
}

func LandBuildingLevelUpHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input LandBuildingLevelUpReq
	var output LandBuildingLevelUpResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		nt := time.Now().Unix()

		if ok, r, bp, _, bd := sl.CheckPlayerBuildingLevelUp(player.Username, input.BuildingID, input.LandID, input.FromLevel, input.FromLevel+1); !ok {
			resp = r
		} else {
			if enough := player.GoldEnough(bp); !enough {
				resp = PIXIE_ERR_PLAYER_CURRENCY_LIMIT
				goto END
			}

			if et, err := land.LandBuildingLevelUp(player.Username, input.LandID, input.BuildingID, nt, bd, input.FromLevel); err != nil {
				resp = PIXIE_ERR_LANG_BUILDING_LEVEL_UP
			} else {
				output.LevelUpStartTime = nt
				output.LevelUpEndTime = et

				player.UseGold(bp)

				if err := player.SaveStatusToDB(""); err != nil {
					resp = PIXIE_ERR_FLUSH_PLAYER
				}
			}
		}

	END:
		output.CurrentGold = player.GetMoney()
		result, parseErr = PixieMarshal(output)
	}
	return
}

func CancelLandBuildingLevelUpHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input LandBuildingLevelUpReq
	var output LandBuildingLevelUpResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if l, err := land.GetOneLand(input.LandID); err != nil {
			resp = PIXIE_ERR_GET_ONE_LAND
		} else {
			nt := time.Now().Unix()

			if l.BuildingLevel == input.FromLevel && l.BuildingID == input.BuildingID && l.LevelUpStartTime > 0 && l.LevelUpStartTime <= nt && l.LevelUpEndTime > 0 && l.LevelUpEndTime >= nt+3 {
				toLevel := input.FromLevel + 1
				var toLevelUpMoney int64
				// var toLevelUpMoneyType int

				if b, exist := building.GetEntertainBuildingByIDLevel(input.BuildingID, toLevel); exist {
					toLevelUpMoney = b.BuildPrice
				} else if b, exist := building.GetRestaurantBuildingByIDLevel(input.BuildingID, toLevel); exist {
					toLevelUpMoney = b.BuildPrice
				} else {
					resp = PIXIE_ERR_LAND_BUILDING_ILLEGAL
					goto END
				}

				if err := land.CancelLandBuildingLevelUp(player.Username, input.LandID, input.BuildingID, input.FromLevel); err != nil {
					resp = PIXIE_ERR_CANCEL_LEVEL_UP_BUILDING
				} else {
					player.AddGold(toLevelUpMoney)

					if err := player.SaveStatusToDB(""); err != nil {
						resp = PIXIE_ERR_FLUSH_PLAYER
					}
				}
			} else {
				resp = PIXIE_ERR_LAND_STATUS_WRONG
			}
		}
	END:
		output.CurrentGold = player.GetMoney()
		result, parseErr = PixieMarshal(output)
	}
	return
}

func LandRentForBusinessHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  LandRentForBusinessReq
		output LandRentForBusinessResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		tNow := time.Now().Unix()
		duration := int64(input.Duration * 3600)
		auctionDuration := int64(input.AuctionDuration * 3600)

		if input.AuctionDuration <= 0 || input.Duration <= 0 || input.MaxPrice < input.MinPrice || input.MinPrice < constants.PIXIE_RENT_MIN_LIMIT || (input.MoneyType == constants.PIXIE_GOLD_TYPE && input.MaxPrice > constants.PIXIE_RENT_CURRENCY_GOLD_MAX_LIMIT) || (input.MoneyType == constants.PIXIE_PXC_TYPE && input.MaxPrice > constants.PIXIE_RENT_CURRENCY_PXC_MAX_LIMIT) || input.MoneyType <= 0 || input.LandID <= 0 {
			resp = PIXIE_ERR_PARAM_ILLEGAL
			goto END
		}

		if err := land.LandRentForBusiness(tNow, auctionDuration, duration, input.MaxPrice, input.MinPrice, input.MoneyType, player.Username, input.LandID); err != nil {
			resp = PIXIE_ERR_RENT_LAND_FOR_BUSINESS
		} else {
			output.AuctionStartTime = tNow
			output.AuctionEndTime = tNow + auctionDuration
		}

	END:
		result, parseErr = PixieMarshal(output)
	}
	return
}

func StopLandRentForBusinessHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  StopLandRentForBusinessReq
		output StopLandRentForBusinessResp
	)
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if ok, path := LockPixieLand(input.LandID); ok {
			defer Unlock(path)
			if err := land.LandStopRentForBusiness(player.Username, input.LandID); err != nil {
				resp = PIXIE_ERR_RENT_LAND_FOR_BUSINESS
			}
		} else {
			resp = PIXIE_ERR_LOCKLAND_FAIL
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func ListLandSaleHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  ListLandSaleReq
		output ListLandSaleResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		nt := time.Now().Unix()
		if l, lmd, err := land.GetLandSale(input.LandID, player.Username, nt); err != nil {
			if err == constants.LandStatusWrong {
				resp = PIXIE_ERR_LAND_NOT_IN_BUSINESS
			} else {
				resp = PIXIE_ERR_LIST_LAND_SALE
			}
		} else {
			output.Result = l
			output.ShopModel = lmd
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}

func DemandStartBusinessHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  DemandStartBusinessReq
		output DemandStartBusinessResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		player.RefreshStreet(false)

		if input.BuildingID <= 0 || input.LandID <= 0 {
			resp = PIXIE_ERR_PARAM_ILLEGAL
		} else if !building.BuildingIsDemand(input.BuildingID) {
			resp = PIXIE_ERR_BUILDING_NOT_DEMAND
		} else if ok, _ := land.LandDemandUpdateStatus(LAND_STATUS_WB_NORMAL, LAND_STATUS_WB_IN_BUSINESS, player.Username, input.LandID, input.BuildingID); !ok {
			resp = PIXIE_ERR_START_DEMAND_BUSINESS
		} else {
			player.RefreshStreet(true)
			output.StreetDetail = player.GetStreetDetail()

			if err := player.SaveStatusToDB(""); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			}
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}

func DemandStopBusinessHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  DemandStopBusinessReq
		output DemandStopBusinessResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		player.RefreshStreet(false)

		if input.BuildingID <= 0 || input.LandID <= 0 {
			resp = PIXIE_ERR_PARAM_ILLEGAL
		} else if !building.BuildingIsDemand(input.BuildingID) {
			resp = PIXIE_ERR_BUILDING_NOT_DEMAND
		} else if ok, _ := land.LandDemandUpdateStatus(LAND_STATUS_WB_IN_BUSINESS, LAND_STATUS_WB_NORMAL, player.Username, input.LandID, input.BuildingID); !ok {
			resp = PIXIE_ERR_STOP_DEMAND_BUSINESS
		} else {
			player.RefreshStreet(true)
			output.StreetDetail = player.GetStreetDetail()

			if err := player.SaveStatusToDB(""); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			}
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}

func LandStartBusinessHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  LandStartBusinessReq
		output LandStartBusinessResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		player.RefreshStreet(false)

		if (input.Type != 1 && input.Type != 2) || len(input.PaperIDList) <= 0 || input.LandID <= 0 {
			resp = PIXIE_ERR_PARAM_ILLEGAL
		} else if err := land.StartBusiness(player.Username, input.ShopName, input.LandID, input.PaperIDList, input.Type); err != nil {
			resp = PIXIE_ERR_START_BUSINESS
		} else {
			player.RefreshStreet(true)
			output.StreetDetail = player.GetStreetDetail()

			if err := player.SaveStatusToDB(""); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			}
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}

func LandStopBusinessHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  LandStopBusinessReq
		output LandStopBusinessResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		player.RefreshStreet(false)

		if (input.Type != 1 && input.Type != 2) || input.LandID <= 0 {
			resp = PIXIE_ERR_PARAM_ILLEGAL
		} else {
			if ok, path := LockPixieLand(input.LandID); ok {
				defer Unlock(path)

				if err := land.StopBusiness(player.Username, input.LandID, input.Type); err != nil {
					resp = PIXIE_ERR_STOP_BUSINESS
				} else {
					player.RefreshStreet(true)
					output.StreetDetail = player.GetStreetDetail()

					if err := player.SaveStatusToDB(""); err != nil {
						resp = PIXIE_ERR_FLUSH_PLAYER
					}
				}
			} else {
				resp = PIXIE_ERR_LOCK_LAND
			}
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}

func RemoveBuildingHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  RemoveBuildingReq
		output RemoveBuildingResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if success, err := land.RemoveBuilding(input.LandID, input.BuildingID, player.Username); !success {
			if err == constants.BuildUserNotMatch {
				resp = PIXIE_ERR_BUILD_USER_NOT_MATCH
			} else if err == constants.NotOccupyLand {
				resp = PIXIE_ERR_NOT_OCCUPY_LAND
			} else {
				resp = PIXIE_ERR_REMOVE_BUILDING
			}
		}

		result, parseErr = PixieMarshal(output)
	}
	return
}

func RentBuildingForBusinessHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  RentBuildingForBusinessReq
		output RentBuildingForBusinessResp
	)
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if ok, path := LockPixieLand(input.LandID); ok {
			defer Unlock(path)

			tNow := time.Now().Unix()
			landModle, err := land.GetOneLand(input.LandID)
			if err != nil {
				resp = PIXIE_ERR_LAND_GET_ERR
				goto END
			}

			if ok, r := sl.CheckRentForBusiness(&landModle, player, &input, tNow); !ok {
				resp = r
				goto END
			}

			if err := land.RentBuildingForBusiness(player.Username, player.Nickname, player.GetHead(), player.GetSex(), tNow, tNow+landModle.RentDuration, input.BuildingID, input.LandID); err != nil {
				resp = PIXIE_ERR_RENT_LAND_ERR
			} else {
				if input.PriceType == constants.PIXIE_GOLD_TYPE {
					player.UseGold(int64(input.Price))
				} else {
					player.UsePxc(input.Price)
				}
				// player.UseCurrency(input.PriceType, input.Price)
				if err := player.SaveStatusToDB(""); err != nil {
					resp = PIXIE_ERR_FLUSH_PLAYER
				} else {
					output.CurrentGold = player.GetMoney()
					output.RentStartTime = tNow
					output.RentEndTime = tNow + landModle.RentDuration
				}
			}
		} else {
			resp = PIXIE_ERR_LOCKLAND_FAIL
		}
	END:
		result, parseErr = PixieMarshal(output)

	}
	return
}
