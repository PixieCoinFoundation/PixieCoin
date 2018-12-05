package player

import (
	"encoding/json"
	"github.com/web"
	"gm/dao"
	. "model"
	. "pixie_contract/api_specification"
	"strings"
	"time"
)

type GetPlayerInfoResp struct {
	Data StatusGM
}

func GetPlayerInfo(ctx *web.Context, username string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	if username == "" {
		ctx.BackError(1001, "查询条件不能为空")
		return
	}
	if err, res := dao.GetPlayerByUsername(username); err != nil {
		ctx.BackError(1001, err)
		return
	} else {
		ctx.BackSuccess(transPlayerInfo(res))
	}
	return
}

func transPlayerInfo(status Status) (statusGM StatusGM) {
	statusGM.Uid = status.Uid
	statusGM.Username = status.Username
	statusGM.Nickname = status.Nickname
	statusGM.StreetName = status.StreetName
	statusGM.Head = status.Head
	statusGM.Sex = status.Sex
	statusGM.Level = status.Level
	statusGM.Exp = status.Exp
	statusGM.Money = status.Money
	statusGM.Tili = status.Tili
	statusGM.MaxTili = status.MaxTili
	//需要处理的字段
	statusGM.TiliRTime = time.Unix(status.TiliRTime, 0).Format("2006-01-02 15:04:05")
	for k, v := range status.Clothes {
		var temp ClothesInfoGM
		temp.ClothesID = k[2:]
		temp.Count = v.Count
		temp.Time = time.Unix(v.Time, 0).Format("2006-01-02 15:04:05")
		if strings.HasPrefix(k, "O-") {
			statusGM.OffClothesList = append(statusGM.OffClothesList, temp)
		} else {
			statusGM.ClothesList = append(statusGM.ClothesList, temp)
		}
	}
	statusGM.HomeShow = status.HomeShow

	for k, v := range status.RecordMap {
		var temp PixieRecordGM
		temp.ID = k
		temp.Score = v.Score
		temp.BestScore = v.BestScore
		temp.Rank = v.Rank
		temp.Clothes = v.Clothes
		temp.BestClothes = v.BestClothes
		hr, _ := json.Marshal(v.HistoryRank)
		temp.HistoryRank = string(hr)
		statusGM.RecordList = append(statusGM.RecordList, temp)
	}

	for k, v := range status.NpcFriendship {
		var temp FriendshipDetailGM
		temp.NpcID = k
		temp.Level = v.Level
		temp.Exp = v.Exp
		statusGM.NpcFriendship = append(statusGM.NpcFriendship, temp)
	}

	statusGM.BanStartTime = status.BanInfo.BanStartTime
	statusGM.BanEndTime = status.BanInfo.BanEndTime
	statusGM.BanReason = status.BanInfo.BanReason

	npc1, _ := json.Marshal(status.NpcBuff)
	statusGM.NpcBuff = string(npc1)
	dayExt, _ := json.Marshal(status.DayExtraInfo)
	statusGM.DayExtraInfo = string(dayExt)
	ex1, _ := json.Marshal(status.ExtraInfo1)
	statusGM.ExtraInfo1 = string(ex1)
	so, _ := json.Marshal(status.StreetOperate)
	statusGM.StreetOperate = string(so)
	sm, _ := json.Marshal(status.SuitMap)
	statusGM.SuitMap = string(sm)
	return
}
