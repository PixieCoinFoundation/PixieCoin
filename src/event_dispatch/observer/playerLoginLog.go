package observer

import (
	"appcfg"
	"constants"
	"encoding/json"
	. "event_dispatch/event"
	. "logger"
	"time"
	. "types"
)

type PlayerLoginLogObserver struct{}

func (o *PlayerLoginLogObserver) OnNotify(e PlayerLoginEvent) {
	le := LoginReqExtra{
		Username:    e.Player.Username,
		AccessToken: e.Player.AccessToken,
		DeviceToken: e.Player.DeviceToken,
		UID:         e.Player.UID,
		// GFToken:             e.GFToken,
		Platform:            e.Player.Platform,
		DeviceId:            e.Input.DeviceId,
		Channel:             e.Input.Channel,
		ThirdChannel:        e.Input.ThirdChannel,
		GuestLogin:          e.Input.GuestLogin,
		DeviceModelType:     e.Input.DeviceModelType,
		KorRevertUnregister: e.Input.KorRevertUnregister,
		RemoteAddr:          e.RemoteAddr,
		ThirdUsername:       e.Player.ThirdUsername,
		ServerAddress:       appcfg.GetAddress(),
		Nickname:            e.Player.Nickname,
		Level:               e.Player.GetLevel(),
		RegTime:             time.Unix(e.Player.GetRegTime(), 0).Format("2006-01-02 15:04:05"),
	}

	leb, _ := json.Marshal(le)
	GMLog(constants.C1_PLAYER, constants.C2_LOGIN, constants.C3_DEFAULT, e.Player.Username, string(leb))

	_, utcOffset := time.Now().Zone()
	de := LoginDataSnapshotExtra{
		Nickname:      e.Player.Nickname,
		Level:         e.Player.GetLevel(),
		Money:         e.Player.GetMoney(),
		ThirdChannel:  e.Player.ThirdChannel,
		UTCOffset:     utcOffset,
		TheatreStatus: e.Player.GetTheatreStatus(),
		EndingStatus:  e.Player.GetEndingStatus(),
		Cloth:         e.Player.GetCloth(),
		Button:        e.Player.GetButton(),
	}
	de.DesignerBoyClothesCnt, de.DesignerGirlClothesCnt, de.OfficialBoyClothesCnt, de.OfficialGirlClothesCnt = e.Player.GetAllBoyGirlClothesCnt()
	deb, _ := json.Marshal(de)
	GMLog(constants.C1_PLAYER, constants.C2_LOGIN_SNAPSHOT, constants.C3_DEFAULT, e.Player.Username, string(deb))
}
