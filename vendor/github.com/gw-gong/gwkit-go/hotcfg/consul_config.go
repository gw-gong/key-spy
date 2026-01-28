package hotcfg

type ConsulConfig interface {
	GetConsulReloadTime() int
	ReadConsulConfig() error
	CalculateConsulConfigHash() string
}

type ConsulConfigOption struct {
	ConsulAddr string `json:"consul_addr" yaml:"consul_addr" mapstructure:"consul_addr"`
	ConsulKey  string `json:"consul_key" yaml:"consul_key" mapstructure:"consul_key"`
	ConfigType string `json:"config_type" yaml:"config_type" mapstructure:"config_type"`
	ReloadTime int    `json:"reload_time" yaml:"reload_time" mapstructure:"reload_time"` // second
}
