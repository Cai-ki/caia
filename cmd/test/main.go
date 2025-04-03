package main

import (
	"fmt"
	"net"
	"time"

	"github.com/Cai-ki/caia/pkg/cnet"
	_ "github.com/Cai-ki/caia/pkg/cnet"
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

	t := time.After(4 * time.Second)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-t:
			cruntime.Stop()
			return
		case <-ticker.C:
			conn.SetDeadline(time.Now().Add(1 * time.Second))
			conn.Write([]byte("hello, world!"))
			data := make([]byte, 1024)
			_, err := conn.Read(data)
			if err != nil {
				fmt.Println("main: read err ", err)
				continue
			}
			//fmt.Printf("main: receive %d byte, data: %s\n", n, string(data[:n]))
		}
	}
}
