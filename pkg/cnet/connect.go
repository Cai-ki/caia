package cnet

import (
	"errors"
	"io"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/Cai-ki/caia/internal/cactor"
	"github.com/Cai-ki/caia/internal/clog"
	"github.com/Cai-ki/caia/internal/ctypes"
	"github.com/Cai-ki/caia/pkg/cprotocol"
	"github.com/Cai-ki/caia/pkg/cruntime"
)

func ConnectHandleFactory(conn *net.TCPConn) ctypes.HandleFunc {
	return func(actor ctypes.Actor, msg ctypes.Message) {
		// config := cruntime.Configs[KeyConfig].(*Config)
		ctx := actor.GetContext()
		// defer conn.Close()

		codec := cprotocol.NewCodec(&cprotocol.TLVHandler{}, func() interface{} {
			return &cprotocol.TLVPacket{Value: make([]byte, 0, 128)}
		})

		cid := ctx.Value(KeyCid).(uint32)
		writer, err := actor.CreateChild(strconv.Itoa(int(cid))+": writer", 10, WriteHandleFactory(conn, codec))
		if err != nil {
			clog.Errorf("net: %s create writer fail", actor.GetName())
			actor.Stop()
			return
		}
		writer.Start()

		reader, err := actor.CreateChild(strconv.Itoa(int(cid))+": reader", 10, ReadHandleFactory(conn, codec), cactor.WithValue("writer", writer))
		if err != nil {
			clog.Errorf("net: %s create reader fail", actor.GetName())
			actor.Stop()
			return
		}
		reader.Start()
		reader.SendMessage(cruntime.MsgStart)

		// <-ctx.Done() //TODO connect actor 除了启动子actor外，还有负责各种msg的处理
	}
}

func ReadHandleFactory(conn *net.TCPConn, codec *cprotocol.Codec) ctypes.HandleFunc {
	// config := cruntime.Configs[KeyConfig].(*Config)
	addTime := time.Duration(config.ReadDeadlineMs) * time.Millisecond
	return func(actor ctypes.Actor, msg ctypes.Message) {
		ctx := actor.GetContext()
		writer := ctx.Value("writer").(ctypes.Actor)
		sandbox := writer.GetMailbox()

		Pool.Submit(func() {
			defer actor.GetParent().StopWithErase()
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
						//actor.GetParent().Stop()
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
	}
}

func WriteHandleFactory(conn *net.TCPConn, codec *cprotocol.Codec) ctypes.HandleFunc {
	// config := cruntime.Configs[KeyConfig].(*Config)
	//addTime := time.Duration(config.WriteDeadlineMs) * time.Millisecond
	return func(actor ctypes.Actor, msg ctypes.Message) {
		//ctx := actor.GetContext()
		pkt, _ := msg.Payload.(ctypes.Result).Data.(*cprotocol.TLVPacket)
		data, err := codec.Encode(pkt)
		codec.Release(pkt)
		if err != nil {
			clog.Error(err)
			return
		}
		defer codec.PutBuffer(data)

		for len(data) > 0 {
			//conn.SetWriteDeadline(time.Now().Add(addTime))
			n, err := conn.Write(data)
			data = data[n:]
			if err != nil {
				if os.IsTimeout(err) {
					// data = data[n:]
					continue
				} else if err == io.ErrClosedPipe || err == io.EOF {
					clog.Debug("net: write close:", err)
				} else if errors.Is(err, syscall.EPIPE) {
					clog.Error(err, "data:", (data), n)
					// actor.Stop()
				}
				return
			}
		}
	}
}
