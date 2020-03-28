package service

import (
	"context"
	"os"
	"os/signal"
)

func signalCtx(ctx context.Context, sig ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		c := make(chan os.Signal, len(sig))
		signal.Notify(c, sig...)

		defer signal.Stop(c)

		select {
		case <-ctx.Done():
		case <-c:
			cancel()
		}
	}()

	return ctx
}
