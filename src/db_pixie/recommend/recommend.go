package recommend

import (
	"constants"
	"dao_pixie"
	"database/sql"
	"db_pixie/paper"
	"encoding/json"
	"fmt"
	. "logger"
	. "pixie_contract/api_specification"
	spaper "service_pixie/paper"
	"tools"
)

func GetSubject() (err error, subject Subject, list []RecommendClothes) {
	var suitStr, spotStr string
	var idList []string

	if err = dao_pixie.GetSubjectStmt.QueryRow().Scan(&subject.ID, &subject.Name, &suitStr, &subject.SuitOwner, &spotStr, &subject.OpenTime, &subject.BannerImg); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("no subject in db")
		} else {
			Err(err)
		}

		return
	}
	json.Unmarshal([]byte(spotStr), &subject.HotSpots)
	json.Unmarshal([]byte(suitStr), &idList)
	if suitStr != "" {
		argList := make([]interface{}, 0)
		for _, v := range idList {
			src, id := tools.GetPixieClothesDetail(v)
			//获取官方衣服
			if src == string(CLOTHES_ORIGIN_OFFICIAL_PAPER) {
				if oPaper, exist := spaper.GetOneOfficialPaperByCID(v); !exist {
					Err("GetSubject not exist", id)
					continue
				} else {
					list = append(list, offPaperTransToRecommendPaper(oPaper))
				}
			} else {
				argList = append(argList, id)
			}
		}
		for len(argList) < constants.PIXIE_RECOMMEND_IN_COUNT_SIZE {
			argList = append(argList, 0)
		}
		var res []RecommendClothes
		if err, res = paper.DesignPaperInList(argList); err != nil {
			Err(err)
			return
		} else {
			list = append(list, getMinPriceToRecomemnd(res)...)
		}
	}
	return
}

func GetSubjectExistHistory() (err error, exist bool) {
	var rows *sql.Rows
	if rows, err = dao_pixie.Get10SubjectStmt.Query(); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		list := []Subject{}
		for rows.Next() {
			var subject Subject
			var suitStr, spotStr string
			if err = rows.Scan(&subject.ID, &subject.Name, &suitStr, &subject.SuitOwner, &spotStr, &subject.OpenTime, &subject.BannerImg); err != nil {
				Err(err)
				return
			}
			list = append(list, subject)
		}
		if len(list) > 1 {
			exist = true
		}
	}
	return
}

func GetItemList(subjectID, offset, pageSize int64) (err error, list []*TopPapers) {
	var rows *sql.Rows
	if subjectID > 0 {
		rows, err = dao_pixie.GetItemBySubjectIDStmt.Query(subjectID)
	} else {
		rows, err = dao_pixie.GetItemByPageStmt.Query(offset, pageSize)
	}
	if err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var topPaper TopPapers
			var temp SubjectItemsTop
			var clothesIDList []int64
			var idList []string

			if err = rows.Scan(&temp.ID, &temp.SubjectID, &temp.ClothesIDList, &temp.BannerImg, &temp.HtmlLink); err != nil {
				Err(err)
				return
			}
			topPaper.SubjectID = temp.SubjectID
			topPaper.BannerImg = temp.BannerImg
			topPaper.HtmlLink = temp.HtmlLink
			json.Unmarshal([]byte(temp.ClothesIDList), &idList)
			for _, v := range idList {
				src, id := tools.GetPixieClothesDetail(v)
				//获取官方衣服
				if src == string(CLOTHES_ORIGIN_OFFICIAL_PAPER) {
					if oPaper, exist := spaper.GetOneOfficialPaperByCID(v); !exist {
						Err("GetSubject not exist", id)
						continue
					} else {
						tempPaper := offPaperTransToRecommendPaper(oPaper)
						tempPaper.ClothesID = v
						topPaper.Papers = append(topPaper.Papers, tempPaper)
					}
				} else {
					clothesIDList = append(clothesIDList, id)
				}
			}
			var res []RecommendClothes
			if err, res = listRecommendByIDList(clothesIDList); err != nil {
				Err(err)
				return
			}
			configNameMap := make(map[int64]string, 0)
			_, recommend := getPaperToRecommend(true, res, nil, configNameMap)
			topPaper.Papers = append(topPaper.Papers, recommend...)
			list = append(list, &topPaper)
		}
	}
	return
}

func GetRecommendPool(status int) (err error, list []*ItemsPaper) {
	var rows *sql.Rows
	var PList []int64
	var res []RecommendClothes

	clothesInfoMap := map[string]SubjectPool{}
	clothesNameMap := make(map[int64]string, 0)
	if rows, err = dao_pixie.GetClothesPoolStmt.Query(1); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var temp SubjectPool
			if err = rows.Scan(&temp.ID, &temp.ClothesID, &temp.OwnerUsername, &temp.BannerImg); err != nil {
				Err(err)
				continue
			}

			src, id := tools.GetPixieClothesDetail(temp.ClothesID)

			//获取官方衣服
			if src == string(CLOTHES_ORIGIN_OFFICIAL_PAPER) {
				if oPaper, exist := spaper.GetOneOfficialPaperByCID(temp.ClothesID); !exist {
					Err("GetOneOfficialPaperByCID not exist", id)
					return
				} else {
					var item ItemsPaper
					item.RecommendClothes = offPaperTransToRecommendPaper(oPaper)
					item.BannerImg = temp.BannerImg
					item.RecommendClothes.ClothesID = temp.ClothesID
					list = append(list, &item)
				}
			} else {
				if temp.OwnerUsername != "" {
					clothesNameMap[id] = temp.OwnerUsername
				}
				PList = append(PList, id)
			}
			clothesInfoMap[temp.ClothesID] = temp
		}
	}
	//设计师衣服
	if err, res = listRecommendByIDList(PList); err != nil {
		Err(err)
		return
	} else {
		items, _ := getPaperToRecommend(true, res, clothesInfoMap, clothesNameMap)
		list = append(list, items...)
	}
	return
}

func GetTopic(page int64, pageSize int64) (err error, topicList []SubjectTopic) {
	var rows *sql.Rows
	if rows, err = dao_pixie.GetTopicStmt.Query(page, pageSize); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var topic SubjectTopic
			var paramsIdList []int64
			var list []RecommendClothes
			var itemList string
			var idList []string

			if err = rows.Scan(&topic.ID, &topic.SubjectID, &itemList, &topic.BannerImg, &topic.HtmlLink, &topic.AddTime); err != nil {
				Err(err)
				return
			}
			if itemList != "" {
				json.Unmarshal([]byte(itemList), &idList)
				for _, v := range idList {
					src, id := tools.GetPixieClothesDetail(v)
					//获取官方衣服
					if src == string(CLOTHES_ORIGIN_OFFICIAL_PAPER) {
						if oPaper, exist := spaper.GetOneOfficialPaperByCID(v); !exist {
							Err("GetOneOfficialPaperByCID not exist", id)
							return
						} else {
							topic.PaperList = append(topic.PaperList, offPaperTransToRecommendPaper(oPaper))
						}
					} else {
						paramsIdList = append(paramsIdList, id)
					}
				}
			}
			if err, list = listRecommendByIDList(paramsIdList); err != nil {
				Err(err)
				return
			} else {
				configNameMap := make(map[int64]string, 0)
				_, recommend := getPaperToRecommend(true, list, nil, configNameMap)
				topic.PaperList = append(topic.PaperList, recommend...)
			}
			topicList = append(topicList, topic)
		}
	}
	return
}

func listRecommendByIDList(idList []int64) (error, []RecommendClothes) {
	argList := make([]interface{}, 0)
	var list []RecommendClothes
	size := len(idList)
	for index, v := range idList {
		argList = append(argList, v)
		if len(argList) >= 100 {
			if err, recommendList := paper.GetRecommendDesignPaperInIDList(argList); err != nil {
				Err(err)
				return err, recommendList
			} else {
				list = append(list, recommendList...)
			}
			argList = make([]interface{}, 0)
		}
		if (size - 1) == index {
			for len(argList) < constants.PIXIE_LAND_MAX_SALE_SIZE {
				argList = append(argList, 0)
			}
			if err, recommendList := paper.GetRecommendDesignPaperInIDList(argList); err != nil {
				Err(err)
				return err, recommendList
			} else {
				list = append(list, recommendList...)
			}
		}
	}
	return nil, list
}

func offPaperTransToRecommendPaper(temp *OfficialPaper) (v RecommendClothes) {
	v.PaperID = temp.PaperID
	v.AuthorUsername = temp.AuthorUsername
	v.AuthorNickname = temp.AuthorNickname
	v.AuthorHead = temp.AuthorHead
	v.AuthorSex = temp.AuthorSex
	v.OwnerUsername = temp.OwnerUsername
	v.OwnerNickname = temp.OwnerNickname
	v.OwnerHead = temp.OwnerHead
	v.OwnerSex = temp.OwnerSex

	v.ClothesType = temp.ClothesType
	v.PartType = temp.PartType
	v.Cname = temp.Cname
	v.Desc = temp.Desc
	v.File = temp.File
	v.Extra = temp.Extra
	v.Status = temp.Status
	v.Star = temp.Star
	//普通标签
	v.Tag1 = temp.Tag1
	v.Tag2 = temp.Tag2

	v.Style = temp.Style
	//暗标签
	v.STag = temp.STag
	v.Price = float64(temp.Price)
	v.PriceType = temp.PriceType
	return
}

func getRecommendClothes(clothesIDList []interface{}, haveBanner bool, clothesInfoMap map[string]SubjectItemsTop) (err error, list []*ItemsPaper) {
	var recommendPaperList []RecommendClothes
	if err, recommendPaperList = paper.GetRecommendDesignPaperInIDList(clothesIDList); err != nil {
		Err(err)
		return
	} else {
		for _, v := range recommendPaperList {
			var item ItemsPaper
			item.RecommendClothes = v
			if haveBanner {
				item.BannerImg = clothesInfoMap[fmt.Sprintf("%d", v.PaperID)].BannerImg
			}
			list = append(list, &item)
		}
	}
	return
}

func getPaperToRecommend(isRecommend bool, source []RecommendClothes, clothesInfoMap map[string]SubjectPool, configNameMap map[int64]string) (res []*ItemsPaper, recList []RecommendClothes) {
	minClothesMap := make(map[int64]RecommendClothes, 0)
	for _, temp := range source {
		if val, exist := minClothesMap[temp.PaperID]; exist {
			if getMinPrice(val.Price, val.PriceType) > getMinPrice(temp.Price, temp.PriceType) {
				minClothesMap[temp.PaperID] = temp
			}
		} else {
			minClothesMap[temp.PaperID] = temp
		}
	}
	for _, v := range source {
		var item ItemsPaper
		if isRecommend {
			item.RecommendClothes = v
			item.BannerImg = clothesInfoMap[v.ClothesID].BannerImg
		}
		//如果配置了OwnerUsername直接显示推荐
		if val, exist := configNameMap[v.PaperID]; exist {
			if val == v.OwnerUsername {
				if isRecommend {
					res = append(res, &item)
				}
				recList = append(recList, v)
			}
			continue
		}

		if minClothesMap[v.PaperID].OwnerUsername == v.OwnerUsername {
			if isRecommend {
				res = append(res, &item)
			}
			recList = append(recList, v)
		}
	}
	return
}

func getMinPriceToRecomemnd(source []RecommendClothes) (res []RecommendClothes) {
	for _, temp := range source {
		minClothesMap := make(map[int64]RecommendClothes, 0)
		if val, exist := minClothesMap[temp.PaperID]; exist {
			if getMinPrice(val.Price, val.PriceType) > getMinPrice(temp.Price, temp.PriceType) {
				minClothesMap[temp.PaperID] = temp
			}
		} else {
			minClothesMap[temp.PaperID] = temp
		}
	}
	return
}

func getMinPrice(price float64, priceType int) (priceRes float64) {
	if priceType == constants.PIXIE_PXC_TYPE {
		priceRes = price * constants.PXCTOGOLD
		return
	}
	priceRes = price
	return
}
