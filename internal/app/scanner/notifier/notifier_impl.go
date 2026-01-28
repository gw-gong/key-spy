package notifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/gw-gong/key-spy/internal/config/scanner/localcfg"
	"github.com/gw-gong/key-spy/internal/pkg/client/wechat"
	"github.com/gw-gong/key-spy/internal/pkg/model"

	"github.com/gw-gong/gwkit-go/log"
)

type notifier struct {
	cfg           *localcfg.Config
	wechatClient  *wechat.WebhookClient
}

// NewNotifier 创建通知器
func NewNotifier(cfg *localcfg.Config) Notifier {
	var wechatClient *wechat.WebhookClient
	if cfg.Notifier != nil && cfg.Notifier.WechatWebhook != nil {
		wechatClient = wechat.NewWebhookClient(cfg.Notifier.WechatWebhook)
	}

	return &notifier{
		cfg:          cfg,
		wechatClient: wechatClient,
	}
}

// Notify 发送扫描完成通知
func (n *notifier) Notify(ctx context.Context, report *model.ScanReport, filePath string) error {
	if n.cfg.Notifier == nil || !n.cfg.Notifier.Enabled {
		log.Debugc(ctx, "notifier is disabled, skip sending notification")
		return nil
	}

	// 发送企业微信通知
	if n.wechatClient != nil {
		if err := n.sendWechatNotification(ctx, report, filePath); err != nil {
			log.Errorc(ctx, "failed to send wechat notification", log.Err(err))
			return err
		}
	}

	return nil
}

// sendWechatNotification 发送企业微信通知
func (n *notifier) sendWechatNotification(ctx context.Context, report *model.ScanReport, filePath string) error {
	content := n.formatMarkdownContent(report, filePath)

	msg := &wechat.MarkdownMessage{
		Content: content,
	}

	if err := n.wechatClient.SendMarkdown(ctx, msg); err != nil {
		return fmt.Errorf("send wechat markdown message failed: %w", err)
	}

	log.Infoc(ctx, "wechat notification sent successfully")
	return nil
}

// formatMarkdownContent 格式化 Markdown 内容
func (n *notifier) formatMarkdownContent(report *model.ScanReport, filePath string) string {
	var sb strings.Builder

	// 标题
	sb.WriteString("## 🔍 Key-Spy 扫描报告\n\n")

	// 基本信息
	sb.WriteString("### 扫描信息\n")
	sb.WriteString(fmt.Sprintf("> 目标网站: **%s**\n", report.TargetURL))
	sb.WriteString(fmt.Sprintf("> 关键词: `%s`\n", strings.Join(report.Keywords, "`, `")))
	sb.WriteString(fmt.Sprintf("> 扫描时间: %s\n", report.StartTime))
	sb.WriteString(fmt.Sprintf("> 耗时: %s\n\n", report.Duration))

	// 统计信息
	sb.WriteString("### 统计摘要\n")

	// 根据匹配结果设置状态颜色
	if report.MatchPages > 0 {
		sb.WriteString(fmt.Sprintf("> <font color=\"warning\">发现 %d 个页面包含关键词</font>\n", report.MatchPages))
	} else {
		sb.WriteString("> <font color=\"info\">未发现包含关键词的页面</font>\n")
	}

	sb.WriteString(fmt.Sprintf("> 扫描页面总数: **%d**\n", report.TotalPages))
	sb.WriteString(fmt.Sprintf("> 匹配页面数: **%d**\n", report.MatchPages))
	if report.ErrorCount > 0 {
		sb.WriteString(fmt.Sprintf("> <font color=\"warning\">错误数: %d</font>\n", report.ErrorCount))
	}
	sb.WriteString("\n")

	// 匹配结果摘要（最多显示 5 条）
	if len(report.Results) > 0 {
		sb.WriteString("### 匹配结果 TOP5\n")
		displayCount := len(report.Results)
		if displayCount > 5 {
			displayCount = 5
		}

		for i := 0; i < displayCount; i++ {
			result := report.Results[i]
			sb.WriteString(fmt.Sprintf("%d. [%s](%s) - 命中 **%d** 次\n",
				i+1, truncateURL(result.URL, 50), result.URL, result.TotalCount))
		}

		if len(report.Results) > 5 {
			sb.WriteString(fmt.Sprintf("\n> 更多结果请查看完整报告（共 %d 条）\n", len(report.Results)))
		}
		sb.WriteString("\n")
	}

	// 报告文件路径
	sb.WriteString(fmt.Sprintf("📄 报告文件: `%s`", filePath))

	return sb.String()
}

// truncateURL 截断 URL 显示
func truncateURL(url string, maxLen int) string {
	if len(url) <= maxLen {
		return url
	}
	return url[:maxLen-3] + "..."
}
