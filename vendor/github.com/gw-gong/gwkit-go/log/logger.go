package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	// ensure that calling the log package at any time will produce output
	_, _ = InitGlobalLogger(NewDefaultLoggerConfig())
}

func InitGlobalLogger(loggerConfig *LoggerConfig) (func(), error) {
	logger, syncGlobalLogger, err := newLogger(loggerConfig)
	if err != nil {
		return syncGlobalLogger, err
	}

	zap.ReplaceGlobals(logger)
	return syncGlobalLogger, nil
}

func newLogger(loggerConfig *LoggerConfig) (*zap.Logger, func(), error) {
	// sync global logger
	var syncGlobalLogger = func() {
		_ = zap.L().Sync()
	}

	loggerConfig = mergeCfgIntoDefault(loggerConfig)

	zapCfg := zap.NewProductionConfig()

	zapCfg.Level.SetLevel(MapLoggerLevel(loggerConfig.Level))

	// create basic encoder configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// create encoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// create WriteSyncer
	var writeSyncer zapcore.WriteSyncer
	if loggerConfig.OutputToFile.Enable && loggerConfig.OutputToFile.FilePath != "" {
		// ensure directory exists
		dir := filepath.Dir(loggerConfig.OutputToFile.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, syncGlobalLogger, fmt.Errorf("failed to create log directory: %w", err)
		}

		// use lumberjack to split logs
		lumber := &lumberjack.Logger{
			Filename:   loggerConfig.OutputToFile.FilePath,
			MaxSize:    loggerConfig.OutputToFile.MaxSize,    // Unit: MB
			MaxBackups: loggerConfig.OutputToFile.MaxBackups, // Maximum number of old files to retain
			MaxAge:     loggerConfig.OutputToFile.MaxAge,     // Maximum number of days to retain old files
			Compress:   loggerConfig.GetCompress(),           // Whether to compress after rotation
		}

		// use buffer if enabled
		if loggerConfig.GetWithBuffer() {
			writeSyncer = &zapcore.BufferedWriteSyncer{
				WS: zapcore.AddSync(lumber),
			}
		} else {
			writeSyncer = zapcore.AddSync(lumber)
		}

		syncGlobalLogger = func() {
			if err := zap.L().Sync(); err != nil {
				if f, err := os.OpenFile(loggerConfig.OutputToFile.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
					defer f.Close()
					fmt.Fprintf(f, "Failed to sync global logger: %v\n", err)
				}
			}
		}
	} else if loggerConfig.OutputToConsole.Enable && IsSupportedEncoding(loggerConfig.OutputToConsole.Encoding) {
		writeSyncer = zapcore.AddSync(os.Stdout)
		if loggerConfig.OutputToConsole.Encoding == OutputEncodingConsole {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		syncGlobalLogger = func() {
			Info("Output to console, no sync needed, please ignore, no need to modify your code")
		}
	} else {
		return nil, syncGlobalLogger, errors.New("no valid output configured: either output_to_file or output_to_console must be enabled with valid settings")
	}

	// create core
	core := zapcore.NewCore(encoder, writeSyncer, zapCfg.Level)

	// create Logger options
	var loggerOptions []zap.Option
	if loggerConfig.GetAddCaller() {
		loggerOptions = append(loggerOptions, zap.AddCaller())
		// add caller skip
		loggerOptions = append(loggerOptions, zap.AddCallerSkip(1)) // due to the log functions are wrapped, so we need to skip one layer
	}

	// create logger
	logger := zap.New(core, loggerOptions...)

	return logger, syncGlobalLogger, nil
}
