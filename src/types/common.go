package types

import (
	. "pixie_contract/api_specification"
)

type ContentJson struct {
	//用于json解析的条目名称对应
	Cid           int        `json:"cid"`
	Username      string     `json:"username"`
	Nickname      string     `json:"nickname"`
	ModelNo       int        `json:"modelNo"`
	Type          int        `json:"type"`
	Cname         string     `json:"cname"`
	Desc          string     `json:"desc"`
	Icon          string     `json:"icon"`
	Clothes       string     `json:"clothes"`
	CloSet        CustomFile `bson:"clothesSet" json:"clothesSet"` // clothesSet的文件位置
	Status        int        `json:"status"`
	UploadTime    int64      `json:"uploadTime"`
	CheckTime     int64      `json:"checkTime"`
	AdminUsername string     `json:"adminUsername"`
	AdminNickname string     `json:"adminNickname"`
	MoneyType     int        `json:"mType"`
	MoneyAmount   int        `json:"mAmount"`
	Hearts        int        `json:"hearts"`
	Sale          bool       `json:"sale"`
	Platform      string     `json:"platform"`
}

type ResponseJson struct {
	//用于json解析的条目名称对应
	Status  int         `json:"status"`
	Content ContentJson `json:"content"`
	Count   int         `json:"count"`
}

type Reward struct {
	GotGold    int
	GotDiamond int
	GotTili    int

	CurrentGold    int
	CurrentDiamond int
	CurrentTili    int

	Clothes []ClothesInfo
}
