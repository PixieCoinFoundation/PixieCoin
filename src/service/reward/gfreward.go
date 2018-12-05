package reward

type GFReward struct {
	RewardID  int    `json:"id"`
	Gold      int    `json:"gold"`
	Diamond   int    `json:"diamond"`
	ClothesID string `json:"clothesID"`
	Tag       string `json:"tag"`
	Desc      string `json:"desc"`
	From      string `json:"from"` // 需要给玩家发送邮件时，显示的邮件的发件人昵称
	Title     string `json:"title"`
}
