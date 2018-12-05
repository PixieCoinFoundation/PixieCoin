package api_specification

//领取临时钱包中的钱
const HTTP_PLAYER_DRAW_STREET_WALLET string = "api/drawWallet"

type DrawStreetWalletReq struct {
	CommonReq
}
type DrawStreetWalletResp struct {
	RewardGold  int64
	CurrentGold int64

	//我的街目前状态
	StreetDetail *StreetOperateResult
}

//雇佣或解雇清洁设施
const HTTP_PLAYER_SET_STREET_CLEAN string = "api/setClean"

type SetStreetCleanReq struct {
	Type          int    //0:雇佣 1:解雇
	CleanObjectID string //清洁设施ID
	Count         int    //雇佣或解雇的数量 必须大于0
}

type SetStreetCleanResp struct {
	//我的街目前状态
	StreetDetail *StreetOperateResult
	CurrentGold  int64
	CurrentPxc   float64
}
