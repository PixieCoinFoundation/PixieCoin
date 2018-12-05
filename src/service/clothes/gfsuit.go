package clothes

type GFSuit struct {
	SID string `json:"id"` // 衣服id
	// Sex int    `json:"sex"` // 性别
	Name string `json:"name"` // 名字
	// Desc string `json:"desc"` // 描述

	Hair     string `json:"hair"`     // 头发
	Coat     string `json:"coat"`     //
	Shirt    string `json:"shirt"`    //
	Pants    string `json:"pants"`    //
	Socks    string `json:"socks"`    //
	Shoes    string `json:"shoes"`    //
	Dress    string `json:"dress"`    //
	Hat      string `json:"hat"`      //
	Glasses  string `json:"glasses"`  //
	Earrings string `json:"earrings"` //
	Necklace string `json:"necklace"` //
	Tie      string `json:"tie"`      //
	Wrist    string `json:"wrist"`    //
	Bag      string `json:"bag"`      //
	Other    string `json:"other"`    //
	Face     string `json:"face"`

	Clos []string `json:"clos"`

	Char int `json:"char"`
}

func (s *GFSuit) IsClothesSuit(cloID string) bool {
	if cloID == "" {
		return false
	}
	if s.Hair == cloID || s.Coat == cloID || s.Shirt == cloID || s.Pants == cloID || s.Socks == cloID || s.Shoes == cloID || s.Dress == cloID || s.Hat == cloID || s.Glasses == cloID || s.Earrings == cloID || s.Necklace == cloID || s.Tie == cloID || s.Wrist == cloID || s.Bag == cloID || s.Other == cloID || s.Face == cloID {
		return true
	}
	return false
}
