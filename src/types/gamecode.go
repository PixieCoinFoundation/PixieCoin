package types

//KorGameCode 韩国游戏码
type KorGameCode struct {
	Code     string
	Username string
	SendTime string `json:"time"`
	Status   bool
}
