// +build linux

package gops

import (
	// "appcfg"
	"bytes"
	// "constants"
	// "fmt"
	"log"
	// "os"
	// "sync"

	// "github.com/google/gops/internal"
	"github.com/google/gops/internal/objfile"
	ps "github.com/keybase/go-ps"

	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/gops/internal"
	"github.com/google/gops/signal"
)

var cmds = map[string](func(addr net.TCPAddr) error){
	"stack":      stackTrace,
	"gc":         gc,
	"memstats":   memStats,
	"version":    version,
	"pprof-heap": pprofHeap,
	"pprof-cpu":  pprofCPU,
	"stats":      stats,
	"trace":      trace,
}

const helpText = `Usage: gops is a tool to list and diagnose Go processes.


Commands:
    stack       Prints the stack trace.
    gc          Runs the garbage collector and blocks until successful.
    memstats    Prints the allocation and garbage collection stats.
    version     Prints the Go version used to build the program.
    stats       Prints the vital runtime stats.
    help        Prints this help text.

Profiling commands:
    trace       Runs the runtime tracer for 5 secs and launches "go tool trace".
    pprof-heap  Reads the heap profile and launches "go tool pprof".
    pprof-cpu   Reads the CPU profile and launches "go tool pprof".


All commands require the agent running on the Go process.
Symbol "*" indicates the process runs the agent.`

// TODO(jbd): add link that explains the use of agent.

// func main() {
// 	if len(os.Args) < 2 {
// 		processes()
// 		return
// 	}

// 	cmd := os.Args[1]
// 	if cmd == "help" {
// 		usage("")
// 	}
// 	if len(os.Args) < 3 {
// 		usage("missing PID or address")
// 	}
// 	fn, ok := cmds[cmd]
// 	if !ok {
// 		usage("unknown subcommand")
// 	}
// 	addr, err := targetToAddr(os.Args[2])
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Couldn't resolve addr or pid %v to TCPAddress: %v\n", os.Args[2], err)
// 		os.Exit(1)
// 	}
// 	if err := fn(*addr); err != nil {
// 		fmt.Fprintf(os.Stderr, "%v\n", err)
// 		os.Exit(1)
// 	}
// }

func Processes(pid string) {
	st := time.Now().UnixNano()
	defer fmt.Println("use time:", float64(time.Now().UnixNano()-st)/float64(time.Millisecond), "ms")

	pss, err := ps.Processes()
	if err != nil {
		log.Fatal(err)
	}

	// var wg sync.WaitGroup
	// wg.Add(len(pss))

	for _, pr := range pss {
		pr := pr
		// go func() {
		// defer wg.Done()

		printIfGo(pr, pid)
		// }()
	}
	// wg.Wait()
}

// printIfGo looks up the runtime.buildVersion symbol
// in the process' binary and determines if the process
// if a Go process or not. If the process is a Go process,
// it reports PID, binary name and full path of the binary.
func printIfGo(pr ps.Process, pid string) {
	if pr.Pid() == 0 {
		// ignore system process
		return
	}
	path, err := pr.Path()
	if err != nil {
		return
	}
	obj, err := objfile.Open(path)
	if err != nil {
		return
	}
	defer obj.Close()

	symbols, err := obj.Symbols()
	if err != nil {
		return
	}

	var ok bool
	for _, s := range symbols {
		if s.Name == "runtime.buildVersion" {
			ok = true
		}
	}

	var agent bool
	pidfile, err := internal.PIDFile(pr.Pid())
	if err == nil {
		_, err := os.Stat(pidfile)
		agent = err == nil
	}

	if ok {
		buf := bytes.NewBuffer(nil)
		fmt.Fprintf(buf, "%d", pr.Pid())
		if agent {
			fmt.Fprint(buf, "*")
		}
		fmt.Fprintf(buf, "\t%v\t(%v)\n", pr.Executable(), path)
		buf.WriteTo(os.Stdout)

		if fmt.Sprintf("%d", pr.Pid()) == pid {
			addr, err := targetToAddr(pid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Couldn't resolve addr or pid %v to TCPAddress: %v\n", os.Args[2], err)
				os.Exit(1)
			}

			// fmt.Println("---------------------------stack trace-------------------------")
			// stackTrace(*addr)
			fmt.Println("---------------------------gc-------------------------")
			gc(*addr)
			fmt.Println("---------------------------mem stats-------------------------")
			memStats(*addr)
			fmt.Println("---------------------------version-------------------------")
			version(*addr)
			// fmt.Println("---------------------------pprof heap-------------------------")
			// pprofHeap(*addr)
			// fmt.Println("---------------------------pprof cpu-------------------------")
			// pprofCPU(*addr)
			fmt.Println("---------------------------stats-------------------------")
			stats(*addr)
			// fmt.Println("---------------------------trace-------------------------")
			// trace(*addr)

		}
	}
}

func usage(msg string) {
	if msg != "" {
		fmt.Printf("gops: %v\n", msg)
	}
	fmt.Fprintf(os.Stderr, "%v\n", helpText)
	os.Exit(1)
}

func stackTrace(addr net.TCPAddr) error {
	return cmdWithPrint(addr, signal.StackTrace)
}

func gc(addr net.TCPAddr) error {
	_, err := cmd(addr, signal.GC)
	return err
}

func memStats(addr net.TCPAddr) error {
	return cmdWithPrint(addr, signal.MemStats)
}

func version(addr net.TCPAddr) error {
	return cmdWithPrint(addr, signal.Version)
}

func pprofHeap(addr net.TCPAddr) error {
	return pprof(addr, signal.HeapProfile)
}

func pprofCPU(addr net.TCPAddr) error {
	fmt.Println("Profiling CPU now, will take 30 secs...")
	return pprof(addr, signal.CPUProfile)
}

func trace(addr net.TCPAddr) error {
	fmt.Println("Tracing now, will take 5 secs...")
	out, err := cmd(addr, signal.Trace)
	if err != nil {
		return err
	}
	if len(out) == 0 {
		return errors.New("nothing has traced")
	}
	tmpfile, err := ioutil.TempFile("", "trace")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())
	if err := ioutil.WriteFile(tmpfile.Name(), out, 0); err != nil {
		return err
	}
	cmd := exec.Command("go", "tool", "trace", tmpfile.Name())
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pprof(addr net.TCPAddr, p byte) error {

	tmpDumpFile, err := ioutil.TempFile("", "profile")
	if err != nil {
		return err
	}
	{
		out, err := cmd(addr, p)
		if err != nil {
			return err
		}
		if len(out) == 0 {
			return errors.New("failed to read the profile")
		}
		defer os.Remove(tmpDumpFile.Name())
		if err := ioutil.WriteFile(tmpDumpFile.Name(), out, 0); err != nil {
			return err
		}
	}
	// Download running binary
	tmpBinFile, err := ioutil.TempFile("", "binary")
	if err != nil {
		return err
	}
	{
		out, err := cmd(addr, signal.BinaryDump)
		if err != nil {
			return fmt.Errorf("failed to read the binary: %v", err)
		}
		if len(out) == 0 {
			return errors.New("failed to read the binary")
		}
		defer os.Remove(tmpBinFile.Name())
		if err := ioutil.WriteFile(tmpBinFile.Name(), out, 0); err != nil {
			return err
		}
	}
	cmd := exec.Command("go", "tool", "pprof", tmpBinFile.Name(), tmpDumpFile.Name())
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func stats(addr net.TCPAddr) error {
	return cmdWithPrint(addr, signal.Stats)
}

func cmdWithPrint(addr net.TCPAddr, c byte) error {
	out, err := cmd(addr, c)
	if err != nil {
		return err
	}
	fmt.Printf("%s", out)
	return nil
}

// targetToAddr tries to parse the target string, be it remote host:port
// or local process's PID.
func targetToAddr(target string) (*net.TCPAddr, error) {
	if strings.Index(target, ":") != -1 {
		// addr host:port passed
		var err error
		addr, err := net.ResolveTCPAddr("tcp", target)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse dst address: %v", err)
		}
		return addr, nil
	}
	// try to find port by pid then, connect to local
	pid, err := strconv.Atoi(target)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse PID: %v", err)
	}
	port, err := internal.GetPort(pid)
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:"+port)
	return addr, nil
}

func cmd(addr net.TCPAddr, c byte) ([]byte, error) {
	conn, err := cmdLazy(addr, c)
	if err != nil {
		return nil, fmt.Errorf("couldn't get port by PID: %v", err)
	}

	all, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func cmdLazy(addr net.TCPAddr, c byte) (io.Reader, error) {
	conn, err := net.DialTCP("tcp", nil, &addr)
	if err != nil {
		return nil, err
	}
	if _, err := conn.Write([]byte{c}); err != nil {
		return nil, err
	}
	return conn, nil
}
