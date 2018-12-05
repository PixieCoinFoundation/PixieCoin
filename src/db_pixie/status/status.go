package status

import (
	"constants"
	"dao"
	"dao_pixie"
	"dao_player"
	"database/sql"
	"db_pixie/land"
	"rao"
)

import (
	. "logger"
	. "pixie_contract/api_specification"
)

const (
	SPLIT                  = "#^#"
	LOGIN_INFO_KEY_PREFIX  = "gf_login_info_"
	LOGIN_INFO_EXPIRE_TIME = 86400 //24hours
)

// var (
// 	startClothes string = appcfg.GetString("start_clothes", "[{\"c\":1000036,\"i\":1},{\"c\":1005002,\"i\":1},{\"c\":2000015,\"i\":1},{\"c\":2005001,\"i\":1},{\"c\":3000006,\"i\":1},{\"c\":3000019,\"i\":1},{\"c\":3005001,\"i\":1},{\"c\":4000000,\"i\":1},{\"c\":4005001,\"i\":1},{\"c\":5000000,\"i\":1},{\"c\":5005001,\"i\":1},{\"c\":6000019,\"i\":1},{\"c\":6005003,\"i\":1}]")
// )

//call this function only when new player comes in
func AddNewPlayer(s *Status) (err error) {
	var tx *sql.Tx
	if tx, err = dao_pixie.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	if err = dao_player.InsertPlayer(s); err == nil {
		if err = initLandForNewPlayer(s.Username, s.Nickname, s.Head, s.Sex, tx); err != nil {
			return
		}
	}

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func InitPlayerLand(un, nn, hd string, sx int) (err error) {
	var tx *sql.Tx
	if tx, err = dao_pixie.BeginTransaction(); err != nil {
		Err(err)
		return
	}
	defer func() {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			Err(re)
		}
	}()

	initLandForNewPlayer(un, nn, hd, sx, tx)

	if err = tx.Commit(); err != nil {
		Err(err)
		return
	}
	return
}

func initLandForNewPlayer(un, nn, hd string, sx int, tx *sql.Tx) (err error) {
	for i := 0; i < constants.PIXIE_INIT_LAND_SIZE; i++ {
		location := i + 1
		if i == 0 {
			//我的店
			if err = land.AddLand(PIXIE_LAND_TYPE_MY_SHOP, location, un, nn, hd, sx, tx); err != nil {
				return
			}
		} else if i == 1 {
			//大型土地
			if err = land.AddLand(PIXIE_LAND_TYPE_BIG, location, un, nn, hd, sx, tx); err != nil {
				return
			}
		} else {
			//默认土地
			if err = land.AddLand(PIXIE_LAND_TYPE_DEFAULT, location, un, nn, hd, sx, tx); err != nil {
				return
			}
		}
	}

	return
}

func updateLastLoginInfo(username string) {
	key := LOGIN_INFO_KEY_PREFIX + username

	conn := rao.GetConn()
	defer conn.Close()

	if _, err := conn.Do("EXPIRE", key, LOGIN_INFO_EXPIRE_TIME); err != nil {
		Err(err)
		return
	}

	return
}

// -----------------------------------------
func Find(username string) (s *Status, err error) {
	return dao_player.GetPlayer(username)
}

func GetThirdUsername(uid int64) (tu string, err error) {
	var dbAddr string
	if err = dao.GetUIDStmt.QueryRow(uid).Scan(&tu, &dbAddr); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Err("can't find player by uid when get third username:", uid)
		} else {
			Err(err)
		}
		return
	}
	return
}

// func UpsertNickname(username string, nickname string) (success bool) {
// 	var res sql.Result
// 	var err error
// 	var effect int64
// 	if res, err = dao.InsertGlobalNicknameStmt.Exec(username, nickname); err != nil {
// 		if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
// 			Info("duplicate:", username, nickname, err.Error())
// 		} else {
// 			Err(err)
// 		}

// 		//try update
// 		if res, err = dao.UpdateGlobalNicknameStmt.Exec(nickname, username); err != nil {
// 			if strings.Contains(err.Error(), constants.SQL_DUPLICATE_KEY_MSG) {
// 				Info("duplicate:", username, nickname, err.Error())
// 			} else {
// 				Err(err)
// 			}
// 			return false
// 		} else if effect, err = res.RowsAffected(); err != nil {
// 			Err(err)
// 			return false
// 		} else if effect > 0 {
// 			return true
// 		}
// 	} else if effect, err = res.RowsAffected(); err != nil {
// 		Err(err)
// 		return false
// 	} else if effect > 0 {
// 		return true
// 	}

// 	return false
// }

// func GetUsernameNickname(name string) (un string, nn string, err error) {
// 	if err = dao.GetUsernameNicknameStmt.QueryRow(name, name).Scan(&un, &nn); err != nil {
// 		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
// 			Info("user not found with name:", name)
// 		} else {
// 			Err(err)
// 		}

// 		return
// 	}
// 	return
// }

// func UpdateNickname(username string, nickname string, ipport string) (success bool) {
// 	stmt := dao_player.UpdateNicknameStmt
// 	if stmt == nil {
// 		Err("can't get stmt:", username, ipport)
// 		return false
// 	}

// 	if res, err := stmt.Exec(nickname, username); err != nil {
// 		Err(err)
// 		return false
// 	} else if effect, err := res.RowsAffected(); err != nil {
// 		Err(err)
// 		return false
// 	} else if effect > 0 {
// 		return true
// 	} else {
// 		return false
// 	}
// }
