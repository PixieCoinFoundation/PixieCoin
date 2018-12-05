package version

import (
	"constants"
	"encoding/json"
	"fmt"
	"github.com/web"
	"gm/dao"
	. "logger"
	. "model"
	"net/http"
	"strings"
	"sync"
	"time"
)

var updateVersionLock sync.Mutex

//获取版本信息
func GetVersionInfo(ctx *web.Context) {
	if err, res := dao.GetVersionInfo(); err != nil {
		ctx.BackError(20001, err)
		return
	} else {
		ctx.BackSuccess(res)
	}
}

type UpdateVersionInfoReq struct {
	CoreVersion   string
	DownLoadUrl   string
	ScriptVersion string
}

//更新版本
func UpdateVersionInfo(ctx *web.Context) {
	updateVersionLock.Lock()
	defer updateVersionLock.Unlock()

	var vn Version
	pack := make(map[string]string)
	var input UpdateVersionInfoReq
	var err error
	tNow := time.Now().UnixNano()

	if err = ctx.GetFormData(&input); err != nil {
		ctx.BackError(10015, "解析数据出错")
		return
	}

	if input.CoreVersion == "" || input.ScriptVersion == "" || input.DownLoadUrl == "" {
		ctx.BackError(10018, "参数不能为空")
		return
	}

	if !strings.HasPrefix(input.DownLoadUrl, "http://") && !strings.HasPrefix(input.DownLoadUrl, "https://") {
		um := make(map[string]string)
		if err := json.Unmarshal([]byte(input.DownLoadUrl), &um); err != nil {
			Err(err)
			ctx.BackError(10020, "游戏下载路径格式错误")
			return
		} else {
			if um["DEFAULT"] == "" || um["IOS"] == "" {
				ctx.BackError(10021, "游戏下载路径错误：url map must contain keys:DEFAULT,IOS")
				return
			}
		}
	}

	//获取当前版本信息
	if err, res := dao.GetVersionInfo(); err != nil {
		ctx.BackError(20001, err)
		return
	} else {
		var nowVersionFile map[string]string
		if err = json.Unmarshal([]byte(res.VersionFile), &nowVersionFile); err != nil {
			ctx.BackError(20002, err)
			return
		}

		var iOSPack, androidPack string
		iOSPack, err = ctx.UploadFormFile("iOSPack", fmt.Sprintf("patch_ios_%d", tNow), true)
		if err == http.ErrMissingFile {
			//没有上传ios更新包 使用原有的
			pack[constants.THIRD_CHANNEL_IOS] = nowVersionFile[constants.THIRD_CHANNEL_IOS]
		} else if err != nil {
			ctx.BackError(10022, "处理ios更新包失败:"+err.Error())
			return
		} else {
			pack[constants.THIRD_CHANNEL_IOS] = iOSPack
		}

		androidPack, err = ctx.UploadFormFile("androidPack", fmt.Sprintf("patch_android_%d", tNow), true)
		if err == http.ErrMissingFile {
			//没有上传安卓更新包 使用原有的
			pack[constants.THIRD_CHANNEL_ANDROID] = nowVersionFile[constants.THIRD_CHANNEL_ANDROID]
		} else if err != nil {
			ctx.BackError(10023, "处理android更新包失败:"+err.Error())
			return
		} else {
			pack[constants.THIRD_CHANNEL_ANDROID] = androidPack
		}
	}

	data, _ := json.Marshal(pack)
	vn.VersionFile = string(data)
	vn.CoreVersion = input.CoreVersion
	vn.ScriptVersion = input.ScriptVersion
	vn.DownLoadUrl = input.DownLoadUrl

	if err := dao.UpdateVersion(vn); err != nil {
		ctx.BackError(20022, err)
		return
	}
}

//更新提审
func UpdateVersionTishen(ctx *web.Context) {
	var vn Version
	vn.Tishen = 1

	if err := dao.UpdateVersion(vn); err != nil {
		ctx.BackError(20004, err)
		return
	}
}
