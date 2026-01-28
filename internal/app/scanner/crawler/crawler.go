package crawler

import (
	"context"

	"github.com/gw-gong/key-spy/internal/pkg/model"
)

// Crawler 爬虫接口
type Crawler interface {
	// Scan 扫描目标网站，返回扫描报告
	Scan(ctx context.Context) (*model.ScanReport, error)
}
