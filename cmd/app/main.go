package main

import (
	"fmt"
	"net"
	"time"

	"github.com/Cai-ki/caia/internal/clog"
	"github.com/Cai-ki/caia/pkg/cnet"
	_ "github.com/Cai-ki/caia/pkg/cnet"
	"github.com/Cai-ki/caia/pkg/cprotocol"
	"github.com/Cai-ki/caia/pkg/cruntime"
)

func main() {
	cruntime.Start()

	<-time.After(time.Second)

	netConfig := cruntime.Configs["NetConfig"].(*cnet.Config)
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", netConfig.Ip, netConfig.Port))
	if err != nil {
		fmt.Println("main: resolve tcp address err: ", err)
		return
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("main: dial tcp err: ", err)
		return
	}

	t := time.After(5 * time.Second)

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
		clog.Error(err)
	}

	fmt.Println(msg)

	for {
		select {
		case <-t:
			cruntime.Stop()
			return
		case <-ticker.C:
			conn.SetDeadline(time.Now().Add(1 * time.Second))
			conn.Write(msg)
			data := make([]byte, 1024)
			n, err := conn.Read(data)
			if err != nil {
				fmt.Println("main: read err ", err)
				continue
			}

			iface, _, err := coder.Decode(data[:n])
			pkg := iface.(*cprotocol.TLVPacket)
			if err == nil {
				fmt.Println(pkg)
			}
			//fmt.Printf("main: receive %d byte, data: %s\n", n, string(data[:n]))
		}
	}
}
