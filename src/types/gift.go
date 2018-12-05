package types

import (
	. "pixie_contract/api_specification"
)

type GiftPack struct {
	ID      int64
	Name    string
	Content string
}

type GiftPackContent struct {
	Gold         int
	Diamond      int
	Clothes      []ClothesInfo
	Item101      int
	Item102      int
	Item103      int
	Item104      int
	Item105      int
	Item106      int
	StartTime    int64
	EndTime      int64
	Type         int
	Code         string
	Tili         int
	ThirdChannel string
}

//GiftExtra 礼物日志
type GiftExtra struct {
	Code          string
	PackID        int64
	PackName      string
	PackType      int
	ServerAddress string
	DeviceId      string
}
