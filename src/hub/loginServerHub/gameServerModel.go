package loginServerHub

import (
	"fmt"
	"net"
	"sync"
	"time"
)

import (
	. "logger"
	. "types"
	. "types/proto"
)

type GFServerStatus int

const (
	RUNNING  GFServerStatus = 0 // 运行
	MAINTAIN GFServerStatus = 1 // 维护
)

type GFServer struct {
	IP             string `json:"ip, omitempty"` // ip
	IP2            string
	IP3            string
	Port           string         `json:"port, omitempty"`   // 端口
	PlayerCount    int            `json:"count, omitempty"`  // 玩家数目
	RequestCount   int            `json:"requestCount"`      //请求数
	Status         GFServerStatus `json:"status, omitempty"` // 运行状态
	Msg            string         `json:"msg, omitempty"`    // 当前服务器状态的描述
	LastReportTime int64

	TmpPlayerCount int `json:"-"` // gfServer同步玩家数量后，又登陆的玩家数量记录，每次同步玩家

	// 各个属性的互斥
	addressMutex  sync.Mutex `json:"-"`
	countMutex    sync.Mutex `json:"-"`
	tmpCountMutex sync.Mutex `json:"-"`
	statusMutex   sync.Mutex `json:"-"`

	conn   *net.TCPConn
	buff   *Buffer
	buffOb chan bool
}

func (self *GFServer) resetStatus() {
	self.PlayerCount = 0
	self.TmpPlayerCount = 0
	self.Status = RUNNING
	self.Msg = ""
}

func (self *GFServer) Initialize(conn *net.TCPConn) {
	self.addressMutex.Lock()
	defer func() {
		self.addressMutex.Unlock()
	}()

	self.conn = conn
	self.buffOb = make(chan bool)
	self.buff = NewBuffer(self.conn, self.buffOb, -1)
	self.LastReportTime = time.Now().Unix()

	go self.Start()
}

func (self *GFServer) Start() {
	defer func() {
		self.Clear()
	}()

	go self.buff.Start()

	for {
		select {
		case <-self.buffOb:
			return
		}
	}
}

func (self *GFServer) Clear() {
	Info("clear self:", self.IP)
	self.buff.Close()

	// 从servers中移除掉自己
	RemoveServer(self)
}

func (self *GFServer) SetAddress(ip string, ip2 string, ip3 string, port string) {
	self.IP = ip
	self.IP2 = ip2
	self.IP3 = ip3
	self.Port = port
}

func (self *GFServer) GetAddress() (string, string, string) {
	self.addressMutex.Lock()
	defer func() {
		self.addressMutex.Unlock()
	}()

	return self.IP, self.IP2, self.IP3
}

func (self *GFServer) GetLocalAddress() string {
	if self.conn != nil {
		return self.conn.RemoteAddr().String()
	} else {
		return ""
	}
}

func (self *GFServer) SetGameServerStatus(playerCount, requestCount int) {
	self.countMutex.Lock()
	defer func() {
		self.countMutex.Unlock()
	}()

	self.PlayerCount = playerCount
	self.RequestCount = requestCount
	self.LastReportTime = time.Now().Unix()

	self.TmpPlayerCount = 0
}

func (self *GFServer) AddTmpPlayerCount() {
	self.tmpCountMutex.Lock()
	defer func() {
		self.tmpCountMutex.Unlock()
	}()
	self.TmpPlayerCount = self.TmpPlayerCount + 1
}

func (self *GFServer) getPlayerCount() int {
	self.countMutex.Lock()
	defer func() {
		self.countMutex.Unlock()
	}()

	count := self.PlayerCount + self.TmpPlayerCount
	return count
}

func (self *GFServer) setStatus(status GFServerStatus, msg string) {
	self.statusMutex.Lock()
	defer func() {
		self.statusMutex.Unlock()
	}()

	self.Status = status
	self.Msg = msg
}

func (self *GFServer) getStatus() (GFServerStatus, string) {
	self.statusMutex.Lock()
	defer func() {
		self.statusMutex.Unlock()
	}()

	status := self.Status
	msg := self.Msg
	return status, msg
}

// ---------------------------------------------------------------------------------------------------
func (self *GFServer) RequestAccessToken(username string, password string, deviceToken string) (err error) {
	var input RequestAccessTokenInfo
	input.Username = username
	input.Password = password
	input.DeviceToken = deviceToken

	packet := self.genPacket(S_L_REQUEST_ACCESS_TOKEN)
	packet.SetData([]byte(input.String()))

	fmt.Println("Request token from:", self.IP, self.Port)

	err = self.buff.Send(packet)
	if err != nil {
		fmt.Println("RequestAccessToken error", err.Error())
	}
	return
}

func (self *GFServer) LogOutUser(username string) (err error) {
	var output LogOutUserRequest
	output.Username = username

	packet := self.genPacket(S_L_LOGOUT_USER)
	packet.SetData([]byte(output.String()))

	if err = self.buff.Send(packet); err != nil {
		fmt.Println("LogOutUser error", err.Error())
		return
	}

	return
	// if err != nil {
	// 	fmt.Println("LogOutUser error", err.Error())
	// }
}

func (self *GFServer) genPacket(msgType uint16) Packet {
	var packet Packet
	packet.Header.Tag = PACKET_DEFAULT_TAG
	packet.Header.HeadLen = 10
	packet.Header.MsgClass = S_L_MSG_CLASS
	packet.Header.MsgType = msgType

	return packet
}
