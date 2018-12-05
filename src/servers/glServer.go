package servers

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
)

import (
	"appcfg"
	"constants"
	"gl_dao"
	"gl_handle/bilibiliHandle"
	"gl_handle/third"
	"gl_handle/userHandle"
	"gl_handle/versionHandle"
	_ "hub/gameServerHub"
	. "logger"
	"time"
)

var (
	gl_handlers map[string]func(http.ResponseWriter, *http.Request)

	ts_gl_handlers map[string]func(http.ResponseWriter, *http.Request)
)

type OtkHandle struct{}
type OtkTsHandle struct{}

func (self *OtkHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handle, ok := gl_handlers[r.URL.String()]; ok {
		handle(w, r)
	} else {
		Err("unknown api call", r.URL.String(), r.URL.Path)
	}
}

func (self *OtkTsHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handle, ok := ts_gl_handlers[r.URL.String()]; ok {
		handle(w, r)
	}
}

func init() {
	if appcfg.GetServerType() != constants.SERVER_TYPE_GL && appcfg.GetServerType() != constants.SERVER_TYPE_GL_TS {
		return
	}

	gl_handlers = make(map[string]func(http.ResponseWriter, *http.Request))
	ts_gl_handlers = make(map[string]func(http.ResponseWriter, *http.Request))

	gl_handlers["/OtkServer/api/Login"] = userHandle.UserLogin
	gl_handlers["/OtkServer/api/Register"] = userHandle.UserRegister
	// gl_handlers["/OtkServer/api/editPassword"] = userHandle.UserEditPassword

	gl_handlers["/OtkServer/api/checkVersion"] = versionHandle.CheckVersion
	// gl_handlers["/OtkServer/updateGFVersion"] = versionHandle.UpdateVersion

	gl_handlers["/OtkServer/api/biliLogin"] = bilibiliHandle.BilibiliLogin
	gl_handlers["/OtkServer/api/biliTestLogin"] = bilibiliHandle.BilibiliTestLogin
	gl_handlers["/OtkServer/api/thirdLogin"] = third.ThirdLogin
	gl_handlers["/OtkServer/test"] = TestHttpDispatch

	ts_gl_handlers["/OtkServer/api/checkVersion"] = versionHandle.CheckVersionTS
	ts_gl_handlers["/OtkServer/test"] = TestHttpDispatch
}

func StartGLServer(ts bool) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	addr := appcfg.GetString("bind_ip", "127.0.0.1") + appcfg.GetString("port", ":29980")

	server := http.Server{
		Addr:         addr,
		ReadTimeout:  constants.TIME_OUT * time.Second,
		WriteTimeout: constants.TIME_OUT * time.Second,
	}
	server.Handler = &OtkHandle{}
	if ts {
		server.Handler = &OtkTsHandle{}
	}

	//heartbeat to db.mark itself working ok.
	if !ts {
		go heartbeat()
	}

	fmt.Println("Server start up OK.")
	var err error
	if appcfg.GetBool("https", false) {
		cerFile := appcfg.GetString("https_cer_file", "")
		keyFile := appcfg.GetString("https_key_file", "")
		err = server.ListenAndServeTLS(cerFile, keyFile)
	} else {
		err = server.ListenAndServe()
	}

	if err != nil {
		fmt.Println("Server start up error:", err)
		return
	}
}

func heartbeat() {
	for {
		var err error

		//local gm server seen
		localAddress := appcfg.GetString("local_ip", "127.0.0.1") + appcfg.GetString("port", "29980")
		//local game server seen
		localForGameAddress := appcfg.GetString("local_ip", "127.0.0.1") + appcfg.GetString("gs_port", "29981")
		now := time.Now().Unix()

		localID := fmt.Sprintf("%s-%d", appcfg.GetString("local_ip", "127.0.0.1"), os.Getpid())

		if _, err = gl_dao.HeatbeatStmt.Exec(localID, localAddress, constants.LOGIN_FOR_LOCAL_SERVER_TYPE, now, constants.LOGIN_FOR_LOCAL_SERVER_TYPE, now); err != nil {
			Err(err)
			//do not panic or break.just log this err msg
		}

		if _, err = gl_dao.HeatbeatStmt.Exec(localID, localForGameAddress, constants.LOGIN_FOR_GAME_SERVER_TYPE, now, constants.LOGIN_FOR_GAME_SERVER_TYPE, now); err != nil {
			Err(err)
			//do not panic or break.just log this err msg
		}

		if _, err = gl_dao.DeleteServerStmt.Exec(now - 2*constants.HEARTBEAT_TIME); err != nil {
			Err(err)
		}

		time.Sleep(time.Duration(constants.HEARTBEAT_TIME) * time.Second)
	}
}

func TestHttpDispatch(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
