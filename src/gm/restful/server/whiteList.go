package server

import (
	"constants"
	"encoding/json"
	"github.com/web"
	"gm/dao"
	. "model"
	"strconv"
)

type TableParseList struct {
	ID string
}

type GetAccountListResp struct {
	WhiteList  []TableParseList
	BlackList  []TableParseList
	ServerList []TableParseList
	AllowPatch bool
}

func GetAccountList(ctx *web.Context) {
	if exist, _ := ctx.CheckToken(); !exist {
		ctx.BackError(1000, "没有权限,请重新登录")
		return
	}
	var output GetAccountListResp
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if err, whiteList := getConfigKeyToArray(constants.MAINTAIN_WHITE_LIST_KEY); err != nil {
		ctx.BackError(30001, "获取白名单出错")
		return
	} else {
		output.WhiteList = whiteList
	}
	if err, blackList := getConfigKeyToArray(constants.BLACK_ACCOUNT_LIST_KEY); err != nil {
		ctx.BackError(30002, "获取黑名单出错")
		return
	} else {
		output.BlackList = blackList
	}
	if err, serverList := getConfigKeyToArray(constants.BLACK_GAME_SERVER_LIST_KEY); err != nil {
		ctx.BackError(30003, "获取服务器名单出错")
		return
	} else {
		output.ServerList = serverList
	}

	if err, val := getConfigKeyToBool(constants.ALLOW_PATCH_CONFIG_KEY); err != nil {
		ctx.BackError(30004, "允许热更新")
		return
	} else {
		output.AllowPatch = val
	}

	ctx.BackSuccess(output)

	return
}

func getConfigKeyToArray(keyName string) (err error, list []TableParseList) {
	var data ConfigInfo
	if _, data, err = dao.GetValByKeyName(keyName); err != nil {
		return err, list
	} else {
		var accountList []string
		json.Unmarshal([]byte(data.Value), &accountList)
		for _, v := range accountList {
			list = append(list, TableParseList{ID: v})
		}
		return nil, list
	}
}

func getConfigKeyToBool(keyName string) (err error, val bool) {
	var data ConfigInfo
	if _, data, err = dao.GetValByKeyName(keyName); err != nil {
		return err, false
	} else {
		val, _ = strconv.ParseBool(data.Value)
	}
	return
}

func addConfigKeyBool(keyName string) (err error) {
	var data ConfigInfo
	var existKey bool
	if existKey, data, err = dao.GetValByKeyName(keyName); err != nil {
		return
	} else {
		valExist, _ := strconv.ParseBool(data.Value)
		str := strconv.FormatBool(!valExist)
		if _, err = dao.AddKeyVal(str, existKey, keyName); err != nil {
			return
		}
	}
	return nil
}

func addConfigKeyCommon(keyName string, account string) (err error) {
	if existKey, data, err := dao.GetValByKeyName(keyName); err != nil {
		return err
	} else {
		var list []string
		var existList bool
		json.Unmarshal([]byte(data.Value), &list)
		for _, v := range list {
			if account == v {
				existList = true
				break
			}
		}
		if !existList {
			list = append(list, account)
			data, _ := json.Marshal(list)
			if _, err = dao.AddKeyVal(string(data), existKey, keyName); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteConfigKeyCommon(keyName string, account string) error {
	if existKey, data, err := dao.GetValByKeyName(keyName); err != nil {
		return err
	} else {
		var list []string
		var res []string
		var existList bool

		json.Unmarshal([]byte(data.Value), &list)
		for _, v := range list {
			if account == v {
				existList = true
			} else {
				res = append(res, v)
			}
		}
		if existList {
			resData, _ := json.Marshal(res)
			if _, err = dao.AddKeyVal(string(resData), existKey, keyName); err != nil {
				return err
			}
		}
	}
	return nil
}

func AddWhiteAccount(ctx *web.Context, account string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if account == "" {
		ctx.BackError(30005, "账号不能为空")
		return
	}
	if err := addConfigKeyCommon(constants.MAINTAIN_WHITE_LIST_KEY, account); err != nil {
		ctx.BackError(30006, "增加key出错")
		return
	}
	ctx.BackSuccess("success")
	return
}

func AddBlackAccount(ctx *web.Context, account string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if account == "" {
		ctx.BackError(30005, "账号不能为空")
		return
	}
	if err := addConfigKeyCommon(constants.BLACK_ACCOUNT_LIST_KEY, account); err != nil {
		ctx.BackError(30006, "增加key出错")
		return
	}
	ctx.BackSuccess("success")
	return
}

func AddServer(ctx *web.Context, account string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if account == "" {
		ctx.BackError(30005, "账号不能为空")
		return
	}
	if err := addConfigKeyCommon(constants.BLACK_GAME_SERVER_LIST_KEY, account); err != nil {
		ctx.BackError(30006, "增加key出错")
		return
	}
	ctx.BackSuccess("success")
	return
}

func DeleteWhiteAccount(ctx *web.Context, account string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if account == "" {
		ctx.BackError(30010, "账号不能为空")
		return
	}
	if err := deleteConfigKeyCommon(constants.MAINTAIN_WHITE_LIST_KEY, account); err != nil {
		ctx.BackError(30011, "删除白名单出错")
		return
	}

	ctx.BackSuccess("success")
	return
}

func DeleteBlackAccount(ctx *web.Context, account string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if account == "" {
		ctx.BackError(30010, "账号不能为空")
		return
	}
	if err := deleteConfigKeyCommon(constants.BLACK_ACCOUNT_LIST_KEY, account); err != nil {
		ctx.BackError(30011, "删除白名单出错")
		return
	}

	ctx.BackSuccess("success")
	return
}

func DeleteServer(ctx *web.Context, account string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if account == "" {
		ctx.BackError(30010, "账号不能为空")
		return
	}
	if err := deleteConfigKeyCommon(constants.BLACK_GAME_SERVER_LIST_KEY, account); err != nil {
		ctx.BackError(30011, "删除白名单出错")
		return
	}
	ctx.BackSuccess("success")
	return
}

func AllowPatch(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	if err := addConfigKeyBool(constants.ALLOW_PATCH_CONFIG_KEY); err != nil {
		ctx.BackError(30013, "允许热更新出错")
		return
	}
	ctx.BackSuccess("success")
	return
}
