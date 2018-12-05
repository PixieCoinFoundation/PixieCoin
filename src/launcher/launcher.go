package main

import (
	"appcfg"
	"constants"
	_ "event_dispatch"
	"fmt"
	"gmlogger"
	"gmreport"
	"log"
	"net/http"
	_ "net/http/pprof"
	. "servers"
	"service/params"
	_ "service_pixie/land"
	_ "service_pixie/recommend"
	_ "simplet"
	"time"
	"tools"
)

func main() {
	cnLocation, _ := time.LoadLocation("Asia/Chongqing")

	nt := time.Now()
	localNow := nt.Format("2006-01-02 15:04:05")
	cnNow := nt.In(cnLocation).Format("2006-01-02 15:04:05")

	if cnNow != localNow && appcfg.GetLanguage() != constants.KOR_LANGUAGE {
		panic("Not in china.but config not kor!")
	}

	fmt.Println("start server type:", appcfg.GetServerType())
	if appcfg.GetBool("test_mail", false) {
		as := []string{"README.txt"}
		tools.SendErrMail("test", "test", as, params.GetErrMailTos())
	}

	if appcfg.GetServerType() == "" || appcfg.GetServerType() == constants.SERVER_TYPE_GL {
		pa := appcfg.GetString("profile_address", "")
		if pa != "" {
			go func() {
				log.Println(http.ListenAndServe(pa, nil))
			}()
		}
	}

	if appcfg.GetServerType() == constants.SERVER_TYPE_GM {
		StartGMServer()
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_GV {
		StartGVServer()
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_GL {
		StartGLServer(false)
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_GL_TS {
		StartGLServer(true)
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_GMLOG {
		gmlogger.Process(appcfg.GetDay())
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_GMLOG_MANUAL {
		gmlogger.ProcessSimple(appcfg.GetOdpsEnvName(), appcfg.GetDay(), appcfg.GetFileName(), appcfg.GetLogPrefix())
	} else if appcfg.GetServerType() == "KOR_REPORT" {
		gmlogger.MailKorOdpsTable1(appcfg.GetDay())
	} else if appcfg.GetServerType() == "MAIL_ODPS_TABLE" {
		gmlogger.MailOdpsTable(appcfg.GetOdpsCmdName(), appcfg.GetOdpsTableName(), appcfg.GetDay(), appcfg.GetReceivers())
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_SNAPSHOT_PLAYER {
		// player_mem.SnapshotPlayers(appcfg.GetFileName(), false)
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_SIMPLE_REPORT {
		// gmreport.ProcessSimpleReport(, false)
		gmreport.Process(appcfg.GetDay(), "scripts/daily_simple_report.sql", true, false, "", "", "")
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_SIMPLE_REPORT_MANUAL {
		// gmreport.ProcessSimpleReport(, false)
		gmreport.Process(appcfg.GetDay(), appcfg.GetFileName(), true, false, appcfg.GetReceivers(), appcfg.GetOdpsEnvName(), appcfg.GetOdpsCmdName())
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_SIMPLE_REPORT_INIT {
		// gmreport.ProcessSimpleReport(, false)
		gmreport.Process(appcfg.GetDay(), "scripts/daily_simple_report_init.sql", false, false, "", "", "")
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_GEN_ORDER {
		// stime, _ := time.ParseInLocation("20060102", appcfg.GetDay(), time.Local)
		// etime, _ := time.ParseInLocation("20060102", appcfg.EndDay, time.Local)
		// iap.GenOrderLog(ORDER_SUCCEED, stime.Unix(), etime.Unix())
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_PERF {
		// gops.Processes(appcfg.Pid)
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_S_FILE_SERVER {
		StartSFileServer()
	} else if appcfg.GetServerType() == "INIT_TEST_ELASTIC_CUSTOM" {
		// InitAndTestCustomFrom(appcfg.GetString("init_elastic_db_address", ""), appcfg.GetString("init_elastic_db_username", ""), appcfg.GetString("init_elastic_db_passwd", ""), appcfg.GetString("init_elastic_db_name", ""), appcfg.GetString("init_elastic_db_tblname", ""), appcfg.GetString("init_elastic_url", ""))
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_SIMPLE_TEST {
		simpleTest()
	} else {
		StartGameServer()
	}
}

func simpleTest() {
	// fmt.Println("check kuaikan result:", iap.CheckKuaikanOrder1("15476505739017420828524151070"))
	fmt.Println("simple test...")
	// m := make(map[string]string)

	// m["DEFAULT"] = "http://yjdyc.biligame.com"
	// m["IOS"] = "http://itunes.apple.com/us/app/%E5%A6%96%E7%B2%BE%E7%9A%84%E8%A1%A3%E6%A9%B1/id1202821460?l=zh&ls=1&mt=8"
	// m["BILI"] = "http://pkgdl.biligame.net/yjdyc/yjdyc_v1.1.1.16_bili_503936.apk"

	// m["0"] = "http://itunes.apple.com/us/app/%E5%A6%96%E7%B2%BE%E7%9A%84%E8%A1%A3%E6%A9%B1/id1202821460?l=zh&ls=1&mt=8"
	// m["1"] = "http://pkgdl.biligame.net/yjdyc/yjdyc_v1.1.1.16_bili_503936.apk"

	// m["8"] = "https://sdtuis.papa91.com/ac/v/nVYppJ"
	// m["9"] = "http://ug.kkmh.com/yjdyc_1.1.1.16_kuaikan_signed.apk"
	// m["10"] = "http://p-r6.sosobook.cn/auto/app/720/20171023213941_20171023213906_apk.apk"
	// m["15"] = "http://tel.down.manhuaren.com/android/game/yjdyf/yiqi-yjdyf-17-10-24.apk"
	// m["13"] = "http://file.market.xiaomi.com/download/AppChannel/0124c95c26c824860118467beff5ee23254c0e52a/com.yjdyc.aoyou.mi.apk"
	// m["7"] = "http://meitu.forgame.com/yjdyc_v1.1.apk"
	// m["5"] = "http://downali.game.uc.cn/wm/15/31/yjdyc-release-uc_10253439_185813a8138f.apk"
	// m["3"] = "http://storedl1.nearme.com.cn/apk/201710/18/5231b903d31feb0c4f73a00cdbfa2839.apk"
	// m["16"] = "http://signd.bce.baidu-mgame.com/service/cloudapk_sign_online/31000/31260/31260_1508404100_BaiduApp_signed.apk"

	// mb, _ := json.Marshal(m)
	// fmt.Println(string(mb))
	// panic("end")
}
