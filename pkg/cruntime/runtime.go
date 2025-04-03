package cruntime

import (
	"context"

	"github.com/Cai-ki/caia/internal/cactor"
	"github.com/Cai-ki/caia/internal/cregistry"
	"github.com/Cai-ki/caia/internal/ctypes"
)

var (
	Root     ctypes.Actor
	Registry ctypes.Registry
	Config   map[string]interface{}
)

var (
	MsgStart ctypes.Message = ctypes.Message{
		Payload: nil,
		ReplyTo: nil,
	}
)

func init() {
	Root = cactor.NewManager("root", 1, context.Background(), func(context.Context, ctypes.Message) {

	})
	Registry = cregistry.NewManager("root")
	Config = map[string]interface{}{}
}

func Start() {
	Root.Start()
	Root.SendMessageToChildren(MsgStart)
}

func Stop() {
	Root.Stop()
}
