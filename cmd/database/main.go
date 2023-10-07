package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"

	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/config"
	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/database"
	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/logger"
	"github.com/anhnmt/golang-grpc-base-project/internal/utils"
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

	ctx := context.Background()
	database, err := database.New(ctx)
	if err != nil {
		log.Fatal().Err(err).
			Msg("New database failed")
	}

	// database = database.Debug()
	//
	// first, err := database.Target.Query().
	// 	Where().
	// 	First(ctx)
	// if err != nil {
	// 	log.Fatal().Err(err).
	// 		Msg("Query failed")
	// 	return
	// }

	// log.Info().Interface("data", first).Msg("Query succeeded")

	// wait for termination signal
	wait := utils.GracefulShutdown(ctx, utils.DefaultShutdownTimeout, map[string]utils.Operation{
		"db": func(c context.Context) error {
			return database.Close()
		},
	})
	<-wait

	log.Info().Msg("Graceful shutdown complete")
}
