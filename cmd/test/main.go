package main

import (
	"fmt"
	"net"
	"time"

	_ "github.com/Cai-ki/caia/pkg/cnet"
	"github.com/Cai-ki/caia/pkg/cruntime"
)

func main() {
	cruntime.Start()

	<-time.After(time.Second)

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "", 60000))
	if err != nil {
		fmt.Println("main: resolve tcp address err: ", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("main: dial tcp err: ", err)
		return
	}

	for {
		conn.Write([]byte("hello, world!"))
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println("main: read err ", err)
			continue
		}
		fmt.Printf("main: receive %d byte, data: %s\n", n, string(data[:n]))
		<-time.After(time.Second)
	}
}
