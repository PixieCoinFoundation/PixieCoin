package guild

import (
	"appcfg"
	// "common"
	"common_db/guild"
	"constants"
	"dao_guild"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	. "logger"
	"math/rand"
	// "rao"
	"strconv"
	"strings"
	// "sync"
	// "sync/atomic"
	"time"
	. "types"
	. "zk_manager"
)

var ErrTransferGuild = errors.New("transfer guild owner error")
var ErrLockOpponent = errors.New("error lock opponent")
var MatchAllErr = errors.New("matched all guild")
var UpdateGuildWarErr = errors.New("not both guild war ok")
var UpdateGuildWarMemberErr = errors.New("not both guild war member ok")
var WarTokenWrong = errors.New("war token wrong")
var WarDateTokenNotMatch = errors.New("war date token not match")
var MemberIgnoreWarNotMatchErr = errors.New("del member ignore war not match")

var warMinMemberCnt int

func init() {
	warMinMemberCnt = appcfg.GetInt("guild_war_min_member", constants.GUILD_WAR_MEMBER_MIN)
}

func ListGuildVicePresidents(gid int64) (res []string, err error) {
	var rows *sql.Rows
	if rows, err = dao_guild.ListVPStmt.Query(gid); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]string, 0)
		for rows.Next() {
			var u string
			if err = rows.Scan(&u); err != nil {
				Err(err)
				return
			}

			res = append(res, u)
		}
	}
	return
}

func ConfigVicePresident(gid int64, username string, configType int) (success bool, err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	var res sql.Result
	var effect int64
	var vpc int
	if err = dao_guild.CountVicePresidentStmt.QueryRow(gid).Scan(&vpc); err != nil {
		Err(err)
		return
	} else if (configType == constants.GUILD_CONFIG_VICE_PRESIDENT_SET && vpc < constants.GUILD_VICE_PRESIDENT_SIZE) || (configType == constants.GUILD_CONFIG_VICE_PRESIDENT_UNSET && vpc > 0) {
		if configType == constants.GUILD_CONFIG_VICE_PRESIDENT_SET {
			res, err = dao_guild.ConfigGuildVicePresidentStmt.Exec(constants.GUILD_CONFIG_VICE_PRESIDENT_SET, gid, username, constants.GUILD_CONFIG_VICE_PRESIDENT_UNSET)
		} else if configType == constants.GUILD_CONFIG_VICE_PRESIDENT_UNSET {
			res, err = dao_guild.ConfigGuildVicePresidentStmt.Exec(constants.GUILD_CONFIG_VICE_PRESIDENT_UNSET, gid, username, constants.GUILD_CONFIG_VICE_PRESIDENT_SET)
		}

		if err != nil {
			Err(err)
			return
		} else if effect, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		} else if effect <= 0 {
			err = constants.GuildVicePresidentConfigEffectErr
			return
		}
	} else {
		err = constants.GuildVicePresidentSizeErr
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	} else {
		success = true
	}
	return
}

func DelGuildBoardMsg(id int64, username string) (err error) {
	if _, err = dao_guild.DelGuildBoardMsgStmt.Exec(id, username); err != nil {
		Err(err)
		return
	}
	return
}

func GetMaxBoardID(gid int64) (id int64) {
	if err := dao_guild.GetGuildMaxBoardIDStmt.QueryRow(gid).Scan(&id); err != nil {
		Err(err)
	}
	return
}

func GetGuildBoardMsgs(gid int64, pageNo, pageCnt int) (fid int64, res []GuildBoardMsg, err error) {
	start := (pageNo - 1) * pageCnt

	var rows *sql.Rows
	if rows, err = dao_guild.ListGuildBoardMsgStmt.Query(gid, start, pageCnt); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]GuildBoardMsg, 0)
		for rows.Next() {
			m := GuildBoardMsg{}
			if err = rows.Scan(&m.ID, &m.Author, &m.AuthorNickname, &m.Head, &m.ReplyTo, &m.Content, &m.Time); err != nil {
				Err(err)
				return
			}
			if fid == 0 {
				fid = m.ID
			}
			res = append(res, m)
		}
	}
	return
}

func AddGuildBoardMsg(gid int64, content, un, nn, head, rt string) (err error) {
	var res sql.Result
	var id int64
	if res, err = dao_guild.AddGuildBoardMsgStmt.Exec(gid, un, nn, head, rt, content, time.Now().Unix()); err != nil {
		Err(err)
		return
	} else if id, err = res.LastInsertId(); err != nil {
		Err(err)
		return
	} else {
		SetGuildBoardNewestIDR(gid, id)
	}
	return
}

func GetSingleWarLog() (fromGID int64, toGID int64, warDate string, err error) {
	if err = dao_guild.GetSingleWarLogStmt.QueryRow().Scan(&fromGID, &toGID, &warDate); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("no war log in db.")
		} else {
			Err(err)
		}
	}
	return
}

func GetSingleWarLogs(fromGID, toGID int64, warDate string) (res []WarLog, err error) {
	var rows *sql.Rows
	if rows, err = dao_guild.GetSingleWarLogsStmt.Query(fromGID, toGID, toGID, fromGID, warDate); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]WarLog, 0)
		var logs string
		for rows.Next() {
			m := WarLog{}
			if err = rows.Scan(&logs); err != nil {
				Err(err)
				return
			}
			if err = json.Unmarshal([]byte(logs), &m); err != nil {
				Err(err)
				return
			}

			res = append(res, m)
		}
	}
	return
}

func DelSingleWarlogs(fromGID, toGID int64, warDate string) (effect int64, err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	var res sql.Result
	if res, err = tx.Stmt(dao_guild.DelSingleWarlogsStmt).Exec(fromGID, toGID, toGID, fromGID, warDate); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	}

	if _, err = tx.Stmt(dao_guild.DelSetDefendLogStmt).Exec(fromGID, toGID, toGID, fromGID, warDate); err != nil {
		Err(err)
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}

	return
}

func GetDisbandedGuilds() (res []int64, err error) {
	var rows *sql.Rows
	if rows, err = dao_guild.GetDisbandedGuildStmt.Query(); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]int64, 0)
		for rows.Next() {
			var id int64
			if err = rows.Scan(&id); err != nil {
				Err(err)
				return
			}

			res = append(res, id)
		}
	}
	return
}

func DelDisbandedGuild(gid int64) (err error) {
	if _, err = dao_guild.DelDisbandedGuildStmt.Exec(gid); err != nil {
		Err(err)
		return
	}
	return
}

func UpdateNewSeasonMedal(gid int64, gwt, lgwt string) (ok bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_guild.UpdateNewSeasonMedalStmt.Exec(constants.GUILD_WAR_NEW_SEASON_MEDAL_MAX, gwt, gid, lgwt); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		ok = true
	} else {
		Err("update new season medal effect 0:", gid, gwt, lgwt)
	}
	return
}

func ListSGuild(start int, limit int) (res []SGuild, err error) {
	var rows *sql.Rows
	if rows, err = dao_guild.ListSGuildStmt.Query(start, limit); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]SGuild, 0)
		for rows.Next() {
			m := SGuild{}
			if err = rows.Scan(&m.ID, &m.MedalCnt, &m.GWT, &m.Zombie); err != nil {
				Err(err)
				return
			}

			res = append(res, m)
		}
	}
	return
}

func AddGuildWarSettleInfo(gwt string, content string) (err error) {
	if _, err = dao_guild.AddGuildWarSettleStmt.Exec(gwt, content); err != nil {
		Err(err)
	}
	return
}

func GuildWarSettled(gwt string) bool {
	var c int
	if err := dao_guild.QueryGuildWarSettleStmt.QueryRow(gwt).Scan(&c); err != nil {
		Err(err)
		return true
	}

	if c > 0 {
		return true
	} else {
		return false
	}
}

func TryLotteryUser(gid int64, username string, index int, week string) (ok bool, err error) {
	var stmt *sql.Stmt
	var standard int
	switch index {
	case 1:
		stmt = dao_guild.UpdateSLotteryUser1Stmt
		standard = constants.GUILD_MEMBER_LOTTERY1_ACTIVITY
	case 2:
		stmt = dao_guild.UpdateSLotteryUser2Stmt
		standard = constants.GUILD_MEMBER_LOTTERY2_ACTIVITY
	case 3:
		stmt = dao_guild.UpdateSLotteryUser3Stmt
		standard = constants.GUILD_MEMBER_LOTTERY3_ACTIVITY
	case 4:
		stmt = dao_guild.UpdateSLotteryUser4Stmt
		standard = constants.GUILD_MEMBER_LOTTERY4_ACTIVITY
	default:
		return
	}

	var res sql.Result
	var effect int64
	if res, err = stmt.Exec(gid, username, standard, week); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		ok = true
	}
	return
}

func GuildLottery(gid int64, index int, username, nickname string, clothes string, week string, ol string) (success bool, us []string, err error) {
	l := GuildLotteryInfo{
		LotteryUsername: username,
		LotteryNickname: nickname,
		LotteryClothes:  clothes,
		Week:            week,
	}
	lb, _ := json.Marshal(l)

	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	switch index {
	case 1:
		success, us, err = GuildLottery1(string(lb), tx, dao_guild.UpdateGuildLottery1Stmt, dao_guild.GetLotteryUser1Stmt, dao_guild.UpdateLotteryUser1Stmt, gid, constants.GUILD_LOTTERY1_ACTIVITY, constants.GUILD_MEMBER_LOTTERY1_ACTIVITY, week, ol)
	case 2:
		success, us, err = GuildLottery1(string(lb), tx, dao_guild.UpdateGuildLottery2Stmt, dao_guild.GetLotteryUser2Stmt, dao_guild.UpdateLotteryUser2Stmt, gid, constants.GUILD_LOTTERY2_ACTIVITY, constants.GUILD_MEMBER_LOTTERY2_ACTIVITY, week, ol)
	case 3:
		success, us, err = GuildLottery1(string(lb), tx, dao_guild.UpdateGuildLottery3Stmt, dao_guild.GetLotteryUser3Stmt, dao_guild.UpdateLotteryUser3Stmt, gid, constants.GUILD_LOTTERY3_ACTIVITY, constants.GUILD_MEMBER_LOTTERY3_ACTIVITY, week, ol)
	case 4:
		success, us, err = GuildLottery1(string(lb), tx, dao_guild.UpdateGuildLottery4Stmt, dao_guild.GetLotteryUser4Stmt, dao_guild.UpdateLotteryUser4Stmt, gid, constants.GUILD_LOTTERY4_ACTIVITY, constants.GUILD_MEMBER_LOTTERY4_ACTIVITY, week, ol)
	default:
		return
	}

	if !success {
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}

	return
}

func GuildLottery1(ginfo string, tx *sql.Tx, ugstmt, ggmstmt, ugmstmt *sql.Stmt, gid int64, guildActivity int, guildMemberActivity int, week string, oginfo string) (success bool, us []string, err error) {
	var res sql.Result
	var rows *sql.Rows
	var effect int64
	if res, err = tx.Stmt(ugstmt).Exec(ginfo, gid, week, oginfo, guildActivity); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect <= 0 {
		Info("guild lottery 1 effect 0:", gid, week, ginfo, oginfo)
		return
	} else if rows, err = tx.Stmt(ggmstmt).Query(gid, guildMemberActivity, week); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		us = make([]string, 0)

		for rows.Next() {
			var u string
			if err = rows.Scan(&u); err != nil {
				Err(err)
				return
			}

			us = append(us, u)
		}
	}

	if res, err = tx.Stmt(ugmstmt).Exec(gid, guildMemberActivity, week); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect == int64(len(us)) {
		success = true
	}
	return

}

func AddPK(fromGID int64, fromUsername, fromNickname string, toGID int64, toUsername, toNickname string, winMedal int, warDate string, addMedal int, gwt string, scs []GuildClothes, lgwt string) (success bool, err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	var res sql.Result
	var effect int64

	if addMedal > 0 {
		//update guild medal
		if res, err = tx.Stmt(dao_guild.AddGuildMedalStmt).Exec(gwt, lgwt, constants.GUILD_WAR_NEW_SEASON_MEDAL_MAX, addMedal, gwt, fromGID, warDate); err != nil {
			Err(err)
			return
		} else if effect, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		} else if effect <= 0 {
			Err("add pk update guild medal effect 0:", fromGID, fromUsername, toGID, toUsername, winMedal, warDate, addMedal)
			return
		}
	}

	//update attack
	if res, err = tx.Stmt(dao_guild.UpdateMemberAttackStmt).Exec(winMedal >= 2, winMedal < 2, fromGID, fromUsername, warDate); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect <= 0 {
		Err("add pk update attack effect 0:", fromGID, fromUsername, toGID, toUsername, winMedal, warDate, addMedal)
		return
	}

	//update defend
	nlm := constants.GUILD_WAR_MEDAL_CNT - winMedal
	if res, err = tx.Stmt(dao_guild.UpdateMemberDefendStmt).Exec(nlm, nlm, winMedal < 2, winMedal >= 2, toGID, toUsername, warDate); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect <= 0 {
		Err("add pk update defend effect 0:", fromGID, fromUsername, toGID, toUsername, winMedal, warDate, addMedal)
		return
	}

	//add log
	wl := WarLog{
		FromGID:      fromGID,
		ToGID:        toGID,
		FromUsername: fromUsername,
		FromNickname: fromNickname,
		ToUsername:   toUsername,
		ToNickname:   toNickname,
		WinMedal:     winMedal,
		Time:         time.Now().Unix(),
	}
	wlb, _ := json.Marshal(wl)
	if _, err = tx.Stmt(dao_guild.AddGuildWarLogStmt).Exec(fromGID, toGID, warDate, string(wlb)); err != nil {
		Err(err)
		return
	}

	if len(scs) > 0 {
		if err = AddClothesShareCnt(tx, fromGID, scs); err != nil {
			return
		}
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	} else {
		success = true
	}
	return
}

func QueryWarLogs(gid int64, pageNo int, pageSize int) (res []WarLog, err error) {
	start := (pageNo - 1) * pageSize
	var rows *sql.Rows
	if rows, err = dao_guild.QueryWarLogStmt.Query(gid, gid, start, pageSize); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]WarLog, 0)
		var logs string
		var firstWD string
		var wd string
		for rows.Next() {
			m := WarLog{}

			if err = rows.Scan(&logs, &wd); err != nil {
				Err(err)
				return
			}

			if firstWD == "" {
				firstWD = wd
			} else if wd != firstWD {
				return
			}

			if err = json.Unmarshal([]byte(logs), &m); err != nil {
				Err(err)
				return
			}

			res = append(res, m)
		}
	}
	return
}

func AddMemberActivity(gid int64, username string, activity int, wt string) (success bool, err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	// if _, err = tx.Stmt(dao_guild.RefreshGuildMemberActivityStmt).Exec(wt, gid, username, wt); err != nil {
	// 	Err(err)
	// 	return
	// }

	var res sql.Result
	var effect int64
	if res, err = tx.Stmt(dao_guild.AddGuildMemberActivityStmt).Exec(wt, activity, wt, wt, wt, wt, wt, wt, gid, username); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		// if _, err = tx.Stmt(dao_guild.RefreshGuildActivityStmt).Exec(wt, gid, wt); err != nil {
		// 	Err(err)
		// 	return
		// }

		if res, err = tx.Stmt(dao_guild.AddGuildActivityStmt).Exec(wt, activity, wt, wt, wt, wt, wt, gid); err != nil {
			Err(err)
			return
		} else if effect, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		} else if effect > 0 {
			success = true
		} else {
			Info("add guild activity effect 0:", gid, username, activity, wt)
			return
		}
	} else {
		Info("add member activity effect 0:", gid, username, activity, wt)
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func GetDefendSize(gid int64, wd string) (size int, doneSize int, err error) {
	if err = dao_guild.GetDefendSizeStmt.QueryRow(gid, wd).Scan(&size, &doneSize); err != nil {
		Err(err)
	}
	return
}

func GetGuildWeekRewardClothes(token string) (res []string, err error) {
	var v string
	if err = dao_guild.GetGuildWeekInfoStmt.QueryRow(token).Scan(&v); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("no guild week info:", token)
		} else {
			Err(err)
		}
		return
	}

	if err = json.Unmarshal([]byte(v), &res); err != nil {
		Err(err)
	}

	return
}

func SetGuildWeekRewardClothes(token string, r []string) (err error) {
	rb, _ := json.Marshal(r)
	if _, err = dao_guild.AddGuildWeekInfoStmt.Exec(token, string(rb)); err != nil {
		if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
			Info("duplicate guild week info:", token)
		} else {
			Err(err)
		}
	}
	return
}

func GetGuildWeekPKSubjects(token string) (res []string, err error) {
	rt := "pk" + token
	var v string
	if err = dao_guild.GetGuildWeekInfoStmt.QueryRow(rt).Scan(&v); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("no guild week info:", rt)
		} else {
			Err(err)
		}
		return
	}

	if err = json.Unmarshal([]byte(v), &res); err != nil {
		Err(err)
	}

	return
}

func SetGuildWeekPKSubjects(token string, r []string) (err error) {
	rt := "pk" + token
	rb, _ := json.Marshal(r)
	if _, err = dao_guild.AddGuildWeekInfoStmt.Exec(rt, string(rb)); err != nil {
		if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
			Info("duplicate guild week info:", rt)
		} else {
			Err(err)
		}
	}
	return
}

func SetDefendClothes(gid, tgid int64, username string, clothes string, wd string, scs []GuildClothes) (success bool, err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	var res sql.Result
	var effect int64
	if res, err = tx.Stmt(dao_guild.SetDefendClothesStmt).Exec(clothes, gid, username, wd); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		scsb, _ := json.Marshal(scs)
		if _, err = tx.Stmt(dao_guild.AddSetDefendLogStmt).Exec(gid, tgid, wd, username, string(scsb), string(scsb)); err != nil {
			Err(err)
			return
		}

		// if err = addClothesShareCnt(-1, tx, gid, ncs); err != nil {
		// 	return
		// }

		// if err = addClothesShareCnt(1, tx, gid, scs); err != nil {
		// 	return
		// }
	} else {
		Info("set defend effect 0:", gid, username, clothes, wd)
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	success = true
	return
}

// func GetShareDefendClothes(fromGID, toGID int64, username string, warDate string) (res []GuildClothes, err error) {
// 	res = make([]GuildClothes, 0)
// 	var cs string
// 	if err = dao_guild.QueryUserShareDefendClothesStmt.QueryRow(fromGID, toGID, username, warDate).Scan(&cs); err != nil {
// 		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
// 			Info("no share clothes used in defend:", fromGID, toGID, username, warDate)
// 		} else {
// 			Err(err)
// 		}
// 		return
// 	} else if err = json.Unmarshal([]byte(cs), &res); err != nil {
// 		Err(err)
// 		return
// 	}
// 	return
// }

func QuerySetDefendUsers(fromGID, toGID int64, warDate string) (res []SetDefendInfo, err error) {
	var rows *sql.Rows
	if rows, err = dao_guild.QuerySetDefendUserStmt.Query(fromGID, toGID, warDate); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]SetDefendInfo, 0)
		for rows.Next() {
			var u SetDefendInfo
			var cs string
			if err = rows.Scan(&u.Username, &cs); err != nil {
				Err(err)
				return
			}

			json.Unmarshal([]byte(cs), &u.ShareClothesUsed)

			res = append(res, u)
		}
	}
	return
}

func MatchGuildWar(fromGid int64, fromMedal int, wd string, pks []string, gwt string) (opponent Guild, err error) {
	for i := 0; i < 3; i++ {
		if opponent, err = matchGuildWarOnce(fromGid, fromMedal, wd, pks, gwt); err == nil && opponent.ID > 0 {
			return
		} else if err == MatchAllErr {
			break
		}
	}
	// Err(MatchCntErr)
	err = constants.MatchCntErr
	return
	// return
}

func matchGuildWarOnce(fromGid int64, fromMedal int, wd string, pks []string, gwt string) (opponent Guild, err error) {
	gap := 10

	var rows *sql.Rows
	var qc int
	for {
		qc++
		start := fromMedal - gap
		if start < 0 {
			start = 0
		}
		end := fromMedal + gap

		if rows, err = dao_guild.MatchGuildWarStmt.Query(start, end, wd, fromGid, warMinMemberCnt); err != nil {
			Err(err)
			return
		} else {
			defer rows.Close()
			opponents := make([]Guild, 0)
			for rows.Next() {
				var m Guild
				if err = rows.Scan(&m.ID, &m.Owner, &m.OwnerNickname, &m.Name, &m.Desc, &m.LogoType, &m.MemberCnt, &m.WarMemberCnt); err != nil {
					Err(err)
					rows.Close()
					return
				}
				opponents = append(opponents, m)
			}
			rows.Close()

			length := len(opponents)

			checkSize := 5
			if !appcfg.GetBool("guild_war_check_time", true) {
				checkSize = 1
			}

			if length >= checkSize {
				opponent = opponents[rand.Intn(length)]
				Info("find opponent:", fromGid, fromMedal, wd, opponent.ID, "query cnt:", qc)
				break
			} else if end > constants.GUILD_WAR_MAX_MEDAL && start <= 0 {
				Info("already find all guilds but no opponent:", fromGid, fromMedal, wd, "query cnt:", qc)
				err = MatchAllErr
				break
			} else {
				// if qc <= 3 {
				// 	gap += 5
				// } else {
				gap = gap * 2
				// }
			}
		}
	}

	if err == nil && opponent.ID > 0 {
		//init war
		if lock, path := LockGuild(opponent.ID); lock {
			defer Unlock(path)

			err = initWar(fromGid, opponent.ID, wd, pks)
		} else {
			err = ErrLockOpponent
		}
	}

	return
}

func initWar(fromGid int64, toGid int64, wd string, pks []string) (err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	if udpateGuildWar(fromGid, toGid, wd, tx) {
		if udpateGuildWarMember(fromGid, toGid, wd, tx, pks) {
			//ok
		} else {
			Err(UpdateGuildWarMemberErr, fromGid, toGid, wd)
			err = UpdateGuildWarMemberErr
			return
		}
	} else {
		Info(UpdateGuildWarErr, fromGid, toGid, wd)
		err = UpdateGuildWarErr
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func udpateGuildWar(fgid, tgid int64, wd string, tx *sql.Tx) (success bool) {
	var res sql.Result
	var effect int64
	var err error

	wt := fmt.Sprintf("%d_%d", fgid, tgid)
	if res, err = tx.Stmt(dao_guild.UpdateGuildWarStmt).Exec(wd, wt, fgid, tgid, wd); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect >= 2 {
		success = true
	}

	if res, err = tx.Stmt(dao_guild.ClearGuildZombieStmt).Exec(fgid, tgid); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		Info("delete zombie in gid:", fgid, tgid, "cnt:", effect)
	}

	return
}

func udpateGuildWarMember(fgid, tgid int64, wd string, tx *sql.Tx, pks []string) (success bool) {
	var res sql.Result
	var effect int64
	var err error

	if res, err = tx.Stmt(dao_guild.UpdateGuildMemberWarStmt).Exec(
		0, constants.GUILD_WAR_ATTACK_CNT, "", constants.GUILD_WAR_MEDAL_CNT, wd,
		pks[0], pks[1], pks[2], pks[3], pks[4], pks[5], pks[6], pks[7], pks[8], pks[9], pks[10], pks[11], pks[12], pks[13], pks[14],
		fgid, tgid); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		success = true
	}

	return
}

func ListApply(start, limit int) (res []*SApply, err error) {
	var rows *sql.Rows
	if rows, err = dao_guild.ListApplyGuildStmt.Query(start, limit); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]*SApply, 0)
		for rows.Next() {
			m := SApply{}
			if err = rows.Scan(&m.ID, &m.ApplyTime); err != nil {
				Err(err)
				return
			}
			res = append(res, &m)
		}
	}
	return
}

func AddClothesShareCnt(tx *sql.Tx, gid int64, cs []GuildClothes) (err error) {
	// var res sql.Result
	// var effect int64

	for _, c := range cs {
		if _, err = tx.Stmt(dao_guild.UpdateGuildClothesShareStmt).Exec(1, c.ClothesID, c.Username, gid); err != nil {
			Err(err)
			return
		}

		if _, err = tx.Stmt(dao_guild.UpdateGuildMemberShareStmt).Exec(1, c.Username, gid); err != nil {
			Err(err)
			return
		}
	}

	return
}

func IsOwner(gid int64, username string) (own bool, err error) {
	var v int
	if err = dao_guild.CheckGuildOwnerStmt.QueryRow(gid, username).Scan(&v); err != nil {
		Err(err)
		return
	} else if v > 0 {
		own = true
	}
	return
}

func IsVP(gid int64, username string) (vp bool, err error) {
	var v int
	if err = dao_guild.CheckVPStmt.QueryRow(gid, username).Scan(&v); err != nil {
		Err(err)
		return
	} else if v > 0 {
		vp = true
	}
	return
}

func UpdateDesc(gid int64, desc, owner string, logo int) (err error) {
	if _, err = dao_guild.UpdateGuildDescStmt.Exec(desc, logo, gid, owner); err != nil {
		Err(err)
	}
	return
}

// func DelClothes(gid int64, username, clothes string) (err error) {
// 	if _, err = dao_guild.DelGuildClothesStmt.Exec(gid, username, clothes); err != nil {
// 		Err(err)
// 	}
// 	return
// }

func AddClothes(gid int64, username, nickname, head, clothes string) (err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	token := time.Now().Format("20060102150405")
	var addSuccess bool
	if addSuccess, err = addMemberClothesCnt(gid, username, 1, tx); addSuccess {
		if _, err = tx.Stmt(dao_guild.AddGuildClothesStmt).Exec(gid, username, nickname, head, clothes, token); err != nil {
			if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
				Info("guild already share clothes:", gid, username, clothes)
			} else {
				Err(err)
			}
			return
		}
	} else if err != nil {
		return
	} else {
		return constants.GuildClothesSizeErr
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func DelClothes(gid int64, username, clothes string) (err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	var delSuccess bool
	if delSuccess, err = addMemberClothesCnt(gid, username, -1, tx); delSuccess {
		if _, err = tx.Stmt(dao_guild.DelGuildClothesStmt).Exec(gid, username, clothes); err != nil {
			Err(err)
			return
		}
	} else if err != nil {
		return
	} else {
		return constants.GuildClothesSizeErr
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func ListClothes(gid int64) (res []GuildClothes, err error) {
	var rows *sql.Rows
	// if username == "" {
	rows, err = dao_guild.ListGuildClothesStmt.Query(gid)
	// } else {
	// rows, err = dao_guild.ListGuildClothesByUserStmt.Query(gid, username)
	// }

	if err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]GuildClothes, 0)
		defer rows.Close()
		for rows.Next() {
			m := GuildClothes{}
			if err = rows.Scan(&m.Username, &m.Nickname, &m.Head, &m.ClothesID, &m.UseCnt, &m.Token); err != nil {
				Err(err)
				return
			}
			res = append(res, m)
		}
	}
	return
}

func transfer(gid int64, fromU, toU, toN string, tx *sql.Tx) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = tx.Stmt(dao_guild.TransferGuildStmt).Exec(toU, toN, gid, fromU); err != nil {
		Err(err)
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
	} else if effect > 0 {
		success = true
	}
	return
}

func Add(owner, ownerNickname, name, desc string, logoType int, head string, vip, level, activity int, gwt string) (id int64, err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	var res sql.Result
	if res, err = tx.Stmt(dao_guild.AddGuildStmt).Exec(owner, ownerNickname, name, desc, logoType, gwt); err != nil {
		if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
			Info("duplicate guild name:", name)
		} else {
			Err(err)
		}
		return
	} else if id, err = res.LastInsertId(); err != nil {
		Err(err)
		return
	}

	if _, err = tx.Stmt(dao_guild.AddGuildBelongStmt).Exec(owner, id); err != nil {
		if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
			Info("player already in guild when create guild:", owner)
			err = constants.AlreadyInGuildErr
			return
		} else {
			Err(err)
			return
		}
	}

	if _, err = tx.Stmt(dao_guild.AddGuildMemberStmt).Exec(id, owner, ownerNickname, head, vip, level, activity); err != nil {
		Err(err)
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	GMLog(constants.C1_SYSTEM, constants.C2_GUILD, constants.C3_CREATE_GUILD, "SYSTEM", fmt.Sprintf("%d", id))
	return
}

func GetUserGuild(username string, wt string) (m GuildMember, err error) {
	var id int64
	var cc, vp int
	if err = dao_guild.GetUserGuildStmt.QueryRow(username).Scan(&id, &m.GID, &m.Username, &m.Nickname, &m.Head, &m.VIP, &m.Level, &m.Activity, &m.IgnoreWar, &m.IgnoreWarNotice, &cc, &m.ClothesShareCnt, &m.LeftAttackCnt, &m.WarClothes, &m.WarSubject, &m.LeftMedalCnt, &m.WarDate, &m.DefendCnt, &m.WinCnt, &m.LoseCnt, &m.ActivityWeek, &m.L1Got, &m.L2Got, &m.L3Got, &m.L4Got, &m.Zombie, &vp); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			// Info("user not in guild:", username)
			err = nil
		} else {
			Err(err)
		}
	} else {
		m.VicePresident = vp == 1
		if m.WarClothes != "" {
			m.WarClothesDone = true
		}

		if m.ActivityWeek != wt {
			m.Activity = 0
			m.L1Got = false
			m.L2Got = false
			m.L3Got = false
			m.L4Got = false
			m.ClothesShareCnt = 0
		}
	}
	return
}

func GetUserHisGuild(username string, gid int64, wt string) (m GuildMember, err error) {
	var id int64
	var cc, vp int
	if err = dao_guild.GetUserHisGuildStmt.QueryRow(username, gid).Scan(&id, &m.GID, &m.Username, &m.Nickname, &m.Head, &m.VIP, &m.Level, &m.Activity, &m.IgnoreWar, &m.IgnoreWarNotice, &cc, &m.ClothesShareCnt, &m.LeftAttackCnt, &m.WarClothes, &m.WarSubject, &m.LeftMedalCnt, &m.WarDate, &m.DefendCnt, &m.WinCnt, &m.LoseCnt, &m.ActivityWeek, &m.L1Got, &m.L2Got, &m.L3Got, &m.L4Got, &m.Zombie, &vp); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			// Info("user not in guild:", username)
			err = nil
		} else {
			Err(err)
		}
	} else {
		m.VicePresident = vp == 1
		if m.WarClothes != "" {
			m.WarClothesDone = true
		}

		if m.ActivityWeek != wt {
			m.Activity = 0
			m.L1Got = false
			m.L2Got = false
			m.L3Got = false
			m.L4Got = false
			m.ClothesShareCnt = 0
		}
	}
	return
}

func GetBothMembers(gid1 int64, gid2 int64, wt string) (res []GuildMember, err error) {
	var rows *sql.Rows
	if rows, err = dao_guild.ListBothGuildMemberStmt.Query(gid1, gid2); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		var id int64
		var cc, vp int
		res = make([]GuildMember, 0)
		for rows.Next() {
			m := GuildMember{}
			if err = rows.Scan(&id, &m.GID, &m.Username, &m.Nickname, &m.Head, &m.VIP, &m.Level, &m.Activity, &m.IgnoreWar, &m.IgnoreWarNotice, &cc, &m.ClothesShareCnt, &m.LeftAttackCnt, &m.WarClothes, &m.WarSubject, &m.LeftMedalCnt, &m.WarDate, &m.DefendCnt, &m.WinCnt, &m.LoseCnt, &m.ActivityWeek, &m.L1Got, &m.L2Got, &m.L3Got, &m.L4Got, &m.Zombie, &vp); err != nil {
				Err(err)
				return
			}
			m.VicePresident = vp == 1
			if m.WarClothes != "" {
				m.WarClothesDone = true
			}
			if m.ActivityWeek != wt {
				m.Activity = 0
				m.L1Got = false
				m.L2Got = false
				m.L3Got = false
				m.L4Got = false
			}
			res = append(res, m)
		}
	}
	return
}

func GetMembers(gid int64, ignore string, includeZombie bool, wt string) (res []GuildMember, err error) {
	var rows *sql.Rows
	if rows, err = dao_guild.ListGuildMemberStmt.Query(gid, ignore); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		var id int64
		var cc, vp int
		res = make([]GuildMember, 0)
		for rows.Next() {
			m := GuildMember{}
			if err = rows.Scan(&id, &m.GID, &m.Username, &m.Nickname, &m.Head, &m.VIP, &m.Level, &m.Activity, &m.IgnoreWar, &m.IgnoreWarNotice, &cc, &m.ClothesShareCnt, &m.LeftAttackCnt, &m.WarClothes, &m.WarSubject, &m.LeftMedalCnt, &m.WarDate, &m.DefendCnt, &m.WinCnt, &m.LoseCnt, &m.ActivityWeek, &m.L1Got, &m.L2Got, &m.L3Got, &m.L4Got, &m.Zombie, &vp); err != nil {
				Err(err)
				return
			}

			m.VicePresident = vp == 1

			if m.WarClothes != "" {
				m.WarClothesDone = true
			}

			if m.ActivityWeek != wt {
				m.Activity = 0
				m.L1Got = false
				m.L2Got = false
				m.L3Got = false
				m.L4Got = false
			}

			if includeZombie || !m.Zombie {
				res = append(res, m)
			}
		}
	}
	return
}

func UpdateMember(gid int64, username, head string, vip, level int) (err error) {
	if _, err = dao_guild.UpdateGuildMemberStmt.Exec(head, vip, level, gid, username); err != nil {
		Err(err)
	}
	return
}

func UpdateMemberWar(gid int64, username string, ignoreWar, ignoreWarNotice bool) (err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	if _, err = tx.Stmt(dao_guild.UpdateGuildMemberWarConfigStmt).Exec(ignoreWar, ignoreWarNotice, gid, username); err != nil {
		Err(err)
		return
	}

	var wmc int
	if err = tx.Stmt(dao_guild.QueryGuildWarMemberCntStmt).QueryRow(gid).Scan(&wmc); err != nil {
		Err(err)
		return
	}

	if _, err = tx.Stmt(dao_guild.UpdateGuildWarMemberCntStmt).Exec(wmc, gid); err != nil {
		Err(err)
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func GetNotFullGuilds(limit int, gwt, lgwt string) (res []Guild, err error) {
	var rows *sql.Rows
	if rows, err = dao_guild.GetNotFullGuildStmt.Query(gwt, lgwt, constants.GUILD_WAR_NEW_SEASON_MEDAL_MAX, constants.GUILD_MAX_MEMBER, limit); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]Guild, 0)
		for rows.Next() {
			m := Guild{}
			if err = rows.Scan(&m.ID, &m.Owner, &m.OwnerNickname, &m.Name, &m.Desc, &m.LogoType, &m.MemberCnt, &m.WarMemberCnt, &m.MedalCnt); err != nil {
				Err(err)
				return
			}
			res = append(res, m)
		}
	}
	return
}

func Disband(id int64) (err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}

	return guild.Disband(id, tx)
}

func Del(id int64) (err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	if _, err = tx.Exec(fmt.Sprintf("delete from gf_guild where id=%d", id)); err != nil {
		Err(err)
		return
	}

	if _, err = tx.Exec(fmt.Sprintf("delete from gf_apply_guild where gid=%d", id)); err != nil {
		Err(err)
		return
	}

	if _, err = tx.Exec(fmt.Sprintf("delete from gf_guild_member where gid=%d", id)); err != nil {
		Err(err)
		return
	}

	if _, err = tx.Exec(fmt.Sprintf("delete from gf_guild_clothes where gid=%d", id)); err != nil {
		Err(err)
		return
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	GMLog(constants.C1_SYSTEM, constants.C2_GUILD, constants.C3_DEL_GUILD, "SYSTEM", fmt.Sprintf("%d", id))
	return
}

func AddApply(gid int64, gname, username, nickname, head string, vip, level int) (err error) {
	if _, err = dao_guild.AddApplyGuildStmt.Exec(gid, gname, username, nickname, head, vip, level, time.Now().Unix()); err != nil {
		if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
			Info("already send apply:", username, gid)
		} else {
			Err(err)
		}
	}
	return
}

func DelApply(applyID int64) (err error) {
	if _, err = dao_guild.DelApplyGuildStmt.Exec(applyID); err != nil {
		Err(err)
	}
	return
}

func DelApplyByUser(username string) (err error) {
	if _, err = dao_guild.DelApplyGuildByUserStmt.Exec(username); err != nil {
		Err(err)
	}
	return
}

func GetApplys(gid int64) (res []GuildApply, err error) {
	var rows *sql.Rows
	if rows, err = dao_guild.GetAllApplyGuildStmt.Query(gid); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]GuildApply, 0)
		for rows.Next() {
			m := GuildApply{}
			if err = rows.Scan(&m.ID, &m.GID, &m.GName, &m.Username, &m.Nickname, &m.VIP, &m.Level, &m.Head, &m.ApplyTime); err != nil {
				Err(err)
				return
			}
			res = append(res, m)
		}
	}
	return
}

func GetOneApply(id int64) (m GuildApply, err error) {
	if err = dao_guild.GetOneApplyGuildStmt.QueryRow(id).Scan(&m.ID, &m.GID, &m.GName, &m.Username, &m.Nickname, &m.VIP, &m.Level, &m.Head, &m.ApplyTime); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("guild apply id not exist:", id)
		} else {
			Err(err)
		}
	}
	return
}

func addMemberCnt(gid int64, cnt int, tx *sql.Tx) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = tx.Stmt(dao_guild.AddMemberCntStmt).Exec(cnt, gid, cnt, constants.GUILD_MAX_MEMBER, cnt); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		success = true
	} else {
		Info("add member cnt failed:", gid, cnt)
	}
	return
}

func addWarMemberCnt(gid int64, cnt int, tx *sql.Tx) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = tx.Stmt(dao_guild.AddWarMemberCntStmt).Exec(cnt, gid, cnt, constants.GUILD_MAX_MEMBER, cnt); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		success = true
	} else {
		Info("add war member cnt failed:", gid, cnt)
	}
	return
}

func addMemberClothesCnt(gid int64, username string, cnt int, tx *sql.Tx) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = tx.Stmt(dao_guild.AddMemberClothesCntStmt).Exec(cnt, gid, username, cnt, cnt, constants.GUILD_MAX_MEMBER_CLOTHES); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		success = true
	}
	return
}

func AddMember(gid int64, username, nickname, head string, vip, level, activity int) (err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	if _, err = tx.Stmt(dao_guild.AddGuildBelongStmt).Exec(username, gid); err != nil {
		if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
			Info("player already in guild when add member:", username)
			return constants.AlreadyInGuildErr
		} else {
			Err(err)
			return
		}
	}

	var addSuccess bool
	if addSuccess, err = addMemberCnt(gid, 1, tx); addSuccess {
		if addSuccess, err = addWarMemberCnt(gid, 1, tx); addSuccess {
			if _, err = tx.Stmt(dao_guild.AddGuildMemberStmt).Exec(gid, username, nickname, head, vip, level, activity); err != nil {
				Err(err)
				return
			}
		} else if err != nil {
			return
		} else {
			return constants.GuildWarSizeErr
		}
	} else if err != nil {
		return
	} else {
		return constants.GuildSizeErr
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func DelMember(gid int64, username string, ignoreWar bool, toUsername string, toNickname string, quitGuild bool) (err error) {
	var tx *sql.Tx
	if tx, err = dao_guild.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	if quitGuild {
		if _, err = tx.Stmt(dao_guild.DelGuildBelongStmt).Exec(username, gid); err != nil {
			Err(err)
			return
		}
	}

	if toUsername != "" && toNickname != "" {
		var transferOK bool
		if transferOK, err = transfer(gid, username, toUsername, toNickname, tx); err != nil {
			return
		} else if !transferOK {
			err = ErrTransferGuild
			return
		}
	}

	if quitGuild {
		var delSuccess bool
		if delSuccess, err = addMemberCnt(gid, -1, tx); delSuccess {
			if !ignoreWar {
				//del war member cnt
				if delSuccess, err = addWarMemberCnt(gid, -1, tx); !delSuccess || err != nil {
					err = constants.GuildWarSizeErr
					return
				}
			}

			var res sql.Result
			var effect int64
			if res, err = tx.Stmt(dao_guild.DelGuildMemberWithWarStmt).Exec(gid, username, ignoreWar); err != nil {
				Err(err)
				return
			} else if effect, err = res.RowsAffected(); err != nil {
				Err(err)
				return
			} else if effect <= 0 {
				Err("del memer ignore war token not match:", gid, username, ignoreWar, toUsername, toNickname, quitGuild)
				err = MemberIgnoreWarNotMatchErr
				return
			}

			if _, err = tx.Stmt(dao_guild.DelGuildMemberClothesStmt).Exec(gid, username); err != nil {
				Err(err)
				return
			}
		} else if err != nil {
			return
		} else {
			return constants.GuildSizeErr
		}
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func GetMemberCnt(gid int64) (cnt int, err error) {
	if err = dao_guild.GetGuildMemberCntStmt.QueryRow(gid).Scan(&cnt); err != nil {
		Err(err)
	}
	return
}

func GetInfo(gid int64, wt string, gwt string, lgwt string) (id int64, owner string, ownerNickname string, name string, desc string, logoType int, memberCnt int, warMemberCnt int, medalCnt int, wd string, wo int64, activity int, lo1 GuildLotteryInfo, lo2 GuildLotteryInfo, lo3 GuildLotteryInfo, lo4 GuildLotteryInfo, zombie bool, l1 string, l2 string, l3 string, l4 string, err error) {
	var wat string
	if err = dao_guild.GetGuildInfoStmt.QueryRow(gwt, lgwt, constants.GUILD_WAR_NEW_SEASON_MEDAL_MAX, wt, wt, wt, wt, wt, gid).Scan(&id, &owner, &ownerNickname, &name, &desc, &logoType, &memberCnt, &warMemberCnt, &medalCnt, &wd, &wat, &activity, &l1, &l2, &l3, &l4, &zombie); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("no guild with id:", gid)
		} else {
			Err(err)
		}
	} else if wat != "" {
		wts := strings.Split(wat, "_")
		if len(wts) != 2 {
			err = WarTokenWrong
			return
		}
		oi1, _ := strconv.ParseInt(wts[0], 10, 64)
		if oi1 != id {
			wo = oi1
		} else {
			oi2, _ := strconv.ParseInt(wts[1], 10, 64)
			if oi2 != id {
				wo = oi2
			}
		}
	}

	if wd != "" && wo == 0 {
		Err("WarDateTokenNotMatch", wd, wo)
		err = WarDateTokenNotMatch
		return
	}

	json.Unmarshal([]byte(l1), &lo1)
	json.Unmarshal([]byte(l2), &lo2)
	json.Unmarshal([]byte(l3), &lo3)
	json.Unmarshal([]byte(l4), &lo4)

	if lo1.Week != wt {
		lo1 = GuildLotteryInfo{}
	}
	if lo2.Week != wt {
		lo2 = GuildLotteryInfo{}
	}
	if lo3.Week != wt {
		lo3 = GuildLotteryInfo{}
	}
	if lo4.Week != wt {
		lo4 = GuildLotteryInfo{}
	}

	return
}

func Search(key string, gwt string) (res []Guild, err error) {
	var rows *sql.Rows
	if rows, err = dao_guild.SearchGuildStmt.Query(gwt, key+"%"); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]Guild, 0)
		for rows.Next() {
			m := Guild{}
			if err = rows.Scan(&m.ID, &m.Owner, &m.OwnerNickname, &m.Name, &m.Desc, &m.LogoType, &m.MemberCnt, &m.WarMemberCnt, &m.MedalCnt); err != nil {
				Err(err)
				return
			}
			res = append(res, m)
		}
	}
	return
}
