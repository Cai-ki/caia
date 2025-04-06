package cnet

import (
	"errors"
	"io"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/Cai-ki/caia/internal/clog"
	"github.com/Cai-ki/caia/internal/ctypes"
	"github.com/Cai-ki/caia/pkg/cprotocol"
	"github.com/panjf2000/ants/v2"
)

func ConnectHandleFactory(conn *net.TCPConn) ctypes.HandleFunc {
	return func(actor ctypes.Actor, msg ctypes.Message) {
		// config := cruntime.Configs[KeyConfig].(*Config)
		ctx := actor.GetContext()

		codec := cprotocol.NewCodec(&cprotocol.TLVHandler{}, func() interface{} {
			return &cprotocol.TLVPacket{Value: make([]byte, 0, 128)}
		})

		// cid := ctx.Value(KeyCid).(uint32)
		addTime := time.Duration(config.ReadDeadlineMs) * time.Millisecond

		sandbox := ctypes.NewMailbox(10)

		ants.Submit(func() {
			defer actor.StopWithErase()
			defer conn.Close()
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
						} else if errors.Is(err, io.EOF) {
							clog.Debugf("net: %s connect closed", conn.RemoteAddr())
						} else {
							clog.Error("net: read error:", err)
						}
						return
					}
					iface, sz, err := codec.Decode(data[:n])
					if err != nil {
						// data not enough
						continue
					}
					_, ok := iface.(*cprotocol.TLVPacket)
					if !ok {
						continue
					}
					sandbox.SendResult(&cprotocol.TLVPacket{
						Type:   0,
						Length: uint32(len([]byte("hello, world!"))),
						Value:  []byte("hello, world!"),
					}, nil)
					clog.Debug("read", sz, "bytes data") //, string(pkt.Value))
				}
			}
		})

		ants.Submit(func() {
			for {
				select {
				case <-ctx.Done():
					return
				case msg, _ := <-sandbox.PriChan():
					pkt, _ := msg.Payload.(ctypes.Result).Data.(*cprotocol.TLVPacket)
					data, err := codec.Encode(pkt)
					codec.Release(pkt)
					if err != nil {
						clog.Error(err)
						return
					}

					for len(data) > 0 {
						//conn.SetWriteDeadline(time.Now().Add(addTime))
						n, err := conn.Write(data)
						data = data[n:]
						if err != nil {
							if os.IsTimeout(err) {
								continue
							} else if err == io.ErrClosedPipe || err == io.EOF {
								clog.Debug("net: write close:", err)
							} else if errors.Is(err, syscall.EPIPE) {
								clog.Error(err, "data:", (data), n)
								actor.Stop()
							}
							return
						}
					}
					codec.PutBuffer(data)
				}
			}
		})
	}
}
