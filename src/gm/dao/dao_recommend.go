package dao

import (
	"constants"
	. "model"
)

//杂志期刊主题
func AddSubject(rsg RecommendSubjectGM) (err error) {
	if _, err = GetEngine().Table("pixie_recommend_subject").Insert(&rsg); err != nil {
		return
	}
	return
}

func DelSubject(subjectID int64) (err error) {
	var rsg RecommendSubjectGM
	rsg.ID = subjectID
	if _, err = GetEngine().Table("pixie_recommend_subject").Where("id=?", subjectID).Delete(&rsg); err != nil {
		return
	}
	return
}

func UpdateSuit(suit, suitOwnerID string, subjectID int64) (err error) {
	var rsg RecommendSubjectGM
	rsg.SuitDetail = suit
	rsg.SuitOwnerID = suitOwnerID
	if _, err = GetEngine().Table("pixie_recommend_subject").Where("id=?", subjectID).Update(&rsg); err != nil {
		return
	}
	return
}

func OclohtesExist(list []interface{}) (err error, exist bool) {
	var res []GMOfficialPaper
	if err = GetEngine().Table("pixie_official_paper").In("id", list...).Find(&res); err != nil {
		return
	}
	if len(res) != len(list) {
		exist = false
	} else {
		exist = true
	}
	return
}

func PclohtesExist(list []interface{}) (err error, exist bool) {
	var res []GMRecommendClothes
	if err = GetEngine().Table("pixie_occupy").In("paper_id", list...).Find(&res); err != nil {
		return
	}
	if len(res) != len(list) {
		exist = false
	} else {
		exist = true
	}
	return
}

func GetSubject(id int64) (err error, rsg RecommendSubjectGM) {
	if _, err = GetEngine().Table("pixie_recommend_subject").Where("id=?", id).Get(&rsg); err != nil {
		return
	}
	return
}

func FindSubject() (err error, rsgs []RecommendSubjectGM) {
	if err = GetEngine().Table("pixie_recommend_subject").Find(&rsgs); err != nil {
		return
	}
	return
}

//杂志推荐跑马灯
func AddClothesWithLink(rig RecommendItemGM) (err error) {
	if _, err = GetEngine().Table("pixie_recommend_item").Insert(&rig); err != nil {
		return
	}
	return
}

func GetClothesWithLink() (err error, rigs []RecommendItemGM) {
	if err = GetEngine().Table("pixie_recommend_item").Find(&rigs); err != nil {
		return
	}
	return
}

func GetOneClothesWithLink(id int64) (err error, rig RecommendItemGM) {
	if _, err = GetEngine().Table("pixie_recommend_item").Where("id=?", id).Get(&rig); err != nil {
		return
	}
	return
}

func DelClothesLink(id int64) (err error) {
	var rig RecommendItemGM
	rig.ID = id
	if _, err = GetEngine().Table("pixie_recommend_item").Where("id=?", id).Delete(&rig); err != nil {
		return
	}
	return
}

//杂志主题+杂志推荐单品(单独推荐不附属任何池,不分页雪娇说了)
func AddTopic(rig RecommendItemsGM) (err error) {
	if _, err = GetEngine().Table("pixie_recommend_items").Insert(&rig); err != nil {
		return
	}
	return
}

func GetTopic() (err error, rigs []RecommendItemsGM) {
	if err = GetEngine().Table("pixie_recommend_items").Find(&rigs); err != nil {
		return
	}
	return
}

func DelTopic(id int64) (err error) {
	var rig RecommendItemsGM
	rig.ID = id
	if _, err = GetEngine().Table("pixie_recommend_items").Where("id=?", id).Delete(&rig); err != nil {
		return
	}
	return
}

//推荐池
func GetClothesPool() (err error, rpgs []RecommendPoolGM) {
	if err = GetEngine().Table("pixie_recommend_pool").Find(&rpgs); err != nil {
		return
	}
	return
}

//推荐单品池
func AddClothesToForPool(rgm RecommendPoolGM) (err error) {
	if _, err = GetEngine().Table("pixie_recommend_pool").Insert(&rgm); err != nil {
		return
	}
	return
}

func DelClothesOneFromPool(clothesID string) (err error) {
	var rgm RecommendPoolGM
	rgm.ClothesID = clothesID
	if _, err = GetEngine().Table("pixie_recommend_pool").Where("clothes_id=?", clothesID).Delete(&rgm); err != nil {
		return
	}
	return
}

func UpdateClothesStatus(paperId string) (err error) {
	var rgm RecommendPoolGM
	rgm.Status = constants.RECOMMEND_POOL
	rgm.ClothesID = paperId
	if _, err = GetEngine().Table("pixie_recommend_pool").Where("clothes_id=?", paperId).Update(&rgm); err != nil {
		return
	}
	return
}
