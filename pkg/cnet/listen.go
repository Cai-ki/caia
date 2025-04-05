package cnet

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/Cai-ki/caia/internal/cactor"
	"github.com/Cai-ki/caia/internal/clog"
	"github.com/Cai-ki/caia/internal/ctypes"
	"github.com/Cai-ki/caia/pkg/cruntime"
	"github.com/panjf2000/ants/v2"
)

var Pool, err = ants.NewPool(100000)

func ListenTCPHandle(actor ctypes.Actor, msg ctypes.Message) {
	// config := cruntime.Configs[KeyConfig].(*Config)
	ctx := actor.GetContext()
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

	ticker := time.NewTicker(2 * time.Second)

	addTime := time.Duration(config.ListenDeadlineMs) * time.Millisecond
	var cid uint32 = 0
	for {
		select {
		case <-ticker.C:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			clog.Info(fmt.Sprintf("Alloc: %v MB, TotalAlloc: %v MB, Sys: %v MB, NumGC: %v",
				memStats.Alloc/1024/1024,
				memStats.TotalAlloc/1024/1024,
				memStats.Sys/1024/1024,
				memStats.NumGC))
			clog.Info("NumGoroutine: ", runtime.NumGoroutine(), "NumCPU: ", runtime.NumCPU(), "ants.Cap: ", Pool.Cap(), "ants.Running: ", Pool.Running(), "son: ", len(actor.GetChildren()))
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

			connActor, err := actor.CreateChild(strconv.Itoa(int(cid)), 10, ConnectHandleFactory(conn), cactor.WithValue("connect", conn), cactor.WithValue("cid", cid))
			if err != nil {
				clog.Errorf("net: failed to create connection with %s: %v", conn.RemoteAddr().String(), err)
			}

			clog.Debugf("net: established connection with %s", conn.RemoteAddr().String())

			connActor.Start()
			connActor.SendMessage(cruntime.MsgStart)
			cid++
		}
	}
}
