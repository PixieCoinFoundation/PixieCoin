package upload

import (
	"constants"
	"encoding/json"
	"fmt"
	"github.com/web"
	"gm/dao"
	. "model"
	"net/http"
	. "pixie_contract/api_specification"
	"service_pixie/paperInfo"
	"service_pixie/record"
	"sync"

	"strconv"
	"strings"
	"time"
)

var UploadReqLock sync.Mutex

type UploadBgReq struct {
	Cname string `json:"cname"` //背景名称
	Desc  string `json:"desc"`  //背景描述

	FileBg   string `json:"fileBg"`   //背景(必传）
	FileIcon string `json:"fileIcon"` //icon(必传）

	FileFront string `json:"fileFront"` //背景前景

	PriceType string `json:"priceType"` // 货币类型
	Pirce     string `json:"price"`     //货币类型
	Star      string `json:"star"`      //星级

	// 风格
	Warm     string `json:"warm"`     //温暖
	Formal   string `json:"formal"`   //正式
	Tight    string `json:"tight"`    //修身
	Bright   string `json:"bright"`   //鲜艳
	Cute     string `json:"cute"`     //可爱
	Man      string `json:"man"`      //帅气
	Tough    string `json:"tough"`    //粗旷
	Sexy     string `json:"sexy"`     //性感
	Gorgeous string `json:"gorgeous"` //华丽
	Noble    string `json:"noble"`    //高贵
	Strange  string `json:"strange"`  //另类
	Sprot    string `json:"sport"`    //运动

	UnlockCondition string `json:"unlock_condition"` //解锁关卡条件
	HiddenTag       string `json:"hiddenTag"`        // 暗标签
	StartTime       string `json:"startTime"`        //上架时间
}

func UploadBg(ctx *web.Context) {
	// if exist, _ := ctx.CheckToken(); !exist {
	// 	ctx.BackError(1000, "没有权限,请重新登录")
	// 	return
	// }
	UploadReqLock.Lock()
	defer UploadReqLock.Unlock()

	var (
		input                       UploadClothesReq
		style                       PaperStyle      //风格
		oPaper                      GMOfficialPaper //上传图纸信息
		styleCount                  int
		tagList                     []string
		fileFront, fileBg, fileIcon string
		err                         error
	)
	tNow := time.Now().Unix()

	//1.处理必填参数
	if err = ctx.GetFormData(&input); err != nil {
		ctx.BackError(10001, "参数错误")
		return
	}
	if input.Cname == "" || input.Desc == "" || strings.TrimSpace(input.Cname) == "" || strings.TrimSpace(input.Desc) == "" {
		ctx.BackError(10002, "必填参数不能为空")
		return
	}
	oPaper.Cname = input.Cname
	oPaper.Desc = input.Desc
	oPaper.PartType = PAPER_TYPE_BEIJING
	// oPaper.ClothesType =
	priceType, _ := strconv.Atoi(input.PriceType)
	price, _ := strconv.ParseInt(input.Pirce, 10, 64)
	star, _ := strconv.Atoi(input.Star)

	if priceType != constants.PIXIE_GOLD_TYPE && priceType != constants.PIXIE_PXC_TYPE {
		ctx.BackError(10003, "货币参数有误")
		return
	}
	//10-10000 pxc限制：10-100000
	if priceType == constants.PIXIE_GOLD_TYPE {
		if price < 10 || price > 10000 {
			ctx.BackError(10004, "价格参数有误")
			return
		}
	}
	if input.StartTime == "" {
		ctx.BackError(10006, "时间参数有误")
		return
	}
	tStart, _ := time.ParseInLocation("2006-01-02 15:04:05", input.StartTime, time.Local)

	if priceType == constants.PIXIE_PXC_TYPE {
		if price < 10 || price > 100000 {
			ctx.BackError(10009, "价格参数有误")
			return
		}
	}
	if !record.LevelIDLegal(input.UnlockCondition) && input.UnlockCondition != "" {
		ctx.BackError(100011, "解锁条件有误")
		return
	}

	if star < 1 || star > 5 {
		ctx.BackError(100012, "星级有误")
		return
	}

	oPaper.PriceType = priceType
	oPaper.Price = price
	oPaper.Star = star

	// 2.风格处理
	if input.Warm != "" {
		styleCount++
		temp, _ := strconv.Atoi(input.Warm)
		style.Warm = temp
	}
	if input.Formal != "" {
		styleCount++
		temp, _ := strconv.Atoi(input.Formal)
		style.Formal = temp
	}
	if input.Tight != "" {
		styleCount++
		temp, _ := strconv.Atoi(input.Tight)
		style.Tight = temp
	}
	if input.Bright != "" {
		styleCount++
		temp, _ := strconv.Atoi(input.Bright)
		style.Bright = temp
	}
	if input.Cute != "" {
		styleCount++
		temp, _ := strconv.Atoi(input.Cute)
		style.Cute = temp
	}
	if input.Man != "" {
		styleCount++
		temp, _ := strconv.Atoi(input.Man)
		style.Man = temp
	}
	if input.Tough != "" {
		styleCount++
		temp, _ := strconv.Atoi(input.Tough)
		style.Tough = temp
	}
	if input.Sexy != "" {
		styleCount++
		temp, _ := strconv.Atoi(input.Sexy)
		style.Sexy = temp
	}
	if input.Gorgeous != "" {
		styleCount++
		temp, _ := strconv.Atoi(input.Gorgeous)
		style.Gorgeous = temp
	}
	if input.Noble != "" {
		styleCount++
		temp, _ := strconv.Atoi(input.Noble)
		style.Noble = temp
	}
	if input.Strange != "" {
		styleCount++
		temp, _ := strconv.Atoi(input.Strange)
		style.Strange = temp
	}
	if input.Sprot != "" {
		styleCount++
		temp, _ := strconv.Atoi(input.Sprot)
		style.Sport = temp
	}
	if styleCount < 1 || styleCount > 2 {
		ctx.BackError(10014, "style数量有误")
		return
	}

	tagMap := paperInfo.GetTagMap()
	//3.tag处理
	for _, v := range ctx.GetTagList() {
		if tagMap[v] != "" {
			tagList = append(tagList, v)
		}
	}
	if len(tagList) > 2 {
		ctx.BackError(10015, " 标签有误")
		return
	}
	//4.上传图片
	if fileFront, err = ctx.UploadFormFile("fileFront", fmt.Sprintf("official_front_%d", tNow), false); err != nil && err != http.ErrMissingFile {
		ctx.BackError(10016, "front图片出错")
		return
	}
	if fileIcon, err = ctx.UploadFormFile("fileIcon", fmt.Sprintf("official_icon_%d", tNow), false); err != nil || fileIcon == "" {
		ctx.BackError(10016, "icon图片出错")
		return
	}
	if fileBg, err = ctx.UploadFormFile("frontBg", fmt.Sprintf("official_back_%d", tNow), false); err != nil || fileBg == "" {
		ctx.BackError(10017, "bg图片出错")
		return
	}
	cloSet := PaperFile{
		Front: fileFront,
		Back:  fileBg,
		Icon:  fileIcon,
	}

	paperFile, _ := json.Marshal(cloSet)
	stypeStr, _ := json.Marshal(style)

	if len(tagList) >= 1 {
		oPaper.Tag1 = tagList[0]
	}
	if len(tagList) >= 2 {
		oPaper.Tag2 = tagList[1]
	}
	oPaper.PaperFile = string(paperFile)
	oPaper.HiddenTag = input.HiddenTag
	oPaper.PaperStyle = string(stypeStr)
	oPaper.UploadTime = tNow
	oPaper.SaleTime = tStart.Unix()
	oPaper.UnlockLevel = input.UnlockCondition
	oPaper.AdminName = "未定"

	if err := dao.InsertOfficialPaper(oPaper); err != nil {
		ctx.BackError(10023, err)
		return
	}
	ctx.BackSuccess("success")
}
