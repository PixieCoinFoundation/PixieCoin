package servers

import (
	"dao"
	"db/global"
	"flag"
	"fmt"
	. "logger"
	"net/http"
	"os"
	"runtime"
	"shutdown"
	"strings"
	"time"
)

import (
	"appcfg"
	"constants"
	"dispatch_pixie"
	// "handle/gmHandle"
	_ "hub/loginServerHub"
	cl "service/clean"
)

const (
	HOT_CUSTOM_PREFIX = "/GFServer/api/hot_custom_page/"
	SALE_RANK_PREFIX  = "/GFServer/api/sale_rank_page/"
)

var (
	handlers    map[string]func(http.ResponseWriter, *http.Request)
	iapHandlers map[string]func(http.ResponseWriter, *http.Request)
	// dataHandlers map[string]func(http.ResponseWriter, *http.Request)
	gmHandlers map[string]func(http.ResponseWriter, *http.Request)

	memprofile = flag.String("memprofile", "", "write memory profile to this file")
)

func init() {
	if appcfg.GetServerType() != "" {
		return
	}

	// 启动http游戏逻辑监听
	handlers = make(map[string]func(http.ResponseWriter, *http.Request))
	for api, _ := range dispatch_pixie.PixiePlayerRequestHandles {
		handlers["/game/"+api] = dispatch_pixie.PixieApiGlobalEntrance
	}
	handlers["/game/api/loginGFServer"] = dispatch_pixie.PixieApiGlobalEntrance

	handlers["/game/test"] = dispatch_pixie.TestApiEntrance

	// if appcfg.GetLanguage() == "" {
	// 	handlers["/GFServer/api/qLotteryLog"] = dataHandle.QueryLotteryLogHandle
	// 	//iap callbacks
	// 	//bili
	// 	handlers["/GFServer/api/biliwebhook"] = iapHandle.BiliWebhook
	// 	//4399
	// 	handlers["/GFServer/api/tPayCallback"] = iapHandle.Third4399Webhook
	// 	//huawei
	// 	handlers["/GFServer/api/hwPayCallback"] = iapHandle.ThirdHuaweiWebhook
	// 	//vivo
	// 	handlers["/GFServer/api/vivoPayCallback"] = iapHandle.ThirdVivoWebhook
	// 	//meitu
	// 	handlers["/GFServer/api/mtPayCallback"] = iapHandle.ThirdMeituWebhook
	// 	//papa
	// 	handlers["/GFServer/api/papaPayCallback"] = iapHandle.ThirdPapaWebhook
	// 	//kuaikan
	// 	handlers["/GFServer/api/kuaikanPayCallback"] = iapHandle.ThirdKuaikanWebhook
	// 	//buka
	// 	handlers["/GFServer/api/bukaPayCallback"] = iapHandle.ThirdBukaWebhook
	// 	handlers["/o/GFServer/api/bukaPayCallback"] = iapHandle.ThirdBukaWebhook
	// 	//uc
	// 	handlers["/GFServer/api/ucPayCallback"] = iapHandle.ThirdUCWebhook
	// 	//xiaomi
	// 	handlers["/GFServer/api/xiaomiPayCallback"] = iapHandle.ThirdXiaomiWebhook
	// 	//sina
	// 	handlers["/GFServer/api/sinaPayCallback"] = iapHandle.ThirdSinaWebhook
	// 	//360
	// 	handlers["/GFServer/api/360PayCallback"] = iapHandle.Third360Webhook
	// 	//oppo
	// 	handlers["/GFServer/api/oppoPayCallback"] = iapHandle.ThirdOppoWebhook
	// 	//mhr
	// 	handlers["/GFServer/api/mhrPayCallback"] = iapHandle.ThirdMHRWebhook
	// 	handlers["/GFServer/api/mhrGenOrder"] = iapHandle.MHRGenOrderHandle
	// 	//baidu
	// 	handlers["/GFServer/api/baiduPayCallback"] = iapHandle.ThirdBaiduWebhook

	// 	handlers["/GFServer/api/getHotCustom"] = customHandle.GetHotCustomHandle
	// 	handlers["/GFServer/api/shareClickLog"] = dataHandle.ShareClickLogHandle
	// 	handlers["/GFServer/api/getSaleRank"] = customHandle.GetTempSaleRankHandle
	// 	handlers["/GFServer/api/getUserSaleRank"] = customHandle.GetTempSaleUserRankHandle

	// 	hd := http.StripPrefix(HOT_CUSTOM_PREFIX, http.FileServer(http.Dir("./hot_custom_page/")))
	// 	handlers[HOT_CUSTOM_PREFIX] = hd.ServeHTTP

	// 	hds := http.StripPrefix(SALE_RANK_PREFIX, http.FileServer(http.Dir("./sale_rank_page/")))
	// 	handlers[SALE_RANK_PREFIX] = hds.ServeHTTP
	// }

	// iapHandlers = make(map[string]func(http.ResponseWriter, *http.Request))
	// iapHandlers["/GFServer/api/pingpppay"] = iapHandle.PingPPPayHandle
	// iapHandlers["/GFServer/api/pingppwebhook"] = iapHandle.PingPPWebhook
	//bili
	// iapHandlers["/GFServer/api/biliwebhook"] = iapHandle.BiliWebhook

	gmHandlers = make(map[string]func(http.ResponseWriter, *http.Request))
	// gmHandlers["/game/gm/api/lockPlayer"] = gmHandle.LockPlayer
	// gmHandlers["/GFServer/gm/api/unlockPlayer"] = gmHandle.UnlockPlayer
}

func StartGameServer() {
	runtime.GOMAXPROCS(runtime.NumCPU()) //running in multicore

	// startServer("iap server", appcfg.GetString("bind_ip", "127.0.0.1")+appcfg.GetString("pay_port", ":29993"), "", "", IAPHandler, true)
	// startServer("data server", appcfg.GetString("data_port", ":29994"), "", "", DataHandler, true)
	startServer("game server", appcfg.GetString("gm_port", "29995"), "", "", GMHandler, true)

	go selfHeartbeat()
	go clean()

	checkInitKorCrossGameReward()

	var cerFile, keyFile string
	if appcfg.GetBool("https", false) {
		cerFile = appcfg.GetString("https_cer_file", "")
		keyFile = appcfg.GetString("https_key_file", "")
	}
	startServer("player server", appcfg.GetString("bind_ip", "127.0.0.1")+appcfg.GetString("port", ":29991"), cerFile, keyFile, PixieHandler, false)
}

func checkInitKorCrossGameReward() {
	if appcfg.GetServerType() == "" && appcfg.GetBool("main_server", false) && appcfg.GetLanguage() == constants.KOR_LANGUAGE && appcfg.GetBool("kor_game_package_code_available", false) {
		//查询是否已经初始化
		row := dao.QueryKorCodeIsInit.QueryRow()
		var id int
		var status bool
		row.Scan(&id, &status)
		//如果已经初始化，直接跳过，如果没有初始化初始化
		if !status {
			//初始礼品包TODO
			err := global.GetCodeInitClear()
			if err != nil {
				panic(err)
			}
			global.GetCodeInitFromFile("scripts/kor/170926_EveryTown.txt", "scripts/kor/170926_VikingIsland.txt", "scripts/kor/170926_Restaurant.txt")
			if _, err := dao.ExecKorCodeInit.Exec(true, true); err != nil {
				Err(err)
			}
		}
	}
}

func selfHeartbeat() {
	for {
		var err error

		//ALL USED FOR GM
		localGMAddress := appcfg.GetString("local_ip", "127.0.0.1") + appcfg.GetString("gm_port", "29995")
		playerGSAddress := appcfg.GetString("player_ip", "127.0.0.1:29991")

		localID := fmt.Sprintf("%s-%d", appcfg.GetString("local_ip", "127.0.0.1"), os.Getpid())

		now := time.Now().Unix()

		if _, err = dao.HeartbeatStmt.Exec(localID, localGMAddress, constants.GAME_LOCAL_FOR_GM_SERVER_TYPE, now, constants.GAME_LOCAL_FOR_GM_SERVER_TYPE, now); err != nil {
			Err(err)
			//do not panic or break.just log this err msg
		}
		if _, err = dao.HeartbeatStmt.Exec(localID, playerGSAddress, constants.GAME_GLOBAL_FOR_PLAYER_SERVER_TYPE, now, constants.GAME_GLOBAL_FOR_PLAYER_SERVER_TYPE, now); err != nil {
			Err(err)
			//do not panic or break.just log this err msg
		}

		// if _, err = dao.DeleteServerStmt.Exec(now - 2*constants.HEARTBEAT_TIME); err != nil {
		// 	Err(err)
		// }

		time.Sleep(time.Duration(constants.HEARTBEAT_TIME) * time.Second)
	}
}

func clean() {
	if appcfg.GetBool("main_server", false) {
		for {
			clean1()
			time.Sleep(60 * time.Second)
		}
	}
}

func clean1() {
	shutdown.AddOneRequest("clean")
	defer shutdown.DoneOneRequest("clean")

	Info("begin clean expired data..")
	// cl.CleanHelp()
	// cl.CleanCustom()
	cl.CleanCosplay()
	// cl.CleanBoard()
	cl.CleanHope()
	cl.CleanMail()
	cl.CleanApplyGuild()
	// cl.CleanDesignMsg()
	// cl.CleanCustomBuy()
	Info("clean expired data done.")
}

func startServer(name, port, certFile, keyFile string, handler func(w http.ResponseWriter, r *http.Request), async bool) {
	server := &http.Server{
		Addr:         port,
		Handler:      http.HandlerFunc(handler),
		ReadTimeout:  constants.TIME_OUT * time.Second,
		WriteTimeout: constants.TIME_OUT * time.Second,
	}

	fmt.Println(name, "begin to serve..")

	if async {
		if certFile != "" && keyFile != "" {
			go listenTLS(name, server, certFile, keyFile)
		} else {
			go listen(name, server)
		}

	} else {
		if certFile != "" && keyFile != "" {
			listenTLS(name, server, certFile, keyFile)
		} else {
			listen(name, server)
		}

	}
}

func listen(name string, server *http.Server) {
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(name, "serve error:", err)
		return
	}
}

func listenTLS(name string, server *http.Server, certFile, keyFile string) {
	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		fmt.Println(name, "serve error:", err)
		return
	}
}

func PixieHandler(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if strings.HasPrefix(p, HOT_CUSTOM_PREFIX) {
		p = HOT_CUSTOM_PREFIX
	} else if strings.HasPrefix(p, SALE_RANK_PREFIX) {
		p = SALE_RANK_PREFIX
	}
	if handle, ok := handlers[p]; ok {
		handle(w, r)
	} else {
		dispatch_pixie.NotSupportApiHandle(w, r)
	}
}

func GMHandler(w http.ResponseWriter, r *http.Request) {
	if handle, ok := gmHandlers[r.URL.String()]; ok {
		handle(w, r)
	}
}
