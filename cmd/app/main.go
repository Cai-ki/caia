package main

import (
	"github.com/Cai-ki/caia/internal/clog"
	_ "github.com/Cai-ki/caia/pkg/cnet"
	"github.com/Cai-ki/caia/pkg/cruntime"
)

func main() {
	clog.SetLevel(clog.WARN)
	cruntime.Start()
	<-cruntime.RootActor.GetContext().Done()
}
