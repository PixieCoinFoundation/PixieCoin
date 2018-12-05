package types

import ()

type Player struct {
	Id       int64
	Password string `bson:"password" json:"password"`
	Username string `bson:"username" json:"username"`
	Nickname string `bson:"nickname" json:"nickname"`
	Udid     string `bson:"udid" json:"udid"`

	LastLoginTime string `bson:"loginTime" json:"loginTime"`
}

type SimpleCart struct {
	Username  string
	ClothesID string
}
