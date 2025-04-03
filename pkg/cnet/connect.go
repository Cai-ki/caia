package cnet

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Cai-ki/caia/internal/clog"
	"github.com/Cai-ki/caia/internal/ctypes"
)

func ConnectHandle(ctx context.Context, msg ctypes.Message) {
	conn := ctx.Value(KeyConnect).(*net.TCPConn)
	defer conn.Close()
	cid := ctx.Value(KeyCid).(uint32)
	for {
		select {
		case <-ctx.Done():
			conn.Close()
			return
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			data := make([]byte, 1024)
			n, err := conn.Read(data)
			if err != nil {
				if os.IsTimeout(err) {
					continue
				}
				clog.Error("net: read error:", err)
				return
			}
			conn.Write([]byte(fmt.Sprintf("id: %d send %d byte, data: %s\n", cid, n, string(data[:n]))))
			clog.Infof("id: %d send %d byte, data: %s", cid, n, string(data[:n]))
		}
	}
}
