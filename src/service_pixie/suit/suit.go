package service_pixie

import (
	"common"
	"db_pixie/paper"
	"encoding/json"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	sp "service_pixie/paper"
	"tools"
)

func GetSuitMD5(suitModel Suit) (err error, md5Str string) {
	suitModel.SuitImg = ""
	suitModel.SuitName = ""
	data, _ := json.Marshal(suitModel)
	if md5Str, err = common.GenMD5Raw(data); err != nil {
		return
	}
	return
}

func CheckSuit(suitModel *Suit, expectModel string, player *GFPlayer) (resp PixieRespInfo, ok bool) {
	idPartMap := suitModel.GenIDPartMap()

	//检查套装内部件数量
	if len(idPartMap) < 2 {
		resp = PIXIE_ERR_SUIT_NOT_LEGAL
		return
	}

	//检查是否拥有
	for id, _ := range idPartMap {
		if !player.HaveClothes(id) {
			resp = PIXIE_ERR_NOT_HAVE_CLOTHES
			return
		}
	}

	//检查套装内容是否合法
	if _, legal := suitPartAndModelLegal(suitModel, expectModel, idPartMap); !legal {
		resp = PIXIE_ERR_SUIT_NOT_LEGAL
		return
	}

	ok = true
	return
}

//只用于检查套装 不能用于检查过关的服装搭配
func suitPartAndModelLegal(suitModel *Suit, expectModel string, idPartMap map[string]string) (err error, legal bool) {
	playerPaperIDList, officialPaperIDList := listPaperID(idPartMap)

	//检查玩家上传衣服
	var designList []DesignPaper

	if len(playerPaperIDList) > 0 {
		if err, designList = paper.GetDesignPaperInIDList(playerPaperIDList); err != nil {
			return
		}

		for _, d := range designList {
			//检查模特和部位是否正确
			nowPart := idPartMap[tools.GenPixiePlayerClothesID(d.PaperID)]
			if nowPart == "" || nowPart != d.PartType || d.ClothesType != expectModel {
				return
			}
		}
	}

	//检查服装数量
	if len(designList)+len(officialPaperIDList) != len(idPartMap) {
		return
	}

	//检查官方服装
	for _, cid := range officialPaperIDList {
		if clo, exist := sp.GetOneOfficialPaperByCID(cid); exist {
			//检查模特和部位是否正确
			nowPart := idPartMap[cid]
			if nowPart == "" || nowPart != clo.PartType || clo.ClothesType != expectModel {
				return
			}
		} else {
			return
		}
	}

	legal = true
	return
}

func listPaperID(idPartMap map[string]string) (playerPaperIDs []interface{}, officialPaperIDs []string) {
	playerPaperIDs = make([]interface{}, 0)
	officialPaperIDs = make([]string, 0)

	if len(idPartMap) > 0 {
		for clothesID, _ := range idPartMap {
			if src, id := tools.GetPixieClothesDetail(clothesID); src == string(CLOTHES_ORIGIN_PLAYER_PAPER) {
				playerPaperIDs = append(playerPaperIDs, id)
			} else if src == string(CLOTHES_ORIGIN_OFFICIAL_PAPER) {
				officialPaperIDs = append(officialPaperIDs, clothesID)
			}
		}
	}

	return
}
