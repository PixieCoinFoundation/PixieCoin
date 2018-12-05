package kakao

import (
	"constants"
	"encoding/json"
	"fmt"
	// "io/ioutil"
	// "log"
	. "logger"
	// "net"
	"net/http"
	"strings"
	"time"
	"tools"
)

type bodyparams struct {
	Currency           string  `json:"currency"`
	Price              float64 `json:"price"`
	Os                 string  `json:"os"`
	Market             string  `json:"market"`
	MarketProductid    string  `json:"marketProductId"`
	Country            string  `json:"country"`
	MarketPurchaseTime int64   `json:"marketPurchaseTime"`
}

type Aoslog struct {
	ID      int
	IOS     string
	Usdcost float64
}
type Korlog struct {
	ID      int
	Aos     string
	Krwcost float64
}

var Aosloglist = []Aoslog{
	Aoslog{
		ID:      1001,
		IOS:     "flerogames.com.ios.for.fw.jewel100",
		Usdcost: 2.19,
	},
	Aoslog{
		ID:      1002,
		IOS:     "flerogames.com.ios.for.fw.jewel250",
		Usdcost: 4.39,
	},
	Aoslog{
		ID:      1003,
		IOS:     "flerogames.com.ios.for.fw.jewel500",
		Usdcost: 9.89,
	},
	Aoslog{
		ID:      1004,
		IOS:     "flerogames.com.ios.for.fw.jewel1500",
		Usdcost: 29.69,
	},
	Aoslog{
		ID:      1005,
		IOS:     "flerogames.com.ios.for.fw.jewel2500",
		Usdcost: 49.49,
	},
	Aoslog{
		ID:      1006,
		IOS:     "flerogames.com.ios.for.fw.jewel5000",
		Usdcost: 98.99,
	},
	Aoslog{
		ID:      1007,
		IOS:     "flerogames.com.ios.for.fw.basicpack",
		Usdcost: 1.09,
	},
	Aoslog{
		ID:      1008,
		IOS:     "flerogames.com.ios.for.fw.dailypack",
		Usdcost: 4.39,
	},
	Aoslog{
		ID:      1009,
		IOS:     "flerogames.com.ios.for.fw.starterpack",
		Usdcost: 9.89,
	},
}
var Korloglist = []Korlog{
	Korlog{
		ID:      1001,
		Aos:     "flerogames.com.aos.kor.fw.jewel100",
		Krwcost: 2200,
	},
	Korlog{
		ID:      1002,
		Aos:     "flerogames.com.aos.kor.fw.jewel250",
		Krwcost: 5500,
	},
	Korlog{
		ID:      1003,
		Aos:     "flerogames.com.aos.kor.fw.jewel500",
		Krwcost: 11000,
	},
	Korlog{
		ID:      1004,
		Aos:     "flerogames.com.aos.kor.fw.jewel1500",
		Krwcost: 33000,
	},
	Korlog{
		ID:      1005,
		Aos:     "flerogames.com.aos.kor.fw.jewel2500",
		Krwcost: 55000,
	},
	Korlog{
		ID:      1006,
		Aos:     "flerogames.com.aos.kor.fw.jewel5000",
		Krwcost: 110000,
	},
	Korlog{
		ID:      1007,
		Aos:     "flerogames.com.ios.for.fw.basicpack",
		Krwcost: 1100,
	},
	Korlog{
		ID:      1008,
		Aos:     "flerogames.com.aos.kor.fw.dailypack",
		Krwcost: 5500,
	},
	Korlog{
		ID:      1009,
		Aos:     "flerogames.com.aos.kor.fw.starterpack",
		Krwcost: 11000,
	},
}

//KorKakaoSendLog 韩国kakao
func KorKakaoSendLog(playerid string, produceid int, os string) error {
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
	var bodySend bodyparams
	bodySend.Country = "kr"
	os = strings.ToLower(os)
	if os == "ios" {
		bodySend.Os = "iOS"
		bodySend.Market = "appStore"
		bodySend.Currency = "USD"
		for _, v := range Aosloglist {
			if v.ID == produceid {
				bodySend.MarketProductid = v.IOS
				bodySend.MarketPurchaseTime = time.Now().Unix()
				bodySend.Price = v.Usdcost
			}
		}
	} else {
		bodySend.Os = "android"
		bodySend.Market = "googlePlay"
		bodySend.Currency = "KRW"
		for _, v := range Korloglist {
			if v.ID == produceid {
				bodySend.MarketProductid = v.Aos
				bodySend.MarketPurchaseTime = time.Now().Unix()
				bodySend.Price = v.Krwcost
			}
		}
	}
	v, _ := json.Marshal(bodySend)

	if reqest, err := http.NewRequest("POST", constants.KAKAOOGAME_SDK_SENDLOG, strings.NewReader(string(v))); err != nil {
		Err(err)
		return err
	} else {
		reqest.Header.Set("appId", constants.KAKAO_APP_ID)
		reqest.Header.Set("appSecret", constants.KAKAP_NATIVE_APP_KEY)
		reqest.Header.Set("playerId", playerid)
		reqest.Header.Set("Content-Type", "application/json;charset=UTF-8")
		reqest.Header.Set("Authorization", fmt.Sprintf("KakaoAK %s", constants.KAKAO_ADMIN_KEY))

		if resp, err := tools.OtakuHttpClient.Do(reqest); err != nil {
			Err(err)
			return err
		} else {
			resp.Body.Close()
		}
		// data, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	Err(err)
		// 	return err
		// }
		// var re RespData
		// err = json.Unmarshal(data, &re)
		// if resp.StatusCode == 200 {
		// 	return nil
		// }
	}

	return nil
}

type RespData struct {
	LogId string `json:"logId"`
}
