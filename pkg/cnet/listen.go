package cnet

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/Cai-ki/caia/internal/cactor"
	"github.com/Cai-ki/caia/internal/clog"
	"github.com/Cai-ki/caia/internal/ctypes"
	"github.com/Cai-ki/caia/pkg/cruntime"
)

func ListenTCPHandle(ctx context.Context, msg ctypes.Message) {
	netConfig := cruntime.Configs[KeyConfig].(*Config)
	netActor := ctx.Value(KeyManager).(ctypes.Actor)
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", netConfig.Ip, netConfig.Port))
	if err != nil {
		clog.Error(fmt.Sprint("net: resolve tcp address error:", err))
		return
	}
	listenner, err := net.ListenTCP("tcp", addr)
	if err != nil {
		clog.Error(fmt.Sprintf("net: listen: %s, error: %s", "tcp", err))
		return
	}

	clog.Info(fmt.Sprintf("net: listen at %s:%d", netConfig.Ip, netConfig.Port))

	var cid uint32 = 0
	for {
		conn, err := listenner.AcceptTCP()
		if err != nil {
			clog.Error(fmt.Sprint("net: accept error", err))
			continue
		}

		connActor, err := netActor.CreateChild(strconv.Itoa(int(cid)), 10, ConnectHandle, cactor.WithValue("connect", conn), cactor.WithValue("cid", cid))
		connActor.Start()
		connActor.SendMessage(cruntime.MsgStart)
		cid++
	}
}
