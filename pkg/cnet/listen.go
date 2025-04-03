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
		clog.Error("net: resolve tcp address error:", err)
		return
	}
	listenner, err := net.ListenTCP("tcp", addr)
	if err != nil {
		clog.Errorf("net: listen: %s, error: %s", "tcp", err)
		return
	}

	clog.Infof("net: listening on %s:%d", netConfig.Ip, netConfig.Port)

	var cid uint32 = 0
	for {
		conn, err := listenner.AcceptTCP()
		if err != nil {
			clog.Error("net: accept error:", err)
			continue
		}

		connActor, err := netActor.CreateChild(strconv.Itoa(int(cid)), 10, ConnectHandle, cactor.WithValue("connect", conn), cactor.WithValue("cid", cid))
		if err != nil {
			clog.Errorf("net: failed to create connection with %s: %v", conn.RemoteAddr().String(), err)
		}

		clog.Infof("net: established connection with %s", conn.RemoteAddr().String())

		connActor.Start()
		connActor.SendMessage(cruntime.MsgStart)
		cid++
	}
}
