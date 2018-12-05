package paper

import (
	"constants"
	"dao_pixie"
	"database/sql"
	. "logger"
	. "pixie_contract/api_specification"
	"sort"
	"time"
)

func LogPaperClothesSale(saleUsername string, landID int64, price float64, priceType int, tNow int64) {
	var gold int64
	var pxc float64
	if priceType == constants.PIXIE_GOLD_TYPE {
		gold = int64(price)
	} else if priceType == constants.PIXIE_PXC_TYPE {
		pxc = price
	}

	date := time.Unix(tNow, 0).Format("20060102")

	if _, err := dao_pixie.InsertClothesSaleInfoStmt.Exec(saleUsername, pxc, gold, date, 1, pxc, gold, 1); err != nil {
		Err(err)
	}

	if _, err := dao_pixie.InsertClothesSaleInfoStmt.Exec(saleUsername, pxc, gold, constants.GET_SALE_INFO_KEY_ALL, 1, pxc, gold, 1); err != nil {
		Err(err)
	}

	if _, err := dao_pixie.InsertLandClothesSaleInfoStmt.Exec(landID, saleUsername, date, pxc, gold, 1, pxc, gold, 1); err != nil {
		Err(err)
	}

	if _, err := dao_pixie.InsertLandClothesSaleInfoStmt.Exec(landID, saleUsername, constants.GET_SALE_INFO_KEY_ALL, pxc, gold, 1, pxc, gold, 1); err != nil {
		Err(err)
	}
}

func SellPaperClothes(paperID int64, logID int64, buyCount int, saleUsername string, price float64, priceType int) (success bool, err error) {

	var tx *sql.Tx
	var res sql.Result
	var effect int64

	if tx, err = dao_pixie.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	if res, err = dao_pixie.SellClothesStmt.Exec(buyCount, buyCount, saleUsername, buyCount, price, priceType, paperID); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect < 1 {
		err = constants.SellClothesEffectZero
		return
	}
	//更新tradelog表的库存
	if logID > 0 {
		if res, err = tx.Stmt(dao_pixie.UpdateOccupyLogSaleStmt).Exec(buyCount, logID); err != nil {
			Err(err)
			return
		} else if effect, err = res.RowsAffected(); err != nil {
			Err(err)
			return
		} else if effect < 1 {
			err = constants.PaperTradeLogSaleCountEffectZero
			return
		}
	}
	success = true

	return
}

func GetPaperOccupyLogID(paperID int64, saleUsername string) (occupyLogID int64, err error) {
	if err = dao_pixie.GetOccupyLogIDStmt.QueryRow(paperID, saleUsername).Scan(&occupyLogID); err != nil {
		Err(err)
		return
	}
	return
}

func GetClothesInSale(username, clothesType, partType string, page, pageSize int) (list []*OwnPaper, err error) {
	var rows *sql.Rows

	if partType != "" {
		if clothesType == "" {
			rows, err = dao_pixie.ListSalePaperByPartTypeOwnerStmt.Query(username, partType, true, (page-1)*pageSize, pageSize)
		} else {
			rows, err = dao_pixie.ListSalePaperByClothesTypePartTypeOwnerStmt.Query(username, clothesType, partType, true, (page-1)*pageSize, pageSize)
		}
	} else {
		if clothesType == "" {
			rows, err = dao_pixie.ListSalePaperByOwnerStmt.Query(username, true, (page-1)*pageSize, pageSize)
		} else {
			rows, err = dao_pixie.ListSalePaperByOwnerClothesTypeStmt.Query(username, clothesType, true, (page-1)*pageSize, pageSize)
		}
	}
	if err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		if err, list = getOwnPapers(rows); err != nil {
			Err(err)
			return
		}
	}
	return
}

func GetMinPricePaperByPaperID(paperID int64) (minPriceOwnerUsername string, err error) {
	var rows *sql.Rows
	if rows, err = dao_pixie.GetPaperListByPaperIDStmt.Query(); err != nil {
		Err(err)
		return
	} else {
		priceMap := make(map[string]float64)
		for rows.Next() {
			var OwnerUsername string
			var ClothPrice float64
			var ClothPriceType int
			if err = rows.Scan(&OwnerUsername, &ClothPrice, &ClothPriceType); err != nil {
				Err(err)
				return
			}
			if ClothPriceType == constants.PIXIE_GOLD_TYPE {
				priceMap[OwnerUsername] = ClothPrice
			} else if ClothPriceType == constants.PIXIE_PXC_TYPE {
				priceMap[OwnerUsername] = ClothPrice * constants.PXCTOGOLD
			}
		}
		if len(priceMap) > 0 {
			list := sortMapByValue(priceMap)
			minPriceOwnerUsername = list[0].Key
		} else {
			return "", nil
		}
	}
	return
}

type Pair struct {
	Key   string
	Value float64
}
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func sortMapByValue(m map[string]float64) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
	}
	sort.Sort(p)
	return p
}
