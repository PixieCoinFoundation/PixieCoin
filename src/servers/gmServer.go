package servers

import (
	// "net/http"
	"appcfg"
	// "crypto/tls"
	"github.com/web"
	_ "gm/router"
	"runtime"
)

// import (
// 	"gm_handle/activityHandle"
// 	"gm_handle/adminHandle"
// 	"gm_handle/cosHandle"
// 	// "gm_handle/customHandle"
// 	"constants"
// 	"gm_handle/elseHandle"
// 	"gm_handle/eventHandle"
// 	"gm_handle/gameParamHandle"
// 	"gm_handle/giftHandle"
// 	"gm_handle/gmlogHandle"
// 	"gm_handle/gmoplogHandle"
// 	"gm_handle/mailHandle"
// 	"gm_handle/maintainHandle"
// 	"gm_handle/orderLogHandle"
// 	"gm_handle/playerHandle"
// 	"gm_handle/publicHandle"
// 	"gm_handle/rtPartyHandle"
// 	// "gm_handle/pushHandle"
// 	"gm_handle/guildHandle"
// 	"gm_handle/partyHandle"
// 	"gm_handle/reportHandle"
// 	"gm_handle/versionHandle"
// 	. "logger"
// 	// _ "player_mem"
// )

func StartGMServer() {
	// 启动服务器

	// 配置最大的go thread
	runtime.GOMAXPROCS(runtime.NumCPU())
	addr := appcfg.GetString("bind_ip", "127.0.0.1") + appcfg.GetString("port", ":29970")

	// cert, err := tls.LoadX509KeyPair("./cert/mycert1.cer", "./cert/mycert1.key")
	// if err != nil {
	// 	panic(err)
	// }
	// config := &tls.Config{Certificates: []tls.Certificate{cert}}

	// web.RunTLS(addr, config)
	web.Run(addr)

	// 管理员页面
	// http.HandleFunc("/GFMServer/register", adminHandle.RegisterPage)
	// http.HandleFunc("/GFMServer/login", adminHandle.LoginPage)
	// http.HandleFunc("/GFMServer/adminList", adminHandle.AdminListHandle)
	// http.HandleFunc("/GFMServer/delAdmin", adminHandle.DelAdminHandle)

	// // 版本管理
	// http.HandleFunc("/GFMServer/version", versionHandle.VersionPage)
	// http.HandleFunc("/GFMServer/setTishenStatus", versionHandle.SetTishenStatusHandle)

	// // 玩家页面
	// http.HandleFunc("/GFMServer/player", playerHandle.PlayerPage)

	// // 邮件相关
	// http.HandleFunc("/GFMServer/mail", mailHandle.MailPage)
	// http.HandleFunc("/GFMServer/sendMailToAll", mailHandle.SendMailToAllPage)
	// http.HandleFunc("/GFMServer/sendMailToOne", mailHandle.SendMailToOnePage)

	// http.HandleFunc("/GFMServer/sendMailMakeup", mailHandle.MailMakeupPage)
	// http.HandleFunc("/GFMServer/doMailMakeup", mailHandle.DoSendMailMakeup)

	// // 推送相关
	// // http.HandleFunc("/GFMServer/push", pushHandle.PushPage)
	// // http.HandleFunc("/GFMServer/pushToAll", pushHandle.PushToAllPage)
	// // http.HandleFunc("/GFMServer/pushToOne", pushHandle.PushToOnePage)

	// // 活动相关
	// http.HandleFunc("/GFMServer/publishCos", cosHandle.PublishCosPage)
	// http.HandleFunc("/GFMServer/findCos", cosHandle.FindCos)
	// http.HandleFunc("/GFMServer/findClosedCos", cosHandle.FindClosedCos)

	// //gm log
	// http.HandleFunc("/GFMServer/gmLog", gmlogHandle.QueryGMLogPage)
	// http.HandleFunc("/GFMServer/gmOPLog", gmoplogHandle.QueryGMOPLogPage)

	// //report
	// http.HandleFunc("/GFMServer/queryReport", reportHandle.QueryReport)
	// http.HandleFunc("/GFMServer/downloadReport", reportHandle.DownloadReport)
	// http.HandleFunc("/GFMServer/onlineReport", reportHandle.QueryOnlineReport)

	// //order log
	// http.HandleFunc("/GFMServer/orderLog", orderLogHandle.OrderLogPage)
	// http.HandleFunc("/GFMServer/doQueryOrderLog", orderLogHandle.FindOrderLog)
	// // http.HandleFunc("/GFMServer/revokeOrder", playerHandle.RevokeOrderHandle)
	// http.HandleFunc("/GFMServer/setOrderPaid", orderLogHandle.SetOrderPaidHandle)

	// //event
	// http.HandleFunc("/GFMServer/event", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/levelSEvent", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/tiliEvent", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/loginCloEvent", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/pkEvent", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/cosEvent", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/monthSignEvent", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/newPlayerEvent", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/levelDropEvent", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/magicBeedClothesEvent", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/diamondEvent", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/diamondIapEvent", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/gameEvent", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/announcement", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/tinyPack", eventHandle.PublishEventHandle)
	// http.HandleFunc("/GFMServer/beilvEvent", eventHandle.PublishEventHandle)

	// http.HandleFunc("/GFMServer/eventManage", eventHandle.EventManageHandle)
	// http.HandleFunc("/GFMServer/maintain", maintainHandle.MaintainHandle)
	// http.HandleFunc("/GFMServer/channelLimit", maintainHandle.ChannelLimitHandle)
	// http.HandleFunc("/GFMServer/setChannelLimit", maintainHandle.SetChannelLimitHandle)
	// http.HandleFunc("/GFMServer/addChannelLimit", maintainHandle.AddChannelLimitHandle)
	// http.HandleFunc("/GFMServer/delChannelLimit", maintainHandle.DelChannelLimitHandle)
	// http.HandleFunc("/GFMServer/clearGame", maintainHandle.ClearGameHandle)
	// http.HandleFunc("/GFMServer/gameParam", gameParamHandle.GameParamHandle)
	// http.HandleFunc("/GFMServer/setPKDiscount", gameParamHandle.SetPKDiscountHandle)

	// //party manage
	// http.HandleFunc("/GFMServer/partyList", partyHandle.PartyManageHandle)
	// http.HandleFunc("/GFMServer/doCloseParty", partyHandle.ClosePartyHandle)
	// http.HandleFunc("/GFMServer/doListPartyPlayer", partyHandle.ListPartyPlayerHandle)

	// //充值部分操作
	// http.HandleFunc("/GFMServer/setOrderStatus", gameParamHandle.SetOrderStatusHandle) //充值状态修改
	// http.HandleFunc("/GFMServer/setMonthCardStatus", gameParamHandle.SetMonthCardStatusHandle)
	// http.HandleFunc("/GFMServer/setGiftPackStatus", gameParamHandle.SetGiftPackStatusHandle)
	// http.HandleFunc("/GFMServer/setGuestLoginStatus", gameParamHandle.SetGuestLoginStatusHandle)
	// http.HandleFunc("/GFMServer/setDelUserDuration", gameParamHandle.SetDelUserDurationHandle)
	// // http.HandleFunc("/GFMServer/setTheater", gameParamHandle.SetTheaterHandle)
	// // http.HandleFunc("/GFMServer/addTheater", gameParamHandle.AddTheaterHandle)
	// http.HandleFunc("/GFMServer/setMonthCardReward", gameParamHandle.SetMonthCardRewardHandle)
	// // http.HandleFunc("/GFMServer/closeTheater", gameParamHandle.CloseTheaterHandle)
	// http.HandleFunc("/GFMServer/deleteMonthCardReward", gameParamHandle.DeleteMonthCardRewardHandle)

	// http.HandleFunc("/GFMServer/rtParty", rtPartyHandle.RTPartyHandle)
	// http.HandleFunc("/GFMServer/addRTPartyDefaultSubject", rtPartyHandle.AddRTPartySubjectHandle)
	// http.HandleFunc("/GFMServer/delRTPartyDefaultSubject", rtPartyHandle.DelRTPartySubjectHandle)

	// http.HandleFunc("/GFMServer/resetPasswd", adminHandle.ResetPasswordPage)
	// http.HandleFunc("/GFMServer/doResetPassword", adminHandle.ResetPasswdHandle)
	// // 逻辑--------------------------------------------
	// // 管理员相关
	// http.HandleFunc("/GFMServer/doRegister", adminHandle.Register)
	// http.HandleFunc("/GFMServer/doLogin", adminHandle.Login)

	// // 版本相关
	// http.HandleFunc("/GFMServer/updateVersion", versionHandle.UpdateVersion)

	// // 玩家相关
	// http.HandleFunc("/GFMServer/findPlayer", playerHandle.FindPlayerPage)
	// http.HandleFunc("/GFMServer/doFindPlayer", playerHandle.FindPlayer)
	// http.HandleFunc("/GFMServer/findCode", playerHandle.FindCodePage)
	// http.HandleFunc("/GFMServer/findGameCode", playerHandle.FindGameCode)

	// //guild
	// http.HandleFunc("/GFMServer/findGuild", guildHandle.FindGuildPage)
	// http.HandleFunc("/GFMServer/doFindGuild", guildHandle.FindGuild)
	// http.HandleFunc("/GFMServer/disbandGuild", guildHandle.DiabandGuildHandle)
	// http.HandleFunc("/GFMServer/setGuildOwner", guildHandle.SetGuildOwnerHandle)
	// http.HandleFunc("/GFMServer/doChangeGuildDesc", guildHandle.ChangeGuildDescHandle)
	// http.HandleFunc("/GFMServer/doChangeGuildActivity", guildHandle.ChangeGuildActivityHandle)
	// http.HandleFunc("/GFMServer/doChangeGuildMedal", guildHandle.ChangeGuildMedalHandle)
	// http.HandleFunc("/GFMServer/kickGuildMember", guildHandle.KickGuildMemberHandle)

	// http.HandleFunc("/GFMServer/doResetMonthCard", playerHandle.ResetMonthCardHandle)
	// http.HandleFunc("/GFMServer/doManagePlayer", playerHandle.ManagePlayer)
	// http.HandleFunc("/GFMServer/doChangePlayer", playerHandle.ChangePlayer)
	// http.HandleFunc("/GFMServer/grantAllClothes", playerHandle.GrantAllClothesHandle)
	// http.HandleFunc("/GFMServer/clearAllRecords", playerHandle.RecordsClearHandle)
	// //背景相关
	// http.HandleFunc("/GFMServer/addBg", playerHandle.AddBG)
	// http.HandleFunc("/GFMServer/deleteBg", playerHandle.DeleteBG)

	// //初始化设计师icon
	// http.HandleFunc("/GFMServer/initIcon", playerHandle.InitPlayerIcon)
	// http.HandleFunc("/GFMServer/resetTheatre", playerHandle.ResetTheatreHandle)
	// http.HandleFunc("/GFMServer/revokeDelUser", playerHandle.RevokeDelUserHandle)
	// http.HandleFunc("/GFMServer/delPlayerHope", playerHandle.DelHopeHandle)
	// http.HandleFunc("/GFMServer/delPlayerBoardMsg", playerHandle.DelBoardMsgHandle)
	// http.HandleFunc("/GFMServer/resetSendHopeTime", playerHandle.ResetSendHopeTimeHandle)
	// http.HandleFunc("/GFMServer/delPlayerPKRecycle", playerHandle.DelPKRecycleHandle)

	// // 邮件相关
	// http.HandleFunc("/GFMServer/doSendMailToAll", mailHandle.SendMailToAll)
	// http.HandleFunc("/GFMServer/doSendMailToOne", mailHandle.SendMailToOne)
	// http.HandleFunc("/GFMServer/doMailConfirm", mailHandle.MailConfirmHandle)
	// http.HandleFunc("/GFMServer/doMailToAllConfirm", mailHandle.MailToAllConfirmHandle)
	// http.HandleFunc("/GFMServer/cronMailManage", mailHandle.CronMailManagePage)
	// http.HandleFunc("/GFMServer/doDeleteCronMail", mailHandle.DeleteCronMailHandle)

	// // 推送相关
	// // http.HandleFunc("/GFMServer/doPushToAll", pushHandle.DoPushToAll)
	// // http.HandleFunc("/GFMServer/doPushToOne", pushHandle.DoPushToOne)

	// // cosplay
	// http.HandleFunc("/GFMServer/doPublishCos", cosHandle.PublishCos)
	// // http.HandleFunc("/GFMServer/doDeleteCos", cosHandle.DeleteCosHandle)
	// // http.HandleFunc("/GFMServer/doCopyCos", cosHandle.CopyCosHandle)

	// //event
	// http.HandleFunc("/GFMServer/doPublishEvent", eventHandle.DoPublishEventHandle)
	// http.HandleFunc("/GFMServer/doPublishMagicBeedClothesEvent", eventHandle.DoPublishMagicBeedClorhesHandle)
	// http.HandleFunc("/GFMServer/doStickEvent", eventHandle.StickEventHandle)
	// http.HandleFunc("/GFMServer/doDelEvent", eventHandle.DeleteEventHandle)
	// http.HandleFunc("/GFMServer/updateEvent", eventHandle.UpdateEventPage)
	// http.HandleFunc("/GFMServer/doUpdateEvent", eventHandle.UpdateEventHandle)
	// http.HandleFunc("/GFMServer/rollingAnnouncement", eventHandle.RollingAnnouncementPageHandle)
	// //时空之门
	// http.HandleFunc("/GFMServer/doDelTimeDoor", activityHandle.TimeDoorDelHandle)
	// http.HandleFunc("/GFMServer/bonusEvent", activityHandle.BonusEventPublishHandle)
	// http.HandleFunc("/GFMServer/updateActivityConfig", activityHandle.DoUpsertActivityConfigHangle)

	// http.HandleFunc("/GFMServer/doPublishRollingAnnouncement", eventHandle.PublishRollingAnnouncementHandle)
	// http.HandleFunc("/GFMServer/delRollingAnnouncement", eventHandle.DelRollingAnnouncementHandle)
	// http.HandleFunc("/GFMServer/changeRollingAnnouncement", eventHandle.ChangeRollingAnnouncementHandle)

	// //push
	// http.HandleFunc("/GFMServer/push", maintainHandle.PushPageHandle)
	// http.HandleFunc("/GFMServer/doPublishPush", maintainHandle.PublishPushHandle)
	// http.HandleFunc("/GFMServer/delPush", maintainHandle.DelPushHandle)
	// http.HandleFunc("/GFMServer/changePush", maintainHandle.ChangePushHandle)

	// http.HandleFunc("/GFMServer/maintainjob", maintainHandle.MaintainJobPageHandle)
	// http.HandleFunc("/GFMServer/doPublishMaintainJob", maintainHandle.PublishMaintainJobHandle)
	// http.HandleFunc("/GFMServer/delMaintainJob", maintainHandle.DelMaintainJobHandle)
	// http.HandleFunc("/GFMServer/changeMaintainJob", maintainHandle.ChangeMaintainJobHandle)

	// http.HandleFunc("/GFMServer/setMaintainStatus", maintainHandle.SetMaintainStatusHandle)
	// http.HandleFunc("/GFMServer/setPatchStatus", maintainHandle.SetPatchStatusHandle)
	// http.HandleFunc("/GFMServer/addWhiteList", maintainHandle.AddWhiteHandle)
	// http.HandleFunc("/GFMServer/delWhiteName", maintainHandle.DelWhiteHandle)
	// http.HandleFunc("/GFMServer/addBlockList", maintainHandle.AddBlockHandle)
	// http.HandleFunc("/GFMServer/delBlockName", maintainHandle.DelBlockHandle)
	// http.HandleFunc("/GFMServer/addErrMailReceiver", maintainHandle.AddErrMailReceiverHandle)
	// http.HandleFunc("/GFMServer/delErrMailReceiver", maintainHandle.DelErrMailReceiverHandle)

	// http.HandleFunc("/GFMServer/addBlackGameServerList", maintainHandle.AddBlackGameServerHandle)
	// http.HandleFunc("/GFMServer/delBlackName", maintainHandle.DelBlackGameServerHandle)

	// http.HandleFunc("/GFMServer/addContractDesigner", gameParamHandle.AddContractDesignerHandle)
	// http.HandleFunc("/GFMServer/delContractDesigner", gameParamHandle.DelContractDesignerHandle)

	// http.HandleFunc("/GFMServer/addBlackAccountList", maintainHandle.AddBlackAccountHandle)
	// http.HandleFunc("/GFMServer/delBlackAccountName", maintainHandle.DelBlackAccountHandle)

	// //gift pack
	// http.HandleFunc("/GFMServer/giftPackList", giftHandle.GiftPackListPage)
	// http.HandleFunc("/GFMServer/addGiftPack", giftHandle.AddGiftPackPage)
	// http.HandleFunc("/GFMServer/doAddGiftPack", giftHandle.AddGiftPackHandle)
	// http.HandleFunc("/GFMServer/doDeleteGiftPack", giftHandle.DeleteGiftPackHandle)
	// http.HandleFunc("/GFMServer/genGiftCode", giftHandle.GenGiftCodePage)
	// http.HandleFunc("/GFMServer/doGenGiftCode", giftHandle.GenGiftCodeHandle)
	// http.HandleFunc("/GFMServer/doQueryGiftCode", giftHandle.QueryGiftCodeHandle)
	// http.HandleFunc("/GFMServer/doDownloadGiftCode", giftHandle.DownGiftCodeHandle)
	// http.HandleFunc("/GFMServer/doQueryGiftCodeStatus", giftHandle.QueryGiftCoseStatusHandle)
	// http.HandleFunc("/GFMServer/doUploadGiftCode", giftHandle.UploadGiftCodeHandle)

	// //internal api
	// http.HandleFunc("/GFMServer/queryRT", reportHandle.QueryRT)
	// http.HandleFunc("/GFMServer/queryRTInit", reportHandle.QueryRTInit)
	// http.HandleFunc("/GFMServer/queryServerHB", reportHandle.QueryServerHB)

	// //public service
	// http.HandleFunc("/GFMServer/iosUseGiftCode", publicHandle.IOSUseGiftCodeHandle)
	// // http.HandleFunc("/GFMServer/yuyueGift", publicHandle.YuYueCanGiftHandle)

	// //else service
	// http.HandleFunc("/GFMServer/else", elseHandle.ElsePage)
	// http.HandleFunc("/GFMServer/queryElse", elseHandle.DownloadElseReport)

	// //gm log
	// http.HandleFunc("/GFMServer/doQueryGMLog", gmlogHandle.FindGMLog)

	// //gm op log
	// http.HandleFunc("/GFMServer/doQueryGMOPLog", gmoplogHandle.FindGMOPLog)

	// // // 图片相关
	// http.Handle("/GFMServer/dev/", http.StripPrefix("/GFMServer/dev/", http.FileServer(http.Dir("./dev/"))))
	// //	图片相关
	// http.Handle("/GFMServer/pic/", http.StripPrefix("/GFMServer/pic/", http.FileServer(http.Dir("./downloads/"))))

	// //test
	// // http.Handle("/GFMServer/test/", http.StripPrefix("/GFMServer/test/", http.FileServer(http.Dir("./test/"))))
	// // http.HandleFunc("/GFMServer/test/test", elseHandle.TestHandle)
	// // 静态资源
	// http.Handle("/GFMServer/web_res/", http.StripPrefix("/GFMServer/web_res/", http.FileServer(http.Dir("./web_res/"))))
	// addr := appcfg.GetString("bind_ip", "127.0.0.1") + appcfg.GetString("http_port", ":29971")
	// Info("Server start listen at:" + addr)
	// // err := http.ListenAndServeTLS(addr, "cert/gfcert/gf.crt", "cert/gfcert/gf.key", nil)
	// if appcfg.GetServerType() == constants.SERVER_TYPE_GM && appcfg.GetLanguage() == constants.KOR_LANGUAGE && appcfg.GetBool("upload_kor_log", false) {
	// 	//go Logfile()
	// }
	// err := http.ListenAndServe(addr, nil)
	// if err != nil {
	// 	Err("Server start up error:", err)
	// 	return
	// }
}

// func serveHTTPS() {
// 	// 启动服务器
// 	addr := appcfg.GetString("bind_ip", "127.0.0.1") + appcfg.GetString("port", ":29970")

// 	Info("Server start listen at :" + addr)
// 	err := http.ListenAndServeTLS(addr, "./cert/mycert1.cer", "./cert/mycert1.key", nil)
// 	if err != nil {
// 		Err("Server start up error:", err)
// 		return
// 	}
// }
