package dispatch_pixie

import (
	"constants"
	"fmt"
	"language"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	suit "service_pixie/suit"
)

func AddSuitHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input AddSuitReq
	var output AddSuitResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if r, ok := suit.CheckSuit(&input.SuitModel, input.ModelType, player); !ok {
			resp = r
			goto END
		}

		if count := player.GetSuitCountByModelType(input.ModelType); count >= constants.PIXIE_SUIT_LIMIT {
			resp = PIXIE_ERR_SUIT_LIMIT_WRONG
		} else {
			if err, md5Str := suit.GetSuitMD5(input.SuitModel); err != nil {
				resp = PIXIE_ERR_SUIT_HANDLE_UNKNOW_WRONG
			} else {
				if exist, suitName := player.ExistSuit(input.ModelType, md5Str); exist {
					PIXIE_ERR_SUIT_ALREADY_EXIST.RespDescCN = fmt.Sprintf(language.L("pixie_suit"), suitName)
					resp = PIXIE_ERR_SUIT_ALREADY_EXIST
				} else {
					player.AddSuit(input.ModelType, input.SuitModel, md5Str)
					if err := player.SaveStatusToDB(constants.PIXIE_SAVE_DB_STATUS_SUIT_MAP); err != nil {
						resp = PIXIE_ERR_FLUSH_PLAYER
					} else {
						output.SuitMD5 = md5Str
					}
				}
			}
		}

	END:
		result, parseErr = PixieMarshal(output)
	}
	return
}

func UpdateSuitHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input UpdateSuitReq
	var output UpdateSuitResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		var err error
		var md5Str string

		if r, ok := suit.CheckSuit(&input.SuitModel, input.ModelType, player); !ok {
			resp = r
		} else {
			if oldExist, _ := player.ExistSuit(input.ModelType, input.SuitMD5); !oldExist {
				resp = PIXIE_ERR_SUIT_NOT_EXIST
			} else {
				if err, md5Str = suit.GetSuitMD5(input.SuitModel); err != nil {
					resp = PIXIE_ERR_SUIT_HANDLE_UNKNOW_WRONG
				} else {
					if exist, suitName := player.ExistSuit(input.ModelType, md5Str); exist {
						PIXIE_ERR_SUIT_ALREADY_EXIST.RespDescCN = fmt.Sprintf(language.L("pixie_suit"), suitName)
						resp = PIXIE_ERR_SUIT_ALREADY_EXIST
					} else {
						player.UpdateSuit(input.ModelType, input.SuitMD5, input.SuitModel, md5Str)
						if err := player.SaveStatusToDB(constants.PIXIE_SAVE_DB_STATUS_SUIT_MAP); err != nil {
							resp = PIXIE_ERR_FLUSH_PLAYER
						} else {
							output.SuitMD5 = md5Str
						}
					}
				}
			}

		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func UpdateSuitNameHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input UpdateSuitNameReq
	var output UpdateSuitNameResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if exist, _ := player.ExistSuit(input.ModelType, input.SuitMD5); !exist {
			resp = PIXIE_ERR_SUIT_NOT_EXIST
		} else {
			player.UpdateSuitName(input.ModelType, input.SuitMD5, input.SuitName)
			if err := player.SaveStatusToDB(constants.PIXIE_SAVE_DB_STATUS_SUIT_MAP); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			}
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func DeleteSuitHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input DeleteSuitReq
	var output DeleteSuitResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if exist, _ := player.ExistSuit(input.ModelType, input.SuitMD5); !exist {
			resp = PIXIE_ERR_SUIT_NOT_EXIST
		} else {
			player.DeleteSuit(input.ModelType, input.SuitMD5)
			if err := player.SaveStatusToDB(constants.PIXIE_SAVE_DB_STATUS_SUIT_MAP); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			}
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func UpdateHomeShowHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input UpdateHomeShowReq
	var output UpdateHomeShowResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		player.SetHomeShow(input.HomeShow)
		if err := player.SaveStatusToDB(constants.PIXIE_SAVE_DB_STATUS_HOME_SHOW); err != nil {
			resp = PIXIE_ERR_FLUSH_PLAYER
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}
