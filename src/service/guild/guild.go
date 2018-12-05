package guild

import (
	"appcfg"
	"constants"
	"db/guild"
	"errors"
	"fmt"
	. "language"
	. "logger"
	"math/rand"
	"service/clothes"
	"service/mails"
	"service/pushlist"
	"strconv"
	"sync"
	. "types"
)

var weekLock sync.RWMutex
var weekToken string
var weekRewadClothes []string
var weekPKSubjects []string

var PKSubjectSizeErr = errors.New("pk subject size wrong")

func IsVP(gid int64, username string) (vp bool) {
	vp, _ = guild.IsVP(gid, username)
	return
}

func ListGuildVicePresidents(gid int64) ([]string, error) {
	return guild.ListGuildVicePresidents(gid)
}

func ConfigGuildVicePresident(gid int64, username string, configType int) (bool, error) {
	return guild.ConfigVicePresident(gid, username, configType)
}

func GetMaxBoardID(gid int64) (id int64) {
	return guild.GetMaxBoardID(gid)
}

func GetGuildBoardNewestIDR(gid int64) (id int64) {
	return guild.GetGuildBoardNewestIDR(gid)
}

func SetGuildBoardNewestIDR(gid int64, id int64) {
	guild.SetGuildBoardNewestIDR(gid, id)
}

func DelGuildBoardMsg(id int64, username string) (err error) {
	return guild.DelGuildBoardMsg(id, username)
}

func GetGuildBoardMsgs(gid int64, pageNo, pageCnt int) (fid int64, res []GuildBoardMsg, err error) {
	return guild.GetGuildBoardMsgs(gid, pageNo, pageCnt)
}

func AddGuildBoardMsg(gid int64, content, un, nn, head, rt string) (err error) {
	return guild.AddGuildBoardMsg(gid, content, un, nn, head, rt)
}

func GetMedalRank(gwt string) []Guild {
	// rankLock.RLock()
	// defer rankLock.RUnlock()
	res, _ := guild.GetGuildMedalRankCache(gwt)
	return res
}

func GetActivityRank(wt string) []Guild {
	// rankLock.RLock()
	// defer rankLock.RUnlock()
	res, _ := guild.GetGuildActivityRankCache(wt)
	return res
}

func RefreshGuildRank(gid int64, medal int, activity int, wt string, gwt string) error {
	if err := guild.RefreshGuildMedalRank(gid, medal, gwt); err != nil {
		return err
	} else {
		return guild.RefreshGuildActivityRank(gid, activity, wt)
	}
}

func RefreshGuildMedalRank(gid int64, medal int, gwt string) error {
	return guild.RefreshGuildMedalRank(gid, medal, gwt)
}

func AddGuildMedalRank(gid int64, addMedal int, gwt string) error {
	return guild.AddGuildMedalR(gid, addMedal, gwt)
}

func RefreshGuildActivityRank(gid int64, activity int, wt string) error {
	return guild.RefreshGuildActivityRank(gid, activity, wt)
}

func AddGuildActivityRank(gid int64, activity int, wt string) error {
	return guild.AddGuildActivityR(gid, activity, wt)
}

func ReadNewApply(gid int64) error {
	return guild.ReadNewApply(gid)
}

func TryLotteryUser(gid int64, username string, index int, week string) (bool, error) {
	return guild.TryLotteryUser(gid, username, index, week)
}

func GuildLottery(gid int64, index int, username, nickname string, clothes string, week string, ol string) (bool, []string, error) {
	return guild.GuildLottery(gid, index, username, nickname, clothes, week, ol)
}

func AddPKData(fromGID int64, fromUsername, fromNickname string, toGID int64, toUsername, toNickname string, winMedal int, warDate string, addMedal int, gwt string, scs []GuildClothes, lgwt string) (bool, error) {
	return guild.AddPK(fromGID, fromUsername, fromNickname, toGID, toUsername, toNickname, winMedal, warDate, addMedal, gwt, scs, lgwt)
}

func GetWeekRewardClothes(wt string) []string {
	weekLock.RLock()

	if wt == weekToken {
		defer weekLock.RUnlock()
		return weekRewadClothes
	} else {
		weekLock.RUnlock()
		setWeekRewardClothes(false)
		return GetWeekRewardClothes(wt)
	}
}

// func GetShareClothesInDefend(fromGID, toGID int64, username, warDate string) ([]GuildClothes, error) {
// 	return guild.GetShareDefendClothes(fromGID, toGID, username, warDate)
// }

func GetWeekRandomRewardClothes(wt string, force35 bool) string {
	cls := GetWeekRewardClothes(wt)
	if force35 {
		ncls := make([]string, 0)
		for _, c := range cls {
			if clo := clothes.GetClothesById(c); clo != nil && clo.Star >= 3 {
				ncls = append(ncls, c)
			}
		}
		return ncls[rand.Intn(len(ncls))]
	} else {
		return cls[rand.Intn(len(cls))]
	}
}

func GetWeekPKSubjects(wt string) []string {
	weekLock.RLock()

	if wt == weekToken {
		defer weekLock.RUnlock()
		return weekPKSubjects
	} else {
		weekLock.RUnlock()
		setWeekPKSubjects(false)
		return GetWeekPKSubjects(wt)
	}
}

func GetWeekRandomPKSubject(wt string) string {
	cls := GetWeekPKSubjects(wt)
	return cls[rand.Intn(len(cls))]
}

func AddMemberActivity(gid int64, username string, activity int, wt string) (bool, error) {
	return guild.AddMemberActivity(gid, username, activity, wt)
}

func GetDefendSize(gid int64, wd string) (int, int, error) {
	return guild.GetDefendSize(gid, wd)
}

func SetDefendClothes(gid, tgid int64, username string, clothes string, wd string, newcs []GuildClothes) (bool, error) {
	return guild.SetDefendClothes(gid, tgid, username, clothes, wd, newcs)
}

func FindWarLogs(fromGid int64) (res []WarLog, err error) {
	return guild.QueryWarLogs(fromGid, 1, 100)
}

func MatchGuildWar(gid int64, medal int, wd string, wt string, gwt string) (opponent Guild, err error) {
	pks := GetWeekPKSubjects(wt)
	if len(pks) != constants.GUILD_WEEK_PK_SUBJECT_SIZE {
		Err(PKSubjectSizeErr)
		err = PKSubjectSizeErr
		return
	}
	return guild.MatchGuildWar(gid, medal, wd, pks, gwt)
}

func IsOwner(gid int64, username string) (bool, error) {
	return guild.IsOwner(gid, username)
}

func UpdateDesc(gid int64, desc, owner string, logo int) (err error) {
	return guild.UpdateDesc(gid, desc, owner, logo)
}

func UpdateMember(gid int64, username, head string, vip, level int) error {
	return guild.UpdateMember(gid, username, head, vip, level)
}

func UpdateMemberWar(gid int64, username string, ignoreWar, ignoreWarNotice bool) error {
	return guild.UpdateMemberWar(gid, username, ignoreWar, ignoreWarNotice)
}

func GetRandomGuilds(size int, gwt, lgwt string) (res []Guild, err error) {
	var gl []Guild
	if gl, err = guild.GetNotFullGuilds(100, gwt, lgwt); err != nil {
		return
	}

	if len(gl) <= size {
		res = gl
		return
	} else {
		res = make([]Guild, 0)

		for len(res) < size {
			length := len(gl)
			idx := rand.Intn(length)
			res = append(res, gl[idx])

			//swap
			tmp := gl[length-1]
			gl[length-1] = gl[idx]
			gl[idx] = tmp

			gl = gl[0 : length-1]
		}
	}
	return
}

func AddClothes(gid int64, username, nickname, head, clothes string) error {
	return guild.AddClothes(gid, username, nickname, head, clothes)
}

// func AddClothesShareCnt(gid int64, cs []string) error {
// 	for len(cs)<3{
// 		cs = append(cs,"")
// 	}
// 	return guild.AddClothesShareCnt(gid, cs)
// }

func DelClothes(gid int64, username, clothes string) error {
	return guild.DelClothes(gid, username, clothes)
}

func ListClothes(gid int64) ([]GuildClothes, error) {
	return guild.ListClothes(gid)
}

func ListMembers(gid int64, ignore string, includeZombie bool, wt string) ([]GuildMember, error) {
	return guild.GetMembers(gid, ignore, includeZombie, wt)
}

func GetBothMembers(gid1 int64, gid2 int64, wt string) ([]GuildMember, error) {
	return guild.GetBothMembers(gid1, gid2, wt)
}

func Transfer(gid int64, fromUsername, toUsername, toNickname string, ignoreWar bool, quitGuild bool) (err error) {
	return guild.DelMember(gid, fromUsername, ignoreWar, toUsername, toNickname, quitGuild)
}

func GetUserGuild(username, wt string) (mInfo GuildMember, err error) {
	return guild.GetUserGuild(username, wt)
}

func GetUserHisGuild(username string, gid int64, wt string) (mInfo GuildMember, err error) {
	return guild.GetUserHisGuild(username, gid, wt)
}

func CreateGuild(owner, ownerNickname, name, desc string, logoType int, head string, vip, level, activity int, gwt string) (int64, error) {
	return guild.Add(owner, ownerNickname, name, desc, logoType, head, vip, level, activity, gwt)
}

func DisbandGuild(id int64, wt, gwt string) error {
	if err := guild.Disband(id); err == nil {
		guild.DelGuildRank(id, wt, gwt)
		return nil
	} else {
		return err
	}

}

// func AddGuildMember(gid int64, username, nickname string) error {
// 	return guild.AddMember(gid, username, nickname)
// }

func DelMember(gid int64, gname string, username string, ignoreWar bool, sendMail bool) error {
	if err := guild.DelMember(gid, username, ignoreWar, "", "", true); err == nil {
		if sendMail {
			//send email
			mails.SendToOneT("", username, L("guild3"), fmt.Sprintf(L("guild4"), gname), 0, 0, "", true, constants.MAIL_TYPE_GUILD_INOUT)
		}

		return nil
	} else {
		return err
	}
}

func GetGuild(gid int64, simple bool, wt, gwt, lgwt string) (g Guild, err error) {
	if g.ID, g.Owner, g.OwnerNickname, g.Name, g.Desc, g.LogoType, g.MemberCnt, g.WarMemberCnt, g.MedalCnt, g.WarDate, g.WarOpponent, g.Activity, g.Lottery1, g.Lottery2, g.Lottery3, g.Lottery4, g.Zombie, g.Lottery1s, g.Lottery2s, g.Lottery3s, g.Lottery4s, err = guild.GetInfo(gid, wt, gwt, lgwt); err != nil {
		return
	}

	if !simple {
		g.NewApplyCnt, _ = guild.GetNewApplyCnt(gid)
		g.ActivityRank, _ = guild.GetGuildActivityRank(gid, wt)
		g.PKRank, _ = guild.GetGuildMedalRank(gid, gwt)
	}

	return
}

func AddApply(gid int64, gname, username, nickname, head string, vip, level int) (err error) {
	if err = guild.AddApply(gid, gname, username, nickname, head, vip, level); err != nil {
		return
	} else {
		guild.AddNewApplyCnt(gid, 1)
	}
	return
}

func DelApply(id int64) (err error) {
	return guild.DelApply(id)
}

func GetApplys(gid int64) ([]GuildApply, error) {
	return guild.GetApplys(gid)
}

func PassApply(id int64) (u string, un string, err error) {
	var ga GuildApply
	if ga, err = guild.GetOneApply(id); err != nil {
		return
	}

	if err = guild.AddMember(ga.GID, ga.Username, ga.Nickname, ga.Head, ga.VIP, ga.Level, 0); err != nil {
		if err == constants.GuildSizeErr {
			guild.DelApply(id)
		} else if err == constants.AlreadyInGuildErr {
			guild.DelApplyByUser(ga.Username)
		}
		return
	} else {
		u = ga.Username
		un = ga.Nickname
		//success
		guild.DelApplyByUser(ga.Username)

		//join chat group
		// ids := fmt.Sprintf("%d", ga.GID)
		// JoinChatGroup(ga.Username, ids)

		//send email
		mails.SendToOneT("", ga.Username, L("guild1"), fmt.Sprintf(L("guild2"), ga.GName), 0, 0, "", true, constants.MAIL_TYPE_GUILD_INOUT)

		//push
		if appcfg.GetLanguage() == "" {
			// push.PushGuild()
			pushlist.AddSinglePushJob(ga.Username, L("guild1"), fmt.Sprintf(L("guild2"), ga.GName), constants.PUSH_TYPE_GUILD)
		}
	}
	return
}

func RejectApply(id int64) (err error) {
	return guild.DelApply(id)
}

func Search(key string, wt string, gwt, lgwt string) (res []Guild, err error) {
	res = make([]Guild, 0)

	//try id first
	if gid, err1 := strconv.ParseInt(key, 10, 64); err1 == nil {
		if g, err1 := GetGuild(gid, true, wt, gwt, lgwt); err1 == nil {
			res = append(res, g)
		}
	}

	//try name
	var nr []Guild
	if nr, err = guild.Search(key, gwt); err != nil {
		return
	} else {
		res = append(res, nr...)
	}

	return
}
