package main

import (
	"log/slog"

	"github.com/spf13/pflag"

	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/config"
	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/logger"
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
}
