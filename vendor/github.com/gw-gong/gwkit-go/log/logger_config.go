package log

const (
	DefaultOutputFilePath string = "logs/app.log"
)

type OutputEncoding string

const (
	OutputEncodingJSON    OutputEncoding = "json"
	OutputEncodingConsole OutputEncoding = "console"
)

type OutputToFileConfig struct {
	Enable     bool   `yaml:"enable" json:"enable" mapstructure:"enable"`                // Whether to enable output
	FilePath   string `yaml:"path" json:"path" mapstructure:"path"`                      // File path
	WithBuffer *bool  `yaml:"with_buffer" json:"with_buffer" mapstructure:"with_buffer"` // Whether to use buffer, will not be immediately flushed to the file
	MaxSize    int    `yaml:"max_size" json:"max_size" mapstructure:"max_size"`          // Maximum file size (MB)
	MaxBackups int    `yaml:"max_backups" json:"max_backups" mapstructure:"max_backups"` // Maximum number of backup files
	MaxAge     int    `yaml:"max_age" json:"max_age" mapstructure:"max_age"`             // Maximum retention days
	Compress   *bool  `yaml:"compress" json:"compress" mapstructure:"compress"`          // Whether to compress after rotation
}

type OutputToConsoleConfig struct {
	Enable   bool           `yaml:"enable" json:"enable" mapstructure:"enable"`       // Whether to enable output
	Encoding OutputEncoding `yaml:"encoding" json:"encoding" mapstructure:"encoding"` // Encoding
}

type LoggerConfig struct {
	Level           LoggerLevel           `yaml:"level" json:"level" mapstructure:"level"`                                     // Log level
	OutputToFile    OutputToFileConfig    `yaml:"output_to_file" json:"output_to_file" mapstructure:"output_to_file"`          // Output to file configuration
	OutputToConsole OutputToConsoleConfig `yaml:"output_to_console" json:"output_to_console" mapstructure:"output_to_console"` // Output to console configuration
	AddCaller       *bool                 `yaml:"add_caller" json:"add_caller" mapstructure:"add_caller"`                      // Whether to add caller information
}

func (config *LoggerConfig) GetWithBuffer() bool {
	if config.OutputToFile.WithBuffer == nil {
		return false
	}
	return *config.OutputToFile.WithBuffer
}

func (config *LoggerConfig) GetCompress() bool {
	if config.OutputToFile.Compress == nil {
		return true
	}
	return *config.OutputToFile.Compress
}

func (config *LoggerConfig) GetAddCaller() bool {
	if config.AddCaller == nil {
		return true
	}
	return *config.AddCaller
}

func IsSupportedEncoding(encoding OutputEncoding) bool {
	return encoding == OutputEncodingJSON || encoding == OutputEncodingConsole
}

func NewDefaultLoggerConfig() *LoggerConfig {
	defaultCompress := true
	defaultAddCaller := true
	defaultWithBuffer := false
	return &LoggerConfig{
		Level: LoggerLevelDebug,
		OutputToFile: OutputToFileConfig{
			Enable:     false,
			FilePath:   DefaultOutputFilePath,
			WithBuffer: &defaultWithBuffer,
			MaxSize:    500,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   &defaultCompress,
		},
		OutputToConsole: OutputToConsoleConfig{
			Enable:   true,
			Encoding: OutputEncodingConsole,
		},
		AddCaller: &defaultAddCaller,
	}
}

func mergeCfgIntoDefault(config *LoggerConfig) *LoggerConfig {
	if config == nil {
		return NewDefaultLoggerConfig()
	}

	defaultConfig := NewDefaultLoggerConfig()

	if config.Level == "" {
		config.Level = defaultConfig.Level
	}

	if config.OutputToFile.Enable {
		if config.OutputToFile.FilePath == "" {
			config.OutputToFile.FilePath = defaultConfig.OutputToFile.FilePath
		}
		if config.OutputToFile.WithBuffer == nil {
			config.OutputToFile.WithBuffer = defaultConfig.OutputToFile.WithBuffer
		}
		if config.OutputToFile.MaxSize <= 0 {
			config.OutputToFile.MaxSize = defaultConfig.OutputToFile.MaxSize
		}
		if config.OutputToFile.MaxBackups <= 0 {
			config.OutputToFile.MaxBackups = defaultConfig.OutputToFile.MaxBackups
		}
		if config.OutputToFile.MaxAge <= 0 {
			config.OutputToFile.MaxAge = defaultConfig.OutputToFile.MaxAge
		}
		if config.OutputToFile.Compress == nil {
			config.OutputToFile.Compress = defaultConfig.OutputToFile.Compress
		}
	}

	if !config.OutputToFile.Enable && !config.OutputToConsole.Enable {
		config.OutputToConsole.Enable = defaultConfig.OutputToConsole.Enable
	}

	if config.OutputToConsole.Enable {
		if config.OutputToConsole.Encoding == "" {
			config.OutputToConsole.Encoding = defaultConfig.OutputToConsole.Encoding
		}
	}

	if config.AddCaller == nil {
		config.AddCaller = defaultConfig.AddCaller
	}

	return config
}
