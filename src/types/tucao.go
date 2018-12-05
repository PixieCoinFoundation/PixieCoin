package types

type Tucao struct {
	ID       int    `bson:"id", json:"id"`
	HelpID   int    `bson:"helpID", json:"helpID"`
	UID      string `bson:"uid", json:"uid"`
	Content  string `bson:"content", json:"content"`
	PostTime int64  `bson:"posttime", json:"posttime"`
}
