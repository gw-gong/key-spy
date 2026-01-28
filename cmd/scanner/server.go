package main

import (
	"context"

	"github.com/gw-gong/key-spy/internal/app/scanner/service"
	"github.com/gw-gong/key-spy/internal/config/scanner/localcfg"

	"github.com/gw-gong/gwkit-go/hotcfg"
	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/setting"
	"github.com/gw-gong/gwkit-go/util"
)

type Server struct {
	cfg            *localcfg.Config
	hlm            hotcfg.HotLoaderManager
	scannerService service.ScannerService
}

func (s *Server) SetupAndRun(ctx context.Context) {
	setting.SetEnv(s.cfg.Env)

	// 初始化全局日志
	syncFn, err := log.InitGlobalLogger(s.cfg.Logger)
	util.ExitOnErr(ctx, err)
	defer syncFn()

	// 启动热加载
	util.ExitOnErr(ctx, s.hlm.RegisterHotLoader(s.cfg))
	util.ExitOnErr(ctx, s.hlm.Watch())

	log.Infoc(ctx, "Key-Spy scanner started",
		log.Str("target_url", s.cfg.Scanner.TargetURL),
		log.Any("keywords", s.cfg.Scanner.Keywords),
		log.Bool("cron_enabled", s.cfg.Cron.Enabled),
	)

	if s.cfg.Cron.Enabled {
		// 定时任务模式
		s.scannerService.RunWithCron(ctx)
	} else {
		// 单次执行模式
		s.scannerService.RunOnce(ctx)
	}
}
