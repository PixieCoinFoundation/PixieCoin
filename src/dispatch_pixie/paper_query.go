package dispatch_pixie

import (
	"constants"
	dp "db_pixie/paper"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service_pixie/paper"
	"time"
	"tools"
)

func GetOneDesignPaperHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input GetOneDesignPaperReq
	var output GetOneDesignPaperResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if err, paper := dp.GetDesignPaperByID(input.PaperID); err == nil {
			output.Paper = paper
			if paper.Status == int(PAPER_STATUS_ADMIN_QUEUE) {
				if cnt, err := dp.GetAdminCheckCountBefore(input.PaperID); err == nil {
					output.AdminCheckBefore = cnt
				} else {
					resp = PIXIE_ERR_GET_ADMIN_CHECK_CNT
				}
			}
		} else {
			resp = PIXIE_ERR_GET_PAPER_BY_ID
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}

func ListDesignPaperHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input ListDesignPaperReq
	var output ListDesignPaperResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		tu := input.TargetUsername
		if tu == "" {
			tu = player.Username
		}

		page := 1
		if input.Page > 0 {
			page = input.Page
		}

		pageSize := constants.GET_PAPER_DEFAULT_PAGE_SIZE
		if input.PageSize > 0 {
			pageSize = input.PageSize
		}

		if err, list := paper.ListPaperByAuthor(tu, page, pageSize); err != nil {
			resp = PIXIE_ERR_LIST_DESIGN_PAPER
		} else {
			output.List = list
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}

func ListOwnPaperHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input ListOwnPaperReq
	var output ListOwnPaperResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if input.PartType != "" && PixiePaperPartTypeMap[input.PartType] == "" {
			resp = PIXIE_ERR_PARAM_ILLEGAL
		} else {
			tu := input.TargetUsername
			if tu == "" {
				tu = player.Username
			}

			page := 1
			if input.Page > 0 {
				page = input.Page
			}

			pageSize := constants.GET_PAPER_DEFAULT_PAGE_SIZE
			if input.PageSize > 0 {
				pageSize = input.PageSize
			}

			nt := time.Now().Unix()
			if err, list := paper.ListPaperByOwner(tu, input.ClothesType, input.PartType, page, pageSize, nt); err != nil {
				resp = PIXIE_ERR_LIST_OWN_PAPER
			} else {
				output.List = list
			}
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}

func ListSpecifiedOwnPaperHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input ListSpecifiedOwnPaperReq
	var output ListSpecifiedOwnPaperResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if len(input.ClothesIDList) > 0 {
			paperIDList := make([]int64, 0)

			for _, cid := range input.ClothesIDList {
				if src, id := tools.GetPixieClothesDetail(cid); src == string(CLOTHES_ORIGIN_PLAYER_PAPER) {
					paperIDList = append(paperIDList, id)
				} else {
					resp = PIXIE_ERR_PARAM_ILLEGAL
					goto END
				}
			}

			tu := input.TargetUsername
			if tu == "" {
				tu = player.Username
			}

			if err, list := dp.ListPaperOwnByOwnerIDList(paperIDList, tu); err != nil {
				resp = PIXIE_ERR_LIST_SPECIFIED_OWN_PAPER
			} else {
				output.List = list
			}
		} else {
			resp = PIXIE_ERR_PARAM_ILLEGAL
		}

	END:
		result, parseErr = PixieMarshal(output)
	}

	return
}

func ListSalePaperHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input ListSalePaperReq
	var output ListSalePaperResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if err, list := paper.ListSalePaper(input.Page, input.PageSize, input.SortType, player.Username); err != nil {
			resp = PIXIE_ERR_LIST_SALE_PAPER
		} else {
			output.List = list
		}

		result, parseErr = PixieMarshal(output)
	}

	return
}

func QueryPaperSaleInfoHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input QueryPaperSaleInfoReq
	var output QueryPaperSaleInfoResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if input.PaperID <= 0 {
			resp = PIXIE_ERR_PARAM_ILLEGAL
		} else if tid, pt, st, dr, maxp, minp, err := paper.QueryPaperSaleInfo(input.PaperID, player.Username, input.Sequence); err != nil {
			if err == constants.TradeNotExist {
				resp = PIXIE_ERR_TRADE_NOT_EXIST
			} else {
				resp = PIXIE_ERR_QUERY_PAPER_SALE
			}

		} else {
			output.TradeID = tid
			output.PriceType = pt
			output.StartTime = st
			output.Duration = dr
			output.MaxPrice = maxp
			output.MinPrice = minp
		}

		result, parseErr = PixieMarshal(output)
	}

	return

}

func ListPaperTradeHistoryHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input ListPaperTradeHistoryReq
	var output ListPaperTradeHistoryResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if input.Sequence < 0 || input.Sequence > constants.PAPER_CIRCULATION_MAX {
			resp = PIXIE_ERR_PARAM_ILLEGAL
		} else if err, res, ts := paper.QueryPaperTradeHistory(input.PaperID, input.Sequence, input.Page); err != nil {
			resp = PIXIE_ERR_QUERY_PAPER_TRADE_HISTORY
		} else {
			output.List = res
			output.TotalSale = ts
		}

		result, parseErr = PixieMarshal(output)
	}

	return

}
