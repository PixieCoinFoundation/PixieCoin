package third

import (
	"appcfg"
	"common"
	"constants"
	// "crypto/sha256"
	// "encoding/base64"
	// "encoding/hex"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "logger"
	"math/rand"
	// "net"
	"net/http"
	"net/url"
	"strings"
	"time"
	"tools"
)

var HUAWEI_LOGIN_KEY = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmKLBMs2vXosqSR2rojMzioTRVt8oc1ox2uKjyZt6bHUK0u+OpantyFYwF3w1d0U3mCF6rGUnEADzXiX/2/RgLQDEXRD22er31ep3yevtL/r0qcO8GMDzy3RJexdLB6z20voNM551yhKhB18qyFesiPhcPKBQM5dnAOdZLSaLYHzQkQKANy9fYFJlLDo11I3AxefCBuoG+g7ilti5qgpbkm6rK2lLGWOeJMrF+Hu+cxd9H2y3cXWXxkwWM1OZZTgTq3Frlsv1fgkrByJotDpRe8SwkiVuRycR0AHsFfIsuZCFwZML16EGnHqm2jLJXMKIBgkZTzL8Z+201RmOheV4AQIDAQAB"

const (
	VIVO_CHECK_URL    = "https://usrsys.vivo.com.cn/sdk/user/auth.do?authtoken=%s&from=yjdyc"
	KUAIKAN_CHECK_URL = "http://api.kkmh.com/v1/game/oauth/check_open_id?app_id=%s&open_id=%s&access_token=%s&sign=%s"
	CHECK_URL_4399    = "http://m.4399api.com/openapi/oauth-check.html?uid=%s&state=%s"
	CHECK_BUKA_URL    = "http://i.tongbulv.com/api/check?supplier_id=%s&t=%s&uid=%s&access_token=%s&sign=%s"
	CHECK_MEITU_URL   = "http://www.91wan.com/api/mobile/sdk_oauth.php?appid=%s&uid=%s&state=%s&flag=%s"
	CHECK_PAPA_URL    = "https://sdkapi.papa91.com/auth/check_token"

	CHECK_UC_URL     = "http://sdk.9game.cn/cp/account.verifySession"
	CHECK_XIAOMI_URL = "http://mis.migc.xiaomi.com/api/biz/service/verifySession.do?appId=%s&session=%s&uid=%s&signature=%s"
	CHECK_SINA_URL   = "http://m.game.weibo.cn/ddsdk/distsvc/1/user/serververify_v2?suid=%s&appkey=%s&token=%s&signature=%s"
	CHECK_360_URL    = "https://openapi.360.cn/user/me.json?access_token=%s"
	CHECK_OPPO_URL   = "http://i.open.game.oppomobile.com/gameopen/user/fileIdInfo?fileId=%s&token=%s"

	CHECK_KAKAO_URL = "https://openapi-zinny3.game.kakao.com:10443/service/v3/zat/validate"
	CHECK_BAIDU_URL = "http://querysdkapi.baidu.com/query/cploginstatequery?"
)

var javaPath string

func init() {
	javaPath = appcfg.GetString("java_path", "java")

	// checkKakao("1", "2")
	// panic("end")
}

type CheckKakaoReq struct {
	Zat string `json:"zat"`
}

// type CheckKakaoResp struct {
// 	Lang    string `json:"lang"`
// 	Market  string `json:"market"`
// 	Country string `json:"country"`
// }

func checkKakao(tu, tt string) bool {
	req := CheckKakaoReq{
		Zat: tt,
	}
	rb, _ := json.Marshal(req)
	Info("check kakao req body:", string(rb))

	toSend := bytes.NewBuffer(rb)
	if reqest, err := http.NewRequest("POST", CHECK_KAKAO_URL, toSend); err != nil {
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
			Info("check kakao result:", string(body))
		}

		if response.StatusCode == 200 {
			return true
		}
	}
	return false
}

type RespBaidu struct {
	ResultCode int
	ResultMsg  string
}

func checkBaidu(tu, tt string) bool {
	sc := fmt.Sprintf("%s%s%s", constants.BAIDU_APP_ID, tt, constants.BAIDU_SECRET_KEY)
	sign, _ := common.GenMD5S(sc)

	data := url.Values{}
	data.Set("AppID", constants.BAIDU_APP_ID)
	data.Set("AccessToken", tt)
	data.Set("Sign", sign)

	if resp, err := http.PostForm(CHECK_BAIDU_URL, data); err != nil {
		Err(err)
		return false
	} else {
		defer resp.Body.Close()
		if res, err := ioutil.ReadAll(resp.Body); err != nil {
			Err(err)
			return false
		} else {
			Info("baidu check login result:", string(res))

			var resp RespBaidu
			if err := json.Unmarshal(res, &resp); err != nil {
				Err(err)
				return false
			}

			if resp.ResultCode == 1 {
				return true
			}
		}
	}

	return false
}

type OppoResp struct {
	ResultCode string
	SSOID      string
}

func checkOPPO(tu, tt string) bool {
	// client := &http.Client{
	// 	Transport: &http.Transport{
	// 		Dial: func(netw, addr string) (net.Conn, error) {
	// 			deadline := time.Now().Add(30 * time.Second)
	// 			c, err := net.DialTimeout(netw, addr, time.Second*30)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			c.SetDeadline(deadline)
	// 			return c, nil
	// 		},
	// 		DisableKeepAlives: true,
	// 	},
	// }
	if req, err := http.NewRequest("GET", fmt.Sprintf(CHECK_OPPO_URL, tu, tt), nil); err != nil {
		Err(err)
		return false
	} else {
		sc := fmt.Sprintf("oauthConsumerKey=%s&oauthToken=%s&oauthSignatureMethod=HMAC-SHA1&oauthTimestamp=%d&oauthNonce=%d&oauthVersion=1.0&", constants.OPPO_APP_KEY, tt, time.Now().Unix(), rand.Intn(100))
		Info(sc)
		sign := common.HmacSha1EncryptBase64(sc, constants.OPPO_APP_SECRET+"&")
		sign = url.QueryEscape(sign)
		Info("check oppo sign:", sign)
		req.Header.Set("param", sc)
		req.Header.Set("oauthSignature", sign)
		if res, err := tools.OtakuHttpClient.Do(req); err != nil {
			Err(err)
			return false
		} else {
			defer res.Body.Close()
			if resb, err := ioutil.ReadAll(res.Body); err == nil {
				Info("check oppo login result:", string(resb))

				var resp OppoResp
				if err := json.Unmarshal(resb, &resp); err != nil {
					Err(err)
					return false
				}

				if resp.ResultCode == "200" && resp.SSOID == tu {
					return true
				}
			} else {
				Err(err)
				return false
			}
		}

	}
	return false
}

type Resp360 struct {
	ID     string
	Name   string
	Avatar string
}

func check360(tt string) (bool, string) {
	if res, err := common.GetUrl(fmt.Sprintf(CHECK_360_URL, tt)); err != nil {
		return false, ""
	} else {
		ress := string(res)
		Info("check 360 login result:", ress)

		var resp Resp360
		if err := json.Unmarshal(res, &resp); err != nil {
			Err(err)
			return false, ""
		}

		if resp.ID != "" {
			return true, resp.ID
		}
	}
	return false, ""
}

type RespSina struct {
	SUID     string `json:"suid"`
	UserID   string `json:"user_id"`
	Username string `json:"user_name"`
	Token    string `json:"token"`
	Avatar   string `json:"avatar"`
}

func checkSina(tu, tt string) bool {
	sc := fmt.Sprintf("appkey|%s|suid|%s|token|%s|%s", constants.SINA_APP_KEY, tu, tt, constants.SINA_SIGNATURE_KEY)
	sign := common.Sha1Encrypt(sc)

	if res, err := common.GetUrl(fmt.Sprintf(CHECK_SINA_URL, tu, constants.SINA_APP_KEY, tt, sign)); err != nil {
		return false
	} else {
		ress := string(res)
		Info("check sina login result:", ress)

		var resp RespSina
		if err := json.Unmarshal(res, &resp); err != nil {
			Err(err)
			return false
		}

		if resp.SUID == tu {
			return true
		}
	}
	return false
}

type RespXiaomi struct {
	ErrorCode int `json:"errcode"`
}

func checkXiaomi(tu, tt string) bool {
	sc := fmt.Sprintf("appId=%s&session=%s&uid=%s", constants.XIAOMI_APP_ID, tt, tu)
	sign := common.HmacSha1Encrypt(sc, constants.XIAOMI_APP_SECRET)

	if res, err := common.GetUrl(fmt.Sprintf(CHECK_XIAOMI_URL, constants.XIAOMI_APP_ID, tt, tu, sign)); err != nil {
		return false
	} else {
		ress := string(res)
		Info("check xiaomi login result:", ress)

		var resp RespXiaomi
		if err := json.Unmarshal(res, &resp); err != nil {
			Err(err)
			return false
		}

		if resp.ErrorCode == 200 {
			return true
		}
	}
	return false
}

type CheckUCReq struct {
	ID   int64             `json:"id"`
	Data map[string]string `json:"data"`
	Game map[string]int    `json:"game"`
	Sign string            `json:"sign"`
}

type CheckUCResp struct {
	ID    int64            `json:"id"`
	State CheckUCRespState `json:"state"`
	Data  CheckUCRespData  `json:"data"`
}

type CheckUCRespState struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CheckUCRespData struct {
	AccountID string `json:"accountId"`
}

func checkUC(tt string) (bool, string) {
	sc := fmt.Sprintf("sid=%s%s", tt, constants.UC_API_KEY)
	Info(sc)
	if sign, err := common.GenMD5S(sc); err != nil {
		return false, ""
	} else {
		req := CheckUCReq{
			ID:   time.Now().Unix(),
			Data: map[string]string{"sid": tt},
			Game: map[string]int{"gameId": constants.UC_GAME_ID},
			Sign: sign,
		}

		rb, _ := json.Marshal(req)
		Info("req body:", string(rb))
		toSend := bytes.NewBuffer(rb)

		if resp, err := http.Post(CHECK_UC_URL, "application/json", toSend); err != nil {
			Err(err)
			return false, ""
		} else {
			defer resp.Body.Close()
			var ucresp CheckUCResp
			if body, err := ioutil.ReadAll(resp.Body); err != nil {
				Err(err)
				return false, ""
			} else {
				Info("check uc result:", string(body))

				if err := json.Unmarshal(body, &ucresp); err != nil {
					Err(err)
					return false, ""
				}

				if ucresp.State.Code == 1 && ucresp.Data.AccountID != "" {
					return true, ucresp.Data.AccountID
				}
			}
		}
	}

	return false, ""
}

type CheckPapaResp struct {
	Data CheckPapaRespData `json:"data"`
}

type CheckPapaRespData struct {
	IsSuccess bool `json:"is_success"`
}

func checkPapa(tu, tt string) bool {
	sc := fmt.Sprintf("%s%sapp_key=%s&token=%s&uid=%s", constants.PAPA_APP_KEY, constants.PAPA_SECRET_KEY, constants.PAPA_APP_KEY, tt, tu)
	if sign, err := common.GenMD5S(sc); err != nil {
		return false
	} else {
		data := url.Values{}
		data.Set("app_key", constants.PAPA_APP_KEY)
		data.Set("token", tt)
		data.Set("uid", tu)
		data.Set("sign", sign)

		if resp, err := http.PostForm(CHECK_PAPA_URL, data); err != nil {
			Err(err)
			return false
		} else {
			defer resp.Body.Close()

			if res, err := ioutil.ReadAll(resp.Body); err != nil {
				Err(err)
				return false
			} else {
				ress := string(res)
				Info("check papa login result:", ress)

				var resp CheckPapaResp
				if err := json.Unmarshal(res, &resp); err != nil {
					Err(err)
					return false
				}

				if resp.Data.IsSuccess {
					return true
				}
			}
		}
	}

	return false
}

type MeituLoginResp struct {
	Ret int `json:"ret"`
}

func checkMeitu(tu, tt string) bool {
	sc := fmt.Sprintf("%s%s%s%s", constants.MEITU_APP_ID, tu, tt, constants.MEITU_LOGIN_KEY)
	if sign, err := common.GenMD5S(sc); err != nil {
		return false
	} else if res, err := common.GetUrl(fmt.Sprintf(CHECK_MEITU_URL, constants.MEITU_APP_ID, tu, tt, sign)); err != nil {
		return false
	} else {
		ress := string(res)
		Info("check meitu login result:", ress)

		var resp MeituLoginResp
		if err := json.Unmarshal(res, &resp); err != nil {
			Err(err)
			return false
		}

		if resp.Ret == 100 {
			return true
		}
	}

	return false
}

type BukaLoginResp struct {
	Code int `json:"code"`
}

func checkBuka(tu string, token string) bool {
	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	sc := fmt.Sprintf("access_token=%s&supplier_id=%s&t=%s&uid=%s%s", token, constants.BUKA_SUPPLY_ID, ts, tu, constants.BUKA_SUPPLY_KEY)
	if sign, err := common.GenMD5S(sc); err != nil {
		return false
	} else {
		if res, err := common.GetUrl(fmt.Sprintf(CHECK_BUKA_URL, constants.BUKA_SUPPLY_ID, ts, tu, token, sign)); err != nil {
			return false
		} else {
			ress := string(res)
			Info("check buba result:", ress)

			var resp BukaLoginResp
			if err := json.Unmarshal(res, &resp); err != nil {
				Err(err)
				return false
			}

			if resp.Code == 1 {
				return true
			}
		}
	}
	return false
}

type LoginResp4399 struct {
	Code string `json:"code"`
}

func check4399(uid string, token string) bool {
	if res, err := common.GetUrl(fmt.Sprintf(CHECK_URL_4399, uid, token)); err != nil {
		return false
	} else {
		ress := string(res)
		Info("4399 check login result:", ress)

		var resp LoginResp4399
		if err := json.Unmarshal(res, &resp); err != nil {
			Err(err)
			return false
		}

		if resp.Code == "100" {
			return true
		}
	}
	return false
}

type KuaikanCheckLoginResp struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Data    KuaikanCheckLoginRespData `json:"data"`
}

type KuaikanCheckLoginRespData struct {
	Ret bool `json:"ret"`
}

func checkKuaikan(tusername string, token string) bool {
	//gen sign
	sc := "access_token=" + token + "&app_id=" + constants.KUAIKAN_APP_ID + "&open_id=" + tusername + "&key=" + constants.KUAIKAN_APP_SECRET
	if sign, err := common.GenMD5SBase64(sc); err != nil {
		return false
	} else {
		sign = url.QueryEscape(sign)
		url := fmt.Sprintf(KUAIKAN_CHECK_URL, constants.KUAIKAN_APP_ID, tusername, token, sign)
		Info("kuaikan check url:", url)
		if res, err := common.GetUrl(url); err != nil {
			return false
		} else {
			ress := string(res)
			Info("kuaikan check login result:", ress)

			var resp KuaikanCheckLoginResp
			if err := json.Unmarshal([]byte(res), &resp); err != nil {
				Err(err)
				return false
			}

			if resp.Data.Ret {
				return true
			}
		}
	}

	return false
}

func checkHuaweiToken(appid string, ts string, playerid string, inToken string) bool {
	s := appid + ts + playerid

	if res, err := common.Execute(javaPath, "bin_java/RSAUtil", s, inToken, HUAWEI_LOGIN_KEY); err != nil {
		return false
	} else {
		res = strings.TrimSpace(res)
		Info("huawei check result:", res)

		if res == "true" {
			return true
		}
	}

	return false
}

type VIVOResp struct {
	RetCode int      `json:"retcode"`
	Data    VIVOData `json:"data"`
}

type VIVOData struct {
	Success bool   `json:"success"`
	OpenID  string `json:"openid"`
}

func checkVIVO(authToken string) bool {
	if res, err := common.GetUrl(fmt.Sprintf(VIVO_CHECK_URL, authToken)); err != nil {
		Err(err)
		return false
	} else {
		Info("vivo check result:", authToken, string(res))
		var resp VIVOResp
		if err := json.Unmarshal(res, &resp); err != nil {
			Err(err)
			return false
		}
		if resp.RetCode == 0 {
			return true
		}
	}
	return false
}
