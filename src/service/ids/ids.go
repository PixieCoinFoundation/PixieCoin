package ids

// import (
// 	"constants"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"sync"
// )

// import (
// 	"appcfg"
// 	"rao"
// )

// var (
// 	serverID         string
// 	serverStartCount int64
// 	prepend          string

// 	helpID     int64
// 	helpIDLock sync.Mutex

// 	orderID     int64
// 	orderIDLock sync.Mutex
// )

func init() {
	// if !appcfg.GetBool("recover", false) && appcfg.GetServerType() != constants.SERVER_TYPE_GV {
	// 	ResetIDs()
	// }
}

// func ResetIDs() {
// 	serverID = appcfg.GetString("server_id", "")
// 	if serverID == "" {
// 		fmt.Println("server id not set！Please check scripts/config.json")
// 		os.Exit(1)
// 		return
// 	}
// 	var err error
// 	serverStartCount, err = rao.GetBootCount()
// 	if err != nil {
// 		fmt.Println("initialize serverStartCount failed！error:", err.Error())
// 		os.Exit(1)
// 	}

// 	prepend = serverID + "_" + strconv.FormatInt(serverStartCount, 10) + "_"

// 	// helpIDLock.Lock()
// 	// helpID, err = rao.GetHelpID()
// 	// helpIDLock.Unlock()

// 	if err != nil {
// 		fmt.Println("initialize help id failed！")
// 		os.Exit(1)
// 	}

// 	// orderIDLock.Lock()
// 	// orderID, err = rao.GetOrderID()
// 	// orderIDLock.Unlock()

// 	if err != nil {
// 		fmt.Println("initialize order id failed！")
// 		os.Exit(1)
// 	}

// 	fmt.Println("server_start_count:", serverStartCount)
// }

// func GetHelpID() string {
// 	helpIDLock.Lock()
// 	helpID++
// 	tmp := helpID
// 	helpIDLock.Unlock()

// 	ret := prepend + strconv.FormatInt(tmp, 10)
// 	return ret
// }

// func GetOrderID() string {
// 	orderIDLock.Lock()
// 	orderID++
// 	tmp := orderID
// 	orderIDLock.Unlock()

// 	ret := prepend + strconv.FormatInt(tmp, 10)
// 	return ret
// }

// func GetCosItemID() (string, error) {
// 	if itemID, err := rao.GetCosItemID(); err != nil {
// 		return "", err
// 	} else {
// 		return strconv.FormatInt(itemID, 10), nil
// 	}
// }
