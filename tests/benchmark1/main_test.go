package benchmark1

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Cai-ki/caia/pkg/cprotocol"
)

func BenchmarkGo(b *testing.B) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "", 9000))
	if err != nil {
		b.Fatal("main: resolve tcp address err: ", err)
		return
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		b.Fatal("main: dial tcp err: ", err)
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	coder := cprotocol.NewCodec(&cprotocol.TLVHandler{}, func() interface{} {
		return &cprotocol.TLVPacket{Value: make([]byte, 0, 128)}
	})

	msg, err := coder.Encode(&cprotocol.TLVPacket{
		Type:   uint32(100),
		Length: uint32(13),
		Value:  []byte("hello, world!"),
	})

	if err != nil {
		b.Error(err)
	}

	data := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		conn.Write(msg)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println("main: read err ", err)
			continue
		}

		_, sz, err := coder.Decode(data[:n])
		if err == nil {
			if sz != len(msg) {
				b.Failed()
			}
		}
	}
	b.StopTimer()
}
