package dispatch_pixie

import (
	"constants"
	"db_pixie/land"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	sl "service_pixie/land"
	"time"
	. "zk_manager"
)

func LandRentForBuildHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  LandRentForBuildReq
		output LandRentForBuildResp
	)
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		tNow := time.Now().Unix()
		auctionDuration := int64(input.AuctionDuration * 3600)

		if input.AuctionDuration <= 0 || input.Duration <= 0 || input.MaxPrice < input.MinPrice || input.MinPrice < constants.PIXIE_RENT_MIN_LIMIT || (input.MoneyType == constants.PIXIE_GOLD_TYPE && input.MaxPrice > constants.PIXIE_RENT_CURRENCY_GOLD_MAX_LIMIT) || (input.MoneyType == constants.PIXIE_PXC_TYPE && input.MaxPrice > constants.PIXIE_RENT_CURRENCY_PXC_MAX_LIMIT) || input.MoneyType <= 0 || input.LandID <= 0 {
			resp = PIXIE_ERR_PARAM_ILLEGAL
			goto END
		}

		if err := land.LandStartRentingForBuild(tNow, auctionDuration, int64(input.Duration*3600), input.MaxPrice, input.MinPrice, input.MoneyType, player.Username, input.LandID); err != nil {
			resp = PIXIE_ERR_LAND_RENT_FOR_BUILD
		} else {
			output.AuctionStartTime = tNow
			output.AuctionEndTime = tNow + auctionDuration
		}

	END:
		result, parseErr = PixieMarshal(output)
	}

	return
}

func StopLandRentForBuildHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  StopLandRentForBuildReq
		output StopLandRentForBuildResp
	)
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if ok, path := LockPixieLand(input.LandID); ok {
			defer Unlock(path)
			if err := land.LandStopRentingForBuild(player.Username, input.LandID); err != nil {
				resp = PIXIE_ERR_LAND_RENT_FOR_BUILD
			}
		} else {
			resp = PIXIE_ERR_LOCKLAND_FAIL
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func LandBuildingHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  LandBuildingReq
		output LandBuildingResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		nt := time.Now().Unix()

		if ok, r, bp, _, bd := sl.CheckPlayerBuilding(player.Username, input.BuildingID, input.LandID); !ok {
			resp = r
		} else {
			if enough := player.GoldEnough(bp); !enough {
				resp = PIXIE_ERR_PLAYER_CURRENCY_LIMIT
				goto END
			}

			if input.BuildType != 1 && input.BuildType != 2 {
				resp = PIXIE_ERR_PARAM_ILLEGAL
			} else if et, err := land.LandStartBuilding(player.Username, input.LandID, input.BuildingID, nt, bd, input.BuildType); err != nil {
				resp = PIXIE_ERR_LANG_BUILDING
			} else {
				output.BuildStartTime = nt
				output.BuildEndTime = et
				player.UseGold(bp)

				if err := player.SaveStatusToDB(""); err != nil {
					resp = PIXIE_ERR_FLUSH_PLAYER
				}
			}
		}

	END:
		output.CurrentGold = player.GetMoney()
		output.CurrentPxc = player.GetPxc()
		result, parseErr = PixieMarshal(output)
	}

	return
}

func StopLandBuildingHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  StopLandBuildingReq
		output StopLandBuildingResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if input.BuildType != 1 && input.BuildType != 2 {
			resp = PIXIE_ERR_PARAM_ILLEGAL
		} else if err := land.StopLandBuilding(player.Username, input.LandID, input.BuildingID, input.BuildType); err != nil {
			resp = PIXIE_ERR_LANG_BUILDING
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}

func RentLandForBuildHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  RentLandForBuildReq
		output RentLandForBuildResp
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

			if ok, r := sl.CheckRentForBuild(&landModle, player, &input, tNow); !ok {
				resp = r
				goto END
			}

			if err := land.RentLandForBuild(player.Username, player.Nickname, player.GetHead(), player.GetSex(), tNow, tNow+landModle.RentDuration, input.LandID); err != nil {
				resp = PIXIE_ERR_RENT_LAND_FOR_BUILD
			} else {
				player.UseGold(input.Price)

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
