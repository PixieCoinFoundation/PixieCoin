package admin

import (
	"encoding/json"
	"fmt"
	"github.com/web"
	"gm/common"
	"gm/dao"
	. "model"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tools"
)

type AdminListResp struct {
	Data  []AdminGM `json:"data"`
	Count int       `json:"count"`
	Code  int       `json:"code"`
	Msg   string    `json:"msg"`
}

func AdminList(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	var output AdminListResp

	if err, list := dao.GetAdminList(); err != nil {
		ctx.BackError(1001, err)
		return
	} else {
		output.Data = list
		output.Count = len(list)
		ctx.BackSuccess(output)
	}
	return
}

func DeleteAdmin(ctx *web.Context, idStr string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	var id int64
	var err error

	if id, err = strconv.ParseInt(idStr, 10, 64); err != nil {
		ctx.BackError(1002, "id有误")
		return
	}
	if err = dao.DeleteAdmin(id); err != nil {
		ctx.BackError(1003, err)
		return
	}
	ctx.BackSuccess("success")
	return
}

type AddAdminReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   string `json:"roleID"`
}

func AddAdmin(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	var input AddAdminReq

	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(1004, err)
		return
	}
	role, err := strconv.ParseInt(input.RoleID, 10, 64)

	if strings.TrimSpace(input.Username) == "" || strings.TrimSpace(input.Password) == "" || err != nil {
		ctx.BackError(1006, "参数有误")
		return
	}

	var agm AdminGM
	agm.Username = input.Username
	agm.Password = input.Password
	agm.Role = role
	agm.LastLoginTime = time.Now().Format("20060102 150405")

	if err := dao.AddAdmin(agm); err != nil {
		ctx.BackError(1007, err)
		return
	}
	ctx.BackSuccess("success")
	return
}

type ChangeRoleReq struct {
	RoleID string `json:"RoleID"`
}

func ChangeRole(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var exist bool
	var id int64

	if exist, id = ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	var input ChangeRoleReq
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(1009, err)
		return
	}
	roleID, _ := strconv.ParseInt(input.RoleID, 10, 64)

	if id != 1 {
		ctx.BackError(1001, "权限不足")
		return
	}
	if err := dao.ChangeRole(id, roleID); err != nil {
		ctx.BackError(1001, err)
		return
	}
	ctx.BackSuccess("success")
	return
}

func GetNavs(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var exist bool
	var id int64
	if exist, id = ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	if err, data := dao.GetOneAdmin(id); err != nil {
		ctx.BackError(1008, err)
		return
	} else {
		if data.Role == 1 {
			data, _ := json.Marshal(GetRootMenu())
			ctx.BackSuccess(string(data))
			return
		} else {
			if _, err, res := dao.GetMenuRolePriv(data.Role); err != nil {
				ctx.BackError(1008, err)
				return
			} else {
				idlist := []int64{}
				json.Unmarshal([]byte(res.MenuIDList), &idlist)
				data, _ := json.Marshal(GetAdminMenu(idlist))
				ctx.BackSuccess(string(data))
			}
		}
	}
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	var input LoginReq
	if err := ctx.GetFormData(&input); err != nil {
		ctx.BackError(1009, err)
		return
	}
	if input.Username == "" || input.Password == "" {
		ctx.BackError(1010, "用户名或者密码不能为空")
		return
	}
	var id int64
	if exist, admin := dao.AdminExist(input.Username, input.Password); !exist {
		ctx.BackError(1011, "用户名或者密码不匹配")
		return
	} else {
		id = admin.ID
	}

	ctx.SetCookie(&http.Cookie{
		Name:   fmt.Sprintf("%s-%d", input.Username, id),
		Value:  tools.Md5([]byte(input.Username + common.GetTokenDate() + "otakugame")),
		MaxAge: 86400,
		Path:   "/",
	})
	ctx.BackSuccess("success")
	return
}
