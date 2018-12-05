package player_mem

import (
	"appcfg"
	"common"
	"constants"
	"cron"
	"dao"
	"db/event"
	"fmt"
	"gmlogger"
	. "logger"
	"service/params"
	"time"
	"tools"
)

const (
	BATCH_SIZE = 1000
	ODPS_BIN   = "./odpscmd/bin/odpscmd"
)

func init() {
	if appcfg.GetServerType() == "" && appcfg.GetBool("main_server", false) && !appcfg.GetBool("debug", false) && appcfg.GetInt("gmlog_server_num", 0) <= 0 {
		panic("parameter gmlog_server_num must > 0")
	}

	spec3 := "0, 1, 0, *, *, *"
	c := cron.New()

	if appcfg.GetServerType() == "" || appcfg.GetServerType() == constants.SERVER_TYPE_GM || appcfg.GetServerType() == constants.SERVER_TYPE_GL || appcfg.GetServerType() == constants.SERVER_TYPE_GV {
		if !appcfg.GetBool("debug", false) {
			spec := appcfg.GetString("log_cron", "0, 5, 0, *, *, *")

			//gm log filter
			if appcfg.GetBool("upload_gmlog", true) {
				Info("add cronjob:gmlogger", spec)
				c.AddFunc(spec, gmlogger.StartGMLogger)
			}

			if appcfg.GetServerType() == "" && appcfg.GetBool("main_server", false) {
				//truncate table
				Info("add cronjob:truncatelog")
				c.AddFunc(spec3, truncateLog)
			}

			if appcfg.GetServerType() == "" || appcfg.GetServerType() == constants.SERVER_TYPE_GL {
				Info("add cronjob:profileReport", spec3)
				c.AddFunc(spec3, addOdpsPartition)
			}

			c.Start()
		}
	}
}

func addOdpsPartition() {
	if appcfg.GetBool("main_server", false) && appcfg.GetBool("upload_gmlog", true) {
		//add partition
		yday := time.Now().AddDate(0, 0, -1).Format("20060102")
		odpsTableName := appcfg.GetString("gm_log_odps_table_name", "")
		if _, err := common.Execute("/bin/bash", ODPS_BIN, "-e", fmt.Sprintf("alter table %s add if not exists partition (ds='%s')", odpsTableName, yday)); err != nil {
			return
		}
	}
}

func truncateLog() {
	var err error
	defer func() {
		if x := recover(); x != nil {
			tools.SendErrMail(constants.MAIL_CRON_ERR_SUBJECT, "truncateLog error:"+x.(string)+","+appcfg.GetString("local_ip", "unknown"), nil, params.GetErrMailTos())
		} else if err != nil {
			tools.SendErrMail(constants.MAIL_CRON_ERR_SUBJECT, "truncateLog error:"+err.Error()+","+appcfg.GetString("local_ip", "unknown"), nil, params.GetErrMailTos())
		}
	}()

	if tools.GetDayTag() == 0 {
		if _, err = dao.TruncateLog2Stmt.Exec(); err != nil {
			Err(err)
			return
		} else {
			Info("truncate day tag 0")
		}
	} else if tools.GetDayTag() == 1 {
		if _, err = dao.TruncateLog1Stmt.Exec(); err != nil {
			Err(err)
			return
		} else {
			Info("truncate day tag 1")
		}
	}

	if _, err = dao.TruncateGMLogStmt.Exec(); err != nil {
		Err(err)
		return
	} else {
		Info("truncate gm log")
	}

	if _, err = dao.TruncateRealTimePlayerLog.Exec(); err != nil {
		Err(err)
		return
	} else {
		Info("truncate real time player cnt table")
	}

	if err = event.DeleteEventPartyLog(); err != nil {
		return
	}
}
