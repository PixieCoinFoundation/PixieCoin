package status

import (
	"constants"
	"dao_player"
	"errors"
)

import (
	. "pixie_contract/api_specification"
	. "types"
)

var stmtErr = errors.New("can't find stmt")

func UpdateStatusDetail(s *Status, name StatusDetailName) error {
	if name == constants.PIXIE_SAVE_DB_STATUS_HOME_SHOW {
		return dao_player.UpdatePlayerHomeShow(s)
	} else if name == constants.PIXIE_SAVE_DB_STATUS_SUIT_MAP {
		return dao_player.UpdatePlayerSuitMap(s)
	} else if name == constants.PIXIE_SAVE_DB_STATUS_STREET_OPERATE {
		return dao_player.UpdatePlayerStreetOperate(s)
	} else if name == constants.PIXIE_SAVE_DB_STATUS_MAP_SKIN {
		return dao_player.UpdatePlayerMapSkin(s)
	} else if name == constants.PIXIE_SAVE_DB_CURRENCY {
		return dao_player.UpdatePlayerCurrency(s)
	} else {
		return dao_player.UpdatePlayer(s)
	}
	return nil
}
