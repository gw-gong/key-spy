package hotcfg

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type BaseConfigCapable interface {
	GetBaseConfig() *BaseConfig
	Unmarshal(v interface{}) error
	AsLocalConfig() LocalConfig
	AsConsulConfig() ConsulConfig
}

type LoadconfigType string

const (
	ConfigTypeLocal  LoadconfigType = "local"
	ConfigTypeConsul LoadconfigType = "consul"
)

type BaseConfig struct {
	Viper              *viper.Viper        `json:"-" yaml:"-" mapstructure:"-"`
	mux                sync.Mutex          `json:"-" yaml:"-" mapstructure:"-"`
	ConfigType         LoadconfigType      `json:"configType" yaml:"configType" mapstructure:"configType"`
	LocalConfigOption  *LocalConfigOption  `json:"localConfigOption" yaml:"localConfigOption" mapstructure:"localConfigOption"`
	ConsulConfigOption *ConsulConfigOption `json:"consulConfigOption" yaml:"consulConfigOption" mapstructure:"consulConfigOption"`
}

func NewLocalBaseConfigCapable(localConfig *LocalConfigOption) (BaseConfigCapable, error) {
	return newBaseConfig(withLocalConfig(localConfig))
}

func NewConsulBaseConfigCapable(consulConfig *ConsulConfigOption) (BaseConfigCapable, error) {
	return newBaseConfig(withConsulConfig(consulConfig))
}

type option func(*BaseConfig)

func withLocalConfig(localConfig *LocalConfigOption) option {
	return func(c *BaseConfig) {
		c.ConfigType = ConfigTypeLocal
		c.LocalConfigOption = localConfig
	}
}

func withConsulConfig(consulConfig *ConsulConfigOption) option {
	return func(c *BaseConfig) {
		c.ConfigType = ConfigTypeConsul
		c.ConsulConfigOption = consulConfig
	}
}

func newBaseConfig(opts ...option) (*BaseConfig, error) {
	c := &BaseConfig{
		Viper: viper.New(),
	}

	for _, opt := range opts {
		opt(c)
	}

	switch c.ConfigType {
	case ConfigTypeLocal:
		if c.LocalConfigOption == nil || c.LocalConfigOption.FilePath == "" || c.LocalConfigOption.FileName == "" || c.LocalConfigOption.FileType == "" {
			return nil, fmt.Errorf("local config is nil or file path, name, and type are required")
		}
		c.Viper.SetConfigType(c.LocalConfigOption.FileType)
		c.Viper.SetConfigName(c.LocalConfigOption.FileName)
		c.Viper.AddConfigPath(c.LocalConfigOption.FilePath)
		if err := c.Viper.ReadInConfig(); err != nil {
			return nil, err
		}
	case ConfigTypeConsul:
		if c.ConsulConfigOption == nil || c.ConsulConfigOption.ConsulAddr == "" || c.ConsulConfigOption.ConsulKey == "" || c.ConsulConfigOption.ConfigType == "" {
			return nil, fmt.Errorf("consul config is nil or addr, key, and type are required")
		}
		if err := c.Viper.AddRemoteProvider("consul", c.ConsulConfigOption.ConsulAddr, c.ConsulConfigOption.ConsulKey); err != nil {
			return nil, err
		}
		c.Viper.SetConfigType(c.ConsulConfigOption.ConfigType)
		if err := c.Viper.ReadRemoteConfig(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid config type: %s", c.ConfigType)
	}

	return c, nil
}

func (c *BaseConfig) GetBaseConfig() *BaseConfig {
	return c
}

func (c *BaseConfig) Unmarshal(v interface{}) error {
	return c.Viper.Unmarshal(v)
}

func (c *BaseConfig) AsLocalConfig() LocalConfig {
	if c.ConfigType == ConfigTypeLocal {
		return c
	}
	return nil
}

func (c *BaseConfig) WatchLocalConfig(loadConfig func()) {
	if c.ConfigType == ConfigTypeLocal {
		c.Viper.WatchConfig()
		c.Viper.OnConfigChange(func(e fsnotify.Event) {
			c.mux.Lock()
			defer c.mux.Unlock()
			loadConfig()
		})
	}
}

func (c *BaseConfig) AsConsulConfig() ConsulConfig {
	if c.ConfigType == ConfigTypeConsul {
		return c
	}
	return nil
}

func (c *BaseConfig) GetConsulReloadTime() int {
	return c.ConsulConfigOption.ReloadTime
}

func (c *BaseConfig) ReadConsulConfig() error {
	if err := c.Viper.ReadRemoteConfig(); err != nil {
		return fmt.Errorf("failed to read remote configuration: %w, consulAddr: %s, consulKey: %s, configType: %s",
			err, c.ConsulConfigOption.ConsulAddr, c.ConsulConfigOption.ConsulKey, c.ConsulConfigOption.ConfigType)
	}
	return nil
}

func (c *BaseConfig) CalculateConsulConfigHash() string {
	return CalculateConfigHash(c.Viper)
}
