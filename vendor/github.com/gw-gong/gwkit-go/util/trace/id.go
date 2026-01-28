package trace

import (
	"github.com/gw-gong/gwkit-go/util/str"
)

func GenerateRequestID() string {
	return str.GenerateULID()
}

func GenerateTraceID() string {
	return str.GenerateULID()
}
