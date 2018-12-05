package lz4

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	//    "fmt"
)

func TestDecodeLZ4Stream(t *testing.T) {
	in, _ := os.Open("cm.log.lz4")
	defer in.Close()
	out, _ := os.Create("testout")
	defer out.Close()
	DecodeLZ4Stream(in, out)
}

func BenchmarkDecodeLZ4Stream(b *testing.B) {
	b.SetBytes(211938580)
	in, err := ioutil.ReadFile("combined.lz4")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//        fmt.Println(i)
		DecodeLZ4Stream(bytes.NewReader(in), ioutil.Discard)
	}
}
