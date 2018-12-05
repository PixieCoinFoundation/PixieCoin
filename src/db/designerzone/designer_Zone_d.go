package designerzone

import (
	// "appcfg"
	"constants"
	"dao_designer"
	"database/sql"
	"encoding/json"
	// "db/custom"
	// "fmt"
	// "rao"
	. "logger"
	"math/rand"
	"time"
	"tools"
)

import (
	. "types"
)

// func getCustomListFromPool() (customMap map[string]*CustomPool, newcustomMap map[string]*CustomPool, recommendCustomList []*CustomPool, recommendNewCustomList []*CustomPool, publicListArray []*CustomPool,err error) {
// var rows *sql.Rows
// if rows, err = dao_designer.GetCustomListFromPoolStmt.Query(); err != nil {
// 	if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
// 		Info("getCustomListFromPool can't find data in db:")
// 		return
// 	} else {
// 		Err(err)
// 		return
// 	}
// } else {
// 	defer rows.Close()
// 	conn := rao.GetConn()
// 	defer conn.Close

// 	customMap = make(map[string]*CustomPool)
// 	newcustomMap = make(map[string]*CustomPool)
// 	recommendCustomList = make([]*CustomPool, 0)
// 	recommendNewCustomList = make([]*CustomPool, 0)
// 	publicListArray = make([]*CustomPool, 0)

// 	for rows.Next() {
// 		var r CustomPool
// 		var id int64
// 		var entTime string
// 		if err = rows.Scan(&id, &r.ClothesIDInGame, &entTime, &r.Poolstatus, &r.Topic, &r.IsStick); err != nil {
// 			Err(err)
// 			return
// 		}

// 		_, _, _, cid := tools.GetInfoFromCustomIDInGame(clothesID)
// 		r.Custom, _ = custom.GetOneCustom(cid, &conn, false)

// 		if r.Username!=""{
// 			tm2, _ := time.ParseInLocation("20060102 150405", entTime, time.Local)
// 			r.EntryTime = tm2.Unix()

// 			if r.Topic>0{
// 				publicListArray = append(publicListArray,&r)
// 			}

// 			if r.Poolstatus == constants.CLOTHESPOOL{
// 				if r.IsStick {
// 					recommendCustomList = append(recommendCustomList, &r)
// 				} else {
// 					customMap[r.ClothesIDInGame] = &r
// 				}
// 			}else{
// 				if v.IsStick {
// 					recommendNewCustomList = append(recommendNewCustomList, &r)
// 				} else {
// 					newcustomMap[r.ClothesIDInGame] = &r
// 				}
// 			}
// 		}
// 	}

//批量获取衣服
// customList := []int64{}
// for cid, _ := range list {
// 	customList = append(customList, cid)

// 	if len(customList)%10 == 0 {
// 		var rows10 *sql.Rows
// 		if rows10, err = dao_designer.GetCustomByIn10Stmt.Query(customList[0], customList[1], customList[2], customList[3], customList[4], customList[5], customList[6], customList[7], customList[8], customList[9]); err != nil {
// 			if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
// 				Info("getCustomListFromPool can't find data in db:")
// 				return list, publicBuyList, err
// 			} else {
// 				Err(err)
// 				return list, publicBuyList, err
// 			}
// 		} else {
// 			for rows10.Next() {
// 				var clothesSet string
// 				var head string
// 				var extra string
// 				var r Custom
// 				err = rows10.Scan(&r.CID, &r.Username, &r.Nickname, &r.ModelNo, &r.Type, &r.CName, &r.Desc, &clothesSet, &r.Status, &r.UploadTime, &r.CheckTime, &r.AdminUsername, &r.MoneyType, &r.MoneyAmount, &r.Hearts, &r.Info, &r.BuyCount, &extra, &r.Inventory, &head, &r.Remark, &r.Tag1, &r.Tag2, &r.GoldProfit, &r.DiamondProfit, &r.Exp, &r.PlayerConfig)

// 				if err != nil {
// 						Err(err)
// 						return
// 				}
// 				if err = json.Unmarshal([]byte(clothesSet), &r.CloSet); err != nil {
// 					Err(err)
// 					continue
// 				}
// 				if err = json.Unmarshal([]byte(extra), &r.Extra); err != nil {
// 					Err(err)
// 					continue
// 				}
// 				ClothesIDInGame := tools.GetCustomIDInGame(r.Type, r.ModelNo, constants.ORIGIN_DESIGNER, r.CID)
// 				r.ClothesIDInGame = ClothesIDInGame
// 				if r.Username != "" {
// 					list[r.CID].Custom = r
// 				}

// 			}
// 			customList = make([]int64, 0)
// 		}
// 	}
// }

// //单个获取衣服
// for _, cid := range customList {
// 	if c, exist := list[cid]; exist {
// 		var cloSetStr string
// 		var extra string
// 		if err = dao_designer.GetCustomByCIDStmt.QueryRow(cid).Scan(&c.CID, &c.Username, &c.Nickname, &c.ModelNo, &c.Type, &c.CName, &c.Desc, &cloSetStr, &c.Status, &c.UploadTime, &c.CheckTime, &c.AdminUsername, &c.MoneyType, &c.MoneyAmount, &c.Hearts, &c.Info, &c.BuyCount, &extra, &c.Inventory, &c.NewSold, &c.StickTime); err != nil {
// 			if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
// 				Info("getCustomListFromPool can't find data loop custom in db:", c.CID)
// 				continue
// 			} else {
// 				Err(err)
// 				return list, publicBuyList, err
// 			}
// 		} else {
// 			if err = json.Unmarshal([]byte(cloSetStr), &c.CloSet); err != nil {
// 				Err(err)
// 				continue
// 			}
// 			if err = json.Unmarshal([]byte(extra), &c.Extra); err != nil {
// 				Err(err)
// 				continue
// 			}
// 			if c.Username != "" {
// 				list[cid] = c
// 			}
// 		}
// 	}
// 	customList = make([]int64, 0)
// }

// for _, v := range list {
// 	if v.Username != "" {
// 		if v.Topic > 0 {
// 			publicBuyList[v.ClothesIDInGame] = v
// 		} else {
// 			result[v.CID] = v
// 		}
// 	}
// }
// 	}
// 	return
// }

func getDesignerListFromPool() (list map[string]*GFDesignerPool, err error) {
	var rows *sql.Rows
	if rows, err = dao_designer.GetDesignerListFromPoolStmt.Query(); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("getDesignerListFromPool can find data in db")
			return
		} else {
			Err(err)
			return
		}
	} else {
		defer rows.Close()
		//获取推荐设计师列表
		list = make(map[string]*GFDesignerPool, 0)
		for rows.Next() {
			var username string
			var id int64
			var designer GFDesignerPool
			var entryTime string
			if err = rows.Scan(&id, &username, &entryTime, &designer.Poolstatus, &designer.IsStick); err != nil {
				Err(err)
				return
			}
			tm2, _ := time.ParseInLocation("2006-01-02 15:04:05", entryTime, time.Local)
			designer.EntryTime = tm2.Unix()
			list[username] = &designer
		}

		designerList := make([]string, 0)
		for username, _ := range list {
			designerList = append(designerList, username)
			//批量获取设计师衣服
			if len(designerList)%10 == 0 {
				var rows10 *sql.Rows
				if rows10, err = dao_designer.GetDesinger10Stmt.Query(designerList[0], designerList[1], designerList[2], designerList[3], designerList[4], designerList[5], designerList[6], designerList[7], designerList[8], designerList[9]); err != nil {
					if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
						Info("not config designer info", username)
						return
					} else {
						Err(err)
						return
					}
				} else {
					defer rows10.Close()
					for rows10.Next() {
						var id int64
						g := GFDesigner{}
						if err = rows10.Scan(&id, &(g.Username), &(g.Nickname), &(g.Diamond), (&g.Gold), &(g.TotalDiamond), &(g.TotalGold), &(g.TotalSale), &(g.Points), &(g.PassedCount), &(g.TotalCount), &(g.FensiCount), &(g.Desc), &(g.Ad_DayReward), &(g.Ad_MonthOneReward), &(g.Ad_AutoConfig), &(g.Ad_MonthTenReward), &(g.Ad_MonthThirtyReward), &(g.Head), &(g.DesignCoin), &(g.TotalDesignCoin), &(g.DesignCoinPriceTokenCount)); err != nil {
							Err(err)
							return
						}

						list[g.Username].GFDesigner = g
					}
				}

				designerList = make([]string, 0)
			}
		}

		//单个获取设计师
		for _, username := range designerList {
			if g, exist := list[username]; exist {
				var id int64
				if err = dao_designer.GetDesignerStmt.QueryRow(username).Scan(&id, &(g.GFDesigner.Username), &(g.GFDesigner.Nickname), &(g.Diamond), (&g.Gold), &(g.TotalDiamond), &(g.TotalGold), &(g.TotalSale), &(g.Points), &(g.PassedCount), &(g.TotalCount), &(g.FensiCount), &(g.GFDesigner.Desc), &(g.Ad_DayReward), &(g.Ad_MonthOneReward), &(g.Ad_AutoConfig), &(g.Ad_MonthTenReward), &(g.Ad_MonthThirtyReward), &(g.Head), &(g.DesignCoin), &(g.TotalDesignCoin), &(g.DesignCoinPriceTokenCount), &g.ETHAccount); err != nil {
					if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
						Info("GetDesignerStmt----------", username)
						return
					} else {
						Err(err)
						return
					}
				}
			}
		}
	}
	return
}

//GetTopicCustomList 获取专题list
func getTopicCustomList(topicID int) (res []CustomPool, err error) {
	var rows *sql.Rows
	if rows, err = dao_designer.GetTopicCustomListStmt.Query(topicID); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info(" getTopicCustomList cannot find data in db")
			return
		} else {
			Err(err)
			return
		}
	} else {
		defer rows.Close()
		list := make(map[string]*CustomPool, 0)
		customList := []string{}
		//解析所有池某个专题的衣服ID
		for rows.Next() {
			var clothesID string
			var r CustomPool
			var poolstatus int
			var id int
			var entTime string
			var stick bool
			var topicID int
			if err = rows.Scan(&id, &clothesID, &entTime, &poolstatus, &stick, &topicID); err != nil {
				Err(err)
				return
			}
			r.IsStick = stick
			r.Poolstatus = poolstatus
			tm2, _ := time.ParseInLocation("2006-01-02 15:04:05", entTime, time.Local)
			r.EntryTime = tm2.Unix()
			r.ClothesIDInGame = clothesID
			list[clothesID] = &r
			customList = append(customList, clothesID)
		}
		//获取某个专题的衣服列表
		for _, clothesIDSrt := range customList {
			if c, exist := list[clothesIDSrt]; exist {
				var cloSetStr string
				var extra string
				_, _, _, id := tools.GetInfoFromCustomIDInGame(clothesIDSrt)
				if err = dao_designer.GetCustomByCIDStmt.QueryRow(id).Scan(&c.CID, &c.Username, &c.Nickname, &c.ModelNo, &c.Type, &c.CName, &c.Desc, &cloSetStr, &c.Status, &c.UploadTime, &c.CheckTime, &c.AdminUsername, &c.MoneyType, &c.MoneyAmount, &c.Hearts, &c.Info, &c.BuyCount, &extra, &c.Inventory, &c.NewSold, &c.StickTime); err != nil {
					if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
						Info("getTopicCustomList can't find data in db:", id)
						return res, err
					} else {
						Err(err)
						return res, err
					}
				} else {
					if err = json.Unmarshal([]byte(cloSetStr), &c.CloSet); err != nil {
						Err(err)
						return
					}
					if err = json.Unmarshal([]byte(extra), &c.Extra); err != nil {
						Err(err)
						return
					}
				}
				res = append(res, *c)
			}
		}
	}
	return
}

func getBannerListFromPool(tNow time.Time) (list []DesignerZoneBanner, err error) {
	var rows *sql.Rows
	if rows, err = dao_designer.GetBannerListFromPoolStmt.Query(); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("getBannerListFromPool can not find data in db")
			return
		} else {
			Err(err)
			return
		}
	} else {
		defer rows.Close()
		list = make([]DesignerZoneBanner, 0)
		for rows.Next() {
			var temp DesignerZoneBanner
			if err = rows.Scan(&temp.ID, &temp.BannerTitle, &temp.BannerType, &temp.URL, &temp.IndexOrRecom, &temp.Shunxu, &temp.IsStick, &temp.CreateTime, &temp.StartTime, &temp.EndTime, &temp.IMGURL); err != nil {
				Err(err)
				return
			}
			//获取所有banner信息
			tmSt, _ := time.ParseInLocation("2006-01-02 15:04:05", temp.StartTime, time.Local)
			tmEnd, _ := time.ParseInLocation("2006-01-02 15:04:05", temp.EndTime, time.Local)
			if tNow.Unix() > tmSt.Unix() && tNow.Unix() < tmEnd.Unix() {
				list = append(list, temp)
			}
		}
	}
	return
}

func getTopicListFromPool(tNow time.Time) (list []Topic, stick []Topic, err error) {
	var rows *sql.Rows
	if rows, err = dao_designer.GetTopicListFromPoolStmt.Query(); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("getTopicListFromPool can not find data in db")
			return
		} else {
			Err(err)
			return
		}
	} else {
		defer rows.Close()
		list = make([]Topic, 0)
		for rows.Next() {
			var temp Topic
			if err = rows.Scan(&temp.ID, &temp.BannerTitle, &temp.OwerTitle, &temp.URL, &temp.ClothesCount, &temp.Shunxu, &temp.IsStick, &temp.CreateTime, &temp.StartTime, &temp.EndTime, &temp.TitleIcon, &temp.MainIcon, &temp.SliderIcon); err != nil {
				Err(err)
				return
			}
			tmSt, _ := time.ParseInLocation("2006-01-02 15:04:05", temp.StartTime, time.Local)
			tmEnd, _ := time.ParseInLocation("2006-01-02 15:04:05", temp.EndTime, time.Local)
			//获取置顶专题和非置顶专题
			if tNow.Unix() > tmSt.Unix() && tNow.Unix() < tmEnd.Unix() {
				if temp.IsStick {
					stick = append(stick, temp)
				}
				list = append(list, temp)
			}
		}
		return list, stick, nil
	}
}

//获取设计师推荐衣服
func getRecommendOrderByBuyCountList(username string) (list []Custom, err error) {
	var rows *sql.Rows
	if rows, err = dao_designer.GetRecommendOrderByBuyCount.Query(username); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("getRecommendOrderByBuyCountList can find  data in db", username)
			return
		} else {
			Err(err)
			return
		}
	} else {
		defer rows.Close()

		medlist := make([]Custom, 0)
		list = make([]Custom, 0)
		for rows.Next() {
			var c Custom
			var clothesSet string
			var extra string
			err = rows.Scan(&c.CID, &c.Username, &c.Nickname, &c.ModelNo, &c.Type, &c.CName, &c.Desc, &clothesSet, &c.Status, &c.UploadTime, &c.CheckTime, &c.AdminUsername, &c.MoneyType, &c.MoneyAmount, &c.Hearts, &c.Info, &c.BuyCount, &extra, &c.Inventory, &c.NewSold, &c.StickTime)
			if err != nil {
				Err(err)
				return
			}

			if err = json.Unmarshal([]byte(clothesSet), &c.CloSet); err != nil {
				Err(err)
				continue
			}
			if err = json.Unmarshal([]byte(extra), &c.Extra); err != nil {
				Err(err)
				continue
			}
			c.ClothesIDInGame = tools.GetCustomIDInGame(c.Type, c.ModelNo, constants.ORIGIN_DESIGNER, c.CID)
			medlist = append(medlist, c)
		}

		//随机取出三件衣服
		size := len(medlist)
		gotMap := make(map[int]int)
		if size <= 3 {
			list = medlist
		} else {
			for i := 0; i < 6; i++ {
				index := rand.Intn(size)
				if gotMap[index] <= 0 {
					gotMap[index] = 1
					list = append(list, medlist[index])
				}
			}
			// var i, index1, index2, index3 int
			// for {
			// 	index1 = rand.Intn(size)
			// 	index2 = rand.Intn(size)
			// 	index3 = rand.Intn(size)
			// 	i++
			// 	if index1 != index2 && index2 != index3 && index1 != index3 {
			// 		break
			// 	}
			// 	if i >= 100 {
			// 		break
			// 	}
			// }
			// list = append(list, medlist[index1], medlist[index2], medlist[index3])
		}
		return
	}
}
