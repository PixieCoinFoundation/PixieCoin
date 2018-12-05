package paper

import (
	"appcfg"
	"bytes"
	"common"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"

	. "logger"
	. "pixie_contract/api_specification"
	file "service/files"
	"strings"
	"sync"
	. "types"
)

var processPaperLock sync.Mutex

func cleanProcessDirectory() {
	common.Execute("rm", "-rf", appcfg.GetProjectPathPrefix()+"GFTP/icons/")
	common.Execute("rm", "-rf", appcfg.GetProjectPathPrefix()+"GFTP/clothes/")
	common.Execute("rm", "-rf", appcfg.GetProjectPathPrefix()+"GFTP/output/")
	common.Execute("rm", "-rf", appcfg.GetProjectPathPrefix()+"GFTP/output_pvr/")
}

func makeProcessDirectory() {
	common.Execute("mkdir", appcfg.GetProjectPathPrefix()+"GFTP/icons")
	common.Execute("mkdir", appcfg.GetProjectPathPrefix()+"GFTP/clothes")
	common.Execute("mkdir", appcfg.GetProjectPathPrefix()+"GFTP/output")
	common.Execute("mkdir", appcfg.GetProjectPathPrefix()+"GFTP/output_pvr")
}

func processLock() {
	processPaperLock.Lock()

	cleanProcessDirectory()
	makeProcessDirectory()
}

func processUnlock() {
	cleanProcessDirectory()

	processPaperLock.Unlock()
}

func ProcessFile(username string, extra *PaperExtra, icon string, main string, bottom string, collar string, shadow string, tNow int64) (err error) {
	processLock()
	defer processUnlock()

	var iconTemp, mainTemp, bottomTemp, collarTemp, shadowTemp string

	iconTemp = fmt.Sprintf("paper_%s_%d_icon", username, tNow)
	iconSavePath := appcfg.GetProjectPathPrefix() + "GFTP/icons/" + iconTemp
	if err = file.DecryptDownloadFile(icon, iconSavePath); err != nil {
		return
	}

	mainTemp = fmt.Sprintf("paper_%s_%d_main", username, tNow)
	mainSavePath := appcfg.GetProjectPathPrefix() + "GFTP/clothes/" + mainTemp
	if err = file.DecryptDownloadFile(main, mainSavePath); err != nil {
		return
	}

	if bottom != "" {
		bottomTemp = fmt.Sprintf("paper_%s_%d_bottom", username, tNow)
		bottomSavePath := appcfg.GetProjectPathPrefix() + "GFTP/clothes/" + bottomTemp
		if err = file.DecryptDownloadFile(bottom, bottomSavePath); err != nil {
			return
		}
	}
	if collar != "" {
		collarTemp = fmt.Sprintf("paper_%s_%d_collar", username, tNow)
		collarSavePath := appcfg.GetProjectPathPrefix() + "GFTP/clothes/" + collarTemp
		if err = file.DecryptDownloadFile(collar, collarSavePath); err != nil {
			return
		}
	}

	if shadow != "" {
		shadowTemp = fmt.Sprintf("paper_%s_%d_shadow", username, tNow)
		shadowSavePath := appcfg.GetProjectPathPrefix() + "GFTP/clothes/" + shadowTemp
		if err = file.DecryptDownloadFile(shadow, shadowSavePath); err != nil {
			return
		}
	}

	//process file
	var res string
	if res, err = common.Execute("/bin/bash", appcfg.GetProjectPathPrefix()+"process_file.sh"); err == nil {
		var infos []ExtraInfo
		if err = json.Unmarshal([]byte("["+strings.Replace(res, "}{", "},{", -1)+"]"), &infos); err != nil {
			Err(err)
			return
		}
		for _, i := range infos {
			if strings.HasSuffix(iconTemp, i.Name) {
				extra.IconX = i.X
				extra.IconY = i.Y
			} else if strings.HasSuffix(mainTemp, i.Name) {
				extra.MainX = i.X
				extra.MainY = i.Y
			} else if strings.HasSuffix(bottomTemp, i.Name) {
				extra.BottomX = i.X
				extra.BottomY = i.Y
			} else if strings.HasSuffix(collarTemp, i.Name) {
				extra.CollarX = i.X
				extra.CollarY = i.Y
			} else if strings.HasSuffix(shadowTemp, i.Name) {
				extra.ShadowX = i.X
				extra.ShadowY = i.Y
			}
		}
	} else {
		return
	}

	if appcfg.GetBool("local_test", false) {
		Info("local test just cp from downloads")
		common.Execute("cp", appcfg.GetProjectPathPrefix()+"GFTP/clothes/"+iconTemp, appcfg.GetProjectPathPrefix()+"GFTP/output")
		common.Execute("cp", appcfg.GetProjectPathPrefix()+"GFTP/clothes/"+iconTemp, appcfg.GetProjectPathPrefix()+"GFTP/output_pvr")
		if mainTemp != "" {
			common.Execute("cp", appcfg.GetProjectPathPrefix()+"GFTP/clothes/"+mainTemp, appcfg.GetProjectPathPrefix()+"GFTP/output")
			common.Execute("cp", appcfg.GetProjectPathPrefix()+"GFTP/clothes/"+mainTemp, appcfg.GetProjectPathPrefix()+"GFTP/output_pvr")
		}
		if bottomTemp != "" {
			common.Execute("cp", appcfg.GetProjectPathPrefix()+"GFTP/clothes/"+bottomTemp, appcfg.GetProjectPathPrefix()+"GFTP/output")
			common.Execute("cp", appcfg.GetProjectPathPrefix()+"GFTP/clothes/"+bottomTemp, appcfg.GetProjectPathPrefix()+"GFTP/output_pvr")
		}
		if collarTemp != "" {
			common.Execute("cp", appcfg.GetProjectPathPrefix()+"GFTP/clothes/"+collarTemp, appcfg.GetProjectPathPrefix()+"GFTP/output")
			common.Execute("cp", appcfg.GetProjectPathPrefix()+"GFTP/clothes/"+collarTemp, appcfg.GetProjectPathPrefix()+"GFTP/output_pvr")
		}
		if shadowTemp != "" {
			common.Execute("cp", appcfg.GetProjectPathPrefix()+"GFTP/clothes/"+shadowTemp, appcfg.GetProjectPathPrefix()+"GFTP/output")
			common.Execute("cp", appcfg.GetProjectPathPrefix()+"GFTP/clothes/"+shadowTemp, appcfg.GetProjectPathPrefix()+"GFTP/output_pvr")
		}
	}

	var md5 string

	// upload
	if shadowTemp != "" {
		shadowFile := appcfg.GetProjectPathPrefix() + "GFTP/output/" + shadowTemp
		if err, w, h := getBounds(shadowFile); err != nil {
			return err
		} else {
			extra.ShadowW = w
			extra.ShadowH = h
		}
		if err = file.Encrypt(shadowFile); err != nil {
			return
		}
		if md5, err = common.GenMD5F(shadowFile); err != nil {
			return
		}
		extra.ShadowG = fmt.Sprintf("%s/%s_shadow_%d_%s", file.GetPlatformOSSDirName(), username, tNow, md5)
		if _, err = file.UploadFileByName(shadowFile, extra.ShadowG, 0); err != nil {
			return
		}
	}

	//process in pvr format
	// if res, err = common.Execute("/bin/bash", appcfg.GetProjectPathPrefix()+"process_file_pvr.sh"); err == nil {
	iconFile := appcfg.GetProjectPathPrefix() + "GFTP/output/" + iconTemp
	if err, w, h := getBounds(iconFile); err != nil {
		return err
	} else {
		extra.IconW = w
		extra.IconH = h
	}

	if err = file.Encrypt(iconFile); err != nil {
		return
	}
	if md5, err = common.GenMD5F(iconFile); err != nil {
		return
	}
	extra.IconG = fmt.Sprintf("%s/%s_icon_%d_%s", file.GetPlatformOSSDirName(), username, tNow, md5)
	if _, err = file.UploadFileByName(iconFile, extra.IconG, 0); err != nil {
		return
	}

	mainFile := appcfg.GetProjectPathPrefix() + "GFTP/output/" + mainTemp
	if err, w, h := getBounds(mainFile); err != nil {
		return err
	} else {
		extra.MainW = w
		extra.MainH = h
	}

	if err = file.Encrypt(mainFile); err != nil {
		return
	}
	if md5, err = common.GenMD5F(mainFile); err != nil {
		return
	}
	extra.MainG = fmt.Sprintf("%s/%s_main_%d_%s", file.GetPlatformOSSDirName(), username, tNow, md5)
	if _, err = file.UploadFileByName(mainFile, extra.MainG, 0); err != nil {
		return
	}

	if bottomTemp != "" {
		bottomFile := appcfg.GetProjectPathPrefix() + "GFTP/output/" + bottomTemp
		if err, w, h := getBounds(bottomFile); err != nil {
			return err
		} else {
			extra.BottomW = w
			extra.BottomH = h
		}
		if err = file.Encrypt(bottomFile); err != nil {
			return
		}
		if md5, err = common.GenMD5F(bottomFile); err != nil {
			return
		}
		extra.BottomG = fmt.Sprintf("%s/%s_bottom_%d_%s", file.GetPlatformOSSDirName(), username, tNow, md5)
		if _, err = file.UploadFileByName(bottomFile, extra.BottomG, 0); err != nil {
			return
		}
	}

	if collarTemp != "" {
		collarFile := appcfg.GetProjectPathPrefix() + "GFTP/output/" + collarTemp
		if err, w, h := getBounds(collarFile); err != nil {
			return err
		} else {
			extra.CollarW = w
			extra.CollarH = h
		}
		if err = file.Encrypt(collarFile); err != nil {
			return
		}
		if md5, err = common.GenMD5F(collarFile); err != nil {
			return
		}
		extra.CollarG = fmt.Sprintf("%s/%s_collar_%d_%s", file.GetPlatformOSSDirName(), username, tNow, md5)
		if _, err = file.UploadFileByName(collarFile, extra.CollarG, 0); err != nil {
			return
		}
	}
	return
}

func getBounds(fileName string) (err error, w, h float64) {
	var data []byte
	var m image.Image
	data, err = ioutil.ReadFile(fileName)
	if err != nil {
		return err, 0, 0
	}
	m, _, err = image.Decode(bytes.NewBuffer(data))
	if err != nil {
		return err, 0, 0
	}
	return nil, float64(m.Bounds().Dx()), float64(m.Bounds().Dy())
}
