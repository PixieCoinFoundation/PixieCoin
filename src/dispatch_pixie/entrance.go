package dispatch_pixie

import (
	"appcfg"
	"constants"
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "logger"
	"net/http"
	"time"
)

import (
	"encrypt/simple"
	"encrypt/xxtea"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	"player_mem"
	"service/params"
	"shutdown"
	"snappy"
	. "types"
)

type inputData struct {
	AccessToken string
	Username    string
	Method      string
	Params      string
	ReqID       string
	PressTest   bool
}

type outputData struct {
	Resp    PixieRespInfo //客户端通常应先读取Resp再决定是否读取content
	Content string
}

type banData struct {
	StartTime int64
	EndTime   int64
	Reason    string
}

type maintainData struct {
	Error string
	URL   string
}

func TestApiEntrance(w http.ResponseWriter, r *http.Request) {
	playerNum, _ := player_mem.GetGameServerStatus(false)
	w.Write([]byte(fmt.Sprint(playerNum)))
}

func NotSupportApiHandle(w http.ResponseWriter, r *http.Request) {
	resp := outputData{
		Resp: PIXIE_ERR_API_NOT_SUPPORT,
	}

	w.Write(EncryptResponse(resp))
}

func PixieApiGlobalEntrance(w http.ResponseWriter, r *http.Request) {
	shutdown.SimpleAddOneRequest()
	defer shutdown.SimpleDoneOneRequest()

	startTime := time.Now().UnixNano()

	//always response
	var resp outputData
	defer func() {
		if x := recover(); x != nil {
			Err("caught panic in HttpDispatch", x)

			resp = outputData{
				Resp: PIXIE_ERR_UNKNOWN,
			}
		}

		w.Write(EncryptResponse(resp))
	}()

	//check if server is shutting down
	if shutdown.IsShutDown() {
		resp = outputData{
			Resp: PIXIE_ERR_SHUTDOWN,
		}

		return
	}

	//read from client
	data, readError := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if readError != nil || len(data) == 0 {
		resp = outputData{
			Resp: PIXIE_ERR_WRONG_PARAMS,
		}
		go statistic(startTime, time.Now().UnixNano(), "UNKNOWN_METHOD_READ", "UNKNOWN_REQ_ID_READ", "UNKNOWN_USERNAME_READ", "UNKNOWN_PARAMS_READ", resp.Resp.RespCodeStr, "")
		return
	}

	data = simple.Unpack(data)
	data = xxtea.Decrypt(data, constants.XXTEA_MSG_KEY)
	uncompressed, snappyErr := snappy.Decode(nil, data)
	if snappyErr != nil {
		Err("snappy error:", snappyErr)
		resp = outputData{
			Resp: PIXIE_ERR_DECRYPT,
		}
		go statistic(startTime, time.Now().UnixNano(), "UNKNOWN_METHOD_DECODE", "UNKNOWN_REQ_ID_DECODE", "UNKNOWN_USERNAME_DECODE", "UNKNOWN_PARAMS_DECODE", resp.Resp.RespCodeStr, "")
		return
	}

	var input inputData
	err := json.Unmarshal([]byte(uncompressed), &input)
	if err != nil {
		// 注意，有可能是一些特殊字符不支持json decode
		Err("json unmarshal in api entrance", err, string(uncompressed))
		resp = outputData{
			Resp: PIXIE_ERR_WRONG_JSON,
		}
		go statistic(startTime, time.Now().UnixNano(), "UNKNOWN_METHOD_JSON", "UNKNOWN_REQ_ID_JSON", "UNKNOWN_USERNAME_JSON", "UNKNOWN_PARAMS_JSO", resp.Resp.RespCodeStr, "")
		return
	}

	defer func() {
		go statistic(startTime, time.Now().UnixNano(), input.Method, input.ReqID, input.Username, input.Params, resp.Resp.RespCodeStr, resp.Content)
	}()

	if !input.PressTest {
		//check if server maintaining..
		ms, mi, murl := params.IsMaintainStatus()
		if ms && !params.IsInWhiteList(input.Username) {
			md := maintainData{
				Error: mi,
				URL:   murl,
			}
			mdb, _ := json.Marshal(md)
			resp = outputData{
				Resp:    PIXIE_ERR_MAINTAIN,
				Content: string(mdb),
			}

			return
		}
	}

	player := player_mem.GetPlayerByUsernameAndToken(input.Username, input.AccessToken)

	if input.Method != HTTP_PLAYER_LOGIN {
		if player == nil {
			aet := AbnormalExtra{
				ServerAddress: appcfg.GetAddress(),
				Info:          input.Method,
			}
			aetb, _ := json.Marshal(aet)
			GMLog(constants.C1_SYSTEM, constants.C2_ANNORMAL, constants.C3_DEFAULT, input.Username, string(aetb))
			resp = outputData{
				Resp: PIXIE_ERR_ACCESS_TOKEN,
			}
			return
		} else if banned, st, et, reason := player.IsBanned(); banned {
			bd := banData{
				StartTime: st,
				EndTime:   et,
				Reason:    reason,
			}
			bdb, _ := json.Marshal(bd)
			resp = outputData{
				Resp:    PIXIE_ERR_PLAYER_BAN,
				Content: string(bdb),
			}
			return
		} else if !input.PressTest && params.ChannelForbided(player.ThirdChannel, player.UID) {
			resp = outputData{
				Resp: PIXIE_ERR_CHANNEL_FORBID,
			}

			return
		} else if blocked, reason := params.Blocked(player.Username); blocked {
			bd := banData{
				StartTime: 0,
				EndTime:   2051193600,
				Reason:    reason,
			}
			bdb, _ := json.Marshal(bd)
			resp = outputData{
				Resp:    PIXIE_ERR_PLAYER_BLOCK,
				Content: string(bdb),
			}
			return
		} else {
			player.Lock()
			defer player.Unlock()
			player.RefreshHeartbeatTime()
		}

		if handle, ok := PixiePlayerRequestHandles[input.Method]; ok {
			content, r, err := handle(player, input.Params)
			if err != nil {
				resp.Resp = PIXIE_ERR_API_UNMARSHAL
			} else {
				resp.Resp = r
				resp.Content = content
			}
		} else {
			resp.Resp = PIXIE_ERR_API_NOT_SUPPORT
		}
	} else {
		content, r, err := LoginHandle(player, input.Params, r.RemoteAddr)
		if err != nil {
			resp.Resp = PIXIE_ERR_API_UNMARSHAL
		} else {
			resp.Resp = r
			resp.Content = content
		}
	}
}

type ProfileExtra struct {
	Output        string
	CostTime      float64
	Data          string
	Method        string
	ReqID         string
	RespStr       string
	ServerAddress string
}

func statistic(begin int64, end int64, methodName string, reqID string, username string, data string, respStr string, respContent string) {
	gap := float64(end-begin) / float64(time.Millisecond)
	length := len(data)

	if length > constants.DATALOG_LEN {
		length = constants.DATALOG_LEN
	}
	if len(respContent) > constants.DATALOG_LEN {
		respContent = respContent[0:constants.DATALOG_LEN]
	}
	p := ProfileExtra{
		Output:        respContent,
		CostTime:      gap,
		Data:          data[0:length],
		Method:        methodName,
		ReqID:         reqID,
		RespStr:       respStr,
		ServerAddress: appcfg.GetAddress(),
	}
	if d, err := json.Marshal(p); err != nil {
		Err(err)
	} else {
		Profile(constants.C1_SYSTEM, constants.C2_PROFILE, constants.C3_API_COST, username, string(d))
	}
}
