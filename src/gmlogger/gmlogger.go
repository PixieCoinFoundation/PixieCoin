package gmlogger

import (
	"appcfg"
	"bufio"
	"bytes"
	"common"
	"constants"
	"db/cronjob"
	"fmt"
	"io"
	"io/ioutil"
	. "logger"
	"os"
	"os/exec"
	"service/params"
	"strings"
	"time"
	"tools"
)

var ODPS_BIN, DOWNLOAD_ODPS_BIN, TMP_DIR, FILTER_LOG, LOG_PREFIX string

var logPath string

func init() {
	if appcfg.GetServerType() == constants.SERVER_TYPE_SIMPLE_TEST {
		return
	}

	TMP_DIR = "/tmp/" + constants.SERVER_GLOBAL_NAME + "/"
	FILTER_LOG = "/tmp/" + constants.SERVER_GLOBAL_NAME + "/filter.log"

	formatLogPath := appcfg.GetString("format_log_path", "")
	flps := strings.Split(formatLogPath, "/")
	if len(flps) > 0 {
		ln := flps[len(flps)-1]
		lns := strings.Split(ln, ".")

		if len(lns) > 0 {
			LOG_PREFIX = lns[0]
		}
	}

	if LOG_PREFIX == "" {
		panic("config format_log_path empty or illegal")
	}

	for _, v := range flps {
		logPath += v + "/"
	}

	if !strings.HasPrefix(logPath, "/") || !strings.HasSuffix(logPath, "/") {
		panic("game_log_path illegal:" + logPath)
	}

	ODPS_BIN = appcfg.GetString("odps_bin_path", "./odpscmd/bin/odpscmd")
	DOWNLOAD_ODPS_BIN = appcfg.GetString("download_odps_bin_path", "./odpscmd/bin/odpscmd")
}

func StartGMLogger() {
	var err error
	defer func() {
		if x := recover(); x != nil {
			tools.SendErrMail(constants.MAIL_CRON_ERR_SUBJECT, "gmlogger error:"+x.(string)+","+appcfg.GetString("local_ip", "unknown"), nil, params.GetErrMailTos())
		} else if err != nil {
			tools.SendErrMail(constants.MAIL_CRON_ERR_SUBJECT, "gmlogger error:"+err.Error()+","+appcfg.GetString("local_ip", "unknown"), nil, params.GetErrMailTos())
		}
	}()

	//process yesterday's logs
	ds := time.Now().AddDate(0, 0, -1).Format("20060102")

	Info("begin gmlogger", ds, appcfg.GetAddress())

	var ok bool
	if ok, err = cronjob.LogCronJob(common.GetGMLogToken(ds, appcfg.GetAddress())); ok {
		err = Process(ds)
	} else if err == nil {
		err = constants.LogUploadDoneErr
	}
}

func Process(yday string) (err error) {
	tmp, _ := time.Parse("20060102", yday)
	tday := tmp.AddDate(0, 0, 1).Format("20060102")

	odpsTableName := appcfg.GetString("gm_log_odps_table_name", "")
	Info("process log date:", yday)

	//clear tmp log
	if err = execute("rm", "-rf", TMP_DIR); err != nil {
		return
	}

	if err = execute("mkdir", TMP_DIR); err != nil {
		return
	}

	if err = execute("touch", FILTER_LOG); err != nil {
		return
	}

	//dest log
	var df *os.File
	df, err = os.OpenFile(FILTER_LOG, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		Err(err)
		return
	}
	defer df.Close()

	var files []string
	files, err = listDir(logPath, LOG_PREFIX+"_"+yday, ".log.tar.gz")
	if err != nil {
		Err(err)
		return
	}

	apath := fmt.Sprintf("%s_%s_0000.log.tar.gz", logPath+LOG_PREFIX, tday)
	if _, err = os.Stat(apath); os.IsNotExist(err) {
		Info("file not exist:", apath)
		tools.SendInternalMail(appcfg.GetString("odps_env_name", "")+"_错误_"+yday+"_"+appcfg.GetAddress(), "file not exist:"+apath, nil)
		err = nil
	} else if err != nil {
		Err(err)
		return
	} else {
		files = append(files, apath)
	}

	Info("logs:", files)

	if len(files) > 0 {

		for _, file := range files {
			if file == fmt.Sprintf("%s_%s_0000.log.tar.gz", logPath+"/"+LOG_PREFIX, yday) {
				continue
			}

			if err = execute("tar", "-xvf", file, "-C", TMP_DIR); err != nil {
				return
			}

			var logs []string
			logs, err = listDir(TMP_DIR+logPath, LOG_PREFIX, "")
			if err != nil {
				Err(err)
				return
			}

			if len(logs) <= 0 {
				// go tools.SendInternalMail(appcfg.GetString("odps_env_name", "")+"_gmlog_abnormal_"+appcfg.GetAddress(), file, nil)
				Info("no log:", file)
			}

			for _, log := range logs {
				if filter(log, df) != nil {
					return
				}
			}

			for _, log := range logs {
				if err = execute("rm", "-rf", log); err != nil {
					return
				}
			}
		}

		opath := logPath + LOG_PREFIX + ".log-" + yday
		if _, err = os.Stat(opath); os.IsNotExist(err) {
			Info("file not exist:", opath)
		} else if err != nil {
			Err(err)
		} else {
			//other log
			filter(opath, df)
		}

		if err = execute("/bin/bash", ODPS_BIN, "-e", fmt.Sprintf("tunnel upload -fd ' %s ' -dbr true %s %s/ds='%s'", FIELD_GAP, FILTER_LOG, odpsTableName, yday)); err != nil {
			return
		}
	} else {
		err = fmt.Errorf("no log found for %s", yday)
	}
	return
}

func ProcessSimple(odpsEnvName string, ds string, dirName string, logPrefix string) (err error) {
	if odpsEnvName == "" || ds == "" || dirName == "" {
		Info("gmlog simple param error")
		return
	}

	Info("process log date:", ds)

	//clear tmp log
	if err = execute("rm", "-rf", TMP_DIR); err != nil {
		return
	}

	if err = execute("mkdir", TMP_DIR); err != nil {
		return
	}

	if err = execute("touch", FILTER_LOG); err != nil {
		return
	}

	lp := LOG_PREFIX
	if logPrefix != "" {
		lp = logPrefix
	}

	//dest log
	var df *os.File
	df, err = os.OpenFile(FILTER_LOG, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		Err(err)
		return
	}
	defer df.Close()

	var files []string
	files, err = listDir(dirName, lp+"-"+ds, ".tar.gz")
	if err != nil {
		Err(err)
		return
	}

	Info("logs:", files)

	if len(files) > 0 {
		for _, file := range files {
			if err = execute("tar", "-xvf", file, "-C", TMP_DIR); err != nil {
				return
			}

			var logs []string
			logs, err = listDir(TMP_DIR, lp, "")
			if err != nil {
				if strings.Contains(err.Error(), "no such file or directory") {
					Err("no log path:", TMP_DIR+logPath)
					err = nil
					continue
				} else {
					Err(err)
					return
				}
			}
			for _, log := range logs {
				if filter(log, df) != nil {
					return
				}
			}

			for _, log := range logs {
				if err = execute("rm", "-rf", log); err != nil {
					return
				}
			}
		}

		odpsTableName := constants.SERVER_GLOBAL_NAME + "_log_" + odpsEnvName
		if err = execute("/bin/bash", ODPS_BIN, "-e", fmt.Sprintf("tunnel upload -bs 50 -fd ' %s ' -dbr true %s %s/ds='%s'", FIELD_GAP, FILTER_LOG, odpsTableName, ds)); err != nil {
			return
		}
	} else {
		err = fmt.Errorf("no log found for %s", ds)
	}
	return
}

func filter(logpath string, df *os.File) (err error) {
	var f *os.File
	f, err = os.Open(logpath)
	if err != nil {
		Err(err)
		return
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	cnt := 0
	gmCnt := 0
	for {
		var line string
		line, err = buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				Info("file:", logpath, "lines:", cnt, "gm log lines:", gmCnt)
				err = nil
			} else {
				Err(err)
				panic(err)
			}
			return
		} else if len(line) <= 2 {
			continue
		}
		cnt++

		var fl string
		if fl = format(line); fl != "" {
			gmCnt++
			if _, err = io.WriteString(df, fl+"\n"); err != nil {
				Err(err)
				return
			}
		}
	}

	return
}

func format(line string) string {
	if len(line) < 10 {
		return ""
	}

	if strings.Contains(line, "Error 1062:") {
		return ""
	}

	ss := strings.Split(line, " ^_^ ")
	if len(ss) < 6 {
		return ""
	}

	start := 2
	if line[1:2] != "[" {
		start = 1
	}
	line = line[start : len(line)-2]

	if line[0] >= 65 && line[0] <= 90 {
		return line
	} else {
		return ""
	}
}

func execute(cmd string, args ...string) (err error) {
	Info(cmd, args)
	// var out []byte

	var cmmd *exec.Cmd
	cmmd = exec.Command(cmd, args...)
	w := bytes.NewBuffer(nil)
	cmmd.Stderr = w
	cmmd.Stdout = w

	if err = cmmd.Run(); err != nil {
		Err(err)
		Err("err:", string(w.Bytes()))
		return
	} else {
		Info(string(w.Bytes()))
	}

	return
}

func listDir(dirPth string, prefix string, suffix string) (files []string, err error) {
	Info("list:", dirPth, "pre:", prefix, "suff:", suffix)
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}

		legal := true
		if prefix != "" && !strings.HasPrefix(fi.Name(), prefix) {
			legal = false
		}

		if suffix != "" && !strings.HasSuffix(fi.Name(), suffix) {
			legal = false
		}
		if legal {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}
