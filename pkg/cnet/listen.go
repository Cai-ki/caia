package cnet

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Cai-ki/caia/internal/cactor"
	"github.com/Cai-ki/caia/internal/clog"
	"github.com/Cai-ki/caia/internal/ctypes"
	"github.com/Cai-ki/caia/pkg/cruntime"
)

func ListenTCPHandle(ctx context.Context, msg ctypes.Message) {
	// config := cruntime.Configs[KeyConfig].(*Config)
	manager := ctx.Value(KeyManager).(ctypes.Actor)
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Ip, config.Port))
	if err != nil {
		clog.Error("net: resolve tcp address error:", err)
		return
	}
	listenner, err := net.ListenTCP("tcp", addr)
	if err != nil {
		clog.Errorf("net: listen: %s, error: %s", "tcp", err)
		return
	}
	defer listenner.Close()

	clog.Infof("net: listening on %s:%d", config.Ip, config.Port)

	addTime := time.Duration(config.ListenDeadlineMs) * time.Millisecond
	var cid uint32 = 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			listenner.SetDeadline(time.Now().Add(addTime))
			conn, err := listenner.AcceptTCP()
			if err != nil {
				if os.IsTimeout(err) {
					continue
				}
				clog.Error("net: accept error:", err)
				continue
			}

			connActor, err := manager.CreateChild(strconv.Itoa(int(cid)), 10, ConnectHandle, cactor.WithValue("connect", conn), cactor.WithValue("cid", cid))
			if err != nil {
				clog.Errorf("net: failed to create connection with %s: %v", conn.RemoteAddr().String(), err)
			}

			clog.Infof("net: established connection with %s", conn.RemoteAddr().String())

			connActor.Start()
			connActor.SendMessage(cruntime.MsgStart)
			cid++
		}
	}
}
