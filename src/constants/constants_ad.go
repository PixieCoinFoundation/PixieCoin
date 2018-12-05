package constants

//SellMod  时间轮模型
type SellMod struct {
	ID     int
	Start  string
	Close  string
	InUse  []string
	OnSale []string
}

var PageInfoMap = map[string]string{
	"showMonthCardIAPLayer()":                          "弹出小额礼包页面",
	"GameController:getInstance():showDesignerLayer()": "跳转到设计师社区主页面",
	"GameController:getInstance():showCosplay()":       "跳转到换装舞会",
	"GameController:getInstance():showTailorShop()":    "跳转到裁缝铺",
	"GameController:getInstance():showPKLayer()":       "跳转到快速pk",
	"GameController:getInstance():showTheaterLayer()":  "跳转到剧场",
}

var KorPageInfoMap = map[string]string{
	"showMonthCardIAPLayer()":                          "Pop up the small package page",
	"GameController:getInstance():showDesignerLayer()": "Jump to the designer community home page",
	"GameController:getInstance():showCosplay()":       "Jump to the dressing ball",
	"GameController:getInstance():showTailorShop()":    "Jump to the tailor shop",
	"GameController:getInstance():showPKLayer()":       "Jump to quick pk",
	"GameController:getInstance():showTheaterLayer()":  "Jump to theatre",
}

var Partlist = []int{100, 200, 300, 400, 500, 600, 701, 702, 703, 704, 705, 706, 707, 708, 709, 800}
var Modellist = []int{101, 103, 201, 202}
var Marginlist = []int{1, 2}

const (
	TRANSFERCHARGE = 0.95
)

var TIME_CIRCLE = map[int]SellMod{
	0: SellMod{
		ID:     0,
		Start:  "4-8",
		Close:  "20-0",
		InUse:  []string{"17-21", "18-22", "19-23", "20-0"},
		OnSale: []string{"21-1", "22-2", "23-3", "0-4"},
	},
	1: SellMod{
		ID:     1,
		Start:  "5-9",
		Close:  "21-1",
		InUse:  []string{"18-22", "19-23", "20-0", "21-1"},
		OnSale: []string{"22-2", "23-3", "0-4", "1-5"},
	},
	2: SellMod{
		ID:     2,
		Start:  "6-10",
		Close:  "22-2",
		InUse:  []string{"19-23", "20-0", "21-1", "22-2"},
		OnSale: []string{"23-3", "0-4", "1-5", "2-6"},
	},
	3: SellMod{
		ID:     3,
		Start:  "7-11",
		Close:  "23-3",
		InUse:  []string{"20-0", "21-1", "22-2", "23-3"},
		OnSale: []string{"0-4", "1-5", "2-6", "3-7"},
	},
	4: SellMod{
		ID:     4,
		Start:  "8-12",
		Close:  "0-4",
		InUse:  []string{"21-1", "22-2", "23-3", "0-4"},
		OnSale: []string{"1-5", "2-6", "3-7", "4-8"},
	},
	5: SellMod{
		ID:     5,
		Start:  "9-13",
		Close:  "1-5",
		InUse:  []string{"22-2", "23-3", "0-4", "1-5"},
		OnSale: []string{"2-6", "3-7", "4-8", "5-9"},
	},
	6: SellMod{
		ID:     6,
		Start:  "10-14",
		Close:  "2-6",
		InUse:  []string{"23-3", "0-4", "1-5", "2-6"},
		OnSale: []string{"3-7", "4-8", "5-9", "6-10"},
	},
	7: SellMod{
		ID:     7,
		Start:  "11-15",
		Close:  "3-7",
		InUse:  []string{"0-4", "1-5", "2-6", "3-7"},
		OnSale: []string{"4-8", "5-9", "6-10", "7-11"},
	},
	8: SellMod{
		ID:     8,
		Start:  "12-16",
		Close:  "4-8",
		InUse:  []string{"1-5", "2-6", "3-7", "4-8"},
		OnSale: []string{"5-9", "6-10", "7-11", "8-12"},
	},
	9: SellMod{
		ID:     9,
		Start:  "13-17",
		Close:  "5-9",
		InUse:  []string{"2-6", "3-7", "4-8", "5-9"},
		OnSale: []string{"6-10", "7-11", "8-12", "9-13"},
	},
	10: SellMod{
		ID:     10,
		Start:  "14-18",
		Close:  "6-10",
		InUse:  []string{"3-7", "4-8", "5-9", "6-10"},
		OnSale: []string{"7-11", "8-12", "9-13", "10-14"},
	},
	11: SellMod{
		ID:     11,
		Start:  "15-19",
		Close:  "7-11",
		InUse:  []string{"4-8", "5-9", "6-10", "7-11"},
		OnSale: []string{"8-12", "9-13", "10-14", "11-15"},
	},
	12: SellMod{
		ID:     12,
		Start:  "16-20",
		Close:  "8-12",
		InUse:  []string{"5-9", "6-10", "7-11", "8-12"},
		OnSale: []string{"9-13", "10-14", "11-15", "12-16"},
	},
	13: SellMod{
		ID:     13,
		Start:  "17-21",
		Close:  "9-13",
		InUse:  []string{"6-10", "7-11", "8-12", "9-13"},
		OnSale: []string{"10-14", "11-15", "12-16", "13-17"},
	},
	14: SellMod{
		ID:     14,
		Start:  "18-22",
		Close:  "10-14",
		InUse:  []string{"7-11", "8-12", "9-13", "10-14"},
		OnSale: []string{"11-15", "12-16", "13-17", "14-18"},
	},
	15: SellMod{
		ID:     15,
		Start:  "19-23",
		Close:  "11-15",
		InUse:  []string{"8-12", "9-13", "10-14", "11-15"},
		OnSale: []string{"12-16", "13-17", "14-18", "15-19"},
	},
	16: SellMod{
		ID:     16,
		Start:  "20-24",
		Close:  "12-16",
		InUse:  []string{"9-13", "10-14", "11-15", "12-16"},
		OnSale: []string{"13-17", "14-18", "15-19", "16-20"},
	},
	17: SellMod{
		ID:     17,
		Start:  "21-1",
		Close:  "13-17",
		InUse:  []string{"10-14", "11-15", "12-16", "13-17"},
		OnSale: []string{"14-18", "15-19", "16-20", "17-21"},
	},
	18: SellMod{
		ID:     18,
		Start:  "22-2",
		Close:  "14-18",
		InUse:  []string{"11-15", "12-16", "13-17", "14-18"},
		OnSale: []string{"15-19", "16-20", "17-21", "18-22"},
	},
	19: SellMod{
		ID:     19,
		Start:  "23-3",
		Close:  "15-19",
		InUse:  []string{"12-16", "13-17", "14-18", "15-19"},
		OnSale: []string{"16-20", "17-21", "18-22", "19-23"},
	},
	20: SellMod{
		ID:     20,
		Start:  "0-4",
		Close:  "16-20",
		InUse:  []string{"13-17", "14-18", "15-19", "16-20"},
		OnSale: []string{"17-21", "18-22", "19-23", "20-0"},
	},
	21: SellMod{
		ID:     21,
		Start:  "1-5",
		Close:  "17-21",
		InUse:  []string{"14-18", "15-19", "16-20", "17-21"},
		OnSale: []string{"18-22", "19-23", "20-21", "21-1"},
	},
	22: SellMod{
		ID:     22,
		Start:  "2-6",
		Close:  "18-22",
		InUse:  []string{"15-19", "16-20", "17-21", "18-22"},
		OnSale: []string{"19-23", "20-0", "21-1", "22-2"},
	},
	23: SellMod{
		ID:     23,
		Start:  "3-7",
		Close:  "19-23",
		InUse:  []string{"16-20", "17-21", "18-22", "19-23"},
		OnSale: []string{"20-0", "21-1", "22-2", "23-3"},
	},
}
