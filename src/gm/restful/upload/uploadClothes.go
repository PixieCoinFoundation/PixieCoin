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
	"service_pixie/paper"
	"service_pixie/paperInfo"
	"service_pixie/record"

	"strconv"
	"strings"
	"time"
)

type UploadClothesReq struct {
	Cname       string `json:"cname"`       //背景名称
	Desc        string `json:"desc"`        //背景描述
	ClothesType string `json:"clothesType"` //模特类型
	PartType    string `json:"partType"`    //服装类别
	CTypeLayer  string `json:"cTypeLayer"`  //服装层级

	FileMain   string `json:"fileMain"`   //服装main (必传）
	FileIcon   string `json:"fileIcon"`   //图标 (必传）
	FileBottom string `json:"fileBottom"` //下层层
	FileCollar string `json:"fileCollar"` //领子
	FileShadow string `json:"fileShadow"` //阴影

	PriceType string `json:"priceType"` // 货币类型
	Pirce     string `json:"price"`     //货币类型
	Star      string `json:"star"`      //星级
	// 标签
	Tag1  string `json:"tag_1"`
	Tag2  string `json:"tag_2"`
	Tag3  string `json:"tag_3"`
	Tag4  string `json:"tag_4"`
	Tag5  string `json:"tag_5"`
	Tag6  string `json:"tag_6"`
	Tag7  string `json:"tag_7"`
	Tag8  string `json:"tag_8"`
	Tag9  string `json:"tag_9"`
	Tag10 string `json:"tag_10"`
	Tag11 string `json:"tag_11"`
	Tag12 string `json:"tag_12"`
	Tag13 string `json:"tag_13"`
	Tag14 string `json:"tag_14"`
	Tag15 string `json:"tag_15"`
	Tag16 string `json:"tag_16"`
	Tag17 string `json:"tag_17"`
	Tag18 string `json:"tag_18"`
	Tag19 string `json:"tag_19"`
	Tag20 string `json:"tag_20"`

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

func UploadClothes(ctx *web.Context) {
	UploadReqLock.Lock()
	defer UploadReqLock.Unlock()
	// if exist, _ := ctx.CheckToken(); !exist {
	// 	ctx.BackError(1000, "没有权限,请重新登录")
	// 	return
	// }
	var (
		input                                                                      UploadClothesReq
		style                                                                      PaperStyle      //风格
		oPaper                                                                     GMOfficialPaper //上传图纸信息
		styleCount                                                                 int
		isExist                                                                    bool
		tagList                                                                    []string
		filePathMain, filePathBottom, filePathCollar, filePathShadow, filePathIcon string
		err                                                                        error
	)
	tNow := time.Now().Unix()

	//1.处理必填参数
	if err = ctx.GetFormData(&input); err != nil {
		ctx.BackError(10001, "参数错误")
		return
	}
	if input.Cname == "" || input.Desc == "" || input.ClothesType == "" || input.PartType == "" || strings.TrimSpace(input.Cname) == "" || strings.TrimSpace(input.Desc) == "" {
		ctx.BackError(10002, "必填参数不能为空")
		return
	}
	oPaper.Cname = input.Cname
	oPaper.Desc = input.Desc
	oPaper.ClothesType = input.ClothesType
	oPaper.PartType = input.PartType

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
	zIndex, _ := strconv.Atoi(input.CTypeLayer)

	zValueList, exist := PAPER_TYPE_ZVALUE_MAP[input.PartType]
	if !exist {
		ctx.BackError(100013, "部位有误")
		return
	} else {
		if len(zValueList) <= 0 {
			ctx.BackError(100013, "服务器内部数据错误")
			return
		} else if len(zValueList) == 1 {
			isExist = true
			zIndex = zValueList[0]
		} else {
			for _, zv := range zValueList {
				if zv == zIndex {
					isExist = true
					break
				}
			}
		}
	}

	if !isExist {
		ctx.BackError(100013, "层级有误")
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
	hiddenTagMap := paperInfo.GetSTagMap()
	if input.HiddenTag != "" && hiddenTagMap[input.HiddenTag] == "" {
		ctx.BackError(10025, "hidden tag 有误")
		return
	}

	//4.上传图片
	if filePathMain, err = ctx.UploadFormFile("fileMain", fmt.Sprintf("official_main_%d", tNow), false); err != nil || filePathMain == "" {
		ctx.BackError(10016, "main图片不能为空")
		return
	}
	if filePathBottom, err = ctx.UploadFormFile("fileBottom", fmt.Sprintf("official_bottom_%d", tNow), false); err != nil && err != http.ErrMissingFile {
		ctx.BackError(10017, "bottom图片出错")
		return
	}
	if filePathCollar, err = ctx.UploadFormFile("fileCollar", fmt.Sprintf("official_collar_%d", tNow), false); err != nil && err != http.ErrMissingFile {
		ctx.BackError(10018, "collar图片出错")
		return
	}
	if filePathShadow, err = ctx.UploadFormFile("fileShadow", fmt.Sprintf("official_shadow_%d", tNow), false); err != nil && err != http.ErrMissingFile {
		ctx.BackError(10019, "shadow图片出错")
		return
	}
	if filePathIcon, err = ctx.UploadFormFile("fileIcon", fmt.Sprintf("official_icon_%d", tNow), false); err != nil || filePathIcon == "" {
		ctx.BackError(10020, "icon图片不能为空")
		return
	}
	var extra PaperExtra
	//todo图层未定

	if err = paper.ProcessFile("main", &extra, filePathIcon, filePathMain, filePathBottom, filePathCollar, filePathShadow, tNow); err != nil {
		ctx.BackError(10021, "上传图片格式、大小有误")
		return
	}

	extra.ZValue = zIndex
	cloSet := PaperFile{
		Icon:   filePathIcon,
		Main:   filePathMain,
		Bottom: filePathBottom,
		Collar: filePathCollar,
		Shadow: filePathShadow,
	}

	paperFile, _ := json.Marshal(cloSet)

	eb, _ := json.Marshal(extra)
	stypeStr, _ := json.Marshal(style)
	if len(tagList) >= 1 {
		oPaper.Tag1 = tagList[0]
	}
	if len(tagList) >= 2 {
		oPaper.Tag2 = tagList[1]
	}
	oPaper.PaperFile = string(paperFile)
	oPaper.PaperExtra = string(eb)
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
