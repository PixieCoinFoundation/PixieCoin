package gameServerHub

import (
	"constants"
	"encoding/json"
	"io"
	"net"
	"sync"
)

import (
	"appcfg"
	"dao"
	. "logger"
	"player_mem"
	"time"
	. "types"
	. "types/proto"
)

const (
	//if login server haven't heartbeat for 5 minutes.game server won't connect.
	LOGIN_SERVER_ALIVE_TIME_GAP = 181 //seconds
	CHECK_LOGIN_SERVER_TIME     = 60  //second
)

var Handles map[uint16]func(string) = map[uint16]func(string){
	S_L_LOGOUT_USER: logOutUser, //gameServer只接受loginServer发出的踢人请求
}

var connectedLoginServers map[string]int
var connectedLoginServerLock sync.Mutex

func init() {
	if appcfg.GetServerType() == "" {
		connectedLoginServers = make(map[string]int)

		loginServers := getLoginServers()

		if len(loginServers) <= 0 {
			panic("no login servers!!!!")
		} else {
			for _, loginServer := range loginServers {
				if !connectIfNotConnected(loginServer) {
					panic("connect login server error")
				}
			}
		}

		go refreshLegalLoginServers()
	}
}

//定时刷新合法的loginServer列表
func refreshLegalLoginServers() {
	for {
		time.Sleep(time.Duration(CHECK_LOGIN_SERVER_TIME) * time.Second)

		res := getLoginServers()
		Info("legal login servers:", res, "current connected login servers:", connectedLoginServers)

		for _, addr := range res {
			connectIfNotConnected(addr)
		}
	}
}

func connectIfNotConnected(serverAddr string) (connected bool) {
	connectedLoginServerLock.Lock()
	defer connectedLoginServerLock.Unlock()

	if _, ok := connectedLoginServers[serverAddr]; ok {
		//already connected.do nothing
		connected = true
	} else {
		//connect
		if connectLoginServer(serverAddr) {
			Info("connected to login server:", serverAddr)
			connectedLoginServers[serverAddr] = 1
			connected = true
		} else {
			Info("err connect to login server:", serverAddr)
			connected = false
		}
	}

	return
}

func connectLoginServer(serverAddr string) bool {
	addr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		Err(err)
		return false
	}

	// 连接玩家服务器
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		Err("game server connect to login server failed:", serverAddr, err)
		return false
	}

	go handleLoginServer(conn)

	return true
}

func handleLoginServer(conn *net.TCPConn) {
	defer func() {
		if x := recover(); x != nil {
			Err("caught panic in handleLoginServer", x)
		}
	}()

	defer conn.Close()

	header := make([]byte, 10)
	msgsFromLoginServer := make(chan Packet, DEFAULT_QUEUE_SIZE)

	go startCommunicateWithLoginServer(conn, msgsFromLoginServer)

	for {
		n, err := io.ReadFull(conn, header)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			Err("error receiving header server:", err)
			break
		}

		var packet Packet
		packet.SetHead(header)

		data := make([]byte, packet.Header.DataLen)
		n, err = io.ReadFull(conn, data)
		if err != nil {
			Err("error receiving payload server:", err)
			break
		}
		packet.SetRawData(data[0:])

		msgsFromLoginServer <- packet
	}

	close(msgsFromLoginServer)
}

func startCommunicateWithLoginServer(conn *net.TCPConn, in chan Packet) {
	defer func() {
		if x := recover(); x != nil {
			Err("caught panic in StartLoginServer", x)
		}
	}()

	ra := conn.RemoteAddr().String()
	//send heart beat to player server
	hbClient := new(HeartbeatClient)
	hbClient.Initialize(conn)

	defer func() {
		hbClient.Clear()

		//mark remote server not online
		Info("login server:", ra, "disconnected..")
		removeConnectedLoginServer(ra)
	}()

	//process msg from player server
	for {
		select {
		case msg, ok := <-in:
			if !ok {
				Err("login server down:", ra)
				return
			}
			dispatch(msg)
		}
	}
}

func getLoginServers() (ret []string) {
	now := time.Now().Unix()
	expireTime := now - LOGIN_SERVER_ALIVE_TIME_GAP

	if rows, err := dao.GetServerStmt.Query(constants.LOGIN_FOR_GAME_SERVER_TYPE, expireTime); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var addr string
			if err = rows.Scan(&addr); err != nil {
				Err(err)
				return
			} else {
				if addr != "" {
					ret = append(ret, addr)
				}
			}
		}

		return
	}
}

func removeConnectedLoginServer(serverAddr string) {
	connectedLoginServerLock.Lock()
	defer connectedLoginServerLock.Unlock()

	delete(connectedLoginServers, serverAddr)
}

func dispatch(msg Packet) {
	if msg.Header.MsgClass != S_L_MSG_CLASS {
		Err("wrong msg class from login server:", msg.Header.MsgClass)
		return
	}
	handle, ok := Handles[msg.Header.MsgType]
	if !ok {
		Err("no handle for msg type:", msg.Header.MsgType)
		return
	}

	handle(string(msg.Data))
}

func logOutUser(data string) {
	var input LogOutUserRequest
	if err := unmarshal(data, &input); err == nil {
		Info("asked to disconnect player:", input.Username)
		player_mem.DisconnectPlayer(input.Username)
	}
}

func unmarshal(data string, v interface{}) (err error) {
	if err = json.Unmarshal([]byte(data), v); err != nil {
		Err("gf dispatch unmarshal error:", data, ", error:", err.Error())
	}
	return
}

func marshal(v interface{}) []byte {
	if jsonbyte, err := json.Marshal(v); err != nil {
		Err("gf dispatch marshal error:", err.Error())
		return nil
	} else {
		return jsonbyte
	}
}
