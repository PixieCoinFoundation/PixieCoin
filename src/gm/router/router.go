package router

import (
	"appcfg"
	"constants"
	"github.com/web"
	"gm/restful/admin"
	"gm/restful/ajax"
	"gm/restful/clothes"
	"gm/restful/player"
	"gm/restful/recommend"
	"gm/restful/server"

	"gm/restful/upload"
	"gm/restful/version"
	"net/http"
)

func init() {
	if appcfg.GetServerType() != constants.SERVER_TYPE_GM {
		return
	}
	//衣服相关
	web.Post("/uploadClothes", upload.UploadClothes)
	web.Post("/uploadBg", upload.UploadBg)
	web.Get("/getOfficialClothes", clothes.GetOfficialClothes)
	web.Get("/getOfficialBG", clothes.GetOfficialBG)
	web.Post("/deleteOfficial", clothes.DeleteOfficial)
	web.Get("/getPlayerClothes/(.*)", clothes.GetPlayerClothes)
	web.Get("/getPlayerBg/(.*)", clothes.GetPlayerBG)
	web.Get("/getSaleClothes", clothes.GetInSaleClothes)
	web.Get("/getPoolClothes", clothes.GetPoolClothes)
	web.Get("/getPoolClothes1", clothes.GetPoolClothes1)

	web.Get("/getInAuctionPaper", clothes.GetInAuctionPaper)
	web.Get("/getRejectPaper", clothes.GetRejectPaper)
	web.Get("/getInAuctionBGPaper", clothes.GetInAuctionBGPaper)
	web.Get("/getRejectBGPaper", clothes.GetRejectBGPaper)

	//衣服层级关系
	web.Get("/missionList", ajax.GetMissionName)
	web.Get("/tagZvalueInfo", ajax.GetPaperTagZvalueInfo)

	//玩家相关
	web.Get("/getPlayer/(.*)", player.GetPlayerInfo)

	//服务器相关
	web.Get("/getVersion", version.GetVersionInfo)
	web.Post("/updateVersion", version.UpdateVersionInfo)
	web.Post("/updateVersionInfoTishen", version.UpdateVersionTishen)
	web.Get("/mainTainJob", server.GetMaintenanceInfo)
	web.Post("/addMaintenance", server.AddMaintenance)
	web.Get("/getAccountList", server.GetAccountList)
	web.Get("/addWhiteAccount/(.*)", server.AddWhiteAccount)
	web.Get("/deleteWhiteAccount/(.*)", server.DeleteWhiteAccount)
	web.Get("/addBlackAccount/(.*)", server.AddBlackAccount)
	web.Get("/deleteBlackAccount/(.*)", server.DeleteBlackAccount)
	web.Get("/addServer/(.*)", server.AddServer)
	web.Get("/deleteServer/(.*)", server.DeleteServer)
	web.Get("/allowPatch", server.AllowPatch)

	//预审
	web.Get("/verifyClothesList", clothes.GetVerifyClothesList)
	web.Get("/verifyBgList", clothes.GetVerifyBgList)
	web.Get("/verifyOne/(.*)", clothes.GetVerifyOne)
	web.Post("/passVerify", clothes.PassVerify)
	web.Post("/rejectVerify", clothes.RejectVerify)

	//管理员相关
	web.Get("/adminList", admin.AdminList)
	web.Put("/deleteAdmin", admin.DeleteAdmin)
	web.Post("/addAdmin", admin.AddAdmin)
	web.Post("/login", admin.Login)
	web.Post("/addRole", admin.AddRole)
	web.Get("/getRoleList", admin.GetRoleList)
	web.Post("/changeRole", admin.ChangeRole)
	web.Get("/getNavs", admin.GetNavs)
	web.Get("/getNodes", admin.GetZnodes)
	web.Post("/addRoleApiPriv", admin.AddRoleApiPriv)

	//推荐相关
	web.Post("/addSubject", recommend.AddSubject)
	web.Post("/updateSuit", recommend.UpdateSuit)
	web.Post("/delSubject", recommend.DeleteSubject)
	web.Get("/getOneSubject/(.*)", recommend.GetOneSubject)
	web.Get("/getSubjects", recommend.GetAllSubject)

	web.Post("/addClothesToPool", recommend.AddClothesToPool)
	web.Post("/addClothesToForPool", recommend.AddClothesToForPool)
	web.Post("/delClothes", recommend.DelClothes)

	web.Post("/addClothesLink", recommend.AddClothesWithLink)
	web.Post("/delClothesLink", recommend.DelClothesLink)
	web.Get("/getAllClothesLink", recommend.GetClothesLink)
	web.Get("/getClothesLink/(.*)", recommend.GetOneClothesLink)

	web.Post("/addTopic", recommend.AddTopic)
	web.Get("/getTopic", recommend.GetTopic)
	web.Post("/delTopic", recommend.DelTopic)

	web.Get("/getClothesPool", recommend.GetClothesPool)
	web.Get("/getSelectInfo", recommend.GetSelectInfo)

	web.Get("/(.*)", http.FileServer(http.Dir("./gm")))
}
