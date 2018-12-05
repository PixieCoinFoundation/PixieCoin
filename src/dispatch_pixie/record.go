package dispatch_pixie

import (
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service_pixie/record"
)

func SaveRecordHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input SaveRecordReq
	var output SaveRecordResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if fnpc, ct := record.GetRecordDetail(input.RecordID); ct > 0 {
			//check tili
			if player.GetTili() >= ct {
				player.RefreshStreet(false)
				player.SaveRecord(input.RecordID, input.Rank, input.Clothes, input.Score, fnpc, input.DeviceId, ct)
				player.RefreshStreet(false)

				if err := player.SaveStatusToDB(""); err != nil {
					resp = PIXIE_ERR_FLUSH_PLAYER
				} else {
					output.Res = player.GetRecordBylevelNo(input.RecordID)
					output.NpcFriendship = player.GetNpcFriendship()
					output.NpcBuff = player.GetNpcBuff()
					output.StreetDetail = player.GetStreetDetail()

					output.Tili = player.GetTili()
					output.MaxTili = player.GetMaxTili()
					output.Level, output.Exp = player.GetLevelAndExp()
				}
			} else {
				resp = PIXIE_ERR_TILI_NOT_ENOUGH
			}
		} else {
			resp = PIXIE_ERR_RECORD_NOT_EXIST
		}

		result, parseErr = PixieMarshal(output)
	}
	return
}

func ViewRecordHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input ViewRecordReq
	var output ViewRecordResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if !record.LevelIDLegal(input.RecordID) {
			resp = PIXIE_ERR_RECORD_NOT_EXIST
			goto END
		}

		player.ViewRecord(input.RecordID)
		if err := player.SaveStatusToDB(""); err != nil {
			resp = PIXIE_ERR_FLUSH_PLAYER
		}
	END:
		result, parseErr = PixieMarshal(output)
	}
	return
}
