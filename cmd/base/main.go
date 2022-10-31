package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/xdorro/golang-grpc-base-project/config"
	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
)

func main() {
	// -env is option for command line
	env := flag.String("env", "local", "environment")
	// -log_file is option for command line
	logFile := flag.String("log_path", "logs/data.log", "log file path")
	flag.Parse()

	logger.NewLogger(*logFile)
	config.NewConfig(*env)

	exit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	srv := initServer()

	// Run server
	go func(srv *server.Server) {
		err := srv.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to run http server")
		}
	}(srv)

	<-exit
	if err := srv.Close(); err != nil {
		log.Err(err).Msg("Failed to close server")
		return
	}

	log.Info().Msg("Graceful shutdown complete")
}
