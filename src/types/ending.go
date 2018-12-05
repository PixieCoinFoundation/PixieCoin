package types

type Ending struct {
	Username string `bson:"username", json:"username"`
	Endings  string `bson:"endings", json:"endings"`
}
