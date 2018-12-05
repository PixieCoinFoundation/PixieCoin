// +build !darwin,!linux

package logger

import (
	"appcfg"
	"fmt"
	"log"
	"runtime/debug"
	// "service/bilibili"
	"strings"
	"time"
)

var (
	isProfiling bool
)

const (
	FIELD_GAP = "^_^"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	isProfiling = appcfg.GetBool("profile", false)
}

func Info(v ...interface{}) {
	// log.Println("[INFO]", v)
	log.Println(fmt.Sprint(v))
}

func Err(v ...interface{}) {
	// log.Println("[ERROR]", v)
	log.Println(fmt.Sprint(v, "\n[stack]", string(debug.Stack())))
}

func Debug(v ...interface{}) {
	// log.Println("[DEBUG]", v)
	log.Println(fmt.Sprint(v))
}

func Profile(v ...interface{}) {
	if isProfiling {
		log.Println(fmt.Sprint(v))
	}
}

func GMLog(c1 string, c2 string, c3 string, username string, extra string) {
	Info(FIELD_GAP, c1, FIELD_GAP, c2, FIELD_GAP, c3, FIELD_GAP, username, FIELD_GAP, time.Now().Format("2006-01-02 15:04:05"), FIELD_GAP, extra)
}

//GMLogKor 韩国版本日志插入
func GMLogKor(c1 string, c2 string, c3 string, username string, extra string) {
	now := time.Now()
	if len(gmlogChan) > appcfg.GetInt("gmlog_buffer", 1000)*9/10 {
		Info(FIELD_GAP, c1, FIELD_GAP, c2, FIELD_GAP, c3, FIELD_GAP, username, FIELD_GAP, now.Format("2006-01-02 15:04:05"), FIELD_GAP, extra)
		logToDBKor(c1, c2, c3, username, now.Unix(), extra)
	} else {
		gmlogChan <- BufferLog{C1: c1, C2: c2, C3: c3, Username: username, Extra: extra, T: now.Format("2006-01-02 15:04:05"), T64: now.Unix()}
	}
}

// func BLog(logType string, platform string, v ...interface{}) {
// 	var params []string
// 	params = append(params, logType)
// 	// bi站Server的id
// 	serverID := bilibili.GetServerID(platform)
// 	params = append(params, serverID)
// 	params = append(params, time.Now().Format("2006-01-02 15:04:05"))
// 	for _, str := range v {
// 		params = append(params, fmt.Sprint(str))
// 	}
// 	param := strings.Join(params, "|")

// 	log.Println(param)
// }
