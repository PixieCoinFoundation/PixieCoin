package clothes

import (
	"fmt"
	"github.com/web"
	"gm/common"
	"gm/dao"
	. "model"
	"service_pixie/paperInfo"
	"strconv"
	"time"
)

type GetOfficialClothesResp struct {
	Data  []GMOfficialPaperTrans `json:"data"`
	Count int                    `json:"count"`
	Code  int                    `json:"code"`
	Msg   string                 `json:"msg"`
}

func GetOfficialClothes(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
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

	var output GetOfficialClothesResp

	err, count, res := dao.GetOfficialPaper(int(offset), int(limit))

	if err != nil {
		ctx.BackError(20001, err)
		return
	}
	output.Data = commonTrans(res, false)
	output.Count = int(count)
	ctx.BackSuccess(output)
}

func commonTrans(res []GMOfficialPaper, isBg bool) (list []GMOfficialPaperTrans) {
	hiddenTagMap := paperInfo.GetSTagMap()
	for _, v := range res {
		var temp GMOfficialPaperTrans
		temp.ID = v.ID
		temp.Cname = v.Cname
		temp.Desc = v.Desc
		//模特类型
		if !isBg {
			cType, exist := common.ClothesTypeStringMap[v.ClothesType]
			if exist {
				temp.ClothesType = cType
			} else {
				temp.ClothesType = "无"
			}
			temp.PartType = common.PartTypeStringMap[v.PartType]
		}
		//价格
		temp.PriceType = common.CurrencyStringMap[v.PriceType]
		temp.Price = fmt.Sprintf("%d", v.Price)
		temp.Star = fmt.Sprintf("%d星", v.Star)

		if len(v.Tag1) > 0 {
			temp.Tag1 = paperInfo.GetTagMap()[v.Tag1]
		}
		if len(v.Tag2) > 0 {
			temp.Tag2 = paperInfo.GetTagMap()[v.Tag2]
		}
		temp.PaperStyle = v.PaperStyle
		temp.UnlockLevel = v.UnlockLevel
		temp.HiddenTag = hiddenTagMap[v.HiddenTag]
		temp.SaleTime = time.Unix(v.SaleTime, 0).Format("2006-01-02 15:04:05")
		temp.UploadTime = time.Unix(v.UploadTime, 0).Format("2006-01-02 15:04:05")
		temp.AdminName = v.AdminName
		list = append(list, temp)
	}
	return
}

func GetOfficialBG(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
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

	var output GetOfficialClothesResp

	err, count, res := dao.GetOfficialBGPaper(int(offset), int(limit))

	if err != nil {
		ctx.BackError(20006, err)
		return
	}
	output.Data = commonTrans(res, true)
	output.Count = int(count)

	ctx.BackSuccess(output)
}

type DeleteOfficialReq struct {
	PaperID string `json:"id"`
}

func DeleteOfficial(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}

	var input DeleteOfficialReq
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(20011, err)
		return
	}
	paperID, _ := strconv.ParseInt(input.PaperID, 10, 64)
	if paperID <= 0 {
		ctx.BackError(20012, "参数不能为空")
		return
	}
	if err, paper := dao.GetOneOfficialPaper(paperID); err != nil {
		ctx.BackError(20013, err)
		return
	} else {
		if paper.SaleTime < time.Now().Unix() || (paper.SaleTime-time.Now().Unix()) < 3600 {
			ctx.BackError(20014, "衣服已经上架或者时间间隔太小,请联系管理员")
			return
		}
	}
	if err := dao.DeleteOfficialPaper(paperID); err != nil {
		ctx.BackError(20015, err)
		return
	}
	ctx.BackSuccess("success")
}
