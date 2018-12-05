package gameServerHub

import (
	"constants"
	"net"
	"time"
)

import (
	"appcfg"
	. "logger"
	"player_mem"
	// "timer"
	. "types"
	. "types/proto"
)

type HeartbeatClient struct {
	buffer     *Buffer
	stopUpdate chan bool
	buffOb     chan bool
}

func (self *HeartbeatClient) Initialize(conn *net.TCPConn) {
	self.stopUpdate = make(chan bool)
	self.buffOb = make(chan bool)

	max := appcfg.GetInt("server_q_size", DEFAULT_QUEUE_SIZE)
	self.buffer = NewBuffer(conn, self.buffOb, max)

	go self.buffer.Start()

	self.Register()
	go updateStatusLoop(self)
}

func (self *HeartbeatClient) Clear() {
	self.buffer.Close()
	self.stopUpdate <- true
}

func (self *HeartbeatClient) Register() {
	packet := self.genPacket(S_L_REGISTER_SERVER)

	serverInfo := ServerInfo{
		IP:   appcfg.GetString("player_ip", "127.0.0.1"),
		IP2:  appcfg.GetString("player_ip_2", "127.0.0.1"),
		IP3:  appcfg.GetString("player_ip_3", "127.0.0.1"),
		Port: appcfg.GetString("port", ":29991"),
	}
	packet.SetData([]byte(serverInfo.String()))

	err := self.buffer.Send(packet)
	if err != nil {
		Err("Send server status error", err.Error())
	}
}

func (self *HeartbeatClient) SendStatus(playerCount int, requestCount int) {
	packet := self.genPacket(S_L_UPDATE_STATUS)
	serverStatus := ServerStatus{
		PlayerCount:  playerCount,
		RequestCount: requestCount,
	}
	packet.SetData([]byte(serverStatus.String()))

	err := self.buffer.Send(packet)
	if err != nil {
		Err("Send server status error", err.Error())
	}
}

func (self *HeartbeatClient) genPacket(msgType uint16) Packet {
	var packet Packet
	packet.Header.Tag = PACKET_DEFAULT_TAG
	packet.Header.HeadLen = 10
	packet.Header.MsgClass = S_L_MSG_CLASS
	packet.Header.MsgType = msgType

	return packet
}

func updateStatusLoop(server *HeartbeatClient) {
	// 定时向登陆服务端发送本服务器的信息
	for {
		select {
		case <-time.After(constants.GAME_SERVER_REPORT_STATUS_INTERVAL * time.Second):
			server.SendStatus(player_mem.GetGameServerStatus(true))
		case <-server.stopUpdate:
			return
		case <-server.buffOb:
			return
		}
	}
}
