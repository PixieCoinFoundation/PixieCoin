package paper

import (
	"constants"
	"dao_pixie"
	"database/sql"
	. "logger"
	"tools"
	. "zk_manager"
)

func SetPaperProduct(username string, paperID int64, startTime int64, produceNum int) (err error) {
	if _, err = dao_pixie.SetPaperProductStmt.Exec(startTime, produceNum, paperID, username); err != nil {
		Err(err)
		return
	}
	return
}

func CancelPaperProduct(paperID int64, username string, tNow int64) (err error) {
	var res sql.Result
	var effect, productStartTime int64
	var productLeftCount, productDoneCount int

	//检查是否有未增加的生产数量
	if err = dao_pixie.GetOnePaperOccupyInProductStmt.QueryRow(paperID, username).Scan(&productStartTime, &productLeftCount, &productDoneCount); err != nil {
		Err(err)
		return
	}

	productCnt := tools.CheckProductNum(tNow, productStartTime, productLeftCount, productDoneCount)
	if productCnt > 0 {
		//更新PaperOccupy 校验库存数量
		if err = DoProductPaper(username, productCnt, productStartTime, productLeftCount-productCnt, productDoneCount+productCnt, paperID, productLeftCount, productDoneCount); err != nil {
			return
		}
	} else {
		//更新PaperOccupy状态
		if res, err = dao_pixie.CancelPaperProductStmt.Exec(paperID, username); err != nil {
			Err(err)
			return
		} else if effect, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		} else if effect < 1 {
			err = constants.CancelPaperProductEffectZero
			return
		}
	}

	return
}

func DoProductPaper(ownerUsername string, inventory int, productStartTime int64, leftCount, doneCount int, paperID int64, originLeft, originDone int) (err error) {
	if success, path := LockPixiePaperProduct(paperID); success {
		defer Unlock(path)

		var in_production bool
		newProductStartTime := productStartTime

		if leftCount <= 0 {
			in_production = false
			newProductStartTime = 0
		} else {
			in_production = true
		}

		var effect int64
		var res sql.Result
		if res, err = dao_pixie.PaperProductStmt.Exec(in_production, inventory, newProductStartTime, leftCount, doneCount, paperID, ownerUsername, originLeft, originDone); err != nil {
			Err(err)
			return
		} else if effect, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		} else if effect < 1 {
			Err("updateProductProcessOne effect 0", paperID, inventory)
			return
		}
	} else {
		Err("lock paper product error", paperID, inventory)
	}

	return
}

func GetPaperInventoryByIDList(idList []interface{}) (err error, inventory int64) {
	var rows *sql.Rows
	if rows, err = dao_pixie.GetPaperInventoryByIDListStmt.Query(idList...); err != nil {
		Err(err)
		return
	} else {
		for rows.Next() {
			var count int64
			if err = rows.Scan(&count); err != nil {
				Err(err)
				return
			}
			inventory += count
		}
	}
	return
}
