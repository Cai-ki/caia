package cnet

import (
	"fmt"

	"github.com/Cai-ki/caia/pkg/cruntime"
)

func init() {
	config, err := LoadConfig(ConfigPath)
	if err != nil {
		panic(fmt.Sprintf("net: config load error: %s", err))
	}
	cruntime.Configs[KeyConfig] = config

	fmt.Println(KeyConfig, ": ", config)

	_, err = cruntime.RootActor.CreateChild(NetActorName, 1, ListenTCPHandle)
	if err != nil {
		panic("net: init fail")
	}
}
