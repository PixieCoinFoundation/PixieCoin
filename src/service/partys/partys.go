package partys

import (
	"appcfg"
	cp "common_db/party"
	"constants"
	"dao_party"
	"db/event"
	"db/global"
	"db/mail"
	"db/party"
	"encoding/json"
	"fmt"
	. "language"
	. "logger"
	"rao"
	"service/files"
	"service/mails"
	"shutdown"
	"time"
	. "types"
	. "zk_manager"

	"github.com/garyburd/redigo/redis"
)

const (
	TICKET_PARTION = 0.15
)

func init() {
	if appcfg.GetServerType() == "" && appcfg.GetBool("main_server", false) {
		go partyLoop()
	}
}

func TryPartyQinmi(u1, u2, ds string) bool {
	return global.TryInviteQinmiFriend(u1, u2, ds)
}

func QueryPartyQinmiUsed(u1, u2, ds string) (bool, error) {
	return global.GetQinmiPartyDone(u1, u2, ds)
}

func GetPrizePartyPlayerRankIds() (res []int64) {
	conn := rao.GetConn()
	defer conn.Close()

	if s, err := redis.String(conn.Do("GET", constants.PARTY_PRIZE_PLAYER_RANK_CACHE_KEY)); err != nil {
		Err(err)
		return
	} else {
		s = s[1 : len(s)-1]

		if err := json.Unmarshal([]byte(s), &res); err != nil {
			Err(err)
			return
		}
	}

	return
}

func GetPrizePartyPoolRankIds() (res []int64) {
	conn := rao.GetConn()
	defer conn.Close()

	if s, err := redis.String(conn.Do("GET", constants.PARTY_PRIZE_POOL_RANK_CACHE_KEY)); err != nil {
		Err(err)
		return
	} else {
		s = s[1 : len(s)-1]

		if err := json.Unmarshal([]byte(s), &res); err != nil {
			Err(err)
			return
		}
	}

	return
}

func partyLoop() {
	i := 0
	for {
		setPartyPrizeRankCache1()

		if i%10 == 0 {
			processParty1()
			i = 0
		}

		i++
		time.Sleep(60 * time.Second)
	}
}

func setPartyPrizeRankCache1() {
	conn := rao.GetConn()
	defer conn.Close()

	//player rank
	if pr, err := party.GetPrizePartyPlayerRank(); err == nil {
		//set cache
		if prb, err := json.Marshal(pr); err != nil {
			Err(err)
		} else {
			if _, err := conn.Do("SET", constants.PARTY_PRIZE_PLAYER_RANK_CACHE_KEY, "'"+string(prb)+"'"); err != nil {
				Err(err)
			}
		}
	}

	//pool rank
	if por, err := party.GetPrizePartyPoolRank(); err == nil {
		//set cache
		if pob, err := json.Marshal(por); err != nil {
			Err(err)
		} else {
			if _, err := conn.Do("SET", constants.PARTY_PRIZE_POOL_RANK_CACHE_KEY, "'"+string(pob)+"'"); err != nil {
				Err(err)
			}
		}
	}
}

// func processParty() {
// 	for {
// 		processParty1()
// 		time.Sleep(600 * time.Second)
// 	}

// }
//processParty1 Party结束后的相关操作
func processParty1() {
	shutdown.AddOneRequest("processParty1")
	defer shutdown.DoneOneRequest("processParty1")

	now := time.Now()
	met := now.Unix() + constants.DEFAULT_MAIL_EXPIRE_TIME

	if ps, err := party.GetPartyByTimeFromDB(0, now.Unix()-60); err != nil {
		return
	} else {
		for _, p := range ps {
			Info("party id:", p.ID, "subject:", p.Subject, "closed.send reward..")
			pc, _ := party.GetPartyItemSizeFromDB(p.ID)
			Info("party:", p.ID, p.Subject, "closed.player size:", pc)
			deleteKeys := make(map[string]int)

			deleteKeys[rao.GetPartySelfKey(p.ID)] = 1
			deleteKeys[rao.GetPartyItemListKey(p.ID)] = 1
			deleteKeys[rao.GetPartyItemRankListKey(p.ID)] = 1
			deleteKeys[rao.GetPartyInviteFlowerKey(p.ID)] = 1
			deleteKeys[rao.GetPartyInviteAttendKey(p.ID)] = 1
			deleteKeys[rao.GetPartyInviteListKey(p.ID)] = 1
			//players 邀请的玩家目录日志
			players, imgs := party.GetPartyPlayers(p.ID)
			for _, u := range players {
				deleteKeys[rao.GetPartyItemSelfKey(p.ID, u)] = 1
				deleteKeys[rao.GetPartyItemFlowerSetKey(p.ID, u)] = 1
				deleteKeys[rao.GetPartyItemCmtListKey(p.ID, u)] = 1
				DelPartyItemLock(p.ID, u)

			}
			//Party相关数据
			if p.Type == constants.PARTY_CASUAL_TYPE {
				//host party prize
				mails.SendToOne("", p.Username, fmt.Sprintf(L("party1"), p.Subject), fmt.Sprintf(L("party2"), p.Subject), 0, 8000, "", true, met)

				//player cnt prize
				if pc > 100 {
					mails.SendToOne("", p.Username, fmt.Sprintf(L("party1"), p.Subject), fmt.Sprintf(L("party3"), p.Subject), 80, 0, "", true, met)
					event.AddEventPartyLog(fmt.Sprintf("%s:%s:%d:%d:%d", "J", "X", 0, pc, 0), p.Username)
				}
				//top prize
				if tops, err := party.GetPartyTopItems(p.ID, 20, p.Type); err != nil {
					return
				} else {
					title := fmt.Sprintf(L("party4"), p.Subject)
					for i, pi := range tops {
						Info("player:", pi.Username, "subject:", p.Subject, "rank:", i+1, "popularity:", pi.Popularity, "partner:", pi.Partner)

						// var wb, pwb string
						// if pi.Partner != "" {
						// 	wb = L("party6") + pi.PartnerNickname
						// 	pwb = L("party6") + pi.Nickname
						// }
						jc := 1.0
						if pi.Partner != "" && pi.PartnerQinmi >= constants.QINMI_FRIEND_MIN {
							jc = 1.5
						}

						var content, pcontent string
						if pi.Partner != "" {
							content = fmt.Sprintf(L("party18"), p.Subject, pi.PartnerNickname, i+1)
							pcontent = fmt.Sprintf(L("party18"), p.Subject, pi.Nickname, i+1)
						} else {
							content = fmt.Sprintf(L("party17"), p.Subject, i+1)
							pcontent = fmt.Sprintf(L("party17"), p.Subject, i+1)
						}
						// content := fmt.Sprintf(L("party12"), p.Subject, wb, i+1)
						// pcontent := fmt.Sprintf(L("party12"), p.Subject, pwb, i+1)

						if i == 0 {
							mails.SendToOneT("", pi.Username, title, content, 0, int64(float64(constants.CASUAL_PARTY_1_GOLD)*jc), "", true, constants.MAIL_TYPE_XX_PARTY_RANK1)
							plog(pi.Username, pi.Nickname, 1, 0, int64(float64(constants.CASUAL_PARTY_1_GOLD)*jc), &p)
							if pi.Partner != "" {
								mails.SendToOneT("", pi.Partner, title, pcontent, 0, int64(float64(constants.CASUAL_PARTY_1_GOLD)*jc), "", true, constants.MAIL_TYPE_XX_PARTY_RANK1)
								plog(pi.Partner, pi.PartnerNickname, 1, 0, int64(float64(constants.CASUAL_PARTY_1_GOLD)*jc), &p)
							}
						} else if i == 1 {
							mails.SendToOneT("", pi.Username, title, content, 0, int64(float64(constants.CASUAL_PARTY_2_GOLD)*jc), "", true, constants.MAIL_TYPE_XX_PARTY_RANK3)
							plog(pi.Username, pi.Nickname, 2, 0, int64(float64(constants.CASUAL_PARTY_2_GOLD)*jc), &p)
							if pi.Partner != "" {
								mails.SendToOneT("", pi.Partner, title, pcontent, 0, int64(float64(constants.CASUAL_PARTY_2_GOLD)*jc), "", true, constants.MAIL_TYPE_XX_PARTY_RANK3)
								plog(pi.Partner, pi.PartnerNickname, 2, 0, int64(float64(constants.CASUAL_PARTY_2_GOLD)*jc), &p)
							}
						} else if i == 2 {
							mails.SendToOneT("", pi.Username, title, content, 0, int64(float64(constants.CASUAL_PARTY_3_GOLD)*jc), "", true, constants.MAIL_TYPE_XX_PARTY_RANK3)
							plog(pi.Username, pi.Nickname, 3, 0, int64(float64(constants.CASUAL_PARTY_3_GOLD)*jc), &p)
							if pi.Partner != "" {
								mails.SendToOneT("", pi.Partner, title, pcontent, 0, int64(float64(constants.CASUAL_PARTY_3_GOLD)*jc), "", true, constants.MAIL_TYPE_XX_PARTY_RANK3)
								plog(pi.Partner, pi.PartnerNickname, 3, 0, int64(float64(constants.CASUAL_PARTY_3_GOLD)*jc), &p)
							}
						} else {
							mails.SendToOne("", pi.Username, title, content, 0, int64(float64(constants.CASUAL_PARTY_20_GOLD)*jc), "", true, met)
							plog(pi.Username, pi.Nickname, i+1, 0, int64(float64(constants.CASUAL_PARTY_20_GOLD)*jc), &p)
							if pi.Partner != "" {
								mails.SendToOne("", pi.Partner, title, pcontent, 0, int64(float64(constants.CASUAL_PARTY_20_GOLD)*jc), "", true, met)
								plog(pi.Partner, pi.PartnerNickname, i+1, 0, int64(float64(constants.CASUAL_PARTY_20_GOLD)*jc), &p)
							}
						}
						event.AddEventPartyLog(fmt.Sprintf("%s:%s:%d:%d:%d", "C", "X", 0, 0, i+1), pi.Username)
					}
				}
			} else if p.Type == constants.PARTY_PRIZE_TYPE {
				var ticketGoldPrize int64
				var ticketDiamondPrize int64
				var goldPool int64
				var diamondPool int64
				dt := "GJ"
				if p.MoneyType == 1 {
					dt = "DJ"
				}

				part := 0.8
				// if p.SingleType == 1 {
				// 	part = 0.85
				// }
				if p.MoneyType == 0 {
					//gold
					ticketGoldPrize = int64(float64(pc) * float64(p.Ticket) * TICKET_PARTION)
					if ticketGoldPrize == 0 {
						ticketGoldPrize = 1
					}
					goldPool = int64(p.StartPool) + int64(float64(pc)*float64(p.Ticket)*part)
				} else if p.MoneyType == 1 {
					ticketDiamondPrize = int64(float64(pc) * float64(p.Ticket) * TICKET_PARTION)
					if ticketDiamondPrize == 0 {
						ticketDiamondPrize = 1
					}
					diamondPool = int64(p.StartPool) + int64(float64(pc)*float64(p.Ticket)*part)
				}

				if pc > 100 {
					event.AddEventPartyLog(fmt.Sprintf("%s:%s:%d:%d:%d", "J", dt, 0, pc, 0), p.Username)
				}

				//ticket prize
				mails.SendToOne("", p.Username, fmt.Sprintf(L("party7"), p.Subject), fmt.Sprintf(L("party8"), p.Subject), ticketDiamondPrize, ticketGoldPrize, "", true, met)

				//top prize
				if tops, err := party.GetPartyTopItems(p.ID, -1, p.Type); err != nil {
					return
				} else {
					title := fmt.Sprintf(L("party4"), p.Subject)

					rank1d := diamondPool / 2
					rank1g := goldPool / 2
					rank2d := diamondPool / 5
					rank2g := goldPool / 5
					rank3d := diamondPool / 10
					rank3g := goldPool / 10
					rank420d := diamondPool / 5
					rank420g := goldPool / 5

					Info("prize diamond pool:", diamondPool, rank1d, rank2d, rank3d, rank420d)
					Info("prize gold pool:", goldPool, rank1g, rank2g, rank3g, rank420g)

					rank1s := make([]*PartyItem, 0)
					rank2s := make([]*PartyItem, 0)
					rank3s := make([]*PartyItem, 0)
					rank420s := make([]*PartyItem, 0)
					var rankNow, popuNow int

					topLen := len(tops)
					for i := 0; i < topLen; i++ {
						if i == 0 {
							rankNow = 1
							popuNow = tops[i].Popularity
							rank1s = append(rank1s, tops[i])
						} else {
							if tops[i].Popularity < popuNow {
								rankNow++
								popuNow = tops[i].Popularity
							}

							if rankNow > 20 {
								break
							}

							switch rankNow {
							case 1:
								rank1s = append(rank1s, tops[i])
							case 2:
								rank2s = append(rank2s, tops[i])
							case 3:
								rank3s = append(rank3s, tops[i])
							default:
								tops[i].Rank = rankNow
								rank420s = append(rank420s, tops[i])
							}
						}
						Info("index:", i, "rank:", rankNow, "player:", tops[i].Username, "partner:", tops[i].Partner, "popu:", tops[i].Popularity)

						if i >= 19 && i < topLen-1 && tops[i+1].Popularity < popuNow {
							break
						}
					}

					l1s := int64(len(rank1s))
					if l1s > 0 {
						rank1d = rank1d / l1s
						rank1g = rank1g / l1s
					}

					l2s := int64(len(rank2s))
					if l2s > 0 {
						rank2d = rank2d / l2s
						rank2g = rank2g / l2s
					}

					l3s := int64(len(rank3s))
					if l3s > 0 {
						rank3d = rank3d / l3s
						rank3g = rank3g / l3s
					}

					l420s := int64(len(rank420s))
					if l420s > 0 {
						rank420d = rank420d / l420s
						rank420g = rank420g / l420s
					}

					for _, pi := range rank1s {
						Info("player:", pi.Username, "subject:", p.Subject, "rank:", 1, "popularity:", pi.Popularity, "partner:", pi.Partner)
						content := fmt.Sprintf(L("party5"), p.Subject, 1, l1s)
						ycontent := fmt.Sprintf(L("party9"), p.Subject, 1, l1s)
						if pi.Partner == "" {
							if diamondPool > 0 && rank1d <= 0 {
								rank1d = 1
							}
							if goldPool > 0 && rank1g <= 0 {
								rank1g = 1
							}
							mails.SendToOne("", pi.Username, title, content, rank1d, rank1g, "", true, met)
							plog(pi.Username, pi.Nickname, 1, rank1d, rank1g, &p)
						} else {
							rank1dj := rank1d * 7 / 10
							rank1gj := rank1g * 7 / 10

							if diamondPool > 0 && rank1dj <= 0 {
								rank1dj = 1
							}

							if goldPool > 0 && rank1gj <= 0 {
								rank1gj = 1
							}

							mails.SendToOne("", pi.Username, title, content, rank1dj, rank1gj, "", true, met)
							mails.SendToOne("", pi.Partner, title, ycontent, rank1d*3/10, rank1g*3/10, "", true, met)

							plog(pi.Username, pi.Nickname, 1, rank1dj, rank1gj, &p)
							plog(pi.Partner, pi.PartnerNickname, 1, rank1d*3/10, rank1g*3/10, &p)
						}
						event.AddEventPartyLog(fmt.Sprintf("%s:%s:%d:%d:%d", "C", dt, 0, 0, 1), pi.Username)
					}

					for _, pi := range rank2s {
						Info("player:", pi.Username, "subject:", p.Subject, "rank:", 2, "popularity:", pi.Popularity, "partner:", pi.Partner)
						content := fmt.Sprintf(L("party5"), p.Subject, 2, l2s)
						ycontent := fmt.Sprintf(L("party9"), p.Subject, 2, l2s)
						if pi.Partner == "" {
							if diamondPool > 0 && rank2d <= 0 {
								rank2d = 1
							}
							if goldPool > 0 && rank2g <= 0 {
								rank2g = 1
							}
							mails.SendToOne("", pi.Username, title, content, rank2d, rank2g, "", true, met)
							plog(pi.Username, pi.Nickname, 1, rank2d, rank2g, &p)
						} else {
							rank2dj := rank2d * 7 / 10
							rank2gj := rank2g * 7 / 10

							if diamondPool > 0 && rank2dj <= 0 {
								rank2dj = 1
							}

							if goldPool > 0 && rank2gj <= 0 {
								rank2gj = 1
							}

							mails.SendToOne("", pi.Username, title, content, rank2dj, rank2gj, "", true, met)
							mails.SendToOne("", pi.Partner, title, ycontent, rank2d*3/10, rank2g*3/10, "", true, met)

							plog(pi.Username, pi.Nickname, 2, rank2dj, rank2gj, &p)
							plog(pi.Partner, pi.PartnerNickname, 2, rank2d*3/10, rank2g*3/10, &p)
						}

						event.AddEventPartyLog(fmt.Sprintf("%s:%s:%d:%d:%d", "C", dt, 0, 0, 2), pi.Username)
					}

					for _, pi := range rank3s {
						Info("player:", pi.Username, "subject:", p.Subject, "rank:", 3, "popularity:", pi.Popularity, "partner:", pi.Partner)
						content := fmt.Sprintf(L("party5"), p.Subject, 3, l3s)
						ycontent := fmt.Sprintf(L("party9"), p.Subject, 3, l3s)

						if pi.Partner == "" {
							if diamondPool > 0 && rank3d <= 0 {
								rank3d = 1
							}
							if goldPool > 0 && rank3g <= 0 {
								rank3g = 1
							}
							mails.SendToOne("", pi.Username, title, content, rank3d, rank3g, "", true, met)
							plog(pi.Username, pi.Nickname, 3, rank3d, rank3g, &p)
						} else {
							rank3dj := rank3d * 7 / 10
							rank3gj := rank3g * 7 / 10

							if diamondPool > 0 && rank3dj <= 0 {
								rank3dj = 1
							}

							if goldPool > 0 && rank3gj <= 0 {
								rank3gj = 1
							}

							mails.SendToOne("", pi.Username, title, content, rank3dj, rank3gj, "", true, met)
							mails.SendToOne("", pi.Partner, title, ycontent, rank3d*3/10, rank3g*3/10, "", true, met)

							plog(pi.Username, pi.Nickname, 3, rank3dj, rank3gj, &p)
							plog(pi.Partner, pi.PartnerNickname, 3, rank3d*3/10, rank3g*3/10, &p)
						}
						event.AddEventPartyLog(fmt.Sprintf("%s:%s:%d:%d:%d", "C", dt, 0, 0, 3), pi.Username)
					}

					for _, pi := range rank420s {
						Info("player:", pi.Username, "subject:", p.Subject, "rank:", pi.Rank, "popularity:", pi.Popularity, "partner:", pi.Partner)
						content := fmt.Sprintf(L("party5"), p.Subject, pi.Rank, l420s)
						ycontent := fmt.Sprintf(L("party9"), p.Subject, pi.Rank, l420s)

						if pi.Partner == "" {
							if diamondPool > 0 && rank420d <= 0 {
								rank420d = 1
							}
							if goldPool > 0 && rank420g <= 0 {
								rank420g = 1
							}
							mails.SendToOne("", pi.Username, title, content, rank420d, rank420g, "", true, met)
							plog(pi.Username, pi.Nickname, pi.Rank, rank420d, rank420g, &p)
						} else {
							rank420dj := rank420d * 7 / 10
							rank420gj := rank420g * 7 / 10

							if diamondPool > 0 && rank420dj <= 0 {
								rank420dj = 1
							}

							if goldPool > 0 && rank420gj <= 0 {
								rank420gj = 1
							}

							mails.SendToOne("", pi.Username, title, content, rank420dj, rank420gj, "", true, met)
							mails.SendToOne("", pi.Partner, title, ycontent, rank420d*3/10, rank420g*3/10, "", true, met)

							plog(pi.Username, pi.Nickname, pi.Rank, rank420dj, rank420gj, &p)
							plog(pi.Partner, pi.PartnerNickname, pi.Rank, rank420d*3/10, rank420g*3/10, &p)
						}
						event.AddEventPartyLog(fmt.Sprintf("%s:%s:%d:%d:%d", "C", dt, 0, 0, pi.Rank), pi.Username)
					}
				}
			} else if p.Type == constants.PARTY_FORBID_TYPE {
				Info("party id:", p.ID, "subject:", p.Subject, "FORBID!.just close!")
			} else if p.Type == constants.PARTY_COSPLAY_TYPE {
				token := fmt.Sprintf("settle_cosplay_%d", p.ID)
				if global.DoOnceJob(token) {
					Info("party is cosplay:", p.ID, "closed.begin process..")
					defer Info("party cosplay process end:", p.ID)

					if partyItems, err := party.GetPartyRankItems(p.ID, 0, ""); err == nil {
						length := len(partyItems)
						Info("party cosplay:", p.ID, p.Subject, "closed.player size db:", pc, "size redis:", length)

						pos1 := int(length / 10) //upper 10%
						if pos1 <= 0 {
							pos1 = 1
						}
						pos3 := int(length / 2) //upper 50%
						if pos3 <= 0 {
							pos3 = 1
						}

						//iterate all items
						for j := 0; j < length; j++ {
							pi := partyItems[j]
							rank := j + 1

							var diamond int64
							var gold int64
							var content string

							if rank <= 10 {
								r := "前十名"
								if rank == 1 {
									r = "第一名"
									diamond = 500
								} else if rank == 2 {
									r = "第二名"
									diamond = 300
								} else if rank == 3 {
									r = "第三名"
									diamond = 200
								} else {
									diamond = 100
								}

								content = fmt.Sprintf("恭喜您获得官方舞会【%s】%s", p.Subject, r)
							} else {
								c1 := "进入"
								c2 := ""
								if rank <= pos1 {
									c2 = "排名前10%"
									gold = 30000
								} else if rank <= pos3 {
									c2 = "排名前50%"
									gold = 15000
								} else {
									c1 = "参与"
									gold = 8000
								}

								content = fmt.Sprintf("恭喜您%s官方舞会活动【%s】%s", c1, p.Subject, c2)
							}

							Info("cosplay:", p.ID, "名次:", rank, "玩家:", pi.Username, "diamond", diamond, "gold", gold)

							m := Mail{
								To:         pi.Username,
								Title:      "官方舞会结算奖励",
								Content:    content,
								Diamond:    diamond,
								Gold:       gold,
								Delete:     true,
								Time:       now.Unix(),
								ExpireTime: now.Unix() + constants.DEFAULT_MAIL_EXPIRE_TIME,
							}
							mail.SendToOne(m, nil)
						}
					}
				}
			} else {
				Info("party id:", p.ID, "subject:", p.Subject, "unknown type!.just close!type:", p.Type)
			}

			if p.Type != constants.PARTY_COSPLAY_TYPE || p.CloseTime <= now.Unix()-32*24*3600 {
				Info("delete party's everything in db", p.ID)
				//delete in db
				cp.DeletePartyInDB(p.ID, false, dao_party.DeletePartyStmt, dao_party.DeletePartyItemStmt, dao_party.DeletePartyFlowerStmt, dao_party.DeletePartyCmtStmt)

				if p.Type == constants.PARTY_COSPLAY_TYPE {
					token := fmt.Sprintf("%d_%s", p.ID, p.Username)
					//delete in redis list
					conn := rao.GetConn()
					defer conn.Close()

					conn.Do("ZREM", rao.GetGlobalCosplayPartyListKey(), token)
				}

				//delete files
				if len(imgs) > 0 {
					Info("delete files:", imgs[0], "...")
				}
				files.DeleteFiles(imgs)

				deleteRedisKeys(deleteKeys)
			}
		}

	}

}

func plog(username, nickname string, rank int, rd int64, rg int64, p *Party) {
	i := AddPartyItemExtra{
		Username:       username,
		Nickname:       nickname,
		PartyID:        p.ID,
		PartyName:      p.Subject,
		PartyHost:      p.Username,
		PartyStartTime: time.Unix(p.StartTime, 0).Format("2006-01-02 15:04:05"),
		PartyCloseTime: time.Unix(p.CloseTime, 0).Format("2006-01-02 15:04:05"),
		ServerAddress:  appcfg.GetAddress(),
		Rank:           rank,
		RewardGold:     rg,
		RewardDiamond:  rd,
	}
	ib, _ := json.Marshal(i)
	//gm log
	GMLog(constants.C1_SYSTEM, constants.C2_PARTY, constants.C3_JOIN_PARTY, username, string(ib))
	GMLog(constants.C1_PLAYER, constants.C2_PARTY, constants.C3_HOST_PARTY_JOIN, p.Username, string(ib))
}

func deleteRedisKeys(keys map[string]int) {
	conn := rao.GetConn()
	defer conn.Close()

	for k, _ := range keys {
		conn.Send("DEL", k)
	}

	if err := conn.Flush(); err != nil {
		Err(err)
	}
}
