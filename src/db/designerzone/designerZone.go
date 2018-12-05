package designerzone

import (
	// "appcfg"
	"constants"
	"dao_designer"
	"database/sql"
	"db/custom"
	. "logger"
	"math/rand"
	"rao"
	"time"
	"tools"
	. "types"
)

//获取衣服信息
func GetCustomIDListFromCustomPool() (customMap map[string]*CustomPool, newcustomMap map[string]*CustomPool, recommendCustomList []*CustomPool, recommendNewCustomList []*CustomPool, publicListArray []*CustomPool, err error) {
	var rows *sql.Rows
	if rows, err = dao_designer.GetCustomListFromPoolStmt.Query(); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("getCustomListFromPool can't find data in db:")
			return
		} else {
			Err(err)
			return
		}
	} else {
		defer rows.Close()
		conn := rao.GetConn()
		defer conn.Close()

		customMap = make(map[string]*CustomPool)
		newcustomMap = make(map[string]*CustomPool)
		recommendCustomList = make([]*CustomPool, 0)
		recommendNewCustomList = make([]*CustomPool, 0)
		publicListArray = make([]*CustomPool, 0)
		tmpPublicListArray := make([]*CustomPool, 0)

		for rows.Next() {
			var r CustomPool
			var id int64
			var entTime string
			if err = rows.Scan(&id, &r.ClothesIDInGame, &entTime, &r.Poolstatus, &r.Topic, &r.IsStick); err != nil {
				Err(err)
				return
			}

			_, _, _, cid := tools.GetInfoFromCustomIDInGame(r.ClothesIDInGame)
			r.Custom, _ = custom.GetOneCustom(cid, &conn, false)

			if r.Custom.Status == PASSED || r.Custom.Status == PASSED_NO_INVENTORY {
				if r.Username != "" {
					tm2, _ := time.ParseInLocation("20060102 150405", entTime, time.Local)
					r.EntryTime = tm2.Unix()

					if r.Topic > 0 {
						tmpPublicListArray = append(tmpPublicListArray, &r)
					}

					if r.Poolstatus == constants.CLOTHESPOOL {
						if r.IsStick {
							recommendCustomList = append(recommendCustomList, &r)
						} else {
							customMap[r.ClothesIDInGame] = &r
						}
					} else {
						if r.IsStick {
							recommendNewCustomList = append(recommendNewCustomList, &r)
						} else {
							newcustomMap[r.ClothesIDInGame] = &r
						}
					}
				}
			}
		}

		publicSize := len(tmpPublicListArray)
		if publicSize > 0 {
			gotMap := make(map[int]int)
			for i := 0; i < 8; i++ {
				r := rand.Intn(publicSize)
				if gotMap[r] <= 0 {
					gotMap[r] = 1
					publicListArray = append(publicListArray, tmpPublicListArray[r])
				}

				if len(publicListArray) >= 4 {
					break
				}
			}
		}
	}

	return

	// customMap = make(map[string]*CustomPool)
	// newcustomMap = make(map[string]*CustomPool)
	// recommendCustomList = make([]*CustomPool, 0)
	// recommendNewCustomList = make([]*CustomPool, 0)
	// list, publicList, e := getCustomListFromPool()
	// for _, v := range publicList {
	// 	if v.Username != "" {
	// 		publicListArray = append(publicListArray, *v)
	// 	}
	// }
	// if e != nil {
	// 	err = e
	// 	return
	// }

	// for _, v := range list {
	// 	var temp CustomPool
	// 	if v.Poolstatus == constants.CLOTHESPOOL {
	// 		temp.Custom = v.Custom
	// 		temp.Poolstatus = v.Poolstatus
	// 		temp.EntryTime = v.EntryTime
	// 		temp.IsStick = v.IsStick
	// 		temp.ClothesIDInGame = tools.GetCustomIDInGame(v.Type, v.ModelNo, constants.ORIGIN_DESIGNER, v.CID)
	// 		if v.IsStick {
	// 			recommendCustomList = append(recommendCustomList, &temp)
	// 		} else {
	// 			customMap[v.ClothesIDInGame] = &temp
	// 		}
	// 	} else {
	// 		temp.Custom = v.Custom
	// 		temp.Poolstatus = v.Poolstatus
	// 		temp.EntryTime = v.EntryTime
	// 		temp.IsStick = v.IsStick
	// 		temp.ClothesIDInGame = tools.GetCustomIDInGame(v.Type, v.ModelNo, constants.ORIGIN_DESIGNER, v.CID)

	// 		if v.IsStick {
	// 			recommendNewCustomList = append(recommendNewCustomList, &temp)
	// 		} else {
	// 			newcustomMap[v.ClothesIDInGame] = &temp
	// 		}
	// 	}
	// }
	// 	return
}

func GetDesignerIDListFromDesignerPool() ([]*GFDesignerPool, []*GFDesignerPool, []*GFDesignerPool, []*GFDesignerPool, error) {
	desingerL := make([]*GFDesignerPool, 0)
	newdesingerL := make([]*GFDesignerPool, 0)

	recommendDesinger := make([]*GFDesignerPool, 0)
	newRecommendDesigner := make([]*GFDesignerPool, 0)

	list, _ := getDesignerListFromPool()
	for _, v := range list {
		var temp GFDesignerPool
		if v.Poolstatus == constants.DESIGNERPOOL {
			//优秀设计师
			temp.GFDesigner = v.GFDesigner
			temp.Poolstatus = v.Poolstatus
			temp.IsStick = v.IsStick
			temp.ClothesList, _ = GetRecommendOrderByCount(v.GFDesigner.Username)

			if v.IsStick {
				recommendDesinger = append(recommendDesinger, &temp)
			} else {
				desingerL = append(desingerL, &temp)
			}
		} else {
			//新人设计师
			temp.GFDesigner = v.GFDesigner
			temp.Poolstatus = v.Poolstatus
			temp.IsStick = v.IsStick
			temp.IsNewDesigner = true
			temp.ClothesList, _ = GetRecommendOrderByCount(v.GFDesigner.Username)
			if v.IsStick {
				newRecommendDesigner = append(newRecommendDesigner, &temp)
			} else {
				newdesingerL = append(newdesingerL, &temp)
			}
		}
	}
	return desingerL, newdesingerL, recommendDesinger, newRecommendDesigner, nil
}
func GetBannerListFromCBannerPool(tNow time.Time) (list []DesignerZoneBanner, err error) {
	return getBannerListFromPool(tNow)
}
func GetTopicListFromTopicPool(tNow time.Time) (list []Topic, stick []Topic, err error) {
	return getTopicListFromPool(tNow)
}

//获取专题详情页
func GetTopicCustomList(topicID int) ([]CustomPool, error) {
	// list, _ := getTopicCustomList(topicID)
	// size := len(list)
	// if size <= 0 {
	// 	return list
	// }
	// for i := 0; i < size; i++ {

	// }
	return getTopicCustomList(topicID)
}
func GetRecommendOrderByCount(username string) (list []Custom, err error) {
	return getRecommendOrderByBuyCountList(username)
}
