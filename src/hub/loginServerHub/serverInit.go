package loginServerHub

import (
	"constants"
	"io"
	"net"
	"os"
)

import (
	"appcfg"
	. "logger"
	. "types"
)

func init() {
	if appcfg.GetServerType() == constants.SERVER_TYPE_GL {
		// 启动监听，等待其它server的连接
		port := appcfg.GetString("gs_port", ":29981")

		// 启动tcp监听
		tcpAddr, err := net.ResolveTCPAddr("tcp4", port)
		checkError(err)

		listener, err := net.ListenTCP("tcp", tcpAddr)
		checkError(err)

		Info("Waiting for game server at port", port)

		go func() {
			for {
				conn, err := listener.AcceptTCP()
				if err != nil {
					continue
				}

				Info("game server:", conn.RemoteAddr().String(), "connected...")
				go handleClient(conn)
			}
		}()
	}
}

func handleClient(conn *net.TCPConn) {
	defer func() {
		if x := recover(); x != nil {
			Info("caught panic in handleClient", x)
		}
	}()
	defer conn.Close()

	header := make([]byte, 10)
	ch := make(chan Packet, 65535)

	go startOneGFServer(conn, ch)

	for {
		n, err := io.ReadFull(conn, header)
		if n == 0 && err == io.EOF {
			Info("game server end:", conn.RemoteAddr().String())
			break
		} else if err != nil {
			Info("error receiving from game server:", err)
			break
		}

		var packet Packet
		packet.SetHead(header)

		data := make([]byte, packet.Header.DataLen)
		n, err = io.ReadFull(conn, data)
		if err != nil {
			Info("error receiving head from game server:", err)
			break
		}
		packet.SetRawData(data[0:])

		ch <- packet
	}

	close(ch)
}

func checkError(err error) {
	if err != nil {
		Info("Fatal error: %v", err)
		os.Exit(-1)
	}
}
