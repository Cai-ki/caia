package cnet

import (
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Cai-ki/caia/internal/cactor"
	"github.com/Cai-ki/caia/internal/clog"
	"github.com/Cai-ki/caia/internal/ctypes"
	"github.com/Cai-ki/caia/pkg/cruntime"
)

func ConnectHandleFactory(conn *net.TCPConn) ctypes.HandleFunc {

	return func(actor ctypes.Actor, msg ctypes.Message) {
		// config := cruntime.Configs[KeyConfig].(*Config)
		ctx := actor.GetContext()
		defer conn.Close()

		cid := ctx.Value(KeyCid).(uint32)
		writer, err := actor.CreateChild(strconv.Itoa(int(cid))+": writer", 10, WriteHandleFactory(conn))
		if err != nil {
			clog.Errorf("net: %s create writer fail", actor.GetName())
			actor.Stop()
			return
		}
		writer.Start()

		reader, err := actor.CreateChild(strconv.Itoa(int(cid))+": reader", 10, ReadHandleFactory(conn), cactor.WithValue("writer", writer))
		if err != nil {
			clog.Errorf("net: %s create reader fail", actor.GetName())
			actor.Stop()
			return
		}
		reader.Start()
		reader.SendMessage(cruntime.MsgStart)

		<-ctx.Done() //TODO connect actor 除了启动子actor外，还有负责各种msg的处理
	}
}

func ReadHandleFactory(conn *net.TCPConn) ctypes.HandleFunc {
	// config := cruntime.Configs[KeyConfig].(*Config)
	addTime := time.Duration(config.ReadDeadlineMs) * time.Millisecond
	return func(actor ctypes.Actor, msg ctypes.Message) {
		ctx := actor.GetContext()
		writer := ctx.Value("writer").(ctypes.Actor)
		sandbox := writer.GetMailbox()
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

func WriteHandleFactory(conn *net.TCPConn) ctypes.HandleFunc {
	// config := cruntime.Configs[KeyConfig].(*Config)
	addTime := time.Duration(config.WriteDeadlineMs) * time.Millisecond
	return func(actor ctypes.Actor, msg ctypes.Message) {
		ctx := actor.GetContext()
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
