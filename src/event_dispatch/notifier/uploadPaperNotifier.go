package notifier

import (
	"constants"
	"encoding/json"
	. "event_dispatch/event"
	. "logger"
	. "types"
)

type UploadPaperObserver interface {
	OnNotify(UploadPaperEvent)
}

type UploadPaperNotifier struct {
	Observers map[UploadPaperObserver]struct{}
}

func (o *UploadPaperNotifier) Register(l UploadPaperObserver) {
	o.Observers[l] = struct{}{}
}

func (p *UploadPaperNotifier) Notify(e UploadPaperEvent) {
	//log
	de := UploadPaperLogExtra{
		UploadPaperReq: e.Input,
		PaperID:        e.PaperID,
	}
	data, _ := json.Marshal(de)
	GMLog(constants.C1P_PLAYER, constants.C2P_PAPER, constants.C3P_UPLOAD, e.Player.Username, string(data))

	for o := range p.Observers {
		o.OnNotify(e)
	}
}
