package cruntime

import (
	"context"

	"github.com/Cai-ki/caia/internal/cfactory"
	"github.com/Cai-ki/caia/internal/cregistry"
	"github.com/Cai-ki/caia/internal/croutine"
	"github.com/Cai-ki/caia/internal/ctypes"
)

var (
	Factory  ctypes.Factory
	Routine  ctypes.Routine
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
	Factory = cfactory.NewManager("root")
	Routine = croutine.NewManager("root", 1, context.Background(), func(context.Context, ctypes.Message) {

	})
	Registry = cregistry.NewManager("root")
	Config = map[string]interface{}{}
}

func Start() {
	Routine.Start()
	Routine.SendMessageToChildren(MsgStart)
}

func Stop() {
	Routine.Stop()
}
