package types

type Follow struct {
	Username string `bson:"username", json:"username"`
	HelpIDs  string `bson:"helpIDs", json:"helpIDs"`
}
