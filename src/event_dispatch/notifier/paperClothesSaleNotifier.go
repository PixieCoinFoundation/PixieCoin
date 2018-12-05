package notifier

import (
	"db_pixie/paper"
	. "event_dispatch/event"
)

type PaperClothesSaleObserver interface {
	OnNotify(PaperClothesSaleEvent)
}

type PaperClothesSaleNotifier struct {
	Observers map[PaperClothesSaleObserver]struct{}
}

func (o *PaperClothesSaleNotifier) Register(l PaperClothesSaleObserver) {
	o.Observers[l] = struct{}{}
}

func (p *PaperClothesSaleNotifier) Notify(e PaperClothesSaleEvent) {
	paper.LogPaperClothesSale(e.SaleUsername, e.LandID, e.Price, e.PriceType, e.Time)

	for o := range p.Observers {
		o.OnNotify(e)
	}
}
