package guild

import (
	"appcfg"
	"common"
	"database/sql"
	"db/guild"
	"service/clothes"
	// "service/pk"
	"constants"
	"dao_guild"
	"encoding/json"
	"fmt"
	. "language"
	. "logger"
	"math/rand"
	"service/mails"
	"shutdown"
	"sync"
	"time"
	. "types"
	"xlsx"
)

var guildPKSubjectList []string

var settleSeasonRunning bool
var settleSeasonLock sync.RWMutex

func init() {
	// if appcfg.GetServerType() != "" {
	// 	return
	// }
	// loadGuildPKSubjectList()
	// setWeekRewardClothes(true)
	// setWeekPKSubjects(true)

	// if appcfg.GetServerType() == "" && appcfg.GetBool("main_server", false) {
	// 	go loop()
	// 	go freqLoop()
	// }
}

func loop() {
	for {
		refreshRank()
		settleSeason()
		delDisbandGuilds()
		time.Sleep(60 * time.Second)
	}
}

func freqLoop() {
	for {
		if settleWar() {
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(60 * time.Second)
		}
	}
}

func SettleSeasonRunning() bool {
	settleSeasonLock.RLock()
	defer settleSeasonLock.RUnlock()

	return settleSeasonRunning
}

func delDisbandGuilds() {
	shutdown.AddOneRequest("delDisbandGuilds")
	defer shutdown.DoneOneRequest("delDisbandGuilds")

	now := time.Now()
	nwarDate := now.Format("20060102")
	ywarDate := now.AddDate(0, 0, -1).Format("20060102")
	wt, gwt := common.GetGuildYWTokenByTime(now, 0, 0)
	_, lgwt := common.GetGuildYWTokenByTime(now, 0, -1)
	if gs, err := guild.GetDisbandedGuilds(); err != nil {
		return
	} else {
		for _, id := range gs {
			if g, err := GetGuild(id, true, wt, gwt, lgwt); err != nil {
				if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
					if err := guild.Del(id); err != nil {
						return
					} else {
						guild.DelDisbandedGuild(id)
					}
				} else {
					continue
				}
			} else if g.WarDate != nwarDate && g.WarDate != ywarDate && g.Zombie {
				if err := guild.Del(id); err != nil {
					return
				} else {
					guild.DelDisbandedGuild(id)
				}
			}
		}
	}
}

func settleWar() bool {
	shutdown.AddOneRequest("settleWar")
	defer shutdown.DoneOneRequest("settleWar")

	now := time.Now()
	nhour := now.Hour()
	nwarDate := now.Format("20060102")
	ywarDate := now.AddDate(0, 0, -1).Format("20060102")

	if fg, tg, wd, err := guild.GetSingleWarLog(); err == nil {
		//check if war already ends
		if !WarDateLegal(wd, nwarDate, ywarDate, nhour) {
			Info("begin settle war:", fg, tg, wd)
			if fsdus, err := guild.QuerySetDefendUsers(fg, tg, wd); err == nil {
				Info("gid:", fg, "defend size:", len(fsdus))
				if tsdus, err := guild.QuerySetDefendUsers(tg, fg, wd); err == nil {
					Info("gid:", tg, "defend size:", len(tsdus))
					if logs, err := guild.GetSingleWarLogs(fg, tg, wd); err == nil {
						Info("war:", fg, tg, wd, "log size:", len(logs))
						//deal with logs to generate result(which guild wins)
						fromWinMedalMap := make(map[string]int) //opponent username:max win medal
						toWinMedalMap := make(map[string]int)   //opponent username:max win medal
						fum := make(map[string]int)             //from attack&defend users
						tum := make(map[string]int)             //to attack&defend users

						for _, u := range fsdus {
							fum[u.Username] = 1
						}

						for _, u := range tsdus {
							tum[u.Username] = 1
						}

						var fl, tl int
						for _, l := range logs {
							if l.FromGID == fg {
								fl++
								fum[l.FromUsername] = 1
								if fromWinMedalMap[l.ToUsername] < l.WinMedal {
									fromWinMedalMap[l.ToUsername] = l.WinMedal
								}
							} else if l.FromGID == tg {
								tl++
								tum[l.FromUsername] = 1
								if toWinMedalMap[l.ToUsername] < l.WinMedal {
									toWinMedalMap[l.ToUsername] = l.WinMedal
								}
							}
						}
						Info("war:", fg, tg, wd, "from log size:", fl, "to log size:", tl)
						Info("war:", fg, tg, wd, "from member size:", len(fum), "to member size:", len(tum))

						var fw, tw int
						for u, w := range fromWinMedalMap {
							Info("fg win from:", u, "medal:", w)
							fw += w
						}

						for u, w := range toWinMedalMap {
							Info("tg win from:", u, "medal:", w)
							tw += w
						}
						Info("war:", fg, tg, wd, "from win medal:", fw, "to win medal:", tw)

						//del this war's logs
						if effect, err := guild.DelSingleWarlogs(fg, tg, wd); err == nil {
							Info("delete log effect:", fg, tg, wd, effect)

							//send reward when there is a winner
							Info("sending reward...")
							if fw > tw {
								for u, _ := range fum {
									Info("guild war reward send to:", u)
									mails.SendToOneT(L("o12"), u, L("guild15"), L("guild16"), int64(constants.GUILD_WAR_WIN_REWARD_DIAMOND), 0, "", true, 0)
								}
							} else if fw < tw {
								for u, _ := range tum {
									Info("guild war reward send to:", u)
									mails.SendToOneT(L("o12"), u, L("guild15"), L("guild16"), int64(constants.GUILD_WAR_WIN_REWARD_DIAMOND), 0, "", true, 0)
								}
							}

							//add clothes share cnt
							Info("adding clothes share cnt...")
							if tx, err := dao_guild.BeginTransaction(); err == nil {
								defer func() {
									if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
										Err(re)
									}
								}()

								//add clothes share cnt
								for _, u := range fsdus {
									Info("guild war clothes share add:", fg, u.Username)
									guild.AddClothesShareCnt(tx, fg, u.ShareClothesUsed)
								}

								for _, u := range tsdus {
									Info("guild war clothes share add:", tg, u.Username)
									guild.AddClothesShareCnt(tx, tg, u.ShareClothesUsed)
								}

								if err = tx.Commit(); err != nil {
									Err(err)
								}
							}

							return true
						}
					}
				}
			}
		}
	}
	return false
}

func settleSeason() {
	shutdown.AddOneRequest("settleSeason")
	defer shutdown.DoneOneRequest("settleSeason")

	now := time.Now()
	wt, gwt := common.GetGuildYWTokenByTime(now, 0, 0)
	lwt, lgwt := common.GetGuildYWTokenByTime(now, 0, -1)
	_, llgwt := common.GetGuildYWTokenByTime(now, 0, -2)
	Info("guild war season settle check settle gwt:", gwt, "lgwt:", lgwt)
	if ShouldSettle(now) {
		Info("guild war season settle begin settle:", gwt, lgwt)
		//check if already settled
		if !guild.GuildWarSettled(lgwt) {
			settleSeasonLock.Lock()
			settleSeasonRunning = true
			defer func() {
				settleSeasonRunning = false
				settleSeasonLock.Unlock()
			}()

			Info("guild war season settle not settle yet:", gwt, lgwt)
			shouldUpdateMedalList := make([]int64, 0)

			//reset rank data in redis from db
			Info("guild war season settle reset guild medal rank in redis with data from db..")
			start := 0
			step := 500
			var sgs []SGuild
			var err error

			normalRewardGuilds := make([]SGuild, 0)
			if sgs, err = guild.ListSGuild(start, step); err != nil {
				Err(err)
				return
			} else if len(sgs) <= 0 {
				Info("guild war season settle no guild in db")
				return
			}

			for {
				if err != nil {
					return
				} else if len(sgs) <= 0 {
					break
				} else {
					for _, g := range sgs {
						if g.Zombie {
							continue
						}

						if g.GWT != lgwt {
							Info("guild war season settle guild gwt not right!maybe a guild never match war:", g.ID, g.GWT, lgwt)
							continue
						}

						if err := RefreshGuildMedalRank(g.ID, g.MedalCnt, lgwt); err != nil {
							return
						} else {
							Info("guild war season settle gid:", g.ID, "lgwt medal:", g.MedalCnt)
							if g.MedalCnt >= 1 {
								shouldUpdateMedalList = append(shouldUpdateMedalList, g.ID)
								normalRewardGuilds = append(normalRewardGuilds, g)
							}
						}
					}
					start += step
					sgs, err = guild.ListSGuild(start, step)
				}
			}

			//settle with data in redis
			Info("guild war season settle get rank data from redis and settle..")
			res := make([]GuildWarSettleInfo, 0)
			gids, _ := guild.ListGuildMedalRank(lgwt)

			gs := make([]Guild, 0)
			for _, id := range gids {
				if id > 0 {
					if g, err := GetGuild(id, true, lwt, lgwt, llgwt); err == nil {
						if !g.Zombie && g.ID == id {
							gs = append(gs, g)
						}
					} else if err.Error() != constants.SQL_NO_DATA_ERR_MSG {
						Err("guild war season settle error when get guild from db with data from rank!")
						return
					} else {
						Info("gid in rank but not exist in db:", id)
					}
				}
			}

			if len(gs) != len(gids) {
				Err("guild war season settle gids and gs size not equal:gids:", len(gids), "gs:", len(gs))
			}

			if len(gs) <= 0 {
				Info("guild war season settle last gwt guild medal rank size 0:", lgwt)
				return
			} else {
				for i := 0; i < len(gs); i++ {
					g := gs[i]
					rank := i + 1
					rd := GetGuildWarMedalRankReward(rank)
					Info("guild war season settle gid:", g.ID, "rank:", rank, "reward diamond:", rd)
					gwsi := GuildWarSettleInfo{
						GID:   g.ID,
						Rank:  rank,
						Medal: g.MedalCnt,
					}
					rewards := make([]MemberReward, 0)

					if mbs, err := ListMembers(g.ID, "", false, wt); err != nil {
						Err("guild war season settle err when get members:", g.ID)
						continue
					} else {
						for _, m := range mbs {
							Info("gid:", g.ID, "username:", m.Username, "reward diamond:", rd)
							mr := MemberReward{
								Username:      m.Username,
								RewardDiamond: rd,
							}
							rewards = append(rewards, mr)
						}
						gwsi.Rewards = rewards
					}

					res = append(res, gwsi)
				}
			}

			//send reward/reset new medal
			//do not return when error happens
			resb, _ := json.Marshal(res)
			if err := guild.AddGuildWarSettleInfo(lgwt, string(resb)); err == nil {
				//send reward
				Info("guild war season settle send rank reward...")
				for _, gwsi := range res {
					for _, mr := range gwsi.Rewards {
						Info("season guild war rank reward send to gid:", gwsi.GID, "player:", mr.Username, "rank:", gwsi.Rank, "medal", gwsi.Medal)
						mails.SendToOneT(L("o12"), mr.Username, L("guild13"), fmt.Sprintf(L("guild14"), gwsi.Rank), int64(mr.RewardDiamond), 0, "", true, 0)
					}

					nm := gwsi.Medal / 6
					if nm > constants.GUILD_WAR_NEW_SEASON_MEDAL_MAX {
						nm = constants.GUILD_WAR_NEW_SEASON_MEDAL_MAX
					}
					//new medal rank
					RefreshGuildMedalRank(gwsi.GID, nm, gwt)
				}

				//reset medal
				Info("guild war season settle reset medal...")
				for _, id := range shouldUpdateMedalList {
					guild.UpdateNewSeasonMedal(id, gwt, lgwt)
				}

				//send normal reward
				for _, g := range normalRewardGuilds {
					var reward int
					if g.MedalCnt < constants.GUILD_MEDAL_LEVEL_2 {
						reward = constants.GUILD_MEDAL_LEVEL_1_GOLD
					} else if g.MedalCnt < constants.GUILD_MEDAL_LEVEL_3 {
						reward = constants.GUILD_MEDAL_LEVEL_2_GOLD
					} else if g.MedalCnt < constants.GUILD_MEDAL_LEVEL_4 {
						reward = constants.GUILD_MEDAL_LEVEL_3_GOLD
					} else if g.MedalCnt < constants.GUILD_MEDAL_LEVEL_5 {
						reward = constants.GUILD_MEDAL_LEVEL_4_GOLD
					} else if g.MedalCnt < constants.GUILD_MEDAL_LEVEL_6 {
						reward = constants.GUILD_MEDAL_LEVEL_5_GOLD
					} else {
						reward = constants.GUILD_MEDAL_LEVEL_6_GOLD
					}
					Info("guild war season settle normal reward gid:", g.ID, "reward gold:", reward)
					if mbs, err := ListMembers(g.ID, "", false, wt); err != nil {
						Err("guild war season settle send normal reward error when get members:", g.ID)
						continue
					} else {
						for _, m := range mbs {
							Info("season guild war medalcnt reward send to gid:", g.ID, "player:", m.Username, "medal", g.MedalCnt, "reward gold:", reward)
							mails.SendToOneT(L("o12"), m.Username, L("guild17"), L("guild18"), 0, int64(reward), "", true, 0)
						}
					}
				}
			} else {
				Err("guild war season settle try add guild war settle info err!!!!!!!!!")
			}
		}
	}
}

func refreshRank() {
	nml := make([]Guild, 0)
	nal := make([]Guild, 0)

	now := time.Now()
	wt, gwt := common.GetGuildYWTokenByTime(now, 0, 0)
	_, lgwt := common.GetGuildYWTokenByTime(now, 0, -1)

	// wt, gwt := common.GetGuildYearWeekToken()
	Info("refresh guild rank:", wt, gwt)
	if mids, err := guild.ListGuildMedalRank(gwt); err == nil {
		for _, id := range mids {
			if id > 0 {
				if g, err := GetGuild(id, true, wt, gwt, lgwt); err == nil {
					if g.Zombie {
						guild.DelGuildRank(g.ID, wt, gwt)
					} else {
						nml = append(nml, g)
					}
				} else if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
					guild.DelGuildRank(id, wt, gwt)
				}
			}
		}
	}

	if aids, err := guild.ListGuildActivityRank(wt); err == nil {
		for _, id := range aids {
			if id > 0 {
				if g, err := GetGuild(id, true, wt, gwt, lgwt); err == nil {
					if g.Zombie {
						guild.DelGuildRank(g.ID, wt, gwt)
					} else {
						nal = append(nal, g)
					}
				} else if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
					guild.DelGuildRank(id, wt, gwt)
				}
			}
		}
	}

	nmlb, _ := json.Marshal(nml)
	nalb, _ := json.Marshal(nal)
	nmls := string(nmlb)
	nals := string(nalb)

	guild.SetGuildMedalRankCache(gwt, &nmls)
	guild.SetGuildActivityRankCache(wt, &nals)
}

func loadGuildPKSubjectList() {
	guildPKSubjectList = make([]string, 0)

	file := "scripts/data/data_pk_team.xlsx"
	if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
		file = "scripts/data/data_pk_team_kor.xlsx"
	}

	xlFile, err := xlsx.OpenFile(file)
	if err != nil {
		panic(err)
		return
	}

	for _, sheet := range xlFile.Sheets {
		if sheet.Name == "Sheet1" {
			for i, row := range sheet.Rows {
				if i > 0 {
					var id string
					for j, cell := range row.Cells {
						if j == 0 {
							id, _ = cell.String()
							break
						}
					}
					guildPKSubjectList = append(guildPKSubjectList, id)
				}
			}
		}
	}
	Info("guild pk subject size:", len(guildPKSubjectList), "sample:", guildPKSubjectList[0])
}

func getGuildWeekPKSubjects() (res []string) {
	res = make([]string, 0)
	resm := make(map[string]int)
	for i := 0; i < constants.GUILD_WEEK_PK_SUBJECT_SIZE; i++ {
		c := guildPKSubjectList[rand.Intn(len(guildPKSubjectList))]
		for resm[c] != 0 {
			c = guildPKSubjectList[rand.Intn(len(guildPKSubjectList))]
		}

		res = append(res, c)
		resm[c] = 1
	}
	return
}

func setWeekRewardClothes(canPanic bool) {
	weekLock.Lock()
	defer weekLock.Unlock()

	weekToken, _ = common.GetGuildYearWeekToken()
	weekRewadClothes = make([]string, 0)

	var err error
	for len(weekRewadClothes) <= 0 {
		if weekRewadClothes, err = guild.GetGuildWeekRewardClothes(weekToken); err != nil {
			rc := clothes.GetGuildWeekRewardClothes()
			if err = guild.SetGuildWeekRewardClothes(weekToken, rc); err != nil {
				if canPanic {
					panic(err)
				}
			} else {
				weekRewadClothes = rc
			}
		} else if len(weekRewadClothes) <= 0 {
			if canPanic {
				panic("guild week reward clothes size wrong")
			}
		}
	}
	Info("guild week reward clohtes:", weekToken, weekRewadClothes)
}

func setWeekPKSubjects(canPanic bool) {
	weekLock.Lock()
	defer weekLock.Unlock()

	weekToken, _ = common.GetGuildYearWeekToken()
	weekPKSubjects = make([]string, 0)

	var err error
	for len(weekPKSubjects) <= 0 {
		if weekPKSubjects, err = guild.GetGuildWeekPKSubjects(weekToken); err != nil {
			rc := getGuildWeekPKSubjects()
			if err = guild.SetGuildWeekPKSubjects(weekToken, rc); err != nil {
				if canPanic {
					panic(err)
				}
			} else {
				weekPKSubjects = rc
			}
		} else if len(weekPKSubjects) <= 0 {
			if canPanic {
				panic("guild week reward clothes size wrong")
			}
		}
	}
	Info("guild week pk subjects:", weekToken, weekPKSubjects)
}
