package observer

import (
	"appcfg"
	"constants"
	"encoding/json"
	. "event_dispatch/event"
	. "logger"
	"time"
	"tools"
	. "types"
)

type PlayerLoginKorLogObserver struct{}

func (o *PlayerLoginKorLogObserver) OnNotify(e PlayerLoginEvent) {
	if appcfg.GetLanguage() == constants.KOR_LANGUAGE && appcfg.GetBool("upload_kor_log", false) {
		var (
			KOROSCode int
			lkl       LoginKorLog
		)

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

		le.KORGameCode = constants.KOR_GAME_CODE
		le.KOROSCode, _, _ = tools.GetKOROrderInfo(e.Player.Platform)
		lkl = LoginKorLog{
			Guid:           e.Player.UID,
			Uuid:           e.Player.Username,
			Cuid:           e.Player.UID,
			Gamecode:       constants.KOR_GAME_CODE,
			Nickname:       e.Player.Nickname,
			LevYn:          "N",
			UserLv:         e.Player.GetLevel(),
			Authtoken:      0,
			Oscode:         KOROSCode,
			Tempuserid:     e.Player.Nickname,
			TokenredisDate: le.RegTime,
			Registdate:     le.RegTime,
			Mappingdate:    le.RegTime,
			IP:             e.RemoteAddr,
		}

		//韩国版本日志插入
		lkldata, _ := json.Marshal(lkl)
		GMLogKor(constants.C1_PLAYER, constants.C2_LOGIN, constants.C3_LOGIN_LOG_KOR, e.Player.Username, string(lkldata))
		//韩国版本日志1插入
		lkl1 := LoginKorLog1{
			Logindate:         time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"),
			Logindateseq:      0,
			Cuid:              e.Player.UID,
			Isloginsuccessdiv: 1,
			Platformuserid:    e.Player.ThirdUsername,
			Loginrootcode:     1,
			GameCode:          constants.KOR_GAME_CODE,
			Deviceid:          0,
			Platformcode:      7,
			Oscode:            KOROSCode,
			Ostopverno:        0,
		}
		lkldata1, _ := json.Marshal(lkl1)
		GMLogKor(constants.C1_PLAYER, constants.C2_LOGIN, constants.C3_LOGIN1_LOG_KOR, e.Player.Username, string(lkldata1))
	}

}
