package recommend

import (
	"appcfg"
	"constants"
	"encoding/json"
	"fmt"
	"github.com/web"
	"gm/common"
	"gm/dao"
	. "model"
	. "pixie_contract/api_specification"
	"strconv"
	"strings"
	"time"
)

type AddSubjectReq struct {
	Name     string `json:"name"`
	OpenTime string `json:"open_time"`

	HotSpot1X     string `json:"hotSpot1X"`
	HotSpot1Y     string `json:"hotSpot1Y"`
	HotSpot1W     string `json:"hotSpot1W"`
	HotSpot1H     string `json:"hotSpot1H"`
	HotSpot1Type  string `json:"jumpType1"`
	HotSpot1Value string `json:"jumpValue1"`

	HotSpot2X     string `json:"hotSpot2X"`
	HotSpot2Y     string `json:"hotSpot2Y"`
	HotSpot2W     string `json:"hotSpot2W"`
	HotSpot2H     string `json:"hotSpot2H"`
	HotSpot2Type  string `json:"jumpType2"`
	HotSpot2Value string `json:"jumpValue2"`

	HotSpot3X     string `json:"hotSpot3X"`
	HotSpot3Y     string `json:"hotSpot3Y"`
	HotSpot3W     string `json:"hotSpot3W"`
	HotSpot3H     string `json:"hotSpot3H"`
	HotSpot3Type  string `json:"jumpType3"`
	HotSpot3Value string `json:"jumpValue3"`

	HotSpot4X     string `json:"hotSpot4X"`
	HotSpot4Y     string `json:"hotSpot4Y"`
	HotSpot4W     string `json:"hotSpot4W"`
	HotSpot4H     string `json:"hotSpot4H"`
	HotSpot4Type  string `json:"jumpType4"`
	HotSpot4Value string `json:"jumpValue4"`
}

func AddSubject(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var input AddSubjectReq
	var saveInput RecommendSubjectGM
	var mainBanner string
	var list []HSpot
	var err error

	if err = ctx.GetFormData(&input); err != nil {
		ctx.BackError(1009, err)
		return
	}
	if legal, spot := paramsHotSpot(input.HotSpot1X, input.HotSpot1Y, input.HotSpot1W, input.HotSpot1H, input.HotSpot1Type, input.HotSpot1Value); legal {
		list = append(list, spot)
	}

	if legal, spot := paramsHotSpot(input.HotSpot2X, input.HotSpot2Y, input.HotSpot2W, input.HotSpot2H, input.HotSpot2Type, input.HotSpot2Value); legal {
		list = append(list, spot)
	}
	if legal, spot := paramsHotSpot(input.HotSpot3X, input.HotSpot3Y, input.HotSpot3W, input.HotSpot3H, input.HotSpot3Type, input.HotSpot3Value); legal {
		list = append(list, spot)
	}
	if legal, spot := paramsHotSpot(input.HotSpot4X, input.HotSpot4Y, input.HotSpot4W, input.HotSpot4H, input.HotSpot4Type, input.HotSpot4Value); legal {
		list = append(list, spot)
	}

	tOpen, _ := time.ParseInLocation("2006-01-02 15:04:05", input.OpenTime, time.Local)
	if mainBanner, err = ctx.UploadFormFile("banner_img", fmt.Sprintf("recommend_subject_%d_main", time.Now().Unix()), false); err != nil || mainBanner == "" {
		ctx.BackError(10020, "banner图片不能为空")
		return
	}
	data, _ := json.Marshal(list)

	saveInput.Name = input.Name
	saveInput.BannerFileID = mainBanner
	saveInput.HotSpot = string(data)
	saveInput.OpenTime = tOpen.Unix()
	if err = dao.AddSubject(saveInput); err != nil {
		ctx.BackError(1009, "参数有误")
		return
	}
	ctx.BackSuccess("success")
	return
}

type DeleteSubjectReq struct {
	SubjectID string `json:"subject_id"`
}

func DeleteSubject(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var input DeleteSubjectReq
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(1009, err)
		return
	}
	subjectId, _ := strconv.ParseInt(input.SubjectID, 10, 64)
	if err := dao.DelSubject(subjectId); err != nil {
		ctx.BackError(1010, err)
		return
	}
	ctx.BackSuccess("success")
}

func paramsHotSpot(x, y, w, h, htype, hvalue string) (legal bool, spot HSpot) {
	if x != "" && y != "" && w != "" && h != "" && htype != "" && hvalue != "" {
		x, _ := strconv.ParseInt(x, 10, 64)
		y, _ := strconv.ParseInt(y, 10, 64)
		w, _ := strconv.ParseInt(w, 10, 64)
		h, _ := strconv.ParseInt(h, 10, 64)
		spot.X = x
		spot.Y = y
		spot.Width = w
		spot.Height = h
		spot.HType = htype
		spot.HLink = hvalue
		legal = true
		return
	}
	return
}

type UpdateSuitClothesReq struct {
	SubjectID   string `json:"subject_id"`
	SuitOwnerID string `json:"suit_owner_id"`
	OSuitIDList string `json:"o_suit_id_list"`
	PSuitIDList string `json:"p_suit_id_list"`
}

func UpdateSuit(ctx *web.Context) {
	var input UpdateSuitClothesReq
	var saveSuit []string
	var oList, pList []interface{}
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(1009, err)
		return
	}
	if input.OSuitIDList == "" && input.PSuitIDList == "" {
		ctx.BackError(1010, "参数不能为空")
		return
	}
	Oidlist := strings.Split(input.OSuitIDList, ",")
	Pidlist := strings.Split(input.PSuitIDList, ",")
	if len(Oidlist) > 1 {
		for _, v := range Oidlist {
			oList = append(oList, v)
			saveSuit = append(saveSuit, fmt.Sprintf("O-%s", v))
		}
		if err, exist := dao.OclohtesExist(oList); err != nil {
			ctx.BackError(1001, err)
			return
		} else if !exist {
			ctx.BackError(1012, "官方衣服不存在")
			return
		}
	}
	if len(Pidlist) > 1 {
		for _, v := range Pidlist {
			pList = append(pList, v)
			saveSuit = append(saveSuit, fmt.Sprintf("P-%s", v))
		}
		if err, exist := dao.PclohtesExist(pList); err != nil {
			ctx.BackError(1013, err)
			return
		} else if !exist {
			ctx.BackError(1014, "设计师衣服不存在")
			return
		}
	}
	if len(Oidlist)+len(Pidlist) < 2 {
		ctx.BackError(1005, "参数有误")
		return
	}

	data, _ := json.Marshal(saveSuit)
	subjectID, _ := strconv.ParseInt(input.SubjectID, 10, 64)

	if err := dao.UpdateSuit(string(data), input.SuitOwnerID, subjectID); err != nil {
		ctx.BackError(1016, err)
		return
	}
	ctx.BackSuccess("success")
}

func legalClothesType(idList []interface{}) (legal bool) {
	if err, papers := dao.FindPaperByIDList(idList); err != nil {
		return
	} else {
		for _, v := range papers {
			fmt.Println(v)
		}
	}
	return
}

func GetOneSubject(ctx *web.Context, idStr string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")

	var id int64
	if err, subject := dao.GetSubject(id); err != nil {
		ctx.BackError(1009, "参数有误")
		return
	} else {
		ctx.BackSuccess(subject)
		return
	}
	return
}

type GetAllSubjectResp struct {
	Data  []RecommendSubjectGM `json:"data"`
	Count int                  `json:"count"`
	Code  int                  `json:"code"`
	Msg   string               `json:"msg"`
}

func GetAllSubject(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")

	var output GetAllSubjectResp
	if err, subjects := dao.FindSubject(); err != nil {
		ctx.BackError(1009, "参数有误")
		return
	} else {
		output.Data = subjects
		output.Count = len(subjects)
		ctx.BackSuccess(output)
		return
	}
	return
}

type AddClothesWithLinkReq struct {
	SubjectID   string `json:"subject_id"`
	OSuitIDList string `json:"o_suit_id_list"`
	PSuitIDList string `json:"p_suit_id_list"`
}

//推荐跑马灯
func AddClothesWithLink(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Request.ParseMultipartForm((1 << 20) * 10)

	var input AddClothesWithLinkReq
	var saveInput RecommendItemGM
	var banner, h5link string
	var err error
	var idList []int64
	var idListStr []string

	if err = ctx.GetFormData(&input); err != nil {
		ctx.BackError(1009, err)
		return
	}
	olist := strings.Split(input.OSuitIDList, ",")
	plist := strings.Split(input.PSuitIDList, ",")

	if (len(olist) + len(plist)) < 1 {
		ctx.BackError(10010, "配置的衣服不能为空")
		return
	}
	for _, v := range olist {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			idList = append(idList, id)
			idListStr = append(idListStr, fmt.Sprintf("O-%s", v))
		}
	}
	for _, v := range plist {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			idList = append(idList, id)
			idListStr = append(idListStr, fmt.Sprintf("P-%s", v))
		}
	}

	subjectID, _ := strconv.ParseInt(input.SubjectID, 10, 64)

	if banner, err = ctx.UploadFormFile("banner_img", fmt.Sprintf("recommend_subject_clothes_link_%d_banner", time.Now().Unix()), false); err != nil || banner == "" {
		ctx.BackError(10020, "banner图片不能为空")
		return
	}

	if h5link, err = ctx.Unzip("h5_zip"); err != nil {
		ctx.BackError(10020, "un_zip不能为空")
		return
	}
	if err, _ := dao.GetPlayerClothesExist(plist); err != nil {
		ctx.BackError(1010, err)
		return
	} else {
		data, _ := json.Marshal(idListStr)
		saveInput.SubjectID = subjectID
		saveInput.ClothesIDList = string(data)
		saveInput.BannerFileID = banner
		saveInput.HtmlZipFileID = appcfg.GetString("h5_link", "http://47.96.67.246:29970") + h5link
	}

	if err := dao.AddClothesWithLink(saveInput); err != nil {
		ctx.BackError(1009, err)
		return
	}
	return
}

type GetClothesLinkResp struct {
	Data  []RecommendItemGM `json:"data"`
	Count int               `json:"count"`
	Code  int               `json:"code"`
	Msg   string            `json:"msg"`
}

func GetClothesLink(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")

	var output GetClothesLinkResp
	if err, list := dao.GetClothesWithLink(); err != nil {
		ctx.BackError(1009, "参数有误")
		return
	} else {
		for _, v := range list {
			var temp RecommendItemGM
			temp.ID = v.ID
			temp.SubjectID = v.SubjectID
			temp.ClothesIDList = v.ClothesIDList
			temp.BannerFileID = common.DownImg(fmt.Sprintf("clothesLink_h5_%d", v.ID), v.BannerFileID)
			temp.HtmlZipFileID = v.HtmlZipFileID
			output.Data = append(output.Data, temp)
		}
		output.Count = len(list)
		ctx.BackSuccess(output)
		return
	}
	return
}

type DelClothesLinkReq struct {
	ID string `json:"id"`
}

func DelClothesLink(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var input DelClothesLinkReq
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(1009, err)
		return
	}
	id, _ := strconv.ParseInt(input.ID, 10, 64)

	if err := dao.DelClothesLink(id); err != nil {
		ctx.BackError(1009, err)
		return
	}
	ctx.BackSuccess("success")
}

func GetOneClothesLink(ctx *web.Context, idStr string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")

	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err, clothes := dao.GetOneClothesWithLink(id); err != nil {
		ctx.BackError(1009, "参数有误")
		return
	} else {
		ctx.BackSuccess(clothes)
		return
	}
	return
}

type AddTopicReq struct {
	SubjectID    string `json:"subject_id"`
	OClothesList string `json:"o_clothes_list"`
	PClothesList string `json:"p_clothes_list"`
}

//推荐专题
func AddTopic(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var inputSave RecommendItemsGM
	var input AddTopicReq
	var banner string
	var listSearch []string
	var h5link string
	var err error

	if err = ctx.GetFormData(&input); err != nil {
		ctx.BackError(1009, err)
		return
	}

	olist := strings.Split(input.OClothesList, ",")
	plist := strings.Split(input.PClothesList, ",")

	size := len(olist) + len(plist)

	subjectID, _ := strconv.ParseInt(input.SubjectID, 10, 64)

	if size < 1 {
		ctx.BackError(1010, "专题衣服为空")
		return
	}
	for _, v := range plist {
		listSearch = append(listSearch, fmt.Sprintf("P-%s", v))
	}
	if len(plist) > 1 {
		if err, clothesList := dao.GetPlayerClothesExist(plist); err != nil {
			ctx.BackError(1010, err)
			return
		} else {
			if len(clothesList) != len(plist) {
				ctx.BackError(1010, "推荐设计师衣服不存在")
				return
			}
		}
	}
	if len(olist) > 1 {
		for _, v := range olist {
			listSearch = append(listSearch, fmt.Sprintf("O-%s", v))
		}
		if err, list := dao.FindOffPaper(olist); err != nil {
			ctx.BackError(1010, err)
			return
		} else {
			if len(list) != len(olist) {
				ctx.BackError(1010, "推荐官方不存在")
				return
			}
		}
	}
	if banner, err = ctx.UploadFormFile("banner_img", fmt.Sprintf("recommend_subject_topic_%d_banner", time.Now().Unix()), false); err != nil || banner == "" {
		ctx.BackError(10021, "banner图片不能为空")
		return
	}
	if h5link, err = ctx.Unzip("h5_zip"); err != nil {
		ctx.BackError(10023, err)
		return
	}

	data, _ := json.Marshal(listSearch)
	inputSave.BannerFileID = banner
	inputSave.SubjectID = subjectID
	inputSave.ItemList = string(data)
	inputSave.HtmlZipFileID = appcfg.GetString("h5_link", "http://47.96.67.246:29970") + h5link
	inputSave.AddTime = time.Now().Unix()
	if err := dao.AddTopic(inputSave); err != nil {
		ctx.BackError(1011, "参数有误")
		return
	}
	ctx.BackSuccess("success")
	return
}

type GetTopicResp struct {
	Data  []TransRecommendItemGM `json:"data"`
	Count int                    `json:"count"`
	Code  int                    `json:"code"`
	Msg   string                 `json:"msg"`
}

func GetTopic(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")

	var output GetTopicResp
	if err, list := dao.GetTopic(); err != nil {
		ctx.BackError(1009, "参数有误")
		return
	} else {
		if err, subjects := dao.FindSubject(); err != nil {
			ctx.BackError(1009, "参数有误")
			return
		} else {
			output.Data = transTopic(list, subjects)
		}
		output.Count = len(list)
		ctx.BackSuccess(output)
		return
	}
	return
}

type DelTopicReq struct {
	TopicID string `json:"topic_id"`
}

func DelTopic(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var input DelTopicReq
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(1009, err)
		return
	}
	id, _ := strconv.ParseInt(input.TopicID, 10, 64)

	if err := dao.DelTopic(id); err != nil {
		ctx.BackError(1009, err)
		return
	}
	ctx.BackSuccess("success")
}

func transTopic(list []RecommendItemsGM, subjects []RecommendSubjectGM) (result []TransRecommendItemGM) {
	for _, v := range list {
		var temp TransRecommendItemGM
		temp.ID = v.ID
		for _, v1 := range subjects {
			if v.SubjectID == v1.ID {
				temp.SubjectName = v1.Name
			}
		}
		temp.ClothesIDList = v.ItemList
		temp.BannerImg = common.DownImg(fmt.Sprintf("%d", v.ID), v.BannerFileID)
		result = append(result, temp)
	}
	return
}

type AddClothesToForPoolReq struct {
	ClothesID     string `json:"clothes_id"`
	OwnerUsername string `json:"owner_username"`
}

//推荐池
func AddClothesToForPool(ctx *web.Context) {
	var saveInput RecommendPoolGM
	var input AddClothesToForPoolReq
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(1009, err)
		return
	}
	saveInput.OwnerUsername = input.OwnerUsername
	saveInput.ClothesID = input.ClothesID
	saveInput.AddTime = time.Now().Unix()
	saveInput.Status = constants.RECOMMEND_FOR_POOL
	if err := dao.AddClothesToForPool(saveInput); err != nil {
		ctx.BackError(1009, "参数有误")
		return
	}
	ctx.BackSuccess("success")
}

type AddClothesToPoolReq struct {
	ClothesID string `json:"clothes_id"`
}

func AddClothesToPool(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")

	var input AddClothesToForPoolReq
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(1009, err)
		return
	}
	if err := dao.UpdateClothesStatus(input.ClothesID); err != nil {
		ctx.BackError(1009, err)
		return
	}
	ctx.BackSuccess("success")
}

type GetClothesPoolResp struct {
	Data  []RecommendPoolGM `json:"data"`
	Count int               `json:"count"`
	Code  int               `json:"code"`
	Msg   string            `json:"msg"`
}

func GetClothesPool(ctx *web.Context) {
	var output GetClothesPoolResp
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if err, list := dao.GetClothesPool(); err != nil {
		ctx.BackError(1009, "参数有误")
		return
	} else {
		output.Data = list
		output.Count = len(list)
		ctx.BackSuccess(output)
		return
	}
	return
}

type DelClothesReq struct {
	ClothesID string `json:"clothes_id"`
}

func DelClothes(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var input DelClothesReq
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(1009, err)
		return
	}
	if err := dao.DelClothesOneFromPool(input.ClothesID); err != nil {
		ctx.BackError(1009, err)
		return
	}
	ctx.BackSuccess("success")
}

type GetSelectInfoResp struct {
	ClothesList []RecommendPoolGM
	TopicList   []RecommendItemsGM
	SuitList    []RecommendSubjectGM
}

func GetSelectInfo(ctx *web.Context) {
	var output GetSelectInfoResp
	var clothesList []RecommendPoolGM
	var topicList []RecommendItemsGM
	var suitList []RecommendSubjectGM
	var err error

	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if err, clothesList = dao.GetClothesPool(); err != nil {
		ctx.BackError(1015, "衣服有误")
		return
	}

	if err, topicList = dao.GetTopic(); err != nil {
		ctx.BackError(1016, "专题有误")
		return
	}

	if err, suitList = dao.FindSubject(); err != nil {
		ctx.BackError(1017, "套装有误")
		return
	}
	output.ClothesList = clothesList
	output.SuitList = suitList
	output.TopicList = topicList
	ctx.BackSuccess(output)
}
