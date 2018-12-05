package request

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

import (
	. "logger"
)

var (
	client http.Client
)

const (
	TIME_OUT time.Duration = time.Duration(30)
)

func init() {
	client = http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(TIME_OUT * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*TIME_OUT)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			DisableKeepAlives: true,
		},
	}
}

func RequestWithValues(reqUrl string, values map[string]string) (ret []byte, err error) {
	data := url.Values{}
	for k, v := range values {
		data.Set(k, v)
	}

	var resp *http.Response
	resp, err = client.PostForm(reqUrl, data)
	if err != nil {
		Err("request with values err:", err)
		return
	}

	defer func() {
		resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		ret, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			Err("request with values read response body error:", err)
			return
		} else {
			return
		}
	} else {
		Err("request with values err, http status code:", resp.StatusCode)
		err = errors.New("http status:" + resp.Status)
		return
	}
}
