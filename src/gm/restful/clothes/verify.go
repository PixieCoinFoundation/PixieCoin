package clothes

import (
	"encoding/json"
	"fmt"
	"github.com/web"
	"gm/common"
	"gm/dao"
	. "model"
	. "pixie_contract/api_specification"
	"service/files"
	"service_pixie/paperInfo"
	"strconv"
	"tools"
)

type GetVerifyClothesListResp struct {
	Data  []DesignPaperGM `json:"data"`
	Count int             `json:"count"`
	Code  int             `json:"code"`
	Msg   string          `json:"msg"`
}

func GetVerifyClothesList(ctx *web.Context) {
	var output GetVerifyClothesListResp
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	if err, data := dao.GetVerifyList(PAPER_STATUS_ADMIN_QUEUE, true); err != nil {
		ctx.BackError(30000, err)
		return
	} else {
		output.Count = len(data)
		output.Data = data
		ctx.BackSuccess(output)
	}
	return
}

type GetVerifyBgListResp struct {
	Data  []DesignPaperGM `json:"data"`
	Count int             `json:"count"`
	Code  int             `json:"code"`
	Msg   string          `json:"msg"`
}

func GetVerifyBgList(ctx *web.Context) {
	var output GetVerifyBgListResp
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	if err, data := dao.GetVerifyList(PAPER_STATUS_ADMIN_QUEUE, false); err != nil {
		ctx.BackError(30010, err)
		return
	} else {
		output.Count = len(data)
		output.Data = data
		ctx.BackSuccess(output)
	}
	return
}

type GetVerifyOneResp struct {
	Data DesignPaperGMTrans `json:"data"`
}

func GetVerifyOne(ctx *web.Context, paperIDStr string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	var output GetVerifyOneResp

	paperID, err := strconv.ParseInt(paperIDStr, 10, 64)
	if paperID <= 0 || err != nil {
		ctx.BackError(30016, "参数不能为空")
		return
	}
	if err, data := dao.GetVerifyOne(paperID); err != nil {
		ctx.BackError(30017, err)
		return
	} else {
		output.Data = transFront(data)
		ctx.BackSuccess(output)
	}
	return
}

type GetInAuctionPaperResp struct {
	Data  []DesignPaperGM `json:"data"`
	Count int             `json:"count"`
	Code  int             `json:"code"`
	Msg   string          `json:"msg"`
}

func GetInAuctionPaper(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var output GetInAuctionPaperResp

	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	if err, data := dao.GetVerifyList(PAPER_STATUS_QUEUE, true); err != nil {
		ctx.BackError(30010, err)
		return
	} else {
		output.Count = len(data)
		output.Data = data
		ctx.BackSuccess(output)
	}
	return
}

type GetRejectPaperResp struct {
	Data  []DesignPaperGM `json:"data"`
	Count int             `json:"count"`
	Code  int             `json:"code"`
	Msg   string          `json:"msg"`
}

func GetRejectPaper(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var output GetRejectPaperResp

	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	if err, data := dao.GetVerifyList(PAPER_STATUS_FAIL, true); err != nil {
		ctx.BackError(30010, err)
		return
	} else {
		output.Count = len(data)
		output.Data = data
		ctx.BackSuccess(output)
	}
	return
}

type GetInAuctionBGPaperResp struct {
	Data  []DesignPaperGM `json:"data"`
	Count int             `json:"count"`
	Code  int             `json:"code"`
	Msg   string          `json:"msg"`
}

func GetInAuctionBGPaper(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var output GetInAuctionBGPaperResp

	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	if err, data := dao.GetVerifyList(PAPER_STATUS_QUEUE, false); err != nil {
		ctx.BackError(30010, err)
		return
	} else {
		output.Count = len(data)
		output.Data = data
		ctx.BackSuccess(output)
	}
	return
}

type GetRejectBGPaperResp struct {
	Data  []DesignPaperGM `json:"data"`
	Count int             `json:"count"`
	Code  int             `json:"code"`
	Msg   string          `json:"msg"`
}

func GetRejectBGPaper(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var output GetRejectBGPaperResp

	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	if err, data := dao.GetVerifyList(PAPER_STATUS_FAIL, false); err != nil {
		ctx.BackError(30010, err)
		return
	} else {
		output.Count = len(data)
		output.Data = data
		ctx.BackSuccess(output)
	}
	return
}

func transFront(dp DesignPaperGM) (res DesignPaperGMTrans) {
	res.PaperID = dp.PaperID
	res.AuthorUsername = dp.AuthorUsername
	cType, exist := common.ClothesTypeStringMap[dp.ClothesType]
	if exist {
		res.ClothesType = cType
	} else {
		res.ClothesType = "无"
	}
	var paperFile PaperFile
	var paperExtra PaperExtra
	json.Unmarshal([]byte(dp.File), &paperFile)
	json.Unmarshal([]byte(dp.Extra), &paperExtra)
	//部位
	res.PartType = common.PartTypeStringMap[dp.PartType]
	res.Desc = dp.Desc
	res.ClothesLayer = paperExtra.ZValue
	prefix := "./gm"
	iconPath := "/images/custom_" + fmt.Sprintf("%d", dp.PaperID) + "_icon_front"
	mainPath := "/images/custom_" + fmt.Sprintf("%d", dp.PaperID) + "_main_front"
	bottomPath := "/images/custom_" + fmt.Sprintf("%d", dp.PaperID) + "_bottom_front"
	collorPath := "/images/custom_" + fmt.Sprintf("%d", dp.PaperID) + "_collor_front"
	shadowPath := "/images/custom_" + fmt.Sprintf("%d", dp.PaperID) + "_shadow_front"
	frontPath := "/images/custom_" + fmt.Sprintf("%d", dp.PaperID) + "_bg_front"
	backPath := "/images/custom_" + fmt.Sprintf("%d", dp.PaperID) + "_bg_back_front"

	//todo exist img
	if !tools.ExistsFile(prefix + iconPath) {
		files.DecryptDownloadFile(paperFile.Icon, prefix+iconPath)
	}
	if !tools.ExistsFile(prefix + mainPath) {
		files.DecryptDownloadFile(paperFile.Main, prefix+mainPath)
	}
	if paperFile.Bottom != "" {
		if !tools.ExistsFile(prefix + bottomPath) {
			files.DecryptDownloadFile(paperFile.Bottom, prefix+bottomPath)
		}
		res.BottomLayer = bottomPath
	}
	if paperFile.Collar != "" {
		if !tools.ExistsFile(prefix + collorPath) {
			files.DecryptDownloadFile(paperFile.Collar, prefix+collorPath)
		}
		res.CollorImg = collorPath
	}
	if paperFile.Shadow != "" {
		if !tools.ExistsFile(prefix + shadowPath) {
			files.DecryptDownloadFile(paperFile.Shadow, prefix+shadowPath)
		}
		res.ShadowImg = shadowPath
	}
	if paperFile.Front != "" {
		if !tools.ExistsFile(prefix + frontPath) {
			files.DecryptDownloadFile(paperFile.Front, prefix+frontPath)
		}
		res.Front = frontPath
	}
	if paperFile.Back != "" {
		if !tools.ExistsFile(prefix + backPath) {
			files.DecryptDownloadFile(paperFile.Back, prefix+backPath)
		}
		res.Back = backPath
	}
	if dp.ClothesType == MODEL_MALE {
		res.ModelPic = "/images/c1000.png"
	} else if dp.ClothesType == MODEL_FEMALE {
		res.ModelPic = "/images/c1001.png"
	} else if dp.ClothesType == MODEL_KID {
		res.ModelPic = "/images/c1002.png"
	}
	res.IconImg = iconPath
	res.MainImg = mainPath
	res.IconIX = paperExtra.IconX
	res.IconIY = paperExtra.IconY
	return
}

type PassVerifyReq struct {
	PaperID   string `json:"paperID"`
	HiddenTag string `json:"hiddenTag"` //暗标签
}

func PassVerify(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	var input PassVerifyReq
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(30030, err)
		return
	}
	id, _ := strconv.ParseInt(input.PaperID, 10, 64)
	if input.PaperID == "" || id <= 0 {
		ctx.BackError(30031, "参数有误")
		return
	}
	if _, exist := paperInfo.GetSTagMap()[input.HiddenTag]; !exist && input.HiddenTag != "" {
		ctx.BackError(30032, "参数有误")
		return
	} else {
		if err := dao.PassVerify(id, input.HiddenTag); err != nil {
			ctx.BackError(30033, err)
			return
		}
	}
	ctx.BackSuccess("success")
}

type RejectVerifyReq struct {
	PaperID      string `json:"paperID"`
	RejectReason string `json:"rejectReason"` //拒绝标签
	Reason       string `json:"reason"`       //附加理由
}

func RejectVerify(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var input RejectVerifyReq
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(30036, err)
		return
	}
	id, _ := strconv.ParseInt(input.PaperID, 10, 64)
	if input.PaperID == "" || id <= 0 {
		ctx.BackError(30037, "参数有误")
		return
	}
	if input.RejectReason == "" && input.Reason == "" {
		ctx.BackError(30038, "参数有误")
		return
	}

	if err := dao.RejectVerify(id, input.RejectReason+input.Reason); err != nil {
		ctx.BackError(30039, err)
		return
	}
	ctx.BackSuccess("success")
}
