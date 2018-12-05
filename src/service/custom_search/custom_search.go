package custom_search

import (
	"appcfg"
	"constants"
	"context"
	// "database/sql"
	"elastic"
	"fmt"
	. "logger"
	// "math/rand"
	"reflect"
	"strings"
	"time"
	. "types"
)

const (
	ELASTIC_INIT_BATCH_SIZE = 1000
	DESIGN_INDEX_NAME       = "design"
	CUSTOM_INDEX_TYPE_NAME  = "custom"
)

var ctx context.Context
var client *elastic.Client

func init() {
	//init
	ctx = context.Background()
	if appcfg.GetBool("support_elastic", true) {
		if appcfg.GetServerType() == "" || appcfg.GetServerType() == constants.SERVER_TYPE_GV {
			var err error

			urls := appcfg.GetString("elastic_urls", "")
			urlList := strings.Split(urls, ";")
			if len(urls) <= 0 {
				panic("no elastic url config")
			}

			urlConfig := elastic.SetURL(urlList...)
			if client, err = elastic.NewClient(urlConfig); err != nil {
				panic(err)
			}

			// if appcfg.GetBool("main_server", false) && appcfg.GetBool("init_custom_elastic", false) {
			// Info("init custom")
			// initCustomElastic()
			// }
		}
	}
}

// func InitAndTestCustomFrom(dbAddress, dbUsername, dbPasswd, dbName, tblName, elasticURLs string) {
// 	urlList := strings.Split(elasticURLs, ";")
// 	urlConfig := elastic.SetURL(urlList...)
// 	if iclient, err := elastic.NewClient(urlConfig); err != nil {
// 		panic(err)
// 	} else {
// 		//delete index first
// 		if _, err := iclient.DeleteIndex(DESIGN_INDEX_NAME).Do(ctx); err != nil {
// 			Err(err)
// 		}

// 		var err error
// 		var db *sql.DB
// 		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?readTimeout=15s&writeTimeout=15s&timeout=10s&charset=utf8mb4", dbUsername, dbPasswd, dbAddress, dbName))
// 		if err != nil {
// 			panic(err)
// 		}

// 		var ElasticListCustomStmt *sql.Stmt
// 		if ElasticListCustomStmt, err = db.Prepare("select cid,username,nickname,model_no,type,cname,`desc`,custom_file,status,upload_time,check_time,m_type,m_amount,hearts,buy_count,custom_extra,inventory,tag1,tag2 from " + tblName + " where status in (3,6) order by cid limit ?,?"); err != nil {
// 			panic(err)
// 		}

// 		start := 0

// 		for {
// 			Info("init custom search query:", start, ELASTIC_INIT_BATCH_SIZE)
// 			if res, err := custom.ElasticListCustom1(ElasticListCustomStmt, start, ELASTIC_INIT_BATCH_SIZE); err != nil {
// 				panic(err)
// 			} else {
// 				if len(res) <= 0 {
// 					break
// 				}

// 				for _, c := range res {
// 					// Info("init custom:", c.CID)
// 					if _, err = iclient.Index().
// 						Index(DESIGN_INDEX_NAME).
// 						Type(CUSTOM_INDEX_TYPE_NAME).
// 						Id(fmt.Sprintf("%d", c.CID)).
// 						BodyJson(c).
// 						Refresh("true").
// 						Do(ctx); err != nil {
// 						panic(err)
// 					}
// 				}
// 			}

// 			start += ELASTIC_INIT_BATCH_SIZE
// 		}
// 		Info("init done")

// 		types := []int{100, 200, 300, 400, 500, 600, 701, 702, 703, 704, 705, 706, 707, 708, 709, 800}
// 		models := []int{101, 103, 201, 202}
// 		testSize := 1000

// 		mt := 0.0
// 		for i := 0; i < testSize; i++ {
// 			// v := rand.Intn(40000)
// 			Info("search", i)
// 			t := types[rand.Intn(len(types))]
// 			m := models[rand.Intn(len(models))]

// 			st := time.Now().UnixNano()
// 			query := genQuery("", t, m, 0, false, 0, 0, 0, "")
// 			searchCustom("", iclient, query, 1, 15)
// 			st1 := time.Now().UnixNano()

// 			mt += float64(st1-st) / float64(time.Millisecond)
// 		}

// 		Info("avg search time in milli second:", mt/float64(testSize))
// 	}
// }

// func initCustomElastic() {
// 	//delete index first
// 	if _, err := client.DeleteIndex(DESIGN_INDEX_NAME).Do(ctx); err != nil {
// 		Err(err)
// 	}

// 	//create index next
// 	if _, err := client.CreateIndex(DESIGN_INDEX_NAME).Do(ctx); err != nil {
// 		Err(err)
// 	}

// 	start := 0

// 	for {
// 		Info("init custom search query:", start, ELASTIC_INIT_BATCH_SIZE)
// 		if res, err := custom.ElasticListCustom(start, ELASTIC_INIT_BATCH_SIZE); err != nil {
// 			panic(err)
// 		} else {
// 			if len(res) <= 0 {
// 				break
// 			}

// 			for _, c := range res {
// 				// Info("init custom:", c.CID)
// 				if _, err = client.Index().
// 					Index(DESIGN_INDEX_NAME).
// 					Type(CUSTOM_INDEX_TYPE_NAME).
// 					Id(fmt.Sprintf("%d", c.CID)).
// 					BodyJson(c).
// 					Refresh("true").
// 					Do(ctx); err != nil {
// 					panic(err)
// 				}
// 			}
// 		}

// 		start += ELASTIC_INIT_BATCH_SIZE
// 	}
// 	Info("init done")
// }

func UpCustomSearch(cid int64, mt, ma, inventory int) {
	Info("update custom search:", cid, mt, ma, inventory)
	if ur, err := client.Update().Index(DESIGN_INDEX_NAME).Type(CUSTOM_INDEX_TYPE_NAME).Id(fmt.Sprintf("%d", cid)).
		Script(elastic.NewScriptInline(fmt.Sprintf("ctx._source.status = %d;ctx._source.mType = %d;ctx._source.mAmount = %d;ctx._source.inventory = %d;ctx._source.upTime = %d", PASSED, mt, ma, inventory, time.Now().Unix()))).
		Do(ctx); err != nil {
		Err(err)
	} else {
		Info("update result index:", ur.Index, "id:", ur.Id, "version:", ur.Version, "result:", ur.Result, "total:", ur.Shards.Total, "success:", ur.Shards.Successful, "fail:", ur.Shards.Failed)
	}
}

func SellOneCustomSearch(cid int64) {
	Info("sell one custom search:", cid)
	if _, err := client.Update().Index(DESIGN_INDEX_NAME).Type(CUSTOM_INDEX_TYPE_NAME).Id(fmt.Sprintf("%d", cid)).
		Script(elastic.NewScriptInline("ctx._source.inventory -= 1")).
		Do(ctx); err != nil {
		Err(err)
	}
}

func SellSomeCustomSearch(cid int64, buyCount int) {
	Info("sell one custom search:", cid)
	if _, err := client.Update().Index(DESIGN_INDEX_NAME).Type(CUSTOM_INDEX_TYPE_NAME).Id(fmt.Sprintf("%d", cid)).
		Script(elastic.NewScriptInline("ctx._source.inventory -= " + fmt.Sprintf("%d", buyCount))).
		Do(ctx); err != nil {
		Err(err)
	}
}

func DownCustomSearch(cid int64) {
	Info("down custom search:", cid)
	if _, err := client.Update().Index(DESIGN_INDEX_NAME).Type(CUSTOM_INDEX_TYPE_NAME).Id(fmt.Sprintf("%d", cid)).
		Script(elastic.NewScriptInline(fmt.Sprintf("ctx._source.status = %d;ctx._source.inventory = 0", CUSTOM_READY))).
		Do(ctx); err != nil {
		Err(err)
	}
}

func BanCustomSearch(cid int64) {
	Info("down custom search:", cid)
	if _, err := client.Update().Index(DESIGN_INDEX_NAME).Type(CUSTOM_INDEX_TYPE_NAME).Id(fmt.Sprintf("%d", cid)).
		Script(elastic.NewScriptInline(fmt.Sprintf("ctx._source.status = %d", CUSTOM_BAN))).
		Do(ctx); err != nil {
		ErrMail(err)
	}
}

func HuifuCustomSearch(cid int64, mt, ma int) {
	Info("huifu custom search:", cid)
	if _, err := client.Update().Index(DESIGN_INDEX_NAME).Type(CUSTOM_INDEX_TYPE_NAME).Id(fmt.Sprintf("%d", cid)).
		Script(elastic.NewScriptInline(fmt.Sprintf("ctx._source.status = %d;ctx._source.mType = %d;ctx._source.mAmount = %d", PASSED, mt, ma))).
		Do(ctx); err != nil {
		ErrMail(err)
	}
}

func InitOneCustomSearch(c ElasticCustom) {
	if _, err := client.Index().
		Index(DESIGN_INDEX_NAME).
		Type(CUSTOM_INDEX_TYPE_NAME).
		Id(fmt.Sprintf("%d", c.CID)).
		BodyJson(c).
		Refresh("true").
		Do(ctx); err != nil {
		ErrMail(err)
	}
}

func RefreshCustomSearch(c1 *Custom, mt, ma, inventory int, toStatus CustomStatus) {
	c := ElasticCustom{
		CID:             c1.CID,
		ModelNo:         c1.ModelNo,
		Type:            c1.Type,
		CName:           c1.CName,
		Desc:            c1.Desc,
		CloSet:          c1.CloSet,
		Status:          toStatus,
		UploadTime:      c1.UploadTime,
		CheckTime:       c1.CheckTime,
		MoneyType:       mt,
		MoneyAmount:     ma,
		Hearts:          c1.Hearts,
		Username:        c1.Username,
		Nickname:        c1.Nickname,
		BuyCount:        c1.BuyCount,
		Inventory:       inventory,
		Extra:           c1.Extra,
		Tag1:            c1.Tag1,
		Tag2:            c1.Tag2,
		ClothesIDInGame: c1.ClothesIDInGame,

		UpTime: time.Now().Unix(),
	}

	if _, err := client.Index().
		Index(DESIGN_INDEX_NAME).
		Type(CUSTOM_INDEX_TYPE_NAME).
		Id(fmt.Sprintf("%d", c.CID)).
		BodyJson(c).
		Refresh("true").
		Do(ctx); err != nil {
		ErrMail(err)
	}
}

func RefreshCustomSearchStickTime(c1 *Custom, st int64) {
	c := ElasticCustom{
		CID:             c1.CID,
		ModelNo:         c1.ModelNo,
		Type:            c1.Type,
		CName:           c1.CName,
		Desc:            c1.Desc,
		CloSet:          c1.CloSet,
		Status:          c1.Status,
		UploadTime:      c1.UploadTime,
		CheckTime:       c1.CheckTime,
		MoneyType:       c1.MoneyType,
		MoneyAmount:     c1.MoneyAmount,
		Hearts:          c1.Hearts,
		Username:        c1.Username,
		Nickname:        c1.Nickname,
		BuyCount:        c1.BuyCount,
		Inventory:       c1.Inventory,
		Extra:           c1.Extra,
		Tag1:            c1.Tag1,
		Tag2:            c1.Tag2,
		ClothesIDInGame: c1.ClothesIDInGame,

		UpTime: st,
	}

	if _, err := client.Index().
		Index(DESIGN_INDEX_NAME).
		Type(CUSTOM_INDEX_TYPE_NAME).
		Id(fmt.Sprintf("%d", c.CID)).
		BodyJson(c).
		Refresh("true").
		Do(ctx); err != nil {
		ErrMail(err)
	}
}

func DelCustomSearch(cid int64) {
	Info("del custom search:", cid)
	if res, err := client.Delete().
		Index(DESIGN_INDEX_NAME).
		Type(CUSTOM_INDEX_TYPE_NAME).
		Id(fmt.Sprintf("%d", cid)).
		Do(ctx); err != nil {
		// Info(err.Error())
		// Info(strings.Contains(err.Error(), constants.ELASTIC_NOT_FOUND_ERR_MSG))
		if strings.Contains(err.Error(), constants.ELASTIC_NOT_FOUND_ERR_MSG) {
			Info("del custom in search not exist:", cid)
		} else {
			Err(err)
		}
	} else if res.Found {
		Info("Custom deleted from search:", cid)
	} else {
		Info("Custom delete but not in search:", cid)
	}
}

func SearchCustom(keyword string, cloType, modelNo, page, pageSize int, moneyType int, ignoreNoInventory bool, qTag1, qTag2, startStar int, designerUsername string) (res []Custom, err error) {
	start := (page - 1) * pageSize

	query := genQuery(keyword, cloType, modelNo, moneyType, ignoreNoInventory, qTag1, qTag2, startStar, designerUsername)

	return searchCustom(keyword, client, query, start, pageSize)
}

func genQuery(keyword string, cloType, modelNo, moneyType int, ignoreNoInventory bool, qTag1, qTag2, startStar int, designerUsername string) (query *elastic.BoolQuery) {
	query = elastic.NewBoolQuery()

	//should
	if qTag1 > 0 || qTag2 > 0 {
		tags := make([]interface{}, 0)

		if qTag1 > 0 {
			tags = append(tags, qTag1)
		}

		if qTag2 > 0 {
			tags = append(tags, qTag2)
		}

		tagQuery := elastic.NewBoolQuery()
		tagQuery.Should(elastic.NewTermsQuery("Tag1", tags...), elastic.NewTermsQuery("Tag2", tags...))
		tagQuery.MinimumNumberShouldMatch(1)
		query.Filter(tagQuery)
	}

	statusQuery := elastic.NewBoolQuery()
	statusQuery.Should(elastic.NewTermQuery("status", fmt.Sprintf("%d", PASSED)), elastic.NewTermQuery("status", fmt.Sprintf("%d", PASSED_NO_INVENTORY)))
	statusQuery.MinimumNumberShouldMatch(1)
	query.Filter(statusQuery)

	//filter last
	query.Filter(elastic.NewTermQuery("type", fmt.Sprintf("%d", cloType)))
	query.Filter(elastic.NewTermQuery("modelNo", fmt.Sprintf("%d", modelNo)))

	if designerUsername != "" {
		query.Filter(elastic.NewTermQuery("username", designerUsername))
	}

	if moneyType != 0 {
		query.Filter(elastic.NewTermQuery("mType", fmt.Sprintf("%d", moneyType)))
	}

	if ignoreNoInventory {
		inventoryQuery := elastic.NewBoolQuery()
		inventoryQuery.Should(elastic.NewRangeQuery("inventory").Gt(0), elastic.NewTermQuery("status", fmt.Sprintf("%d", PASSED_NO_INVENTORY)))
		inventoryQuery.MinimumNumberShouldMatch(1)
		query.Filter(inventoryQuery)
	}

	if startStar > 0 {
		query.Filter(elastic.NewRangeQuery("hearts").Gte(startStar))
	}

	//must
	if keyword != "" {
		query.Must(elastic.NewMultiMatchQuery(keyword, "cname", "desc"))
	}

	return
}

func searchCustom(keyword string, cli *elastic.Client, query *elastic.BoolQuery, start, pageSize int) (res []Custom, err error) {
	var searchResult *elastic.SearchResult
	s := cli.Search().Index(DESIGN_INDEX_NAME).Query(query).From(start).Size(pageSize)

	if keyword == "" {
		s.Sort("upTime", false)
	}

	if searchResult, err = s.Do(ctx); err != nil {
		Err(err)
		return
	}

	// Info("Default query custom took", searchResult.TookInMillis, "milliseconds", "hits:", searchResult.TotalHits())

	res = make([]Custom, 0)
	var ttyp Custom
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Custom); ok {
			res = append(res, t)
		}
	}

	return
}
