package clean

import (
	"common"
	"constants"
	"db/board"
	"db/cosplay"
	"db/cosplay/items"
	"db/guild"
	"db/hope"
	"db/mail"
	"encoding/json"
	"fmt"
	. "language"
	. "logger"
	"rao"
	"service/files"
	"service/mails"
	"time"
	"tools"
	. "types"
)

const (
	HELP_PRESERVE_SIZE = 1
	MAX_MAIL_TIME      = 99 * 24 * 3600 //99 day in seconds
)

// func CleanCustomBuy() {
// 	et := time.Now().Unix() - 15*24*3600
// 	if delID := custom.GetExpireCustomBuyID(100, et); delID > 0 {
// 		Info("del custom buy id<=", delID)
// 		custom.BatchDelCustomBuyLog(delID)
// 	}
// }

//CleadKorlog 清除韩国日志
// func CleadKorlog() {
// 	et := time.Now().Unix() - 7*24*3600
// 	if timePoint := orders.QueryKorOrderlog(100, et); len(timePoint) > 0 {
// 		orders.CleanKorOrderLog(timePoint)
// 	}
// 	if timePoint := kor.QueryKorLoginlog(100, et); len(timePoint) > 0 {
// 		kor.CleanKorLoginLog(timePoint)
// 	}
// 	if timePoint := kor.QueryKorLogin1log(100, et); len(timePoint) > 0 {
// 		kor.CleanKorLogin1Log(timePoint)
// 	}
// }
// func CleanDesignMsg() {
// 	conn := rao.GetConn()
// 	defer conn.Close()

// 	et := time.Now().Unix() - 7*24*3600
// 	if msgs, err := custom.GetDesignMsgDBForClean(100); err != nil {
// 		return
// 	} else {
// 		for _, m := range msgs {
// 			if m.Time < et {
// 				Info("del design msg:", m.ID)
// 				conn.Send("LREM", constants.GLOBAL_DESIGN_MSG_LIST_KEY, 0, fmt.Sprintf("%d:%s", m.ID, m.DesignerUsername))
// 				conn.Send("DEL", rao.GetDesignMsgSelfKey(m.ID))

// 				//new list
// 				timeToken := tools.GetDesignMsgTimeToken(time.Unix(m.Time, 0))
// 				listKeyValue := rao.GetUserDesignMsgListKey(m.DesignerUsername, m.Type, timeToken)
// 				dmik := rao.GetUserDesignMsgTimeInfoKey(listKeyValue)

// 				conn.Send("ZREM", rao.GetNewDesignMsgListkey(), listKeyValue)
// 				conn.Send("DEL", listKeyValue)
// 				conn.Send("DEL", dmik)

// 				conn.Flush()
// 				custom.DelDesignMsg(m.ID)
// 			}
// 		}
// 	}
// }

func CleanApplyGuild() {
	ext := time.Now().Unix() - 24*3600
	if res, err := guild.ListApply(0, 50); err != nil {
		return
	} else if len(res) > 0 {
		for _, a := range res {
			if a.ApplyTime <= ext {
				Info("delete guild apply:", a.ID)
				guild.DelApply(a.ID)
			}
		}
	}
}

func CleanMail() {
	now := time.Now().Unix()
	if res, err := mail.List(0, 50); err != nil {
		return
	} else if len(res) > 0 {
		for _, m := range res {
			if (m.ExpireTime > 0 && m.ExpireTime <= now) || (m.Time > 0 && now-m.Time >= MAX_MAIL_TIME) {
				Info("delete mail:", m.To, m.MailID)
				mail.DeleteMail(m.To, m.MailID)
				continue
			}
		}
	}
}

func CleanHope() {
	now := time.Now().Unix()
	expireTime := now - 7*3600*24 //7 day

	if rs, err := hope.SelectHope(1000); err != nil {
		return
	} else {
		for _, hp := range rs {
			if hp.SendTime < expireTime {
				if hp.Status == constants.HOPE_STATUS_SENDED {
					if tools.ClothesIDLegal(hp.OfferClothes) {
						if ok, _ := hope.UpdateHopeStatus(hp.ID, constants.HOPE_STATUS_CLEAN, constants.HOPE_STATUS_SENDED); ok {
							//send clo piece
							mc := MailContent{
								Content:   L("cp1"),
								CloPieces: map[string]int{fmt.Sprintf("%s:%d", hp.OfferClothes, hp.OfferPart): hp.OfferNum},
							}
							mcb, _ := json.Marshal(mc)
							mails.SendToOne("", hp.Sender, L("cp2"), string(mcb), 0, 0, "", true, now+constants.DEFAULT_MAIL_EXPIRE_TIME)
						}
					} else {
						Info("clean SENDED hope:", hp.ID)
						hope.DelHope(hp.ID)
					}
				} else if hp.Status == constants.HOPE_STATUS_DONE {
					if tools.ClothesIDLegal(hp.NeedClothes) {
						if ok, _ := hope.UpdateHopeStatus(hp.ID, constants.HOPE_STATUS_CLEAN, constants.HOPE_STATUS_DONE); ok {
							//send clo piece
							mc := MailContent{
								Content:   L("cp3"),
								CloPieces: map[string]int{fmt.Sprintf("%s:%d", hp.NeedClothes, hp.NeedPart): hp.NeedNum},
							}
							mcb, _ := json.Marshal(mc)
							mails.SendToOne("", hp.Sender, L("cp4"), string(mcb), 0, 0, "", true, now+constants.DEFAULT_MAIL_EXPIRE_TIME)
						}
					} else {
						Info("clean DONE hope:", hp.ID)
						hope.DelHope(hp.ID)
					}
				} else if hp.Status == constants.HOPE_STATUS_DONE_READ {
					Info("clean DONE_READ hope:", hp.ID)
					hope.DelHope(hp.ID)
				} else if hp.Status == constants.HOPE_STATUS_CLEAN {
					Info("clean CLEAN hope:", hp.ID)
					hope.DelHope(hp.ID)
				}
			}
		}
	}
}

func CleanBoard() {
	var mid int64
	expireTime := time.Now().Unix() - 24*3600*32*6 //one month

	if res, err := board.BatchQueryBoardMsgs(0, 10000); err != nil {
		return
	} else {
		for _, bm := range res {
			if bm.Time <= expireTime {
				mid = bm.ID
			} else {
				break
			}
		}

		if mid > 0 {
			board.BatchDelBoardMsgs(mid)
		}
	}
}

func CleanCosplay() {
	//32 day
	expireTime := time.Now().Unix() - constants.COSPLAY_LAST_SECOND
	stop := false
	start := 0
	size := 1000
	conn := rao.GetConn()
	defer conn.Close()
	for !stop {
		if res, err := cosplay.GetOldCosplayFromDB(start, size); err != nil {
			return
		} else if len(res) > 0 {
			for _, c := range res {
				if c.CloseTime > expireTime {
					stop = true
					break
				} else if c.Status == COS_CLOSED {
					Info("delete cosplay id:", c.CosplayID, "file:", c.Icon, c.CosBg, c.ListBg)
					//delete in db
					cosplay.DeleteCosplay(c.CosplayID)

					//delete in redis
					if s := common.GetCosplayStrInRedis(c.CosplayID, c.Title, c.Keyword, c.Type, c.OpenTime, c.CloseTime, c.AdminUsername, c.AdminNickname, c.Icon, c.CosBg, c.ListBg, c.Params); s != "" {
						conn.Send("ZREM", constants.COSPLAY_TYPE_LIST_KEY, s)
					}
					conn.Send("DEL", rao.GetCosItemRankSetKey(c.CosplayID))

					//delete cosplay item
					cleanCosplayItem(c.CosplayID)

					//delete cosplay item cmt
					cosplay.DeleteCosComment(c.CosplayID)

					//delete file
					files.DeleteFile(c.Icon, "")
					files.DeleteFile(c.CosBg, "")
					files.DeleteFile(c.ListBg, "")
				}
			}
		} else {
			stop = true
		}
		start += size
	}

	if err := conn.Flush(); err != nil {
		Err(err)
	}
}

func cleanCosplayItem(cosplayID int64) {
	stop := false
	start := 0
	size := 1000
	conn := rao.GetConn()
	defer conn.Close()
	for !stop {
		if res, err := items.GetCosItemByCosplay(cosplayID, start, size); err != nil {
			return
		} else if len(res) > 0 {
			dc := 0
			for _, fc := range res {
				dc++
				Info("delete cos item id:", fc.ItemID, "icon file:", fc.Icon, "img file:", fc.Img)
				//delete in db
				items.DeleteCosItem(fc.ItemID)

				//delete in redis
				conn.Send("DEL", rao.GetCosItemUserListKey(cosplayID, fc.Username))
				conn.Send("DEL", rao.GetCosItemCmtListKey(fc.ItemID))
				conn.Send("DEL", rao.GetCosItemHasCmtKey(fc.ItemID))
				conn.Send("DEL", rao.GetCosItemSelfKey(fc.ItemID))
				conn.Send("DEL", rao.GetCosItemUserMarkListKey(cosplayID, fc.ItemID))

				//delete file
				files.DeleteFile(fc.Icon, "")
				files.DeleteFile(fc.Img, "")
			}

			if dc == 0 {
				stop = true
			}
		} else {
			stop = true
		}
		start += size
	}

	if err := conn.Flush(); err != nil {
		Err(err)
	}
}
