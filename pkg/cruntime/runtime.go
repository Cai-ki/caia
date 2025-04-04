package cruntime

import (
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

var config *Config

func init() {
	RootActor = cactor.NewManager(RootActorName, 1, nil, func(ctypes.Actor, ctypes.Message) {})
	rootRegistry := cregistry.NewManager(RootActorName)
	Registrys = map[string]ctypes.Registry{}
	Registrys[RootActorName] = rootRegistry
	Configs = map[string]interface{}{}

	c, err := LoadConfig(ConfigPath)
	if err != nil {
		clog.Fatal("runtime: config load error:", err)
	}
	config = c
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
