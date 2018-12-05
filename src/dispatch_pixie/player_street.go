package dispatch_pixie

import (
	// "constants"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service_pixie/building"
)

func DrawStreetWalletHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input DrawStreetWalletReq
	var output DrawStreetWalletResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		player.RefreshStreet(false)

		output.RewardGold = player.ClearStreetWallet()
		player.AddGold(output.RewardGold)

		output.CurrentGold = player.GetMoney()
		output.StreetDetail = player.GetStreetDetail()

		if err := player.SaveStatusToDB(""); err != nil {
			resp = PIXIE_ERR_FLUSH_PLAYER
		}

		result, parseErr = PixieMarshal(output)
	}
	return
}

func SetStreetCleanHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input SetStreetCleanReq
	var output SetStreetCleanResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if p, s, ep := building.GetCleanObjectDetail(input.CleanObjectID); p <= 0 || s <= 0 {
			resp = PIXIE_ERR_CLEAN_OBJECT_NOT_EXIST
		} else {
			employCostGold := int64(ep)
			player.RefreshStreet(false)

			if input.Type == 0 {
				if player.GoldEnough(employCostGold) {
					if !player.SetStreetClean(input.CleanObjectID, input.Count) {
						resp = PIXIE_ERR_EMPLOY_CLEAN_OBJECT_FAIL
						goto END
					} else {
						player.UseGold(employCostGold)
					}
				} else {
					resp = PIXIE_ERR_PLAYER_CURRENCY_LIMIT
				}
			} else if input.Type == 1 {
				if !player.SetStreetClean(input.CleanObjectID, -input.Count) {
					resp = PIXIE_ERR_UNEMPLOY_CLEAN_OBJECT_FAIL
					goto END
				}
			} else {
				resp = PIXIE_ERR_PARAM_ILLEGAL
				goto END
			}

			player.RefreshStreet(true)
			if err := player.SaveStatusToDB(""); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			}
		}
	END:
		output.CurrentGold = player.GetMoney()
		output.CurrentPxc = player.GetPxc()
		output.StreetDetail = player.GetStreetDetail()
		result, parseErr = PixieMarshal(output)
	}
	return
}
