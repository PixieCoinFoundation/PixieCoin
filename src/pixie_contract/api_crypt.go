package pixie_contract

import (
	"constants"
	"encoding/json"
	"encrypt/simple"
	"encrypt/xxtea"
	. "logger"
	"snappy"
)

func EncryptResponse(resp interface{}) []byte {
	jsonByte, _ := json.Marshal(resp)
	compressed := snappy.Encode(nil, jsonByte)
	encrypted := xxtea.Encrypt(compressed, constants.XXTEA_MSG_KEY)
	sendData := simple.Pack(encrypted)

	return sendData
}

func PixieGetInfoFromGFToken(token string) (db string) {
	return simple.SimpleDecrypt(token)
}

func PixieUnmarshal(data string, v interface{}) (err error) {
	if err = json.Unmarshal([]byte(data), v); err != nil {
		Err("gf http dispatch unmarshal error, data:", err, data)
	}
	return
}

func PixieMarshal(v interface{}) (string, error) {
	if jsonbyte, err := json.Marshal(v); err != nil {
		Err("gf http dispatch marshal error:", err.Error())
		return "", err
	} else {
		return string(jsonbyte), nil
	}
}
