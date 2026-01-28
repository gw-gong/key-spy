package wechat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gw-gong/gwkit-go/log"
)

// WebhookClient 企业微信 Webhook 客户端
type WebhookClient struct {
	webhookURL string
	httpClient *http.Client
}

// WebhookConfig 企业微信 Webhook 配置
type WebhookConfig struct {
	URL     string `yaml:"url" mapstructure:"url"`         // 完整的 Webhook URL
	Timeout int    `yaml:"timeout" mapstructure:"timeout"` // 请求超时时间（毫秒），默认 5000
}

// NewWebhookClient 创建企业微信 Webhook 客户端
func NewWebhookClient(cfg *WebhookConfig) *WebhookClient {
	if cfg == nil || cfg.URL == "" {
		return nil
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 5000
	}

	return &WebhookClient{
		webhookURL: cfg.URL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Millisecond,
		},
	}
}

// TextMessage 文本消息
type TextMessage struct {
	Content             string   `json:"content"`                         // 文本内容，最长不超过 2048 个字节
	MentionedList       []string `json:"mentioned_list,omitempty"`        // userid 列表，@指定用户
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"` // 手机号列表，@指定用户
}

// MarkdownMessage Markdown 消息
type MarkdownMessage struct {
	Content string `json:"content"` // Markdown 内容，最长不超过 4096 个字节
}

// webhookRequest Webhook 请求结构
type webhookRequest struct {
	MsgType  string           `json:"msgtype"`
	Text     *TextMessage     `json:"text,omitempty"`
	Markdown *MarkdownMessage `json:"markdown,omitempty"`
}

// webhookResponse Webhook 响应结构
type webhookResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// SendText 发送文本消息
func (c *WebhookClient) SendText(ctx context.Context, msg *TextMessage) error {
	if c == nil {
		return fmt.Errorf("webhook client is not initialized")
	}

	req := &webhookRequest{
		MsgType: "text",
		Text:    msg,
	}
	return c.send(ctx, req)
}

// SendMarkdown 发送 Markdown 消息
func (c *WebhookClient) SendMarkdown(ctx context.Context, msg *MarkdownMessage) error {
	if c == nil {
		return fmt.Errorf("webhook client is not initialized")
	}

	req := &webhookRequest{
		MsgType:  "markdown",
		Markdown: msg,
	}
	return c.send(ctx, req)
}

// send 发送消息
func (c *WebhookClient) send(ctx context.Context, req *webhookRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request failed: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %w", err)
	}

	var result webhookResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("unmarshal response failed: %w", err)
	}

	if result.ErrCode != 0 {
		log.Errorc(ctx, "wechat webhook error",
			log.Int("errcode", result.ErrCode),
			log.Str("errmsg", result.ErrMsg),
		)
		return fmt.Errorf("wechat webhook error: code=%d, msg=%s", result.ErrCode, result.ErrMsg)
	}

	log.Debugc(ctx, "wechat webhook message sent successfully")
	return nil
}
