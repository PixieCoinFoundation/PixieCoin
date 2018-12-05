package gmreport

import (
	"appcfg"
	"bytes"
	"common"
	"constants"
	"database/sql"
	"db/cronjob"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	. "logger"
	"math"
	"os/exec"
	"service/params"
	"strconv"
	"strings"
	"time"
	"tools"
)

const (
	DEVICE_DAU_SQL = `
		select 
			count(distinct username),
			count(distinct get_json_object(extra,'$.deviceId'))
  		from %s 
  		where ds='%s' and category_1='%s' and category_2='%s'
	`

	DEVICE_DAU_FULL_SQL = `
		select 
			count(distinct username),
			count(distinct get_json_object(extra,'$.deviceId'))
  		from %s 
  		where ds='%s' and category_1='%s' and category_2='%s' and category_3='%s'
	`

	SP_REPORT_NAME  = "SP"
	DAU_REPORT_NAME = "DAU"

	GAP4 = "+------------+------------+"

	ODPS_BIN = "./odpscmd/bin/odpscmd"
)

type DAU struct {
	//DAU
	Dau     float64
	UserDau float64
}

type SPDau struct {
	//cosplay dau
	CosDau     float64
	UserCosDau float64
	//pk dau
	PkDau     float64
	UserPkDau float64
	//designer dau
	DesignerDau     float64
	UserDesignerDau float64
	//gold lottery dau
	GlodLotteryDau     float64
	UserGoldLotteryDau float64
	//diamond lottery dau
	DiamondLotteryDau     float64
	UserDiamondLotteryDau float64
	//help dau
	HelpDau     float64
	UserHelpDau float64
	//party dau
	PartyDau     float64
	UserPartyDau float64
}

func StartGMReport() {
	var err error
	defer func() {
		if x := recover(); x != nil {
			tools.SendErrMail(constants.MAIL_CRON_ERR_SUBJECT, "gmreport error:"+x.(string)+","+appcfg.GetString("local_ip", "unknown"), nil, params.GetErrMailTos())
		} else if err != nil {
			tools.SendErrMail(constants.MAIL_CRON_ERR_SUBJECT, "gmreport error:"+err.Error()+","+appcfg.GetString("local_ip", "unknown"), nil, params.GetErrMailTos())
		}
	}()

	yday := time.Now().AddDate(0, 0, -1).Format("20060102")

	if cronjob.QueryCronJobProgress(common.GetGMLogLikeToken(yday)) >= appcfg.GetInt("gmlog_server_num", 0) {
		ProcessGMReport(yday)
	} else {
		err = constants.GMLogServerErr
	}
}

func ProcessGMReport(yday string) {
	Info("begin generate report..", time.Now())

	db_url := appcfg.GetString("db_ipport", "")
	db_username := appcfg.GetString("db_username", "")
	db_password := appcfg.GetString("db_password", "")
	db_name := appcfg.GetString("db_name", "")
	odps_table := appcfg.GetString("gm_log_odps_table_name", "")

	var db *sql.DB
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?readTimeout=15s&writeTimeout=15s&timeout=10s", db_username, db_password, db_url, db_name))
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	// yday := time.Now().AddDate(0, 0, -1).Format("20060102")

	Info("process date:", yday)

	res := SPDau{}
	dres := DAU{}

	//dau computed by distinct device id

	if dres.UserDau, dres.Dau, err = getCount(ODPS_BIN, fmt.Sprintf(DEVICE_DAU_SQL, odps_table, yday, "PLAYER", "LOGIN")); err != nil {
		Info("err occured when query dau.quit report.")
		return
	} else {
		Info("dau", dres.Dau, "userDau", dres.UserDau)
	}

	if dres.Dau > 0 && dres.UserDau > 0 {
		var userCosDau, cosDau float64
		if userCosDau, cosDau, err = getCount(ODPS_BIN, fmt.Sprintf(DEVICE_DAU_SQL, odps_table, yday, "SYSTEM", "COSPLAY")); err != nil {
			Info("err occured when query cosplay dau.quit report.")
			return
		} else {
			Info("cosplay dau", cosDau, "user cosplay dau", userCosDau)
			res.UserCosDau = round(userCosDau/dres.UserDau, 2) * 100
			res.CosDau = round(cosDau/dres.Dau, 2) * 100
		}

		var userPkDau, pkDau float64
		if userPkDau, pkDau, err = getCount(ODPS_BIN, fmt.Sprintf(DEVICE_DAU_SQL, odps_table, yday, "SYSTEM", "PK")); err != nil {
			Info("err occured when query cosplay dau.quit report.")
			return
		} else {
			Info("pk dau", pkDau, "user pk dau", userPkDau)
			res.PkDau = round(pkDau/dres.Dau, 2) * 100
			res.UserPkDau = round(userPkDau/dres.UserDau, 2) * 100
		}

		var userDesignerDau, designerDau float64
		if userDesignerDau, designerDau, err = getCount(ODPS_BIN, fmt.Sprintf(DEVICE_DAU_FULL_SQL, odps_table, yday, "SYSTEM", "COSPLAY", "DESIGNER_UPLOAD")); err != nil {
			Info("err occured when query cosplay dau.quit report.")
			return
		} else {
			Info("designer dau", designerDau, "user designer dau", userDesignerDau)
			res.DesignerDau = round(designerDau/dres.Dau, 2) * 100
			res.UserDesignerDau = round(userDesignerDau/dres.UserDau, 2) * 100
		}

		var userGoldLotteryDau, glodLotteryDau float64
		if userGoldLotteryDau, glodLotteryDau, err = getCount(ODPS_BIN, fmt.Sprintf(DEVICE_DAU_FULL_SQL, odps_table, yday, "SYSTEM", "LOTTERY", "GOLD")); err != nil {
			Info("err occured when query cosplay dau.quit report.")
			return
		} else {
			Info("gold lottery dau", glodLotteryDau, "user gold lottery dau:", userGoldLotteryDau)
			res.GlodLotteryDau = round(glodLotteryDau/dres.Dau, 2) * 100
			res.UserGoldLotteryDau = round(userGoldLotteryDau/dres.UserDau, 2) * 100
		}

		var userDiamondLotteryDau, diamondLotteryDau float64
		if userDiamondLotteryDau, diamondLotteryDau, err = getCount(ODPS_BIN, fmt.Sprintf(DEVICE_DAU_FULL_SQL, odps_table, yday, "SYSTEM", "LOTTERY", "DIAMOND")); err != nil {
			Info("err occured when query cosplay dau.quit report.")
			return
		} else {
			Info("diamond lottery dau", diamondLotteryDau, "user diamond lottery dau:", userDiamondLotteryDau)
			res.DiamondLotteryDau = round(diamondLotteryDau/dres.Dau, 2) * 100
			res.UserDiamondLotteryDau = round(userDiamondLotteryDau/dres.UserDau, 2) * 100
		}

		var userHelpDau, helpDau float64
		if userHelpDau, helpDau, err = getCount(ODPS_BIN, fmt.Sprintf(DEVICE_DAU_SQL, odps_table, yday, "SYSTEM", "HELP")); err != nil {
			Info("err occured when query cosplay dau.quit report.")
			return
		} else {
			Info("help dau", helpDau, "user help dau:", userHelpDau)
			res.HelpDau = round(helpDau/dres.Dau, 2) * 100
			res.UserHelpDau = round(userHelpDau/dres.UserDau, 2) * 100
		}

		var userPartyDau, partyDau float64
		if userPartyDau, partyDau, err = getCount(ODPS_BIN, fmt.Sprintf(DEVICE_DAU_SQL, odps_table, yday, "SYSTEM", "PARTY")); err != nil {
			Info("err occured when query party dau.quit report.")
			return
		} else {
			Info("party dau", partyDau, "user party dau:", userPartyDau)
			res.PartyDau = round(partyDau/dres.Dau, 2) * 100
			res.UserPartyDau = round(userPartyDau/dres.UserDau, 2) * 100
		}
	}

	//insert to db
	var data, ddata []byte
	if data, err = json.Marshal(res); err != nil {
		Err(err)
		return
	} else {
		if _, err = db.Exec("insert into system_report values(0,?,?,?) on duplicate key update report=?", yday, SP_REPORT_NAME, string(data), string(data)); err != nil {
			Err(err)
			return
		}

	}

	if ddata, err = json.Marshal(dres); err != nil {
		Err(err)
		return
	} else {
		if _, err = db.Exec("insert into system_report values(0,?,?,?) on duplicate key update report=?", yday, DAU_REPORT_NAME, string(ddata), string(ddata)); err != nil {
			Err(err)
			return
		}
	}

	Info("report success:", res, dres, time.Now())
}

func getCount(odpsPath string, sql string) (ures float64, res float64, err error) {
	var c string
	if c, err = execute("/bin/bash", odpsPath, "-e", sql); err != nil {
		Err(err)
		return
	} else {
		return formatResult(c)
	}
}

func getSqlResult(odpsPath string, sql string) (res string, err error) {
	if res, err = execute("/bin/bash", odpsPath, "-e", sql); err != nil {
		Err(err)
	}
	return
}

func execute(cmd string, args ...string) (res string, err error) {
	Info(cmd, args)

	var cmmd *exec.Cmd
	cmmd = exec.Command(cmd, args...)
	errOutput := bytes.NewBuffer(nil)
	cmmd.Stderr = errOutput

	resOutput := bytes.NewBuffer(nil)
	cmmd.Stdout = resOutput

	if err = cmmd.Run(); err != nil {
		Err(err, string(errOutput.Bytes()))
		return
	} else {
		res = string(resOutput.Bytes())
	}

	return
}

func formatResult(res string) (ur float64, r float64, err error) {
	Info(res)
	parts := strings.Split(res, GAP4)
	fparts := strings.Split(parts[2], "|")

	if ur, err = strconv.ParseFloat(strings.TrimSpace(fparts[1]), 32); err != nil {
		Err(err)
		return
	}

	if r, err = strconv.ParseFloat(strings.TrimSpace(fparts[2]), 32); err != nil {
		Err(err)
		return
	}

	return
}

func round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}
