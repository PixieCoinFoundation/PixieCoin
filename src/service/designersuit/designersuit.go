package designersuit

import (
	"db/designersuit"
	. "types"
)

func AddSuit(username string, nickname string, position int, suit string, suitDesc string) error {
	return designersuit.AddSuit(username, nickname, position, suit, suitDesc)
}
func GetSuit(username string) ([]DesignerSuit, error) {
	return designersuit.GetSuit(username)
}
func DeleteSuit(username string, position int) error {
	return designersuit.ClearSuit(username, position)
}
func UpdateSuitDesc(username string, position int, suitdesc string) (err error) {
	return designersuit.UpdateSuitDesc(username, position, suitdesc)
}
