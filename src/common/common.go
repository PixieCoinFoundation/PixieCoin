package common

import (
	// "appcfg"
	"bytes"
	"constants"
	"crypto/hmac"
	"crypto/md5"
	cr "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	. "language"
	. "logger"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	. "pixie_contract/api_specification"
	"rao"
	"service/clothes"
	"strconv"
	"strings"
	"time"
	. "types"
)

const (
	SECRET      = "#^_^#"
	NAME_PREFIX = "gift_name_"
	CNT_PREFIX  = "gift_cnt_"

	SUIT_PREFIX     = "suit_name_"
	SUIT_CNT_PREFIX = "suit_cnt_"

	CP_NAME_PREFIX = "p_gift_name_"
	CP_CNT_PREFIX  = "p_gift_cnt_"

	GMLOG_TOKEN_PREFIX = "GMLOG"
)

var st time.Time
var gst time.Time
var ClothesZeroErr = errors.New("clothes cnt less than 0")
var ClothesIllegalErr = errors.New("clothes reward illegal!")

func init() {
	st, _ = time.ParseInLocation("20060102150405", constants.PK_START_DAY, time.Local)
	gst, _ = time.ParseInLocation("20060102150405", constants.GUILD_START_DAY, time.Local)

	// RecurDecryptFiles("/Users/zzd/Downloads/zihan")
	// panic("end")
}

func RecurDecryptFiles(dirName string) {
	if err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// don't process folders themselves
		if path == dirName {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".png") {
			Execute("scripts/DecryptFile_osx", path)
		}

		return err
	}); err != nil {
		Err(err)
	}
}

func Sha1Encrypt(content string) string {
	h := sha1.New()
	h.Write([]byte(content))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func HmacSha1Encrypt(content, keys string) string {
	//hmac ,use sha1
	key := []byte(keys)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(content))
	res := fmt.Sprintf("%x", mac.Sum(nil))
	return res
}

func HmacSha1EncryptBase64(content, keys string) string {
	key := []byte(keys)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(content))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func RsaEncrypt(plaintext []byte, pub_key []byte) ([]byte, error) {
	block, _ := pem.Decode(pub_key)
	if block == nil {
		return nil, errors.New("invalid rsa public key")
	}

	pubInf, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInf.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(cr.Reader, pub, plaintext)
}

func GetGMLogToken(ds string, address string) string {
	return fmt.Sprintf("%s_%s_%s", GMLOG_TOKEN_PREFIX, ds, address)
	// return GMLOG_TOKEN_PREFIX + ds
}

func GetGMLogLikeToken(ds string) string {
	return fmt.Sprintf("%s_%s_", GMLOG_TOKEN_PREFIX, ds) + `%`
	// return GMLOG_TOKEN_PREFIX + ds
}

func GetMonthToken() string {
	return time.Now().Format("200601")
}

func GetMonthTokenByGapDay(gapDay int) string {
	return time.Now().AddDate(0, 0, gapDay).Format("200601")
}

func GenVIP(buyDiamond int) int {
	if buyDiamond < constants.V1_BUY_DIAMOND {
		return 0
	} else if buyDiamond < constants.V2_BUY_DIAMOND {
		return 1
	} else if buyDiamond < constants.V3_BUY_DIAMOND {
		return 2
	} else if buyDiamond < constants.V4_BUY_DIAMOND {
		return 3
	} else if buyDiamond < constants.V5_BUY_DIAMOND {
		return 4
	} else if buyDiamond < constants.V6_BUY_DIAMOND {
		return 5
	} else if buyDiamond < constants.V7_BUY_DIAMOND {
		return 6
	} else if buyDiamond < constants.V8_BUY_DIAMOND {
		return 7
	} else if buyDiamond < constants.V9_BUY_DIAMOND {
		return 8
	} else {
		return 9
	}
}

func GapTiliVIP(fv, nv int) int {
	return getVIPAddTili(nv) - getVIPAddTili(fv)
}

func getVIPAddTili(vip int) int {
	if vip == 9 {
		return constants.V9_ADD_TILI
	} else if vip == 8 {
		return constants.V8_ADD_TILI
	} else if vip == 7 {
		return constants.V7_ADD_TILI
	} else if vip == 6 {
		return constants.V6_ADD_TILI
	} else if vip == 5 {
		return constants.V5_ADD_TILI
	} else if vip == 4 {
		return constants.V4_ADD_TILI
	} else if vip == 3 {
		return constants.V3_ADD_TILI
	} else if vip == 2 {
		return constants.V2_ADD_TILI
	} else if vip == 1 {
		return constants.V1_ADD_TILI
	}
	return 0
}

func GetUrl(url string) ([]byte, error) {
	Info(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetYearWeekToken() (string, string) {
	return GetYWToken(0, 0)
}

func GetYWToken(gap int, ggap int) (string, string) {
	n := time.Now()

	return GetYWTokenByTime(n, gap, ggap)
}

func GetYWTokenByTime(n time.Time, gap int, ggap int) (string, string) {
	//week
	w := n.AddDate(0, 0, 7*gap).Sub(st).Hours() / 7.0 / 24.0
	if w < 0 {
		w = 0
	}

	//gweek
	gapWeek := n.AddDate(0, 0, ggap*12*7).Sub(st).Hours() / 12.0 / 7.0 / 24.0
	if gapWeek < 0 {
		gapWeek = 0
	}

	return fmt.Sprintf("%s_%d", constants.PK_START_DAY, int(w)), fmt.Sprintf("%s_%d", constants.PK_START_DAY, int(gapWeek))
}

func GetGuildYearWeekToken() (string, string) {
	return GetGuildYWToken(0, 0)
}

func GetGuildYWToken(gap int, ggap int) (string, string) {
	n := time.Now()

	return GetGuildYWTokenByTime(n, gap, ggap)
}

func GetGuildYWTokenByTime(n time.Time, gap int, ggap int) (string, string) {
	//week
	w := n.AddDate(0, 0, 7*gap).Sub(gst).Hours() / 7.0 / 24.0
	if w < 0 {
		w = 0
	}

	//gweek
	gapWeek := n.AddDate(0, 0, ggap*12*7).Sub(gst).Hours() / 12.0 / 7.0 / 24.0
	if gapWeek < 0 {
		gapWeek = 0
	}

	return fmt.Sprintf("%s_%d", constants.GUILD_START_DAY, int(w)), fmt.Sprintf("%s_%d", constants.GUILD_START_DAY, int(gapWeek))
}

func GetClothesGifts(req *http.Request) (res []ClothesInfo, clo5s []string, err error) {
	res = make([]ClothesInfo, 0)
	clo5s = make([]string, 0)

	for k, v := range req.Form {
		if strings.HasPrefix(k, NAME_PREFIX) && len(v) > 0 && v[0] != "" {
			index := k[len(NAME_PREFIX):]
			for k1, v1 := range req.Form {
				if strings.HasPrefix(k1, CNT_PREFIX) && len(v1) > 0 && v1[0] != "" && k1[len(CNT_PREFIX):] == index {
					// Info("clothes id:", v[0], "cnt:", v1[0])
					idStr := strings.Split(v[0], "-")[0]

					var cnt int
					if cnt, err = strconv.Atoi(v1[0]); err != nil {
						Err(err)
						return
					} else if cnt < 0 {
						err = ClothesZeroErr
						return
					} else {
						exist, star := clothes.ClothesLegal(idStr)
						if !exist {
							err = ClothesIllegalErr
							return
						} else {
							c := ClothesInfo{
								ClothesID: idStr,
								Count:     cnt,
							}
							if star == 5 && cnt > 10 {
								clo5s = append(clo5s, v[0])
							}
							res = append(res, c)
						}
					}
				}
			}
		} else if strings.HasPrefix(k, SUIT_PREFIX) && len(v) > 0 && v[0] != "" {
			index := k[len(SUIT_PREFIX):]
			for k1, v1 := range req.Form {
				if strings.HasPrefix(k1, SUIT_CNT_PREFIX) && len(v1) > 0 && v1[0] != "" && k1[len(SUIT_CNT_PREFIX):] == index {
					// Info("clothes id:", v[0], "cnt:", v1[0])
					idStr := strings.Split(v[0], "-")[0]

					var cnt int
					if cnt, err = strconv.Atoi(v1[0]); err != nil {
						Err(err)
						return
					} else if cnt < 0 {
						err = ClothesZeroErr
						return
					} else {
						clos := clothes.GetSuitClothesList(idStr)
						for _, cs := range clos {
							exist, star := clothes.ClothesLegal(cs)
							if !exist {
								err = ClothesIllegalErr
								return
							} else {
								c := ClothesInfo{
									ClothesID: cs,
									Count:     cnt,
								}
								if star == 5 && cnt > 10 {
									clo5s = append(clo5s, v[0])
								}
								res = append(res, c)
							}
						}
					}
				}
			}
		}
	}
	Info(res, clo5s)
	return
}

func GetClothesPGifts(req *http.Request) (res map[string]int, clo5s map[string]int, err error) {
	res = make(map[string]int)
	clo5s = make(map[string]int)

	for k, v := range req.Form {
		// Info(k, v)
		if strings.HasPrefix(k, CP_NAME_PREFIX) && len(v) > 0 && v[0] != "" {
			index := k[len(CP_NAME_PREFIX):]
			for k1, v1 := range req.Form {
				if strings.HasPrefix(k1, CP_CNT_PREFIX) && len(v1) > 0 && v1[0] != "" && k1[len(CP_CNT_PREFIX):] == index {
					// Info("clothes id:", v[0], "cnt:", v1[0])
					idStr := strings.Split(v[0], "-")[0]
					partStr := strings.Split(v[0], "-")[2]

					var cnt int
					if cnt, err = strconv.Atoi(v1[0]); err != nil {
						Err(err)
						return
					} else if cnt < 0 {
						err = ClothesZeroErr
						return
					} else {
						exist, star := clothes.ClothesLegal(idStr)
						if !exist {
							err = ClothesIllegalErr
							return
						} else {
							key := fmt.Sprintf("%s:%s", idStr, partStr)
							res[key] = cnt
							if star == 5 && cnt > 10 {
								clo5s[v[0]] = cnt
							}
						}
					}
				}
			}
		}
	}
	Info(res, clo5s)
	return
}

func GetCosplayStrInRedis(cosplatID int64, title string, keyword string, typee int, openTime int64, closeTime int64, adminUsername string, adminNickname string, icon string, cosbg string, listbg string, params CosParams) string {
	c := GetCosplayInRedis(cosplatID, title, keyword, typee, openTime, closeTime, adminUsername, adminNickname, icon, cosbg, listbg, params)
	if cb, err := json.Marshal(c); err != nil {
		Err(err)
		return ""
	} else {
		return "'" + string(cb) + "'"
	}
}

//IsTimeZone 判断是否在输入的时间段
func IsTimeZone(nowtime int64, startTime, endTime string) bool {
	timeSta, err := time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		Err(err)
		return false
	}
	timeEnd, err := time.Parse("2006-01-02 15:04:05", endTime)
	if err != nil {
		Err(err)
		return false
	}
	if nowtime > timeSta.Unix() && nowtime < timeEnd.Unix() {
		return true
	} else {
		return false
	}

}
func GetCosplayInRedis(cosplatID int64, title string, keyword string, typee int, openTime int64, closeTime int64, adminUsername string, adminNickname string, icon string, cosbg string, listbg string, params CosParams) (res Cosplay) {
	res.CosplayID = cosplatID
	res.Title = title
	res.Keyword = keyword
	res.Type = typee
	res.OpenTime = openTime
	res.CloseTime = closeTime
	res.AdminUsername = adminUsername
	res.AdminNickname = adminNickname
	res.Icon = icon
	res.CosBg = cosbg
	res.ListBg = listbg
	res.Params = params

	return
}

func GetHelpInRedis(username string, nickname string, lid string, content string, postTime int64, image string, bc string, gc string) (res Help) {
	// res.HelpID = helpID
	res.Username = username
	res.Nickname = nickname
	res.LevelID = lid
	res.Content = content
	res.PostTime = postTime
	res.Image = image
	res.BoyClothes = bc
	res.GirlClothes = gc

	return
}

func GetHelpStrInRedis(h *Help) (res string) {
	s := GetHelpInRedis(h.Username, h.Nickname, h.LevelID, h.Content, h.PostTime, h.Image, h.BoyClothes, h.GirlClothes)
	s.HelpID = h.HelpID
	if sb, err := json.Marshal(s); err != nil {
		Err(err)
		return
	} else {
		res = "'" + string(sb) + "'"
	}
	return
}

func GetGapDay(st string, et string) int {
	stime, _ := time.ParseInLocation("2006-01-02 15:04:05", st, time.Local)
	etime, _ := time.ParseInLocation("2006-01-02 15:04:05", et, time.Local)

	sd := time.Date(stime.Year(), stime.Month(), stime.Day(), 0, 0, 0, 0, time.Local)
	ed := time.Date(etime.Year(), etime.Month(), etime.Day(), 0, 0, 0, 0, time.Local)

	return int(ed.Sub(sd).Hours() / 24)
}

func GetEventBaseProgress(st int64, daycnt int) int {
	now := time.Now()
	stime := time.Unix(st, 0)
	loc, _ := time.LoadLocation("Local")
	nd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	sd := time.Date(stime.Year(), stime.Month(), stime.Day(), 0, 0, 0, 0, loc)

	gapDay := int(nd.Sub(sd).Hours() / 24)

	return gapDay * daycnt
}

func PrettyJson(b string) string {
	var out bytes.Buffer
	if err := json.Indent(&out, []byte(b), "", "\t"); err != nil {
		Err(err)
		return b
	}
	pre := string(out.Bytes())
	pre = strings.Replace(pre, "\n", "<br>", -1)
	pre = strings.Replace(pre, "\t", "      ", -1)
	return pre
}

func SimplePrettyJson(b string) string {
	var out bytes.Buffer
	if err := json.Indent(&out, []byte(b), "", "\t"); err != nil {
		Err(err)
		return b
	}
	pre := string(out.Bytes())
	// pre = strings.Replace(pre, "\n", "<br>", -1)
	// pre = strings.Replace(pre, "\t", "      ", -1)
	return pre
}

// func GetQueueCnt() int {
// 	mServer := appcfg.GetString("m_server_address", "192.168.1.253:29970/GFMServer/")

// 	postUrl := "http://" + mServer + "getQueueCnt"

// 	i := 0
// 	for {
// 		resp, err := http.Get(postUrl)
// 		if err != nil {
// 			//may need retry
// 			if i >= 3 {
// 				//do not retry again
// 				return 0
// 			} else {
// 				//retry
// 				i++
// 			}
// 		} else {
// 			//success
// 			defer resp.Body.Close()
// 			bodyByte, readError := ioutil.ReadAll(resp.Body)
// 			if readError != nil {
// 				Err(readError)
// 				return 0
// 			} else {
// 				// Info(string(bodyByte))
// 				var responseJson ResponseJson
// 				err = json.Unmarshal(bodyByte, &responseJson)
// 				if err != nil {
// 					Err(err)
// 					return 0
// 				} else {
// 					return responseJson.Count
// 				}
// 			}
// 		}
// 	}
// }

func Execute(cmd string, args ...string) (res string, err error) {
	Info(cmd, args)

	var cmmd *exec.Cmd
	cmmd = exec.Command(cmd, args...)
	errOutput := bytes.NewBuffer(nil)
	cmmd.Stderr = errOutput

	resOutput := bytes.NewBuffer(nil)
	cmmd.Stdout = resOutput

	if err = cmmd.Run(); err != nil {
		Err(err, string(errOutput.Bytes()))
		return
	} else {
		res = string(resOutput.Bytes())
	}

	return
}

func TrimRedisList(key string, start int, end int) {
	Info("trim list:", key, "start:", start, "end:", end)
	conn := rao.GetConn()
	defer conn.Close()

	_, err := conn.Do("LTRIM", key, start, end)
	if err != nil {
		Err(err)
		return
	}
}

func GenMD5Multipart(f multipart.File) (res string, err error) {
	md5h := md5.New()
	if _, err = io.Copy(md5h, f); err != nil {
		Err(err)
		return
	}

	return hex.EncodeToString(md5h.Sum([]byte(""))), err
}

func GenMD5F(filepath string) (res string, err error) {
	var file *os.File
	if file, err = os.Open(filepath); err != nil {
		Err(err)
		return
	} else {
		defer file.Close()
		md5h := md5.New()
		if _, err = io.Copy(md5h, file); err != nil {
			Err(err)
			return
		}
		return hex.EncodeToString(md5h.Sum([]byte(""))), err
	}
}

func GenMD5Raw(content []byte) (res string, err error) {
	md5h := md5.New()
	if _, err = md5h.Write(content); err != nil {
		Err(err)
		return
	}
	res = hex.EncodeToString(md5h.Sum([]byte("")))
	return
}

func GenMD5S(s string) (res string, err error) {
	md5h := md5.New()
	if _, err = md5h.Write([]byte(s)); err != nil {
		Err(err)
		return
	}
	res = hex.EncodeToString(md5h.Sum([]byte("")))
	return
}

func GenMD5SBase64(s string) (res string, err error) {
	md5h := md5.New()
	if _, err = md5h.Write([]byte(s)); err != nil {
		Err(err)
		return
	}
	// res = hex.EncodeToString(md5h.Sum([]byte("")))
	res = base64.StdEncoding.EncodeToString(md5h.Sum(nil))
	return
}

func GenPKLevel(point int) (int, string) {
	if point <= 4 {
		return 1, L("pk13")
	} else if point <= 19 {
		return 2, L("pk14")
	} else if point <= 49 {
		return 3, L("pk15")
	} else if point <= 99 {
		return 4, L("pk16")
	} else if point <= 199 {
		return 5, L("pk17")
	} else if point <= 399 {
		return 6, L("pk18")
	} else if point <= 799 {
		return 7, L("pk19")
	} else if point <= 1199 {
		return 8, L("pk20")
	} else if point <= 1499 {
		return 9, L("pk21")
	} else {
		return 10, L("pk22")
	}
}

func ComputePKPointRange(point int) (start int, end int) {
	l, _ := GenPKLevel(point)
	if l == 1 {
		return -999999, 4
	} else if l == 2 {
		return 5, 19
	} else if l == 3 {
		return 20, 49
	} else if l == 4 {
		return 50, 99
	} else if l == 5 {
		return 100, 199
	} else if l == 6 {
		return 200, 399
	} else if l == 7 {
		return 400, 799
	} else if l == 8 {
		return 800, 1199
	} else if l == 9 {
		return 1200, 1499
	} else if l == 10 {
		return 1500, 999999
	}
	return -999999, 999999
}

func GenDesignerRate(points int) float64 {
	var rate float64
	if points < constants.DESIGNER_LEVEL_2_POINT {
		rate = 0.02
	} else if points < constants.DESIGNER_LEVEL_3_POINT {
		rate = 0.05
	} else if points < constants.DESIGNER_LEVEL_4_POINT {
		rate = 0.10
	} else if points < constants.DESIGNER_LEVEL_5_POINT {
		rate = 0.15
	} else if points < constants.DESIGNER_LEVEL_6_POINT {
		rate = 0.20
	} else {
		rate = 0.25
	}

	return rate
}

func GenDesignerLevel(points int) int {
	if points < constants.DESIGNER_LEVEL_2_POINT {
		return 1
	} else if points < constants.DESIGNER_LEVEL_3_POINT {
		return 2
	} else if points < constants.DESIGNER_LEVEL_4_POINT {
		return 3
	} else if points < constants.DESIGNER_LEVEL_5_POINT {
		return 4
	} else if points < constants.DESIGNER_LEVEL_6_POINT {
		return 5
	} else {
		return 6
	}
}

func GetOneCloPieceReward(part1Cnt int, part2Cnt int, part3Cnt int, part4Cnt int) (part int) {
	//set up min cnt part set
	cntSet := []int{part1Cnt, part2Cnt, part3Cnt, part4Cnt}
	minPartSet := []int{1}

	if part2Cnt < cntSet[minPartSet[0]-1] {
		minPartSet = []int{2}
	} else if part2Cnt == cntSet[minPartSet[0]-1] {
		minPartSet = append(minPartSet, 2)
	}

	if part3Cnt < cntSet[minPartSet[0]-1] {
		minPartSet = []int{3}
	} else if part3Cnt == cntSet[minPartSet[0]-1] {
		minPartSet = append(minPartSet, 3)
	}

	if part4Cnt < cntSet[minPartSet[0]-1] {
		minPartSet = []int{4}
	} else if part4Cnt == cntSet[minPartSet[0]-1] {
		minPartSet = append(minPartSet, 4)
	}

	//chooes which part is 10%--change to 19%
	randNum := rand.Intn(len(minPartSet))
	part10 := minPartSet[randNum]
	var part30 []int
	if part10 == 1 {
		part30 = []int{2, 3, 4}
	} else if part10 == 2 {
		part30 = []int{1, 3, 4}
	} else if part10 == 3 {
		part30 = []int{1, 2, 4}
	} else {
		part30 = []int{1, 2, 3}
	}

	randNum = rand.Intn(100)

	if randNum <= 18 {
		//10%--change to 19
		return part10
	} else {
		randNum = rand.Intn(3)
		return part30[randNum]
	}
}

//CalculateDays 计算天数
func CalculateDays(year int, month int) (days int) {
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30

		} else {
			days = 31
			fmt.Fprintln(os.Stdout, "The month has 31 days")
		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			days = 29
		} else {
			days = 28
		}
	}
	return
}
