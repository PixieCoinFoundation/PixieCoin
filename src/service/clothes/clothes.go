package clothes

import (
	"appcfg"
	"constants"
	// "encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	// "sync"
	"time"
	"tools"
	. "types"
	"xlsx"
)

import (
	. "logger"
	. "model"
	// . "pixie_contract/api_specification"
)

var (
	AllSuits    []*GFSuit // 所有的套装
	SuitClothes []RawClothes
	RClothes    []RawClothes
	RSuites     []RawSuit

	CloChances map[string]int

	AllClothes       []*GFClothes // 所有的衣服
	AllClothesMap    map[string]*GFClothes
	allClothesCntMap map[int]int

	shopClothesMap map[string]int

	VIPClothes []*GFClothes // VIP独有衣服
	COSClothes []*GFClothes // COS独有衣服

	LotteryClothesLevel1 []*GFClothes // 抽奖服装level
	LotteryClothesLevel2 []*GFClothes // 抽奖服装level
	LotteryClothesLevel3 []*GFClothes // 抽奖服装level
	LotteryClothesLevel4 []*GFClothes // 抽奖服装level
	LotteryClothesLevel5 []*GFClothes // 抽奖服装level

	LotteryPoolActivity map[string]*GFClothes //时空之门奖池
	LotteryPool         map[string][]*GFClothes

	// GuildLevel1     []*GFClothes
	// GuildLevel2     []*GFClothes
	GuildLevel3     []*GFClothes
	GuildLevel4     []*GFClothes
	GuildLevel5     []*GFClothes
	AllGuildClothes []*GFClothes

	suitCloMap map[string][]string //sid:clo list

	// DiamondLotteryClothes []*GFClothes // 钻石抽奖服装

	PKClothes1 []*GFClothes // PK服装level1
	PKClothes2 []*GFClothes // PK服装level2
	PKClothes3 []*GFClothes // PK服装level3
	PKClothes4 []*GFClothes // PK服装level4
	PKClothes5 []*GFClothes // PK服装level5

	// reloadLock sync.RWMutex

	FormulaMap map[string]*Formula
)

func init() {
	// parseMakeup()
	// panic("end")
	// if appcfg.GetServerType() == constants.SERVER_TYPE_SIMPLE_TEST {
	// 	parseDesignerReward()
	// 	panic("end")
	// }

	// parseRDB("/Users/zzd/Downloads/memory.xlsx")
	// panic("end")
	if appcfg.GetServerType() != "" && appcfg.GetServerType() != constants.SERVER_TYPE_GM && appcfg.GetServerType() != constants.SERVER_TYPE_REVOKE_CLOTHES && appcfg.GetServerType() != constants.SERVER_TYPE_SIMPLE_TEST {
		return
	}
	rand.Seed(time.Now().Unix())

	loadClothesByXLSX()
	loadSuitsByXLSX()
	loadCloChancesByXLSX()
	loadFormulasByXLSX()

	Info("all clothes cnt:", len(AllClothes))
	Info("all suits cnt:", len(AllSuits))
	Info("vip clothes cnt:", len(VIPClothes))
	Info("cos clothes cnt:", len(COSClothes))
	// Info("lottery1 clothes cnt:", len(LotteryClothesLevel1))
	// Info("lottery2 clothes cnt:", len(LotteryClothesLevel2))
	// Info("lottery3 clothes cnt:", len(LotteryClothesLevel3))
	// Info("lottery4 clothes cnt:", len(LotteryClothesLevel4))
	// Info("lottery5 clothes cnt:", len(LotteryClothesLevel5))

	for k, v := range LotteryPool {
		Info("lottery pool:", k, "size:", len(v))
	}

	// Info("guild level1 size:", len(GuildLevel1))
	// Info("guild level2 size:", len(GuildLevel2))
	Info("guild level3 size:", len(GuildLevel3))
	Info("guild level4 size:", len(GuildLevel4))
	Info("guild level5 size:", len(GuildLevel5))

	Info("pk1 clothes cnt:", len(PKClothes1))
	Info("pk2 clothes cnt:", len(PKClothes2))
	Info("pk3 clothes cnt:", len(PKClothes3))
	Info("pk4 clothes cnt:", len(PKClothes4))
	Info("pk5 clothes cnt:", len(PKClothes5))

	// if appcfg.GetServerType() == "" {
	// 	checkStartClothes()
	// }
	// Info("guild clothes:", GetGuildWeekRewardClothes())
	// test(100)
	// Info("hahaha:", GetSuitClothesList("c900101000000000054"))
}

// type Content struct {
// 	ClothesID   string
// 	Amount      int
// 	Gold        int      `json:"gold"`
// 	Diamond     int      `json:"diamond"`
// 	ClothesList []string `json:"clothes"`
// }

// type Makeup struct {
// 	ClothesList []string
// 	Diamond     int
// 	Money       int
// }

// func parseMakeup() {
// 	path := "/Users/zzd/Desktop/makeup_20170911.xlsx"

// 	c1Idx := 1
// 	c2Idx := 2
// 	c3Idx := 3
// 	unIdx := 4
// 	contentIdx := 6

// 	res := make(map[string]*Makeup)

// 	xlFile, err := xlsx.OpenFile(path)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, sheet := range xlFile.Sheets {
// 		for ri, row := range sheet.Rows {
// 			if ri > 0 {
// 				var un, c1, c2, c3 string
// 				var con Content
// 				for j, cell := range row.Cells {
// 					if j == c1Idx {
// 						c1, _ = cell.String()
// 					} else if j == c2Idx {
// 						c2, _ = cell.String()
// 					} else if j == c3Idx {
// 						c3, _ = cell.String()
// 					} else if j == unIdx {
// 						un, _ = cell.String()
// 					} else if j == contentIdx {
// 						cons, _ := cell.String()
// 						if err := json.Unmarshal([]byte(cons), &con); err != nil {
// 							panic(err)
// 						}
// 					}
// 				}

// 				if un != "" {
// 					if res[un] == nil {
// 						res[un] = &Makeup{}
// 						res[un].ClothesList = make([]string, 0)
// 					}

// 					if c2 == "DIAMOND" {
// 						if con.Amount < 0 {
// 							panic(c1 + c2 + c3 + un + "diamond<0")
// 						}
// 						res[un].Diamond += con.Amount
// 					} else if c2 == "GOLD" {
// 						if con.Amount < 0 {
// 							panic(c1 + c2 + c3 + un + "gold<0")
// 						}
// 						res[un].Money += con.Amount
// 					} else if c2 == "CLOTHES" {
// 						res[un].ClothesList = append(res[un].ClothesList, con.ClothesID)
// 					} else if c2 == "MAIL" {
// 						res[un].Diamond += con.Diamond
// 						res[un].Money += con.Gold
// 						if con.ClothesList != nil {
// 							res[un].ClothesList = append(res[un].ClothesList, con.ClothesList...)
// 						}
// 					}
// 				} else {
// 					Info(un, c1, c2, c3)
// 				}
// 			}

// 			// memMap[pre] += size
// 		}
// 	}

// 	for k, v := range res {
// 		Info(k, v.Money, v.Diamond, v.ClothesList)
// 	}
// }

func GetAllRawSuits() []RawSuit {
	res := make([]RawSuit, 0)
	for _, s := range AllSuits {
		r := RawSuit{
			ID:   s.SID,
			Name: s.Name,
		}
		res = append(res, r)
	}
	return res
}

func parseDesignerReward() {
	path := "/Users/zzd/Desktop/backup/bili_rank_201805.xlsx"
	etime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2018-06-20 00:00:00", time.Local)
	et := etime.Unix()

	nt := time.Now().Unix()
	title := "设计师排行榜奖励"
	content := "亲爱的设计师您好，恭喜您进入5月设计师排行榜。排名信息如下：%s 您的领奖码为：%s。请在每月20日前到【妖精的衣橱手游】官方微博私信官博（注：此为唯一认证方式）提供您的：游戏昵称、游戏ID、领奖码、收货地址以及手机号。妖精币榜单获奖玩家需额外提供真实姓名、银行卡卡号、银行卡开户行。逾期将为放弃奖励。"

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		panic(err)
	}

	unIdx := 1
	codeIdx := 2
	diamondIdx := 3
	goldIdx := 4
	passIdx := 5
	saleIdx := 6

	for _, sheet := range xlFile.Sheets {
		for ri, row := range sheet.Rows {
			if ri > 0 {
				var un, code string
				var di, gi, pi, si int
				for j, cell := range row.Cells {
					if j == unIdx {
						un, _ = cell.String()
					} else if j == codeIdx {
						code, _ = cell.String()
					} else if j == diamondIdx {
						di, _ = cell.Int()
					} else if j == goldIdx {
						gi, _ = cell.Int()
					} else if j == passIdx {
						pi, _ = cell.Int()
					} else if j == saleIdx {
						si, _ = cell.Int()
					}
				}
				if un != "" {
					cn := getContent(content, di, gi, pi, si, code)
					Info(fmt.Sprintf("insert into gf_mail values(0,0,'妖精管理局','%s','%s','%s',0,0,'[]',%d,0,0,%d);", un, title, cn, nt, et))
				}
			}

			// memMap[pre] += size
		}
	}
}

func getContent(c string, di, gi, pi, si int, code string) string {
	var ps string
	if di > 0 {
		ps += "钻石榜排名：" + fmt.Sprintf("%d", di) + ";"
	}
	if gi > 0 {
		ps += "金币榜排名：" + fmt.Sprintf("%d", gi) + ";"
	}
	if pi > 0 {
		ps += "高产榜排名：" + fmt.Sprintf("%d", pi) + ";"
	}
	if si > 0 {
		ps += "妖精币榜排名：" + fmt.Sprintf("%d", si) + ";"
	}

	return fmt.Sprintf(c, ps, code)
}

// func parseRDB(file string) {
// 	xlFile, err := xlsx.OpenFile(file)
// 	if err != nil {
// 		panic(err)
// 	}

// 	memMap := make(map[string]int64)

// 	for _, sheet := range xlFile.Sheets {
// 		for _, row := range sheet.Rows {
// 			var pre string
// 			var size int64
// 			for j, cell := range row.Cells {
// 				if j == 2 {
// 					key, _ := cell.String()
// 					keys := strings.Split(key, "_")
// 					if len(keys) >= 2 {
// 						pre = keys[0]
// 					} else {
// 						pre = key
// 					}
// 				} else if j == 3 {
// 					size, _ = cell.Int64()
// 				}
// 			}
// 			memMap[pre] += size
// 		}
// 	}

// 	for k, v := range memMap {
// 		Info(k, v)
// 	}
// }

// func test(lcnt int) {
// 	//poolid:2 diamond:10
// 	totalCnt := len(LotteryPool["2:3"]) + len(LotteryPool["2:4"]) + len(LotteryPool["2:5"])

// 	c95t := 0
// 	c100t := 0
// 	for i := 0; i < lcnt; i++ {
// 		cnt95 := 0
// 		gotMap := make(map[string]int)
// 		for len(gotMap) < totalCnt*95/100 {
// 			cnt95++
// 			cs := testGotClothes(2)
// 			for _, c := range cs {
// 				gotMap[c] = 1
// 			}
// 		}

// 		cnt100 := 0
// 		gotMap = make(map[string]int)
// 		for len(gotMap) < totalCnt {
// 			cnt100++
// 			cs := testGotClothes(2)
// 			for _, c := range cs {
// 				gotMap[c] = 1
// 			}
// 		}

// 		// Info("95% cnt:", cnt95, "100% cnt:", cnt100)
// 		c95t += cnt95
// 		c100t += cnt100
// 	}

// 	Info("avg 95%:", c95t/lcnt, "100%:", c100t/lcnt)
// }

// func testGotClothes(poolID int) []string {
// 	clothesArray := make([]string, 0)
// 	contain5 := false
// 	for i := 0; i < 10; i++ {
// 		force5 := false
// 		if !contain5 {
// 			if i == 9 {
// 				force5 = true
// 			} else {
// 				r := rand.Intn(10)
// 				if r == 0 {
// 					force5 = true
// 				}
// 			}

// 		}

// 		reward, level := GetRandomLotteryClothes(true, force5, poolID)

// 		if level == 5 {
// 			contain5 = true
// 		}

// 		clothesArray = append(clothesArray, reward)
// 	}
// 	return clothesArray
// }

func GetGuildWeekRewardClothes() (res []string) {
	var c3, c4, c5 bool
	res = make([]string, 0)
	resm := make(map[string]int)
	for i := 0; i < constants.GUILD_WEEK_CLOTHES_SIZE; i++ {
		var star int

		if i == 7 && !c3 {
			star = 3
		} else if i == 8 && !c4 {
			star = 4
		} else if i == 9 && !c5 {
			star = 5
		}

		c, s := getRandomGuildClothes(star)
		for resm[c] != 0 {
			c, s = getRandomGuildClothes(star)
		}

		res = append(res, c)
		resm[c] = 1
		if s <= 3 {
			c3 = true
		} else if s == 4 {
			c4 = true
		} else if s == 5 {
			c5 = true
		}
	}
	return
}

func getRandomGuildClothes(star int) (string, int) {
	if star > 0 {
		if star <= 3 {
			index := rand.Intn(len(GuildLevel3))
			return GuildLevel3[index].CId, GuildLevel3[index].Star
		} else if star == 4 {
			index := rand.Intn(len(GuildLevel4))
			return GuildLevel4[index].CId, GuildLevel4[index].Star
		} else if star == 5 {
			index := rand.Intn(len(GuildLevel5))
			return GuildLevel5[index].CId, GuildLevel5[index].Star
		}
	} else {
		index := rand.Intn(len(AllGuildClothes))
		return AllGuildClothes[index].CId, AllGuildClothes[index].Star
	}

	return "", 0
}

func GetClothesName(cid string) string {
	if c := GetClothesById(cid); c != nil {
		return c.Name
	} else if _, _, ori, _ := tools.GetInfoFromCustomIDInGame(cid); ori == constants.ORIGIN_DESIGNER {
		return "designer custom"
	} else {
		return "unknown"
	}
}

func ClothesInShop(cid string) bool {
	return shopClothesMap[cid] > 0
}

func GetClothCnt(cid string) int {
	if c := GetClothesById(cid); c != nil {
		if ClothesInShop(cid) {
			return c.Star * constants.CLOTH_STAR_MULTI / 10
		} else {
			return c.Star * constants.CLOTH_STAR_MULTI
		}
	}
	return 0
}

func GetAllClothesCntMap() map[int]int {
	return allClothesCntMap
}

func GetFormula(targetClothes string) *Formula {
	return FormulaMap[targetClothes]
}

// func checkStartClothes() {
// 	var sc []ClothesInfo
// 	if err := json.Unmarshal([]byte(appcfg.GetString("start_clothes", "")), &sc); err != nil {
// 		panic(err)
// 	} else {
// 		for _, sci := range sc {
// 			if c := GetClothesById(sci.ClothesID); c == nil {
// 				panic("start clothes not found:" + sci.ClothesID)
// 			}
// 		}
// 	}
// }

// func ReloadClothes() bool {
// 	return false
// }

func ClothesLegal(id string) (exist bool, star int) {
	if c := GetClothesById(id); c != nil {
		exist = true
		star = c.Star
		return
	} else {
		return
	}
}

func CombinePieceCnt(star int) int {
	cnt := 1
	if star == 2 {
		cnt = 2
	} else if star == 3 {
		cnt = 3
	} else if star == 4 {
		cnt = 4
	} else if star == 5 {
		cnt = 5
	}
	return cnt
}

func ClothesInOneSuit(clom map[string]string, keyIsClo bool) bool {
	cloLen := 0

	for k, v := range clom {
		if !keyIsClo && strings.HasPrefix(v, "c") && len(v) == constants.CLO_LEN {
			cloLen++
		} else if keyIsClo && strings.HasPrefix(k, "c") && len(k) == constants.CLO_LEN {
			cloLen++
		}
	}

	if cloLen <= 0 {
		return false
	}

	//sid : cnt
	m := make(map[string]int)

	for k, v := range clom {
		for _, s := range AllSuits {
			if !keyIsClo && s.IsClothesSuit(v) {
				m[s.SID] = m[s.SID] + 1
			} else if keyIsClo && s.IsClothesSuit(k) {
				m[s.SID] = m[s.SID] + 1
			}
		}
	}

	for _, cnt := range m {
		if cnt >= cloLen {
			return true
		}
	}
	return false
}

func ClothesInOneSuit1(clom map[string]int) bool {
	cloLen := 0

	for k, _ := range clom {
		if strings.HasPrefix(k, "c") && len(k) == constants.CLO_LEN {
			cloLen++
		}
	}

	if cloLen <= 0 {
		return false
	}

	//sid : cnt
	m := make(map[string]int)

	for k, _ := range clom {
		for _, s := range AllSuits {
			if s.IsClothesSuit(k) {
				m[s.SID] = m[s.SID] + 1
			}
		}
	}

	for _, cnt := range m {
		if cnt >= cloLen {
			return true
		}
	}
	return false
}

func GetClothesById(clothesId string) (clothes *GFClothes) {
	// reloadLock.RLock()
	// defer reloadLock.RUnlock()

	clothes, _ = AllClothesMap[clothesId]
	return
}

func GetRandomVIPClothes() string {
	// reloadLock.RLock()
	// defer reloadLock.RUnlock()

	vipClothesCount := len(VIPClothes)
	if vipClothesCount != 0 {
		index := rand.Intn(vipClothesCount)
		return VIPClothes[index].CId
	} else {
		return ""
	}
}

func GetRandomCOSClothes() string {
	// reloadLock.RLock()
	// defer reloadLock.RUnlock()

	cosClothesCount := len(COSClothes)
	if cosClothesCount != 0 {
		index := rand.Intn(cosClothesCount)
		return COSClothes[index].CId
	} else {
		return ""
	}
}

func getRandomLevel(diamondLottery bool, force5 bool) int {
	if force5 {
		return 5
	}

	randValue := rand.Intn(100)
	if diamondLottery {
		// Info("diamond lottery", randValue)
		//diamond lottery
		if randValue <= 79 {
			return 3
		} else if randValue <= 94 {
			return 4
		} else {
			return 5
		}
	} else {
		// Info("gold lottery", randValue)
		//gold lottery
		if randValue <= 94 {
			return rand.Intn(3) + 1
		} else if randValue <= 98 {
			return 4
		} else {
			return 5
		}
	}
}
func GetRandomLotteryClothes(diamondLottery bool, force5 bool, poolID int) (string, int, string) {
	level := getRandomLevel(diamondLottery, force5)
	var tmpClothes []*GFClothes

	if poolID == 0 {
		if level == 1 {
			tmpClothes = LotteryClothesLevel1
		} else if level == 2 {
			tmpClothes = LotteryClothesLevel2
		} else if level == 3 {
			tmpClothes = LotteryClothesLevel3
		} else if level == 4 {
			tmpClothes = LotteryClothesLevel4
		} else {
			tmpClothes = LotteryClothesLevel5
		}
	} else if poolID > 0 {
		if !diamondLottery && level <= 3 {
			tmpClothes = LotteryPool[fmt.Sprintf("%d:%d", poolID, 13)]
		} else {
			tmpClothes = LotteryPool[fmt.Sprintf("%d:%d", poolID, level)]
		}

	}

	lotteryClothesCount := len(tmpClothes)
	if lotteryClothesCount != 0 {
		index := rand.Intn(lotteryClothesCount)
		return tmpClothes[index].CId, tmpClothes[index].Star, tmpClothes[index].Name
	} else {
		return "", 0, ""
	}
}

func GetClothesNameAndStar(cid string) []string {
	if c := GetClothesById(cid); c != nil {
		return []string{c.Name, fmt.Sprintf("%d❤", c.Star)}
	}
	return []string{"", "0❤"}
}

//GetClothesStar 获取衣服星级
func GetClothesStar(cid string) string {
	c := GetClothesById(cid)
	return fmt.Sprintf("%d", c.Star)
}

func GetRandomPKClothesByLevel(level int) string {
	// reloadLock.RLock()
	// defer reloadLock.RUnlock()

	if level == 1 {
		clothesCount := len(PKClothes1)
		if clothesCount != 0 {
			index := rand.Intn(clothesCount)
			return PKClothes1[index].CId
		} else {
			return ""
		}
	} else if level == 2 {
		clothesCount := len(PKClothes2)
		if clothesCount != 0 {
			index := rand.Intn(clothesCount)
			return PKClothes2[index].CId
		} else {
			return ""
		}
	} else if level == 3 {
		clothesCount := len(PKClothes3)
		if clothesCount != 0 {
			index := rand.Intn(clothesCount)
			return PKClothes3[index].CId
		} else {
			return ""
		}
	} else if level == 4 {
		clothesCount := len(PKClothes4)
		if clothesCount != 0 {
			index := rand.Intn(clothesCount)
			return PKClothes4[index].CId
		} else {
			return ""
		}
	} else if level == 5 {
		clothesCount := len(PKClothes5)
		if clothesCount != 0 {
			index := rand.Intn(clothesCount)
			return PKClothes5[index].CId
		} else {
			return ""
		}
	} else {
		return ""
	}
}

func IsCloSuit(cloID string) bool {
	// reloadLock.RLock()
	// defer reloadLock.RUnlock()

	for _, v := range AllSuits {
		if v.IsClothesSuit(cloID) {
			return true
		}
	}

	return false
}

func IsSuit(suitID string) bool {
	// reloadLock.RLock()
	// defer reloadLock.RUnlock()

	for _, v := range AllSuits {
		if v.SID == suitID {
			return true
		}
	}

	return false
}

func GetRandomSuitClothes(suitID string) (clothes *GFClothes) {
	// reloadLock.RLock()
	// defer reloadLock.RUnlock()

	var suit *GFSuit

	for _, v := range AllSuits {
		if v.SID == suitID {
			suit = v
			break
		}
	}

	index := rand.Intn(len(suit.Clos))

	return GetClothesById(suit.Clos[index])
}

func GetSuitClothesList(suitID string) (ret []string) {
	// reloadLock.Lock()
	// reloadLock.Unlock()

	for _, v := range AllSuits {
		if v.SID == suitID {
			return v.Clos
		}
	}

	return
}

func GetSuitClothesMap(suitID string) (ret map[string]int) {
	// reloadLock.Lock()
	// reloadLock.Unlock()

	for _, v := range AllSuits {
		if v.SID == suitID {
			ret = make(map[string]int)
			if v.Hair != "" {
				ret[v.Hair] = 1
			}
			if v.Coat != "" {
				ret[v.Coat] = 1
			}
			if v.Shirt != "" {
				ret[v.Shirt] = 1
			}
			if v.Pants != "" {
				ret[v.Pants] = 1
			}
			if v.Socks != "" {
				ret[v.Socks] = 1
			}
			if v.Shoes != "" {
				ret[v.Shoes] = 1
			}
			if v.Dress != "" {
				ret[v.Dress] = 1
			}
			if v.Hat != "" {
				ret[v.Hat] = 1
			}
			if v.Glasses != "" {
				ret[v.Glasses] = 1
			}
			if v.Earrings != "" {
				ret[v.Earrings] = 1
			}
			if v.Necklace != "" {
				ret[v.Necklace] = 1
			}
			if v.Tie != "" {
				ret[v.Tie] = 1
			}
			if v.Wrist != "" {
				ret[v.Wrist] = 1
			}
			if v.Bag != "" {
				ret[v.Bag] = 1
			}
			if v.Other != "" {
				ret[v.Other] = 1
			}
			if v.Face != "" {
				ret[v.Face] = 1
			}
		}
	}

	return
}

func loadFormulasByXLSX() {
	FormulaMap = make(map[string]*Formula)

	// f := "scripts/data/data_tailor.xlsx"
	// if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
	// 	f = "scripts/data/data_tailor_kor.xlsx"
	// }
	// xlFile, err := xlsx.OpenFile(f)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, sheet := range xlFile.Sheets {
	// 	if sheet.Name == "Sheet1" {
	// 		for i, row := range sheet.Rows {
	// 			if i > 0 {
	// 				var target, origin string
	// 				var gold int64
	// 				var cloth, diamond, item1, item2, item3, item4, item5 int
	// 				for j, cell := range row.Cells {
	// 					if j == 0 {
	// 						target, _ = cell.String()
	// 					} else if j == 1 {
	// 						origin, _ = cell.String()
	// 					} else if j == 2 {
	// 						cloth, _ = cell.Int()
	// 					} else if j == 3 {
	// 						gold, _ = cell.Int64()
	// 					} else if j == 4 {
	// 						diamond, _ = cell.Int()
	// 					} else if j == 5 {
	// 						item1, _ = cell.Int()
	// 					} else if j == 6 {
	// 						item2, _ = cell.Int()
	// 					} else if j == 7 {
	// 						item3, _ = cell.Int()
	// 					} else if j == 8 {
	// 						item4, _ = cell.Int()
	// 					} else if j == 9 {
	// 						item5, _ = cell.Int()
	// 					}
	// 				}
	// 				if target != "" {
	// 					FormulaMap[target] = &Formula{
	// 						OriginClothes: origin,
	// 						ClothCnt:      cloth,
	// 						Gold:          gold,
	// 						Diamond:       diamond,
	// 						Item1:         item1,
	// 						Item2:         item2,
	// 						Item3:         item3,
	// 						Item4:         item4,
	// 						Item5:         item5,
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}

func loadClothesByXLSX() {
	AllClothes = make([]*GFClothes, 0)
	AllClothesMap = make(map[string]*GFClothes, 0)
	RClothes = make([]RawClothes, 0)
	allClothesCntMap = make(map[int]int)

	shopClothesMap = make(map[string]int)

	VIPClothes = make([]*GFClothes, 0)
	COSClothes = make([]*GFClothes, 0)

	LotteryClothesLevel1 = make([]*GFClothes, 0)
	LotteryClothesLevel2 = make([]*GFClothes, 0)
	LotteryClothesLevel3 = make([]*GFClothes, 0)
	LotteryClothesLevel4 = make([]*GFClothes, 0)
	LotteryClothesLevel5 = make([]*GFClothes, 0)
	LotteryPool = make(map[string][]*GFClothes)
	LotteryPoolActivity = make(map[string]*GFClothes)
	// GuildLevel1 = make([]*GFClothes, 0)
	// GuildLevel2 = make([]*GFClothes, 0)
	GuildLevel3 = make([]*GFClothes, 0)
	GuildLevel4 = make([]*GFClothes, 0)
	GuildLevel5 = make([]*GFClothes, 0)
	AllGuildClothes = make([]*GFClothes, 0)

	PKClothes1 = make([]*GFClothes, 0)
	PKClothes2 = make([]*GFClothes, 0)
	PKClothes3 = make([]*GFClothes, 0)
	PKClothes4 = make([]*GFClothes, 0)
	PKClothes5 = make([]*GFClothes, 0)

	// tm := make(map[string]int)

	// f := "scripts/data/data_clothes.xlsx"
	// if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
	// 	f = "scripts/data/data_clothes_kor.xlsx"
	// }

	// xlFile, err := xlsx.OpenFile(f)
	// if err != nil {
	// 	panic(err)
	// }
	// var idIdx, nameIdx, typeIdx, smallTypeIdx, priceIdx, priceTypeIdx, sellPriceIdx, starIdx, dropCosIdx, dropPKIdx, dropVIPIdx, dropLotteryIdx, charIdx int
	// var warmIdx, formalIdx, tightIdx, brightIdx, darkIdx, cuteIdx, manIdx, toughIdx, nobleIdx, strangeIdx, sexyIdx, sportIdx, storeIdx, guildIdx, doorIdx int
	// for _, sheet := range xlFile.Sheets {
	// 	if sheet.Name == "Sheet1" {
	// 		for i, row := range sheet.Rows {
	// 			if i == 0 {
	// 				for j, cell := range row.Cells {
	// 					v, _ := cell.String()
	// 					if v == "id" {
	// 						idIdx = j
	// 					} else if v == "name" {
	// 						nameIdx = j
	// 					} else if v == "type" {
	// 						typeIdx = j
	// 					} else if v == "small_type" {
	// 						smallTypeIdx = j
	// 					} else if v == "price" {
	// 						priceIdx = j
	// 					} else if v == "price_type" {
	// 						priceTypeIdx = j
	// 					} else if v == "sell_price" {
	// 						sellPriceIdx = j
	// 					} else if v == "star" {
	// 						starIdx = j
	// 					} else if v == "cos_drop" {
	// 						dropCosIdx = j
	// 					} else if v == "pk_drop" {
	// 						dropPKIdx = j
	// 					} else if v == "vip_drop" {
	// 						dropVIPIdx = j
	// 					} else if v == "lottery_drop" {
	// 						dropLotteryIdx = j
	// 					} else if v == "char" {
	// 						charIdx = j
	// 					} else if v == "warm" {
	// 						warmIdx = j
	// 					} else if v == "formal" {
	// 						formalIdx = j
	// 					} else if v == "tight" {
	// 						tightIdx = j
	// 					} else if v == "dark" {
	// 						darkIdx = j
	// 					} else if v == "cute" {
	// 						cuteIdx = j
	// 					} else if v == "man" {
	// 						manIdx = j
	// 					} else if v == "tough" {
	// 						toughIdx = j
	// 					} else if v == "noble" {
	// 						nobleIdx = j
	// 					} else if v == "strange" {
	// 						strangeIdx = j
	// 					} else if v == "sexy" {
	// 						sexyIdx = j
	// 					} else if v == "sport" {
	// 						sportIdx = j
	// 					} else if v == "store" {
	// 						storeIdx = j
	// 					} else if v == "team" {
	// 						guildIdx = j
	// 					} else if v == "door" {
	// 						doorIdx = j
	// 					}
	// 				}
	// 				continue
	// 			}
	// 			c := GFClothes{}
	// 			rc := RawClothes{}
	// 			store := 0
	// 			guild := 0
	// 			door := 0
	// 			for j, cell := range row.Cells {
	// 				v, _ := cell.String()

	// 				if j == idIdx {
	// 					//levelid
	// 					// indexlist = append(indexlist, v)
	// 					c.CId = v
	// 					rc.ID = v
	// 				} else if j == nameIdx {
	// 					c.Name = v
	// 					rc.Name = v
	// 				} else if j == typeIdx {
	// 					c.Type = getInt(v, "clothes type")
	// 					rc.Type = c.Type
	// 				} else if j == smallTypeIdx {
	// 					c.SmallType = getInt(v, "clothes small type")
	// 				} else if j == priceIdx {
	// 					c.Price = getInt(v, "clothes price")
	// 				} else if j == priceTypeIdx {
	// 					c.PriceType = getInt(v, "clothes price type")
	// 				} else if j == sellPriceIdx {
	// 					c.SellPrice = getInt(v, "clothes sell price")
	// 				} else if j == starIdx {
	// 					c.Star = getInt(v, "clothes star")
	// 					rc.Star = c.Star
	// 				} else if j == dropCosIdx {
	// 					c.DropCos = getInt(v, "clothes cos drop")
	// 				} else if j == dropPKIdx {
	// 					c.DropPK = getInt(v, "clothes pk drop")
	// 				} else if j == dropVIPIdx {
	// 					c.DropVIP = getInt(v, "clothes vip drop")
	// 				} else if j == dropLotteryIdx {
	// 					c.DropLottery = getInt(v, "clothes lottery drop")
	// 				} else if j == charIdx {
	// 					c.Char = getInt(v, "clothes char")
	// 					rc.ModelNo = c.Char
	// 				} else if j == warmIdx {
	// 					c.Warm = getFloat64(v, "clothes warm")
	// 				} else if j == formalIdx {
	// 					c.Formal = getFloat64(v, "clothes formal")
	// 				} else if j == tightIdx {
	// 					c.Tight = getFloat64(v, "clothes tight")
	// 				} else if j == brightIdx {
	// 					c.Bright = getFloat64(v, "clothes bright")
	// 				} else if j == darkIdx {
	// 					c.Dark = getFloat64(v, "clothes dark")
	// 				} else if j == cuteIdx {
	// 					c.Cute = getFloat64(v, "clothes cute")
	// 				} else if j == manIdx {
	// 					c.Man = getFloat64(v, "clothes man")
	// 				} else if j == toughIdx {
	// 					c.Tough = getFloat64(v, "clothes tough")
	// 				} else if j == nobleIdx {
	// 					c.Noble = getFloat64(v, "clothes noble")
	// 				} else if j == strangeIdx {
	// 					c.Strange = getFloat64(v, "clothes strange")
	// 				} else if j == sexyIdx {
	// 					c.Sexy = getFloat64(v, "clothes sexy")
	// 				} else if j == sportIdx {
	// 					c.Sport = getFloat64(v, "clothes sport")
	// 				} else if j == storeIdx {
	// 					store = getInt(v, "store")
	// 				} else if j == guildIdx {
	// 					guild = getInt(v, "guild")
	// 				} else if j == doorIdx {
	// 					door = getInt(v, "door")
	// 				}
	// 			}
	// 			if tm[c.CId] > 0 {
	// 				panic("clothes id duplicate:" + c.CId)
	// 			} else {
	// 				tm[c.CId] = 1
	// 			}

	// 			t, m, _, _ := tools.GetInfoFromCustomIDInGame(c.CId)
	// 			if t != c.Type {
	// 				if c.Type == 700 && t >= 700 && t < 800 {
	// 					//ok
	// 				} else {
	// 					panic("clothes type wrong:" + c.CId)
	// 				}
	// 			}
	// 			if m != c.Char {
	// 				panic("clothes model wrong:" + c.CId)
	// 			}

	// 			if i == 1 {
	// 				Info("clothes sample:", c)
	// 			}
	// 			if store == 1 {
	// 				shopClothesMap[c.CId] = 1
	// 			}
	// 			if guild == 1 {
	// 				AllGuildClothes = append(AllGuildClothes, &c)
	// 				if c.Star <= 3 {
	// 					GuildLevel3 = append(GuildLevel3, &c)
	// 				} else if c.Star == 4 {
	// 					GuildLevel4 = append(GuildLevel4, &c)
	// 				} else if c.Star == 5 {
	// 					GuildLevel5 = append(GuildLevel5, &c)
	// 				}
	// 			}
	// 			//时空之门活动奖池
	// 			if door == 2 {
	// 				LotteryPoolActivity[c.CId] = &c
	// 			}
	// 			RClothes = append(RClothes, rc)
	// 			AllClothes = append(AllClothes, &c)
	// 			allClothesCntMap[c.Char] += 1
	// 		}
	// 	}
	// }

	// for _, cloth := range AllClothes {
	// 	AllClothesMap[cloth.CId] = cloth

	// 	if cloth.DropVIP == 1 {
	// 		VIPClothes = append(VIPClothes, cloth)
	// 	}
	// 	if cloth.DropPK == 1 {
	// 		if cloth.Star == 1 {
	// 			PKClothes1 = append(PKClothes1, cloth)
	// 		}
	// 		if cloth.Star == 2 {
	// 			PKClothes2 = append(PKClothes2, cloth)
	// 		}
	// 		if cloth.Star == 3 {
	// 			PKClothes3 = append(PKClothes3, cloth)
	// 		}
	// 		if cloth.Star == 4 {
	// 			PKClothes4 = append(PKClothes4, cloth)
	// 		}
	// 		if cloth.Star == 5 {
	// 			PKClothes5 = append(PKClothes5, cloth)
	// 		}
	// 	}
	// 	if cloth.DropCos == 1 {
	// 		COSClothes = append(COSClothes, cloth)
	// 	}
	// 	if cloth.DropLottery > 0 {
	// 		if cloth.Star == 1 {
	// 			LotteryClothesLevel1 = append(LotteryClothesLevel1, cloth)
	// 		} else if cloth.Star == 2 {
	// 			LotteryClothesLevel2 = append(LotteryClothesLevel2, cloth)
	// 		} else if cloth.Star == 3 {
	// 			LotteryClothesLevel3 = append(LotteryClothesLevel3, cloth)
	// 		} else if cloth.Star == 4 {
	// 			LotteryClothesLevel4 = append(LotteryClothesLevel4, cloth)
	// 		} else if cloth.Star == 5 {
	// 			LotteryClothesLevel5 = append(LotteryClothesLevel5, cloth)
	// 		}

	// 		//lottery pool
	// 		key := fmt.Sprintf("%d:%d", cloth.DropLottery, cloth.Star)
	// 		if LotteryPool[key] == nil {
	// 			LotteryPool[key] = make([]*GFClothes, 0)
	// 		}

	// 		LotteryPool[key] = append(LotteryPool[key], cloth)

	// 		if cloth.Star <= 3 {
	// 			key13 := fmt.Sprintf("%d:%d", cloth.DropLottery, 13)
	// 			if LotteryPool[key13] == nil {
	// 				LotteryPool[key13] = make([]*GFClothes, 0)
	// 			}
	// 			LotteryPool[key13] = append(LotteryPool[key13], cloth)
	// 		}
	// 	}
	// }

	// fmt.Println("Load clothes data over...")
}

func loadCloChancesByXLSX() {
	CloChances = make(map[string]int)

	// f := "scripts/data/data_chance.xlsx"
	// if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
	// 	f = "scripts/data/data_chance_kor.xlsx"
	// }
	// xlFile, err := xlsx.OpenFile(f)
	// if err != nil {
	// 	Err(err)
	// 	return
	// }

	// for _, sheet := range xlFile.Sheets {
	// 	if sheet.Name == "Sheet1" {
	// 		for i, row := range sheet.Rows {
	// 			if i > 0 {
	// 				var cloID string
	// 				var chance int
	// 				for j, cell := range row.Cells {
	// 					v, _ := cell.String()
	// 					if j == 0 {
	// 						//clo id
	// 						cloID = v
	// 					} else if j == 1 {
	// 						//chance
	// 						chance, _ = strconv.Atoi(v)
	// 					}
	// 				}
	// 				CloChances[cloID] = chance
	// 			}
	// 		}
	// 	}
	// }
	// Info("clo chances:", CloChances)
}

func CloPieceChance(cloID string) int {
	return CloChances[cloID]
}

func loadSuitsByXLSX() {
	AllSuits = make([]*GFSuit, 0)
	SuitClothes = make([]RawClothes, 0)
	RSuites = make([]RawSuit, 0)
	// suitCloMap = make(map[string][]string)

	// f := "scripts/data/data_suites.xlsx"
	// if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
	// 	f = "scripts/data/data_suites_kor.xlsx"
	// }

	// xlFile, err := xlsx.OpenFile(f)
	// if err != nil {
	// 	Err(err)
	// 	return
	// }

	// var sidIdx, charIdx, nameIdx, hairIdx, coatIdx, shirtIdx, pantIdx, sockIdx, shoeIdx int
	// var dressIdx, hatIdx, glassIdx, errIdx, neckIdx, tieIdx, wristIdx, bagIdx, otherIdx, faceIdx int
	// m := make(map[string]int)
	// for _, sheet := range xlFile.Sheets {
	// 	if sheet.Name == "Sheet1" {
	// 		for i, row := range sheet.Rows {
	// 			if i == 0 {
	// 				for j, cell := range row.Cells {
	// 					v, _ := cell.String()
	// 					if v == "id" {
	// 						sidIdx = j
	// 					} else if v == "name" {
	// 						nameIdx = j
	// 					} else if v == "char" {
	// 						charIdx = j
	// 					} else if v == "hair" {
	// 						hairIdx = j
	// 					} else if v == "coat" {
	// 						coatIdx = j
	// 					} else if v == "shirt" {
	// 						shirtIdx = j
	// 					} else if v == "pants" {
	// 						pantIdx = j
	// 					} else if v == "socks" {
	// 						sockIdx = j
	// 					} else if v == "shoes" {
	// 						shoeIdx = j
	// 					} else if v == "dress" {
	// 						dressIdx = j
	// 					} else if v == "hat" {
	// 						hatIdx = j
	// 					} else if v == "glasses" {
	// 						glassIdx = j
	// 					} else if v == "necklace" {
	// 						neckIdx = j
	// 					} else if v == "tie" {
	// 						tieIdx = j
	// 					} else if v == "wrist" {
	// 						wristIdx = j
	// 					} else if v == "bag" {
	// 						bagIdx = j
	// 					} else if v == "other" {
	// 						otherIdx = j
	// 					} else if v == "earrings" {
	// 						errIdx = j
	// 					} else if v == "face" {
	// 						faceIdx = j
	// 					}
	// 				}
	// 				Info("suits idx:", sidIdx, charIdx, nameIdx, hairIdx, coatIdx, shirtIdx, pantIdx, sockIdx, shoeIdx, dressIdx, hatIdx, glassIdx, errIdx, neckIdx, tieIdx, wristIdx, bagIdx, otherIdx)
	// 				continue
	// 			}
	// 			s := GFSuit{
	// 				Clos: make([]string, 0),
	// 			}
	// 			for j, cell := range row.Cells {
	// 				v, _ := cell.String()

	// 				if j == sidIdx {
	// 					s.SID = v
	// 				} else if j == charIdx {
	// 					s.Char = getInt(v, "suit char")
	// 				} else if j == nameIdx {
	// 					s.Name = v
	// 				} else if j == hairIdx {
	// 					s.Hair = v
	// 				} else if j == coatIdx {
	// 					s.Coat = v
	// 				} else if j == shirtIdx {
	// 					s.Shirt = v
	// 				} else if j == pantIdx {
	// 					s.Pants = v
	// 				} else if j == sockIdx {
	// 					s.Socks = v
	// 				} else if j == shoeIdx {
	// 					s.Shoes = v
	// 				} else if j == dressIdx {
	// 					s.Dress = v
	// 				} else if j == hatIdx {
	// 					s.Hat = v
	// 				} else if j == glassIdx {
	// 					s.Glasses = v
	// 				} else if j == errIdx {
	// 					s.Earrings = v
	// 				} else if j == neckIdx {
	// 					s.Necklace = v
	// 				} else if j == tieIdx {
	// 					s.Tie = v
	// 				} else if j == wristIdx {
	// 					s.Wrist = v
	// 				} else if j == bagIdx {
	// 					s.Bag = v
	// 				} else if j == otherIdx {
	// 					s.Other = v
	// 				} else if j == faceIdx {
	// 					s.Face = v
	// 				}

	// 				if j >= 5 && j <= 20 && v != "" {
	// 					s.Clos = append(s.Clos, v)
	// 					m[v] = 1
	// 				}
	// 			}
	// 			if i == 1 {
	// 				Info("suit sample:", s)
	// 			}
	// 			AllSuits = append(AllSuits, &s)
	// 			RSuites = append(RSuites, RawSuit{ID: s.SID, Name: s.Name})
	// 		}
	// 	}
	// }

	// for k, _ := range m {
	// 	c := GetClothesById(k)
	// 	if c != nil {
	// 		rc := RawClothes{
	// 			ID:   c.CId,
	// 			Name: c.Name,
	// 		}
	// 		SuitClothes = append(SuitClothes, rc)
	// 	} else {
	// 		panic("suit clothes:" + k + " not in data_clothes")
	// 	}
	// }
}

func getInt(v string, info string) int {
	if r, err := strconv.Atoi(v); err != nil {
		Err(err, v, info)
		panic(err)
	} else {
		return r
	}
}

func getFloat64(v string, info string) float64 {
	if r, err := strconv.ParseFloat(v, 64); err != nil {
		Err(err, v, info)
		panic(err)
	} else {
		return r
	}
}
