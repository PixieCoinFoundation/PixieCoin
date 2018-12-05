package event_dispatch

import (
	. "event_dispatch/notifier"
	. "event_dispatch/observer"
)

//上传paper事件
var uploadPaperNotifier UploadPaperNotifier

//删除paper事件
var deletePaperNotifier DeletePaperNotifier

//玩家登录事件
var loginNotifier PlayerLoginNotifier

//服装销售事件
var paperClothesSaleNotifier PaperClothesSaleNotifier

func init() {
	initGlobalEvents()
	initPaperEvents()
}

func initGlobalEvents() {
	loginNotifier = PlayerLoginNotifier{
		Observers: map[PlayerLoginObserver]struct{}{},
	}
	//国服日志
	loginNotifier.Register(&PlayerLoginLogObserver{})
	//韩服日志
	loginNotifier.Register(&PlayerLoginKorLogObserver{})
}

func initPaperEvents() {
	uploadPaperNotifier = UploadPaperNotifier{
		Observers: map[UploadPaperObserver]struct{}{},
	}

	deletePaperNotifier = DeletePaperNotifier{
		Observers: map[DeletePaperObserver]struct{}{},
	}

	paperClothesSaleNotifier = PaperClothesSaleNotifier{
		Observers: map[PaperClothesSaleObserver]struct{}{},
	}
}
