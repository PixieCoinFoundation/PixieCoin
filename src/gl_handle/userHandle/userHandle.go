package userHandle

import (
	"encoding/json"
	"fmt"
	. "gl_db/server"
	"io/ioutil"
	"math/rand"
	"net/http"
	// "os"
	// "profile"
	// "strconv"
	"time"
)

import (
	"appcfg"
	"constants"
	"encrypt/simple"
	"encrypt/xxtea"
	"gl_db/users"
	gf "hub/loginServerHub"
	. "logger"
	"service/version"
	"snappy"
	. "types"
)

// 初始化时，从数据库读取一个全局唯一的自动注册用户名的id
func init() {
	if appcfg.GetServerType() != constants.SERVER_TYPE_GL {
		return
	}

	rand.Seed(time.Now().UnixNano())
}

type registerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 处理用户注册
type registerResp struct {
	Status   ErrorCode
	Content  string
	Username string
	Password string
}

// 玩家注册
func UserRegister(rw http.ResponseWriter, req *http.Request) {
	var responseData registerResp
	defer func() {
		jsonByte, _ := json.Marshal(responseData)
		rw.Write(encode(jsonByte))
	}()

	defer func() {
		req.Body.Close()
	}()

	data, readError := ioutil.ReadAll(req.Body)
	if readError != nil {
		responseData = registerResp{
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
		responseData = registerResp{
			Status:  ERR_UNKNOW,
			Content: "format err",
		}
		return
	}

	var input registerReq
	err := json.Unmarshal([]byte(uncompressed), &input)
	if err != nil {
		Err(err.Error())
		responseData = registerResp{
			Status:  ERR_UNKNOW,
			Content: "format error",
		}
		return
	}

	username := input.Username
	password := input.Password

	responseData = register(username, password)

	return
}

func register(username string, password string) (resp registerResp) {
	resp = registerResp{
		Status:  NO_ERROR,
		Content: "OK",
	}

	if checkEmpty(username) || checkEmpty(password) {
		resp.Status = ERR_REG_FORMAT_WRONG
		resp.Content = "invalid username or password"
		fmt.Println("username:", username, ",password:", password)
		return
	}

	err := users.Register(username, password)
	if err != nil {
		resp.Status = ERR_REG_USERNAME_EXISTS
	}

	return
}

// type editPasswordRequest struct {
// 	Username    string `json:"username"`
// 	OldPassword string `json:"oldPassword"`
// 	NewPassword string `json:"newPassword"`
// }

// // 用户修改密码
// type editPasswordResp struct {
// 	Status  ErrorCode
// 	Content string
// }

// // 玩家修改密码
// func UserEditPassword(rw http.ResponseWriter, req *http.Request) {
// 	// startTime := time.Now().UnixNano()

// 	var responseData editPasswordResp
// 	defer func() {
// 		jsonByte, _ := json.Marshal(responseData)
// 		rw.Write(encode(jsonByte))
// 	}()

// 	defer func() {
// 		req.Body.Close()
// 	}()

// 	data, readError := ioutil.ReadAll(req.Body)
// 	if readError != nil {
// 		responseData = editPasswordResp{
// 			Status:  ERR_UNKNOW,
// 			Content: "read data error",
// 		}

// 		return
// 	}

// 	data = simple.Unpack(data)
// 	data = xxtea.Decrypt(data, constants.XXTEA_MSG_KEY)
// 	uncompressed, snappyErr := snappy.Decode(nil, data)
// 	if snappyErr != nil {
// 		Err("snappy error:", snappyErr.Error())
// 		responseData = editPasswordResp{
// 			Status:  ERR_UNKNOW,
// 			Content: "format err",
// 		}
// 		return
// 	}

// 	var input editPasswordRequest
// 	err := json.Unmarshal([]byte(uncompressed), &input)
// 	if err != nil {
// 		Err(err.Error())
// 		responseData = editPasswordResp{
// 			Status:  ERR_EDIT_PWD_DATA_ERR,
// 			Content: "format error",
// 		}
// 		return
// 	}

// 	responseData = editPasswordResp{
// 		Status: NO_ERROR,
// 	}
// 	responseData.Status = editPassword(input.Username, input.OldPassword, input.NewPassword)

// 	return
// }

// func editPassword(username string, oldPassword string, newPassword string) ErrorCode {
// 	if len(newPassword) < 8 || len(newPassword) > 20 {
// 		fmt.Println("oldPassword:", oldPassword, ",newPassword:", newPassword)
// 		return ERR_EDIT_PWD_NEW_PWD_LENGTH
// 	}

// 	err := users.EditPassword(username, oldPassword, newPassword)
// 	if err != nil {
// 		fmt.Println("edit password error:", err.Error())
// 		return ERR_EDIT_PWD
// 	} else {
// 		return NO_ERROR
// 	}
// }

type loginReq struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DeviceToken string `json:"deviceToken"`
}

// 处理用户登陆
type loginResp struct {
	Status             ErrorCode `json:"status"`
	Content            string    `json:"content"`
	Address            string    `json:"address"`            // 逻辑服务器地址
	Address2           string    `json:"address2"`           // 逻辑服务器地址
	Address3           string    `json:"address3"`           // 逻辑服务器地址
	FileServerAddress  string    `json:"fileServerAddress"`  // 文件服务器地址
	FileServerAddress2 string    `json:"fileServerAddress2"` // 文件服务器地址
	FileServerAddress3 string    `json:"fileServerAddress3"` // 文件服务器地址

	PX_TOKEN string
	Uid      string

	PX_UID int64
}

// 玩家登陆
func UserLogin(w http.ResponseWriter, r *http.Request) {
	var responseData loginResp
	defer func() {
		jsonByte, _ := json.Marshal(responseData)
		w.Write(encode(jsonByte))
	}()

	defer func() {
		r.Body.Close()
	}()

	data, readError := ioutil.ReadAll(r.Body)
	if readError != nil {
		Err(readError)
		responseData = loginResp{
			Status:  ERR_UNKNOW,
			Content: "read data error",
		}

		return
	}

	data = simple.Unpack(data)
	data = xxtea.Decrypt(data, constants.XXTEA_MSG_KEY)
	uncompressed, snappyErr := snappy.Decode(nil, data)
	if snappyErr != nil {
		Err("snappy error:", snappyErr)
		responseData = loginResp{
			Status:  ERR_UNKNOW,
			Content: "format err",
		}
		return
	}

	var input loginReq
	err := json.Unmarshal([]byte(uncompressed), &input)
	if err != nil {
		Err(err)
		responseData = loginResp{
			Status:  ERR_UNKNOW,
			Content: "format error",
		}
		return
	}

	username := input.Username
	password := input.Password

	var server, server2, server3, lastServer, dbAddr string
	uid := int64(-1)
	var lastLoginTime int64

	// 获取空闲server,如果传入了server，就登陆指定server，用于测试
	var succeed bool
	var sp *gf.GFServer
	server, server2, server3, sp, succeed = gf.GetIdleServer(username)
	if succeed {
		responseData, lastServer, uid, dbAddr, lastLoginTime = login(username, password, server)
		gfUsername := fmt.Sprintf("%s%d", appcfg.GetString("un_prefix", ""), uid)
		timeDiff := time.Now().Unix() - lastLoginTime
		if responseData.Status == 0 && (timeDiff < -5 || timeDiff > 5) {
			// 如果数据库校验成功，则向对应的游戏服务器请求访问token
			resp := fmt.Sprintf("%s_%d", username, time.Now().UnixNano())

			// 通知之前登陆的server清掉玩家数据
			if lastServer != server && lastServer != "" {
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
			}

			Info("player:", username, gfUsername, "login success:", server)
			responseData.Status = NO_ERROR
			responseData.Content = resp
			responseData.Address = server3
			responseData.Address2 = server
			responseData.Address3 = server2
			responseData.Uid = gfUsername
			responseData.PX_UID = uid
			responseData.PX_TOKEN = dbAddr

			responseData.FileServerAddress, responseData.FileServerAddress2, responseData.FileServerAddress3 = version.GetFileServers()

			// 将idleServer的玩家数加1
			if sp != nil {
				sp.AddTmpPlayerCount()
			}
			gf.CheckIdleServer()
		} else if timeDiff >= -5 && timeDiff <= 5 {
			responseData.Status = ERR_LOGIN_TOO_FAST
			responseData.Content = "login too fast"
		}
	} else {
		responseData.Status = ERR_GAME_SERVER_DOWN
		responseData.Content = "server down"
	}

	return
}

func login(username string, password string, server string) (resp loginResp, lastLoginServer string, uid int64, dbAddr string, lastLoginTime int64) {
	resp = loginResp{0, "OK", "", "", "", "", "", "", "", "", 0}

	if checkEmpty(username) || checkEmpty(password) {
		resp.Status = ERR_LOGIN_FORMAT_WRONG
		resp.Content = "invalid username or password"
		return
	}

	var err error
	lastLoginServer, lastLoginTime, uid, dbAddr, err = users.Login(username, password, server)
	if err != nil {
		if err.Error() != constants.SQL_NO_DATA_ERR_MSG {
			Err("user login failed:", username, err)
		} else {
			Info("no such player:", username)
		}
		resp.Status = ERR_LOGIN_FIELD_WRONG
		resp.Content = err.Error()
	}

	return
}

// 检查是否字符串为空
func checkEmpty(content string) bool {
	if content == "" {
		return true
	} else {
		return false
	}
}

func encode(input []byte) []byte {
	compressed := snappy.Encode(nil, input)
	encrypted := xxtea.Encrypt(compressed, constants.XXTEA_MSG_KEY)
	encoded := simple.Pack(encrypted)

	return encoded
}
