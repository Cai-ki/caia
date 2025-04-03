package cnet

import (
	"context"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Cai-ki/caia/internal/cactor"
	"github.com/Cai-ki/caia/internal/clog"
	"github.com/Cai-ki/caia/internal/ctypes"
	"github.com/Cai-ki/caia/pkg/cruntime"
)

func ConnectHandle(ctx context.Context, msg ctypes.Message) {
	// config := cruntime.Configs[KeyConfig].(*Config)
	conn := ctx.Value(KeyConnect).(*net.TCPConn)
	defer conn.Close()

	cid := ctx.Value(KeyCid).(uint32)
	manager := ctx.Value(KeyManager).(ctypes.Actor)
	writer, err := manager.CreateChild(strconv.Itoa(int(cid))+": writer", 10, WriteHandleFactory(ctx))
	if err != nil {
		clog.Errorf("net: %s create writer fail", manager.GetName())
		manager.Stop()
		return
	}
	writer.Start()

	reader, err := manager.CreateChild(strconv.Itoa(int(cid))+": reader", 10, ReadHandleFactory(ctx), cactor.WithValue("sandbox", writer.GetMailbox()))
	if err != nil {
		clog.Errorf("net: %s create reader fail", manager.GetName())
		manager.Stop()
		return
	}
	reader.Start()
	reader.SendMessage(cruntime.MsgStart)

	<-ctx.Done()
}

func ReadHandleFactory(ctx context.Context) ctypes.HandleFunc {
	// config := cruntime.Configs[KeyConfig].(*Config)
	conn := ctx.Value(KeyConnect).(*net.TCPConn)
	addTime := time.Duration(config.ReadDeadlineMs) * time.Millisecond
	return func(ctx context.Context, msg ctypes.Message) {
		sandbox := ctx.Value("sandbox").(ctypes.Mailbox)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				conn.SetReadDeadline(time.Now().Add(addTime))
				data := make([]byte, 1024)
				n, err := conn.Read(data)
				if err != nil {
					if os.IsTimeout(err) {
						continue
					}
					clog.Error("net: read error:", err)
					return
				}

				sandbox.SendResult(data[:n], nil)
				clog.Info("read", n, "bytes data:", string(data[:n]))
			}
		}
	}
}

func WriteHandleFactory(ctx context.Context) ctypes.HandleFunc {
	// config := cruntime.Configs[KeyConfig].(*Config)
	addTime := time.Duration(config.WriteDeadlineMs) * time.Millisecond
	conn := ctx.Value(KeyConnect).(*net.TCPConn)

	return func(ctx context.Context, msg ctypes.Message) {
		data := msg.Payload.(ctypes.Result).Data.([]byte)
		for len(data) > 0 {
			select {
			case <-ctx.Done():
				return
			default:
				conn.SetWriteDeadline(time.Now().Add(addTime))
				n, err := conn.Write(data)
				if err != nil {
					if os.IsTimeout(err) {
						data = data[n:]
						continue
					}
					clog.Error("net: write failed:", err)
				}
				data = data[n:]
			}
		}
	}
}
