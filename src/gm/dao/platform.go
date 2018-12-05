package dao

import (
	"appcfg"
	"database/sql"
	"fmt"
	"time"
	. "types"

	"github.com/garyburd/redigo/redis"
)

type Platform struct {
	//platform's default db
	DBIPPort string
	DB       *sql.DB

	//other dbs
	DBDesigner *sql.DB
	DBFriend   *sql.DB
	DBGuild    *sql.DB
	DBParty    *sql.DB

	//platform's default redis
	rpool *redis.Pool

	//platform's player dbs
	// DBMap map[string]*DBConn

	//odps report table
	OdpsTable string

	//not important.
	dcs []*DBConfig

	QueryDesignerStmt         *sql.Stmt
	UpdateDesignerPointStmt   *sql.Stmt
	UpdateDesignerDiamondStmt *sql.Stmt
	UpdateDesignerDescStmt    *sql.Stmt
	UpdateDesignerGoldStmt    *sql.Stmt
	//初始化设计师的icon
	UpdateDesignerIcon *sql.Stmt
	//初始化设计师stat的icon
	// UpdateDesignerStatIcon *sql.Stmt
	//增加玩家背景
	// AddBGStmt *sql.Stmt
	//删除玩家列表
	// DeleteBGStmt *sql.Stmt

	// QUsernameStmt *sql.Stmt
	GetUIDStmt *sql.Stmt

	//order
	QueryOrderStmt               *sql.Stmt
	UpdateOrderStatusFromStmt    *sql.Stmt
	UpdateOrderStatusAndInfoStmt *sql.Stmt
	GetOneOrderStmt              *sql.Stmt
	QuerySingleOrderStmt         *sql.Stmt
	ListAllOrderStmt             *sql.Stmt

	//cosplay
	AddCosplayStmt    *sql.Stmt
	DeleteCosplayStmt *sql.Stmt
	GetCosplayStmt    *sql.Stmt

	//mail
	AddMailStmt    *sql.Stmt
	QueryMailStmt  *sql.Stmt
	DeleteMailStmt *sql.Stmt
	GetAllMailStmt *sql.Stmt

	//help
	QueryHelpStmt  *sql.Stmt
	DeleteHelpStmt *sql.Stmt

	QueryGMServerStmt *sql.Stmt

	//system
	GetLocalLoginServerStmt *sql.Stmt
	// GetUsernameStmt         *sql.Stmt
	QUsernameByNicknameStmt *sql.Stmt

	// UpdateGlobalNicknameStmt *sql.Stmt
	// InsertGlobalNicknameStmt *sql.Stmt

	//system report
	GetSPReportStmt *sql.Stmt

	GetAddressStmt *sql.Stmt
	GetLocalIPStmt *sql.Stmt

	//event
	AddEventStmt           *sql.Stmt
	StickEventStmt         *sql.Stmt
	QueryEventsStmt        *sql.Stmt
	QueryNotOpenEventsStmt *sql.Stmt
	QueryClosedEventsStmt  *sql.Stmt
	DeleteEventStmt        *sql.Stmt
	GetEventStmt           *sql.Stmt
	UpdateEventStmt        *sql.Stmt

	GetConfigStmt *sql.Stmt
	SetConfigStmt *sql.Stmt

	//log
	QueryLogStmt              *sql.Stmt
	AddPlayerLog1Stmt         *sql.Stmt
	AddPlayerLog2Stmt         *sql.Stmt
	QueryPlayerLog1Stmt       *sql.Stmt
	QueryPlayerDetailLog1Stmt *sql.Stmt
	QueryPlayerDefaultlogStmt *sql.Stmt
	QueryPlayerAllLogStmt     *sql.Stmt
	QueryPlayerLog2Stmt       *sql.Stmt
	QueryPlayerDetailLog2Stmt *sql.Stmt
	QueryPlayerHopeStmt       *sql.Stmt
	//cron mail
	AddCronMailStmt    *sql.Stmt
	GetCronMailStmt    *sql.Stmt
	UpdateCronMailStmt *sql.Stmt
	DeleteCronMailStmt *sql.Stmt

	//gift
	AddGiftPackStmt    *sql.Stmt
	DeleteGiftPackStmt *sql.Stmt
	GetGiftPackStmt    *sql.Stmt
	AddGiftCodeStmt    *sql.Stmt
	AddGiftCode2Stmt   *sql.Stmt

	GetGiftCode2Stmt         *sql.Stmt
	GetGiftCodeListStmt      *sql.Stmt
	GetGift2CodeListStmt     *sql.Stmt
	DeleteCodeStmt           *sql.Stmt
	DeleteCode2Stmt          *sql.Stmt
	DeleteCode1Stmt          *sql.Stmt
	QueryGiftCodeStatusStmt  *sql.Stmt
	QueryGift2CodeStatusStmt *sql.Stmt
	GetGiftPackContentStmt   *sql.Stmt
	UpdateGiftCodeStmt       *sql.Stmt
	UpdateGiftCode2Stmt      *sql.Stmt
	UseOneCodeStmt           *sql.Stmt
	GetGiftPackIDStmt        *sql.Stmt
	GetGiftPackID2Stmt       *sql.Stmt

	QueryGiftCodeStmt *sql.Stmt

	UpdateVersionStmt *sql.Stmt

	//UpdateShikongzhimenAvailableStmt 更新时空之门
	UpsertShiKongzhimenStmt *sql.Stmt
	//GetShikongzhimenStmt 更新时空之门
	GetShikongzhimenStmt *sql.Stmt

	DelShiKongzhimenStmt *sql.Stmt

	//InsertShikongzhimenStmt 更新时空之门
	InsertShikongzhimenStmt *sql.Stmt
	QueryVersionStmt        *sql.Stmt

	UpdateTishenStmt              *sql.Stmt
	GetRollingAnnouncementStmt    *sql.Stmt
	AddRollingAnnouncementStmt    *sql.Stmt
	DelRollingAnnouncementStmt    *sql.Stmt
	ChangeRollingAnnouncementStmt *sql.Stmt

	AddPushStmt    *sql.Stmt
	DelPushStmt    *sql.Stmt
	ChangePushStmt *sql.Stmt
	GetPushStmt    *sql.Stmt

	AddMaintainJobStmt    *sql.Stmt
	DelMaintainJobStmt    *sql.Stmt
	ChangeMaintainJobStmt *sql.Stmt
	GetMaintainJobStmt    *sql.Stmt

	GetOnlinePlayerCntStmt *sql.Stmt
	QueryServerHBStmt      *sql.Stmt

	GetAllConfigStmt *sql.Stmt

	//guild
	QueryGuildInfoStmt            *sql.Stmt
	UpdateGuildMemberActivityStmt *sql.Stmt
	UpdateGuildShareStmt          *sql.Stmt
	UpdateGuildDescStmt           *sql.Stmt
	UpdateGuildActivityStmt       *sql.Stmt
	UpdateGuildMedalStmt          *sql.Stmt

	QueryGiftCodeExistStmt  *sql.Stmt
	QueryGiftCode2ExistStmt *sql.Stmt

	ListAddressStmt *sql.Stmt
	// GetUsername1Stmt *sql.Stmt

	// QBusernameStmt   *sql.Stmt
	ListHopeStmt     *sql.Stmt
	ListBoardMsgStmt *sql.Stmt
	DelHopeStmt      *sql.Stmt
	DelBoardMsgStmt  *sql.Stmt

	GetNicknameStmt *sql.Stmt

	ClosePartyStmt *sql.Stmt

	//kor
	QueryKorLogStmt                 *sql.Stmt
	QueryKorLogDetailStmt           *sql.Stmt
	QueryKorLogNoUsernameDetailStmt *sql.Stmt

	AddKorLogStmt     *sql.Stmt
	DelUnregisterStmt *sql.Stmt
	FindGuildStmt     *sql.Stmt

	ListGuildMemberStmt  *sql.Stmt
	UpdateGuildOwnerStmt *sql.Stmt

	// //QueryKorLoginStmt 查询韩国login
	// QueryKorLoginStmt *sql.Stmt

	// //QueryKorLogin1Stmt 查询韩国login1
	// QueryKorLogin1Stmt *sql.Stmt

	// //QueryKorOrderStmt 查询韩国版本订单日志
	// QueryKorOrderStmt *sql.Stmt

	// //QueryKorlogUploadProgress 查询韩国日志上传进度
	// QueryKorlogUploadProgress *sql.Stmt

	// //AddKorLogProgress 增加韩国日志上传进度
	// AddKorLogProgress *sql.Stmt

	//QueryKorGameCode1 查询游戏码1
	QueryKorGameCode1 *sql.Stmt

	//QueryKorGameCode3 查询游戏码2
	QueryKorGameCode2 *sql.Stmt

	//QueryKorGameCode3 查询游戏码3
	QueryKorGameCode3 *sql.Stmt

	//QueryKorGameCode1 查询游戏码1
	QueryKorGameCode1ByCode *sql.Stmt

	//QueryKorGameCode3 查询游戏码2
	QueryKorGameCode2ByCode *sql.Stmt

	//QueryKorGameCode3 查询游戏码3
	QueryKorGameCode3ByCode *sql.Stmt

	AddPartyStmt           *sql.Stmt
	AddPartyItemStmt       *sql.Stmt
	AddPartyItemFlowerStmt *sql.Stmt
	AddPartyItemCmtStmt    *sql.Stmt
	DeletePartyStmt        *sql.Stmt
	DeletePartyItemStmt    *sql.Stmt
	DeletePartyCmtStmt     *sql.Stmt
	DeletePartyFlowerStmt  *sql.Stmt
}

func NewPlatform(
	db_ipport, dbName string,
	redis_ipport string, odpsTable string, reids_db_index int, redisPasswd string,
	dbUsername, dbPasswd string,
	designDBIPPort, friendDBIPPort, guildDBIPPort, partyDBIPPort string) (p *Platform) {
	fmt.Println("gm platform:", db_ipport, redis_ipport, odpsTable, reids_db_index, dbUsername, dbPasswd)
	if odpsTable == "" {
		panic("odps table config empty")
	}

	p = &Platform{
		OdpsTable: odpsTable,
	}
	var err error
	if p.DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?readTimeout=15s&writeTimeout=15s&timeout=10s&charset=utf8mb4", dbUsername, dbPasswd, db_ipport, dbName)); err != nil {
		panic(err)
	}
	p.DBIPPort = db_ipport

	if p.DBDesigner, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?readTimeout=15s&writeTimeout=15s&timeout=10s&charset=utf8mb4", dbUsername, dbPasswd, designDBIPPort, dbName)); err != nil {
		panic(err)
	}
	if p.DBFriend, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?readTimeout=15s&writeTimeout=15s&timeout=10s&charset=utf8mb4", dbUsername, dbPasswd, friendDBIPPort, dbName)); err != nil {
		panic(err)
	}
	if p.DBGuild, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?readTimeout=15s&writeTimeout=15s&timeout=10s&charset=utf8mb4", dbUsername, dbPasswd, guildDBIPPort, dbName)); err != nil {
		panic(err)
	}
	if p.DBParty, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?readTimeout=15s&writeTimeout=15s&timeout=10s&charset=utf8mb4", dbUsername, dbPasswd, partyDBIPPort, dbName)); err != nil {
		panic(err)
	}

	p.rpool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		MaxActive:   appcfg.GetInt("redis_per_max_conn", 2500),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redis_ipport, redis.DialConnectTimeout(30*time.Second), redis.DialReadTimeout(30*time.Second), redis.DialWriteTimeout(30*time.Second))
			if err != nil {
				return nil, err
			}

			if redisPasswd != "" {
				if _, err = c.Do("AUTH", redisPasswd); err != nil {
					c.Close()
					return nil, err
				}
			}

			if _, err = c.Do("SELECT", reids_db_index); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	p.DB.SetMaxIdleConns(10)

	p.QueryOrderStmt, err = p.DB.Prepare("select diamond_id,order_time,pay_time,status,diamond,code,rmb_cost,channel,info,order_uid,reward from gf_orders where username=?")
	printErr(err)
	p.QuerySingleOrderStmt, err = p.DB.Prepare("select diamond_id,order_time,pay_time,status,diamond,code,rmb_cost,channel,info,order_uid from gf_orders where order_uid=?")
	printErr(err)
	p.UpdateOrderStatusFromStmt, err = p.DB.Prepare("update gf_orders set status=? where order_uid=? and status=?")
	printErr(err)
	p.UpdateOrderStatusAndInfoStmt, err = p.DB.Prepare("update gf_orders set status=?,info=? where order_uid=? and status=?")
	printErr(err)
	p.GetOneOrderStmt, err = p.DB.Prepare("select username,diamond_id,reward from gf_orders where order_uid=?")
	printErr(err)
	p.ListAllOrderStmt, err = p.DB.Prepare("select diamond_id,order_time,pay_time,status,diamond,code,rmb_cost,channel,info,order_uid from gf_orders order by order_id desc limit 10000")
	printErr(err)

	p.GetNicknameStmt, err = p.DB.Prepare("select nickname from gf_nickname where username=?")
	printErr(err)

	p.SetConfigStmt, err = p.DB.Prepare("insert into gf_config values(0,?,?) on duplicate key update value=?")
	printErr(err)
	p.GetConfigStmt, err = p.DB.Prepare("select value from gf_config where `key`=?")
	printErr(err)
	p.GetAllConfigStmt, err = p.DB.Prepare("select `key`,value from gf_config")
	printErr(err)
	p.QueryDesignerStmt, err = p.DBDesigner.Prepare("select points,diamond,gold,`desc` from gf_designer where username=?")
	printErr(err)
	p.UpdateDesignerPointStmt, err = p.DBDesigner.Prepare("update gf_designer set points=? where username=?")
	printErr(err)
	p.UpdateDesignerDiamondStmt, err = p.DBDesigner.Prepare("update gf_designer set diamond=? where username=?")
	printErr(err)
	p.UpdateDesignerDescStmt, err = p.DBDesigner.Prepare("update gf_designer set `desc`=? where username=?")
	printErr(err)
	p.UpdateDesignerIcon, err = p.DBDesigner.Prepare("update gf_designer set head = ? where username=?")
	printErr(err)
	// p.UpdateDesignerStatIcon, err = p.DBDesigner.Prepare("update gf_status set head='head-101.png' where username=?")
	// printErr(err)
	// p.AddBGStmt, err = p.DBDesigner.Prepare("update gf_status set extra1 = ? where username =?")
	// printErr(err)
	// p.DeleteBGStmt, err = p.DBDesigner.Prepare("update gf_status set extra1 = ? where username = ?")
	// printErr(err)
	p.UpdateDesignerGoldStmt, err = p.DBDesigner.Prepare("update gf_designer set gold=? where username=?")
	printErr(err)

	// p.QUsernameStmt, err = p.DB.Prepare("select bilibili_username,coalesce(first_login_time,0) from pixie_player limit ?,?")
	// printErr(err)

	p.AddCosplayStmt, err = p.DBParty.Prepare("insert into gf_cosplay values(0,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	printErr(err)
	p.DeleteCosplayStmt, err = p.DBParty.Prepare("delete from gf_cosplay where cosplay_id=?")
	printErr(err)
	p.GetCosplayStmt, err = p.DBParty.Prepare("select title,keyword,type,cos_params from gf_cosplay where cosplay_id=?")
	printErr(err)

	p.AddMailStmt, err = p.DB.Prepare("insert into gf_mail values(0,?,?,?,?,?,?,?,?,?,?,?,?)")
	printErr(err)
	p.QueryMailStmt, err = p.DB.Prepare("select mail_id,`from`,content,gold,diamond,clothes,time from gf_mail where `to`=? and `read`=0 order by mail_id desc")
	printErr(err)
	p.DeleteMailStmt, err = p.DB.Prepare("delete from gf_mail where mail_id=?")
	printErr(err)
	p.GetAllMailStmt, err = p.DB.Prepare("select * from gf_mail where `to`=? order by mail_id desc limit ?")
	printErr(err)

	p.GetUIDStmt, err = p.DB.Prepare("select last_server,coalesce(db_addr,'') from pixie_player where uid=?")
	printErr(err)
	// p.GetUsernameStmt, err = p.DB.Prepare("select gf_username from pixie_player where uid=?")
	// printErr(err)
	// p.GetUsername1Stmt, err = p.DB.Prepare("select gf_username from pixie_player where username=?")
	// printErr(err)

	p.QueryHelpStmt, err = p.DB.Prepare("select help_id,content,image,post_time from gf_help where username=? order by help_id desc")
	printErr(err)
	p.DeleteHelpStmt, err = p.DB.Prepare("delete from gf_help where help_id=?")
	printErr(err)
	p.QueryGMServerStmt, err = p.DB.Prepare("select address from gf_servers where type=? and address like ? order by heartbeat_time desc limit 1")
	printErr(err)

	p.GetLocalLoginServerStmt, err = p.DB.Prepare("select address from gf_servers where type=? order by heartbeat_time desc limit 1")
	printErr(err)

	p.QUsernameByNicknameStmt, err = p.DB.Prepare("select username from gf_nickname where nickname=?")
	printErr(err)

	// p.UpdateGlobalNicknameStmt, err = p.DB.Prepare("update gf_nickname set nickname=? where username=?")
	// printErr(err)
	// p.InsertGlobalNicknameStmt, err = p.DB.Prepare("insert into gf_nickname values(0,?,?)")
	// printErr(err)

	p.GetSPReportStmt, err = p.DB.Prepare("select ds,report from system_report where report_name=? and ds>=? and ds<=?")
	printErr(err)

	p.GetAddressStmt, err = p.DB.Prepare("select address from gf_servers where local_ip=? and type=? order by heartbeat_time desc limit 1")
	printErr(err)
	p.ListAddressStmt, err = p.DB.Prepare("select distinct address from gf_servers where type=?")
	printErr(err)
	p.GetLocalIPStmt, err = p.DB.Prepare("select local_ip from gf_servers where address=? order by heartbeat_time desc limit 1")
	printErr(err)

	p.AddEventStmt, err = p.DB.Prepare("insert into gf_event values(0,?,?,?,?,?,?,?,?,?,?,?,?)")
	printErr(err)
	p.UpdateEventStmt, err = p.DB.Prepare("update gf_event set type=?,title=?,content=?,start_time=?,end_time=?,redirect_type=?,auto_pop=?,content_type=?,is_long=?,info=? where id=?")
	printErr(err)
	p.StickEventStmt, err = p.DB.Prepare("update gf_event set refresh_time=? where id=?")
	printErr(err)
	p.QueryEventsStmt, err = p.DB.Prepare("select * from gf_event where (start_time<=? and end_time>=?) or is_long=1 order by refresh_time desc")
	printErr(err)
	p.QueryNotOpenEventsStmt, err = p.DB.Prepare("select * from gf_event where start_time>? and is_long=0 order by refresh_time desc")
	printErr(err)
	p.QueryClosedEventsStmt, err = p.DB.Prepare("select * from gf_event where end_time<? and is_long=0 order by refresh_time desc limit 10")
	printErr(err)
	p.DeleteEventStmt, err = p.DB.Prepare("update gf_event set end_time=0,start_time=0,is_long=0 where id=?")
	printErr(err)
	p.GetEventStmt, err = p.DB.Prepare("select * from gf_event where id=?")
	printErr(err)

	p.QueryLogStmt, err = p.DB.Prepare("select * from gf_gmlog where username=? and time>=? and time<=? order by time desc limit 10000")
	printErr(err)
	p.AddPlayerLog1Stmt, err = p.DB.Prepare("insert into gf_player_log1 values(0,?,?,?,?,?,?)")
	printErr(err)
	p.AddPlayerLog2Stmt, err = p.DB.Prepare("insert into gf_player_log2 values(0,?,?,?,?,?,?)")
	printErr(err)
	p.QueryPlayerLog1Stmt, err = p.DB.Prepare("select * from gf_player_log1 where username=? and time>=? and time<=? order by time desc limit 10000")
	printErr(err)
	p.QueryPlayerDetailLog1Stmt, err = p.DB.Prepare("select * from gf_player_log1 where username=? and c1=? and c2=? and time>=? and time<=? order by time desc limit 10000")
	printErr(err)
	//查询韩国版本日志某个时间日志
	p.QueryPlayerDefaultlogStmt, err = p.DB.Prepare("select * from gf_player_log_kor where username=? and c3=? and c2=? and c1=? and time>=? and time<=? order by id desc limit 10000")
	printErr(err)
	//查询某个用户所有日志
	p.QueryPlayerAllLogStmt, err = p.DB.Prepare("select * from gf_player_log_kor where username=? order by time desc limit 10000")
	printErr(err)
	//查询某个用户的愿望
	p.QueryPlayerHopeStmt, err = p.DB.Prepare("select * from gf_player_log_kor where username=? and C2=?")
	printErr(err)

	p.QueryPlayerLog2Stmt, err = p.DB.Prepare("select * from gf_player_log2 where username=? and time>=? and time<=? order by time desc limit 10000")
	printErr(err)
	p.QueryPlayerDetailLog2Stmt, err = p.DB.Prepare("select * from gf_player_log2 where username=? and c1=? and c2=? and time>=? and time<=? order by time desc limit 10000")
	printErr(err)

	p.AddCronMailStmt, err = p.DB.Prepare("insert into gf_cron_mail values(0,?,?,?,?,?,?,?,?,?,?,?)")
	printErr(err)
	p.GetCronMailStmt, err = p.DB.Prepare("select * from gf_cron_mail where time<=?+30000000")
	printErr(err)
	p.UpdateCronMailStmt, err = p.DB.Prepare("update gf_cron_mail set status=?,info=? where id=?")
	printErr(err)
	p.DeleteCronMailStmt, err = p.DB.Prepare("delete from gf_cron_mail where id=?")
	printErr(err)

	//party
	p.ClosePartyStmt, err = p.DB.Prepare("update gf_party set type=?,close_time=? where id=? and close_time>?")
	printErr(err)

	//gift
	p.AddGiftPackStmt, err = p.DB.Prepare("insert into gf_gift_pack values(0,?,?)")
	printErr(err)
	p.DeleteGiftPackStmt, err = p.DB.Prepare("delete from gf_gift_pack where id=?")
	printErr(err)
	p.DeleteCodeStmt, err = p.DB.Prepare("delete from gf_code where gift_pack_id=?")
	printErr(err)
	p.DeleteCode1Stmt, err = p.DB.Prepare("delete from gf_code_1 where gift_pack_id=?")
	printErr(err)
	p.DeleteCode2Stmt, err = p.DB.Prepare("delete from gf_code_2 where gift_pack_id=?")
	printErr(err)
	p.GetGiftPackStmt, err = p.DB.Prepare("select * from gf_gift_pack")
	printErr(err)
	p.AddGiftCodeStmt, err = p.DB.Prepare("insert into gf_code(code,gift_pack_id) values(?,?)")
	printErr(err)
	p.AddGiftCode2Stmt, err = p.DB.Prepare("insert into gf_code_2(code,gift_pack_id) values(?,?)")
	printErr(err)

	p.GetGiftCode2Stmt, err = p.DB.Prepare("select gift_pack_id,coalesce(username,'') from gf_code_2 where code=?")
	printErr(err)
	p.GetGiftCodeListStmt, err = p.DB.Prepare("select code from gf_code where gift_pack_id=? and (username='' or username is null) limit ?,?")
	printErr(err)
	p.GetGift2CodeListStmt, err = p.DB.Prepare("select code from gf_code_2 where gift_pack_id=? and (username='' or username is null) limit ?,?")
	printErr(err)
	p.GetGiftPackContentStmt, err = p.DB.Prepare("select name,content from gf_gift_pack where id=?")
	printErr(err)

	p.QueryGiftCodeStatusStmt, err = p.DB.Prepare("select count(1),sum(case when username='' or username is null then 0 else 1 end) from gf_code where gift_pack_id=?")
	printErr(err)
	p.QueryGift2CodeStatusStmt, err = p.DB.Prepare("select count(1),sum(case when username='' or username is null then 0 else 1 end) from gf_code_2 where gift_pack_id=?")
	printErr(err)

	p.QueryGiftCodeStmt, err = p.DB.Prepare("select gift_pack_id,coalesce(username,'') from gf_code where code=?")
	printErr(err)

	p.UpdateGiftCodeStmt, err = p.DB.Prepare("update gf_code set username=? where code=? and (username='' or username is null)")
	printErr(err)
	p.UpdateGiftCode2Stmt, err = p.DB.Prepare("update gf_code_2 set username=? where code=? and (username='' or username is null)")
	printErr(err)
	p.UseOneCodeStmt, err = p.DB.Prepare("insert into gf_code_1(gift_pack_id,username) values(?,?)")
	printErr(err)

	p.GetGiftPackIDStmt, err = p.DB.Prepare("select gift_pack_id from gf_code where code=? limit 1")
	printErr(err)
	p.GetGiftPackID2Stmt, err = p.DB.Prepare("select gift_pack_id from gf_code_2 where code=? limit 1")
	printErr(err)

	p.QueryGiftCodeExistStmt, err = p.DB.Prepare("select gift_pack_id,coalesce(username,'') from gf_code where code=?")
	printErr(err)
	p.QueryGiftCode2ExistStmt, err = p.DB.Prepare("select gift_pack_id,coalesce(username,'') from gf_code_2 where code=?")
	printErr(err)

	p.UpdateVersionStmt, err = p.DB.Prepare("insert into version(id,game,core_version,download_url,download_url2,script_version,version_file) values(0,?,?,?,?,?,?) on duplicate key update core_version=?,download_url=?,download_url2=?,script_version=?,version_file=?")
	printErr(err)

	//更新时空之门
	p.UpsertShiKongzhimenStmt, err = p.DB.Prepare("insert into gf_config values(0,?,?) on duplicate key update `value`=?")
	printErr(err)
	//获取时空之门
	p.GetShikongzhimenStmt, err = p.DB.Prepare("select * from  gf_config  where `key`=?")
	printErr(err)
	//删除时空之门活动
	p.DelShiKongzhimenStmt, err = p.DB.Prepare("delete from gf_config where `key`= ?")
	printErr(err)
	//插入时空之门参数
	p.InsertShikongzhimenStmt, err = p.DB.Prepare("insert into gf_config values(0,?,?)")
	printErr(err)

	p.QueryVersionStmt, err = p.DB.Prepare("select game,core_version,download_url,download_url2,script_version,version_file,tishen from version")
	printErr(err)

	p.UpdateTishenStmt, err = p.DB.Prepare("update version set tishen=?")
	printErr(err)

	p.GetRollingAnnouncementStmt, err = p.DB.Prepare("select * from gf_rolling_info")
	printErr(err)
	p.AddRollingAnnouncementStmt, err = p.DB.Prepare("insert into gf_rolling_info values(0,?,?,?)")
	printErr(err)
	p.DelRollingAnnouncementStmt, err = p.DB.Prepare("delete from gf_rolling_info where id=?")
	printErr(err)
	p.ChangeRollingAnnouncementStmt, err = p.DB.Prepare("update gf_rolling_info set content=?,start_time=?,end_time=? where id=?")
	printErr(err)

	p.AddPushStmt, err = p.DB.Prepare("insert into gf_push_list values(0,?,?,?,?,?)")
	printErr(err)
	p.DelPushStmt, err = p.DB.Prepare("delete from gf_push_list where id=?")
	printErr(err)
	p.ChangePushStmt, err = p.DB.Prepare("update gf_push_list set content=?,push_time=?,`to`=? where id=?")
	printErr(err)
	p.GetPushStmt, err = p.DB.Prepare("select * from gf_push_list")
	printErr(err)

	p.AddMaintainJobStmt, err = p.DB.Prepare("insert into gf_maintain_job values(0,?,?,?,?,?)")
	printErr(err)
	p.DelMaintainJobStmt, err = p.DB.Prepare("delete from gf_maintain_job where id=?")
	printErr(err)
	p.ChangeMaintainJobStmt, err = p.DB.Prepare("update gf_maintain_job set content=?,start_time=?,end_time=?,url=?,show_time=? where id=?")
	printErr(err)
	p.GetMaintainJobStmt, err = p.DB.Prepare("select id,content,start_time,end_time,url,show_time from gf_maintain_job")
	printErr(err)

	p.GetOnlinePlayerCntStmt, err = p.DB.Prepare("select token,player_cnt,`desc` from gf_online_player_cnt where token>=? and token<=?")
	printErr(err)
	p.QueryServerHBStmt, err = p.DB.Prepare("select id,address,type,heartbeat_time from gf_servers")
	printErr(err)
	p.QueryGuildInfoStmt, err = p.DB.Prepare("select gid,activity,clothes_share_cnt from gf_guild_member where username=? and zombie=0 order by id desc limit 1")
	printErr(err)

	p.UpdateGuildMemberActivityStmt, err = p.DBGuild.Prepare("update gf_guild_member set activity=? where gid=? and username=?")
	printErr(err)
	p.UpdateGuildShareStmt, err = p.DBGuild.Prepare("update gf_guild_member set clothes_share_cnt=? where gid=? and username=?")
	printErr(err)
	// p.QBusernameStmt, err = p.DB.Prepare("select username from pixie_player where gf_username=?")
	// printErr(err)
	p.UpdateGuildDescStmt, err = p.DBGuild.Prepare("update gf_guild set description=? where id=?")
	printErr(err)
	p.UpdateGuildActivityStmt, err = p.DBGuild.Prepare("update gf_guild set activity=? where id=?")
	printErr(err)
	p.UpdateGuildMedalStmt, err = p.DBGuild.Prepare("update gf_guild set medal_cnt=? where id=?")
	printErr(err)

	p.ListHopeStmt, err = p.DB.Prepare("select id,offer_clothes,offer_clo_part,offer_num,hope_clothes,hope_clo_part,hope_num,send_time,status,helper,helper_nickname from gf_hope where username=? order by id desc limit 10")
	printErr(err)
	p.ListBoardMsgStmt, err = p.DB.Prepare("select id,author,author_nickname,reply_to,content,time from gf_board where owner=? order by id desc limit 100")
	printErr(err)
	p.DelHopeStmt, err = p.DBFriend.Prepare("delete from gf_hope where id=?")
	printErr(err)
	p.DelBoardMsgStmt, err = p.DBFriend.Prepare("delete from gf_board where id=?")
	printErr(err)

	//kor
	p.QueryKorLogStmt, err = p.DB.Prepare("select * from gf_player_log_kor where username=? and time>=? and time<=? order by time desc limit 10000")
	printErr(err)
	p.QueryKorLogDetailStmt, err = p.DB.Prepare("select * from gf_player_log_kor where username=? and c1=? and c2=? and time>=? and time<=? order by time desc limit 10000")
	printErr(err)
	p.QueryKorLogNoUsernameDetailStmt, err = p.DB.Prepare("select * from gf_player_log_kor where c1=? and c2=? and time>=? and time<=? order by time desc limit 10000")
	printErr(err)
	p.AddKorLogStmt, err = p.DB.Prepare("insert into gf_player_log_kor values(0,?,?,?,?,?,?)")
	printErr(err)
	p.DelUnregisterStmt, err = p.DB.Prepare("delete from gf_kor_unregister where username=?")
	printErr(err)
	p.FindGuildStmt, err = p.DB.Prepare("select id,owner,owner_nickname,name,description,activity,war_date,medal_cnt from gf_guild where (id=? or name=?) and zombie=0")
	printErr(err)
	p.ListGuildMemberStmt, err = p.DB.Prepare("select username,nickname,vip,level,activity,ignore_war from gf_guild_member where gid=? and zombie=0")
	printErr(err)
	p.UpdateGuildOwnerStmt, err = p.DB.Prepare("update gf_guild set owner=?,owner_nickname=? where id=? and owner=?")
	printErr(err)

	// //QueryKorLoginStmt 查询韩国login
	// p.QueryKorLoginStmt, err = p.DB.Prepare("select * from gf_login_kor_csv limit ?,?")
	// printErr(err)
	// //QueryKorLogin1Stmt 查询韩国login1
	// p.QueryKorLogin1Stmt, err = p.DB.Prepare("select * from gf_login1_kor_csv where ds >= ? and ds <= ? limit ?,?")
	// printErr(err)
	// //QueryKorOrderStmt 查询韩国订单
	// p.QueryKorOrderStmt, err = p.DB.Prepare("select * from gf_order_kor_csv where ds >= ? and ds <= ? limit ?,?")
	// printErr(err)

	//QueryKorGameCode1 查询韩国游戏码1
	p.QueryKorGameCode1, err = p.DB.Prepare("select * from code1 where username=?")
	printErr(err)

	//QueryKorGameCode2 查询韩国游戏码2
	p.QueryKorGameCode2, err = p.DB.Prepare("select * from code2 where username=?")
	printErr(err)

	//QueryKorGameCode2 查询韩国游戏码3
	p.QueryKorGameCode3, err = p.DB.Prepare("select * from code3 where username=?")
	printErr(err)

	//QueryKorGameCode1 查询韩国游戏码1
	p.QueryKorGameCode1ByCode, err = p.DB.Prepare("select * from code1 where code=?")
	printErr(err)

	//QueryKorGameCode2 查询韩国游戏码2
	p.QueryKorGameCode2ByCode, err = p.DB.Prepare("select * from code2 where code=?")
	printErr(err)

	//QueryKorGameCode2 查询韩国游戏码3
	p.QueryKorGameCode3ByCode, err = p.DB.Prepare("select * from code3 where code=?")
	printErr(err)

	// p.QueryKorlogUploadProgress, err = p.DB.Prepare("select * from gf_kor_log_progress where ds= ?")
	// printErr(err)

	// p.AddKorLogProgress, err = p.DB.Prepare("insert into gf_kor_log_progress values(0,?,?) ")
	// printErr(err)
	//party
	p.AddPartyStmt, err = p.DBParty.Prepare("insert into gf_party values(0,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	printErr(err)
	p.AddPartyItemStmt, err = p.DBParty.Prepare("insert into gf_party_item(id,party_id,username,nickname,head,img,like_cnt,unlike_cnt,upload_time) values(0,?,?,?,?,?,?,?,?)")
	printErr(err)
	p.AddPartyItemFlowerStmt, err = p.DBParty.Prepare("insert into gf_party_flower values(0,?,?,?,?)")
	printErr(err)
	p.AddPartyItemCmtStmt, err = p.DBParty.Prepare("insert into gf_party_cmt values(0,?,?,?,?,?,?)")
	printErr(err)
	p.DeletePartyStmt, err = p.DBParty.Prepare("delete from gf_party where id=?")
	printErr(err)
	p.DeletePartyItemStmt, err = p.DBParty.Prepare("delete from gf_party_item where party_id=?")
	printErr(err)
	p.DeletePartyCmtStmt, err = p.DBParty.Prepare("delete from gf_party_cmt where party_id=?")
	printErr(err)
	p.DeletePartyFlowerStmt, err = p.DBParty.Prepare("delete from gf_party_flower where party_id=?")
	printErr(err)

	//init sharding dbs
	// var dbConfigStr string
	// if dbConfigStr, err = p.getDBConfigStr(); err != nil || dbConfigStr == "" {
	// 	panic(err)
	// }
	// if p.DBMap, p.dcs, err = getDBMap(dbConfigStr, dbUsername, dbPasswd, designDBIPPort, friendDBIPPort, guildDBIPPort, partyDBIPPort); err != nil {
	// 	panic(err)
	// }
	return
}

func (self *Platform) GetRConn() (c redis.Conn) {
	return self.rpool.Get()
}

// func (self *Platform) GetIPPort(uid int64) string {
// 	for _, v := range self.dcs {
// 		if uid >= v.Start && uid <= v.End {
// 			value := uid % 2
// 			if value == 0 {
// 				return v.IPPort
// 			} else {
// 				return v.IPPort2
// 			}
// 		}
// 	}
// 	fmt.Println("can't find db for user:", uid)
// 	return ""
// }

// func (self *Platform) getDBConfigStr() (str string, err error) {
// 	if err = self.DB.QueryRow("select value from gf_config where `key`=?", DB_SHARDING_KEY).Scan(&str); err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	return
// }
