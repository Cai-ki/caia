package cruntime

import (
	"context"

	"github.com/Cai-ki/caia/internal/cactor"
	"github.com/Cai-ki/caia/internal/cregistry"
	"github.com/Cai-ki/caia/internal/ctypes"
)

var (
	RootActor    ctypes.Actor
	RootRegistry ctypes.Registry
	Config       map[string]interface{}
)

var (
	MsgStart ctypes.Message = ctypes.Message{
		Payload: nil,
		ReplyTo: nil,
	}
)

func init() {
	RootActor = cactor.NewManager("root", 1, context.Background(), func(context.Context, ctypes.Message) {

	})
	RootRegistry = cregistry.NewManager("root")
	Config = map[string]interface{}{}
}

func Start() {
	RootActor.Start()
	RootActor.SendMessageToChildren(MsgStart)
}

func Stop() {
	RootActor.Stop()
}
