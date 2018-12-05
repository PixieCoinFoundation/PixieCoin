// +build linux

package logger

import (
	"appcfg"
	"constants"
	"dao_korlog"
	glog "db/log"
	"encoding/json"
	"fmt"
	// gmlog "gm_db/log"
	// gvlog "gv_db/log"
	"types"

	"github.com/sirupsen/logrus"
	. "model"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
	"tools"
)

var (
	infoLogPath   string
	formatLogPath string
	infoF         *os.File
	formatF       *os.File

	isProfiling bool

	yace     bool
	logDebug bool

	gmlogChan    chan BufferLog
	shutdownChan chan struct{}
	wg           sync.WaitGroup

	sLog      *logrus.Logger
	formatLog *logrus.Logger
)

const (
	FIELD_GAP = "^_^"
	TMP_SIZE  = 10
)

func init() {
	var err error
	infoLogPath = appcfg.GetString("info_log_path", "")
	formatLogPath = appcfg.GetString("format_log_path", "")

	if infoLogPath == "" || formatLogPath == "" || !strings.HasSuffix(infoLogPath, ".log") || !strings.HasSuffix(formatLogPath, ".log") {
		panic("config info_log_path or json_log_path empty or illegal")
	}

	if infoF, err = os.OpenFile(infoLogPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666); err != nil {
		panic(err)
	}
	if formatF, err = os.OpenFile(formatLogPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666); err != nil {
		panic(err)
	}

	sLog = logrus.New()
	formatLog = logrus.New()

	sLog.Out = infoF
	formatLog.Out = formatF
	formatLog.Formatter = new(GMTextFormatter)

	yace = appcfg.GetBool("yace", false)
	logDebug = appcfg.GetBool("log_debubg", false)

	isProfiling = appcfg.GetBool("profile", false)

	gmlogChan = make(chan BufferLog, appcfg.GetInt("gmlog_buffer", 1000))

	shutdownChan = make(chan struct{})

	for i := 0; i < appcfg.GetInt("gmlog_worker", 10); i++ {
		wg.Add(1)
		go startGMLogWorker(i)
	}

	if appcfg.GetLanguage() == constants.KOR_LANGUAGE && appcfg.GetServerType() == "" && appcfg.GetBool("main_server", false) {
		//韩国日志文件定时上传
		go loopCleanKorLog()
	}
}

func loopCleanKorLog() {
	for {
		et := time.Now().Unix() - 32*24*3600

		var delId int64
		if logs, err := glog.ListKorLog(); err != nil {
			return
		} else {
			for _, l := range logs {
				if l.Time < et && l.ID > delId {
					delId = l.ID
				}
			}
		}

		if delId > 0 {
			glog.DelKorLog(delId)
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(60 * time.Second)
		}
	}
}

func End() {
	Info("gm log workers closing...")
	close(shutdownChan)
	wg.Wait()

	if infoF != nil {
		infoF.Close()
	}
	if formatF != nil {
		formatF.Close()
	}
	Info("gm log workers closed")
}

func startGMLogWorker(i int) {
	sLog.Info(fmt.Sprint("log worker: ", i, " started"))
	tmpLogs := make([]BufferLog, 0)
	for {
		select {
		case <-time.After(1 * time.Minute):
			if len(tmpLogs) > 0 {
				for _, l := range tmpLogs {
					logToDB(l.C1, l.C2, l.C3, l.Username, l.T64, l.Extra)
				}
				tmpLogs = make([]BufferLog, 0)
			}
		case <-shutdownChan:
			Info("close gm log worker:", i)
			for _, l := range tmpLogs {
				logToDB(l.C1, l.C2, l.C3, l.Username, l.T64, l.Extra)
			}

			wg.Done()
			return
		case l := <-gmlogChan:
			GMInfo(l.C1, FIELD_GAP, l.C2, FIELD_GAP, l.C3, FIELD_GAP, l.Username, FIELD_GAP, l.T, FIELD_GAP, l.Extra)
			if l.C3 != constants.C3_API_COST && l.C2 != constants.C2_CLIENT_LOG {
				if len(tmpLogs) >= TMP_SIZE {
					logToDB10(tmpLogs)
					tmpLogs = make([]BufferLog, 0)
				}

				tmpLogs = append(tmpLogs, l)
			}
		}
	}
}

//用这个方法记录日志必须精简 绝对不允许刷屏！！！！
func Info(v ...interface{}) {
	if !yace {
		sLog.Info(appcfg.GetServerType() + ":" + fmt.Sprint(v))
	}
}

//这个方法记录日志 正常只在非正式环境记录
func DebugInfo(v ...interface{}) {
	if logDebug {
		sLog.Info(appcfg.GetServerType() + ":" + fmt.Sprint(v))
	}
}

func GMInfo(v ...interface{}) {
	if !yace {
		formatLog.Info(fmt.Sprint(v))
	}
}

func Err(v ...interface{}) {
	msg := appcfg.GetServerType() + ":" + fmt.Sprint(v, "[stack]\n", string(debug.Stack()))
	sLog.Warn(msg)
	// fmt.Println(msg)
	if appcfg.GetBool("main_server", false) || appcfg.GetServerType() == constants.SERVER_TYPE_GM || appcfg.GetServerType() == constants.SERVER_TYPE_GV {
		go tools.SendInternalMail(appcfg.GetString("odps_env_name", "")+"_err_log_"+appcfg.GetAddress(), msg, nil)
	}
}

func ErrMail(v ...interface{}) {
	msg := appcfg.GetServerType() + ":" + fmt.Sprint(v, "[stack]\n", string(debug.Stack()))
	sLog.Warn(msg)
	// fmt.Println(msg)
	go tools.SendInternalMail(appcfg.GetString("odps_env_name", "")+"_err_log_"+appcfg.GetAddress(), msg, nil)
}

func OrderLog(channel string, v ...interface{}) {
	msg := fmt.Sprint(v)
	sLog.Info(appcfg.GetServerType() + ":" + msg)
	go tools.SendInternalMail(appcfg.GetString("odps_env_name", "")+"_"+channel+"_order_"+appcfg.GetAddress(), msg, nil)
}

func Debug(v ...interface{}) {
	if !yace {
		sLog.Debug(appcfg.GetServerType() + ":" + fmt.Sprint(v))
	}
}

func Profile(c1 string, c2 string, c3 string, username string, extra string) {
	if isProfiling {
		GMInfo(c1, FIELD_GAP, c2, FIELD_GAP, c3, FIELD_GAP, username, FIELD_GAP, time.Now().Format("2006-01-02 15:04:05"), FIELD_GAP, extra)
	}
}

func GMLog(c1 string, c2 string, c3 string, username string, extra string) {
	now := time.Now()
	clen := len(gmlogChan)
	if clen > appcfg.GetInt("gmlog_buffer", 1000)*9/10 {
		sLog.Info(appcfg.GetServerType() + ":" + fmt.Sprint("log channel size too big:", clen, c1, c2, c3))
		GMInfo(c1, FIELD_GAP, c2, FIELD_GAP, c3, FIELD_GAP, username, FIELD_GAP, now.Format("2006-01-02 15:04:05"), FIELD_GAP, extra)
		// go logToDB(c1, c2, c3, username, now.Unix(), extra)
	} else {
		gmlogChan <- BufferLog{C1: c1, C2: c2, C3: c3, Username: username, Extra: extra, T: now.Format("2006-01-02 15:04:05"), T64: now.Unix()}
	}

}

//GMLogKor 韩国版本日志插入
func GMLogKor(c1 string, c2 string, c3 string, username string, extra string) {
	now := time.Now()
	// Info(FIELD_GAP, c1, FIELD_GAP, c2, FIELD_GAP, c3, FIELD_GAP, username, FIELD_GAP, now.Format("2006-01-02 15:04:05"), FIELD_GAP, extra)
	logToDBKor(c1, c2, c3, username, now.Unix(), extra)
}

func logToDB(c1 string, c2 string, c3 string, username string, t int64, extra string) {
	if !yace {
		if appcfg.GetServerType() == "" {
			glog.AddLog(c1, c2, c3, username, t, extra)
		} else if appcfg.GetServerType() == constants.SERVER_TYPE_GM {
			// gmlog.AddLog(c1, c2, c3, username, t, extra)
		} else if appcfg.GetServerType() == constants.SERVER_TYPE_GV {
			// gvlog.AddLog(c1, c2, c3, username, t, extra)
		}
	}
}

func logToDB10(logs []BufferLog) {
	if !yace {
		if appcfg.GetServerType() == "" {
			glog.AddLog10(logs)
		} else if appcfg.GetServerType() == constants.SERVER_TYPE_GM {
			// for _, l := range logs {
			// 	// gmlog.AddLog(l.C1, l.C2, l.C3, l.Username, l.T64, l.Extra)
			// }
		}

	}
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
func korOrderlog(jsonStr string) {
	var OrderKorLogModule types.OrderKorLog
	if err := json.Unmarshal([]byte(jsonStr), &OrderKorLogModule); err != nil {
		Err(err)
		return
	}
	stmt := dao_korlog.AddKorOrderLogStmt
	_, err := stmt.Exec(OrderKorLogModule.Eventdate, OrderKorLogModule.Logtype, OrderKorLogModule.Guid, OrderKorLogModule.Platformuserid, OrderKorLogModule.Platformcode, OrderKorLogModule.Oscode, OrderKorLogModule.Storetype, OrderKorLogModule.Orderid, OrderKorLogModule.Paydate, OrderKorLogModule.Paytool, OrderKorLogModule.Gamecode, OrderKorLogModule.Gamename, OrderKorLogModule.Itemcode, OrderKorLogModule.Itemname, OrderKorLogModule.Currentcycode, OrderKorLogModule.Amount, OrderKorLogModule.Paycountry)
	if err != nil {
		Err(err)
		return
	}
}

//logKorAllLog 插入登录日志
func logKorAllLog(jsonStr string) {
	var loginModule types.LoginKorLog
	if err := json.Unmarshal([]byte(jsonStr), &loginModule); err != nil {
		Err(err)
		return

	}
	stmt := dao_korlog.AddKorLoginLogStmt
	_, err := stmt.Exec(loginModule.Guid, loginModule.Uuid, loginModule.Cuid, loginModule.Gamecode, loginModule.Nickname, loginModule.Sex, loginModule.Age, loginModule.LevYn, loginModule.UserLv, loginModule.Authtoken, loginModule.Oscode, loginModule.Ostopverno, loginModule.Tempuserid, loginModule.TokenredisDate, loginModule.Registdate, loginModule.Mappingdate)
	if err != nil {
		Err(err)
		return
	}
}

//korloginLog 插入会员登日志
func korLoginLog(jsonStr string) {
	var login1Module types.LoginKorLog1
	if err := json.Unmarshal([]byte(jsonStr), &login1Module); err != nil {
		Err(err)
		return

	}
	stmt := dao_korlog.AddKorLogin1LogStmt
	_, err := stmt.Exec(login1Module.Logindate, login1Module.Cuid, login1Module.Isloginsuccessdiv, login1Module.Platformuserid, login1Module.Loginrootcode, login1Module.GameCode, login1Module.Deviceid, login1Module.Platformcode, login1Module.Oscode, login1Module.Ostopverno)
	if err != nil {
		Err(err)
		return

	}
}
