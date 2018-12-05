package notifier

import (
	"constants"
	"encoding/json"
	. "event_dispatch/event"
	. "logger"
	. "types"
	// . "zk_manager"
)

type DeletePaperObserver interface {
	OnNotify(DeletePaperEvent)
}

type DeletePaperNotifier struct {
	Observers map[DeletePaperObserver]struct{}
}

func (o *DeletePaperNotifier) Register(l DeletePaperObserver) {
	o.Observers[l] = struct{}{}
}

func (p *DeletePaperNotifier) Notify(e DeletePaperEvent) {
	//del zk lock
	// DelPixiePaperLock(e.Input.PaperID)

	//log
	de := DeletePaperLogExtra{
		DeletePaperReq: e.Input,
	}
	data, _ := json.Marshal(de)
	GMLog(constants.C1P_PLAYER, constants.C2P_PAPER, constants.C3P_DELETE, e.Player.Username, string(data))

	for o := range p.Observers {
		o.OnNotify(e)
	}
}
