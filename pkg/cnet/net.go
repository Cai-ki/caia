package cnet

import (
	"github.com/Cai-ki/caia/internal/clog"
	"github.com/Cai-ki/caia/pkg/cruntime"
)

func init() {
	config, err := LoadConfig(ConfigPath)
	if err != nil {
		clog.Fatal("net: config load error:", err)
	}
	cruntime.Configs[KeyConfig] = config

	clog.Info(KeyConfig, ": ", *config)

	_, err = cruntime.RootActor.CreateChild(NetActorName, 1, ListenTCPHandle)
	if err != nil {
		clog.Fatal("net: init fail")
	}
}
