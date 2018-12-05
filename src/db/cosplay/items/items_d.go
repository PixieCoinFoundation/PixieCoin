package items

import (
	"dao_party"
	"database/sql"
	"encoding/json"
	. "logger"
)

import (
	. "types"
)

func countPlayerCos(cosplayID int64) (cnt int, err error) {
	if err = dao_party.CountPlayerCosStmt.QueryRow(cosplayID).Scan(&cnt); err != nil {
		Err(err)
	}
	return
}

func addScoreToDB(itemId int64, username string, score int, totalScore int64) (err error) {
	if _, err = dao_party.AddCosItemScoreStmt.Exec(itemId, username, score); err != nil {
		Err(err)
		return
	}

	if _, err = dao_party.SetCosItemScoreStmt.Exec(totalScore, itemId); err != nil {
		Err(err)
		return
	}

	return
}

func getScoresFromDB(itemID int64) (res []CosItemScore, err error) {
	var rows *sql.Rows
	if rows, err = dao_party.GetCosItemScoreStmt.Query(itemID); err != nil {
		Err(err)
		return
	}
	defer rows.Close()

	res = make([]CosItemScore, 0)
	for rows.Next() {
		c := CosItemScore{
			CosItemID: itemID,
		}
		if err = rows.Scan(&(c.Username), &(c.Score)); err != nil {
			Err(err)
			return
		}

		res = append(res, c)
	}

	return
}

func addCosItemToDB(ci *CosItem) (id int64, err error) {
	var clob []byte
	if clob, err = json.Marshal(ci.Clothes); err != nil {
		Err(err)
		return
	}

	var res sql.Result
	if res, err = dao_party.AddCosItemStmt.Exec(ci.Username, ci.UID, ci.Nickname, ci.CosplayID, ci.Title, ci.Desc, ci.ModelNo, string(clob), ci.Score, ci.SysScore, ci.Icon, ci.Head, ci.CosBg, ci.CosPoints, ci.UploadTime, ci.SetTopTime, ci.Rank, ci.Img); err != nil {
		Err(err)
		return
	} else if id, err = res.LastInsertId(); err != nil {
		Err(err)
		return
	}

	return
}

func DeleteCosItem(itemID int64) (err error) {
	if _, err = dao_party.DeleteCosItemStmt.Exec(itemID); err != nil {
		Err(err)
	}
	return
}

func GetCosItemByCosplay(cosplayID int64, start int, size int) (res []FCosItem, err error) {
	var rows *sql.Rows
	if rows, err = dao_party.GetCosItemByCosplayStmt.Query(cosplayID, start, size); err != nil {
		Err(err)
		return
	}
	defer rows.Close()

	res = make([]FCosItem, 0)
	for rows.Next() {
		f := FCosItem{}

		if err = rows.Scan(&f.ItemID, &f.Username, &f.Icon, &f.Img); err != nil {
			Err(err)
			return
		}

		res = append(res, f)
	}

	return
}
