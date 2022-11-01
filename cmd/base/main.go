package main

import (
	"context"
	"errors"
	"flag"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/xdorro/golang-grpc-base-project/config"
	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
	"github.com/xdorro/golang-grpc-base-project/utils"
)

func main() {
	// -env is option for command line
	env := flag.String("env", "local", "environment")
	// -log_file is option for command line
	logFile := flag.String("log_path", "logs/data.log", "log file path")
	flag.Parse()

	logger.NewLogger(*logFile)
	config.NewConfig(*env)

	srv := initServer()

	// Run server
	go func(srv *server.Server) {
		err := srv.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to run http server")
		}
	}(srv)

	// wait for termination signal and register client & http server clean-up operations
	wait := utils.GracefulShutdown(context.Background(), utils.DefaultShutdownTimeout, map[string]utils.Operation{
		"server": func(ctx context.Context) error {
			return srv.Close()
		},
	})
	<-wait

	log.Info().Msg("Graceful shutdown complete")
}
