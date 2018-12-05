package npc

import (
	"appcfg"
	"constants"
	. "logger"
	"tools"
)

var npcMap map[string]map[int]*NpcDetail

type NpcDetail struct {
	ID                 string
	FriendLevel        int
	NextFriendLevelExp int

	BuffDesc string

	//buff 1
	Buff1Type   string
	Buff1Param1 string
	Buff1Param2 string

	//buff 2
	Buff2Type   string
	Buff2Param1 string
	Buff2Param2 string
}

func init() {
	if appcfg.GetServerType() != "" {
		return
	}

	loadNpc()
}

func GetNpcFriendshipLevelupExp(npcID string, level int) int {
	if npcMap[npcID] != nil {
		if npcMap[npcID][level] != nil {
			return npcMap[npcID][level].NextFriendLevelExp
		}
	}

	return 0
}

func GetNpcBuff(npcID string, level int) (string, string, string, string, string, string) {
	if npcMap[npcID] != nil {
		if npcMap[npcID][level] != nil {
			fd := npcMap[npcID][level]
			return fd.Buff1Type, fd.Buff1Param1, fd.Buff1Param2, fd.Buff2Type, fd.Buff2Param1, fd.Buff2Param2
		}
	}

	return "", "", "", "", "", ""
}

func loadNpc() {
	npcMap = make(map[string]map[int]*NpcDetail)
	path := "scripts/data_pixie/friendship.csv"

	records := tools.LoadCSV(path)

	for k, v := range records {
		if k > 0 {
			var nd NpcDetail
			for k1, v1 := range v {
				if k1 == 0 {
					nd.ID = v1
				} else if k1 == 1 {
					nd.FriendLevel = tools.LoadGetInt(v1, "friend level")
				} else if k1 == 2 {
					nd.NextFriendLevelExp = tools.LoadGetInt(v1, "next friend level exp")
				} else if k1 == 3 {
					nd.BuffDesc = v1
				} else if k1 == 4 {
					nd.Buff1Type = v1
				} else if k1 == 5 {
					nd.Buff1Param1 = v1
				} else if k1 == 6 {
					nd.Buff1Param2 = v1
				} else if k1 == 7 {
					nd.Buff2Type = v1
				} else if k1 == 8 {
					nd.Buff2Param1 = v1
				} else if k1 == 9 {
					nd.Buff2Param2 = v1
				}
			}

			if (nd.Buff1Type != "" && !constants.NPC_BUFF_MAP[nd.Buff1Type]) || (nd.Buff2Type != "" && !constants.NPC_BUFF_MAP[nd.Buff2Type]) {
				panic("unknown buff type " + nd.Buff1Type + " or " + nd.Buff2Type)
			}

			if nd.ID == "" {
				continue
			}

			if npcMap[nd.ID] == nil {
				npcMap[nd.ID] = make(map[int]*NpcDetail)
			}

			npcMap[nd.ID][nd.FriendLevel] = &nd
		}
	}

	size := len(npcMap)
	if size < 1 {
		panic("load npc friendship size 0")
	} else {
		Info("all npc friendship config size", size)
	}
}
