package dispatch_pixie

import (
	"common"
	// "constants"
	"fmt"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service/files"
	"time"
	// . "types"
)

func FileUploadHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input FileUploadReq
	var output FileUploadResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		canUpload := false
		success := true
		et := 0
		// n := time.Now().Format("20060102")
		//检查上传次数
		if input.Type == 1 {
			// if player.GetUploadCustomDayCnt(n) < constants.UPLOAD_CUSTOM_DAY_LIMIT {
			canUpload = true
			// } else {
			// resp = PIXIE_ERR_UPLOAD_LIMIT
			// goto END
			// }

			//举报上传
		} else if input.Type == 2 {
			canUpload = true

		} else if input.Type == 3 {
			canUpload = true

		} else {
			resp = PIXIE_ERR_UPLOAD_TYPE_WRONG
			goto END
		}

		if canUpload {
			output.Filenames = make([]string, 0)
			for i := 0; i < len(input.Content); i++ {
				md5, _ := common.GenMD5Raw(input.Content[i])
				key := fmt.Sprintf("%s/%s_%d_%d_%d_%s", files.GetPlatformOSSDirName(), player.Username, time.Now().UnixNano(), i, input.Type, md5)
				if _, err := files.UploadFileByBuffer(input.Content[i], key, et); err != nil {
					success = false
				} else {
					output.Filenames = append(output.Filenames, key)
				}
			}

			if success {
				if input.Type == 1 {
					// player.AddUploadCustomDayCnt(n)
					// if err := player.SaveStatusToDB(STATUS_DETAIL_DAY_EXTRA); err != nil {
					// 	resp = PIXIE_ERR_SAVE_DB_WORNG
					// }
				}
			} else {
				resp = PIXIE_ERR_UPLOAD_FILE_FAIL
			}
		}

	END:
		result, parseErr = PixieMarshal(output)
	}
	return
}
