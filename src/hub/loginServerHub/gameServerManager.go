package loginServerHub

import (
	"appcfg"
	"constants"
	"encoding/json"
	"fmt"
	"gl_dao"
	. "logger"
	"sync"
	"time"
)

var (
	servers map[string]*GFServer
	sMutex  sync.RWMutex

	idleServer *GFServer
	idleMutex  sync.RWMutex

	blackServerList     map[string]int
	blackServerListLock sync.Mutex

	blackAccountList     map[string]int
	blackAccountListLock sync.Mutex
)

var gameServerLBAddress1 string
var gameServerLBAddress2 string
var gameServerLBAddress3 string

func init() {
	servers = make(map[string]*GFServer)

	if appcfg.GetServerType() == constants.SERVER_TYPE_GL {
		if appcfg.GetBool("use_server_lb", false) {
			gameServerLBAddress1 = appcfg.GetString("game_server_lb_1", "")
			gameServerLBAddress2 = appcfg.GetString("game_server_lb_2", "")
			gameServerLBAddress3 = appcfg.GetString("game_server_lb_3", "")
		}

		if gameServerLBAddress1 == "" && gameServerLBAddress2 == "" && gameServerLBAddress3 == "" {
			go refreshBlack()
		} else {
			nea := gameServerLBAddress1
			if nea == "" {
				nea = gameServerLBAddress2
			}
			if nea == "" {
				nea = gameServerLBAddress3
			}

			if nea == "" {
				panic("game server lb address all empty")
			}

			if gameServerLBAddress1 == "" {
				gameServerLBAddress1 = nea
			}
			if gameServerLBAddress2 == "" {
				gameServerLBAddress2 = nea
			}
			if gameServerLBAddress3 == "" {
				gameServerLBAddress3 = nea
			}

			Info("game server lb:", gameServerLBAddress1, gameServerLBAddress2, gameServerLBAddress3)
		}
	}
}

func GetAllUserCnt() (int, string) {
	sMutex.RLock()
	defer sMutex.RUnlock()

	var desc string
	cnt := 0
	for _, s := range servers {
		c := s.getPlayerCount()

		sn, _, _ := s.GetAddress()
		desc += fmt.Sprintf("[%s]  :  [%d]<br>", sn, c)
		cnt += c
	}

	return cnt, desc
}

func RegisterServer(server *GFServer) {
	addr, _, _ := server.GetAddress()
	setServer(addr, server)
	CheckIdleServer()
}

func GetIdleServer(username string) (string, string, string, *GFServer, bool) {
	if gameServerLBAddress1 != "" {
		return gameServerLBAddress1, gameServerLBAddress2, gameServerLBAddress3, nil, true
	}

	if !blackGameServerEmpty() && AccountIsBlack(username) {
		sMutex.RLock()
		for _, v := range servers {
			if serverInblackServerList(v.IP) {
				defer sMutex.RUnlock()
				return v.IP, v.IP2, v.IP3, v, true
			}
		}
		sMutex.RUnlock()
	}

	idleMutex.RLock()
	defer idleMutex.RUnlock()

	if idleServer != nil {
		return idleServer.IP, idleServer.IP2, idleServer.IP3, idleServer, true
	} else {
		return "", "", "", nil, false
	}
}

func GetServerByAddress(address string) *GFServer {
	return getServer(address)
}

func LogOutAll(username string) (err error) {
	sMutex.RLock()
	defer sMutex.RUnlock()
	for _, v := range servers {
		if err = v.LogOutUser(username); err != nil {
			return
		}
	}
	return
}

func CheckIdleServer() {
	var server *GFServer

	sMutex.RLock()
	for _, v := range servers {
		if serverInblackServerList(v.IP) {
			continue
		}

		gap := time.Now().Unix() - v.LastReportTime
		if gap > 2*constants.GAME_SERVER_REPORT_STATUS_INTERVAL {
			address, _, _ := v.GetAddress()
			Err("game server not available:", address, "for seconds:", gap)
			continue
		}

		if server == nil {
			server = v
		} else {
			if server.getPlayerCount() > v.getPlayerCount() {
				server = v
			}
		}
	}
	sMutex.RUnlock()

	idleMutex.Lock()
	idleServer = server
	idleMutex.Unlock()
}

func RemoveServer(server *GFServer) {
	addr, _, _ := server.GetAddress()
	removeServer(addr, server)
	CheckIdleServer()
}

func setServer(address string, server *GFServer) {
	Info("set server:", address, server.IP)
	sMutex.Lock()
	defer sMutex.Unlock()

	servers[address] = server
}

func removeServer(address string, server *GFServer) {
	Info("remove server:", address, server.IP)
	sMutex.Lock()
	defer sMutex.Unlock()

	delete(servers, address)
}

func getServer(address string) *GFServer {
	sMutex.RLock()
	defer sMutex.RUnlock()

	return servers[address]
}

func serverInblackServerList(s string) bool {
	blackServerListLock.Lock()
	defer blackServerListLock.Unlock()

	if len(blackServerList) <= 0 {
		return false
	}

	if _, ok := blackServerList[s]; ok {
		return true
	}

	return false
}

func blackGameServerEmpty() bool {
	blackServerListLock.Lock()
	defer blackServerListLock.Unlock()

	return len(blackServerList) <= 0
}

func AccountIsBlack(username string) bool {
	blackAccountListLock.Lock()
	defer blackAccountListLock.Unlock()

	if username == "" {
		return false
	}

	if len(blackAccountList) <= 0 {
		return false
	}

	if _, ok := blackAccountList[username]; ok {
		return true
	}

	return false
}

func refreshBlack() {
	for {
		nm := getBlackGameServerList()
		am := getBlackAccountList()

		blackServerListLock.Lock()
		blackServerList = nm
		blackServerListLock.Unlock()

		blackAccountListLock.Lock()
		blackAccountList = am
		blackAccountListLock.Unlock()

		time.Sleep(60 * time.Second)
	}
}

func getBlackGameServerList() map[string]int {
	res := make([]string, 0)
	resm := make(map[string]int)

	var s string
	if err := gl_dao.GetConfigStmt.QueryRow(constants.BLACK_GAME_SERVER_LIST_KEY).Scan(&s); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("no black game server")
		} else {
			Err(err)
		}
	} else if err := json.Unmarshal([]byte(s), &res); err != nil {
		Err(err)
	} else {
		for _, s := range res {
			Info("black game server:", s)
			resm[s] = 1
		}
	}

	return resm
}

func getBlackAccountList() map[string]int {
	res := make([]string, 0)
	resm := make(map[string]int)

	var s string
	if err := gl_dao.GetConfigStmt.QueryRow(constants.BLACK_ACCOUNT_LIST_KEY).Scan(&s); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("no black account")
		} else {
			Err(err)
		}
	} else if err := json.Unmarshal([]byte(s), &res); err != nil {
		Err(err)
	} else {
		for _, s := range res {
			Info("black account:", s)
			resm[s] = 1
		}
	}

	return resm
}
