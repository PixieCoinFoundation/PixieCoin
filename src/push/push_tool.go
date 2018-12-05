package push

import (
	// "appcfg"
	"common"
	// "constants"
	// "crypto/md5"
	// "encoding/hex"
	"encoding/json"
	"fmt"
	// "io/ioutil"
	// "net/http"
	// "strconv"
	// "strings"
	. "logger"
	"net/url"
	"time"
)

// func tPushTagIOS(title, content, tag string) {
// 	tagList := []string{tag}
// 	tb, _ := json.Marshal(tagList)

// 	m := TPushIOSMsg{
// 		Aps: TPushIOSAps{Alert: content},
// 	}
// 	mb, _ := json.Marshal(m)

// 	now := time.Now().In(cnLocation).Unix()

// 	src := fmt.Sprintf("%s%saccess_id=%senvironment=%smessage=%smessage_type=%dtags_list=%stags_op=%stimestamp=%d%s", "GET", TPUSH_TAG_ADDR, tpush_ios_access_id, ios_environment, string(mb), 1, string(tb), "OR", now, tpush_ios_secret_key)
// 	sign, _ := common.GenMD5S(src)
// 	url := fmt.Sprintf("http://%s?access_id=%s&environment=%s&message=%s&message_type=%d&timestamp=%d&tags_list=%s&tags_op=%s&sign=%s", TPUSH_TAG_ADDR, tpush_ios_access_id, ios_environment, url.QueryEscape(string(mb)), 1, now, url.QueryEscape(string(tb)), "OR", sign)

// 	if r, e := common.GetUrl(url); e != nil {
// 		Err(e, string(r))
// 	} else {
// 		Info("push to ios tag:", tag, "resp:", string(r))
// 	}
// }

// func tPushTagAndroid(title, content, tag string) {
// 	tagList := []string{tag}
// 	tb, _ := json.Marshal(tagList)

// 	m := TPushAndroidMsg{
// 		Title:   title,
// 		Content: content,
// 	}
// 	mb, _ := json.Marshal(m)

// 	now := time.Now().In(cnLocation).Unix()

// 	src := fmt.Sprintf("%s%saccess_id=%smessage=%smessage_type=%dtags_list=%stags_op=%stimestamp=%d%s", "GET", TPUSH_TAG_ADDR, tpush_android_access_id, string(mb), 1, string(tb), "OR", now, tpush_android_secret_key)
// 	sign, _ := common.GenMD5S(src)
// 	url := fmt.Sprintf("http://%s?access_id=%s&message=%s&message_type=%d&timestamp=%d&tags_list=%s&tags_op=%s&sign=%s", TPUSH_TAG_ADDR, tpush_android_access_id, url.QueryEscape(string(mb)), 1, now, url.QueryEscape((string(tb))), "OR", sign)

// 	if r, e := common.GetUrl(url); e != nil {
// 		Err(e, string(r))
// 	} else {
// 		Info("push to android tag:", tag, "resp:", string(r))
// 	}
// }

func tPushOneIOS(username string, title, content string) {
	m := TPushIOSMsg{
		Aps: TPushIOSAps{Alert: content},
	}
	mb, _ := json.Marshal(m)

	now := time.Now().In(cnLocation).Unix()

	src := fmt.Sprintf("%s%saccess_id=%saccount=%senvironment=%smessage=%smessage_type=%dtimestamp=%d%s", "GET", TPUSH_ONE_ADDR, tpush_ios_access_id, username, ios_environment, string(mb), 1, now, tpush_ios_secret_key)
	sign, _ := common.GenMD5S(src)
	url := fmt.Sprintf("http://%s?access_id=%s&account=%s&environment=%s&message=%s&message_type=%d&timestamp=%d&sign=%s", TPUSH_ONE_ADDR, tpush_ios_access_id, username, ios_environment, url.QueryEscape(string(mb)), 1, now, sign)

	if r, e := common.GetUrl(url); e != nil {
		Err(e, string(r))
	} else {
		Info("push to ios:", username, "resp:", string(r))
	}
}

func tPushOneAndroid(username string, title, content string) {
	m := TPushAndroidMsg{
		Title:   title,
		Content: content,
	}
	mb, _ := json.Marshal(m)

	now := time.Now().In(cnLocation).Unix()

	src := fmt.Sprintf("%s%saccess_id=%saccount=%smessage=%smessage_type=%dtimestamp=%d%s", "GET", TPUSH_ONE_ADDR, tpush_android_access_id, username, string(mb), 1, now, tpush_android_secret_key)
	sign, _ := common.GenMD5S(src)
	url := fmt.Sprintf("http://%s?access_id=%s&account=%s&message=%s&message_type=%d&timestamp=%d&sign=%s", TPUSH_ONE_ADDR, tpush_android_access_id, username, url.QueryEscape(string(mb)), 1, now, sign)

	if r, e := common.GetUrl(url); e != nil {
		Err(e, string(r))
	} else {
		Info("push to android:", username, "resp:", string(r))
	}
}

func tPushAllAndroid(title, content string) {
	m := TPushAndroidMsg{
		Title:   title,
		Content: content,
	}
	mb, _ := json.Marshal(m)

	now := time.Now().In(cnLocation).Unix()

	src := fmt.Sprintf("%s%saccess_id=%smessage=%smessage_type=%dtimestamp=%d%s", "GET", TPUSH_ALL_ADDR, tpush_android_access_id, string(mb), 1, now, tpush_android_secret_key)
	sign, _ := common.GenMD5S(src)
	url := fmt.Sprintf("http://%s?access_id=%s&message=%s&message_type=%d&timestamp=%d&sign=%s", TPUSH_ALL_ADDR, tpush_android_access_id, url.QueryEscape(string(mb)), 1, now, sign)

	if r, e := common.GetUrl(url); e != nil {
		Err(e, string(r))
	} else {
		Info("push to all android resp:", string(r))
	}
}

func tPushAllIOS(title, content string) {
	m := TPushIOSMsg{
		Aps: TPushIOSAps{Alert: content},
	}
	mb, _ := json.Marshal(m)

	now := time.Now().In(cnLocation).Unix()

	src := fmt.Sprintf("%s%saccess_id=%senvironment=%smessage=%smessage_type=%dtimestamp=%d%s", "GET", TPUSH_ALL_ADDR, tpush_ios_access_id, ios_environment, string(mb), 1, now, tpush_ios_secret_key)
	sign, _ := common.GenMD5S(src)
	url := fmt.Sprintf("http://%s?access_id=%s&environment=%s&message=%s&message_type=%d&timestamp=%d&sign=%s", TPUSH_ALL_ADDR, tpush_ios_access_id, ios_environment, url.QueryEscape(string(mb)), 1, now, sign)

	if r, e := common.GetUrl(url); e != nil {
		Err(e, string(r))
	} else {
		Info("push to all ios resp:", string(r))
	}
}

// func setTagForAccount(username, tag string) {
// 	if tpush_android_secret_key != "" && tpush_android_access_id != "" {
// 		setTagForAccountAndroid(username, tag)
// 	}

// 	if tpush_ios_secret_key != "" && tpush_ios_access_id != "" {
// 		setTagForAccountIOS(username, tag)
// 	}
// }

// func delTagForAccount(username, tag string) {
// 	if tpush_android_secret_key != "" && tpush_android_access_id != "" {
// 		delTagForAccountAndroid(username, tag)
// 	}

// 	if tpush_ios_secret_key != "" && tpush_ios_access_id != "" {
// 		delTagForAccountIOS(username, tag)
// 	}
// }

// func setTagForAccountIOS(username, tag string) (err error) {
// 	var tokens []string
// 	if tokens, err = getIOSTokensFromAccount(username); err != nil {
// 		return
// 	}

// 	if len(tokens) > 0 {
// 		var tagTokenList [][]string
// 		for _, token := range tokens {
// 			tagTokenList = append(tagTokenList, []string{tag, token})
// 		}

// 		tb, _ := json.Marshal(tagTokenList)

// 		err = setTagForAccountIOS1(string(tb))
// 	} else {
// 		Info("no ios tokens:", username)
// 	}
// 	return
// }

// func delTagForAccountIOS(username, tag string) (err error) {
// 	var tokens []string
// 	if tokens, err = getIOSTokensFromAccount(username); err != nil {
// 		return
// 	}

// 	if len(tokens) > 0 {
// 		var tagTokenList [][]string
// 		for _, token := range tokens {
// 			tagTokenList = append(tagTokenList, []string{tag, token})
// 		}

// 		tb, _ := json.Marshal(tagTokenList)

// 		err = delTagForAccountIOS1(string(tb))
// 	} else {
// 		Info("no ios tokens:", username)
// 	}
// 	return
// }

// func setTagForAccountIOS1(tagTokenListJson string) (err error) {
// 	now := time.Now().In(cnLocation).Unix()

// 	src := fmt.Sprintf("%s%saccess_id=%stag_token_list=%stimestamp=%d%s", "GET", TPUSH_BATCH_SET_TAG_ADDR, tpush_ios_access_id, tagTokenListJson, now, tpush_ios_secret_key)
// 	sign, _ := common.GenMD5S(src)
// 	url := fmt.Sprintf("http://%s?access_id=%s&tag_token_list=%s&timestamp=%d&sign=%s", TPUSH_BATCH_SET_TAG_ADDR, tpush_ios_access_id, tagTokenListJson, now, sign)

// 	if r, e := common.GetUrl(url); e != nil {
// 		Err(e, string(r))
// 		err = e
// 	} else {
// 		Info("setTagForAccountIOS1:", tagTokenListJson, "resp:", string(r))
// 	}
// 	return
// }

// func delTagForAccountIOS1(tagTokenListJson string) (err error) {
// 	now := time.Now().In(cnLocation).Unix()

// 	src := fmt.Sprintf("%s%saccess_id=%stag_token_list=%stimestamp=%d%s", "GET", TPUSH_BATCH_DEL_TAG_ADDR, tpush_ios_access_id, tagTokenListJson, now, tpush_ios_secret_key)
// 	sign, _ := common.GenMD5S(src)
// 	url := fmt.Sprintf("http://%s?access_id=%s&tag_token_list=%s&timestamp=%d&sign=%s", TPUSH_BATCH_DEL_TAG_ADDR, tpush_ios_access_id, tagTokenListJson, now, sign)

// 	if r, e := common.GetUrl(url); e != nil {
// 		Err(e, string(r))
// 		err = e
// 	} else {
// 		Info("delTagForAccountIOS1:", tagTokenListJson, "resp:", string(r))
// 	}
// 	return
// }

// func setTagForAccountAndroid(username, tag string) (err error) {
// 	var tokens []string
// 	if tokens, err = getAndroidTokensFromAccount(username); err != nil {
// 		return
// 	}

// 	if len(tokens) > 0 {
// 		var tagTokenList [][]string
// 		for _, token := range tokens {
// 			tagTokenList = append(tagTokenList, []string{tag, token})
// 		}

// 		tb, _ := json.Marshal(tagTokenList)

// 		err = setTagForAccountAndroid1(string(tb))
// 	} else {
// 		Info("no android tokens:", username)
// 	}
// 	return
// }

// func delTagForAccountAndroid(username, tag string) (err error) {
// 	var tokens []string
// 	if tokens, err = getAndroidTokensFromAccount(username); err != nil {
// 		return
// 	}

// 	if len(tokens) > 0 {
// 		var tagTokenList [][]string
// 		for _, token := range tokens {
// 			tagTokenList = append(tagTokenList, []string{tag, token})
// 		}

// 		tb, _ := json.Marshal(tagTokenList)

// 		err = delTagForAccountAndroid1(string(tb))
// 	} else {
// 		Info("no android tokens:", username)
// 	}
// 	return
// }

// func setTagForAccountAndroid1(tagTokenListJson string) (err error) {
// 	now := time.Now().In(cnLocation).Unix()

// 	src := fmt.Sprintf("%s%saccess_id=%stag_token_list=%stimestamp=%d%s", "GET", TPUSH_BATCH_SET_TAG_ADDR, tpush_android_access_id, tagTokenListJson, now, tpush_android_secret_key)
// 	sign, _ := common.GenMD5S(src)
// 	url := fmt.Sprintf("http://%s?access_id=%s&tag_token_list=%s&timestamp=%d&sign=%s", TPUSH_BATCH_SET_TAG_ADDR, tpush_android_access_id, tagTokenListJson, now, sign)

// 	if r, e := common.GetUrl(url); e != nil {
// 		Err(e, string(r))
// 		err = e
// 	} else {
// 		Info("setTagForAccountAndroid1:", tagTokenListJson, "resp:", string(r))
// 	}
// 	return
// }

// func delTagForAccountAndroid1(tagTokenListJson string) (err error) {
// 	now := time.Now().In(cnLocation).Unix()

// 	src := fmt.Sprintf("%s%saccess_id=%stag_token_list=%stimestamp=%d%s", "GET", TPUSH_BATCH_DEL_TAG_ADDR, tpush_android_access_id, tagTokenListJson, now, tpush_android_secret_key)
// 	sign, _ := common.GenMD5S(src)
// 	url := fmt.Sprintf("http://%s?access_id=%s&tag_token_list=%s&timestamp=%d&sign=%s", TPUSH_BATCH_DEL_TAG_ADDR, tpush_android_access_id, tagTokenListJson, now, sign)

// 	if r, e := common.GetUrl(url); e != nil {
// 		Err(e, string(r))
// 		err = e
// 	} else {
// 		Info("delTagForAccountAndroid1:", tagTokenListJson, "resp:", string(r))
// 	}
// 	return
// }

// func getTokendsFromAccount(username string) (res []string, err error) {
// 	var ta, ti []string
// 	res = make([]string, 0)

// 	if ta, err = getAndroidTokensFromAccount(username); err != nil {
// 		return
// 	} else if len(ta) > 0 {
// 		res = append(res, ta...)
// 	}

// 	if ti, err = getIOSTokensFromAccount(username); err != nil {
// 		return
// 	} else if len(ti) > 0 {
// 		res = append(res, ti...)
// 	}

// 	return
// }

// type ListAccountTokensResp struct {
// 	RetCode int                    `json:"ret_code"`
// 	Result  ListAccountTokensResp1 `json:"result"`
// }

// type ListAccountTokensResp1 struct {
// 	Tokens []string `json:"tokens"`
// }

// func getAndroidTokensFromAccount(username string) (tokens []string, err error) {
// 	now := time.Now().In(cnLocation).Unix()

// 	src := fmt.Sprintf("%s%saccess_id=%saccount=%stimestamp=%d%s", "GET", TPUSH_LIST_ACCOUNT_TOKEN_ADDR, tpush_android_access_id, username, now, tpush_android_secret_key)
// 	sign, _ := common.GenMD5S(src)
// 	url := fmt.Sprintf("http://%s?access_id=%s&account=%s&timestamp=%d&sign=%s", TPUSH_LIST_ACCOUNT_TOKEN_ADDR, tpush_android_access_id, username, now, sign)

// 	if r, e := common.GetUrl(url); e != nil {
// 		Err(e, string(r))
// 		err = e
// 		return
// 	} else {
// 		Info("getAndroidTokensFromAccount:", username, "resp:", string(r))

// 		var resp ListAccountTokensResp
// 		if err = json.Unmarshal(r, &resp); err != nil {
// 			Err(err)
// 			return
// 		}

// 		if resp.RetCode == 0 {
// 			tokens = resp.Result.Tokens
// 		}

// 		return
// 	}
// }

// func getIOSTokensFromAccount(username string) (tokens []string, err error) {
// 	now := time.Now().In(cnLocation).Unix()

// 	src := fmt.Sprintf("%s%saccess_id=%saccount=%stimestamp=%d%s", "GET", TPUSH_LIST_ACCOUNT_TOKEN_ADDR, tpush_ios_access_id, username, now, tpush_ios_secret_key)
// 	sign, _ := common.GenMD5S(src)
// 	url := fmt.Sprintf("http://%s?access_id=%s&account=%s&timestamp=%d&sign=%s", TPUSH_LIST_ACCOUNT_TOKEN_ADDR, tpush_ios_access_id, username, now, sign)

// 	if r, e := common.GetUrl(url); e != nil {
// 		Err(e, string(r))
// 		err = e
// 		return
// 	} else {
// 		Info("getIOSTokensFromAccount:", username, "resp:", string(r))

// 		var resp ListAccountTokensResp
// 		if err = json.Unmarshal(r, &resp); err != nil {
// 			Err(err)
// 			return
// 		}

// 		if resp.RetCode == 0 {
// 			tokens = resp.Result.Tokens
// 		}

// 		return
// 	}
// }
