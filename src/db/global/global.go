package global

import (
	"appcfg"
	"constants"
	"dao"
	"database/sql"
	"fmt"
	. "logger"
	"rao"
	"strings"
	"time"
	"tools"
	// "xlsx"
	"io/ioutil"
	"os"

	"github.com/garyburd/redigo/redis"
)

const (
	HOUR_6_EXPIRE = 6 * 3600
	MONTH_EXPIRE  = 32 * 24 * 3600
)

func init() {
	if appcfg.GetBool("main_server", false) && appcfg.GetServerType() == "" {
		rcs := appcfg.GetString("cron_report_receiver", "")
		envName := appcfg.GetString("odps_env_name", "")
		if rcs != "" {
			Info("cron report receiver:", rcs)
			go loopCronReport(envName, rcs)
		}
	}
}

func GetQinmiPartyDone(u1, u2, ds string) (res bool, err error) {
	nu1, nu2 := tools.GetUsernameSeq(u1, u2)
	token := fmt.Sprintf("qinmi_party_%s_%s_%s", nu1, nu2, ds)

	res = CheckOnceJob(token)

	return
}

func TryInviteQinmiFriend(u1, u2, ds string) bool {
	nu1, nu2 := tools.GetUsernameSeq(u1, u2)
	token := fmt.Sprintf("qinmi_party_%s_%s_%s", nu1, nu2, ds)
	return DoOnceJob(token)
}

func AddClientLogType(id int, name string, t int64) {
	dao.AddClientLogTypeStmt.Exec(id, name, t, t)
}

func DoOnceJob(token string) bool {
	if _, err := dao.TryDoOnceJobStmt.Exec(token, time.Now().Unix()); err != nil {
		if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
			Info("once job already done:", token)
		} else {
			Err(err)
		}
		return false
	}

	return true
}

func CheckOnceJob(token string) bool {
	var v int
	if err := dao.CheckOnceJobStmt.QueryRow(token).Scan(&v); err != nil {
		Err(err)
		return true
	}

	return v > 0
}

func loopCronReport(env, tos string) {
	for {
		n := time.Now()
		ns := n.Format("2006_01_02_15")
		if npr, e := GetNewPlayerCntReport(n); e == nil {
			go tools.SendSimpleReport(env+"_new_player_"+ns, npr, tos)
		}

		if nor, e := GetNewOrderReport(n); e == nil {
			go tools.SendSimpleReport(env+"_cronorder_"+ns, nor, tos)
		}

		time.Sleep(2 * time.Hour)
	}
}

func AddOfficialClothesReport(cid string, username, contact, reason, evidence string) {
	if _, err := dao.ReportOfficialClothesStmt.Exec(cid, username, contact, reason, evidence); err != nil {
		Err(err)
	}
}

//GetCodeFromRedisQueue 连接数据库
func GetCodeFromRedisQueue() (code1, code2, code3 string, err1, err2, err3 error) {
	conn := rao.GetConn()
	defer conn.Close()
	//获取队列名称
	queName1 := rao.GetActivationCode1Key()
	conn.Send("SPOP", queName1)
	queName2 := rao.GetActivationCode2Key()
	conn.Send("SPOP", queName2)
	queName3 := rao.GetActivationCode3Key()
	conn.Send("SPOP", queName3)

	if err := conn.Flush(); err != nil {
		Err(err)
		err1 = err
		return
	}

	code1, err1 = redis.String(conn.Receive())
	if err1 != nil {
		if err1.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err1)
			return "", "", "", err1, nil, nil
		}
		err1 = nil

	}
	code2, err2 = redis.String(conn.Receive())
	if err2 != nil {
		if err1.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err2)
			return "", "", "", err2, nil, nil
		}
		err2 = nil

	}
	code3, err3 = redis.String(conn.Receive())
	if err3 != nil {
		if err3.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err3)
			return "", "", "", err3, nil, nil
		}
		err3 = nil

	}
	return
}

//GetCodeInitFromFile 初始化礼包激活码
func GetCodeInitFromFile(path1, path2, path3 string) error {
	conn := rao.GetConn()
	defer conn.Close()

	//TODO openFile init codeslice
	var game1CodeList, game2CodeList, game3CodeList []string
	var err error
	if len(path1) < 1 || len(path2) < 1 || len(path3) < 1 {
		Err(constants.PARAMS_CAN_NOT_EMPTY)
	}

	//解析数据
	game1CodeList = ParseFileToCode(path1)
	game2CodeList = ParseFileToCode(path2)
	game3CodeList = ParseFileToCode(path3)

	clearQueue1 := rao.GetActivationCode1Key()
	if err = initCode(game1CodeList, clearQueue1, dao.AddKorCode1Init, dao.AddKorCode110Init, &conn); err != nil {
		return err
	}

	clearQueue2 := rao.GetActivationCode2Key()
	if err = initCode(game2CodeList, clearQueue2, dao.AddKorCode2Init, dao.AddKorCode210Init, &conn); err != nil {
		return err
	}

	clearQueue3 := rao.GetActivationCode3Key()
	if err = initCode(game3CodeList, clearQueue3, dao.AddKorCode3Init, dao.AddKorCode310Init, &conn); err != nil {
		return err
	}

	return nil
}

func initCode(gameCodeList []string, redisKey string, stmt, stmt10 *sql.Stmt, conn *redis.Conn) error {
	i := 0
	cs := make([]string, 0)
	var err error
	for _, v := range gameCodeList {
		i++
		(*conn).Send("SADD", redisKey, v)
		cs = append(cs, v)

		if i%20 == 0 {
			Info("init code:", i, redisKey)
			if err = (*conn).Flush(); err != nil {
				Err(err)
				return err
			}

			if _, err = stmt10.Exec(cs[0], cs[1], cs[2], cs[3], cs[4], cs[5], cs[6], cs[7], cs[8], cs[9], cs[10], cs[11], cs[12], cs[13], cs[14], cs[15], cs[16], cs[17], cs[18], cs[19]); err != nil {
				Err(err)
				return err
			}

			cs = make([]string, 0)
		}
	}

	if len(cs) > 0 {
		for _, v := range cs {
			(*conn).Send("SADD", redisKey, v)

			if _, err = stmt.Exec(v); err != nil {
				Err(err)
				return err
			}
		}

		if err = (*conn).Flush(); err != nil {
			Err(err)
			return err
		}
	}

	return nil
}

//GetCodeInitClear 清除redis 表中的数据
func GetCodeInitClear() error {
	conn := rao.GetConn()
	defer conn.Close()

	clearQueue1 := rao.GetActivationCode1Key()
	conn.Send("DEL", clearQueue1)

	clearQueue2 := rao.GetActivationCode2Key()
	conn.Send("DEL", clearQueue2)

	clearQueue3 := rao.GetActivationCode3Key()
	conn.Send("DEL", clearQueue3)
	if err := conn.Flush(); err != nil {
		Err(err)
		return err
	}
	return nil
}

//ParseFileToCode 解析数据
func ParseFileToCode(path string) []string {
	// var game1CodeList []string
	// data, err := xlsx.OpenFile(path)
	// if err != nil {
	// 	Err(err)
	// 	return nil
	// }
	// for _, sheet := range data.Sheet {
	// 	for _, row := range sheet.Rows {
	// 		for i, cell := range row.Cells {
	// 			if i == 1 {
	// 				v, _ := cell.String()
	// 				game1CodeList = append(game1CodeList, v)
	// 			}
	// 		}
	// 	}
	// }
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	ldata := strings.Split(string(data), "\n")
	Info("len:", len(ldata), path)
	return ldata
}

//GetCodeLeftFromRedisQueue 获取剩余量
func GetCodeLeftFromRedisQueue() (num1, num2, num3 int64, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	//获取游戏激活码
	queName1 := rao.GetActivationCode1Key()
	conn.Send("SCARD", queName1)
	queName2 := rao.GetActivationCode2Key()
	conn.Send("SCARD", queName2)
	queName3 := rao.GetActivationCode3Key()
	conn.Send("SCARD", queName3)

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	if num1, err = redis.Int64(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		}
		err = nil
	}
	if num2, err = redis.Int64(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		}
		err = nil
	}
	if num3, err = redis.Int64(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		}
		err = nil
	}
	return
}

func AddNewPlayerCnt(channel string, t time.Time) {
	conn := rao.GetConn()
	defer conn.Close()

	ck := rao.GetHourChannelNewPlayerCntKey(t, channel)
	ack := rao.GetHourAllChannelNewPlayerCntKey(t)
	dnpk := rao.GetDayPlayerNewCntKey(t)

	conn.Send("INCR", ck)
	conn.Send("INCR", ack)
	conn.Send("INCR", dnpk)
	// conn.Send("SADD",dapk,username)

	conn.Send("EXPIRE", dnpk, MONTH_EXPIRE)
	// conn.Send("EXPIRE",dapk,MONTH_EXPIRE)
	conn.Send("EXPIRE", ck, HOUR_6_EXPIRE)
	conn.Send("EXPIRE", ack, HOUR_6_EXPIRE)

	if err := conn.Flush(); err != nil {
		Err(err)
	}
}

func AddOrderCnt(rmb float64, currenty string, channel string, t time.Time) {
	conn := rao.GetConn()
	defer conn.Close()

	ck := rao.GetHourChannelNewOrderCntKey(t, channel)
	ack := rao.GetHourAllChannelNewOrderCntKey(t)

	ak := rao.GetHourChannelNewOrderAmountKey(t, channel, currenty)
	aak := rao.GetHourAllChannelNewOrderAmountKey(t, currenty)

	dock, doak := rao.GetDayOrderCntAndAmountKey(t)

	conn.Send("INCR", ck)
	conn.Send("INCR", ack)
	conn.Send("INCRBY", ak, rmb)
	conn.Send("INCRBY", aak, rmb)
	conn.Send("INCR", dock)
	conn.Send("INCRBY", doak, rmb)

	conn.Send("EXPIRE", ck, HOUR_6_EXPIRE)
	conn.Send("EXPIRE", ack, HOUR_6_EXPIRE)
	conn.Send("EXPIRE", ak, HOUR_6_EXPIRE)
	conn.Send("EXPIRE", aak, HOUR_6_EXPIRE)
	conn.Send("EXPIRE", dock, MONTH_EXPIRE)
	conn.Send("EXPIRE", doak, MONTH_EXPIRE)

	if err := conn.Flush(); err != nil {
		Err(err)
	}
}

func GetNewOrderReport(nowt time.Time) (rpt string, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	nh := nowt
	lh := nh.Add(-1 * time.Hour)
	ldh := nh.Add(-24 * time.Hour)

	dock, doak := rao.GetDayOrderCntAndAmountKey(nh)
	ldock, ldoak := rao.GetDayOrderCntAndAmountKey(ldh)

	conn.Send("GET", dock)
	conn.Send("GET", ldock)
	conn.Send("GET", doak)
	conn.Send("GET", ldoak)

	for i := 2; i <= 16; i++ {
		channel := fmt.Sprintf("%d", i)

		conn.Send("GET", rao.GetHourChannelNewOrderCntKey(nh, channel))
		conn.Send("GET", rao.GetHourChannelNewOrderAmountKey(nh, channel, ""))
		conn.Send("GET", rao.GetHourChannelNewOrderCntKey(lh, channel))
		conn.Send("GET", rao.GetHourChannelNewOrderAmountKey(lh, channel, ""))
	}

	conn.Send("GET", rao.GetHourAllChannelNewOrderCntKey(nh))
	conn.Send("GET", rao.GetHourAllChannelNewOrderAmountKey(nh, ""))
	conn.Send("GET", rao.GetHourAllChannelNewOrderCntKey(lh))
	conn.Send("GET", rao.GetHourAllChannelNewOrderAmountKey(lh, ""))

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	var nc, lc, na, la int
	if nc, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		} else {
			err = nil
		}
	}
	if lc, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		} else {
			err = nil
		}
	}
	if na, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		} else {
			err = nil
		}
	}
	if la, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		} else {
			err = nil
		}
	}

	rpt += fmt.Sprintf("今日总订单数：%d,昨日总订单数：%d<br>", nc, lc)
	rpt += fmt.Sprintf("今日总订单额：%d,昨日总订单额：%d<br>", na, la)

	rpt += `<table>
			<thead>
				<tr>
					<th>渠道</th>
					<th>本小时新增订单数</th>
					<th>本小时新增交易额</th>
					<th>上小时新增订单数</th>
					<th>上小时新增交易额</th>
				</tr>
			</thead>
			<tbody>`

	for i := 2; i <= 16; i++ {
		if nc, err = redis.Int(conn.Receive()); err != nil {
			if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
				Err(err)
				return
			} else {
				err = nil
			}
		}
		if lc, err = redis.Int(conn.Receive()); err != nil {
			if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
				Err(err)
				return
			} else {
				err = nil
			}
		}
		if na, err = redis.Int(conn.Receive()); err != nil {
			if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
				Err(err)
				return
			} else {
				err = nil
			}
		}
		if la, err = redis.Int(conn.Receive()); err != nil {
			if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
				Err(err)
				return
			} else {
				err = nil
			}
		}
		channel := GetChannelName(fmt.Sprintf("%d", i))
		rpt += fmt.Sprintf("<tr><th>%s</th><th>%d</th><th>%d</th><th>%d</th><th>%d</th></tr>", channel, nc, lc, na, la)
	}

	if nc, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		} else {
			err = nil
		}
	}
	if lc, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		} else {
			err = nil
		}
	}
	if na, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		} else {
			err = nil
		}
	}
	if la, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		} else {
			err = nil
		}
	}

	rpt += fmt.Sprintf("<tr><th>all channel</th><th>%d</th><th>%d</th><th>%d</th><th>%d</th></tr></tbody></table>", nc, lc, na, la)

	return
}

func GetNewPlayerCntReport(nowt time.Time) (rpt string, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	nh := nowt
	lh := nh.Add(-1 * time.Hour)
	ldh := nh.Add(-24 * time.Hour)

	dnpk := rao.GetDayPlayerNewCntKey(nh)
	ldnpk := rao.GetDayPlayerNewCntKey(ldh)

	var nc, lc int

	conn.Send("GET", dnpk)
	conn.Send("GET", ldnpk)

	for i := 2; i <= 16; i++ {
		channel := fmt.Sprintf("%d", i)
		nck := rao.GetHourChannelNewPlayerCntKey(nh, channel)
		lck := rao.GetHourChannelNewPlayerCntKey(lh, channel)

		conn.Send("GET", nck)
		conn.Send("GET", lck)
	}

	ack := rao.GetHourAllChannelNewPlayerCntKey(nh)
	lack := rao.GetHourAllChannelNewPlayerCntKey(lh)
	conn.Send("GET", ack)
	conn.Send("GET", lack)

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	if nc, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		} else {
			err = nil
		}
	}
	if lc, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		} else {
			err = nil
		}
	}

	rpt += fmt.Sprintf("今日新增玩家：%d,昨日新增玩家：%d<br>", nc, lc)

	rpt += `<table>
			<thead>
				<tr>
					<th>渠道</th>
					<th>本小时新增玩家数</th>
					<th>上小时新增玩家数</th>
				</tr>
			</thead>
			<tbody>`

	for i := 2; i <= 16; i++ {
		if nc, err = redis.Int(conn.Receive()); err != nil {
			if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
				Err(err)
				return
			} else {
				err = nil
			}
		}
		if lc, err = redis.Int(conn.Receive()); err != nil {
			if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
				Err(err)
				return
			} else {
				err = nil
			}
		}
		channel := GetChannelName(fmt.Sprintf("%d", i))
		rpt += fmt.Sprintf("<tr><th>%s</th><th>%d</th><th>%d</th></tr>", channel, nc, lc)
	}

	if nc, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		} else {
			err = nil
		}
	}
	if lc, err = redis.Int(conn.Receive()); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
			return
		} else {
			err = nil
		}
	}

	rpt += fmt.Sprintf("<tr><th>all channel</th><th>%d</th><th>%d</th></tr></tbody></table>", nc, lc)

	return
}

func GetChannelName(channel string) string {
	switch channel {
	case "2":
		return "HUAWEI"
	case "3":
		return "OPPO"
	case "4":
		return "VIVO"
	case "5":
		return "UC"
	case "6":
		return "4399"
	case "7":
		return "MEITU"
	case "8":
		return "PAPA"
	case "9":
		return "KUAIKAN"
	case "10":
		return "BUKA"
	case "11":
		return "TOPTOP"
	case "12":
		return "SINA"
	case "13":
		return "XIAOMI"
	case "14":
		return "360"
	case "15":
		return "MANHUAREN"
	case "16":
		return "BAIDU"
	default:
		return channel
	}
}
