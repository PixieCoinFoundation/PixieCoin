package party

import (
	"constants"
	"dao_party"
	"database/sql"
	. "logger"
	"sort"
	"strings"
	"time"
	. "types"
)

func joinedParty(partyID int64, username string) (string, error) {
	var joiner string
	if err := dao_party.GetPartyItemCntStmt.QueryRow(partyID, username, partyID, username).Scan(&joiner); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			return "", nil
		} else {
			Err(err)
			return "", err
		}

	}
	return joiner, nil
}

func refreshPartyItem(partyID int64, username string, partner, partnerNickname string, qm int) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_party.RefreshPartyItemStmt.Exec(partner, partnerNickname, qm, partyID, username); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		success = true
		addPartyPlayerCnt(partyID, 1)
	}
	return
}

func addPartyPlayerCnt(partyID int64, cnt int) (err error) {
	if _, err = dao_party.AddPartyPlayerCntStmt.Exec(cnt, partyID); err != nil {
		Err(err)
	}
	return
}

func getPrizePartyRank(rows *sql.Rows) (res []int64, err error) {
	res = make([]int64, 0)
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			Err(err)
			return
		}

		res = append(res, id)
	}

	return
}

func GetPartyPlayers(partyID int64) (res []string, files []string) {
	res = make([]string, 0)
	files = make([]string, 0)
	if rows, err := dao_party.GetPartyPlayerStmt.Query(partyID); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var u string
			var img string
			if err = rows.Scan(&u, &img); err != nil {
				Err(err)
				return
			}

			if u != "" {
				res = append(res, u)
			}
			if img != "" {
				files = append(files, img)
			}
		}
	}

	return
}

func GetPrizePartyPlayerRank() (res []int64, err error) {
	now := time.Now().Unix()
	var rows *sql.Rows
	if rows, err = dao_party.GetPrizePartyPlayerRankStmt.Query(now, now*2); err != nil {
		Err(err)
		return
	}
	defer rows.Close()
	return getPrizePartyRank(rows)
}

func GetPrizePartyPoolRank() (res []int64, err error) {
	now := time.Now().Unix()
	var rows *sql.Rows
	if rows, err = dao_party.GetPrizePartyPoolRankStmt.Query(now, now*2); err != nil {
		Err(err)
		return
	}
	defer rows.Close()
	return getPrizePartyRank(rows)
}

// func addPartyToDB(typee int, singleType int, startTime int64, closeTime int64, username string, nickname string, subject string, desc string, moneyType int, startPool int, ticket int, bgBannerType int) (id int64, err error) {
// 	var result sql.Result
// 	if result, err = dao_party.AddPartyStmt.Exec(typee, singleType, startTime, closeTime, username, nickname, subject, desc, 0, moneyType, startPool, ticket, bgBannerType); err != nil {
// 		Err(err)
// 		return
// 	} else {
// 		if id, err = result.LastInsertId(); err != nil {
// 			Err(err)
// 			return
// 		}
// 	}
// 	return
// }

func addPartyItemFlowerToDB(partyID int64, username string, sendUsername string) (err error) {
	if _, err = dao_party.AddPartyItemFlowerStmt.Exec(partyID, username, sendUsername, time.Now().Unix()); err != nil {
		if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
			Info("player already send flower:", partyID, username, sendUsername)
		} else {
			Err(err)
		}
	}
	return
}

// func DeletePartyInDB(partyID int64, partyOnly bool) {
// 	if _, err := dao_party.DeletePartyStmt.Exec(partyID); err != nil {
// 		Err(err)
// 	}

// 	if !partyOnly {
// 		if _, err := dao_party.DeletePartyItemStmt.Exec(partyID); err != nil {
// 			Err(err)
// 		}

// 		if _, err := dao_party.DeletePartyFlowerStmt.Exec(partyID); err != nil {
// 			Err(err)
// 		}

// 		if _, err := dao_party.DeletePartyCmtStmt.Exec(partyID); err != nil {
// 			Err(err)
// 		}
// 	}
// }

func addPartyItemToDB(partyID int64, username string, nickname string, head string, img string, uploadTime int64) (id int64, err error) {
	var result sql.Result
	if result, err = dao_party.AddPartyItemStmt.Exec(partyID, username, nickname, head, img, 0, 0, uploadTime); err != nil {
		Err(err)
		return
	} else {
		if id, err = result.LastInsertId(); err != nil {
			Err(err)
			return
		}

		addPartyPlayerCnt(partyID, 1)
	}
	return
}

func addPartyItemLikeCnt(partyID int64, username string, likeCnt int, unlikeCnt int) (err error) {
	if _, err = dao_party.AddPartyItemLikeCntStmt.Exec(likeCnt, unlikeCnt, partyID, username); err != nil {
		Err(err)
	}
	return
}

func addPartyCmtToDB(partyID int64, username string, sendUsername string, sendNickname string, content string, t int64) (id int64, err error) {
	var result sql.Result
	if result, err = dao_party.AddPartyItemCmtStmt.Exec(partyID, username, sendUsername, sendNickname, content, t); err != nil {
		Err(err)
		return
	} else {
		if id, err = result.LastInsertId(); err != nil {
			Err(err)
			return
		}
	}
	return
}

func getPartyFromDB(partyID int64) (res Party, err error) {
	if err = dao_party.GetPartyStmt.QueryRow(partyID).Scan(&res.ID, &res.Type, &res.StartTime, &res.CloseTime, &res.Username, &res.Nickname, &res.Subject, &res.Desc, &res.PlayerCnt, &res.MoneyType, &res.StartPool, &res.Ticket, &res.BgBannerType); err != nil {
		Err(err)
		return
	}
	return
}

func getPartyItemFromDB(partyItemID int64) (pi PartyItem, err error) {
	if err = dao_party.GetPartyItemStmt.QueryRow(partyItemID).Scan(&pi.ID, &pi.PartyID, &pi.Username, &pi.Nickname, &pi.Img, &pi.LikeCnt, &pi.UnlikeCnt, &pi.UploadTime, &pi.Partner); err != nil {
		Err(err)
		return
	}
	return
}

func GetPartyItemListFromDB(partyID int64) (res []PartyItem, err error) {
	res = make([]PartyItem, 0)

	var rows *sql.Rows
	start := 0
	size := 10000

	for {
		if rows, err = dao_party.GetPartyItemListStmt.Query(partyID, start, size); err != nil {
			Err(err)
			return
		} else {
			defer rows.Close()
			cnt := 0
			for rows.Next() {
				cnt++
				pi := PartyItem{}
				if err = rows.Scan(&pi.ID, &pi.PartyID, &pi.Username, &pi.Nickname, &pi.Head, &pi.Img, &pi.LikeCnt, &pi.UnlikeCnt, &pi.UploadTime, &pi.Partner, &pi.PartnerNickname, &pi.PartnerQinmi); err != nil {
					Err(err)
					rows.Close()
					return
				}

				res = append(res, pi)
			}
			rows.Close()
			if cnt <= 0 {
				break
			} else {
				start += size
			}
		}
	}
	return
}

func getPartyCmtFromDB(partyID int64, username string, start int, size int, asc bool) (r []PartyComment, err error) {
	r = make([]PartyComment, 0)

	var rows *sql.Rows
	var stmt *sql.Stmt

	if !asc {
		stmt = dao_party.GetPartyItemCmtStmt
	} else {
		stmt = dao_party.GetPartyItemCmtAscStmt
	}

	if rows, err = stmt.Query(partyID, username, start, size); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			pc := PartyComment{}
			if err = rows.Scan(&pc.ID, &pc.SendUsername, &pc.SendNickname, &pc.Content, &pc.Time); err != nil {
				Err(err)
				return
			}

			r = append(r, pc)
		}
	}
	return
}

func GetPartyByTimeFromDB(st int64, et int64) (r []Party, err error) {
	r = make([]Party, 0)

	var rows *sql.Rows
	if rows, err = dao_party.GetPartyByTimeStmt.Query(st, et); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			p := Party{}
			if err = rows.Scan(&p.ID, &p.Type, &p.SingleType, &p.StartTime, &p.CloseTime, &p.Username, &p.Nickname, &p.Subject, &p.Desc, &p.PlayerCnt, &p.MoneyType, &p.StartPool, &p.Ticket, &p.BgBannerType); err != nil {
				Err(err)
				return
			}

			r = append(r, p)
		}
	}
	return
}

func GetPartyItemSizeFromDB(partyID int64) (size int, err error) {
	if err = dao_party.GetPartyItemSizeStmt.QueryRow(partyID).Scan(&size); err != nil {
		Err(err)
	}
	return
}

func getPartyItemFlowerMap(partyID int64) (r map[string]int, err error) {
	r = make(map[string]int)

	var rows *sql.Rows
	if rows, err = dao_party.GetPartyItemFlowerMapStmt.Query(partyID); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var username string
			var fc int
			if err = rows.Scan(&username, &fc); err != nil {
				Err(err)
				return
			}

			// if fc > constants.PARTY_FLOWER_MAX_CNT {
			// 	fc = constants.PARTY_FLOWER_MAX_CNT
			// }

			r[username] = fc
		}
	}

	return
}

func GetPartyTopItems(partyID int64, size int, partyType int) (r PartyItemList, err error) {
	// return getPartyItemRankListFromRedis(partyID, size, "")
	var pis []PartyItem
	if pis, err = GetPartyItemListFromDB(partyID); err != nil {
		return
	}

	var fm map[string]int
	if fm, err = getPartyItemFlowerMap(partyID); err != nil {
		return
	}

	for i, pi := range pis {
		fc := fm[pi.Username]

		// if partyType == constants.PARTY_CASUAL_TYPE {
		// 	if fc > constants.PARTY_FLOWER_MAX_CNT {
		// 		fc = constants.PARTY_FLOWER_MAX_CNT
		// 	}
		// } else {
		// 	if fc > constants.PRIZE_PARTY_FLOWER_MAX_CNT {
		// 		fc = constants.PRIZE_PARTY_FLOWER_MAX_CNT
		// 	}
		// }

		if fc > constants.PARTY_FLOWER_MAX {
			fc = constants.PARTY_FLOWER_MAX
		}

		pis[i].Popularity = 100 + pi.LikeCnt*5 + fc*5 - pi.UnlikeCnt*5
		if pi.Partner != "" {
			if pi.PartnerQinmi >= 60 {
				pis[i].Popularity = pis[i].Popularity * 7 / 5
			} else {
				pis[i].Popularity = pis[i].Popularity * 6 / 5
			}
		}

		r = append(r, &pis[i])
	}

	sort.Sort(r)

	if size > 0 {
		if size > len(r) {
			return
		} else {
			r = r[0:size]
			return
		}
	}

	return
}
