package types

type Hope struct {
	ID             int64
	Sender         string
	SenderNickname string
	OfferClothes   string
	OfferPart      int
	OfferNum       int
	NeedClothes    string
	NeedPart       int
	NeedNum        int
	SendTime       int64
	Status         int    //1:发布但未达成 2:已达成但未领取 3:已达成且已领取奖励(客户端忽略此状态)
	Helper         string //帮助者用户名(若有)
	HelperNickname string //帮助者昵称(若有)
}

type HopeView struct {
	Hope
	SendTimes  string
	StatusInfo string
}
