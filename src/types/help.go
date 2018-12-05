package types

import (
// "gopkg.in/mgo.v2/bson"
)

type HelpComment struct {
	CommentID int64  `bson:"commentID" json:"commentID"`
	HelpID    int64  `bson:"helpID" json:"helpID"`
	Username  string `bson:"username" json:"username"`
	Nickname  string `bson:"nickname" json:"nickname"`
	Content   string `bson:"content" json:"content"`
	ReplyTime int64  `bson:"replytime" json:"replytime"`
}

type Help struct {
	HelpID      int64  `bson:"helpID", json:"helpID"`
	Username    string `bson:"username" json:"username"`
	Nickname    string `bson:"nickname" json:"nickname"`
	LevelID     string `bson:"levelID" json:"levelID"`
	Content     string `bson:"content" json:"content"`
	PostTime    int64  `bson:"posttime" json:"posttime"`
	ReplyCount  int    `bson:"replycount" json:"replycount"`
	ReplyTime   int64  `bson:"replytime" json:"replytime"`
	Image       string `bson:"image" json:"image"`
	BoyClothes  string `bson:"bc" json:"bc"`
	GirlClothes string `bson:"gc" json:"gc"`

	Followed bool  `bson:"followed" json:"followed"`
	FollowID int64 `bson:"followID" json:"followID"`
}

type RawHelp struct {
	HelpID      int64  `bson:"helpID", json:"helpID"`
	Username    string `bson:"username" json:"username"`
	Nickname    string `bson:"nickname" json:"nickname"`
	LevelID     string `bson:"levelID" json:"levelID"`
	Content     string `bson:"content" json:"content"`
	PostTime    int64  `bson:"posttime" json:"posttime"`
	ReplyCount  int    `bson:"replycount" json:"replycount"`
	ReplyTime   int64  `bson:"replytime" json:"replytime"`
	Image       string `bson:"image" json:"image"`
	BoyClothes  string `bson:"bc" json:"bc"`
	GirlClothes string `bson:"gc" json:"gc"`
	Type        int    `bson:"type" json:"type"`
}

type GMHelp struct {
	HelpID   int64
	Content  string
	Image    string
	PostTime string
}
