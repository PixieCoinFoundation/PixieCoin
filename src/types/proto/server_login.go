package proto

import (
	"encoding/json"
)

const (
	S_L_MSG_CLASS uint16 = 1000
)

const (
	S_L_REGISTER_SERVER uint16 = S_L_MSG_CLASS + 1
	S_L_UPDATE_STATUS   uint16 = S_L_MSG_CLASS + 2

	S_L_REQUEST_ACCESS_TOKEN uint16 = S_L_MSG_CLASS + 3
	S_L_LOGOUT_USER          uint16 = S_L_MSG_CLASS + 4
)

// 设置服务器地址
type ServerInfo struct {
	IP   string `json:"ip"`
	IP2  string
	IP3  string
	Port string `json:"port"`
}

func (self *ServerInfo) String() string {
	str, _ := json.Marshal(self)

	return string(str)
}

// 更新服务器状态
type ServerStatus struct {
	PlayerCount  int `json:"playerCount"`
	RequestCount int `json:"requestCount"`
}

func (self *ServerStatus) String() string {
	str, _ := json.Marshal(self)

	return string(str)
}

// 获取AccessToken
type RequestAccessTokenInfo struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DeviceToken string `json:"deviceToken"`
}

func (self *RequestAccessTokenInfo) String() string {
	str, _ := json.Marshal(self)

	return string(str)
}

type RequestAccessTokenResp struct {
	Username    string `json:"username"`
	AccessToken string `json:"token"`
}

func (self *RequestAccessTokenResp) String() string {
	str, _ := json.Marshal(self)

	return string(str)
}

type LogOutUserRequest struct {
	Username string `json:"username"`
}

func (self *LogOutUserRequest) String() string {
	str, _ := json.Marshal(self)

	return string(str)
}
