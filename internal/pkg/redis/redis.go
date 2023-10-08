package redis

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/config"
)

func New(ctx context.Context) (redis.UniversalClient, error) {
	if !config.RedisEnabled() {
		return nil, nil
	}

	addrs := strings.Split(config.RedisAddress(), ",")

	log.Info().
		Strs("address", addrs).
		// Int("db", r.db).
		Msg("Connecting to Redis")

	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    addrs,
		Password: config.RedisPassword(),
		DB:       config.RedisDB(),
		PoolSize: 100,
	})

	newCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := client.Ping(newCtx).Err()
	if err != nil {
		return nil, err
	}

	log.Info().Msg("Connecting to Redis successfully.")
	return client, nil
}
