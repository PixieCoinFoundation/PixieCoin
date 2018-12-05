package dispatch_pixie

import (
	"constants"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service_pixie/building"
)

func ChangeStreetSkinHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input ChangeStreetSkinReq
	var output ChangeStreetSkinResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if _, exist := building.GetSkinByID(input.SkinID); !exist {
			resp = PIXIE_ERR_CHANGE_SKINID_NOT_EXIST
		} else {
			if !player.CheckSkinExist(input.SkinID) {
				resp = PIXIE_ERR_SKIN_ID_NOT_EXIST
				goto END
			}
			player.SetMapID(input.SkinID)
			if err := player.SaveStatusToDB(constants.PIXIE_SAVE_DB_STATUS_MAP_SKIN); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			}
		}
	END:
		result, parseErr = PixieMarshal(output)
	}
	return
}

func BuySkinHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input BuySkinReq
	var output BuySkinResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if skin, exist := building.GetSkinByID(input.SkinID); !exist {
			resp = PIXIE_ERR_CHANGE_SKINID_NOT_EXIST
		} else {
			player.UseGold(skin.Price)
			player.AddSkinMap(skin.ID)
			if err := player.SaveStatusToDB(""); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			} else {
				output.CurrentGold = player.GetMoney()
				output.CurrentPxc = player.GetPxc()
			}
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}
