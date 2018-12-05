package paper

import (
	"common"
	"constants"
	"dao_pixie"
	"database/sql"
	"encoding/json"
	. "logger"
	. "pixie_contract/api_specification"
	"tools"
)

func ListAllOfficialPaper() (list []OfficialPaper, err error) {
	var rows *sql.Rows

	if rows, err = dao_pixie.ListAllOfficialPaperStmt.Query(); err != nil {
		Err(err)
		return
	} else {
		list = make([]OfficialPaper, 0)
		for rows.Next() {
			var op OfficialPaper

			if err = rows.Scan(&op.PaperID, &op.Cname, &op.Desc, &op.ClothesType, &op.PartType, &op.File, &op.Extra, &op.Style, &op.PriceType, &op.Price, &op.Star, &op.Tag1, &op.Tag2, &op.UnlockLevel, &op.STag, &op.SaleTime, &op.UploadTime, &op.AdminName); err != nil {
				Err(err)
				return
			}

			list = append(list, op)
		}
	}

	return
}

func ListAllOfficialPaperMap() (m map[string]*OfficialPaper, l []*OfficialPaper, md5 string, err error) {
	var rows *sql.Rows

	if rows, err = dao_pixie.ListAllOfficialPaperStmt.Query(); err != nil {
		Err(err)
		return
	} else {
		l = make([]*OfficialPaper, 0)
		m = make(map[string]*OfficialPaper)
		for rows.Next() {
			var op OfficialPaper

			if err = rows.Scan(&op.PaperID, &op.Cname, &op.Desc, &op.ClothesType, &op.PartType, &op.File, &op.Extra, &op.Style, &op.PriceType, &op.Price, &op.Star, &op.Tag1, &op.Tag2, &op.UnlockLevel, &op.STag, &op.SaleTime, &op.UploadTime, &op.AdminName); err != nil {
				Err(err)
				return
			}

			m[tools.GenPixieOfficialClothesID(op.PaperID)] = &op
			l = append(l, &op)
		}
	}

	var b []byte
	if b, err = json.Marshal(m); err != nil {
		Err(err)
		return
	} else if md5, err = common.GenMD5Raw(b); err != nil {
		Err(err)
		return
	}

	return
}

func GetOneOfficialPaper(id int64) (op OfficialPaper, err error) {
	if err = dao_pixie.GetOneOfficialPaperStmt.QueryRow(id).Scan(&op.PaperID, &op.Cname, &op.Desc, &op.ClothesType, &op.PartType, &op.File, &op.Extra, &op.Style, &op.PriceType, &op.Price, &op.Star, &op.Tag1, &op.Tag2, &op.UnlockLevel, &op.STag, &op.SaleTime, &op.UploadTime, &op.AdminName); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Err("no official paper when get from db", id)
		} else {
			Err(err)
		}

		return
	}

	return
}
