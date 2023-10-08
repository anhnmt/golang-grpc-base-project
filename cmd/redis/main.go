package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"

	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/config"
	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/logger"
	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/redis"
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
	redis, err := redis.New(ctx)
	if err != nil {
		log.Fatal().Err(err).
			Msg("New redis failed")
	}

	err = redis.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		log.Fatal().Err(err).
			Msg("Set redis key failed")
	}

	val := redis.Get(ctx, "key").Val()

	log.Info().
		Interface("val", val).
		Msg("Query succeeded")

	// wait for termination signal
	wait := utils.GracefulShutdown(ctx, utils.DefaultShutdownTimeout, map[string]utils.Operation{
		"redis": func(c context.Context) error {
			return redis.Close()
		},
	})
	<-wait

	log.Info().Msg("Graceful shutdown complete")
}
