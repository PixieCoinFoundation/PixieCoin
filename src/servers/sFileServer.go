package servers

import (
	"flag"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

import (
	// "appcfg"
	. "logger"
)

type Sizer interface {
	Size() int64
}

var addr = flag.String("a", ":29951", "serve address")

const (
	LOCAL_PATH = "./file_server_files/"
)

func StartSFileServer() {
	flag.Parse()

	// 配置最大的go thread
	runtime.GOMAXPROCS(runtime.NumCPU()) //running in multicore

	// 页面--------------------------------------------
	// 管理员页面
	http.HandleFunc("/upload", uploadHandle)

	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir(LOCAL_PATH))))

	// addr := appcfg.GetString("bind_ip", "127.0.0.1") + appcfg.GetString("http_port", ":29971")

	Info("Server start listen at:" + *addr)
	// err := http.ListenAndServeTLS(addr, "cert/gfcert/gf.crt", "cert/gfcert/gf.key", nil)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		Err("Server start up error:", err)
		return
	}
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

func uploadHandle(w http.ResponseWriter, r *http.Request) {
	if "POST" == r.Method {
		file, _, err := r.FormFile("file_content")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer file.Close()

		file_name := r.FormValue("file_name")
		dc := strings.Count(file_name, "/")
		if dc > 1 {
			http.Error(w, "file_name illegal", 500)
			return
		} else if dc == 1 {
			dn := strings.Split(file_name, "/")[0]
			os.MkdirAll(LOCAL_PATH+dn, os.ModePerm)
		}

		f, err := os.Create(LOCAL_PATH + file_name)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		// fmt.Fprintf(w, "上传文件的大小为: %d", file.(Sizer).Size())
		Info("uploaded file:%s,size:%d", file_name, file.(Sizer).Size())
		return
	}

	http.Error(w, "error", 500)
}
