package types

type Scripts struct {
	Username  string `bson:"username", json:"username"`
	ScriptIDs string `bson:"scriptIDs", json:"scriptIDs"`
}
