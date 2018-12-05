package event

import (
	. "pixie_contract/api_specification"
	. "player_mem"
)

type UploadPaperEvent struct {
	Input   *UploadPaperReq
	Player  *GFPlayer
	PaperID int64
}

type DeletePaperEvent struct {
	Input  *DeletePaperReq
	Player *GFPlayer
}

type PlayerLoginEvent struct {
	Input      *PixieLoginReq
	RemoteAddr string
	Player     *GFPlayer
}

type PaperClothesSaleEvent struct {
	SaleUsername string
	LandID       int64
	PaperID      int64
	BuyCount     int

	PriceType int
	Price     float64

	Player *GFPlayer
	Time   int64
}
