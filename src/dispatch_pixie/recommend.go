package dispatch_pixie

import (
	"constants"
	db "db_pixie/recommend"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service_pixie/recommend"
)

func GetSubjectHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input GetSubjectReq
	var output GetSubjectResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		output.SubjectInfo = append(output.SubjectInfo, recommend.GetSubject())
		result, parseErr = PixieMarshal(output)
	}
	return
}

func GetSubjectDetailHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input GetSubjectDetailReq
	var output GetSubjectDetailResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if input.Page > 1 {
			output.ItemsDown = recommend.GetDownPaper(input.IDList)
		} else {
			if err, topItem := db.GetItemList(input.SubjectID, 0, 0); err != nil {
				resp = PIXIE_ERR_GET_RECOMMEND_DETAIL
				return
			} else {
				output.ItemsTop = topItem
			}
			if err, exist := db.GetSubjectExistHistory(); err != nil {
				resp = PIXIE_ERR_GET_RECOMMEND_DETAIL
				return
			} else {
				output.HasHistory = exist
			}
			output.ItemsMedium = recommend.GetMediumPaper()
			output.ItemsDown = recommend.GetDownPaper([]string{})
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func GetHistorySubjectItemHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input GetHistorySubjectItemReq
	var page int64
	var output GetHistorySubjectItemResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if input.Page < 1 {
			page = 1
		} else {
			page = input.Page - 1
		}
		offset := page * constants.PIXIE_RECOMEMND_ITEM_PAGE_SIZE
		if err, topItem := db.GetItemList(0, offset, constants.PIXIE_RECOMEMND_ITEM_PAGE_SIZE); err != nil {
			resp = PIXIE_ERR_GET_RECOMMEND_DETAIL
			return
		} else {
			output.ItemsTop = topItem
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func GetSubjectSuitHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input GetSubjectSuitReq
	var output GetSubjectSuitResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		output.SuitList = recommend.GetSubjectSuit()
		result, parseErr = PixieMarshal(output)
	}
	return
}

func GetSubjectTopicHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input GetSubjectTopicReq
	var output GetSubjectTopicResp
	var page int64
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if input.Page < 1 {
			page = 0
		} else {
			page = input.Page - 1
		}
		offset := page * constants.RECOMMEND_TOPIC_PAGESIZE
		if err, topicList := db.GetTopic(offset, constants.RECOMMEND_TOPIC_PAGESIZE); err != nil {
			resp = PIXIE_ERR_GET_RECOMMEND_TOPIC
			return
		} else {
			output.Topic = topicList
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}
