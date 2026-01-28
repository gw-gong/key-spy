package util

import (
	"context"
	"runtime/debug"

	"github.com/gw-gong/gwkit-go/log"
)

func ExitOnErr(ctx context.Context, err error) {
	if err != nil {
		log.Errorc(ctx, "ExitOnErr", log.Str("stack", string(debug.Stack())), log.Err(err))
		panic("ExitOnErr: " + err.Error())
	}
}
