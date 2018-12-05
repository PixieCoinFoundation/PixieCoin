package dao

import (
	"appcfg"
	"constants"
	"database/sql"
	"fmt"
	"runtime/debug"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

// var platforms map[string]*Platform

var gmRootPlatform Platform
var DB *sql.DB

//admin
var FindAdminStmt *sql.Stmt
var RegisAdminStmt *sql.Stmt
var UpdateAdminLoginStmt *sql.Stmt
var FindAdminUPStmt *sql.Stmt
var FindAllAdminStmt *sql.Stmt
var DelAdminStmt *sql.Stmt
var ResetPasswdStmt *sql.Stmt

var ListClientLogTypeStmt *sql.Stmt
var AddLogStmt *sql.Stmt

var GetGiftCodeStmt *sql.Stmt
var UpdateRegCodeStmt *sql.Stmt

var QueryDetailLogStmt *sql.Stmt

var YuyueStmt *sql.Stmt

func init() {
	if appcfg.GetServerType() != constants.SERVER_TYPE_GM {
		return
	}

	// platforms = make(map[string]*Platform)

	var err error

	dbUsername := appcfg.GetString("db_username", "")
	dbPassword := appcfg.GetString("db_password", "")
	dbIPPort := appcfg.GetString("db_ipport", "")
	dbName := appcfg.GetString("db_name", "")

	if DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?readTimeout=15s&writeTimeout=15s&timeout=10s&charset=utf8mb4", dbUsername, dbPassword, dbIPPort, dbName)); err != nil {
		panic(err)
	}

	// DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(10)

	FindAdminStmt, err = DB.Prepare("select * from gf_admin where username=? and password=?")
	printErr(err)
	FindAllAdminStmt, err = DB.Prepare("select * from gf_admin")
	printErr(err)
	DelAdminStmt, err = DB.Prepare("delete from gf_admin where username=?")
	printErr(err)
	RegisAdminStmt, err = DB.Prepare("insert into gf_admin(username,password,nickname,role) values(?,?,?,?)")
	printErr(err)
	UpdateAdminLoginStmt, err = DB.Prepare("update gf_admin set login_time=? where username=?")
	printErr(err)
	FindAdminUPStmt, err = DB.Prepare("select * from gf_admin where username=? and password=?")
	printErr(err)
	ResetPasswdStmt, err = DB.Prepare("update gf_admin set password=? where username=?")
	printErr(err)

	AddLogStmt, err = DB.Prepare("insert into gf_gmlog values(0,?,?,?,?,?,?)")
	printErr(err)

	ListClientLogTypeStmt, err = DB.Prepare("select log_id,log_name,last_refresh_time from gf_client_log_type order by log_id asc,log_name asc")
	printErr(err)

	GetGiftCodeStmt, err = DB.Prepare("select gift_pack_id,coalesce(username,'') from gf_code where code=?")
	printErr(err)
	UpdateRegCodeStmt, err = DB.Prepare("update gf_code set username=? where code=? and (username='' or username is null)")
	printErr(err)

	QueryDetailLogStmt, err = DB.Prepare("select * from gf_gmlog where username=? and c1=? and c2=? and time>=? and time<=? order by time desc")
	printErr(err)

	YuyueStmt, err = DB.Prepare("insert into gf_yuyue values(0,?)")
	printErr(err)

	// initPlatforms()
}

func initPlatforms() {
	gmRootPlatform = *NewPlatform(
		appcfg.GetString("db_ipport", ""),
		appcfg.GetString("db_name", ""),
		appcfg.GetString("redis_ipport", ""),
		appcfg.GetString("odps_table", ""),
		appcfg.GetInt("redis_db_index", 0),
		appcfg.GetString("redis_password", ""),
		appcfg.GetString("db_username", ""),
		appcfg.GetString("db_password", ""),
		appcfg.GetString("db_designer", ""),
		appcfg.GetString("db_friend", ""),
		appcfg.GetString("db_guild", ""),
		appcfg.GetString("db_party", ""))
}

func GetPlatformOdpsTable() string {
	return gmRootPlatform.OdpsTable
}

// func GetPlatformConns() (res []*sql.DB) {
// 	res = make([]*sql.DB, 0)
// 	for _, dc := range gmRootPlatform.DBMap {
// 		res = append(res, dc.DB)
// 	}

// 	return
// }

func GetRConn() redis.Conn {
	return gmRootPlatform.GetRConn()
}

func AddPartyStmt() *sql.Stmt {
	return gmRootPlatform.AddPartyStmt
}

func AddPartyItemStmt() *sql.Stmt {
	return gmRootPlatform.AddPartyItemStmt
}

func AddPartyItemCmtStmt() *sql.Stmt {
	return gmRootPlatform.AddPartyItemCmtStmt
}

func AddPartyItemFlowerStmt() *sql.Stmt {
	return gmRootPlatform.AddPartyItemFlowerStmt
}

func DeletePartyStmt() *sql.Stmt {
	return gmRootPlatform.DeletePartyStmt
}

func DeletePartyItemStmt() *sql.Stmt {
	return gmRootPlatform.DeletePartyItemStmt
}

func DeletePartyCmtStmt() *sql.Stmt {
	return gmRootPlatform.DeletePartyCmtStmt
}

func GetPlatformGuildTx() (tx *sql.Tx) {
	var err error
	if tx, err = gmRootPlatform.DBGuild.Begin(); err != nil {
		fmt.Println(err)
		return
	}
	return
}

func ClosePartyStmt() *sql.Stmt {
	return gmRootPlatform.ClosePartyStmt
}

func GetNicknameStmt() *sql.Stmt {
	return gmRootPlatform.GetNicknameStmt
}

func UpdateOrderStatusAndInfoStmt() *sql.Stmt {
	return gmRootPlatform.UpdateOrderStatusAndInfoStmt
}

func UpdateDesignerDescStmt() *sql.Stmt {
	return gmRootPlatform.UpdateDesignerDescStmt
}

func UpdateDesignerDiamondStmt() *sql.Stmt {
	return gmRootPlatform.UpdateDesignerDiamondStmt
}

//UpdateDesignerIconStmt 更新设计师designer表的icon
func UpdateDesignerIconStmt() *sql.Stmt {
	return gmRootPlatform.UpdateDesignerIcon
}

func UpdateDesignerGoldStmt() *sql.Stmt {
	return gmRootPlatform.UpdateDesignerGoldStmt
}

func UpdateGuildMedalStmt() *sql.Stmt {
	return gmRootPlatform.UpdateGuildMedalStmt
}

func UpdateGuildDescStmt() *sql.Stmt {
	return gmRootPlatform.UpdateGuildDescStmt
}

func UpdateGuildOwnerStmt() *sql.Stmt {
	return gmRootPlatform.UpdateGuildOwnerStmt
}

func FindGuildStmt() *sql.Stmt {
	return gmRootPlatform.FindGuildStmt
}

func ListGuildMemberStmt() *sql.Stmt {
	return gmRootPlatform.ListGuildMemberStmt
}

func DelUnregisterStmt() *sql.Stmt {
	return gmRootPlatform.DelUnregisterStmt
}

func AddKorLogStmt() *sql.Stmt {
	return gmRootPlatform.AddKorLogStmt
}
func QueryKorLogStmt() *sql.Stmt {
	return gmRootPlatform.QueryKorLogStmt
}

func QueryKorLogDetailStmt() *sql.Stmt {
	return gmRootPlatform.QueryKorLogDetailStmt
}

//QueryKorNoUsernameLogDetailStmt 查询表没有昵称
func QueryKorNoUsernameLogDetailStmt() *sql.Stmt {
	return gmRootPlatform.QueryKorLogNoUsernameDetailStmt
}

// //QueryKorLoginStmt 查询韩国login
// func QueryKorLoginStmt() *sql.Stmt {
// 	return gmRootPlatform.QueryKorLoginStmt

// }

// //QueryKorLogin1Stmt 查询韩国login1
// func QueryKorLogin1Stmt() *sql.Stmt {
// 	return gmRootPlatform.QueryKorLogin1Stmt

// }

// //QueryKorOrderStmt 查询韩国版本订单日志
// func QueryKorOrderStmt() *sql.Stmt {
// 	return gmRootPlatform.QueryKorOrderStmt

// }

// //QueryKorlogUploadProgressStmt 查询韩国日志进度
// func QueryKorlogUploadProgressStmt() *sql.Stmt {
// 	return gmRootPlatform.QueryKorlogUploadProgress
// }

// //AddKorLogProgressStmt 增加韩国日志进度
// func AddKorLogProgressStmt() *sql.Stmt {
// 	return gmRootPlatform.AddKorLogProgress
// }

//QueryKorGameCodeByCode1  查询game游戏码
func QueryKorGameCodeByCode1() *sql.Stmt {
	return gmRootPlatform.QueryKorGameCode1ByCode
}

//QueryKorGameCodeByCode2  查询game游戏码
func QueryKorGameCodeByCode2() *sql.Stmt {
	return gmRootPlatform.QueryKorGameCode2ByCode
}

//QueryKorGameCodeByCode3 查询game游戏码
func QueryKorGameCodeByCode3() *sql.Stmt {
	return gmRootPlatform.QueryKorGameCode3ByCode
}

//QueryKorGameCodeByUsername1  查询game游戏码
func QueryKorGameCodeByUsername1() *sql.Stmt {
	return gmRootPlatform.QueryKorGameCode1
}

//QueryKorGameCodeByUsername2  查询game游戏码
func QueryKorGameCodeByUsername2() *sql.Stmt {
	return gmRootPlatform.QueryKorGameCode2
}

//QueryKorGameCodeByUsername3 查询game游戏码
func QueryKorGameCodeByUsername3() *sql.Stmt {
	return gmRootPlatform.QueryKorGameCode3
}

func DelHopeStmt() *sql.Stmt {
	return gmRootPlatform.DelHopeStmt
}

func DelBoardMsgStmt() *sql.Stmt {
	return gmRootPlatform.DelBoardMsgStmt
}

func ListHopeStmt() *sql.Stmt {
	return gmRootPlatform.ListHopeStmt
}

func ListBoardMsgStmt() *sql.Stmt {
	return gmRootPlatform.ListBoardMsgStmt
}

func ListAllOrderStmt() *sql.Stmt {
	return gmRootPlatform.ListAllOrderStmt
}

func QuerySingleOrderStmt() *sql.Stmt {
	return gmRootPlatform.QuerySingleOrderStmt
}

//GetHopeInfoStmt 获取某个用户的愿望
func GetHopeInfoStmt() *sql.Stmt {
	return gmRootPlatform.QueryPlayerHopeStmt
}

func DeletePartyFlowerStmt() *sql.Stmt {
	return gmRootPlatform.DeletePartyFlowerStmt
}

// func GetUsername1Stmt() *sql.Stmt {
// 	return gmRootPlatform.GetUsername1Stmt
// }

func QueryGiftCodeExistStmt() *sql.Stmt {
	return gmRootPlatform.QueryGiftCodeExistStmt
}

func QueryGiftCode2ExistStmt() *sql.Stmt {
	return gmRootPlatform.QueryGiftCode2ExistStmt
}

func ListAddressStmt() *sql.Stmt {
	return gmRootPlatform.ListAddressStmt
}

func AddPlayerLog1Stmt() *sql.Stmt {
	return gmRootPlatform.AddPlayerLog1Stmt
}

func AddPlayerLog2Stmt() *sql.Stmt {
	return gmRootPlatform.AddPlayerLog2Stmt
}

func UpdateGuildMemberActivityStmt() *sql.Stmt {
	return gmRootPlatform.UpdateGuildMemberActivityStmt
}

func UpdateGuildActivityStmt() *sql.Stmt {
	return gmRootPlatform.UpdateGuildActivityStmt
}

func UpdateGuildShareStmt() *sql.Stmt {
	return gmRootPlatform.UpdateGuildShareStmt
}

func QueryGuildInfoStmt() *sql.Stmt {
	return gmRootPlatform.QueryGuildInfoStmt
}

func QueryServerHBStmt() *sql.Stmt {
	return gmRootPlatform.QueryServerHBStmt
}

func GetAllConfigStmt() *sql.Stmt {
	return gmRootPlatform.GetAllConfigStmt
}

func GetOneOrderStmt() *sql.Stmt {
	return gmRootPlatform.GetOneOrderStmt
}

func UpdateOrderStatusFromStmt() *sql.Stmt {
	return gmRootPlatform.UpdateOrderStatusFromStmt
}
func GetOnlinePlayerCntStmt() *sql.Stmt {
	return gmRootPlatform.GetOnlinePlayerCntStmt
}
func GetMaintainJobStmt() *sql.Stmt {
	return gmRootPlatform.GetMaintainJobStmt
}

func AddMaintainJobStmt() *sql.Stmt {
	return gmRootPlatform.AddMaintainJobStmt
}
func DelMaintainJobStmt() *sql.Stmt {
	return gmRootPlatform.DelMaintainJobStmt
}
func ChangeMaintainJobStmt() *sql.Stmt {
	return gmRootPlatform.ChangeMaintainJobStmt
}

func GetPushStmt() *sql.Stmt {
	return gmRootPlatform.GetPushStmt
}

func AddPushStmt() *sql.Stmt {
	return gmRootPlatform.AddPushStmt
}
func DelPushStmt() *sql.Stmt {
	return gmRootPlatform.DelPushStmt
}
func ChangePushStmt() *sql.Stmt {
	return gmRootPlatform.ChangePushStmt
}

func DelRollingAnnouncementStmt() *sql.Stmt {
	return gmRootPlatform.DelRollingAnnouncementStmt
}
func ChangeRollingAnnouncementStmt() *sql.Stmt {
	return gmRootPlatform.ChangeRollingAnnouncementStmt
}

func AddRollingAnnouncementStmt() *sql.Stmt {
	return gmRootPlatform.AddRollingAnnouncementStmt
}
func GetRollingAnnouncementStmt() *sql.Stmt {
	return gmRootPlatform.GetRollingAnnouncementStmt
}

func UpdateTishenStmt() *sql.Stmt {
	return gmRootPlatform.UpdateTishenStmt
}

// func GetUsernameStmt() *sql.Stmt {
// 	return gmRootPlatform.GetUsernameStmt
// }
func QueryOrderStmt() *sql.Stmt {
	return gmRootPlatform.QueryOrderStmt
}
func QueryPlayerLog1Stmt() *sql.Stmt {
	return gmRootPlatform.QueryPlayerLog1Stmt
}
func QueryPlayerLog2Stmt() *sql.Stmt {
	return gmRootPlatform.QueryPlayerLog2Stmt
}
func QueryPlayerDetailLog1Stmt() *sql.Stmt {
	return gmRootPlatform.QueryPlayerDetailLog1Stmt
}

//QueryPlayerDefaultlogStmt 查询韩国版本的日志
func QueryPlayerDefaultlogStmt() *sql.Stmt {
	return gmRootPlatform.QueryPlayerDefaultlogStmt
}

//QueryPlayerAllLogStmt 查询某个用户的所有日志
func QueryPlayerAllLogStmt() *sql.Stmt {
	return gmRootPlatform.QueryPlayerAllLogStmt
}

func QueryPlayerDetailLog2Stmt() *sql.Stmt {
	return gmRootPlatform.QueryPlayerDetailLog2Stmt
}

func QueryLogStmt() *sql.Stmt {
	return gmRootPlatform.QueryLogStmt
}
func GetGift2CodeListStmt() *sql.Stmt {
	return gmRootPlatform.GetGift2CodeListStmt
}
func GetGiftCodeListStmt() *sql.Stmt {
	return gmRootPlatform.GetGiftCodeListStmt
}
func QueryGift2CodeStatusStmt() *sql.Stmt {
	return gmRootPlatform.QueryGift2CodeStatusStmt
}
func QueryGiftCodeStatusStmt() *sql.Stmt {
	return gmRootPlatform.QueryGiftCodeStatusStmt
}
func GetGiftCode2Stmt() *sql.Stmt {
	return gmRootPlatform.GetGiftCode2Stmt
}

func GetGiftPackStmt() *sql.Stmt {
	return gmRootPlatform.GetGiftPackStmt
}
func GetGiftPackContentStmt() *sql.Stmt {
	return gmRootPlatform.GetGiftPackContentStmt
}
func GetEventStmt() *sql.Stmt {
	return gmRootPlatform.GetEventStmt
}
func UpdateEventStmt() *sql.Stmt {
	return gmRootPlatform.UpdateEventStmt
}

//UpdateShiKongzhimenStmt 更新时空之
func UpsertShiKongzhimenStmt() *sql.Stmt {
	return gmRootPlatform.UpsertShiKongzhimenStmt
}

//DelShiKongzhimenStmt 更新时空之
func DelShiKongzhimenStmt() *sql.Stmt {
	return gmRootPlatform.DelShiKongzhimenStmt
}

//InsertShiKongzhimenStmt 插入时空之门
func InsertShiKongzhimenStmt() *sql.Stmt {
	return gmRootPlatform.GetShikongzhimenStmt
}

//GetShiKongzhimenStmt 获取时空之门
func GetShiKongzhimenStmt() *sql.Stmt {
	return gmRootPlatform.GetShikongzhimenStmt
}

func AddEventStmt() *sql.Stmt {
	return gmRootPlatform.AddEventStmt
}

func DeleteEventStmt() *sql.Stmt {
	return gmRootPlatform.DeleteEventStmt
}
func StickEventStmt() *sql.Stmt {
	return gmRootPlatform.StickEventStmt
}
func QueryEventsStmt() *sql.Stmt {
	return gmRootPlatform.QueryEventsStmt
}
func QueryNotOpenEventsStmt() *sql.Stmt {
	return gmRootPlatform.QueryNotOpenEventsStmt
}
func QueryClosedEventsStmt() *sql.Stmt {
	return gmRootPlatform.QueryClosedEventsStmt
}
func GetConfigStmt() *sql.Stmt {
	return gmRootPlatform.GetConfigStmt
}
func GetSPReportStmt() *sql.Stmt {
	return gmRootPlatform.GetSPReportStmt
}

func AddGiftPackStmt() *sql.Stmt {
	return gmRootPlatform.AddGiftPackStmt
}
func DeleteGiftPackStmt() *sql.Stmt {
	return gmRootPlatform.DeleteGiftPackStmt
}
func DeleteCodeStmt() *sql.Stmt {
	return gmRootPlatform.DeleteCodeStmt
}
func DeleteCode1Stmt() *sql.Stmt {
	return gmRootPlatform.DeleteCode1Stmt
}
func DeleteCode2Stmt() *sql.Stmt {
	return gmRootPlatform.DeleteCode2Stmt
}
func AddGiftCode2Stmt() *sql.Stmt {
	return gmRootPlatform.AddGiftCode2Stmt
}
func AddGiftCodeStmt() *sql.Stmt {
	return gmRootPlatform.AddGiftCodeStmt
}
func DeleteCronMailStmt() *sql.Stmt {
	return gmRootPlatform.DeleteCronMailStmt
}
func GetCronMailStmt() *sql.Stmt {
	return gmRootPlatform.GetCronMailStmt
}
func AddCronMailStmt() *sql.Stmt {
	return gmRootPlatform.AddCronMailStmt
}
func AddMailStmt() *sql.Stmt {
	return gmRootPlatform.AddMailStmt
}
func GetAllMailStmt() *sql.Stmt {
	return gmRootPlatform.GetAllMailStmt
}
func SetConfigStmt() *sql.Stmt {
	return gmRootPlatform.SetConfigStmt
}

func GetCosplayStmt() *sql.Stmt {
	return gmRootPlatform.GetCosplayStmt
}

func GetUIDStmt() *sql.Stmt {
	return gmRootPlatform.GetUIDStmt
}
func GetAddressStmt() *sql.Stmt {
	return gmRootPlatform.GetAddressStmt
}
func GetLocalIPStmt() *sql.Stmt {
	return gmRootPlatform.GetLocalIPStmt
}

// func UpdateGlobalNicknameStmt() *sql.Stmt {
// 	return gmRootPlatform.UpdateGlobalNicknameStmt
// }
// func InsertGlobalNicknameStmt() *sql.Stmt {
// 	return gmRootPlatform.InsertGlobalNicknameStmt
// }

func QueryGiftCodeStmt() *sql.Stmt {
	return gmRootPlatform.QueryGiftCodeStmt
}

func AddCosplayStmt() *sql.Stmt {
	return gmRootPlatform.AddCosplayStmt
}
func DeleteCosplayStmt() *sql.Stmt {
	return gmRootPlatform.DeleteCosplayStmt
}

func UpdateDesignerPointStmt() *sql.Stmt {
	return gmRootPlatform.UpdateDesignerPointStmt
}
func DeleteMailStmt() *sql.Stmt {
	return gmRootPlatform.DeleteMailStmt
}
func DeleteHelpStmt() *sql.Stmt {
	return gmRootPlatform.DeleteHelpStmt
}
func QUsernameByNicknameStmt() *sql.Stmt {
	return gmRootPlatform.QUsernameByNicknameStmt
}
func QueryDesignerStmt() *sql.Stmt {
	return gmRootPlatform.QueryDesignerStmt
}
func QueryMailStmt() *sql.Stmt {
	return gmRootPlatform.QueryMailStmt
}

func QueryHelpStmt() *sql.Stmt {
	return gmRootPlatform.QueryHelpStmt
}

func UpdateGiftCodeStmt() *sql.Stmt {
	return gmRootPlatform.UpdateGiftCodeStmt
}
func UpdateGiftCode2Stmt() *sql.Stmt {
	return gmRootPlatform.UpdateGiftCode2Stmt
}
func UseOneCodeStmt() *sql.Stmt {
	return gmRootPlatform.UseOneCodeStmt
}

func GetGiftPackIDStmt() *sql.Stmt {
	return gmRootPlatform.GetGiftPackIDStmt
}
func GetGiftPackID2Stmt() *sql.Stmt {
	return gmRootPlatform.GetGiftPackID2Stmt
}

func UpdateVersionStmt() *sql.Stmt {
	return gmRootPlatform.UpdateVersionStmt
}

func QueryVersionStmt() *sql.Stmt {
	return gmRootPlatform.QueryVersionStmt
}

// func QBusernameStmt() *sql.Stmt {
// 	return gmRootPlatform.QBusernameStmt
// }

// func QueryMailStmt() *sql.Stmt {
// 	return gmRootPlatform.QueryMailStmt
// }
// func QueryMailStmt() *sql.Stmt {
// 	return gmRootPlatform.QueryMailStmt
// }
// func QueryMailStmt() *sql.Stmt {
// 	return gmRootPlatform.QueryMailStmt
// }

func printErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		debug.PrintStack()
		panic(err)
	}
}
