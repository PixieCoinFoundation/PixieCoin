package gmreport

import (
	"appcfg"
	"common"
	"db/cronjob"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	. "logger"
	"os"
	"service/guild"
	"service/params"
	"strings"
	"time"
	"tools"
)

func CronSimpleReport() {
	ds := time.Now().AddDate(0, 0, -1).Format("20060102")
	// ds7 := time.Now().AddDate(0, 0, -7).Format("20060102")
	// ds30 := time.Now().AddDate(0, 0, -30).Format("20060102")

	progress := cronjob.QueryCronJobProgress(common.GetGMLogLikeToken(ds))
	if progress >= appcfg.GetInt("gmlog_server_num", 0) {
		if progress != appcfg.GetInt("gmlog_server_num", 0) {
			tools.SendInternalMail(appcfg.GetString("odps_env_name", "")+"_gmlog_cnt_wrong_"+ds+"_"+appcfg.GetAddress(), "gmlog_server_num not perfect", nil)
		}
		ProcessSimpleReport(ds, true)
	} else {
		Err("last day's cron simple report failed!:", ds)
	}
}

func ProcessSimpleReport(ds string, cron bool) {
	if e := Process(ds, "scripts/daily_simple_report_init.sql", false, cron, "", "", ""); e == nil {
		Process(ds, "scripts/daily_simple_report.sql", true, cron, "", "", "")
	}
}

func Process(ds, path string, doreport, cron bool, rvs string, oen string, odpscmdName string) error {
	Info(ds, path, doreport, cron, rvs, oen, odpscmdName)
	t, _ := time.ParseInLocation("20060102", ds, time.Local)
	ds1 := t.AddDate(0, 0, -1).Format("20060102")
	ds2 := t.AddDate(0, 0, -2).Format("20060102")
	ds3 := t.AddDate(0, 0, -3).Format("20060102")
	ds4 := t.AddDate(0, 0, -4).Format("20060102")
	ds7 := t.AddDate(0, 0, -7).Format("20060102")
	ds15 := t.AddDate(0, 0, -15).Format("20060102")
	ds30 := t.AddDate(0, 0, -30).Format("20060102")
	Info("process simple report:", ds, ds1, ds3, ds7, ds15, ds30, path, doreport, cron, rvs, oen)

	var df *os.File
	fileName := "reports/report_" + ds + ".csv"
	var err error

	if doreport {
		if _, err = common.Execute("rm", "-rf", appcfg.GetProjectPathPrefix()+fileName); err != nil {
			return err
		}

		if _, err = common.Execute("touch", appcfg.GetProjectPathPrefix()+fileName); err != nil {
			return err
		}

		if df, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0666); err != nil {
			Err(err)
			return err
		} else {
			defer df.Close()
		}
	}

	var sql string
	var cb []byte

	if cb, err = ioutil.ReadFile(path); err != nil {
		Err(err)
		return err
	} else {
		sql = string(cb)
	}

	sql = strings.Replace(sql, "###ds###", ds, -1)
	sql = strings.Replace(sql, "###1ds###", ds1, -1)
	sql = strings.Replace(sql, "###2ds###", ds2, -1)
	sql = strings.Replace(sql, "###3ds###", ds3, -1)
	sql = strings.Replace(sql, "###4ds###", ds4, -1)
	sql = strings.Replace(sql, "###7ds###", ds7, -1)
	sql = strings.Replace(sql, "###15ds###", ds15, -1)
	sql = strings.Replace(sql, "###30ds###", ds30, -1)
	if oen != "" {
		sql = strings.Replace(sql, "#ENV_NAME#", oen, -1)
	} else {
		sql = strings.Replace(sql, "#ENV_NAME#", appcfg.GetString("odps_env_name", ""), -1)
	}

	sqls := strings.Split(sql, ";;")

	ob := ODPS_BIN
	if odpscmdName != "" {
		ob = "./odpscmd/bin/" + odpscmdName
	}

	for _, s := range sqls {
		if len(s) > 1 {
			if r, e := getSqlResult(ob, s); e != nil {
				return e
			} else if doreport {
				lines := strings.Split(r, "\n")

				c := 0
				for _, line := range lines {
					if strings.HasPrefix(line, "+") || len(line) < 2 {
						continue
					}

					c++
					nl := strings.TrimSpace(line[1 : len(line)-1])
					nl = strings.Replace(nl, " | ", ",", -1)

					if _, err = io.WriteString(df, nl+"\n"); err != nil {
						Err(err)
						return err
					}
				}

				if c > 0 {
					if _, err = io.WriteString(df, "\n"); err != nil {
						Err(err)
						return err
					}
				}
			}
		}
	}

	if cron && doreport {
		if _, err = io.WriteString(df, "guild medal rank snapshot\n"); err != nil {
			Err(err)
			return err
		}

		_, gwt := common.GetGuildYWTokenByTime(t, 0, 0)
		guildRank := guild.GetMedalRank(gwt)

		for _, g := range guildRank {
			if _, err = io.WriteString(df, fmt.Sprintf("id:%d name:%s medal:%d\n", g.ID, g.Name, g.MedalCnt)); err != nil {
				Err(err)
				return err
			}
		}
	}

	if doreport {
		var title string
		if oen == "" {
			title = appcfg.GetString("odps_env_name", "") + "_日常数据报表_" + ds + "_" + appcfg.GetAddress()
		} else {
			title = oen + "_日常数据报表_" + ds
		}

		content := "请查看附件"
		if cron {
			tools.SendErrMail(title, content, []string{fileName}, params.GetErrMailTos())
		} else {
			if rvs == "" {
				tools.SendInternalMail(title, content, []string{fileName})
			} else {
				tools.Send(title, content, rvs, []string{fileName})
			}
		}
	}
	return nil
}
