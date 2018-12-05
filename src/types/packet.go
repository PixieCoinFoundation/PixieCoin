package types

import (
	"encoding/binary"
	"log"
	"snappy"
)

import ()

const (
	PACKET_DEFAULT_TAG uint16 = (('O' & 0x00ff) | 'T'<<8)
)

type Head struct {
	Tag      uint16 // tag
	HeadLen  uint16 // 报文头长度
	DataLen  uint16 // 数据长度，不包括头的长度
	MsgClass uint16 // 消息的大类型
	MsgType  uint16 // 消息的小类型
}

type Packet struct {
	Header Head   // 报文头
	Data   []byte // 数据，json
}

func (self *Packet) SetHead(headbytes []byte) bool {
	if len(headbytes) != 10 {
		log.Println("packet header wrong size")
		return false
	}

	self.Header.Tag = binary.LittleEndian.Uint16(headbytes[0:])
	self.Header.HeadLen = binary.LittleEndian.Uint16(headbytes[2:])
	self.Header.DataLen = binary.LittleEndian.Uint16(headbytes[4:])
	self.Header.MsgClass = binary.LittleEndian.Uint16(headbytes[6:])
	self.Header.MsgType = binary.LittleEndian.Uint16(headbytes[8:])

	return true
}

func (self *Packet) SetData(data []byte) {
	encoded := snappy.Encode(nil, data)
	self.Header.DataLen = uint16(len(encoded))

	self.Data = encoded
}

func (self *Packet) SetRawData(data []byte) {
	decoded, _ := snappy.Decode(nil, data)
	self.Header.DataLen = uint16(len(decoded))

	self.Data = decoded
}
