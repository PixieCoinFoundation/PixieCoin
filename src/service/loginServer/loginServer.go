package loginServer

import (
	"net"
	"time"
)

import (
	"appcfg"
	. "logger"
	"player_mem"
	"timer"
	. "types"
	. "types/proto"
)

type LoginServer struct {
	conn       *net.TCPConn
	buffer     *Buffer
	stopUpdate chan bool
	buffOb     chan bool
}

func (self *LoginServer) Initialize(conn *net.TCPConn) {
	self.conn = conn

	self.stopUpdate = make(chan bool)
	self.buffOb = make(chan bool)

	max := appcfg.GetInt("server_q_size", DEFAULT_QUEUE_SIZE)
	self.buffer = NewBuffer(conn, self.buffOb, max)

	go self.buffer.Start()

	self.Register()
	go updateStatus(self)
}

func (self *LoginServer) Clear() {
	self.buffer.Close()
	self.stopUpdate <- true
}

func (self *LoginServer) Register() {
	packet := self.genPacket(S_L_REGISTER_SERVER)

	serverInfo := ServerInfo{
		IP:   appcfg.GetString("player_ip", "127.0.0.1"),
		Port: appcfg.GetString("port", ":29991"),
	}
	packet.SetData([]byte(serverInfo.String()))

	err := self.buffer.Send(packet)
	if err != nil {
		Err("Send server status error", err.Error())
	}
}

func (self *LoginServer) SendStatus(count int) {
	packet := self.genPacket(S_L_UPDATE_STATUS)
	serverStatus := ServerStatus{
		PlayerCount: count,
	}
	packet.SetData([]byte(serverStatus.String()))

	err := self.buffer.Send(packet)
	if err != nil {
		Err("Send server status error", err.Error())
	}
}

func (self *LoginServer) SendAccessToken(username string, token string) {
	packet := self.genPacket(S_L_REQUEST_ACCESS_TOKEN)
	resp := RequestAccessTokenResp{
		Username:    username,
		AccessToken: token,
	}
	packet.SetData([]byte(resp.String()))
	err := self.buffer.Send(packet)
	if err != nil {
		Err("Send server status error", err.Error())
	}
}

func (self *LoginServer) genPacket(msgType uint16) Packet {
	var packet Packet
	packet.Header.Tag = PACKET_DEFAULT_TAG
	packet.Header.HeadLen = 10
	packet.Header.MsgClass = S_L_MSG_CLASS
	packet.Header.MsgType = msgType

	return packet
}

func updateStatus(server *LoginServer) {
	// 定时向登陆服务端发送本服务器的信息
	// luaTimer := make(chan int32, 1)
	// timer.Add(int32(TIMER_UPDATE_SERVER_STATUS), time.Now().Unix()+updateInterval, luaTimer)
	for {
		select {
		case <-time.Sleep(constants.GAME_SERVER_REPORT_STATUS_INTERVAL * time.Second):
			server.SendStatus(player_mem.GetPlayerCount(true))

			// timer.Add(int32(TIMER_UPDATE_SERVER_STATUS), time.Now().Unix()+updateInterval, luaTimer)
		case <-server.stopUpdate:
			return
		case <-server.buffOb:
			return
		}
	}
}
