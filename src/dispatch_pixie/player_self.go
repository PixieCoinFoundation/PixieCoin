package dispatch_pixie

import (
	"constants"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"regexp"
	"strings"
)

func ChangeNicknameHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input ChangeNicknameReq
	var output ChangeNicknameResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if matched, err := regexp.MatchString("^"+constants.DEFAULT_NICKNAME_PREFIX+"\\d+$", input.NewNickname); matched || err != nil {
			if matched {
				resp = PIXIE_ERR_NICKNAME_ILLEGAL
			} else {
				resp = PIXIE_ERR_SET_NICKNAME
			}
		} else if player.GetNickname() == input.NewNickname {
			resp = PIXIE_ERR_NICKNAME_NOT_CHANGE
		} else if e := player.SetNickname(input.NewNickname); e != nil {
			if strings.Contains(e.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
				resp = PIXIE_ERR_NICKNAME_ALREADY_EXIST
			} else {
				resp = PIXIE_ERR_SET_NICKNAME
			}
		} else {
			if err := player.SaveStatusToDB(""); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			}
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}

func ChangeStreetNameHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input ChangeStreetNameReq
	var output ChangeStreetNameResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if player.GetStreetName() == input.NewStreetName {
			resp = PIXIE_ERR_STREET_NAME_NOT_CHANGE
		} else {
			player.SetStreetName(input.NewStreetName)

			if err := player.SaveStatusToDB(""); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			}
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}
