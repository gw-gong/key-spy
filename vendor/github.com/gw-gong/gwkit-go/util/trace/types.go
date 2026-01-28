package trace

type (
	ContextKeyRequestID = struct{}
	ContextKeyTraceID   = struct{}
)

const (
	LoggerFieldRequestID = "rid"
	LoggerFieldTraceID   = "tid"
)

const (
	HttpHeaderRequestID = "X-Request-Id"
	HttpHeaderTraceID   = "X-Trace-Id"
)
