package lz4

import (
	"io"
	"os"
	"testing"

	"fmt"
)

func printBytes(name string, p []byte) {
	fmt.Println(name)
	for i := 0; i < len(p); i++ {
		if i%8 == 0 {
			fmt.Printf("|[%v]| ", i)
		}
		fmt.Printf("%02x ", p[i])
	}
	fmt.Printf("\n")
}
func deepEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			fmt.Printf("not equal at: %v, %02x %02x\n", i, a[i], b[i])
			return false
		}
	}
	return true
}
func TestEncodeLZ4SingleBlock(t *testing.T) {
	in, _ := os.Open("cm.log")
	defer in.Close()

	inBlock := make([]byte, 65536)
	compBlock := make([]byte, 65536)
	outBlock := make([]byte, 65536)
	in.Read(inBlock)

	_, ln, err := EncodeLZ4SingleBlock(inBlock, compBlock)
	if err != nil {
		t.Fatalf("enocde error in test")
	}
	_, l, err := DecodeLZ4SingleBlock(compBlock, outBlock, ln)
	if err != nil {
		t.Fatalf("deocde error in test, %v", err)
	}
	if l != 65536 {
		t.Fatalf("deocde len error in test, len=%v", l)
	}
	if deepEqual(inBlock, outBlock) != true {
		printBytes("inBlock", inBlock)
		printBytes("outBlock", outBlock)
		printBytes("compBlock", compBlock)
		t.Fatalf("input output not equal")
	}
}

func BenchmarkEncodeLZ4SingleBlock(b *testing.B) {
	in, _ := os.Open("combined.bin")
	defer in.Close()

	inBlock := make([]byte, 32768*6*1024)
	outBlock := make([]byte, 42768)
	_, err := in.Read(inBlock)
	if err != io.EOF {
		in.Seek(0, 0)
	}
	b.SetBytes(32768)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := (32768 * i) % (32768 * 6 * 1023)
		_, _, err = EncodeLZ4SingleBlock(inBlock[start:start+32768], outBlock)
		if err != nil {
			b.Fatalf("enocde error in test")
		}
	}
	/*
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			_, err := in.Read(inBlock)
			if err != io.EOF {
				in.Seek(0, 0)
			}
			b.StartTimer()

			_, _, err = EncodeLZ4SingleBlock(inBlock, outBlock)
			if err != nil {
				b.Fatalf("enocde error in test")
			}
		}
	*/
}

func BenchmarkDecodeLZ4SingleBlock(b *testing.B) {
	in, _ := os.Open("combined.bin")
	defer in.Close()

	inBlock := make([]byte, 32768*6*1024)
	outBlockSlice := make([]([]byte), 6*1024)
	outBlockSize := make([]uint32, 6*1024)
	_, err := in.Read(inBlock)
	if err != io.EOF {
		in.Seek(0, 0)
	}
	for i := 0; i < 6*1024; i++ {
		outBlockSlice[i] = make([]byte, 42768)
		start := (32768 * i) % (32768 * 6 * 1023)
		_, outBlockSize[i], err = EncodeLZ4SingleBlock(inBlock[start:start+32768], outBlockSlice[i])
		if err != nil {
			b.Fatalf("enocde error in test")
		}
	}
	b.SetBytes(32768)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err = DecodeLZ4SingleBlock(outBlockSlice[i%(6*1024)], inBlock, outBlockSize[i%(6*1024)])
		if err != nil {
			b.Fatalf("enocde error in test")
		}
	}
}
