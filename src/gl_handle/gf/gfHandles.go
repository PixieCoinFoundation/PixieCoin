package gf

import (
	"encoding/json"
	"fmt"
)

import ()

var RequestHandles map[string]func(string, string, chan string, chan int) = map[string]func(string, string, chan string, chan int){
// servers
// "registerServer":          RegisterServer,
// "updateServerPlayerCount": UpdateServerPlayerCount,
}

func unmarshal(data string, v interface{}, errflag chan int) (err error) {
	if err = json.Unmarshal([]byte(data), v); err != nil {
		fmt.Println("gf handle unmarshal error")
		errflag <- 1
	}
	return
}

func marshal(v interface{}, errflag chan int) string {
	if jsonbyte, err := json.Marshal(v); err != nil {
		fmt.Println("gf handle marshal error")
		errflag <- 1
		return ""
	} else {
		return string(jsonbyte)
	}
}
