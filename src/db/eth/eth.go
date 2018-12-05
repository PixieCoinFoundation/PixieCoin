package eth

import (
	"constants"
	"dao_designer"
	"database/sql"
	"encoding/json"
	. "logger"
	. "types"
)

func AddTransaction(
	t int64,
	hash string,
	fu, fn, fa string,
	tu, tn, ta string,
	clothesID, clothesName string,
	price, buyCount int,
	amount string,
	amountType int,
	cloListJson string) {
	if _, err := dao_designer.AddEthereumTransactionStmt.Exec(t, hash, fu, fn, fa, tu, tn, ta, clothesID, clothesName, price, buyCount, amount, amountType, constants.TRANSACTION_STATUS_SUBMIT, cloListJson); err != nil {
		ErrMail("buy pxc clothes transaction add error!!!", hash, fu, tu, clothesID, price, err)
	}
}

func UpdateTransaction(hash string, status int, receipt, data string) {
	if _, err := dao_designer.UpdateEthereumTransactionStmt.Exec(status, receipt, data, hash); err != nil {
		ErrMail("transaction update error!!!", hash, status, receipt, err)
	}
}

func ListTransaction(status, page, pageSize int) (res []EthereumTransaction, err error) {
	start := (page - 1) * pageSize

	var rows *sql.Rows
	if rows, err = dao_designer.ListEthereumTransactionStmt.Query(status, start, pageSize); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		res = make([]EthereumTransaction, 0)
		for rows.Next() {
			var t EthereumTransaction
			var cj string
			if err = rows.Scan(&t.ID, &t.SubmitTime, &t.Hash, &t.FromUsername, &t.FromNickname, &t.FromAccount, &t.ToUsername, &t.ToNickname, &t.ToAccount, &t.ClothesID, &t.ClothesName, &t.Price, &t.BuyCount, &t.Amount, &t.AmountType, &t.Status, &t.Data, &t.Receipt, &cj); err != nil {
				Err(err)
				return
			}

			if cj != "" {
				cl := make([]SimpleClothes1, 0)
				if e := json.Unmarshal([]byte(cj), &cl); e != nil {
					ErrMail(e)
					continue
				} else {
					t.ClothesList = cl
				}
			}

			res = append(res, t)
		}
	}

	return
}

func ListAllEthereumTransaction(page, pageSize int) (res []EthereumTransaction, err error) {
	start := (page - 1) * pageSize

	var rows *sql.Rows
	if rows, err = dao_designer.ListAllEthereumTransactionStmt.Query(start, pageSize); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		res = make([]EthereumTransaction, 0)
		for rows.Next() {
			var t EthereumTransaction
			var cj string
			if err = rows.Scan(&t.ID, &t.SubmitTime, &t.Hash, &t.FromUsername, &t.FromNickname, &t.FromAccount, &t.ToUsername, &t.ToNickname, &t.ToAccount, &t.ClothesID, &t.ClothesName, &t.Price, &t.BuyCount, &t.Amount, &t.AmountType, &t.Status, &t.Data, &t.Receipt, &cj); err != nil {
				Err(err)
				return
			}

			if cj != "" {
				cl := make([]SimpleClothes1, 0)
				if e := json.Unmarshal([]byte(cj), &cl); e != nil {
					ErrMail(e)
					continue
				} else {
					t.ClothesList = cl
				}
			}

			res = append(res, t)
		}
	}

	return
}

func ListUserEthereumTransaction(username string, amountType int, page, pageSize int) (res []EthereumTransaction, err error) {
	start := (page - 1) * pageSize

	var rows *sql.Rows
	if rows, err = dao_designer.ListUserEthereumTransactionStmt.Query(username, amountType, start, pageSize); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()

		res = make([]EthereumTransaction, 0)
		for rows.Next() {
			var t EthereumTransaction
			var cj string
			if err = rows.Scan(&t.ID, &t.SubmitTime, &t.Hash, &t.FromUsername, &t.FromNickname, &t.FromAccount, &t.ToUsername, &t.ToNickname, &t.ToAccount, &t.ClothesID, &t.ClothesName, &t.Price, &t.BuyCount, &t.Amount, &t.AmountType, &t.Status, &t.Data, &t.Receipt, &cj); err != nil {
				Err(err)
				return
			}

			if cj != "" {
				cl := make([]SimpleClothes1, 0)
				if e := json.Unmarshal([]byte(cj), &cl); e != nil {
					ErrMail(e)
					continue
				} else {
					t.ClothesList = cl
				}
			}

			res = append(res, t)
		}
	}

	return
}
