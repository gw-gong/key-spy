package reporter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gw-gong/key-spy/internal/config/scanner/localcfg"
	"github.com/gw-gong/key-spy/internal/pkg/model"

	"github.com/gw-gong/gwkit-go/log"
)

type reporter struct {
	cfg *localcfg.Config
}

func NewReporter(cfg *localcfg.Config) Reporter {
	return &reporter{
		cfg: cfg,
	}
}

func (r *reporter) GenerateReport(ctx context.Context, report *model.ScanReport) (filePath string, err error) {
	// 确保输出目录存在
	if err := os.MkdirAll(r.cfg.Output.Dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s.txt", r.cfg.Output.FilePrefix, timestamp)
	filePath = filepath.Join(r.cfg.Output.Dir, fileName)

	// 生成报告内容
	content := r.formatReport(report)

	// 写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write report file: %w", err)
	}

	log.Infoc(ctx, "Report generated", log.Str("file_path", filePath))

	return filePath, nil
}

func (r *reporter) formatReport(report *model.ScanReport) string {
	var sb strings.Builder

	// 报告头部
	sb.WriteString("=" + strings.Repeat("=", 79) + "\n")
	sb.WriteString("                           KEY-SPY 扫描报告\n")
	sb.WriteString("=" + strings.Repeat("=", 79) + "\n\n")

	// 基本信息
	sb.WriteString("【扫描信息】\n")
	sb.WriteString(fmt.Sprintf("  目标网站: %s\n", report.TargetURL))
	sb.WriteString(fmt.Sprintf("  搜索关键词: %s\n", strings.Join(report.Keywords, ", ")))
	sb.WriteString(fmt.Sprintf("  开始时间: %s\n", report.StartTime))
	sb.WriteString(fmt.Sprintf("  结束时间: %s\n", report.EndTime))
	sb.WriteString(fmt.Sprintf("  耗时: %s\n", report.Duration))
	sb.WriteString("\n")

	// 统计信息
	sb.WriteString("【统计摘要】\n")
	sb.WriteString(fmt.Sprintf("  扫描页面总数: %d\n", report.TotalPages))
	sb.WriteString(fmt.Sprintf("  匹配关键词页面数: %d\n", report.MatchPages))
	sb.WriteString(fmt.Sprintf("  扫描错误数: %d\n", report.ErrorCount))
	sb.WriteString("\n")

	// 匹配结果
	if len(report.Results) > 0 {
		sb.WriteString("-" + strings.Repeat("-", 79) + "\n")
		sb.WriteString("                           匹配结果详情\n")
		sb.WriteString("-" + strings.Repeat("-", 79) + "\n\n")

		// 按出现次数排序
		sortedResults := make([]*model.ScanResult, len(report.Results))
		copy(sortedResults, report.Results)
		sort.Slice(sortedResults, func(i, j int) bool {
			return sortedResults[i].TotalCount > sortedResults[j].TotalCount
		})

		for i, result := range sortedResults {
			sb.WriteString(fmt.Sprintf("[%d] URL: %s\n", i+1, result.URL))
			sb.WriteString(fmt.Sprintf("    页面深度: %d\n", result.Depth))
			sb.WriteString(fmt.Sprintf("    关键词总出现次数: %d\n", result.TotalCount))
			sb.WriteString(fmt.Sprintf("    出现的关键词: %s\n", strings.Join(result.Keywords, ", ")))
			sb.WriteString("    各关键词统计:\n")
			for keyword, count := range result.KeywordCounts {
				sb.WriteString(fmt.Sprintf("      - %s: %d 次\n", keyword, count))
			}
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("【匹配结果】\n")
		sb.WriteString("  未找到包含关键词的页面。\n\n")
	}

	// 报告尾部
	sb.WriteString("=" + strings.Repeat("=", 79) + "\n")
	sb.WriteString("                           报告结束\n")
	sb.WriteString("=" + strings.Repeat("=", 79) + "\n")

	return sb.String()
}
