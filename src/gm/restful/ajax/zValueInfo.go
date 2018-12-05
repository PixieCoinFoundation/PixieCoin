package ajax

import (
	"appcfg"
	"constants"
	"github.com/web"
	. "pixie_contract/api_specification"
	"service_pixie/paperInfo"
	"strconv"
	"strings"
)

var zValueList = []PartOfZvalue{}
var partMap = make(map[string]string)
var pixiePaperPartTypeMap = make(map[string]string)

func init() {
	if appcfg.GetServerType() != "" && appcfg.GetServerType() != constants.SERVER_TYPE_GM {
		return
	}
	for k, v := range PAPER_TYPE_ZVALUE_MAP {
		var temp PartOfZvalue
		if len(v) > 1 {
			temp.Part = k
			for _, val := range v {
				var zval ZvalueWithName
				zval.Name = ZValueMap[val]
				zval.Zvalue = strconv.Itoa(val)
				temp.PartofList = append(temp.PartofList, zval)
			}
			zValueList = append(zValueList, temp)
		}
	}
	for k, v := range PixiePaperPartTypeMap {
		if !strings.EqualFold(k, PAPER_TYPE_BEIJING) {
			pixiePaperPartTypeMap[k] = v
		}
		if k == PAPER_TYPE_HAND || k == PAPER_TYPE_WAZI || k == PAPER_TYPE_MAKEUP || PAPER_TYPE_BEIJING == k {
			continue
		}
		partMap[k] = v
	}

}

type PaperInfoResp struct {
	TagList     map[string]string
	STagList    map[string]string
	StyleList   map[string]string
	OriginPart  map[string]string
	PartNameMap map[string]string
	ZvalueList  []PartOfZvalue
}

func GetPaperTagZvalueInfo(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var output PaperInfoResp
	output.ZvalueList = zValueList
	styleMap := make(map[string]string)

	output.PartNameMap = partMap
	output.STagList = paperInfo.GetSTagMap()
	for k, v := range paperInfo.GetStyleMap() {
		styleMap[k[:strings.Index(k, `_`)]] = v
	}
	output.OriginPart = PixiePaperPartTypeMap
	output.StyleList = styleMap
	output.TagList = paperInfo.GetTagMap()
	ctx.BackSuccess(output)
	return
}
