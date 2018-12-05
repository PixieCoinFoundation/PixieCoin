package shutdown

import (
	"appcfg"
	"constants"
	"fmt"
	"io/ioutil"
	. "logger"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var shutdown bool
var globalServiceWaitGroup sync.WaitGroup

//request num being processing
var processingNum int
var numLock sync.Mutex

func AddOneRequest(name string) {
	globalServiceWaitGroup.Add(1)

	numLock.Lock()
	processingNum++
	numLock.Unlock()
}

func DoneOneRequest(name string) {
	globalServiceWaitGroup.Done()

	numLock.Lock()
	processingNum--
	numLock.Unlock()
}

func SimpleAddOneRequest() {
	globalServiceWaitGroup.Add(1)

	numLock.Lock()
	processingNum++
	numLock.Unlock()
}

func SimpleDoneOneRequest() {
	globalServiceWaitGroup.Done()

	numLock.Lock()
	processingNum--
	numLock.Unlock()
}

func GetProcessingRequestsNum() int {
	return processingNum
}

func IsShutDown() bool {
	return shutdown
}

type shutdown_task struct {
	tag  string
	task func()
}

// var task_map map[string]*shutdown_task

func init() {
	if appcfg.GetServerType() == "" || appcfg.GetServerType() == constants.SERVER_TYPE_GL || appcfg.GetServerType() == constants.SERVER_TYPE_GM || appcfg.GetServerType() == constants.SERVER_TYPE_GV {
		go startListenSignal()
	}
}

func startListenSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		switch sig := <-c; sig {
		case os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("server closing..")
			buf := make([]byte, 1<<20)
			runtime.Stack(buf, true)
			ioutil.WriteFile(fmt.Sprintf("/tmp/stack_%s.txt", time.Now().Format("2006-01-02-15-04-05")), buf, os.ModeAppend)
			shutdown = true
			doTasks()
			os.Exit(0)
		default:
			fmt.Println("Other")
		}
	}
}

// func Add(tag string, task func()) {
// 	if _, ok := task_map[tag]; ok {
// 		fmt.Println("add shutdown task error, task tag [", tag, "] exist.")
// 		os.Exit(1)
// 	}

// 	task_map[tag] = &shutdown_task{
// 		tag:  tag,
// 		task: task,
// 	}

// 	return
// }

func doTasks() {
	fmt.Println("shuting down...")
	fmt.Println("wait for processing requests to be processed..")
	globalServiceWaitGroup.Wait()
	fmt.Println("all requests processed.")
	// var wg sync.WaitGroup
	// defer wg.Wait()
	// for _, v := range task_map {
	// 	wg.Add(1)
	// 	go func(t *shutdown_task) {
	// 		fmt.Println("do task ", t.tag)
	// 		t.task()
	// 		fmt.Println("do task ", t.tag, " over")
	// 		wg.Done()
	// 	}(v)
	// }
	End()
	return
}
