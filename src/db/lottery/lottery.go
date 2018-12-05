package lottery

import (
	// "appcfg"
	"dao"
	"database/sql"
	"fmt"
	. "logger"
	"service/clothes"
	// "time"
	. "types"
)

// func init() {
// 	go func() {
// 		for {
// 			AddLog("我", "c400201000000000101")
// 			AddLog("我啊", "c400201000000000101")
// 			time.Sleep(10 * time.Second)
// 		}
// 	}()

// }

func AddLog(name string, cid string, sname string) (err error) {
	var rn string
	for i, ch := range name {
		if i == 0 {
			rn += fmt.Sprintf("%s", string(ch))
			break
		}
	}
	rn += "****"

	rs := clothes.GetClothesNameAndStar(cid)
	if _, err = dao.AddLotteryLogStmt.Exec(rn, rs[0], rs[1], sname); err != nil {
		Err(err)
	}
	return
}

// func Add10Logs(name string, cs []*ClothesInfo) (err error) {
// 	if len(cs) != 10 {
// 		Err("add 10 lottery log size:", len(cs))
// 		return
// 	}

// 	var rn string
// 	for i, ch := range name {
// 		if i == 0 {
// 			rn += fmt.Sprintf("%s", string(ch))
// 		} else {
// 			rn += "*"
// 		}
// 	}

// 	ps := make([]interface{}, 0)
// 	for _, c := range cs {
// 		ps = append(ps, rn)
// 		rs := clothes.GetClothesNameAndStar(c.ClothesID)
// 		ps = append(ps, rs[0], rs[1])
// 	}
// 	if _, err = dao.AddLotteryLog10Stmt.Exec(ps...); err != nil {
// 		Err(err)
// 	}
// 	return
// }

func QueryLog(limit int) (res []LotteryPubRespData, err error) {
	var rows *sql.Rows
	if rows, err = dao.ListLotteryLog1Stmt.Query(limit); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]LotteryPubRespData, 0)
		for rows.Next() {
			// m := LotteryPubRespData{
			// 	SName: appcfg.GetString("server_group_name", ""),
			// }
			var m LotteryPubRespData
			var id int64

			if err = rows.Scan(&id, &m.RName, &m.Info, &m.Star, &m.SName); err != nil {
				Err(err)
				return
			}
			m.ID = fmt.Sprintf("%d", id)

			res = append(res, m)
		}
	}
	return
}

func QueryLog1(id int64, limit int) (res []LotteryPubRespData, err error) {
	var rows *sql.Rows
	if rows, err = dao.ListLotteryLog2Stmt.Query(id, limit); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]LotteryPubRespData, 0)
		for rows.Next() {
			// m := LotteryPubRespData{
			// 	SName: appcfg.GetString("server_group_name", ""),
			// }
			var m LotteryPubRespData
			var id int64
			if err = rows.Scan(&id, &m.RName, &m.Info, &m.Star, &m.SName); err != nil {
				Err(err)
				return
			}
			m.ID = fmt.Sprintf("%d", id)

			res = append(res, m)
		}
	}
	return
}
