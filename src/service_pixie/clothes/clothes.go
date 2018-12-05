package clothes

import (
	dl "db_pixie/land"
	"encoding/json"
	. "logger"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
)

func CheckSellClothes(landID, paperID, nt int64, saleUsername string, price float64, priceType, buyCount int, player *GFPlayer) (ok bool, resp PixieRespInfo) {
	if saleUsername == player.Username {
		resp = PIXIE_ERR_BUY_SELF
		return
	}

	if land, err := dl.GetOneLand(landID); err != nil {
		resp = PIXIE_ERR_GET_ONE_LAND
		return
	} else {
		if land.SaleInfo != "" {
			var salePaperIDList []int64
			if err := json.Unmarshal([]byte(land.SaleInfo), &salePaperIDList); err != nil {
				Err(err)
				resp = PIXIE_ERR_LAND_SALE_FORMAT
				return
			}

			if land.Status != int(LAND_STATUS_WB_IN_BUSINESS) {
				resp = PIXIE_ERR_LAND_NOT_IN_BUSINESS
				return
			}

			if saleUsername == land.RenterUsername && land.RentEndTime <= nt {
				resp = PIXIE_ERR_LAND_RENT_END
				return
			}

			if (saleUsername != land.OwnerUsername && saleUsername != land.RenterUsername) || landID <= 0 || paperID <= 0 || buyCount <= 0 || price <= 0 || priceType <= 0 {
				resp = PIXIE_ERR_PARAM_ILLEGAL
				return
			}
		} else {
			resp = PIXIE_ERR_CLOTHES_NOT_SALE
			return
		}
	}

	ok = true
	return
}

func CheckBatchSellClothes(landID int64, simplePapers []*SimplePaper, nt int64, saleUsername string, player *GFPlayer) (ok bool, canBuyPapers []*SimplePaper, resp PixieRespInfo) {
	if saleUsername == player.Username {
		resp = PIXIE_ERR_BUY_SELF
		return
	}

	if land, err := dl.GetOneLand(landID); err != nil {
		resp = PIXIE_ERR_GET_ONE_LAND
		return
	} else {
		if land.SaleInfo != "" {
			var salePaperIDList []int64
			if err := json.Unmarshal([]byte(land.SaleInfo), &salePaperIDList); err != nil {
				Err(err)
				resp = PIXIE_ERR_LAND_SALE_FORMAT
				return
			}

			if land.Status != int(LAND_STATUS_WB_IN_BUSINESS) {
				resp = PIXIE_ERR_LAND_NOT_IN_BUSINESS
				return
			}

			if saleUsername == land.RenterUsername && land.RentEndTime <= nt {
				resp = PIXIE_ERR_LAND_RENT_END
				return
			}

			if (saleUsername != land.OwnerUsername && saleUsername != land.RenterUsername) || landID <= 0 {
				resp = PIXIE_ERR_PARAM_ILLEGAL
				return
			}

			canBuyPapers = make([]*SimplePaper, 0)
			for _, sp := range simplePapers {
				if sp.PriceType > 0 && sp.Price > 0 && sp.BuyCount > 0 {
					for _, pid := range salePaperIDList {
						if pid == sp.PaperID {
							canBuyPapers = append(canBuyPapers, sp)
							break
						}
					}
				}
			}
		} else {
			resp = PIXIE_ERR_CLOTHES_NOT_SALE
			return
		}
	}

	ok = true
	return
}
