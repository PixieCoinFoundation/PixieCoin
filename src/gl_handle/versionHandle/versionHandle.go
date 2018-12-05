package versionHandle

import (
	"appcfg"
	"constants"
	"encoding/json"
	// "fmt"
	"gl_dao"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	// "strings"
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
)

type versionInfo struct {
	Username             string `json:"username"`      //用户名
	CoreVersion          string `json:"coreVersion"`   // 核心版本
	ScriptVersion        string `json:"scriptVersion"` // 脚本版本
	ThirdChannel         string `json:"thirdChannel"`  // 第三方渠道ID
	VersionFile          string `json:"versionFile"`
	KorSpecialTishenTag  bool   `json:"korSpecialTishenTag"`
	KorSpecialTishenTag2 bool   `json:"korSpecialTishenTag2"`
}

type versionResp struct {
	Status             ErrorCode `json:"status"`
	CoreVersion        string    `json:"coreVersion"`
	CoreDownloadURL    string    `json:"coreDownloadURL"`
	CoreDownloadURL2   string    `json:"coreDownloadURL2"`
	ScriptVersion      string    `json:"scriptVersion"`
	VersionFile        string    `json:"versionFile"`
	FileServerAddress  string    `json:"address"`
	FileServerAddress2 string    `json:"address2"`
	FileServerAddress3 string    `json:"address3"`
	Tishen             bool      `json:"tishen"` //true代表大版本不同时连接提审服 false表示进行强更
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 检查总更新
func CheckVersion(w http.ResponseWriter, r *http.Request) {
	var responseData versionResp
	defer func() {
		jsonByte, _ := json.Marshal(responseData)
		w.Write(encode(jsonByte))
	}()

	defer func() {
		r.Body.Close()
	}()

	data, readError := ioutil.ReadAll(r.Body)
	if readError != nil {
		responseData = versionResp{
			Status: ERR_UNKNOW,
		}

		return
	}

	data = simple.Unpack(data)
	data = xxtea.Decrypt(data, constants.XXTEA_MSG_KEY)
	uncompressed, snappyErr := snappy.Decode(nil, data)
	if snappyErr != nil {
		Err("snappy error:", snappyErr.Error())
		responseData = versionResp{
			Status: ERR_UNKNOW,
		}
		return
	}

	Info("client check version:", string(uncompressed))
	var version versionInfo
	err := json.Unmarshal(uncompressed, &version)
	if err == nil {
		responseData = checkVersion(version.Username, version.CoreVersion, version.ScriptVersion, version.ThirdChannel, version.VersionFile, version.KorSpecialTishenTag2)
	} else {
		responseData = versionResp{
			Status: ERR_HANDLE_PARAMS_WRONG,
		}
	}
	return
}

//提审专用
func CheckVersionTS(w http.ResponseWriter, r *http.Request) {
	responseData := versionResp{
		Status:             2002,
		CoreVersion:        "0.16",
		ScriptVersion:      "1",
		CoreDownloadURL:    "http://yjdyc.biligame.com/yxz/",
		CoreDownloadURL2:   "http://yjdyc.biligame.com/yxz/",
		VersionFile:        "group1/M00/00/00/wKhtg1be0CmALMHBAAAUi9533Gg6302986",
		FileServerAddress:  "http://fg-fw.cdn.wemade.com/",
		FileServerAddress2: "http://fg-fw.cdn.wemade.com/",
		FileServerAddress3: "http://fg-fw.cdn.wemade.com/",
		Tishen:             true,
	}

	jsonByte, _ := json.Marshal(responseData)
	w.Write(encode(jsonByte))

	return
}

func checkVersion(username, coreVersion, scriptVersion, thirdChannel, inputVersionFile string, korSpecialTag2 bool) (resp versionResp) {
	gfVersion, tvf, err := version.GetCurrentVersion(thirdChannel)

	if err != nil {
		resp.Status = GAME_GET_VERSION_ERROR
	} else {
		resp.CoreVersion = gfVersion.CoreVersion
		resp.ScriptVersion = gfVersion.ScriptVersion
		resp.CoreDownloadURL = gfVersion.CoreDownloadURL
		resp.CoreDownloadURL2 = gfVersion.CoreDownloadURL2
		resp.VersionFile = gfVersion.VersionFile
		resp.FileServerAddress = gfVersion.FileServerAddress
		resp.FileServerAddress2 = gfVersion.FileServerAddress2
		resp.FileServerAddress3 = gfVersion.FileServerAddress3
		resp.Tishen = gfVersion.Tishen

		if coreVersion != gfVersion.CoreVersion {
			// 先检查core version

			// core version分为两部分， 小数点前和小数点后
			clientVersion, err0 := strconv.ParseFloat(coreVersion, 64)

			//server version
			serverVersion, err1 := strconv.ParseFloat(gfVersion.CoreVersion, 64)
			// Info(clientVersion,serverVersion)

			if err0 != nil || err1 != nil {
				resp.Status = ERR_PARSE_SERVER_GAME_VERSION
			} else {
				if clientVersion < serverVersion {
					resp.Status = GAME_UPDATE_CORE_VERSION
					resp.Tishen = false
				} else {
					if gfVersion.Tishen {
						resp.Status = GAME_UPDATE_CORE_VERSION
					}
				}
			}
		} else if scriptVersion != gfVersion.ScriptVersion {
			// 再检查script version

			si, _ := strconv.Atoi(scriptVersion)
			sni, _ := strconv.Atoi(gfVersion.ScriptVersion)
			if si >= sni {
				resp.Status = 0
				resp.ScriptVersion = scriptVersion
			} else {
				if version.AllowPatch || gf.AccountIsBlack(username) {
					resp.Status = GAME_UPDATE_SCRIPT_VERSION
				} else {
					resp.Status = 0
					resp.ScriptVersion = scriptVersion
					if tvf != "" {
						resp.VersionFile = tvf
					} else {
						resp.VersionFile = inputVersionFile
					}
				}
			}
		} else {
			resp.Status = 0
		}

		if appcfg.GetLanguage() == constants.KOR_LANGUAGE && korSpecialTag2 && checkKorSpecialTag2() {
			resp.Status = GAME_UPDATE_CORE_VERSION
			resp.Tishen = true
		}
	}

	return
}

func checkKorSpecialTag2() bool {
	var checkResult int
	if err := gl_dao.GetCheckKorSpecialTagStmt.QueryRow().Scan(&checkResult); err != nil {
		Err(err)
		return false
	} else {
		return checkResult == 1
	}
}

type GetVersionResp struct {
	Status        int    `json:"status"`
	CoreVersion   string `json:"coreVersion"`
	DownloadURL   string `json:"downloadUrl"`
	DownloadURL2  string `json:"downloadUrl2"`
	ScriptVersion string `json:"scriptVersion"`
	VersionFile   string `json:"versionFile"`
}

// type UpdateVersionRequest struct {
// 	CoreVersion   string `json:"coreVersion"`
// 	DownloadURL   string `json:"downloadUrl"`
// 	DownloadURL2  string
// 	ScriptVersion string `json:"scriptVersion"`
// 	VersionFile   string `json:"versionFile"`
// }

// type UpdateVersionResp struct {
// 	Status int `json:"status"`
// }

// func UpdateVersion(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()

// 	var input UpdateVersionRequest
// 	var ret UpdateVersionResp
// 	ret.Status = 0

// 	input.CoreVersion = r.PostFormValue("coreVersion")
// 	input.DownloadURL = r.PostFormValue("downloadURL")
// 	input.DownloadURL2 = r.PostFormValue("downloadURL2")
// 	input.ScriptVersion = r.PostFormValue("scriptVersion")
// 	input.VersionFile = r.PostFormValue("versionFile")

// 	if input.CoreVersion == "" || input.DownloadURL == "" || input.DownloadURL2 == "" || input.ScriptVersion == "" || input.VersionFile == "" {
// 		ret.Status = 2004
// 	}

// 	if ret.Status == 0 {
// 		err := version.SetNewVersion(input.CoreVersion, input.DownloadURL, input.DownloadURL2, input.ScriptVersion, input.VersionFile)
// 		if err != nil {
// 			ret.Status = 2005
// 			fmt.Println("error:", err.Error())
// 		} else {
// 			ret.Status = 0
// 		}
// 	}

// 	jsonByte, _ := json.Marshal(ret)

// 	fmt.Fprintf(w, string(jsonByte))
// }

func encode(input []byte) []byte {
	compressed := snappy.Encode(nil, input)
	encrypted := xxtea.Encrypt(compressed, constants.XXTEA_MSG_KEY)
	data := simple.Pack(encrypted)

	return data
}
