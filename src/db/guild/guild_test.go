package guild

import (
	"appcfg"
	"common"
	"constants"
	"dao"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	. "logger"
	"math/rand"
	"rao"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	. "types"
)

func init() {
	if appcfg.GetServerType() != "" {
		return
	}

	// if appcfg.GetBool("local_test", false) {
	// initTestData()
	// }

	// for i := 0; i < 50; i++ {
	// 	testMatch(testCnt)
	// }
}

func initTestData() {
	testCnt := 100

	dao.DefaultDB.Exec("truncate gf_guild")
	dao.DefaultDB.Exec("truncate gf_guild_member")
	dao.DefaultDB.Exec("truncate gf_guild_war_log")
	dao.DefaultDB.Exec("truncate gf_guild_clothes")
	dao.DefaultDB.Exec("truncate gf_apply_guild")
	_, gwt := common.GetGuildYearWeekToken()
	for i := 0; i < testCnt; i++ {
		is := fmt.Sprintf("%d", i)
		Add("a"+is, "an"+is, "aname"+is, "adesc"+is, 1, "head-101.png", 0, 1, 0, gwt)
	}

	conn := rao.GetConn()
	defer conn.Close()
	wt, gwt := common.GetGuildYearWeekToken()

	conn.Do("DEL", rao.GetGuildMedalGWeekKey(gwt))
	conn.Do("DEL", rao.GetGuildActivityWeekKey(wt))
}

func testMatch(testCnt int) {
	dao.DefaultDB.Exec("update gf_guild set war_date='',war_token=''")
	dao.DefaultDB.Exec("update gf_guild_member set war_date=''")

	pks := []string{"pk_team0001", "pk_team0001", "pk_team0001", "pk_team0001", "pk_team0001", "pk_team0001", "pk_team0001", "pk_team0001", "pk_team0001", "pk_team0001",
		"pk_team0001", "pk_team0001", "pk_team0001", "pk_team0001", "pk_team0001"}

	var wg sync.WaitGroup
	var fc, ufc int32

	var mlock sync.Mutex
	var ut int64
	var running int32
	m := make(map[int64]int)
	_, gwt := common.GetGuildYearWeekToken()
	for i := 1; i < testCnt+1; i++ {
		wg.Add(1)
		atomic.AddInt32(&running, 1)

		for atomic.LoadInt32(&running) >= int32(50) {
			time.Sleep(10 * time.Millisecond)
		}

		go func(n int64, f *int32, uf *int32) {
			defer atomic.AddInt32(&running, -1)
			defer wg.Done()

			st := time.Now().UnixNano()
			o, e := matchGuildWarOnce(n, 0, "20170531", pks, gwt)
			et := time.Now().UnixNano()

			if e == nil {
				atomic.AddInt32(f, 1)
				// Info("f:", n, "o:", o.ID)
				mlock.Lock()
				ut += et - st
				m[n] = 1
				m[o.ID] = 1
				mlock.Unlock()
			} else {
				atomic.AddInt32(uf, 1)
				// Info("n:", n)
			}
		}(int64(i), &fc, &ufc)
	}
	wg.Wait()
	Info("fc:", fc, "ufc:", ufc, "len map:", len(m), "avg time:", float64(ut)/float64(testCnt)/float64(time.Millisecond))

	if fc+ufc != int32(testCnt) {
		panic("fc+ufc!=testCnt")
	}

	if fc > int32(testCnt/2) {
		panic("fc>testCnt/2")
	}

	if fc*2 != int32(len(m)) {
		panic("fc*2!=len(m)")
	}
}
