package land

import (
	dl "db_pixie/land"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"tools"
)

func CheckLandOccupy(landID, buildingID int64, username string, nt int64) (ok, own bool, resp PixieRespInfo) {
	if land, err := dl.GetOneLand(landID); err != nil {
		resp = PIXIE_ERR_GET_ONE_LAND
	} else {
		ok = true

		if land.BuildingID != buildingID {
			return
		}

		if land.RenterUsername == username && land.RentStartTime <= nt && land.RentEndTime >= nt+3 {
			//rent own
			own = true
		} else if land.OwnerUsername == username {
			//owner own
			own = true
		}
	}

	return
}

func CheckRentForBuild(land *Land, player *GFPlayer, input *RentLandForBuildReq, tNow int64) (ok bool, resp PixieRespInfo) {
	if input.PriceType != land.RentPriceType {
		resp = PIXIE_ERR_RENT_PRICE_TYPE_WRONG
		return
	}
	if enough := player.GoldEnough(int64(input.Price)); !enough {
		resp = PIXIE_ERR_PLAYER_CURRENCY_LIMIT
		return
	}

	// if enough := player.CurrencyEnough(input.PriceType, input.Price); !enough {
	// 	resp = PIXIE_ERR_PLAYER_CURRENCY_LIMIT
	// 	return
	// }

	if player.Username == land.OwnerUsername {
		resp = PIXIE_ERR_RENT_SELF
		return
	}

	if land.Status != int(LAND_STATUS_NB_RENTING_FOR_BUILD) {
		resp = PIXIE_ERR_LAND_STATUS_WRONG
		return
	}

	if land.RenterUsername != "" && land.RentEndTime > tNow {
		resp = PIXIE_ERR_RENTING
		return
	}

	if legal, _ := tools.AuctionPriceLegal(land.RentMaxPrice, land.RentMinPrice, land.AuctionStartTime, land.AuctionEndTime-land.AuctionStartTime, tNow, float64(input.Price)); !legal {
		resp = PIXIE_ERR_PRICE_ILLEGAL
		return
	}

	ok = true
	return
}

func CheckRentForBusiness(land *Land, player *GFPlayer, input *RentBuildingForBusinessReq, tNow int64) (ok bool, resp PixieRespInfo) {
	if input.PriceType != land.RentPriceType {
		resp = PIXIE_ERR_RENT_PRICE_TYPE_WRONG
		return
	}

	// if enough := player.CurrencyEnough(land.RentPriceType, input.Price); !enough {
	// 	resp = PIXIE_ERR_PLAYER_CURRENCY_LIMIT
	// 	return
	// }

	if enough := player.GoldEnough(int64(input.Price)); !enough {
		resp = PIXIE_ERR_PLAYER_CURRENCY_LIMIT
		return
	}

	if player.Username == land.OwnerUsername {
		resp = PIXIE_ERR_RENT_SELF
		return
	}

	if land.Status != int(LAND_STATUS_WB_RENTING_FOR_BUSINESS) {
		resp = PIXIE_ERR_LAND_STATUS_WRONG
		return
	}

	if land.RenterUsername != "" && land.RentEndTime > tNow {
		resp = PIXIE_ERR_RENTING
		return
	}

	if legal, _ := tools.AuctionPriceLegal(land.RentMaxPrice, land.RentMinPrice, land.AuctionStartTime, land.AuctionEndTime-land.AuctionStartTime, tNow, float64(input.Price)); !legal {
		resp = PIXIE_ERR_PRICE_ILLEGAL
		return
	}

	ok = true
	return
}
