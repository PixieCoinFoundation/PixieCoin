package api_specification

import (
	"encoding/json"
	"fmt"
)

func jsonEqual(v interface{}, v1 interface{}) (bool, string) {
	b1, _ := json.Marshal(v)
	b2, _ := json.Marshal(v1)

	if string(b1) != string(b2) {
		return false, string(b1) + "【---】" + string(b2)
	}
	return true, ""
}

func (s Status) String() string {
	jsonbyte, _ := json.Marshal(s)
	return string(jsonbyte)
}

func Equal(self *Status, s *Status) (bool, string) {
	if self.Username != s.Username {
		return false, "username"
	}
	if self.Uid != s.Uid {
		return false, "uid"
	}
	if self.Nickname != s.Nickname {
		return false, "nickname"
	}
	if self.Head != s.Head {
		return false, "head"
	}
	if self.Money != s.Money {
		return false, "money"
	}

	if self.Tili != s.Tili {
		return false, fmt.Sprintf("Tili %d %d", self.Tili, s.Tili)
	}
	if self.MaxTili != s.MaxTili {
		return false, "MaxTili"
	}
	if self.TiliRTime != s.TiliRTime {
		return false, fmt.Sprintf("TiliRTime", self.TiliRTime, s.TiliRTime)
	}

	//clothes
	if e, r := jsonEqual(self.Clothes, s.Clothes); !e {
		return false, "Clothes:" + r
	}
	if e, r := jsonEqual(self.Scripts, s.Scripts); !e {
		return false, "scripts:" + r
	}
	if e, r := jsonEqual(self.Ending, s.Ending); !e {
		return false, "ending:" + r
	}
	if e, r := jsonEqual(self.Tasks, s.Tasks); !e {
		return false, "tasks:" + r
	}
	if e, r := jsonEqual(self.RecordMap, s.RecordMap); !e {
		return false, "records:" + r
	}
	if e, r := jsonEqual(self.DayExtraInfo, s.DayExtraInfo); !e {
		return false, "day_extra_info:" + r
	}
	if e, r := jsonEqual(self.ExtraInfo1, s.ExtraInfo1); !e {
		return false, "Extra1:" + r
	}
	if e, r := jsonEqual(self.ExtraInfo2, s.ExtraInfo2); !e {
		return false, "Extra2:" + r
	}
	if e, r := jsonEqual(self.HomeShow, s.HomeShow); !e {
		return false, "HomeShow:" + r
	}
	if e, r := jsonEqual(self.SuitMap, s.SuitMap); !e {
		return false, "SuitMap:" + r
	}
	if e, r := jsonEqual(self.NpcFriendship, s.NpcFriendship); !e {
		return false, "NpcFriendship:" + r
	}
	if e, r := jsonEqual(self.NpcBuff, s.NpcBuff); !e {
		return false, "NpcBuff:" + r
	}
	if e, r := jsonEqual(self.StreetOperate, s.StreetOperate); !e {
		return false, "StreetOperate:" + r
	}

	return true, ""
}
