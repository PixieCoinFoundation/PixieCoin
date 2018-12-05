package types

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
)

import ()

type Buffer struct {
	ctrl    chan bool    // receive exit signal
	ob      chan bool    // notify observer exit
	pending chan []byte  // pending packet
	max     int          // max queue size
	conn    *net.TCPConn // connection
}

const (
	DEFAULT_QUEUE_SIZE = 15
)

// send packet
func (buf *Buffer) Send(packet Packet) (err error) {
	defer func() {
		if x := recover(); x != nil {
			log.Println("buffer.Send failed", x)
		}
	}()

	if !checkPacket(packet) {
		return errors.New(fmt.Sprintf("Send buffer error, Remote: %v", buf.conn.RemoteAddr()))
	}

	data := make([]byte, packet.Header.HeadLen)
	binary.LittleEndian.PutUint16(data[0:], packet.Header.Tag)
	binary.LittleEndian.PutUint16(data[2:], packet.Header.HeadLen)
	binary.LittleEndian.PutUint16(data[4:], packet.Header.DataLen)
	binary.LittleEndian.PutUint16(data[6:], packet.Header.MsgClass)
	binary.LittleEndian.PutUint16(data[8:], packet.Header.MsgType)

	data = append(data, packet.Data...)

	buf.pending <- data

	return
}

// packet sender goroutine
func (buf *Buffer) Start() {
	defer func() {
		if x := recover(); x != nil {
			log.Println("caught panic in buffer goroutine", x)
		}
	}()

	for {
		select {
		case data := <-buf.pending:
			buf.raw_send(data)
		case <-buf.ctrl:
			// close(buf.pending)
			log.Println("close buffer")
			for data := range buf.pending {
				buf.raw_send(data)
			}
			buf.conn.Close()
			buf.ob <- true
			return
		}
	}
}

// send packet
func (buf *Buffer) raw_send(data []byte) {
	_, err := buf.conn.Write(data)
	if err != nil {
		log.Println("Error send reply: ", err)
		buf.Close()

		return
	}
}

func (buf *Buffer) SetConnection(conn *net.TCPConn) {
	buf.conn = conn
}

func (buf *Buffer) Close() {
	buf.ctrl <- true
}

// create a new write buffer
func NewBuffer(conn *net.TCPConn, buffOb chan bool, maxQueueSize int) *Buffer {
	size := maxQueueSize

	if size == -1 {
		size = DEFAULT_QUEUE_SIZE
	}

	buf := new(Buffer)
	buf.conn = conn
	buf.pending = make(chan []byte, size)
	buf.ctrl = make(chan bool)
	buf.ob = buffOb
	buf.max = size

	return buf
}

func checkPacket(packet Packet) bool {
	if packet.Header.Tag != PACKET_DEFAULT_TAG {
		return false
	}
	if packet.Header.HeadLen != 10 {
		return false
	}
	if packet.Header.DataLen != (uint16)(len(packet.Data)) {
		return false
	}

	return true
}
