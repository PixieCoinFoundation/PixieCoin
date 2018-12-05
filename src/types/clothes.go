package types

import (
	"encoding/json"
)

type RawSuit struct {
	ID   string
	Name string
}

type Clothes struct {
	Username string `bson:"username" json:"-"`
	Clothes  string `json:"clothes"`
}

func (c Clothes) String() string {
	jsonbyte, _ := json.Marshal(c)
	return string(jsonbyte)
}

type ClothesDetailInfo struct {
	ClothesID string `json:"c"`
	Count     int    `json:"i"`
	Star      int    `json:"s"`
	Name      string `json:"n"`
}

type ClothesInfo1 struct {
	ClothesID string `json:"c"`
	Count     int    `json:"i"`
}

type GMClothesInfo struct {
	ClothesID   string `bson:"c" json:"c"`
	ClothesName string `bson:"cn" json:"cn"`
	Count       int    `bson:"i" json:"i"`
}

type GMClothesPieceInfo struct {
	ClothesPieceID string `bson:"c" json:"c"`
	ClothesName    string `bson:"cn" json:"cn"`
	Count          int    `bson:"i" json:"i"`
}

type RawClothes struct {
	ID      string
	Name    string
	Type    int
	ModelNo int
	Star    int
}

type Formula struct {
	OriginClothes string
	ClothCnt      int
	Gold          int64
	Diamond       int
	Item1         int //纽扣
	Item2         int //魔术针
	Item3         int //五彩线
	Item4         int
	Item5         int
}
