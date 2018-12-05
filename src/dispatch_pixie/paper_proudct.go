package dispatch_pixie

import (
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service_pixie/paper"
	"time"
	. "zk_manager"
)

// 生产衣服
func PaperProductHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input PaperProductReq
	var output PaperProductResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		tNow := time.Now().Unix()

		if input.ProductNum <= 0 {
			resp = PIXIE_ERR_PARAM_ILLEGAL
			goto END
		}

		if success, path := LockPixiePaper(input.PaperID); success {
			defer Unlock(path)

			if err := paper.SetPaperProduct(player.Username, input.PaperID, tNow, input.ProductNum); err != nil {
				resp = PIXIE_ERR_PRODUCT_PAPER_FAIL
			} else {
				output.ProductStartTime = tNow
			}
		} else {
			resp = PIXIE_ERR_PAPER_PRODUCT_LOCK_FAIL
		}

	END:
		result, parseErr = PixieMarshal(output)
	}
	return
}

// 取消生产
func CancelPaperProductHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input CancelPaperProductReq
	var output CancelPaperProductResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		tNow := time.Now().Unix()

		if err := paper.CancelPaperProduct(input.PaperID, player.Username, tNow); err != nil {
			resp = PIXIE_ERR_CANCEL_PRODUCT_PAPER_FAIL
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}
