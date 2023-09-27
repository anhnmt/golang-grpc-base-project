package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"

	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/config"
	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/logger"
	"github.com/anhnmt/golang-grpc-base-project/internal/server"
	"github.com/anhnmt/golang-grpc-base-project/internal/utils"
	"github.com/anhnmt/golang-grpc-base-project/internal/wire"
)

var (
	logFile string
	env     string
)

func init() {
	pflag.StringVar(&logFile, "log-file", "", "log file path, ex: logs/data.log")
	pflag.StringVar(&env, "env", "", "environment")
	pflag.Parse()
}

func main() {
	logger.New(logFile)
	config.New(env)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := wire.InitServer()
	if err != nil {
		log.Panic().Msg("initial server failed")
	}

	go func(srv *server.Server) {
		if err = srv.Start(); err != nil {
			log.Panic().Msg("start server failed")
		}
	}(srv)

	// wait for termination signal
	wait := utils.GracefulShutdown(ctx, utils.DefaultShutdownTimeout, map[string]utils.Operation{
		"server": func(c context.Context) error {
			return srv.Close(c)
		},
	})
	<-wait

	log.Info().Msg("graceful shutdown complete")
}
