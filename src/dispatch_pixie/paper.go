package dispatch_pixie

import (
	"constants"
	dp "db_pixie/paper"
	. "event_dispatch"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service_pixie/paper"
	"strings"
	"time"

	. "zk_manager"
)

func BatchGetPaperHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  BatchGetPaperReq
		output BatchGetPaperResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if len(input.PaperIDList) > 0 {
			//非官方
			if err, list := dp.BatchGetPaperFromRedis(input.PaperIDList); err != nil {
				resp = PIXIE_ERR_BATCH_GET_PAPER
			} else {
				output.Papers = list
			}
		} else if len(input.ClothesIDList) > 0 {
			//官方及非官方
			if list, err := paper.BatchGetPaperFromClothesIDs(input.ClothesIDList); err != nil {
				resp = PIXIE_ERR_BATCH_GET_PAPER
			} else {
				output.Papers = list
			}
		}

		result, parseErr = PixieMarshal(output)
	}
	return
}

func UploadPaperHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  UploadPaperReq
		output UploadPaperResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		st := time.Now().Unix()
		if strings.TrimSpace(input.CName) == "" || strings.TrimSpace(input.Desc) == "" {
			resp = PIXIE_ERR_PARAM_ILLEGAL
			goto END
		}

		if PixiePaperPartTypeMap[input.PartType] == "" || PixiePaperClothesTypeMap[input.ClothesType] <= 0 {
			resp = PIXIE_ERR_PARAM_ILLEGAL
			goto END
		}

		if err, pid := paper.AddPaper(player.Username, player.Nickname, player.GetHead(), player.GetSex(), input.ClothesType, input.PartType, input.CName, input.Desc, st, input.Icon, input.Main, input.Bottom, input.Collar, input.Shadow, input.BeijingFileType, input.Front, input.Back, input.ZValue); err != nil {
			if err == constants.ProcessFileError {
				resp = PIXIE_ERR_UPLOAD_PAPER_FORMAT_WRONG
			} else {
				resp = PIXIE_ERR_UPLOAD_PAPER
			}
		} else {
			NotifyUploadPaper(player, &input, pid)
		}
	END:
		result, parseErr = PixieMarshal(output)
	}
	return
}

func DeletePaperHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input DeletePaperReq
	var output DeletePaperResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if ok, path := LockPixiePaper(input.PaperID); ok {
			deleted := false

			defer func() {
				Unlock(path)

				if deleted {
					NotifyDeletePaper(player, &input)
				}
			}()

			if ok, _ := paper.DeletePaper(input.PaperID, player.Username); ok {
				deleted = true
			} else {
				resp = PIXIE_ERR_DELETE_PAPER_ERR
			}
		} else {
			resp = PIXIE_ERR_LOCKPAPER_FAIL
		}

		result, parseErr = PixieMarshal(output)
	}
	return
}
