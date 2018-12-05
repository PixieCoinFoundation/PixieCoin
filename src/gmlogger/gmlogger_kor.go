package gmlogger

import (
	"appcfg"
	// "bufio"
	// "bytes"
	// "common"
	"constants"
	// "db/cronjob"
	"fmt"
	"gopkg.in/goftp"
	// "io"
	// "io/ioutil"
	. "logger"
	"os"
	// "os/exec"
	// "service/params"
	// "strings"
	"time"
	"tools"
)

func init() {
	if appcfg.GetBool("test_upload_kor_log", false) {
		MailKorOdpsTable()
	}
}

func MailOdpsTable(odpsCmtName, tableName, ds, rcs string) {
	fn := fmt.Sprintf("%s_%s_%d.csv", tableName, ds, time.Now().Unix())
	fp := "downloads/download_table/" + fn
	if err := execute("/bin/bash", odpsCmtName, "-e", fmt.Sprintf("tunnel download -h true %s/ds='%s' %s", tableName, ds, fp)); err != nil {
		return
	}

	tools.SendErrMail(tableName+"_"+ds, "", []string{fp}, rcs)
}

func MailKorOdpsTable() {
	ds := time.Now().AddDate(0, 0, -1).Format("20060102")

	MailKorOdpsTable1(ds)
}

func MailKorOdpsTable1(ds string) {
	now := time.Now().Unix()
	fn := fmt.Sprintf("login_kor_%s_%d.csv", ds, now)
	fp := "downloads/download_table/" + fn
	if err := execute("/bin/bash", DOWNLOAD_ODPS_BIN, "-e", fmt.Sprintf("tunnel download -h true login_kor/ds='%s' %s", ds, fp)); err != nil {
		return
	}

	fn1 := fmt.Sprintf("login_kor_1_%s_%d.csv", ds, now)
	fp1 := "downloads/download_table/" + fn1
	if err := execute("/bin/bash", DOWNLOAD_ODPS_BIN, "-e", fmt.Sprintf("tunnel download -h true login_kor_1/ds='%s' %s", ds, fp1)); err != nil {
		return
	}

	fn2 := fmt.Sprintf("order_kor_%s_%d.csv", ds, now)
	fp2 := "downloads/download_table/" + fn2
	if err := execute("/bin/bash", DOWNLOAD_ODPS_BIN, "-e", fmt.Sprintf("tunnel download -h true order_kor/ds='%s' %s", ds, fp2)); err != nil {
		return
	}

	// tools.SendErrMail(appcfg.GetString("odps_env_name", "")+"_login_report_"+ds, appcfg.GetString("odps_env_name", "")+"_log_report_"+ds, []string{fp, fp1, fp2}, params.GetErrMailTos())

	//upload kor ftp
	if ftp, err := goftp.Connect("112.175.37.100:21"); err != nil {
		Err(err)
		return
	} else {
		defer ftp.Close()

		// Username / password authentication
		if err = ftp.Login("bi_user02", "bi_user02!%"); err != nil {
			Err(err)
			return
		}

		var file *os.File
		if file, err = os.Open(fp); err != nil {
			Err(err)
			return
		}
		if err = ftp.Stor(fmt.Sprintf("gamelogininfo/%d.user_gamelogininfo_%s.csv", constants.KOR_GAME_CODE, ds), file); err != nil {
			Err(err)
			return
		}

		if file, err = os.Open(fp1); err != nil {
			Err(err)
			return
		}
		if err = ftp.Stor(fmt.Sprintf("cmemberloginlog/%d.user_cmemberloginlog_%s.csv", constants.KOR_GAME_CODE, ds), file); err != nil {
			Err(err)
			return
		}

		if file, err = os.Open(fp2); err != nil {
			Err(err)
			return
		}
		if err = ftp.Stor(fmt.Sprintf("paylog/%d.PayLog_%s.csv", constants.KOR_GAME_CODE, ds), file); err != nil {
			Err(err)
			return
		}
	}
}
