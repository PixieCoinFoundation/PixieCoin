package bilibiliHandle

import (
	"appcfg"
	"constants"
	"encoding/json"
	"fmt"
	. "gl_db/server"
	"gl_db/users"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

import (
	"encrypt/simple"
	"encrypt/xxtea"
	gf "hub/loginServerHub"
	. "logger"
	"service/bilibili"
	"service/version"
	"snappy"
	. "types"
	. "zk_manager"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type biliLoginReq struct {
	PlayerInfo BiliPlayer `json:"biliplayer"`
	Platform   string     `json:"platform"`
}

// 处理用户登陆
type biliLoginResp struct {
	Status             ErrorCode `json:"status"`
	Content            string    `json:"content"`
	Address            string    `json:"address"`            // 逻辑服务器地址
	Address2           string    `json:"address2"`           // 逻辑服务器地址
	Address3           string    `json:"address3"`           // 逻辑服务器地址
	FileServerAddress  string    `json:"fileServerAddress"`  // 文件服务器地址
	FileServerAddress2 string    `json:"fileServerAddress2"` // 文件服务器地址
	FileServerAddress3 string    `json:"fileServerAddress3"` // 文件服务器地址

	BiliCode int    `json:"biliCode"`     // SDKServer返回的code
	Uid      string `json:"biliUid"`      // SDKServer返回的b站用户uid
	Nickname string `json:"biliNickname"` // SDKServer返回的b站用户昵称

	PX_TOKEN string

	PX_UID int64
}

// 玩家在客户端调用sdk登陆之后，将获取到的玩家信息发送到服务端
// 服务端向bilibili的SDK服务器验证accessToken，并将验证结果返回客户端
func BilibiliLogin(w http.ResponseWriter, r *http.Request) {
	var responseData biliLoginResp
	defer func() {
		jsonByte, _ := json.Marshal(responseData)
		w.Write(encode(jsonByte))
	}()

	defer func() {
		r.Body.Close()
	}()

	if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
		responseData = biliLoginResp{
			Status:  ERR_CHANNEL_NOT_SUPPORT,
			Content: "channel not support",
		}
		return
	}

	data, readError := ioutil.ReadAll(r.Body)
	if readError != nil {
		responseData = biliLoginResp{
			Status:  ERR_UNKNOW,
			Content: "read data error",
		}

		return
	}

	data = simple.Unpack(data)
	data = xxtea.Decrypt(data, constants.XXTEA_MSG_KEY)
	uncompressed, snappyErr := snappy.Decode(nil, data)
	if snappyErr != nil {
		Err("snappy error:", snappyErr.Error())
		responseData = biliLoginResp{
			Status:  ERR_BILI_SESSION_VERIFY_FORMAT,
			Content: "format err",
		}
		return
	}

	var input biliLoginReq
	err := json.Unmarshal([]byte(uncompressed), &input)
	if err != nil {
		Err(err.Error())
		responseData = biliLoginResp{
			Status:  ERR_BILI_SESSION_VERIFY_FORMAT,
			Content: "format error",
		}
		return
	}

	var result SessionVerifyResult
	if input.Platform != "android_channel" {
		var err error
		result, err = bilibili.SessionVerify(input.Platform, input.PlayerInfo.AccessToken, input.PlayerInfo.Uid)
		if err != nil {
			Err("bilibili session verify err:", err)
			responseData = biliLoginResp{
				Status:  ERR_BILI_SESSION_VERIFY,
				Content: "bilibili server error",
			}
			return
		}
	}

	now := time.Now()
	responseData.BiliCode = result.Code
	if result.Code == 0 {
		if input.Platform != "android_channel" && input.PlayerInfo.Uid != strconv.FormatInt(result.OpenID, 10) {
			responseData.Status = ERR_BILI_SESSION_VERIFY_FAILED
		} else {
			responseData.Status = NO_ERROR

			if ok, path := LockPlayerLogin(input.PlayerInfo.Uid); ok {
				defer Unlock(path)

				// 返回相对空闲的服务器ip
				server, server2, server3, sp, succeed := gf.GetIdleServer(input.PlayerInfo.Uid)
				if succeed {
					// todo:
					// 获取上一次玩家登陆时的游戏服务器ip
					lastServer, lastLoginTime, uid, db_addr, info2, newPlayer, err := users.QueryOrRegisterPlayer(input.PlayerInfo.Uid, input.PlayerInfo.Username, server)
					if err != nil {
						responseData.Status = ERR_QUERY_LAST_LOGIN_INFO
					} else if !newPlayer && now.Unix()-lastLoginTime <= 5 && now.Unix()-lastLoginTime >= -5 {
						responseData.Status = ERR_LOGIN_TOO_FAST
					} else {
						gfUsername := fmt.Sprintf("%s%d", appcfg.GetString("un_prefix", ""), uid)

						if sp == nil {
							if err = gf.LogOutInAllServers(gfUsername); err != nil {
								responseData.Status = ERR_LOGOUT_ALL
							}
						} else if lastServer != server && lastServer != "" {
							if err = gf.LogOutUser(lastServer, gfUsername); err != nil {
								//check if game server alive
								if t, err := QueryHeartbeatTime(lastServer); err != nil {
									//err
									responseData.Status = ERR_LOGOUT_USER
								} else if time.Now().Unix()-t < constants.HEARTBEAT_GAP {
									//LAST SERVER ALIVE.ERR!
									responseData.Status = ERR_LOGOUT_LAST_SERVER_ALIVE
								} else {
									//LAST SERVER DOWN.LOGIN OK.DO NOTHING
									Info("log out user failed at last server:", lastServer, "it's dead")
								}
							}
						}

						if responseData.Status == NO_ERROR {
							//DO LOGIN
							if _, err = users.LoginBilibili(input.PlayerInfo.Uid, input.PlayerInfo.Username, server, lastServer, lastLoginTime, info2); err != nil {
								responseData.Status = ERR_LOGIN_GF
							} else {
								Info("bilibili player:", input.PlayerInfo.Uid, gfUsername, "login success:", server)
								responseData.Content = input.PlayerInfo.AccessToken
								responseData.Address = server3
								responseData.Address2 = server
								responseData.Address3 = server2
								responseData.Uid = gfUsername
								responseData.Nickname = result.Uname
								responseData.PX_UID = uid

								responseData.PX_TOKEN = db_addr

								responseData.FileServerAddress, responseData.FileServerAddress2, responseData.FileServerAddress3 = version.GetFileServers()

								// 将idleServer的玩家数加1
								if sp != nil {
									sp.AddTmpPlayerCount()
								}
								gf.CheckIdleServer()
							}
						}
					}

				} else {
					responseData.Status = ERR_GAME_SERVER_DOWN
				}
			} else {
				responseData.Status = ERR_LOGIN_PLAYER_LOCK
			}

		}
	} else {
		responseData.Status = ERR_BILI_SESSION_VERIFY_FAILED
	}

	return
}

func encode(input []byte) []byte {
	compressed := snappy.Encode(nil, input)
	encrypted := xxtea.Encrypt(compressed, constants.XXTEA_MSG_KEY)
	encoded := simple.Pack(encrypted)

	return encoded
}

func BilibiliTestLogin(w http.ResponseWriter, r *http.Request) {
	// startTime := time.Now().UnixNano()
	var responseData biliLoginResp
	defer func() {
		jsonByte, _ := json.Marshal(responseData)
		w.Write(encode(jsonByte))
	}()

	defer func() {
		r.Body.Close()
	}()

	data, readError := ioutil.ReadAll(r.Body)
	if readError != nil {
		responseData = biliLoginResp{
			Status:  ERR_UNKNOW,
			Content: "read data error",
		}

		return
	}

	data = simple.Unpack(data)
	data = xxtea.Decrypt(data, constants.XXTEA_MSG_KEY)
	uncompressed, snappyErr := snappy.Decode(nil, data)
	if snappyErr != nil {
		Err("snappy error:", snappyErr.Error())
		responseData = biliLoginResp{
			Status:  ERR_BILI_SESSION_VERIFY_FORMAT,
			Content: "format err",
		}
		return
	}

	var input biliLoginReq
	err := json.Unmarshal([]byte(uncompressed), &input)
	if err != nil {
		Err(err.Error())
		responseData = biliLoginResp{
			Status:  ERR_BILI_SESSION_VERIFY_FORMAT,
			Content: "format error",
		}
		return
	}

	responseData.BiliCode = 0
	responseData.Status = NO_ERROR

	// 返回相对空闲的服务器ip
	server, server2, server3, sp, succeed := gf.GetIdleServer(input.PlayerInfo.Uid)
	if succeed {
		// todo:
		// 获取上一次玩家登陆时的游戏服务器ip
		lastServer, lastLoginTime, uid, db_addr, linfo2, newPlayer, err := users.QueryOrRegisterPlayer(input.PlayerInfo.Uid, input.PlayerInfo.Username, server)
		if err != nil {
			responseData.Status = ERR_QUERY_LAST_LOGIN_INFO
		} else {
			timeDiff := time.Now().Unix() - lastLoginTime
			if newPlayer || timeDiff < -5 || timeDiff > 5 {
				gfUsername := fmt.Sprintf("%s%d", appcfg.GetString("un_prefix", ""), uid)
				if sp == nil {
					if err = gf.LogOutInAllServers(gfUsername); err != nil {
						responseData.Status = ERR_LOGOUT_ALL
					}
				} else if lastServer != server && lastServer != "" {
					if err = gf.LogOutUser(lastServer, gfUsername); err != nil {
						//check if game server alive
						if t, err := QueryHeartbeatTime(lastServer); err != nil {
							//err
							responseData.Status = ERR_LOGOUT_USER
						} else if time.Now().Unix()-t < constants.HEARTBEAT_GAP {
							//LAST SERVER ALIVE.ERR!
							responseData.Status = ERR_LOGOUT_LAST_SERVER_ALIVE
						} else {
							//LAST SERVER DOWN.LOGIN OK.DO NOTHING
							Info("log out user failed at last server:", lastServer, "it's dead")
						}
					}
				}

				if responseData.Status == NO_ERROR {
					//DO LOGIN
					var info2 string
					if info2, err = users.LoginBilibili(input.PlayerInfo.Uid, input.PlayerInfo.Username, server, lastServer, lastLoginTime, linfo2); err != nil {
						responseData.Status = ERR_LOGIN_GF
					} else {
						responseData.Content = info2
						responseData.Address = server3
						responseData.Address2 = server
						responseData.Address3 = server2
						responseData.Uid = gfUsername
						responseData.PX_UID = uid
						responseData.PX_TOKEN = db_addr

						responseData.FileServerAddress, responseData.FileServerAddress2, responseData.FileServerAddress3 = version.GetFileServers()

						// 将idleServer的玩家数加1
						if sp != nil {
							sp.AddTmpPlayerCount()
						}
						gf.CheckIdleServer()
					}
				}
			} else {
				Debug("time diff:", timeDiff, ",lastLoginTime:", lastLoginTime, ", curTime:", time.Now().Unix())
				responseData.Status = ERR_LOGIN_TOO_FAST
			}
		}

	} else {
		responseData.Status = ERR_GAME_SERVER_DOWN
	}

	return
}
