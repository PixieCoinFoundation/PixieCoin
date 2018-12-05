package dispatch_pixie

import (
	"constants"
	"db_pixie/paper"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service_pixie/clothes"
	spaper "service_pixie/paper"

	. "event_dispatch"
	"time"
	"tools"
	. "zk_manager"
)

func ListOfficialClothesHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var output ListOfficialClothesResp

	output.Res = spaper.GetAllOfficialPaper()

	result, parseErr = PixieMarshal(output)
	return
}

func BuyOfficialClothesHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input BuyOfficialClothesReq
	var output BuyOfficialClothesResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		totalCurrency := input.Price * float64(input.BuyCount)

		if op, exist := spaper.GetOneOfficialPaperByID(input.PaperID); exist {
			if op.PriceType == input.PriceType {
				if op.PriceType == constants.PIXIE_GOLD_TYPE {
					if !player.GoldEnough(int64(totalCurrency)) {
						resp = PIXIE_ERR_PLAYER_CURRENCY_LIMIT
						goto END
					} else {
						player.UseGold(int64(totalCurrency))
					}
				} else if op.PriceType == constants.PIXIE_PXC_TYPE {
					if !player.PxcEnough(totalCurrency) {
						resp = PIXIE_ERR_PLAYER_CURRENCY_LIMIT
						goto END
					} else {
						player.UsePxc(totalCurrency)
					}
				} else {
					resp = PIXIE_ERR_OFFICIAL_CLOTHES_PRICE_TYPE_WRONG
					goto END
				}

				cc := player.AddPixieClothes(CLOTHES_ORIGIN_OFFICIAL_PAPER, input.PaperID, input.BuyCount, time.Now().Unix(), constants.SRC_OFFICIAL_SHOP)

				if err := player.SaveStatusToDB(""); err != nil {
					resp = PIXIE_ERR_FLUSH_PLAYER
				} else {
					output.CurrentClothes = cc
				}
			} else {
				resp = PIXIE_ERR_OFFICIAL_CLOTHES_PRICE_TYPE_WRONG
			}
		} else {
			resp = PIXIE_ERR_OFFICIAL_CLOTHES_NOT_EXIST
		}
	END:
		output.CurrentGold = player.GetMoney()
		output.CurrentPxc = player.GetPxc()
		result, parseErr = PixieMarshal(output)
	}

	return
}

func BuyClothesHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  BuyClothesReq
		output BuyClothesResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		tNow := time.Now()
		nt := tNow.Unix()

		occupyLogID, _ := paper.GetPaperOccupyLogID(input.PaperID, input.SaleUsername)

		if ok, r := clothes.CheckSellClothes(input.LandID, input.PaperID, nt, input.SaleUsername, input.Price, input.PriceType, input.BuyCount, player); !ok {
			resp = r
			goto END
		}

		if buyOK, ot, r := buyOnePaperClothes(input.SaleUsername, input.LandID, occupyLogID, input.PaperID, input.Price, nt, input.PriceType, input.BuyCount, player); buyOK {
			output = ot

			if err := player.SaveStatusToDB(""); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			}
		} else {
			resp = r
			goto END
		}
	END:
		result, parseErr = PixieMarshal(output)
	}
	return
}

func BatchBuyClothesHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  BatchBuyClothesReq
		output BatchBuyClothesResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		tNow := time.Now()
		nt := tNow.Unix()

		if ok, canBuyList, r := clothes.CheckBatchSellClothes(input.LandID, input.PaperList, nt, input.SaleUsername, player); !ok {
			resp = r
			goto END
		} else {
			output.SuccessClothesList = make([]ClothesInfo, 0)
			for _, cp := range canBuyList {
				if buyOK, ot, _ := buyOnePaperClothes(input.SaleUsername, input.LandID, 0, cp.PaperID, cp.Price, nt, cp.PriceType, cp.BuyCount, player); buyOK {
					output.SuccessClothesList = append(output.SuccessClothesList, ot.CurrentClothes)
				}
			}
		}

		if err := player.SaveStatusToDB(""); err != nil {
			resp = PIXIE_ERR_FLUSH_PLAYER
		}

	END:
		output.CurrentGold = player.GetMoney()
		output.CurrentPxc = player.GetPxc()
		result, parseErr = PixieMarshal(output)
	}
	return
}

func BuyClothesListHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var (
		input  BuyClothesListReq
		output BuyClothesListResp
	)

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		tNow := time.Now()
		nt := tNow.Unix()
		output.SuccessClothesList = make([]ClothesInfo, 0)
		nonOfficialClothesMap := make(map[int64]*SimpleClothes)

		for _, c := range input.ClothesList {
			if src, id := tools.GetPixieClothesDetail(c.ClothesID); src == string(CLOTHES_ORIGIN_OFFICIAL_PAPER) {
				//buy official paper
				if op, exist := spaper.GetOneOfficialPaperByID(id); exist {
					if op.PriceType == constants.PIXIE_GOLD_TYPE && player.GoldEnough(op.Price) {
						player.UseGold(op.Price)
					} else if op.PriceType == constants.PIXIE_PXC_TYPE && player.PxcEnough(float64(op.Price)) {
						player.UsePxc(float64(op.Price))
					} else {
						continue
					}

					cc := player.AddPixieClothes(CLOTHES_ORIGIN_OFFICIAL_PAPER, id, c.BuyCount, nt, constants.SRC_SHOP)
					output.SuccessClothesList = append(output.SuccessClothesList, cc)
				}
			} else {
				nonOfficialClothesMap[id] = c
			}
		}

		//buy non official clothes
		if err, nonOfficialPaperList := paper.ListPaperOwnByIDMap(nonOfficialClothesMap); err == nil {
			for _, op := range nonOfficialPaperList {
				if buyOK, ot, _ := buyOnePaperClothes(op.OwnerUsername, 0, op.OccupyLogID, op.PaperID, op.ClothPrice, nt, op.ClothPriceType, nonOfficialClothesMap[op.PaperID].BuyCount, player); buyOK {
					output.SuccessClothesList = append(output.SuccessClothesList, ot.CurrentClothes)
				}
			}
		}

		if len(output.SuccessClothesList) > 0 {
			if err := player.SaveStatusToDB(""); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			}
		}

		output.CurrentGold = player.GetMoney()
		output.CurrentPxc = player.GetPxc()
		result, parseErr = PixieMarshal(output)
	}
	return
}

func buyOnePaperClothes(saleUsername string, landID, occupyLogID, paperID int64, price float64, nt int64, priceType, buyCount int, player *GFPlayer) (buyOK bool, output BuyClothesResp, resp PixieRespInfo) {
	if ok, path := LockPixieLand(landID); ok {
		defer Unlock(path)
		if ok, path := LockPixiePaper(paperID); ok {
			defer Unlock(path)
			totalCurrency := price * float64(buyCount)
			if priceType == constants.PIXIE_GOLD_TYPE {
				if !player.GoldEnough(int64(totalCurrency)) {
					resp = PIXIE_ERR_PLAYER_CURRENCY_LIMIT
					return
				}
			} else {
				if !player.PxcEnough(totalCurrency) {
					resp = PIXIE_ERR_PLAYER_CURRENCY_LIMIT
					return
				}
			}

			if success, _ := paper.SellPaperClothes(paperID, occupyLogID, buyCount, saleUsername, price, priceType); success {
				if priceType == constants.PIXIE_GOLD_TYPE {
					player.UseGold(int64(totalCurrency))
				} else {
					player.UsePxc(totalCurrency)
				}
				cc := player.AddPixieClothes(CLOTHES_ORIGIN_PLAYER_PAPER, paperID, buyCount, nt, constants.SRC_SHOP)

				output.CurrentGold = player.GetMoney()
				output.CurrentPxc = player.GetPxc()
				output.CurrentClothes = cc

				//event dispatch
				NotifyNonofficialPaperClothesSale(saleUsername, landID, paperID, buyCount, priceType, price, player, nt)
				buyOK = true
			} else {
				resp = PIXIE_ERR_BUY_CLOTHES
			}
		} else {
			resp = PIXIE_ERR_LOCKPAPER_FAIL
		}
	} else {
		resp = PIXIE_ERR_LOCK_LAND
	}

	return
}

func PaperClothesPricingHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input PaperClothesPricingReq
	var output PaperClothesPricingResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {

		if input.Price < constants.PIXIE_PRICING_MIN_LIMIT || (input.PriceType == constants.PIXIE_GOLD_TYPE && input.Price > constants.PIXIE_PRICING_CURRENCY_GOLD_MAX_LIMIT) || (input.PriceType == constants.PIXIE_PXC_TYPE && input.Price > constants.PIXIE_PRICING_CURRENCY_PXC_MAX_LIMIT) {
			resp = PIXIE_ERR_PARAM_ILLEGAL
			goto END
		}

		if ok, path := LockPixiePaper(input.PaperID); ok {
			defer Unlock(path)

			if input.Type == 1 {
				if err := spaper.ClothesPricing(input.Price, input.PriceType, input.PaperID, player.Username); err != nil {
					resp = PIXIE_ERR_PAPER_CLOTHES_PRICING_SET
				}
			} else if input.Type == 2 {
				if err := spaper.ClothesPricingChangePrice(input.Price, input.PaperID, player.Username, input.PriceType); err != nil {
					resp = PIXIE_ERR_PAPER_CLOTHES_PRICING_CHANGE
				}
			} else {
				resp = PIXIE_ERR_PAPER_CLOTHES_PRICING_TYPE_WRONG
			}
		} else {
			resp = PIXIE_ERR_LOCKPAPER_FAIL
		}
	END:
		result, parseErr = PixieMarshal(output)
	}
	return
}
