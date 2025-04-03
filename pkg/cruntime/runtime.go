package cruntime

import (
	"context"

	"github.com/Cai-ki/caia/internal/cactor"
	"github.com/Cai-ki/caia/internal/clog"
	"github.com/Cai-ki/caia/internal/cregistry"
	"github.com/Cai-ki/caia/internal/ctypes"
)

var (
	RootActor ctypes.Actor
	Registrys map[string]ctypes.Registry
	Configs   map[string]interface{}
)

var (
	MsgStart ctypes.Message = ctypes.Message{
		Payload: nil,
		ReplyTo: nil,
	}
)

func init() {
	RootActor = cactor.NewManager(RootActorName, 1, context.WithValue(context.Background(), KeyManager, nil), func(context.Context, ctypes.Message) {
	})
	rootRegistry := cregistry.NewManager(RootActorName)
	Registrys = map[string]ctypes.Registry{}
	Registrys[RootActorName] = rootRegistry
	Configs = map[string]interface{}{}

	config, err := LoadConfig(ConfigPath)
	if err != nil {
		clog.Fatal("runtime: config load error:", err)
	}
	Configs[KeyConfig] = config

	clog.Info(KeyConfig, ": ", *config)
}

func Start() {
	RootActor.Start()
	RootActor.SendMessageToChildren(MsgStart)
}

func Stop() {
	RootActor.Stop()
}
