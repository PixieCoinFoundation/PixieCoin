package clothes

import (
	"fmt"
	"github.com/web"
	"gm/dao"
	. "model"
	. "pixie_contract/api_specification"
	"strconv"
	"strings"
)

type GetPlayerClothesResp struct {
	List []DesignPaper
}

func GetPlayerClothes(ctx *web.Context, page string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}

	var output GetPlayerClothesResp
	if err, list := dao.GetPlayerClothes(1); err != nil {
		ctx.BackError(30001, err)
		return
	} else {
		output.List = list
	}
	ctx.BackSuccess(output)
}

type GetPlayerBGResp struct {
	List []DesignPaper
}

func GetPlayerBG(ctx *web.Context, page string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	var output GetPlayerBGResp
	if err, list := dao.GetPlayerBG(1); err != nil {
		ctx.BackError(30006, err)
		return
	} else {
		output.List = list
	}
	ctx.BackSuccess(output)
}

type GetInSaleClothesResp struct {
	Data  []SalePaperGM `json:"data"`
	Count int           `json:"count"`
	Code  int           `json:"code"`
	Msg   string        `json:"msg"`
}

func GetInSaleClothes(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var output GetInSaleClothesResp
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	page, _ := strconv.Atoi(ctx.Params["page"])
	limit, _ := strconv.Atoi(ctx.Params["limit"])
	if page < 1 {
		page = 1
	}
	if limit < 10 {
		limit = 10
	}
	offset := (page - 1) * limit

	if err, count, list := dao.GetSaleClothes(offset, limit); err != nil {
		ctx.BackError(30006, err)
		return
	} else {
		output.Data = list
		output.Count = int(count)
	}
	ctx.BackSuccess(output)
	return
}

type GetPoolClothesResp struct {
	Data  []SalePaperGM `json:"data"`
	Count int           `json:"count"`
	Code  int           `json:"code"`
	Msg   string        `json:"msg"`
}

func GetPoolClothes(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var output GetPoolClothesResp
	var OidList, PidList []string

	// if exist, _ := ctx.CheckToken(); !exist {
	// 	ctx.BackError(1000, "没有权限,请重新登录")
	// 	return
	// }
	page, _ := strconv.ParseInt(ctx.Params["page"], 10, 64)
	limit, _ := strconv.ParseInt(ctx.Params["limit"], 10, 64)
	if page < 1 {
		page = 1
	}
	if limit < 10 {
		limit = 10
	}
	offset := (page - 1) * limit

	configNameMap := map[string]string{}
	if err, list := dao.GetPoolStatusClothes(0, offset, limit); err != nil {
		ctx.BackError(30006, err)
		return
	} else {
		for _, v := range list {
			if strings.HasPrefix(v.ClothesID, "O-") {
				OidList = append(OidList, v.ClothesID[2:])
			} else {
				PidList = append(PidList, v.ClothesID[2:])
				configNameMap[v.ClothesID] = v.OwnerUsername
			}
		}
		if err, data := dao.FindSalePaper(PidList); err != nil {
			ctx.BackError(30007, err)
			return
		} else {
			for _, v := range data {
				v.PaperID = fmt.Sprintf("P-%s", v.PaperID)
				if val, exist := configNameMap[v.PaperID]; exist {
					if val == v.OwnerUsername {
						output.Data = append(output.Data, v)
					}
				} else {
					output.Data = append(output.Data, v)
				}
			}
			if err, data1 := dao.FindOffPaper(OidList); err != nil {
				ctx.BackError(30008, err)
				return
			} else {
				for _, v := range data1 {
					var temp SalePaperGM
					temp.ID = v.ID
					temp.PaperID = fmt.Sprintf("O-%d", v.ID)
					temp.Cname = v.Cname
					temp.ClothesType = v.ClothesType
					temp.PartType = v.PartType
					temp.Price = v.Price
					temp.PriceType = v.PriceType
					temp.Star = v.Star
					temp.Tag1 = v.Tag1
					temp.Tag2 = v.Tag2
					temp.Style = v.Style
					temp.STag = v.STag
					temp.Extra = v.Extra
					output.Data = append(output.Data, temp)
				}
			}
			output.Count = len(data)
		}
	}

	ctx.BackSuccess(output)
	return
}

func GetPoolClothes1(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var output GetPoolClothesResp
	var OidList, PidList []string
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	page, _ := strconv.ParseInt(ctx.Params["page"], 10, 64)
	limit, _ := strconv.ParseInt(ctx.Params["limit"], 10, 64)
	if page < 1 {
		page = 1
	}
	if limit < 10 {
		limit = 10
	}
	offset := (page - 1) * limit
	if err, list := dao.GetPoolStatusClothes(1, offset, limit); err != nil {
		ctx.BackError(30006, err)
		return
	} else {
		configNameMap := map[string]string{}
		for _, v := range list {
			if strings.HasPrefix(v.ClothesID, "O-") {
				OidList = append(OidList, v.ClothesID[2:])
			} else {
				PidList = append(PidList, v.ClothesID[2:])
				if v.OwnerUsername != "" {
					configNameMap[v.ClothesID] = v.OwnerUsername
				}
			}
		}
		if err, data := dao.FindSalePaper(PidList); err != nil {
			ctx.BackError(30007, err)
			return
		} else {
			for _, v := range data {
				v.PaperID = fmt.Sprintf("P-%s", v.PaperID)
				if val, exist := configNameMap[v.PaperID]; exist {
					if val == v.OwnerUsername {
						output.Data = append(output.Data, v)
					}
				} else {
					output.Data = append(output.Data, v)
				}
			}
			if err, data1 := dao.FindOffPaper(OidList); err != nil {
				ctx.BackError(30008, err)
				return
			} else {
				for _, v := range data1 {
					var temp SalePaperGM
					temp.ID = v.ID
					temp.PaperID = fmt.Sprintf("O-%d", v.ID)
					temp.Cname = v.Cname
					temp.ClothesType = v.ClothesType
					temp.PartType = v.PartType
					temp.Price = v.Price
					temp.PriceType = v.PriceType
					temp.Star = v.Star
					temp.Tag1 = v.Tag1
					temp.Tag2 = v.Tag2
					temp.Style = v.Style
					temp.STag = v.STag
					temp.Extra = v.Extra
					output.Data = append(output.Data, temp)
				}
			}
			output.Count = len(data)
		}
	}
	ctx.BackSuccess(output)
	return
}
