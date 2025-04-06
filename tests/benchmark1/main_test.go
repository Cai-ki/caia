package benchmark2

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/Cai-ki/caia/pkg/cprotocol"
)

func BenchmarkConcurrentConnections(b *testing.B) {
	var (
		targetAddr  = "localhost:9000"    // 被测服务地址
		concurrency = 10000               // 并发连接数
		message     = make([]byte, 256-8) //"hello, world!"  // 测试消息内容
	)
	//TODO 超过1024字节数据会有问题

	var wg sync.WaitGroup
	wg.Add(concurrency)

	b.ResetTimer() // 开始计时

	// 每个goroutine代表一个独立连接
	for ii := 0; ii < concurrency; ii++ {
		go func() {
			defer wg.Done()
			coder := cprotocol.NewCodec(
				&cprotocol.TLVHandler{},
				func() interface{} { return &cprotocol.TLVPacket{} },
			)

			msg, err := coder.Encode(&cprotocol.TLVPacket{
				Type:   uint32(0),
				Length: uint32(len(message)),
				Value:  []byte(message),
			})
			if err != nil {
				b.Error("编码失败:", err)
			}
			// 建立独立连接
			conn, err := net.DialTimeout("tcp", targetAddr, 10*time.Second)
			if err != nil {
				b.Error("连接失败:", err)
				return
			}
			defer conn.Close()

			buf := make([]byte, 1024*1024) // 每个连接独立缓冲区

			// 每个连接执行 N/concurrency 次请求
			for i := 0; i < b.N/concurrency; i++ {
				writeLoop := func() bool {
					n, err := conn.Write(msg)
					if err != nil {
						b.Error("发送失败:", err)
						return true
					}
					if n < len(msg) {
						msg = msg[:n]
						return false
					}
					return true
				}

				for !writeLoop() {
				}

				readLoop := func() bool {
					n, err := conn.Read(buf)
					if err != nil {
						b.Error("接收失败:", err)
						return true
					}
					// 解码验证
					ifs, _, err := coder.Decode(buf[:n])
					if err != nil {
						return false
					}

					if string(ifs.(*cprotocol.TLVPacket).Value) != "hello, world!" {
						b.Error("验证失败:", err, "get:", string(ifs.(*cprotocol.TLVPacket).Value), "want:", "hello, world!")
						return true
					}
					return true
				}

				for !readLoop() {
				}
			}
		}()
	}

	wg.Wait()

}
