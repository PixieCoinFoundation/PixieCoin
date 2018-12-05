package common

import (
	"fmt"
	"service/files"
	"time"
	"tools"
)

func DownImg(saveNameParam string, imgPath string) (savePath string) {
	prefix := "./gm"
	savePath = "/images/recommend_" + saveNameParam + "_banner"
	if !tools.ExistsFile(prefix + savePath) {
		files.DecryptDownloadFile(imgPath, prefix+savePath)
	}

	return
}

func GetTokenDate() string {
	var hour int
	if time.Now().Hour()%2 == 0 {
		hour = time.Now().Hour()
	} else {
		hour = time.Now().Hour() - 1
	}
	return fmt.Sprintf("%d-%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day(), hour)
}
