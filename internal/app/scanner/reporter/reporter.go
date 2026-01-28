package reporter

import (
	"context"

	"github.com/gw-gong/key-spy/internal/pkg/model"
)

// Reporter 报告生成器接口
type Reporter interface {
	// GenerateReport 生成扫描报告并保存到文件
	GenerateReport(ctx context.Context, report *model.ScanReport) (filePath string, err error)
}
