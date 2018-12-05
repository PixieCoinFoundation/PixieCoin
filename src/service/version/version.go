package version

import (
	"appcfg"
	"constants"
	"fmt"
	"gl_db/config"
	"gl_db/version"
	. "logger"
	"os"
	"time"
	. "types"
)

var f1 string
var f2 string
var f3 string

var AllowPatch bool

func loadPatchStatus() {
	if config.GetConfig(constants.ALLOW_PATCH_CONFIG_KEY) == "true" {
		AllowPatch = true
	} else {
		Info("forbid online patch!")
		AllowPatch = false
	}
}

func refreshPatchStatus() {
	for {
		time.Sleep(60 * time.Second)
		loadPatchStatus()
	}
}

func init() {
	if appcfg.GetServerType() == constants.SERVER_TYPE_GL {
		f1 = appcfg.GetString("file_server", "")
		f2 = appcfg.GetString("file_server2", "")
		f3 = appcfg.GetString("file_server3", "")
		if f1 == "" || f2 == "" || f3 == "" {
			fmt.Println("wrong file_servers.")
			os.Exit(1)
		}

		loadPatchStatus()

		go refreshPatchStatus()
	}
}

// func SetNewVersion(coreVersion string, downloadUrl string, downloadUrl2, scriptVersion string, versionFile string) error {
// 	return version.UpdateVersion(coreVersion, downloadUrl, downloadUrl2, scriptVersion, versionFile)
// }

func GetCurrentVersion(thirdChannel string) (gfVersion GFVersion, tvf string, err error) {
	if gfVersion, tvf, err = version.GetCurrentVersion(thirdChannel); err != nil {
		return
	}
	gfVersion.FileServerAddress = f1
	gfVersion.FileServerAddress2 = f2
	gfVersion.FileServerAddress3 = f3
	return
}

func GetFileServers() (string, string, string) {
	return f3, f1, f2
}
