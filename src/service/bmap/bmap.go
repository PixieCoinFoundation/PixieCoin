package bmap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "logger"
	"math/rand"
	// "net"
	"net/http"
	"strconv"
	"strings"
	// "time"
	"tools"
	"xlsx"
)

var url string = "http://api.map.baidu.com/place/v2/search?page_size=20&scope=1&&output=json&ak=bZ1adGcGHQ2v50VyKWXMgEroSXSGrHrV"

//bounds=39.915,116.404,39.975,116.414
//query=美食
//page_num=0

type BMResult struct {
	Results []BMResult1
}

type BMResult1 struct {
	Name     string
	Location LI
}

type LI struct {
	Lat float64
	Lng float64
}

func GenExcel(code string, config map[string][]string, startPage, endPage int, lat1, lng1, lat2, lng2 float64, hp, zp int) (fileName string) {
	fileName = fmt.Sprintf("dev/%s.xlsx", code)
	file := xlsx.NewFile()
	sheet1, _ := file.AddSheet("Sheet1")
	sheet2, _ := file.AddSheet("Sheet2")

	apicnt := 0
	incluCnt := 0
	for key, exu := range config {
		turl := url + "&query=" + key
		zincr := (lat2 - lat1) / float64(zp)
		hincr := (lng2 - lng1) / float64(hp)

		for z := 0; z < zp; z++ {
			for h := 0; h < hp; h++ {
				lat1Now := lat1 + float64(z)*zincr
				lng1Now := lng1 + float64(h)*hincr

				lat2Now := lat1 + float64(z+1)*zincr
				lng2Now := lng1 + float64(h+1)*hincr

				for i := startPage; i <= endPage; i++ {
					tturl := turl + "&page_num=" + strconv.Itoa(i) + fmt.Sprintf("&bounds=%f,%f,%f,%f", lat1Now, lng1Now, lat2Now, lng2Now)
					Info(tturl)
					apicnt++
					req, reqErr := http.NewRequest("GET", tturl, nil)
					if reqErr != nil {
						Err(reqErr)
						continue
					}
					var err error
					var resp *http.Response

					// DefaultClient := http.Client{
					// 	Transport: &http.Transport{
					// 		Dial: func(netw, addr string) (net.Conn, error) {
					// 			deadline := time.Now().Add(60 * time.Second)
					// 			c, eerr := net.DialTimeout(netw, addr, time.Second*60)
					// 			if eerr != nil {
					// 				Err(eerr)
					// 				return nil, err
					// 			}
					// 			c.SetDeadline(deadline)
					// 			return c, nil
					// 		},
					// 		DisableKeepAlives: true,
					// 	},
					// }
					resp, err = tools.OtakuHttpClient.Do(req)
					if err != nil {
						Err(err)
						continue
					}

					defer func() {
						resp.Body.Close()
					}()

					var ret []byte
					ret, err = ioutil.ReadAll(resp.Body)
					if err != nil {
						Err(err)
						continue
					}
					// Info(string(ret))

					var res BMResult
					json.Unmarshal(ret, &res)

					if len(res.Results) <= 0 {
						break
					}
					for _, v := range res.Results {
						if !nameLegal(key, v.Name, exu) {
							Info("exclude:", key, v.Name)
							row := sheet2.AddRow()

							cell1 := row.AddCell()
							cell1.Value = v.Name
							cell2 := row.AddCell()
							cell2.Value = fmt.Sprintf("%f", v.Location.Lat)
							cell3 := row.AddCell()
							cell3.Value = fmt.Sprintf("%f", v.Location.Lng)
						} else {
							incluCnt++
							Info("include:", key, v.Name)

							t, p, f, rp := getPropDetail(key)

							row := sheet1.AddRow()

							cell11 := row.AddCell()
							cell11.Value = fmt.Sprintf("%s-%d", code, incluCnt)

							cell1 := row.AddCell()
							cell1.Value = v.Name

							cell12 := row.AddCell()
							cell12.Value = t

							cell2 := row.AddCell()
							cell2.Value = fmt.Sprintf("%f", v.Location.Lat)

							cell3 := row.AddCell()
							cell3.Value = fmt.Sprintf("%f", v.Location.Lng)

							cell4 := row.AddCell()
							cell4.Value = fmt.Sprintf("%d", p)

							cell5 := row.AddCell()
							cell5.Value = fmt.Sprintf("%f", f)

							cell6 := row.AddCell()
							cell6.Value = fmt.Sprintf("%d", rp)
						}

					}

				}
			}
		}
	}
	Info("total query cnt:", apicnt)
	file.Save(fileName)

	return
}

func nameLegal(key, name string, exus []string) bool {
	for _, exu := range exus {
		if strings.Contains(name, exu) {
			return false
		}
	}
	return true
}

func getPropDetail(key string) (typee string, price int64, factor float64, realPrice int64) {
	factor = float64(93+rand.Intn(15)) / 100.0
	if key == "写字楼" || key == "商厦" || key == "商业中心" {
		r := rand.Intn(100)
		if r <= 49 {
			typee = "普通写字楼"
			price = int64(25000000)
			realPrice = int64(float64(price) * factor)
		} else {
			typee = "高档写字楼"
			price = int64(90000000)
			realPrice = int64(float64(price) * factor)
		}
	} else if key == "公寓" || key == "住宅小区" {
		r := rand.Intn(100)
		if r <= 49 {
			typee = "普通民居"
			price = int64(2500000)
			realPrice = int64(float64(price) * factor)
		} else {
			typee = "高档民居"
			price = int64(8000000)
			realPrice = int64(float64(price) * factor)
		}
	} else if key == "商场" {
		r := rand.Intn(100)
		if r <= 49 {
			typee = "普通商场"
			price = int64(60000000)
			realPrice = int64(float64(price) * factor)
		} else {
			typee = "高档商场"
			price = int64(120000000)
			realPrice = int64(float64(price) * factor)
		}
	} else if key == "公园" {
		r := rand.Intn(100)
		if r <= 49 {
			typee = "普通公园"
			price = int64(50000000)
			realPrice = int64(float64(price) * factor)
		} else {
			typee = "高档公园"
			price = int64(100000000)
			realPrice = int64(float64(price) * factor)
		}
	}
	return
}
