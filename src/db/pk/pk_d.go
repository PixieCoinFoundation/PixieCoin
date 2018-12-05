package pk

// import (
// 	"constants"
// 	"dao"
// 	"database/sql"
// 	. "logger"
// 	. "types"
// )

// func addPKInDB(pk *GFPK) (err error) {
// 	if _, err = dao.AddPKStmt.Exec(pk.Username, pk.LevelID, pk.Nickname, pk.Head, pk.PKPoint, pk.PKLevel, pk.GirlClothesCount, pk.BoyClothesCount, pk.Operations, pk.UploadTime); err != nil {
// 		Err(err)
// 		return
// 	}
// 	return
// }

// func getPKsFromDB(lid string) (res []*GFPK, err error) {
// 	var rows *sql.Rows
// 	if rows, err = dao.GetPKStmt.Query(lid, constants.GF_PK_MAX_SIZE); err != nil {
// 		Err(err)
// 		return
// 	}

// 	res = make([]*GFPK, 0)
// 	for rows.Next() {
// 		p := GFPK{}
// 		var id int64
// 		if err = rows.Scan(&id, &p.Username, &p.LevelID, &p.Nickname, &p.Head, &p.PKPoint, &p.PKLevel, &p.GirlClothesCount, &p.BoyClothesCount, &p.Operations, &p.UploadTime); err != nil {
// 			Err(err)
// 			return
// 		}

// 		res = append(res, &p)
// 	}
// 	return
// }
