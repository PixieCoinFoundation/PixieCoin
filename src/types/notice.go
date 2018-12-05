package types

type Notice struct {
	NoticeID     int64  `json:"NoticeID, omitempty"`
	FromUsername string `json:"FromUsername, omitempty"`
	FromUID      int64  `json:"FromUID, omitempty"`
	FromNickname string `json:"FromNickname, omitempty"`
	Type         int    `json:"Type, omitempty"` //1:add friend 2:clothes piece swap 3.clothes piece swap response 4:add friend response
	Content      string `json:"Content, omitempty"`
	SendTime     int64
}

type SwapCloPieceContent struct {
	FromCloID       string
	FromCloPart     int
	FromCloPieceCnt int
	ToCloID         string
	ToCloPart       int
	ToCloPieceCnt   int
}

// type NewPieceInfo struct {
// 	FromUsername string
// 	GotCloID     int
// 	GotPart      int
// 	GotCount     int
// 	LossCloID    int
// 	LossPart     int
// 	LossCount    int
// }
