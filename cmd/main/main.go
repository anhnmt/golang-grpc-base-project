package main

import (
	"log/slog"

	"github.com/spf13/pflag"

	"github.com/anhnmt/golang-grpc-base-project/internal/config"
	"github.com/anhnmt/golang-grpc-base-project/pkg/logger"
)

var (
	env     string
	logFile string
)

func init() {
	pflag.StringVarP(&env, "env", "e", "local", "environment")
	pflag.StringVarP(&logFile, "log-file", "l", "logs/data.log", "log file path, ex: logs/data.log")
	pflag.Parse()
}

func main() {
	err := logger.NewLogger(logFile)
	if err != nil {
		slog.Error("New logger failed",
			slog.Any("err", err),
		)
		return
	}

	err = config.NewConfig(env)
	if err != nil {
		slog.Error("Load config failed",
			slog.Any("err", err),
		)
		return
	}

	slog.Info("Hello world",
		slog.String("app_name", config.AppName()),
	)
}
