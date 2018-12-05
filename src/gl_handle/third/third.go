package third

import (
	"appcfg"
	"constants"
	// "crypto/x509"
	"encoding/json"
	"fmt"
	. "gl_db/server"
	"gl_db/users"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

import (
	"encrypt/simple"
	"encrypt/xxtea"
	gf "hub/loginServerHub"
	. "logger"
	"service/version"
	"snappy"
	. "types"
	. "zk_manager"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	//100026925 1498025199977 28785605
	// checkHuaweiToken("100026925", "1498025199977", "28785605", "Yta+PSQ4S7tfT3enTU8nIgfjBf8ym8qXzCLOHviKF0LR8c3NWoQt1sIR0zB5REU/sn91plS/jrhEUOt16OE2RNQ7AVS3uwGynATkjmXdH3cS8HD3GtV1v8wjuODi/95YCUKFISD5wH9uTHgO0qwL+t1YW5wSSnYo5mS89LrSu7fzhQaMxbBPbKwHc/PV8FtVE4hDNB7CLz9Jol9WtQmng4m34akBc/WEstfzTImKcDiPNDqmUYr3jsfzLOOEDAvGPYSRDIFmFgADnOepGrbjxfyzsC1NQaiwc8t+ysbx/CKiZgEXR27z6jyQ0NjWcQrgKmI4a3Dtg2wdUgStcp20SA==")
}

type thirdLoginReq struct {
	Platform      string
	ThirdUsername string
	ThirdToken    string
	ThirdExtra1   string
	ThirdExtra2   string
	ThirdExtra3   string

	DelPlayer bool //韩国专用参数，删除原有账号映射，新增账号
}

// 处理用户登陆
type thirdLoginResp struct {
	Status ErrorCode

	Address            string // 逻辑服务器地址
	Address2           string // 逻辑服务器地址
	Address3           string // 逻辑服务器地址
	FileServerAddress  string // 文件服务器地址
	FileServerAddress2 string // 文件服务器地址
	FileServerAddress3 string // 文件服务器地址

	PX_Username string //game username
	PX_Token    string //game token
	PX_UID      int64  //game uid
}

// 玩家在客户端调用sdk登陆之后，将获取到的玩家信息发送到服务端
// 服务端向bilibili的SDK服务器验证accessToken，并将验证结果返回客户端
func ThirdLogin(w http.ResponseWriter, r *http.Request) {
	// startTime := time.Now().UnixNano()
	var responseData thirdLoginResp
	defer func() {
		jsonByte, _ := json.Marshal(responseData)
		w.Write(encode(jsonByte))
	}()

	defer func() {
		r.Body.Close()
	}()

	data, readError := ioutil.ReadAll(r.Body)
	if readError != nil {
		responseData.Status = ERR_UNKNOW

		return
	}

	data = simple.Unpack(data)
	data = xxtea.Decrypt(data, constants.XXTEA_MSG_KEY)
	uncompressed, snappyErr := snappy.Decode(nil, data)
	if snappyErr != nil {
		Err("snappy error:", snappyErr.Error())
		responseData = thirdLoginResp{
			Status: ERR_UNKNOW,
		}
		return
	}

	var input thirdLoginReq
	err := json.Unmarshal([]byte(uncompressed), &input)
	if err != nil {
		Err(err.Error())
		responseData = thirdLoginResp{
			Status: ERR_UNKNOW,
		}
		return
	}

	if appcfg.GetLanguage() == constants.KOR_LANGUAGE && input.Platform != constants.THIRD_CHANNEL_KOR_IOS && input.Platform != constants.THIRD_CHANNEL_KOR_GOOGLE_PLAY && input.Platform != constants.THIRD_CHANNEL_KOR_ONE_STORE {
		responseData = thirdLoginResp{
			Status: ERR_CHANNEL_NOT_SUPPORT,
		}
		return
	}

	//verify third user
	Info("third login:", input.Platform, input.ThirdUsername, input.ThirdToken)
	tn := input.ThirdUsername
	switch input.Platform {
	case constants.THIRD_CHANNEL_HUAWEI:
		if !checkHuaweiToken(input.ThirdExtra1, input.ThirdExtra2, input.ThirdUsername, input.ThirdToken) {
			responseData.Status = ERR_HUAWEI_CHECK_FAIL
			return
		}
	case constants.THIRD_CHANNEL_VIVO:
		if !checkVIVO(input.ThirdToken) {
			responseData.Status = ERR_VIVO_CHECK_FAIL
			return
		}
	case constants.THIRD_CHANNEL_UC:
		var ok bool
		if ok, tn = checkUC(input.ThirdToken); !ok {
			responseData.Status = ERR_UC_CHECK_FAIL
			return
		}
	case constants.THIRD_CHANNEL_4399:
		if !check4399(input.ThirdUsername, input.ThirdToken) {
			responseData.Status = ERR_4399_CHECK_FAIL
			return
		}
	case constants.THIRD_CHANNEL_MEITU:
		if !checkMeitu(input.ThirdUsername, input.ThirdToken) {
			responseData.Status = ERR_MEITU_CHECK_FAIL
			return
		}
	case constants.THIRD_CHANNEL_PAPA:
		if !checkPapa(input.ThirdUsername, input.ThirdToken) {
			responseData.Status = ERR_PAPA_CHECK_FAIL
			return
		}
	case constants.THIRD_CHANNEL_KUAIKAN:
		if !checkKuaikan(input.ThirdUsername, input.ThirdToken) {
			responseData.Status = ERR_KUAIKAN_CHECK_FAIL
			return
		}
	case constants.THIRD_CHANNEL_BUKA:
		if !checkBuka(input.ThirdUsername, input.ThirdToken) {
			responseData.Status = ERR_BUKA_CHECK_FAIL
			return
		}
	case constants.THIRD_CHANNEL_XIAOMI:
		if !checkXiaomi(input.ThirdUsername, input.ThirdToken) {
			responseData.Status = ERR_XIAOMI_CHECK_FAIL
			return
		}
	case constants.THIRD_CHANNEL_SINA:
		if !checkSina(input.ThirdUsername, input.ThirdToken) {
			responseData.Status = ERR_SINA_CHECK_FAIL
			return
		}
	case constants.THIRD_CHANNEL_360:
		var ok bool
		if ok, tn = check360(input.ThirdToken); !ok {
			responseData.Status = ERR_360_CHECK_FAIL
			return
		}
	case constants.THIRD_CHANNEL_OPPO:
		if !checkOPPO(input.ThirdUsername, input.ThirdToken) {
			responseData.Status = ERR_OPPO_CHECK_FAIL
			return
			// canReg = false
		}
	case constants.THIRD_CHANNEL_BAIDU:
		if !checkBaidu(input.ThirdUsername, input.ThirdToken) {
			responseData.Status = ERR_BAIDU_CHECK_FAIL
			return
		}
	case constants.THIRD_CHANNEL_MHR:
		//no check!
	case constants.THIRD_CHANNEL_KOR_IOS:
		fallthrough
	case constants.THIRD_CHANNEL_KOR_GOOGLE_PLAY:
		fallthrough
	case constants.THIRD_CHANNEL_KOR_ONE_STORE:
		if !checkKakao(input.ThirdUsername, input.ThirdToken) {
			responseData.Status = ERR_KAKAO_CHECK_FAIL
			return
		}
	default:
		responseData.Status = ERR_THIRD_CHANNEL_EMPTY
		return
	}

	ru := tn

	if input.Platform != constants.THIRD_CHANNEL_KOR_IOS && input.Platform != constants.THIRD_CHANNEL_KOR_ONE_STORE && input.Platform != constants.THIRD_CHANNEL_KOR_GOOGLE_PLAY {
		ru = "T" + input.Platform + constants.THIRD_USERNAME_GAP + ru
	}

	now := time.Now()

	if ok, path := LockPlayerLogin(ru); ok {
		defer Unlock(path)

		if input.DelPlayer && (input.Platform == constants.THIRD_CHANNEL_KOR_IOS || input.Platform == constants.THIRD_CHANNEL_KOR_ONE_STORE || input.Platform == constants.THIRD_CHANNEL_KOR_GOOGLE_PLAY) {
			if err := users.DeletePlayer(ru); err != nil {
				responseData.Status = ERR_DELETE_PLAYER
				return
			} else {
				Info("delete kor player succeed.gen new player:", ru)
			}
		}

		// 返回相对空闲的服务器ip
		server, server2, server3, sp, succeed := gf.GetIdleServer(ru)
		if succeed {
			// todo:
			// 获取上一次玩家登陆时的游戏服务器ip
			lastServer, lastLoginTime, uid, db_addr, info2, newPlayer, err := users.QueryOrRegisterPlayer(ru, ru, server)
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
					if _, err = users.LoginBilibili(ru, ru, server, lastServer, lastLoginTime, info2); err != nil {
						responseData.Status = ERR_LOGIN_GF
					} else {
						Info("platform:", input.Platform, "player:", ru, gfUsername, "login success:", server)
						responseData.Address = server3
						responseData.Address2 = server
						responseData.Address3 = server2
						responseData.PX_Username = gfUsername
						responseData.PX_UID = uid

						responseData.PX_Token = db_addr

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

	return
}

func encode(input []byte) []byte {
	compressed := snappy.Encode(nil, input)
	encrypted := xxtea.Encrypt(compressed, constants.XXTEA_MSG_KEY)
	encoded := simple.Pack(encrypted)

	return encoded
}
