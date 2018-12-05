package notifier

import (
	"db_pixie/global"
	. "event_dispatch/event"
	// . "zk_manager"
)

type PlayerLoginObserver interface {
	OnNotify(PlayerLoginEvent)
}

type PlayerLoginNotifier struct {
	Observers map[PlayerLoginObserver]struct{}
}

func (o *PlayerLoginNotifier) Register(l PlayerLoginObserver) {
	o.Observers[l] = struct{}{}
}

func (p *PlayerLoginNotifier) Notify(e PlayerLoginEvent) {
	// DelPlayerLoginLock(e.Player.GameThirdUsername)

	global.AddPlayerLoginForRandomVisit(e.Player.Username)

	for o := range p.Observers {
		o.OnNotify(e)
	}
}
