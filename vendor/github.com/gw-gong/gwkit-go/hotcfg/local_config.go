package hotcfg

type LocalConfig interface {
	WatchLocalConfig(loadConfig func())
}

type LocalConfigOption struct {
	FilePath string `json:"filePath" yaml:"filePath" mapstructure:"filePath"`
	FileName string `json:"fileName" yaml:"fileName" mapstructure:"fileName"`
	FileType string `json:"fileType" yaml:"fileType" mapstructure:"fileType"`
}
