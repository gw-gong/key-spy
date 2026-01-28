package localcfg

import (
	"github.com/gw-gong/gwkit-go/hotcfg"
	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/setting"
	"github.com/gw-gong/key-spy/internal/pkg/client/wechat"
)

type Config struct {
	hotcfg.BaseConfigCapable
	Env      setting.Env           `yaml:"env" mapstructure:"env"`
	Scanner  *ScannerConfig        `yaml:"scanner" mapstructure:"scanner"`
	Cron     *CronConfig           `yaml:"cron" mapstructure:"cron"`
	Output   *OutputConfig         `yaml:"output" mapstructure:"output"`
	Notifier *NotifierConfig       `yaml:"notifier" mapstructure:"notifier"`
	Logger   *log.LoggerConfig     `yaml:"logger" mapstructure:"logger"`
}

type ScannerConfig struct {
	TargetURL         string   `yaml:"target_url" mapstructure:"target_url"`
	Keywords          []string `yaml:"keywords" mapstructure:"keywords"`
	MaxDepth          int      `yaml:"max_depth" mapstructure:"max_depth"`
	RequestTimeoutMs  int      `yaml:"request_timeout_ms" mapstructure:"request_timeout_ms"`
	RequestIntervalMs int      `yaml:"request_interval_ms" mapstructure:"request_interval_ms"`
	MaxConcurrent     int      `yaml:"max_concurrent" mapstructure:"max_concurrent"`
	UserAgent         string   `yaml:"user_agent" mapstructure:"user_agent"`
}

type CronConfig struct {
	Spec    string `yaml:"spec" mapstructure:"spec"`
	Enabled bool   `yaml:"enabled" mapstructure:"enabled"`
}

type OutputConfig struct {
	Dir        string `yaml:"dir" mapstructure:"dir"`
	FilePrefix string `yaml:"file_prefix" mapstructure:"file_prefix"`
}

type NotifierConfig struct {
	Enabled       bool                  `yaml:"enabled" mapstructure:"enabled"`               // 是否启用通知
	WechatWebhook *wechat.WebhookConfig `yaml:"wechat_webhook" mapstructure:"wechat_webhook"` // 企业微信 Webhook 配置
}

func (c *Config) LoadConfig() {
	if err := c.Unmarshal(&c); err != nil {
		log.Error("unmarshal config failed", log.Err(err))
		return
	}

	log.Info("LoadConfig", log.Any("config", c))
}

func NewConfig(cfgOption *hotcfg.LocalConfigOption) (config *Config, err error) {
	config = &Config{}
	config.BaseConfigCapable, err = hotcfg.NewLocalBaseConfigCapable(cfgOption)
	config.LoadConfig()
	return config, err
}
