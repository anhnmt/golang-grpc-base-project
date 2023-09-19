package main

import (
	"log/slog"

	"github.com/anhnmt/golang-grpc-base-project/pkg/logger"
)

func main() {
	logger.NewLogger("logs/data.log")

	slog.Info("Hello world")

	slog.Error(
		"An error occurred while processing the request",
		slog.String("url", "https://example.com"),
	)
}
