package benchmark2

import (
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/Cai-ki/caia/pkg/cprotocol"
)

func BenchmarkConcurrentConnections(b *testing.B) {
	var (
		targetAddr  = "localhost:9000" // 被测服务地址
		concurrency = 1000             // 并发连接数
		message     = "hello, world!"  // 测试消息内容
	)

	// // 预编码测试消息（所有连接共用）
	// coder := cprotocol.NewCodec(
	// 	&cprotocol.TLVHandler{},
	// 	func() interface{} { return &cprotocol.TLVPacket{} },
	// )

	// msg, err := coder.Encode(&cprotocol.TLVPacket{
	// 	Type:   100,
	// 	Length: uint32(len(message)),
	// 	Value:  []byte(message),
	// })
	// if err != nil {
	// 	b.Fatal("编码失败:", err)
	// }

	var wg sync.WaitGroup
	wg.Add(concurrency)

	b.ResetTimer() // 开始计时

	// 每个goroutine代表一个独立连接
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			coder := cprotocol.NewCodec(
				&cprotocol.TLVHandler{},
				func() interface{} { return &cprotocol.TLVPacket{} },
			)

			msg, err := coder.Encode(&cprotocol.TLVPacket{
				Type:   100,
				Length: uint32(len(message)),
				Value:  []byte(message),
			})
			if err != nil {
				b.Fatal("编码失败:", err)
			}
			// 建立独立连接
			conn, err := net.DialTimeout("tcp", targetAddr, 2*time.Second)
			if err != nil {
				b.Fatal("连接失败:", err)
				return
			}
			defer conn.Close()

			buf := make([]byte, 1024*1024) // 每个连接独立缓冲区

			// 每个连接执行 N/concurrency 次请求
			for i := 0; i < b.N/concurrency; i++ {
				// 发送请求
				if _, err := conn.Write(msg); err != nil {
					b.Fatal("发送失败:", err)
					return
				}

				err = errors.New("start")
				for err != nil {
					// 接收响应
					n, err := conn.Read(buf)
					if err != nil {
						b.Fatal("接收失败:", err)
						return
					}
					// 解码验证
					_, sz, err := coder.Decode(buf[:n])
					if err != nil {
						continue
					}

					if sz != len(msg) {
						b.Fatal("验证失败:", err, sz, len(msg))
						return
					}
					//fmt.Println(string(ifa.(*cprotocol.TLVPacket).Value))
					break
				}
			}
		}()
	}

	wg.Wait()
}
