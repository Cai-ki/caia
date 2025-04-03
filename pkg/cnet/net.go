package cnet

import (
	"github.com/Cai-ki/caia/internal/clog"
	"github.com/Cai-ki/caia/pkg/cruntime"
)

var config *Config

func init() {
	c, err := LoadConfig(ConfigPath)
	if err != nil {
		clog.Fatal("net: config load error:", err)
	}
	config = c
	cruntime.Configs[KeyConfig] = config

	clog.Info(KeyConfig, ": ", *config)

	_, err = cruntime.RootActor.CreateChild(NetActorName, 1, ListenTCPHandle)
	if err != nil {
		clog.Fatal("net: init fail")
	}
}
