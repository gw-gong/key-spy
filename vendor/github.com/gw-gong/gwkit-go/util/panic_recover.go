package util

import (
	"context"
	"runtime/debug"

	"github.com/gw-gong/gwkit-go/log"
)

type PanicHandler func(err interface{})

// WithRecover is a function that recovers from a panic and calls the panic handler.
// !!! Only one panicHandler needs to be passed in; if multiple are provided, only the first one will be used. !!!
func WithRecover(f func(), panicHandlers ...PanicHandler) {
	defer func() {
		if err := recover(); err != nil {
			if len(panicHandlers) > 0 && panicHandlers[0] != nil {
				panicHandlers[0](err)
			} else {
				DefaultPanicHandler(err)
			}
		}
	}()

	f()
}

func DefaultPanicHandler(err interface{}) {
	log.Error("panic", log.Any("err", err), log.Str("stack", string(debug.Stack())))
}

func DefaultPanicWithCtx(ctx context.Context, err interface{}) {
	log.Errorc(ctx, "panic", log.Any("err", err), log.Str("stack", string(debug.Stack())))
}
