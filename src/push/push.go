package push

import (
	"appcfg"
	"constants"
	"time"
)

const (
	TPUSH_ONE_ADDR = "openapi.xg.qq.com/v2/push/single_account"
	TPUSH_ALL_ADDR = "openapi.xg.qq.com/v2/push/all_device"
)

var ios_environment = "2"
var tpush_android_access_id, tpush_ios_access_id string
var tpush_android_secret_key, tpush_ios_secret_key string

var cnLocation *time.Location

func init() {
	if appcfg.GetServerType() == "" || appcfg.GetServerType() == constants.SERVER_TYPE_GV {
		if appcfg.GetBool("push_production_environment", false) {
			ios_environment = "1"
		}

		tpush_android_access_id = appcfg.GetString("tpush_android_access_id", "")
		tpush_ios_access_id = appcfg.GetString("tpush_ios_access_id", "")
		tpush_android_secret_key = appcfg.GetString("tpush_android_secret_key", "")
		tpush_ios_secret_key = appcfg.GetString("tpush_ios_secret_key", "")

		cnLocation, _ = time.LoadLocation("Asia/Chongqing")
	}
}

func TPushToOne(username, title, content string) {
	if tpush_ios_secret_key != "" && tpush_ios_access_id != "" {
		tPushOneIOS(username, title, content)
	}

	if tpush_android_secret_key != "" && tpush_android_access_id != "" {
		tPushOneAndroid(username, title, content)
	}
}

func TPushToAll(title, content string) {
	if tpush_ios_secret_key != "" && tpush_ios_access_id != "" {
		tPushAllIOS(title, content)
	}
	if tpush_android_secret_key != "" && tpush_android_access_id != "" {
		tPushAllAndroid(title, content)
	}
}

func TPushToAllAndroid(title, content string) {
	if tpush_android_secret_key != "" && tpush_android_access_id != "" {
		tPushAllAndroid(title, content)
	}
}

func TPushToAllIOS(title, content string) {
	if tpush_ios_secret_key != "" && tpush_ios_access_id != "" {
		tPushAllIOS(title, content)
	}
}
