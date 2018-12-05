package api_specification

const HTTP_PLAYER_BUY_SKIN string = "api/buySkin"

type BuySkinReq struct {
	SkinID string //皮肤ID
}

type BuySkinResp struct {
	CurrentGold int64
	CurrentPxc  float64
}

const HTTP_PLAYER_CHANGE_STREET_SKIN string = "api/changeSkin"

type ChangeStreetSkinReq struct {
	SkinID string //皮肤ID
}

type ChangeStreetSkinResp struct {
	SkinID string //皮肤ID
}
