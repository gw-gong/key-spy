package service

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gw-gong/key-spy/internal/app/scanner/crawler"
	"github.com/gw-gong/key-spy/internal/app/scanner/notifier"
	"github.com/gw-gong/key-spy/internal/app/scanner/reporter"
	"github.com/gw-gong/key-spy/internal/config/scanner/localcfg"

	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/util/trace"
	"github.com/robfig/cron/v3"
)

type scannerService struct {
	cfg      *localcfg.Config
	crawler  crawler.Crawler
	reporter reporter.Reporter
	notifier notifier.Notifier
}

// NewScannerService 创建扫描服务
func NewScannerService(
	cfg *localcfg.Config,
	crawler crawler.Crawler,
	reporter reporter.Reporter,
	notifier notifier.Notifier,
) ScannerService {
	return &scannerService{
		cfg:      cfg,
		crawler:  crawler,
		reporter: reporter,
		notifier: notifier,
	}
}

// RunOnce 执行单次扫描
func (s *scannerService) RunOnce(ctx context.Context) {
	ctx = trace.WithLogFieldTraceID(ctx, trace.GenerateTraceID())
	log.Infoc(ctx, "Running single scan...")
	s.doScan(ctx)
	log.Infoc(ctx, "Single scan completed")
}

// RunWithCron 启动定时扫描任务
func (s *scannerService) RunWithCron(ctx context.Context) {
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc(s.cfg.Cron.Spec, func() {
		ctx = trace.WithLogFieldTraceID(ctx, trace.GenerateTraceID())
		log.Infoc(ctx, "Cron job triggered, starting scan...")
		s.doScan(ctx)
	})
	if err != nil {
		log.Errorc(ctx, "Failed to add cron job", log.Err(err))
		return
	}

	c.Start()
	log.Infoc(ctx, "Cron scheduler started", log.Str("spec", s.cfg.Cron.Spec))

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infoc(ctx, "Shutting down cron scheduler...")
	c.Stop()
	log.Infoc(ctx, "Cron scheduler stopped")
}

// doScan 执行扫描任务
func (s *scannerService) doScan(ctx context.Context) {
	// 执行扫描
	report, err := s.crawler.Scan(ctx)
	if err != nil {
		log.Errorc(ctx, "Scan failed", log.Err(err))
		return
	}

	// 生成报告
	filePath, err := s.reporter.GenerateReport(ctx, report)
	if err != nil {
		log.Errorc(ctx, "Failed to generate report", log.Err(err))
		return
	}

	log.Infoc(ctx, "Scan completed successfully",
		log.Str("report_file", filePath),
		log.Int("total_pages", report.TotalPages),
		log.Int("match_pages", report.MatchPages),
	)

	// 发送通知
	if err := s.notifier.Notify(ctx, report, filePath); err != nil {
		log.Errorc(ctx, "Failed to send notification", log.Err(err))
	}
}
