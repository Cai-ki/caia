package cnet

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/Cai-ki/caia/internal/cactor"
	"github.com/Cai-ki/caia/internal/ctypes"
	"github.com/Cai-ki/caia/pkg/cruntime"
)

func ListenTCPHandle(ctx context.Context, msg ctypes.Message) {
	netConfig := cruntime.Configs[KeyConfig].(*Config)
	netActor := ctx.Value(KeyManager).(ctypes.Actor)
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", netConfig.Ip, netConfig.Port))
	if err != nil {
		fmt.Println("net: resolve tcp address err: ", err)
		return
	}
	listenner, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Printf("net: listen: %s, err: %s\n", "tcp", err)
		return
	}

	fmt.Printf("net: listen at %s:%d\n", netConfig.Ip, netConfig.Port)

	var cid uint32 = 0
	for {
		conn, err := listenner.AcceptTCP()
		if err != nil {
			fmt.Println("net: accept err ", err)
			continue
		}

		connActor, err := netActor.CreateChild(strconv.Itoa(int(cid)), 10, ConnectHandle, cactor.WithValue("connect", conn), cactor.WithValue("cid", cid))
		connActor.Start()
		connActor.SendMessage(cruntime.MsgStart)
		cid++
	}
}
