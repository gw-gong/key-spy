package main

import (
	"github.com/gw-gong/gwkit-go/setting"
	"github.com/gw-gong/gwkit-go/util"
)

func main() {
	ctx := setting.GetServiceContext()

	cfgOption, err := initFlags()
	util.ExitOnErr(ctx, err)

	server, cleanup, err := InitServer(cfgOption)
	util.ExitOnErr(ctx, err)
	defer cleanup()

	server.SetupAndRun(ctx)
}
