package utils

import (
	"context"
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
	DefaultShutdownTimeout = 10 * time.Second
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

		log.Info().Msg("Shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Panic().Msgf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for innerKey, innerOp := range ops {
			wg.Add(1)
			func() {
				defer wg.Done()

				log.Info().Msgf("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Error().Msgf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Info().Msgf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}
