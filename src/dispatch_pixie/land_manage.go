package dispatch_pixie

import (
	dland "db_pixie/land"
	"db_pixie/paper"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service_pixie/land"
	spaper "service_pixie/paper"
	"strconv"
	"strings"
	"time"
)

func LandManageHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input LandManageReq
	var output LandManageResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		tNow := time.Now()
		yesterday := tNow.AddDate(0, 0, -1).Format("20060102")
		if err, inSaleCount, inSaleInventory, inSaleShop, ownShop, buildingShop, waitBusiness, waitBuild, todayPxcSale, yesterdayPxcSale, allPxcSale, todayGoldSale, yesterdayGoldSale, allGoldSale, todaySaleClothes, yesterdaySaleClothes, allSaleClothes, list := land.GetSaleInfo(player.Username, yesterday, input.Page, input.PageSize, tNow); err != nil {
			resp = PIXIE_ERR_GET_MY_LAND_WRONG
		} else {
			output.InSaleClothesNum = inSaleCount
			output.InSaleClothesInventory = inSaleInventory

			output.InSaleShopCount = inSaleShop
			output.OwnShopCount = ownShop
			output.BuildingLandCount = buildingShop

			//错误
			output.WBRentedBusinessCount = waitBusiness
			output.NBRentedBuildCount = waitBuild
			//正确
			output.WaitBusinessCount = waitBusiness
			output.WaitBuildCount = waitBuild

			output.TodayGoldSale = todayGoldSale
			output.TodayPxcSale = todayPxcSale
			output.TodaySaleClothes = todaySaleClothes

			output.YesterdayGoldSale = yesterdayGoldSale
			output.YesterdayPxcSale = yesterdayPxcSale
			output.YesterdaySaleClothes = yesterdaySaleClothes

			output.AllGoldSale = allGoldSale
			output.AllPxcSale = allPxcSale
			output.AllSaleClothes = allSaleClothes

			output.List = list
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func GetStatusLandHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input GetStatusLandReq
	var output GetStatusLandResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		yd := time.Now().AddDate(0, 0, -1).Format("20060102")
		if input.LandType == 1 {
			//营业中
			if err, list := land.GetLandByStatusMastername(LAND_STATUS_WB_IN_BUSINESS, player.Username, yd); err != nil {
				resp = PIXIE_ERR_GET_STATUS_LAND_WRONG
			} else {
				output.List = list
			}
		} else if input.LandType == 2 {
			//装修中
			if err, list := land.GetLandByStatusMastername(LAND_STATUS_NB_BUILDING, player.Username, yd); err != nil {
				resp = PIXIE_ERR_GET_STATUS_LAND_WRONG
			} else {
				output.List = list
			}
		} else if input.LandType == 3 {
			//有建筑 没开业
			if err, list := land.GetLandWBWaitBusinessByMastername(player.Username, yd); err != nil {
				resp = PIXIE_ERR_GET_STATUS_LAND_WRONG
			} else {
				output.List = list
			}
		} else if input.LandType == 4 {
			//空闲 还没装修
			if err, list := land.GetLandNBWaitBuildByMastername(player.Username, yd); err != nil {
				resp = PIXIE_ERR_GET_STATUS_LAND_WRONG
			} else {
				output.List = list
			}
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func GoStoreHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input GoStoreReq
	var output GoStoreResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		//设计师商店 如果是在售卖
		if strings.HasPrefix(input.ClothesID, "O-") {
			nt := time.Now().Unix()
			paperIDStr := input.ClothesID[2:]
			paperID, _ := strconv.ParseInt(paperIDStr, 10, 64)
			if ownerUsername, err := paper.GetMinPricePaperByPaperID(paperID); err != nil {

			} else {
				dland.GetLandByRenterUsername(ownerUsername, paperID, nt)
			}
		} else {
			output.Res = spaper.GetAllOfficialPaper()
		}

		//官方商店

		result, parseErr = PixieMarshal(output)
	}
	return
}

func GetSaleClothesHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input GetSaleClothesReq
	var output GetSaleClothesResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if list, err := paper.GetClothesInSale(player.Username, input.ClothesType, input.PartType, input.Page, input.PageSize); err != nil {
			resp = PIXIE_ERR_GET_SALE_CLOTHES
		} else {
			output.List = list
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}
