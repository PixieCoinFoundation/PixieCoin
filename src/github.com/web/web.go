// Package web is a lightweight web framework for Go. It's ideal for
// writing simple, performant backend web services.
package web

import (
	"common"
	"constants"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	gmc "gm/common"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"reflect"
	file "service/files"
	"strconv"
	"strings"
	"time"
)

// A Context object is created for every incoming HTTP request, and is
// passed to handlers as an optional first argument. It provides information
// about the request, including the http.Request object, the GET and POST params,
// and acts as a Writer for the response.
type AuthInfo struct {
	accessToken string
	roleID      int
}

type Context struct {
	Request  *http.Request
	Params   map[string]string
	tagList  []string
	AuthInfo AuthInfo
	Server   *Server
	http.ResponseWriter
}

func NewAuthInfo() *AuthInfo {
	return &AuthInfo{
		roleID: 0,
	}
}

func (ctx *Context) CheckToken() (bool, int64) {
	cook := ctx.Request.Cookies()
	if len(cook) < 1 {
		return false, 0
	}
	id, _ := strconv.ParseInt(strings.Split(cook[0].Name, "-")[1], 10, 64)
	return cook[0].Value == Md5([]byte(strings.Split(cook[0].Name, "-")[0]+gmc.GetTokenDate()+"otakugame")), id
}

func Md5(ori []byte) string {
	return fmt.Sprintf("%x", md5.Sum(ori))
}

func (ctx *Context) GetTagList() []string {
	return ctx.tagList
}

// WriteString writes string data into the response object.
func (ctx *Context) WriteString(content string) {
	ctx.ResponseWriter.Write([]byte(content))
}

// func (ctx *Context) GetBody() ([]byte, error) {
// 	return ioutil.ReadAll(ctx.Request.Body)
// }

func (ctx *Context) GetFormData(structs interface{}) error {
	ctx.Request.ParseMultipartForm((1 << 20) * 10)
	ctx.Header().Set("Access-Control-Allow-Origin", "*")

	medMap := make(map[string]interface{})
	for k, v := range ctx.Request.PostForm {
		if strings.HasPrefix(k, "tag") {
			ctx.tagList = append(ctx.tagList, v[0])
		}
		medMap[k] = v[0]
	}
	data, err := json.Marshal(medMap)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, structs); err != nil {
		return err
	}
	return nil
}

// func (ctx *Context) GetBodyToStruct(structs interface{}) error {
// 	data, err := ctx.GetBody()
// 	if err == nil {
// 		if err = json.Unmarshal(data, structs); err == nil {
// 			return nil
// 		}
// 	}
// 	return err
// }

func (ctx *Context) UploadFormFile(fileName, saveName string, appendSize bool) (string, error) {
	var resName string

	_, handler, err := ctx.Request.FormFile(fileName)
	if err != nil {
		return "", err
	}

	if handler == nil || handler.Size < 1 {
		return "", constants.UploadFileErr
	}

	File, _ := handler.Open()
	tmpFile, _ := ioutil.TempFile(os.TempDir(), "tmp_")
	defer os.Remove(tmpFile.Name())
	io.Copy(tmpFile, File)

	sn := fmt.Sprintf("%s/%s", file.GetPlatformOSSDirName(), saveName)

	if appendSize {
		sn += fmt.Sprintf("_%d", handler.Size)
	}

	resName, err = file.EncryptUploadFile1(tmpFile, sn)
	if err != nil {
		return "", err
	}

	return resName, nil
}

func (ctx *Context) Unzip(fileName string) (string, error) {
	_, handler, err := ctx.Request.FormFile(fileName)
	if err != nil {
		return "", err
	}
	if handler == nil || handler.Size < 1 {
		return "", constants.UploadFileErr
	}
	var saveDir string
	var tPath []string
	File, _ := handler.Open()
	tmpFile, _ := ioutil.TempFile(os.TempDir(), "tmp_")
	defer os.Remove(tmpFile.Name())
	io.Copy(tmpFile, File)
	randFile := time.Now().Unix()
	name := fmt.Sprintf("%d.zip", randFile)

	if _, err := common.Execute("cp", tmpFile.Name(), name); err != nil {
		return "", err
	}
	if _, err := common.Execute("unzip", "-u", name); err != nil {
		return "", err
	}
	tempDir := handler.Filename[:len(handler.Filename)-4]
	if res, err := common.Execute("find", fmt.Sprintf("./%s", tempDir), "-name", "*.html"); err != nil {
		return "", err
	} else if res == "" {
		return "", errors.New("未找到文件")
	} else {
		res = strings.Replace(res, "\n", "", -1)
		tPath = strings.Split(res, "/")

		saveDir = fmt.Sprintf("%d_%s", time.Now().Unix(), tPath[len(tPath)-1])
	}
	if _, err := common.Execute("rm", "-rf", "__MACOSX"); err != nil {
		return "", err
	}
	if _, err := common.Execute("rm", "-rf", name); err != nil {
		return "", err
	}
	saveSuff := saveDir[:len(saveDir)-6]
	if _, err := common.Execute("mv", tempDir, fmt.Sprintf("./gm/zipFile/%s", saveSuff)); err != nil {
		return "", err
	}
	savePath := "/zipFile/" + saveSuff + "/" + tPath[2]
	return savePath, nil
}

// Abort is a helper method that sends an HTTP header and an optional
// body. It is useful for returning 4xx or 5xx errors.
// Once it has been called, any return value from the handler will
// not be written to the response.
func (ctx *Context) Abort(status int, body string) {
	ctx.SetHeader("Content-Type", "text/html; charset=utf-8", true)
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write([]byte(body))
}

// Redirect is a helper method for 3xx redirects.
func (ctx *Context) Redirect(status int, url_ string) {
	ctx.ResponseWriter.Header().Set("Location", url_)
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write([]byte("Redirecting to: " + url_))
}

//BadRequest writes a 400 HTTP response
func (ctx *Context) BadRequest() {
	ctx.ResponseWriter.WriteHeader(400)
}

func (ctx *Context) BackError(code int, msg interface{}) {
	restResp := NewRestData()
	restResp.SetData(map[string]interface{}{"errmsg": msg, "errcode": code}).SetStatus(400)
	ctx.WriteRest(restResp)
}

func (ctx *Context) BackSuccess(data interface{}) {
	restResp := NewRestData()
	if v, ok := data.(string); ok {
		restResp.SetData(map[string]interface{}{"errmsg": v, "errcode": 200}).SetStatus(200)
	} else {
		restResp.SetData(data).SetStatus(200)
	}
	ctx.WriteRest(restResp)
}

func (ctx *Context) WriteRest(data *restData) {
	//设置统一为json返回
	ctx.ContentType("application/json")
	ctx.ResponseWriter.WriteHeader(data.Status())

	ctx.WriteJson(data.Data())
}

func (ctx *Context) WriteJson(data interface{}, closego ...bool) {

	json, err := json.Marshal(data)
	if err != nil {
		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 当请求为jsonp时
	if callback := ctx.Request.Form.Get("callback"); callback != "" {
		json = append([]byte(callback+"("), json...)
		json = append(json, ')')
		ctx.ResponseWriter.Write(json)
	} else {
		ctx.ResponseWriter.Write(json)
	}
	return
}

// Notmodified writes a 304 HTTP response
func (ctx *Context) NotModified() {
	ctx.ResponseWriter.WriteHeader(304)
}

//Unauthorized writes a 401 HTTP response
func (ctx *Context) Unauthorized() {
	ctx.ResponseWriter.WriteHeader(401)
}

//Forbidden writes a 403 HTTP response
func (ctx *Context) Forbidden() {
	ctx.ResponseWriter.WriteHeader(403)
}

// NotFound writes a 404 HTTP response
func (ctx *Context) NotFound(message string) {
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte(message))
}

// ContentType sets the Content-Type header for an HTTP response.
// For example, ctx.ContentType("json") sets the content-type to "application/json"
// If the supplied value contains a slash (/) it is set as the Content-Type
// verbatim. The return value is the content type as it was
// set, or an empty string if none was found.
func (ctx *Context) ContentType(val string) string {
	var ctype string
	if strings.ContainsRune(val, '/') {
		ctype = val
	} else {
		if !strings.HasPrefix(val, ".") {
			val = "." + val
		}
		ctype = mime.TypeByExtension(val)
	}
	if ctype != "" {
		ctx.Header().Set("Content-Type", ctype)
	}
	return ctype
}

// SetHeader sets a response header. If `unique` is true, the current value
// of that header will be overwritten . If false, it will be appended.
func (ctx *Context) SetHeader(hdr string, val string, unique bool) {
	if unique {
		ctx.Header().Set(hdr, val)
	} else {
		ctx.Header().Add(hdr, val)
	}
}

// SetCookie adds a cookie header to the response.
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	ctx.SetHeader("Set-Cookie", cookie.String(), false)
}

// small optimization: cache the context type instead of repeteadly calling reflect.Typeof
var contextType reflect.Type

var defaultStaticDirs []string

func init() {
	contextType = reflect.TypeOf(Context{})
	//find the location of the exe file
	wd, _ := os.Getwd()
	arg0 := path.Clean(os.Args[0])
	var exeFile string
	if strings.HasPrefix(arg0, "/") {
		exeFile = arg0
	} else {
		//TODO for robustness, search each directory in $PATH
		exeFile = path.Join(wd, arg0)
	}
	parent, _ := path.Split(exeFile)
	defaultStaticDirs = append(defaultStaticDirs, path.Join(parent, "static"))
	defaultStaticDirs = append(defaultStaticDirs, path.Join(wd, "static"))
	return
}

// Process invokes the main server's routing system.
func Process(c http.ResponseWriter, req *http.Request) {
	mainServer.Process(c, req)
}

// Run starts the web application and serves HTTP requests for the main server.
func Run(addr string) {
	mainServer.Run(addr)
}

// RunTLS starts the web application and serves HTTPS requests for the main server.
func RunTLS(addr string, config *tls.Config) {
	mainServer.RunTLS(addr, config)
}

// RunScgi starts the web application and serves SCGI requests for the main server.
func RunScgi(addr string) {
	mainServer.RunScgi(addr)
}

// RunFcgi starts the web application and serves FastCGI requests for the main server.
func RunFcgi(addr string) {
	mainServer.RunFcgi(addr)
}

// Close stops the main server.
func Close() {
	mainServer.Close()
}

// Get adds a handler for the 'GET' http method in the main server.
func Get(route string, handler interface{}) {
	mainServer.Get(route, handler)
}

// Post adds a handler for the 'POST' http method in the main server.
func Post(route string, handler interface{}) {
	mainServer.addRoute(route, "POST", handler)
}

// Put adds a handler for the 'PUT' http method in the main server.
func Put(route string, handler interface{}) {
	mainServer.addRoute(route, "PUT", handler)
}

// Delete adds a handler for the 'DELETE' http method in the main server.
func Delete(route string, handler interface{}) {
	mainServer.addRoute(route, "DELETE", handler)
}

// Match adds a handler for an arbitrary http method in the main server.
func Match(method string, route string, handler interface{}) {
	mainServer.addRoute(route, method, handler)
}

// Add a custom http.Handler. Will have no effect when running as FCGI or SCGI.
func Handle(route string, method string, httpHandler http.Handler) {
	mainServer.Handle(route, method, httpHandler)
}

//Adds a handler for websockets. Only for webserver mode. Will have no effect when running as FCGI or SCGI.
func Websocket(route string, httpHandler websocket.Handler) {
	mainServer.Websocket(route, httpHandler)
}

// SetLogger sets the logger for the main server.
func SetLogger(logger *log.Logger) {
	mainServer.Logger = logger
}

// Config is the configuration of the main server.
var Config = &ServerConfig{
	RecoverPanic: true,
	ColorOutput:  true,
}

var mainServer = NewServer()
