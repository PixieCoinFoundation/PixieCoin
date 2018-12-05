package api_specification

const HTTP_PLAYER_LIST_LAND = "api/listLand"

type ListLandReq struct {
	TargetUsername string
}

type ListLandResp struct {
	Lands []Land
}

const HTTP_PLAYER_GET_PLAYER_MAP = "api/getPlayerMap"

type GetPlayerMapReq struct {
	TargetUsername string
}

type GetPlayerMapResp struct {
	UserAttributes *Status
	Lands          []Land
}

const HTTP_PLAYER_GET_RANDOM_MAP = "api/getRandomMap"

type GetRandomMapReq struct {
}

type GetRandomMapResp struct {
	UserAttributes *Status
	Lands          []Land
}
