package handle

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// import (
// 	"encrypt/simple"
// 	"gl_handle/gf"
// 	. "types"
// )

// type logicData struct {
// 	Address string `json:"address"`
// 	Method  string `json:"method"`
// 	Params  string `json:"params"`
// }

// type logicResp struct {
// 	Status  ErrorCode
// 	Content string
// }

// const (
// 	AES_KEY string = "aes_key_huangjundabendan_u3829us"
// )

// func Direct(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()

// 	data := r.PostFormValue("data")
// 	decryptedData := decrypt(data)

// 	var logic logicData
// 	err := json.Unmarshal([]byte(decryptedData), &logic)
// 	if err != nil {
// 		// 注意，有可能是一些特殊字符不支持json decode
// 		fmt.Println("json decode错误")
// 		fmt.Println(err.Error())
// 		resp := logicResp{
// 			Status:  ERR_WRONG_PARAMS,
// 			Content: "format error",
// 		}
// 		jsonByte, _ := json.Marshal(resp)
// 		fmt.Fprintf(w, string(jsonByte))
// 		return
// 	}

// 	var resp logicResp
// 	if handle, ok := gf.RequestHandles[logic.Method]; ok {
// 		result := make(chan string)
// 		handleErr := make(chan int)
// 		go handle(logic.Address, logic.Params, result, handleErr)

// 		for {
// 			select {
// 			case content := <-result:
// 				resp.Content = content
// 				resp.Status = NO_ERROR
// 			case <-handleErr:
// 				resp.Content = "func error"
// 				resp.Status = ERR_HANDLE_PARAMS_WRONG
// 			}
// 			break
// 		}

// 	} else {
// 		resp.Status = ERR_NO_HANDLE_FOUND
// 		resp.Content = ""
// 	}

// 	jsonByte, _ := json.Marshal(resp)
// 	fmt.Fprintf(w, string(jsonByte))
// }

// func encrypt(input string) string {
// 	result := simple.Encode(input)
// 	return result
// }

// func decrypt(input string) string {
// 	result := simple.Decode(input)
// 	return result
// }
