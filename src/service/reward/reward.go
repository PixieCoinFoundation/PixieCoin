package reward

import (
	"appcfg"
	"constants"
	"fmt"
	"strings"
	// "sync"
	. "pixie_contract/api_specification"
	"time"
)

import (
	"common"
	jl "jsonloader"
	. "language"
	. "logger"
	"service/clothes"
	"service/mails"
	"xlsx"
)

var (
	AllRewards []*GFReward
	// allRewardsLock     sync.RWMutex
	IAPReturnRewardMap map[string]int
	OtherIAPReturns    map[string]int
)

func init() {
	if appcfg.GetServerType() != "" {
		return
	}

	loadRewards()
	loadIAPReturnRewards()

	// load7Reward()
}

// func load7Reward() {
// 	xlFile, err := xlsx.OpenFile("scripts/data/7reward.xlsx")
// 	if err != nil {
// 		panic(err)
// 	}

// 	for j, sheet := range xlFile.Sheets {
// 		if j == 0 {
// 			for i, row := range sheet.Rows {
// 				if i > 0 {
// 					var username, rewardCode string
// 					var diamondRank, goldRank, saleRank, passRank int
// 					var diamondRankReward, goldRankReward, saleRankReward, passRankReward string
// 					for j, cell := range row.Cells {
// 						if j == 0 {
// 							username, _ = cell.String()
// 							username = strings.TrimSpace(username)
// 						} else if j == 1 {
// 							rewardCode, _ = cell.String()
// 						} else if j == 2 {
// 							diamondRank, _ = cell.Int()
// 						} else if j == 3 {
// 							goldRank, _ = cell.Int()
// 						} else if j == 4 {
// 							passRank, _ = cell.Int()
// 						} else if j == 5 {
// 							saleRank, _ = cell.Int()
// 						} else if j == 6 {
// 							diamondRankReward, _ = cell.String()
// 						} else if j == 7 {
// 							goldRankReward, _ = cell.String()
// 						} else if j == 8 {
// 							passRankReward, _ = cell.String()
// 						} else if j == 9 {
// 							saleRankReward, _ = cell.String()
// 						}
// 					}
// 					// Info(un, rwd)
// 					// IAPReturnRewardMap[un] += rwd
// 					content := "设计师您好，7月设计师排行榜奖励已经开始发放，恭喜您进入获奖名单。你的排名信息："
// 					if diamondRank > 0 {
// 						if diamondRankReward != "" {
// 							content += fmt.Sprintf("钻石榜排名：%d，奖励：%s;", diamondRank, diamondRankReward)
// 						} else {
// 							panic(username + ",diamond")
// 						}
// 					}

// 					if goldRank > 0 {
// 						if goldRankReward != "" {
// 							content += fmt.Sprintf("金币榜排名：%d，奖励：%s;", goldRank, goldRankReward)
// 						} else {
// 							panic(username + ",gold")
// 						}
// 					}

// 					if saleRank > 0 {
// 						if saleRankReward != "" {
// 							content += fmt.Sprintf("销量榜排名：%d，奖励：%s;", saleRank, saleRankReward)
// 						} else {
// 							panic(username + ",sale")
// 						}
// 					}

// 					if passRank > 0 {
// 						if passRankReward != "" {
// 							content += fmt.Sprintf("高产榜排名：%d，奖励：%s;", passRank, passRankReward)
// 						} else {
// 							panic(username + ",pass")
// 						}
// 					}

// 					content += "您的领奖码：" + rewardCode + ";"
// 					content += "请到【妖精的衣橱手游】官方微博，将【领奖码，游戏昵称，游戏id，收货地址】私信给官博娘，我们将尽快给您邮寄奖励。"
// 					sql := fmt.Sprintf("insert into gf_mail values(0,0,'妖精管理局','%s','7月设计师排行榜奖励信息','%s',0,0,'[]',unix_timestamp(),0,0,1503244800);", username, content)
// 					fmt.Println(sql)
// 				}
// 			}
// 		}
// 	}
// 	panic("end")
// }

func GetIAPReturnDiamond(username string) int {
	if IAPReturnRewardMap[username] > 0 {
		return IAPReturnRewardMap[username]
	} else if OtherIAPReturns[username] > 0 {
		return OtherIAPReturns[username]
	}
	return 0
}

func SendRewardByPKPoints(username string, points int, isWeek bool) {
	title := L("pk10")
	if !isWeek {
		title = L("pk11")
	}
	level, t := common.GenPKLevel(points)
	var rm, rd int64

	if isWeek {
		if level == 1 {
			rm = 10000
		} else if level == 2 {
			rm = 15000
		} else if level == 3 {
			rm = 20000
		} else if level == 4 {
			rm = 25000
		} else if level == 5 {
			rd = 150
		} else if level == 6 {
			rd = 200
		} else if level == 7 {
			rd = 250
		} else if level == 8 {
			rd = 300
		} else if level == 9 {
			rd = 400
		} else if level == 10 {
			rd = 500
		}
	} else {
		if level == 1 {
			rm = 20000
		} else if level == 2 {
			rm = 30000
		} else if level == 3 {
			rm = 40000
		} else if level == 4 {
			rm = 50000
		} else if level == 5 {
			rd = 300
		} else if level == 6 {
			rd = 400
		} else if level == 7 {
			rd = 500
		} else if level == 8 {
			rd = 600
		} else if level == 9 {
			rd = 800
		} else if level == 10 {
			rd = 1000
		}
	}

	mails.SendToOne("", username, title+L("pk9"), fmt.Sprintf(L("pk12"), title, t), rd, rm, "", true, time.Now().Unix()+constants.DEFAULT_MAIL_EXPIRE_TIME)
}

func GetRewardByTag(tag string) (ret *GFReward) {
	// allRewardsLock.RLock()
	// defer func() {
	// 	allRewardsLock.RUnlock()
	// }()
	for _, reward := range AllRewards {
		if reward.Tag == tag {
			ret = reward
			break
		}
	}
	return
}

func GetRewardByID(rewardID int) (ret *GFReward) {
	// allRewardsLock.RLock()
	// defer func() {
	// 	allRewardsLock.RUnlock()
	// }()
	for _, reward := range AllRewards {
		if reward.RewardID == rewardID {
			ret = reward
			break
		}
	}
	return
}

func SendRewardByTag(tag string, to string) {
	// 根据tag获取奖励
	reward := GetRewardByTag(tag)
	if reward == nil {
		Info("reward tag:", tag, "not exist")
		return
	}

	if !strings.HasPrefix(tag, "vip") {
		clothesID := reward.ClothesID
		// 判断服装是否为-1: 随机分配一套vip的服装
		if reward.ClothesID == "-1" {
			clothesID = clothes.GetRandomVIPClothes()
		} else if reward.ClothesID == "-2" {
			// 判断服装是否为-2: 随机分配一套cos的服装
			clothesID = clothes.GetRandomCOSClothes()
		}

		if clothesID != "" {
			// 给玩家发送邮件
			mails.SendToOne(reward.From, to, reward.Title, reward.Desc, int64(reward.Diamond), int64(reward.Gold), clothesID, true, 0)
		}
	} else {
		var clos []string
		if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
			if tag == "vip1" {
				clos = clothes.GetSuitClothesList("c900101000000000010")
			} else if tag == "vip2" {
				clos = clothes.GetSuitClothesList("c900101000000000062")
			} else if tag == "vip3" {
				clos = clothes.GetSuitClothesList("c900101000000000015")
			} else if tag == "vip4" {
				clos = clothes.GetSuitClothesList("c900101000000000017")
			} else if tag == "vip5" {
				clos = clothes.GetSuitClothesList("c900101000000000006")
			} else if tag == "vip6" {
				clos = clothes.GetSuitClothesList("c900101000000000051")
			} else if tag == "vip7" {
				clos = clothes.GetSuitClothesList("c900101000000000064")
			} else if tag == "vip8" {
				clos = clothes.GetSuitClothesList("c900101000000000054")
			} else if tag == "vip9" {
				clos = clothes.GetSuitClothesList("c900101000000000052")
			}
		} else {
			if tag == "vip1" {
				clos = clothes.GetSuitClothesList("c900101000000000010")
			} else if tag == "vip2" {
				clos = clothes.GetSuitClothesList("c900101000000000014")
			} else if tag == "vip3" {
				clos = clothes.GetSuitClothesList("c900101000000000015")
			} else if tag == "vip4" {
				clos = clothes.GetSuitClothesList("c900101000000000017")
			} else if tag == "vip5" {
				clos = clothes.GetSuitClothesList("c900101000000000006")
			} else if tag == "vip6" {
				clos = clothes.GetSuitClothesList("c900101000000000051")
			} else if tag == "vip7" {
				clos = clothes.GetSuitClothesList("c900101000000000053")
			} else if tag == "vip8" {
				clos = clothes.GetSuitClothesList("c900101000000000054")
			} else if tag == "vip9" {
				clos = clothes.GetSuitClothesList("c900101000000000052")
			}
		}

		if len(clos) > 0 {
			rcs := make([]ClothesInfo, 0)
			for _, c := range clos {
				ci := ClothesInfo{
					ClothesID: c,
					Count:     1,
				}
				rcs = append(rcs, ci)
			}
			mails.SendToOneCT(reward.From, to, reward.Title, reward.Desc, int64(reward.Diamond), int64(reward.Gold), rcs, true, 0)
		}
	}
}

func SendRewardByID(id int, to string) {
	// 根据tag获取奖励
	reward := GetRewardByID(id)

	clothesID := reward.ClothesID
	// 判断服装是否为-1: 随机分配一套vip的服装
	if reward.ClothesID == "-1" {
		clothesID = clothes.GetRandomVIPClothes()
	}
	// 给玩家发送邮件
	if exist, _ := clothes.ClothesLegal(clothesID); !exist {
		Err("clothes illegal:", clothesID)
	} else {
		err := mails.SendToOne(reward.From, to, reward.Title, reward.Desc, int64(reward.Diamond), int64(reward.Gold), clothesID, true, 0)
		if err != nil {
			fmt.Println("发送奖励错误：", err.Error())
		}
	}
}

func loadRewards() {
	// allRewardsLock.Lock()
	// defer allRewardsLock.Unlock()
	f := "scripts/data/data_reward.json"

	if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
		f = "scripts/data/data_reward_kor.json"
	}

	if err := jl.LoadFile(f, &AllRewards); err != nil {
		panic(err)
	}
	return
}

func loadIAPReturnRewards() {
	OtherIAPReturns = make(map[string]int)
	if err := jl.LoadFile("scripts/data/other_iap_return.json", &OtherIAPReturns); err != nil {
		panic(err)
	}

	IAPReturnRewardMap = make(map[string]int)

	xlFile, err := xlsx.OpenFile("scripts/data/data_iap_return.xlsx")
	if err != nil {
		panic(err)
	}

	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i > 0 {
				var un string
				var rwd int
				for j, cell := range row.Cells {
					if j == 0 {
						//bili sdk username
						un, _ = cell.String()
						un = strings.TrimSpace(un)
					} else if j == 1 {
						rwd, _ = cell.Int()
					}
				}
				// Info(un, rwd)
				IAPReturnRewardMap[un] += rwd
			}
		}
	}

	xlFile, err = xlsx.OpenFile("scripts/data/data_iap_return_1.xlsx")
	if err != nil {
		panic(err)
	}

	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i > 0 {
				var un string
				var rwd int
				for j, cell := range row.Cells {
					if j == 0 {
						//bili sdk username
						un, _ = cell.String()
						un = strings.TrimSpace(un)
					} else if j == 1 {
						rwd, _ = cell.Int()
					}
				}
				// Info(un, rwd)
				IAPReturnRewardMap[un] += rwd
			}
		}
	}

	Info("iap return rwd size:", len(IAPReturnRewardMap))
}
