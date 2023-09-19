package main

import (
	"log/slog"

	"github.com/anhnmt/golang-grpc-base-project/internal/config"
	"github.com/anhnmt/golang-grpc-base-project/pkg/logger"
)

var (
	env     string
	logFile string
)

func main() {
	err := logger.NewLogger("logs/data.log")
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

	slog.Info("Hello world")

	slog.Error(
		"An error occurred while processing the request",
		slog.String("url", "https://example.com"),
	)
}
