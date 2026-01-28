//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/gw-gong/key-spy/internal/app/scanner/crawler"
	"github.com/gw-gong/key-spy/internal/app/scanner/notifier"
	"github.com/gw-gong/key-spy/internal/app/scanner/reporter"
	"github.com/gw-gong/key-spy/internal/app/scanner/service"
	"github.com/gw-gong/key-spy/internal/config/scanner/localcfg"

	"github.com/google/wire"
	"github.com/gw-gong/gwkit-go/hotcfg"
)

var ConfigSet = wire.NewSet(
	localcfg.NewConfig,
	hotcfg.NewHotLoaderManager,
)

var BizSet = wire.NewSet(
	crawler.NewCrawler,
	reporter.NewReporter,
	notifier.NewNotifier,
	service.NewScannerService,
)

var ServerSet = wire.NewSet(
	wire.Struct(new(Server), "*"),
)

func InitServer(cfgOption *hotcfg.LocalConfigOption) (*Server, func(), error) {
	wire.Build(
		ConfigSet,
		BizSet,
		ServerSet,
	)
	return nil, nil, nil
}
