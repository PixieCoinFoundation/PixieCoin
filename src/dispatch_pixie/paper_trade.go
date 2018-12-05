package dispatch_pixie

import (
	"constants"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service_pixie/paper"
	"time"
	"tools"
	. "zk_manager"
)

// 发布图纸
func SubmitPaperTradeHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input PublishPaperReq
	var output PublishPaperResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		tNow := time.Now().Unix()
		if input.Sequence < 0 || input.Sequence > constants.PAPER_CIRCULATION_MAX || input.CirculationMax < 0 || input.CirculationMax > constants.PAPER_CIRCULATION_MAX || input.Duration <= 0 || input.MaxPrice < input.MinPrice || input.MinPrice < constants.PIXIE_AUCTION_MIN_LIMIT || (input.MoneyType == constants.PIXIE_GOLD_TYPE && input.MaxPrice > constants.PIXIE_AUCTION_CURRENCY_GOLD_MAX_LIMIT) || (input.MoneyType == constants.PIXIE_PXC_TYPE && input.MaxPrice > constants.PIXIE_AUCTION_CURRENCY_PXC_MAX_LIMIT) {
			resp = PIXIE_ERR_PARAM_ILLEGAL
			goto END
		}

		if success, path := LockPixiePaper(input.PaperID); success {
			defer Unlock(path)

			if tid, err := paper.SubmitPaperTrade(player.Username, player.Nickname, player.GetHead(), player.GetSex(), input.MaxPrice, input.MinPrice, tNow, input.Duration, input.MoneyType, input.PaperID, input.CirculationMax, input.Sequence); err != nil {
				resp = PIXIE_ERR_TRADE_FAIL
			} else {
				output.TradeID = tid
				output.StartTime = tNow
			}
		} else {
			resp = PIXIE_ERR_LOCKPAPER_FAIL
		}
	END:
		result, parseErr = PixieMarshal(output)
	}
	return
}

func CancelPaperTradeHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input CancelPaperTradeReq
	var output CancelPaperTradeResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if input.TradeID <= 0 || input.PaperID <= 0 {
			resp = PIXIE_ERR_PARAM_ILLEGAL
		} else {
			if ok, path := LockPixiePaper(input.PaperID); ok {
				defer Unlock(path)

				now := time.Now().Unix()
				if err := paper.CancelPaperTrade(input.PaperID, player.Username, now, input.TradeID); err != nil {
					if err == constants.TradeNotExist {
						resp = PIXIE_ERR_TRADE_NOT_EXIST
					} else {
						resp = PIXIE_ERR_CANCEL_PAPER_TRADE_ERR
					}

				}
			} else {
				resp = PIXIE_ERR_LOCKPAPER_FAIL
			}
		}

		result, parseErr = PixieMarshal(output)
	}
	return
}

func MarkAuctionFailReadHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input MarkAuctionFailUnReadReq
	var output MarkAuctionFailUnReadResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if input.Type <= 0 {
			resp = PIXIE_ERR_PARAM_ILLEGAL
		} else {
			if err := paper.UpdateAuctionFailUnreadStatus(input.PaperID, false, player.Username, input.Type); err != nil {
				resp = PIXIE_ERR_AUCTION_STATUS_ERR
			}
		}

		result, parseErr = PixieMarshal(output)
	}
	return
}

// 购买图纸
func BuyPaperHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input BuyPaperReq
	var output BuyPaperResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if success, path := LockPixiePaper(input.PaperID); success {
			defer Unlock(path)

			tNow := time.Now().Unix()

			if input.AuctionPriceType == constants.PIXIE_GOLD_TYPE {
				if !player.GoldEnough(int64(input.AuctionPrice)) {
					resp = PIXIE_ERR_PLAYER_CURRENCY_LIMIT
					goto END
				}
			} else {
				if !player.PxcEnough(input.AuctionPrice) {
					resp = PIXIE_ERR_PLAYER_CURRENCY_LIMIT
					goto END
				}
			}

			if input.TradeID <= 0 || input.PaperID <= 0 {
				resp = PIXIE_ERR_PARAM_ILLEGAL
				goto END
			}

			if err := paper.BuyPaper(input.TradeID, input.PaperID, input.Sequence, tNow, player.Username, player.Nickname, player.GetHead(), player.GetSex(), input.AuctionPrice, input.AuctionPriceType); err != nil {
				if err == constants.TradeNotExist {
					resp = PIXIE_ERR_TRADE_NOT_EXIST
				} else {
					resp = PIXIE_ERR_BUY_PAPER
				}
			} else {
				if !player.HaveClothes(tools.GenPixiePlayerClothesID(input.PaperID)) {
					output.CurrentClothes = player.AddPixieClothes(CLOTHES_ORIGIN_PLAYER_PAPER, input.PaperID, 1, tNow, constants.SRC_PAPER_TRADE)
				}

				// player.UseCurrency(input.AuctionPriceType, input.AuctionPrice)

				if input.AuctionPriceType == constants.PIXIE_GOLD_TYPE {
					player.UseGold(int64(input.AuctionPrice))
				} else {
					player.UsePxc(input.AuctionPrice)
				}

				if err := player.SaveStatusToDB(""); err != nil {
					resp = PIXIE_ERR_FLUSH_PLAYER
				} else {
					output.CurrentGold = player.GetMoney()
				}
			}
		} else {
			resp = PIXIE_ERR_LOCKPAPER_FAIL
		}
	END:
		result, parseErr = PixieMarshal(output)
	}
	return
}
