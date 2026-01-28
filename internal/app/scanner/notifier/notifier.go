package notifier

import (
	"context"

	"github.com/gw-gong/key-spy/internal/pkg/model"
)

// Notifier 通知器接口
type Notifier interface {
	// Notify 发送扫描完成通知
	Notify(ctx context.Context, report *model.ScanReport, filePath string) error
}
