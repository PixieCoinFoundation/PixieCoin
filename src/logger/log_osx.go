// +build darwin

package logger

import (
	"appcfg"
	"constants"
	"dao_korlog"
	glog "db/log"
	"encoding/json"
	"fmt"
	"tools"
	"types"

	"gopkg.in/goftp"
	"log"
	. "model"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

var (
	// sLog        *syslog.Writer
	// bilibiliLog *syslog.Writer

	isProfiling bool

	// serverId string
	yace     bool
	logDebug bool

	gmlogChan chan BufferLog

	shutdownChan chan struct{}
	wg           sync.WaitGroup
)

const (
	FIELD_GAP = "^_^"
	TMP_SIZE  = 10
)

func init() {
	yace = appcfg.GetBool("yace", false)
	logDebug = appcfg.GetBool("log_debubg", false)

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	isProfiling = appcfg.GetBool("profile", false)
	gmlogChan = make(chan BufferLog, appcfg.GetInt("gmlog_buffer", 1000))
	shutdownChan = make(chan struct{})

	for i := 0; i < appcfg.GetInt("gmlog_worker", 10); i++ {
		wg.Add(1)
		go startGMLogWorker(i)
	}
}

//SendKorLogftp 发送ftp
func SendKorLogftp(fileTimename, localfile string) {
	fplogin1 := fmt.Sprintf("cmemberloginlog/%d.user_cmemberloginlog_%s.csv", constants.KOR_GAME_CODE, fileTimename)
	fporder := fmt.Sprintf("paylog/%d.PayLog_%s.csv", constants.KOR_GAME_CODE, fileTimename)
	localfplogin1 := constants.DefaultPath + "korlogin1log-" + localfile + ".csv"
	localfporder := constants.DefaultPath + "kororderlog-" + localfile + ".csv"
	if ftp, err := goftp.Connect("112.175.37.100:21"); err != nil {
		Err(err)
		return
	} else {
		defer ftp.Close()
		if err = ftp.Login("bi_user02", "bi_user02!%"); err != nil {
			Err(err)
			return
		}

		var filelogin1 *os.File
		if filelogin1, err = os.Open(localfplogin1); err != nil {
			Err(err)
			return
		}
		if err = ftp.Stor(fplogin1, filelogin1); err != nil {
			Err(err)
			return
		}
		var fileorder *os.File
		if fileorder, err = os.Open(localfporder); err != nil {
			Err(err)
			return
		}
		if err = ftp.Stor(fporder, fileorder); err != nil {
			Err(err)
			return
		}
		filelogin1.Close()
		fileorder.Close()
	}
}
func writeData(logtype string, data []string) error {
	file, err, exist := logFile(logtype)
	if !exist {
		if logtype == "korlogin1log" {
			filedata := fmt.Sprintf("%s\n", constants.KOR_LOGIN1_TABLEHEAD)
			file.WriteString(filedata)

		} else if logtype == "kororderlog" {
			filedata := fmt.Sprintf("%s\n", constants.KOR_ORDER_TABLEHEAD)
			file.WriteString(filedata)
		}
	}
	//var loginModule types.LoginKorLog
	var login1Module types.LoginKorLog1
	var orderModuel types.OrderKorLog
	for _, v := range data {
		var medData string
		if logtype == "korlogin1log" {
			json.Unmarshal([]byte(v), &login1Module)
			medData = fmt.Sprintf("%s,%d,%d,%d,%s,%d,%d,%d,%d,%d,%d\n", login1Module.Logindate, login1Module.Logindateseq, login1Module.Cuid, login1Module.Isloginsuccessdiv, login1Module.Platformuserid, login1Module.Loginrootcode, login1Module.GameCode, login1Module.Deviceid, login1Module.Platformcode, login1Module.Oscode, login1Module.Ostopverno)
		} else if logtype == "kororderlog" {
			json.Unmarshal([]byte(v), &orderModuel)
			medData = fmt.Sprintf("%d,%s,%d,%d,%s,%d,%d,%d,%d,%s,%d,%d,%s,%d,%d,%s,%f\n", orderModuel.Paylogid, orderModuel.Eventdate, orderModuel.Logtype, orderModuel.Guid, orderModuel.Platformuserid, orderModuel.Platformcode, orderModuel.Oscode, orderModuel.Storetype, orderModuel.Orderid, orderModuel.Paydate, orderModuel.Paytool, orderModuel.Gamecode, orderModuel.Gamename, orderModuel.Itemcode, orderModuel.Itemname, orderModuel.Currentcycode, orderModuel.Amount)
		}
		_, err = file.Write([]byte(medData))
		if err != nil {
			fmt.Println("filepath not exist ")
		}
	}
	return err
}
func logFile(filename string) (*os.File, error, bool) {
	rootDir := fmt.Sprintf("/tmp/GMKorlog/")
	logFilePath := fmt.Sprintf("%s%s-%s-%d.csv", rootDir, filename, time.Now().Format("2006-01-02"), time.Now().Hour())
	exist := checkFileIsExist(logFilePath)
	fileinfo, err := os.Stat(rootDir)
	// 当文件夹不存在时
	if err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(rootDir, 0700); err != nil {
		}
	} else if err != nil { // 文件夹存在,不过出现了错误
		content := fmt.Sprintf("%s", err)
		tools.SendInternalMail("create folder err", content, nil)
	} else if err == nil { // 文件夹存在,无错误
		// 当日志目录不是目录时
		if !fileinfo.IsDir() {
			content := fmt.Sprintf("%s", err)
			tools.SendInternalMail("path is not a  folder", content, nil)
		}
	}
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	return file, err, exist

}
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
func End() {
	Info("gm log workers closing...")
	close(shutdownChan)
	wg.Wait()
	Info("gm log workers closed")
}

func startGMLogWorker(i int) {
	Info(fmt.Sprint("log worker:", i, "started"))
	tmpLogs := make([]BufferLog, 0)
	for {
		select {
		case <-shutdownChan:
			Info("close gm log worker:", i)
			//clear all log
			for _, l := range tmpLogs {
				if l.C3 == constants.C3_LOGIN_LOG_KOR || l.C3 == constants.C3_LOGIN1_LOG_KOR || l.C3 == constants.C3_ORDER_LOG_KOR {
					logToDBKor(l.C1, l.C2, l.C3, l.Username, l.T64, l.Extra)
				} else {
					logToDB(l.C1, l.C2, l.C3, l.Username, l.T64, l.Extra)
				}
			}

			wg.Done()
			return
		case l := <-gmlogChan:
			Info(l.C1, FIELD_GAP, l.C2, FIELD_GAP, l.C3, FIELD_GAP, l.Username, FIELD_GAP, l.T, FIELD_GAP, l.Extra)
			if l.C3 != constants.C3_API_COST {
				if len(tmpLogs) == TMP_SIZE {
					logToDB10(tmpLogs)
					tmpLogs = make([]BufferLog, 0)
				}

				tmpLogs = append(tmpLogs, l)
			}

		}
	}
}

func Info(v ...interface{}) {
	if !yace {
		log.Println(fmt.Sprint(v))
	}
}

//这个方法记录日志 正常只在非正式环境记录
func DebugInfo(v ...interface{}) {
	if logDebug {
		log.Println(fmt.Sprint(v))
	}
}

func Err(v ...interface{}) {
	log.Println(fmt.Sprint(v, "\n[stack]", string(debug.Stack())))
}

func ErrMail(v ...interface{}) {
	log.Println(fmt.Sprint(v, "\n[stack]", string(debug.Stack())))
}

func OrderLog(channel string, v ...interface{}) {
	log.Println(channel + ":" + fmt.Sprint(v))
}

func Debug(v ...interface{}) {
	if !yace {
		log.Println(fmt.Sprint(v))
	}

}

func Profile(c1 string, c2 string, c3 string, username string, extra string) {
	if isProfiling {
		info(FIELD_GAP, c1, FIELD_GAP, c2, FIELD_GAP, c3, FIELD_GAP, username, FIELD_GAP, time.Now().Format("2006-01-02 15:04:05"), FIELD_GAP, extra)
	}
}

func info(v ...interface{}) {
	log.Println(fmt.Sprint(v))
}

func GMLog(c1 string, c2 string, c3 string, username string, extra string) {
	now := time.Now()
	if len(gmlogChan) > appcfg.GetInt("gmlog_buffer", 1000)*9/10 {
		Info(FIELD_GAP, c1, FIELD_GAP, c2, FIELD_GAP, c3, FIELD_GAP, username, FIELD_GAP, now.Format("2006-01-02 15:04:05"), FIELD_GAP, extra)
		logToDB(c1, c2, c3, username, now.Unix(), extra)
	} else {
		gmlogChan <- BufferLog{C1: c1, C2: c2, C3: c3, Username: username, Extra: extra, T: now.Format("2006-01-02 15:04:05"), T64: now.Unix()}
	}
}

func logToDB(c1 string, c2 string, c3 string, username string, t int64, extra string) {
	// if !yace {
	// 	if appcfg.GetServerType() == "" && c3 != constants.C3_API_COST {
	// 		glog.AddLog(c1, c2, c3, username, t, extra)
	// 	} else if appcfg.GetServerType() == constants.SERVER_TYPE_GM {
	// 		gmlog.AddLog(c1, c2, c3, username, t, extra)
	// 	}
	// }
}

//GMLogKor 韩国版本日志插入
func GMLogKor(c1 string, c2 string, c3 string, username string, extra string) {
	now := time.Now()
	Info(FIELD_GAP, c1, FIELD_GAP, c2, FIELD_GAP, c3, FIELD_GAP, username, FIELD_GAP, now.Format("2006-01-02 15:04:05"), FIELD_GAP, extra)
	logToDBKor(c1, c2, c3, username, now.Unix(), extra)

}

//韩国版本日志
func logToDBKor(c1 string, c2 string, c3 string, username string, t int64, extra string) {
	if c3 == constants.C3_LOGIN_LOG_KOR {
		logKorAllLog(extra)
	} else if c3 == constants.C3_LOGIN1_LOG_KOR {
		korLoginLog(extra)
	} else if c3 == constants.C3_ORDER_LOG_KOR {
		korOrderlog(extra)
	}
}

//logKorAllLog 插入登录日志
func logKorAllLog(jsonStr string) {
	var loginModule types.LoginKorLog
	if err := json.Unmarshal([]byte(jsonStr), &loginModule); err != nil {
		Err(err)
	}
	stmt := dao_korlog.AddKorLogin1LogStmt
	_, err := stmt.Exec(loginModule.Guid, loginModule.Uuid, loginModule.Cuid, loginModule.Gamecode, loginModule.Nickname, loginModule.Sex, loginModule.Age, loginModule.LevYn, loginModule.UserLv, loginModule.Authtoken, loginModule.Oscode, loginModule.Ostopverno, loginModule.Tempuserid, loginModule.TokenredisDate, loginModule.Registdate, loginModule.Mappingdate)
	if err != nil {
		Err(err)
	}

}

//korloginLog 插入会员登日志
func korLoginLog(jsonStr string) {
	var login1Module types.LoginKorLog1
	if err := json.Unmarshal([]byte(jsonStr), &login1Module); err != nil {
		Err(err)
	}
	stmt := dao_korlog.AddKorLoginLogStmt

	_, err := stmt.Exec(login1Module.Logindate, login1Module.Cuid, login1Module.Isloginsuccessdiv, login1Module.Platformuserid, login1Module.Loginrootcode, login1Module.GameCode, login1Module.Deviceid, login1Module.Platformcode, login1Module.Oscode, login1Module.Ostopverno)
	if err != nil {
		Err(err)
	}
}
func korOrderlog(jsonStr string) {
	var OrderKorLogModule types.OrderKorLog
	if err := json.Unmarshal([]byte(jsonStr), &OrderKorLogModule); err != nil {
		Err(err)
	}
	stmt := dao_korlog.AddKorOrderLogStmt
	_, err := stmt.Exec(OrderKorLogModule.Eventdate, OrderKorLogModule.Logtype, OrderKorLogModule.Guid, OrderKorLogModule.Platformuserid, OrderKorLogModule.Platformcode, OrderKorLogModule.Oscode, OrderKorLogModule.Storetype, OrderKorLogModule.Orderid, OrderKorLogModule.Paydate, OrderKorLogModule.Paytool, OrderKorLogModule.Gamecode, OrderKorLogModule.Gamename, OrderKorLogModule.Itemcode, OrderKorLogModule.Itemname, OrderKorLogModule.Currentcycode, OrderKorLogModule.Amount, OrderKorLogModule.Paycountry)
	if err != nil {
		Err(err)
	}
}
func logToDB10(logs []BufferLog) {
	if !yace {
		if appcfg.GetServerType() == "" {
			glog.AddLog10(logs)
		} else if appcfg.GetServerType() == constants.SERVER_TYPE_GM {
			// for _, l := range logs {
			// 	gmlog.AddLog(l.C1, l.C2, l.C3, l.Username, l.T64, l.Extra)
			// }
		}
	}
}
