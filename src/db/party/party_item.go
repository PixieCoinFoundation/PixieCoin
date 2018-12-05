package party

import (
	. "types"
)

type PartyItemList []*PartyItem

func (list PartyItemList) Len() int {
	return len(list)
}

func (list PartyItemList) Less(i, j int) bool {
	if list[i].Popularity < list[j].Popularity {
		return false
	} else if list[i].Popularity > list[j].Popularity {
		return true
	} else {
		return list[i].UploadTime >= list[j].UploadTime
	}
}

func (list PartyItemList) Swap(i, j int) {
	var temp *PartyItem = list[i]
	list[i] = list[j]
	list[j] = temp
}
