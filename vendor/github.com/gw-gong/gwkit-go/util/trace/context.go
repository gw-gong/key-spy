package trace

import (
	"context"

	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/setting"
)

func SetRequestIDToCtx(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKeyRequestID{}, requestID)
}

func GetRequestIDFromCtx(ctx context.Context) string {
	if value := ctx.Value(ContextKeyRequestID{}); value != nil {
		if requestID, ok := value.(string); ok {
			return requestID
		}
	}
	return ""
}

func SetTraceIDToCtx(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ContextKeyTraceID{}, traceID)
}

func GetTraceIDFromCtx(ctx context.Context) string {
	if value := ctx.Value(ContextKeyTraceID{}); value != nil {
		if traceID, ok := value.(string); ok {
			return traceID
		}
	}
	return ""
}

func WithLogFieldRequestID(ctx context.Context, requestID string) context.Context {
	return log.WithFields(ctx, log.Str(LoggerFieldRequestID, requestID))
}

func WithLogFieldTraceID(ctx context.Context, traceID string) context.Context {
	return log.WithFields(ctx, log.Str(LoggerFieldTraceID, traceID))
}

func CopyCtx(ctx context.Context) context.Context {
	newCtx := setting.GetServiceContext()
	if requestID := GetRequestIDFromCtx(ctx); requestID != "" {
		newCtx = SetRequestIDToCtx(newCtx, requestID)
		newCtx = WithLogFieldRequestID(newCtx, requestID)
	}
	if traceID := GetTraceIDFromCtx(ctx); traceID != "" {
		newCtx = SetTraceIDToCtx(newCtx, traceID)
		newCtx = WithLogFieldTraceID(newCtx, traceID)
	}
	return newCtx
}
