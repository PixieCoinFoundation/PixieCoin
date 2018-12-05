package loginServerHub

import (
	"appcfg"
	"constants"
	"encoding/json"
	"gl_db/server"
	"net"
	"strings"
	"sync"
	"time"
	"tools"
)

import (
	. "logger"
	. "types"
	. "types/proto"
)

var GFRequestHandles map[uint16]func(*GFServer, string) = map[uint16]func(*GFServer, string){
	//game server report its ip/port
	S_L_REGISTER_SERVER: Register,

	//game server report its player count
	S_L_UPDATE_STATUS: UpdateGameServerStatus,
}

// 缓存玩家的Token
var (
	tokenMapLock sync.RWMutex
	tokenMap     map[string](chan string)
)

var lastLogPlayerCntToken string

func init() {
	tokenMap = make(map[string](chan string))
}

func startOneGFServer(conn *net.TCPConn, in chan Packet) {
	defer func() {
		if x := recover(); x != nil {
			Err("caught panic in StartOneGFServer", x)
		}
	}()

	addr := conn.RemoteAddr().String()
	server := new(GFServer)
	server.Initialize(conn)

	defer func() {
		Info("begin clear server:", addr)
		server.Clear()
	}()

	for {
		select {
		case msg, ok := <-in:
			if !ok {
				Info("game server down:", addr)
				go tools.SendInternalMail(appcfg.GetString("odps_env_name", "")+"_game_server_down_"+addr, "", nil)
				return
			}
			dispatch(server, msg)
		}
	}
}

func dispatch(server *GFServer, msg Packet) {
	if msg.Header.MsgClass != S_L_MSG_CLASS {
		Err("gf server dispath:wrong msg class")
		return
	}
	handle, ok := GFRequestHandles[msg.Header.MsgType]
	if !ok {
		Err("gf server dispath:wrong msg type")
		return
	}

	handle(server, string(msg.Data))
}

func Register(server *GFServer, data string) {
	var input ServerInfo

	if err := unmarshal(data, &input); err == nil {
		server.SetAddress(input.IP, input.IP2, input.IP3, input.Port)
		RegisterServer(server)
	}
}

func UpdateGameServerStatus(server *GFServer, data string) {
	var input ServerStatus

	if err := unmarshal(data, &input); err == nil {
		addr, _, _ := server.GetAddress()

		if GetServerByAddress(addr) == nil {
			RegisterServer(server)
		}

		server.SetGameServerStatus(input.PlayerCount, input.RequestCount)
		CheckIdleServer()

		ida, _, _, _, _ := GetIdleServer("")
		Info("player:", input.PlayerCount, "request:", input.RequestCount, "server:", addr, server.GetLocalAddress(), "idle:", ida)
	}

	go logPlayerCnt()
}

func LogOutUser(lastServer string, username string) (err error) {
	server := GetServerByAddress(lastServer)
	if server != nil {
		Info("log out:", username, lastServer)
		return server.LogOutUser(username)
	} else {
		Err("can't find server:", lastServer)
		return nil
	}
}

func LogOutInAllServers(username string) (err error) {
	return LogOutAll(username)
}

func logPlayerCnt() {
	token := time.Now().Format("200601021504")
	if lastLogPlayerCntToken != token && (strings.HasSuffix(token, "0") || strings.HasSuffix(token, "5")) {
		v, d := GetAllUserCnt()

		//gm log
		c := OnlineExtra{
			OnlineNum:     v,
			ServerAddress: appcfg.GetAddress(),
		}
		data, _ := json.Marshal(c)
		GMLog(constants.C1_SYSTEM, constants.C2_ONLINE, constants.C3_DEFAULT, "", string(data))

		server.AddPlayerCntLog(token, v, d)
		lastLogPlayerCntToken = token
	}
}

func unmarshal(data string, v interface{}) (err error) {
	if err = json.Unmarshal([]byte(data), v); err != nil {
		Err("gf dispatch unmarshal error:", data, err)
	}
	return
}

func marshal(v interface{}) string {
	if jsonbyte, err := json.Marshal(v); err != nil {
		Err("gf dispatch marshal error", err)
		return ""
	} else {
		return string(jsonbyte)
	}
}
