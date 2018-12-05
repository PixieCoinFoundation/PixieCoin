package designersuit

import (
	"dao_designer"
	"database/sql"
	. "logger"
	. "types"
)

//添加套装
func AddSuit(username string, nickname string, position int, suit string, suitDesc string) (err error) {
	if _, err = dao_designer.InsertSuitStmt.Exec(username, nickname, position, suit, suitDesc, suit); err != nil {
		Err(err)
		return
	}
	return
}

//获取套装
func GetSuit(username string) (list []DesignerSuit, err error) {
	var rows *sql.Rows
	if rows, err = dao_designer.GetSuitByUsernameStmt.Query(username); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var id int
			var model DesignerSuit

			if err := rows.Scan(&id, &model.Username, &model.Nickname, &model.Position, &model.Suit, &model.SuitDesc); err != nil {
				Err(err)
				return list, err
			} else {
				list = append(list, model)
			}
		}
		return
	}
}

//清除套装
func ClearSuit(username string, position int) (err error) {
	if _, err = dao_designer.DeleteSuitByUnique.Exec(username, position); err != nil {
		Err(err)
		return
	}
	return
}
func UpdateSuitDesc(username string, position int, suitDesc string) (err error) {
	if _, err = dao_designer.UpdateSuitDesc.Exec(suitDesc, username, position); err != nil {
		Err(err)
		return
	}
	return
}
