package kor

import (
	"appcfg"
	"bytes"
	"constants"
	"db/kor"
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "logger"
	// "net"
	"net/http"
	"os"
	"strings"
	"time"
	"tools"
)

const (
	KAKAO_UNREGISTER_URL = "https://openapi-zinny3.game.kakao.com:10443/service/v3/player/remove"
)

//Ctblist 玩家列表
// var Cbtlist = map[string]string{}

var reward1Map map[string]int
var reward2Map map[string]int
var reward3Map map[string]int

//CBTPlayer 预约玩家信息
type CBTPlayer struct {
	KaKaoID    string
	RewardType string
}

func init() {
	if appcfg.GetServerType() != "" && appcfg.GetServerType() != constants.SERVER_TYPE_SIMPLE_TEST {
		return
	}

	// if appcfg.GetServerType() == constants.SERVER_TYPE_SIMPLE_TEST {
	// 	Info(1)
	// 	loadFile("scripts/kor/cbt_kakao_id.txt")
	// 	Info(2)
	// 	// Info(reward1Map)
	// 	// Info(reward2Map)
	// 	// Info(reward3Map)
	// 	panic("end")
	// }

	if appcfg.GetServerType() == "" && appcfg.GetLanguage() == constants.KOR_LANGUAGE && appcfg.GetBool("main_server", false) {
		go unregisterLoop()
	}

	if appcfg.GetServerType() == "" && appcfg.GetLanguage() == constants.KOR_LANGUAGE && appcfg.GetBool("game_gift_is_available", false) {
		loadFile("scripts/kor/cbt_kakao_id.txt")
	}
}

func loadFile(filename string) {
	reward1Map = make(map[string]int)
	reward2Map = make(map[string]int)
	reward3Map = make(map[string]int)

	if len(filename) < 1 {
		panic(constants.ErrorNilData)
	}
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	ldata := strings.Split(string(data), "\n")
	if len(ldata) < 1 {
		panic(constants.ErrorNilData)
	} else {
		Info("kor init player reward size:", len(ldata))
	}
	for _, v := range ldata {
		detildata := strings.Split(v, ",")
		if len(detildata) != 2 {
			panic(constants.ErrorNilData)
		}

		rt := strings.TrimSpace(detildata[1])
		kakaoID := strings.TrimSpace(detildata[0])
		if rt == "1" {
			reward1Map[kakaoID] = 1
		} else if rt == "2" {
			reward2Map[kakaoID] = 1
		} else if rt == "3" {
			reward3Map[kakaoID] = 1
		} else {
			panic("unknown reward type:" + rt + ":" + v)
		}
	}
}

func GetKorInitPlayerReward(kakaoID string) (reward1 bool, reward2 bool, reward3 bool) {
	return reward1Map[kakaoID] > 0, reward2Map[kakaoID] > 0, reward3Map[kakaoID] > 0
}

func DelUserUnregister(username string) (err error) {
	return kor.DelUserUnregister(username)
}

func unregisterLoop() {
	for {
		now := time.Now().Unix()
		if us, err := kor.ListUnregister(); err == nil {
			for _, up := range us {
				if up.DelTime <= now {
					//unregister player
					if unregister(up.ThirdUsername) {
						upb, _ := json.Marshal(up)
						GMLog(constants.C1_SYSTEM, constants.C2_KOR_UNREGISTER, constants.C3_DEFAULT, up.Username, string(upb))
						kor.DelUnregister(up.ID)
					}
				}
			}
		}

		time.Sleep(600 * time.Second)
	}
}

type UnregisterKakaoReq struct {
	IdpUnlink bool `json:"idpUnlink"`
}

func unregister(tu string) bool {
	req := UnregisterKakaoReq{
		IdpUnlink: true,
	}
	rb, _ := json.Marshal(req)
	Info("kor unregister kakao:", string(rb), tu)

	toSend := bytes.NewBuffer(rb)
	if reqest, err := http.NewRequest("POST", KAKAO_UNREGISTER_URL, toSend); err != nil {
		Err(err)
		return false
	} else {
		reqest.Header.Set("appId", constants.KAKAO_APP_ID)
		reqest.Header.Set("appSecret", constants.KAKAP_NATIVE_APP_KEY)
		reqest.Header.Set("playerId", tu)
		reqest.Header.Set("Content-Type", "application/json;charset=UTF-8")
		reqest.Header.Set("Authorization", fmt.Sprintf("KakaoAK %s", constants.KAKAO_ADMIN_KEY))

		// client := &http.Client{
		// 	Transport: &http.Transport{
		// 		Dial: func(netw, addr string) (net.Conn, error) {
		// 			deadline := time.Now().Add(constants.DEFAULT_HTTP_TIMEOUT * time.Second)
		// 			c, err := net.DialTimeout(netw, addr, time.Second*constants.DEFAULT_HTTP_TIMEOUT)
		// 			if err != nil {
		// 				return nil, err
		// 			}
		// 			c.SetDeadline(deadline)
		// 			return c, nil
		// 		},
		// 		DisableKeepAlives: true,
		// 	},
		// }
		response, err := tools.OtakuHttpClient.Do(reqest)
		if err != nil {
			Err(err)
			return false
		}
		defer response.Body.Close()

		if body, err := ioutil.ReadAll(response.Body); err != nil {
			Err(err)
			return false
		} else {
			Info("unregister kakao result:", string(body))
		}

		if response.StatusCode == 200 {
			return true
		}
	}

	return false
}
