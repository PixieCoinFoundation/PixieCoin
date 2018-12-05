package event_dispatch

import (
	. "event_dispatch/event"
	. "pixie_contract/api_specification"
	. "player_mem"
)

func NotifyUploadPaper(p *GFPlayer, in *UploadPaperReq, pid int64) {
	uploadPaperNotifier.Notify(UploadPaperEvent{Player: p, Input: in, PaperID: pid})
}

func NotifyDeletePaper(p *GFPlayer, in *DeletePaperReq) {
	deletePaperNotifier.Notify(DeletePaperEvent{Input: in, Player: p})
}

func NotifyPlayerLogin(req *PixieLoginReq, remoteAddr string, p *GFPlayer) {
	loginNotifier.Notify(PlayerLoginEvent{Input: req, RemoteAddr: remoteAddr, Player: p})
}

func NotifyNonofficialPaperClothesSale(saleUsername string, landID, paperID int64, buyCount, priceType int, price float64, p *GFPlayer, t int64) {
	pcse := PaperClothesSaleEvent{
		SaleUsername: saleUsername,
		LandID:       landID,
		PaperID:      paperID,
		BuyCount:     buyCount,
		PriceType:    priceType,
		Price:        price,

		Player: p,
		Time:   t,
	}
	paperClothesSaleNotifier.Notify(pcse)
}
