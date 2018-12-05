package player_mem

import (
	"constants"
	"errors"
	"fmt"
	"runtime"
	"time"
	"tools"
)

import (
	"appcfg"
	cm "concurrent_map"
	. "logger"
	"shutdown"
)

var (
	players cm.ConcurrentMap
	tokens  cm.ConcurrentMap

	roundTime   int64
	timeOutTime int64

	playerCntCache int64
)

func init() {
	if appcfg.GetServerType() != "" {
		return
	}

	roundTime = appcfg.GetInt64("check_off_line_interval", 10)
	timeOutTime = appcfg.GetInt64("offline_interval", 600)

	players = cm.New()
	tokens = cm.New()

	go checkTimeOutPlayers()
}

func GeneratePlayerByToken(username, deviceToken string, uid int64, platform, deviceId, channel, biliUsername, accessToken, thirdChannel, oldChannel string, korRevertUnregister bool) (p *GFPlayer, err error) {
	// logout之前的玩家，如果有的话
	prePlayer := getPlayerByUsername(username)
	if prePlayer != nil {
		logOutUser1(prePlayer)
	}

	p = GetGFPlayer()

	if success, banned, limited, deleting, deled := p.initialize(username, uid, deviceToken, platform, deviceId, channel, biliUsername, accessToken, thirdChannel, oldChannel, korRevertUnregister); success {
		setPlayer(username, p)
		return
	} else if limited {
		err = errors.New(constants.PLAYER_LIMITED_MSG)
		PutGFPlayer(p)
		return
	} else if banned {
		err = errors.New(constants.PLAYER_BAN_MSG)
		PutGFPlayer(p)
		return
	} else if deleting {
		err = errors.New(constants.PLAYER_DELETING_MSG)
		PutGFPlayer(p)
		return
	} else if deled {
		err = errors.New(constants.PLAYER_DELED_MSG)
		PutGFPlayer(p)
		return
	} else {
		err = errors.New("initialize player err, username:" + username)
		PutGFPlayer(p)
		return
	}

}

func GetGameServerStatus(real bool) (pc int, rc int) {
	requestNum := shutdown.GetProcessingRequestsNum()

	if real {
		count := players.Count()
		playerCntCache = int64(count)
		a, i, d := GetPoolInfo()
		goroutineNum := runtime.NumGoroutine()

		if goroutineNum >= 500 || requestNum >= 10 {
			go tools.SendInternalMail("server wrong", fmt.Sprintf("server %s goroutine %d requests %d", appcfg.GetAddress(), goroutineNum, requestNum), nil)
		}

		Info("player:", count, "goroutine:", goroutineNum, "pool-active:", a, "pool-idle:", i, "pool-destroyed:", d, "requests:", requestNum)
		return count, requestNum
	} else {
		return int(playerCntCache), requestNum
	}
}

func GetPlayerByUsername(username string) *GFPlayer {
	player := getPlayerByUsername(username)

	return player

}

func GetPlayerByUsernameAndToken(username, accessToken string) *GFPlayer {
	player := getPlayerByUsername(username)

	if player != nil && player.AccessToken == accessToken {
		return player
	} else {
		return nil
	}

}

func getPlayerByUsername(username string) *GFPlayer {
	if player, ok := players.Get(username); ok {
		return player.(*GFPlayer)
	} else {
		return nil
	}
}

func setPlayer(username string, player *GFPlayer) {
	players.Set(username, player)
}

func deletePlayer(username string) {
	players.Remove(username)
}

func logOutUser(username string) {
	if player := getPlayerByUsername(username); player != nil {
		Info("log out:", username, "now-heartbeat:", time.Now().Unix()-player.HeartbeatTime)
		logOutUser1(player)
	}
}

func logOutUser1(p *GFPlayer) {
	deletePlayer(p.Username)
	p.Clear()
	PutGFPlayer(p)
}

func DisconnectPlayer(username string) {
	prePlayer := getPlayerByUsername(username)
	if prePlayer != nil {
		logOutUser1(prePlayer)
	}
}

func checkTimeOutPlayers() {
	for {
		for _, v := range players.Items() {
			vp := v.(*GFPlayer)
			if time.Now().Unix()-vp.HeartbeatTime >= timeOutTime {
				logOutUser1(vp)
			}
		}
		time.Sleep(time.Second * time.Duration(roundTime))
	}
}
