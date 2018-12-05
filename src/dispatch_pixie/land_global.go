package dispatch_pixie

import (
	"constants"
	"dao_player"
	"db_pixie/global"
	"db_pixie/land"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
)

func ListLandHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  ListLandReq
		output ListLandResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		tu := input.TargetUsername
		if tu == "" {
			tu = player.Username
		}

		if ls, err := land.ListPlayerLand(tu); err != nil {
			resp = PIXIE_ERR_LIST_LAND
		} else {
			output.Lands = ls
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}

func GetPlayerMapHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  GetPlayerMapReq
		output GetPlayerMapResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if input.TargetUsername == "" || input.TargetUsername == player.Username {
			resp = PIXIE_ERR_PARAM_ILLEGAL
			goto END
		}

		if status, err := dao_player.GetPlayer(input.TargetUsername); err != nil {
			if err == constants.PlayerNotFoundErr {
				resp = PIXIE_ERR_PLAYER_NOT_FOUND
				goto END
			} else {
				resp = PIXIE_ERR_LIST_LAND
				goto END
			}
		} else {
			if ls, err := land.ListPlayerLand(input.TargetUsername); err != nil {
				resp = PIXIE_ERR_LIST_LAND
			} else {
				output.Lands = ls
				output.UserAttributes = status
			}
		}
	END:
		result, parseErr = PixieMarshal(output)
	}

	return
}

func GetRandomMapHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  GetRandomMapReq
		output GetRandomMapResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if tu, err := global.GetRandomPlayerForVisit(player.Username, player.LastVistUsername); err != nil || tu == "" {
			resp = PIXIE_ERR_GET_RANDOM_MAP
		} else {
			if status, err := dao_player.GetPlayer(tu); err != nil {
				if err == constants.PlayerNotFoundErr {
					resp = PIXIE_ERR_PLAYER_NOT_FOUND
				} else {
					resp = PIXIE_ERR_LIST_LAND
				}
			} else {
				if ls, err := land.ListPlayerLand(tu); err != nil {
					resp = PIXIE_ERR_LIST_LAND
				} else {
					output.Lands = ls
					output.UserAttributes = status
					player.LastVistUsername = tu
				}
			}
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}
