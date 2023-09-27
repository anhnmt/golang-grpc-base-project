package utils

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
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

		log.Info().Msg("shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Panic().Msg(fmt.Sprintf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds()))
		})
		defer timeoutFunc.Stop()

		var wg sync.WaitGroup
		// Do the operations asynchronously to save time
		for innerKey, innerOp := range ops {
			wg.Add(1)
			func() {
				defer wg.Done()

				log.Info().Msg(fmt.Sprintf("cleaning up: %s", innerKey))
				if err := innerOp(ctx); err != nil {
					log.Panic().Msg(fmt.Sprintf("%s: clean up failed: %s", innerKey, err.Error()))
				}

				log.Info().Msg(fmt.Sprintf("%s was shutdown gracefully", innerKey))
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}
