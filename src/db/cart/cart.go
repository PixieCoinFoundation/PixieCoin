package cart

import (
	"constants"
	"dao_designer"
	"database/sql"
	. "logger"
	. "pixie_contract/api_specification"
	"tools"
	. "types"
)

//购物车衣服库存紧张推送时用
// func ListCartByPage(page, pageSize int) (res []SimpleCart, err error) {
// 	start := (page - 1) * pageSize

// 	var rows *sql.Rows
// 	if rows, err = dao_designer.ListShopCartByPageStmt.Query(start, pageSize); err != nil {
// 		Err(err)
// 		return
// 	} else {
// 		res = make([]SimpleCart, 0)
// 		for rows.Next() {
// 			var sc SimpleCart
// 			if err = rows.Scan(&sc.Username, &sc.ClothesID); err != nil {
// 				Err(err)
// 				return
// 			}

// 			res = append(res, sc)
// 		}
// 	}

// 	return
// }

//获取已上架且库存紧张的衣服
func ListSoldClothes(page, pageSize int) (res []*SimpleCustom, err error) {
	start := (page - 1) * pageSize

	var rows *sql.Rows
	if rows, err = dao_designer.ListSoldClothesStmt.Query(start, pageSize); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]*SimpleCustom, 0)
		for rows.Next() {
			var sc SimpleCustom
			var cid int64
			var model, typee int
			var name string
			if err = rows.Scan(&cid, &model, &typee, &name); err != nil {
				Err(err)
				return
			}

			sc.ClothesIDInGame = tools.GetCustomIDInGame(typee, model, constants.ORIGIN_DESIGNER, cid)
			sc.Name = name

			res = append(res, &sc)
		}
	}

	return
}

//购物车补货推送时用
func ListCartUserByPage(cloid string, page, pageSize int) (res []string, err error) {
	start := (page - 1) * pageSize

	var rows *sql.Rows
	if rows, err = dao_designer.ListShopCartUserByPageStmt.Query(cloid, start, pageSize); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]string, 0)
		for rows.Next() {
			var u string
			if err = rows.Scan(&u); err != nil {
				Err(err)
				return
			}

			res = append(res, u)
		}
	}

	return
}

func AddToCart(username, clothesID string, nt int64) (err error) {
	if _, err = dao_designer.AddShopCartStmt.Exec(username, clothesID, nt, 1, nt); err != nil {
		Err(err)
	}

	return
}

func IncreCart(username, clothesID string) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_designer.IncreShopCartStmt.Exec(1, username, clothesID); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		success = true
	}

	return
}

func DecreCart(username, clothesID string) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = dao_designer.DecreShopCartStmt.Exec(1, username, clothesID, 1); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		success = true
	}

	return
}

func DelCart(username, clothesID string) (err error) {
	if _, err = dao_designer.DelShopCartStmt.Exec(username, clothesID); err != nil {
		Err(err)
	}

	return
}

func ListCart(username string) (cart map[string]*ShopCartDetail, err error) {
	var rows *sql.Rows
	if rows, err = dao_designer.ListShopCartStmt.Query(username); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		cart = make(map[string]*ShopCartDetail)
		var cloid string
		var ut int64
		var count int
		for rows.Next() {
			if err = rows.Scan(&cloid, &ut, &count); err != nil {
				Err(err)
				return
			}

			cart[cloid] = &ShopCartDetail{Count: count, AddTime: ut}
		}
	}

	return
}
