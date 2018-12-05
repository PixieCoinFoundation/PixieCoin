package server

import (
	"github.com/web"
	"gm/dao"
	. "model"
	"time"
)

type GetMaintenanceResp struct {
	MainTainAgo   []Maintenance
	MainTainIng   []Maintenance
	MainTainAfter []Maintenance
}

func GetMaintenanceInfo(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	tNow := time.Now().Unix()
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}

	var output GetMaintenanceResp
	if err, res := dao.GetMaintenanceList(); err != nil {
		ctx.BackError(20001, err)
		return
	} else {
		for _, v := range res {
			if v.StartTime > tNow && v.EndTime > tNow {
				output.MainTainAgo = append(output.MainTainAgo, v)
			} else if v.StartTime < tNow && v.EndTime > tNow {
				output.MainTainIng = append(output.MainTainIng, v)
			} else if v.StartTime < tNow && v.EndTime < tNow {
				output.MainTainAfter = append(output.MainTainAfter, v)
			} else {
				ctx.BackError(20003, "活动有误,请联系管理员")
				return
			}
		}
		ctx.BackSuccess(output)
	}
	return
}

type AddMaintenanceReq struct {
	Content   string `json:"textarea"`
	StartTime string
	EndTime   string
	URL       string `json:"url"`
	ShowTime  string
}

func AddMaintenance(ctx *web.Context) {
	var input AddMaintenanceReq
	var mainMod Maintenance
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(20006, err)
		return
	}
	if input.StartTime == "" || input.EndTime == "" {
		ctx.BackError(20007, "维护时间不能为空")
		return
	}
	tStart, _ := time.ParseInLocation("2006-01-02 15:04:05", input.StartTime, time.Local)
	tEnd, _ := time.ParseInLocation("2006-01-02 15:04:05", input.StartTime, time.Local)

	if input.ShowTime != "" {
		tShow, _ := time.ParseInLocation("2006-01-02 15:04:05", input.ShowTime, time.Local)
		mainMod.ShowTime = tShow.Unix()
	}
	mainMod.Content = input.Content
	mainMod.StartTime = tStart.Unix()
	mainMod.EndTime = tEnd.Unix()
	if input.URL != "" {
		mainMod.URL = input.URL
	}
	if _, err := dao.AddMaintenance(mainMod); err != nil {
		ctx.BackError(20008, err)
		return
	}
	return
}
