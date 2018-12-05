package ajax

import (
	"github.com/web"
	"service_pixie/record"
)

func GetMissionName(ctx *web.Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.BackSuccess(record.GetLevelIDList())
	return
}
