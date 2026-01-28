package log

import (
	"context"

	"go.uber.org/zap"
)

type ctxKeyLogger struct{}

func setLoggerToCtx(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger{}, logger)
}

func getLoggerFromCtx(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ctxKeyLogger{}).(*zap.Logger)
	if !ok {
		return zap.L()
	}
	return logger
}
