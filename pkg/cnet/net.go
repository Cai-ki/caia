package cnet

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/Cai-ki/caia/internal/ctypes"
	"github.com/Cai-ki/caia/pkg/cruntime"
)

func init() {
	netActor, err := cruntime.RootActor.CreateChild("net", 10, nil)
	if err != nil {
		panic("net: init fail")
	} else {
		netActor.SetHandle(func(context.Context, ctypes.Message) {
			addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "", 60000))
			if err != nil {
				fmt.Println("net: resolve tcp address err: ", err)
				return
			}
			listenner, err := net.ListenTCP("tcp", addr)
			if err != nil {
				fmt.Printf("net: listen: %s, err: %s\n", "tcp", err)
				return
			}
			var cid uint32 = 0
			for {
				conn, err := listenner.AcceptTCP()
				if err != nil {
					fmt.Println("net: accept err ", err)
					continue
				}

				connRoot, err := netActor.CreateChild(strconv.Itoa(int(cid)), 10, func(ctx context.Context, msg ctypes.Message) {
					for {
						data := make([]byte, 1024)
						n, err := conn.Read(data)
						if err != nil {
							fmt.Println("net: read err ", err)
							continue
						}
						conn.Write([]byte(fmt.Sprintf("id: %d send %d byte, data: %s\n", cid, n, string(data[:n]))))
					}
				})
				connRoot.Start()
				connRoot.SendMessage(cruntime.MsgStart)
				cid++
			}
		})
	}
}
