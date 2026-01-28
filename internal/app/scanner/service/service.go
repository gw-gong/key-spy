package service

import (
	"context"
)

// ScannerService 扫描服务接口
type ScannerService interface {
	// RunOnce 执行单次扫描
	RunOnce(ctx context.Context)
	// RunWithCron 启动定时扫描任务
	RunWithCron(ctx context.Context)
}
