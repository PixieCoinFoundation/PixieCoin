package cosplay

import (
	cc "common_db/cosplay"
	"database/sql"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"rao"
)

import (
	"dao_party"
	. "logger"
	. "types"
)

//cosplay with score -1 means closed
func CloseCosplay(cosType int, cosplayId int64) (err error) {
	conn := rao.GetConn()
	defer conn.Close()
	return cc.CloseCosplay(dao_party.UpdateCosplayStatusStmt, cosplayId, &conn)
}

//rconn not nil means recover mode
func Add(cosplay *Cosplay, rconn *redis.Conn) (err error) {
	conn := rao.GetConn()
	defer conn.Close()
	return cc.AddCosplay(dao_party.AddCosplayStmt, cosplay, rconn, &conn)
}

func FindCosplayByStatus(status CosplayStatus) (result []Cosplay, err error) {
	conn := rao.GetConn()
	defer conn.Close()
	return cc.FindCosplayByStatus(status, &conn, dao_party.CheckCosStmt)
}

func GetOldCosplayFromDB(start int, size int) (res []FCosplay, err error) {
	var rows *sql.Rows
	if rows, err = dao_party.GetOldCosplayStmt.Query(COS_CLOSED, start, size); err != nil {
		Err(err)
		return
	}
	defer rows.Close()

	res = make([]FCosplay, 0)
	for rows.Next() {
		f := FCosplay{}

		var paramStr string
		if err = rows.Scan(&f.CosplayID, &f.Title, &f.Keyword, &f.Type, &f.Status, &f.OpenTime, &f.CloseTime, &f.AdminUsername, &f.AdminNickname, &f.Icon, &f.CosBg, &f.ListBg, &paramStr); err != nil {
			Err(err)
			return
		}

		if err = json.Unmarshal([]byte(paramStr), &f.Params); err != nil {
			Err(err)
			return
		}

		res = append(res, f)
	}
	return
}

func DeleteCosplay(cosplayID int64) (err error) {
	if _, err = dao_party.DeleteCosplayStmt.Exec(cosplayID); err != nil {
		Err(err)
	}

	return
}
