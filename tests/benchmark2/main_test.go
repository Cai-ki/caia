package benchmark2

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Cai-ki/caia/pkg/cprotocol"
)

func BenchmarkConcurrentConnections(b *testing.B) {
	var (
		targetAddr  = "localhost:9000"         // 被测服务地址
		concurrency = 100                      // 并发询问数
		message     = make([]byte, 1024*256-8) //"hello, world!"  // 测试消息内容
	)
	//TODO 超过1024字节数据会有问题

	var wg sync.WaitGroup
	wg.Add(b.N)

	c := atomic.Int32{}

	b.ResetTimer() // 开始计时

	fmt.Println(b.N)
	// 每个goroutine代表一个独立连接
	for i := 0; i < b.N; i++ {
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
				b.Error("编码失败:", err)
			}
			// 建立独立连接
			conn, err := net.DialTimeout("tcp", targetAddr, 2*time.Second)
			if err != nil {
				b.Error("连接失败:", err)
				return
			}
			defer conn.Close()

			buf := make([]byte, 1024*1024) // 每个连接独立缓冲区

			// 每个连接执行 N/concurrency 次请求
			for i := 0; i < concurrency/b.N; i++ {
				// 发送请求
				if _, err := conn.Write(msg); err != nil {
					b.Error("发送失败:", err)
					return
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
					c.Add(1)
					return true
				}

				for !readLoop() {
				}
				// for err != nil {
				// 	// 接收响应
				// 	n, err := conn.Read(buf)
				// 	if err != nil {
				// 		b.Fatal("接收失败:", err)
				// 		return
				// 	}
				// 	// 解码验证
				// 	ifs, _, err := coder.Decode(buf[:n])
				// 	if err != nil {
				// 		continue
				// 	}

				// 	if string(ifs.(*cprotocol.TLVPacket).Value) != "hello, world!" {
				// 		b.Fatal("验证失败:", err, "get:", string(ifs.(*cprotocol.TLVPacket).Value), "want:", "hello, world!")
				// 		return
				// 	}
				// 	//fmt.Println(string(ifa.(*cprotocol.TLVPacket).Value))
				// 	break
				// }
			}
		}()
	}

	wg.Wait()

	if c.Load() != int32(concurrency) {
		b.Error("fail want:", concurrency, "but:", c.Load())
	} else {
		fmt.Println("ok want:", concurrency, "and:", c.Load())
	}

}
