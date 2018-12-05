package admin

import (
	"encoding/json"
	"github.com/web"
	"gm/common"
	"gm/dao"
	. "model"
	"strconv"
	"strings"
)

func SaveAPIToDB(apiMap map[string]string) {
	// dao.ClearApi()
	// dao.AddApiList(apiMap)
}

var Menus []MenuInfo

func init() {
	LoadMenu()
}

func LoadMenu() {
	if err := json.Unmarshal([]byte(common.MenuList), &Menus); err != nil {
		panic(err)
	}
}

func GetRootMenu() []MenuInfo {
	return Menus
}

func GetAdminMenu(idList []int64) (list []MenuInfo) {
	if len(idList) < 1 {
		return
	}
	for _, men := range Menus {
		var temp MenuInfo
		temp.ID = men.ID
		temp.Title = men.Title
		temp.Icon = men.Icon
		temp.ParentID = men.ParentID
		temp.Spread = men.Spread
		temp.Href = men.Href
		for _, child := range men.Children {
			for _, v := range idList {
				if child.ID == v {
					temp.Children = append(temp.Children, child)
				}
			}
		}
		if len(temp.Children) > 0 {
			list = append(list, temp)
		}
	}
	return
}

type Nodes struct {
	ID      int64  `json:"id"`
	PID     int64  `json:"pId"`
	Name    string `json:"name"`
	Checked bool   `json:"checked"`
	Open    bool   `json:"open"`
}

func GetZnodes(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var result []Nodes
	for _, v := range GetRootMenu() {

		var temp Nodes
		temp.ID = v.ID
		temp.PID = v.ParentID
		temp.Name = v.Title
		temp.Checked = true

		for _, v1 := range v.Children {
			var temp1 Nodes
			temp1.ID = v1.ID
			temp1.PID = v1.ParentID
			temp1.Name = v1.Title
			temp1.Checked = true
			result = append(result, temp1)
		}
		result = append(result, temp)
	}
	ctx.BackSuccess(result)
	return
}

type AddApiReq struct {
	Name      string `json:"name"`
	MenuLayer string `json:"menuLayer"`
	MenuUrl   string `json:"menuUrl"`
	MenuType  string `json:"menuType"`
	ParentID  string `json:"parentID"`
}

type AddRoleReq struct {
	RoleName string `json:"roleName"`
	Desc     string `json:"roleDesc"` //角色描述
}

func AddRole(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	var input AddRoleReq
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(2008, err)
		return
	}

	if _, err := dao.AddRole(input.RoleName, input.Desc); err != nil {
		ctx.BackError(2009, err)
		return
	}
	ctx.BackSuccess("success")
	return

}

type GetRoleListResp struct {
	Data  []AdminRole `json:"data"`
	Count int         `json:"count"`
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
}

func GetRoleList(ctx *web.Context) {
	var output GetRoleListResp
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	if list, err := dao.GetRoleList(); err != nil {
		ctx.BackError(2010, err)
		return
	} else {
		output.Data = list
		output.Count = len(list)
		ctx.BackSuccess(output)
		return
	}
}

type AddRoleApiPrivReq struct {
	RoleID string `json:"RoleID"`
	IDList string `json:"IDList"`
}

func AddRoleApiPriv(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	var input AddRoleApiPrivReq
	var apiList []int64
	var err error
	var exist bool

	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(2013, err)
		return
	}
	roleID, _ := strconv.ParseInt(input.RoleID, 10, 64)
	apiListStr := strings.Split(input.IDList, ",")
	if len(apiListStr) < 1 {
		ctx.BackError(2008, "没有权限要增加")
		return
	}
	for _, v := range apiListStr {
		apiID, _ := strconv.ParseInt(v, 10, 64)
		apiList = append(apiList, apiID)
	}

	if exist, err, _ = dao.GetMenuRolePriv(roleID); err != nil {
		ctx.BackError(2009, err)
		return
	}
	if err := dao.AddRoleApiRriv(exist, roleID, apiList); err != nil {
		ctx.BackError(2014, err)
		return
	}
	ctx.BackSuccess("success")
}
