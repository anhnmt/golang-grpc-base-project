package main

import (
	"context"
	"log/slog"
	"os"

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

	slog.Info("Hello world",
		slog.String("app_name", config.AppName()),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := wire.InitServer()
	if err != nil {
		slog.Error("Initial server failed")
	}

	go func(srv *server.Server) {
		if err = srv.Start(); err != nil {
			slog.Error("Start server failed")
			os.Exit(0)
		}
	}(srv)

	// wait for termination signal
	wait := utils.GracefulShutdown(ctx, utils.DefaultShutdownTimeout, map[string]utils.Operation{
		"server": func(c context.Context) error {
			return srv.Close(c)
		},
	})
	<-wait

	slog.Info("Graceful shutdown complete")
}
