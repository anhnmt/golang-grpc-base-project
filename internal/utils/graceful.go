package utils

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Operation is a cleanup function on shutting down
type Operation func(ctx context.Context) error

const (
	DefaultShutdownTimeout = 20 * time.Second
)

// GracefulShutdown is a utility to shut down a server gracefully
func GracefulShutdown(
	ctx context.Context, timeout time.Duration, ops map[string]Operation,
) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscall that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		slog.Info("Shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			slog.Error(fmt.Sprintf("Timeout %d ms has been elapsed, force exit", timeout.Milliseconds()))
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for innerKey, innerOp := range ops {
			wg.Add(1)
			func() {
				defer wg.Done()

				slog.Info(fmt.Sprintf("Cleaning up: %s", innerKey))
				if err := innerOp(ctx); err != nil {
					slog.Error(fmt.Sprintf("%s: clean up failed: %s", innerKey, err.Error()))
					return
				}

				slog.Info(fmt.Sprintf("%s was shutdown gracefully", innerKey))
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}
