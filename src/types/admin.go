package types

import (
// "gopkg.in/mgo.v2/bson"
)

// 管理员
type Admin struct {
	Id       int64  `bson:"_id",json:"id"`
	Username string `bson:"username",json:"username"`
	Password string `bson:"password",json:"password"`
	Nickname string `bson:"nickname",json:"nickname"`

	LastLoginTime string `bson:"loginTime",json:"loginTime"`
	Role          int    `bson:"role",json:"role"`
}
