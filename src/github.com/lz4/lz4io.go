package lz4

import (
	//    "errors"
	"encoding/binary"
	"fmt"
	"io"
	//    "github.com/vova616/xxhash"
)

const MAGICNUM = 0X184D2204

type streamDesc struct {
	FLG        byte
	BD         byte
	StreamSize []byte
	DictID     []byte
	HC         byte
}

func (s *streamDesc) parse(in io.Reader) (e error) {
	buf := make([]byte, 8)
	_, err := in.Read(buf[0:1])
	if err != nil {
		return err
	}
	s.FLG = buf[0]
	_, err = in.Read(buf[0:1])
	if err != nil {
		return err
	}
	s.BD = buf[0]
	_, err = in.Read(buf[0:1])
	if err != nil {
		return err
	}
	s.HC = buf[0]
	//h32 := xxhash.Checksum32Seed([]byte{s.FLG, s.BD}, 0)
	//fmt.Printf("Head checksum (xxhash) 0x%02X\n", h32)
	//fmt.Printf("FLG: %08b, BD %08b, HD 0x%02X\n", s.FLG, s.BD, s.HC)
	return nil
}
func DecodeLZ4Stream(in io.Reader, out io.Writer) (e error) {
	var magicNumber uint32
	err := binary.Read(in, binary.LittleEndian, &magicNumber)
	if err != nil {
		return err
	}

	if magicNumber != MAGICNUM {
		return fmt.Errorf("wrong magic number: 0x%x", magicNumber)
	}
	var s streamDesc
	s.parse(in)
	err = DecodeLZ4Blocks(in, out)
	if err != nil {
		return err
	}
	return nil
}
func DecodeLZ4Blocks(in io.Reader, out io.Writer) (e error) {
	ibuf := make([]byte, 4*1024*1024)
	obuf := make([]byte, 4*1024*1024)
	for {
		// read block size
		var blockSize uint32
		err := binary.Read(in, binary.LittleEndian, &blockSize)
		if err != nil {
			return err
		}

		//last block
		if blockSize == 0 {
			break
		}

		// data is uncompressed if highest bit is 1
		if (blockSize & 0x80000000) != 0 {
			io.CopyN(out, in, int64(blockSize&0x7FFFFFFF))
		} else {
			in.Read(ibuf[:blockSize])
			proccessed, outLen, err := DecodeLZ4SingleBlock(ibuf[:blockSize], obuf, blockSize)
			out.Write(obuf[0:outLen])
			if err != nil {
				fmt.Printf("proccessed %d bytes, output %d bytes, err %s", proccessed, outLen, err)
				return err
			}
		}
	}
	return nil
}
