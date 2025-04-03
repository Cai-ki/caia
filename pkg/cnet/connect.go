package cnet

import (
	"context"
	"fmt"
	"net"

	"github.com/Cai-ki/caia/internal/ctypes"
)

func ConnectHandle(ctx context.Context, msg ctypes.Message) {
	conn := ctx.Value(KeyConnect).(*net.TCPConn)
	cid := ctx.Value(KeyCid).(uint32)
	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println("net: read err ", err)
			continue
		}
		conn.Write([]byte(fmt.Sprintf("id: %d send %d byte, data: %s\n", cid, n, string(data[:n]))))
	}
}
