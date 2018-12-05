package gf

import (
// . "logger"
// "service/servers"
// . "types"
)

// type RegisterServerRequest struct {
// 	IP   string `json:"ip, omitempty"`
// 	Port string `json:"port, omitempty"`
// }

// type RegisterServerResp struct {
// 	Status ErrorCode `json:"status, omitempty"`
// }

// func RegisterServer(address string, params string, result chan string, errflag chan int) {
// 	var input RegisterServerRequest
// 	var retData RegisterServerResp

// 	if err := unmarshal(params, &input, errflag); err == nil {
// 		servers.RegisterServer(input.IP, input.Port)
// 		retData.Status = NO_ERROR
// 		Debug("服务器：", address, "连接")
// 	}

// 	result <- marshal(retData, errflag)
// }

// type UpdateServerPlayerCountRequest struct {
// 	Count int `json:"count, omitempty"`
// }

// type UpdateServerPlayerCountResp struct {
// 	Status ErrorCode `json:"status, omitempty"`
// }

// func UpdateServerPlayerCount(address string, params string, result chan string, errflag chan int) {
// 	var input UpdateServerPlayerCountRequest
// 	var retData UpdateServerPlayerCountResp

// 	if err := unmarshal(params, &input, errflag); err == nil {
// 		servers.UpdateServerPlayerCount(address, input.Count)
// 		retData.Status = NO_ERROR

// 		Debug("Server：", address, "updated, player count:", input.Count)
// 	}

// 	result <- marshal(retData, errflag)
// }
