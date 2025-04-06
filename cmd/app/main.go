package main

import (
	"runtime"

	"github.com/Cai-ki/caia/internal/clog"
	_ "github.com/Cai-ki/caia/pkg/cnet"
	"github.com/Cai-ki/caia/pkg/cruntime"
)

func main() {
	clog.Info("CPU核心数:", runtime.NumCPU())
	clog.Info("GOMAXPROCS:", runtime.GOMAXPROCS(0))
	//clog.SetLevel(clog.WARN)
	cruntime.Start()
	<-cruntime.RootActor.GetContext().Done()
}
