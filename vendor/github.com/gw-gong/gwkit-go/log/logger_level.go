package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerLevel string

const (
	LoggerLevelDebug LoggerLevel = "debug"
	LoggerLevelInfo  LoggerLevel = "info"
	LoggerLevelWarn  LoggerLevel = "warn"
	LoggerLevelError LoggerLevel = "error"
)

func MapLoggerLevel(level LoggerLevel) zapcore.Level {
	switch level {
	case LoggerLevelDebug:
		return zap.DebugLevel
	case LoggerLevelInfo:
		return zap.InfoLevel
	case LoggerLevelWarn:
		return zap.WarnLevel
	case LoggerLevelError:
		return zap.ErrorLevel
	}
	Errorf("unknown logger level: %s", level)
	return zap.InfoLevel
}
