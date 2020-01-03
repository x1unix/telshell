package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var once = &sync.Once{}

var (
	ctx        context.Context
	cancelFunc context.CancelFunc
)

// GetApplicationContext returns application context for graceful shutdown
func GetApplicationContext() (context.Context, context.CancelFunc) {
	once.Do(func() {
		ctx, cancelFunc = context.WithCancel(context.Background())

		go func() {
			signals := []os.Signal{syscall.SIGTERM, syscall.SIGINT, os.Interrupt}
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, signals...)
			defer signal.Reset(signals...)
			<-sigChan
			cancelFunc()
		}()
	})

	return ctx, cancelFunc
}
