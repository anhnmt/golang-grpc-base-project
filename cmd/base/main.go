package main

import (
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/xdorro/golang-grpc-base-project/config"
	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
)

func init() {
	// -env is option for command line
	env := flag.String("env", "local", "environment")
	// -log_file is option for command line
	logFile := flag.String("log_path", "logs/data.log", "log file path")
	flag.Parse()

	logger.NewLogger(*logFile)
	config.NewConfig(*env)
}

func main() {
	log.Info().Msg("Hello world")
}
