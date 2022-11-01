package redis

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var _ IRedis = (*Redis)(nil)

// IRedis is the interface that must be implemented by a redis.
type IRedis interface {
	Close() error
}

// Redis is a redis struct.
type Redis struct {
	mu       sync.Mutex
	addrs    []string
	password string
	db       int

	client redis.UniversalClient
}

// NewRedis is new redis.
func NewRedis() *Redis {
	redisURL := strings.Split(viper.GetString("redis.url"), " ")
	r := &Redis{
		addrs:    redisURL,
		password: viper.GetString("redis.password"), // no password set
		db:       viper.GetInt("redis.db"),          // use default DB
	}

	log.Info().
		Strs("redis_url", r.addrs).
		Int("redis_db", r.db).
		Msg("Connecting to Redis")

	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    r.addrs,
		Password: r.password,
		DB:       r.db,
		PoolSize: 1000,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Panic().Err(err).Msg("Connecting to Redis failed")
	}

	// Add client to redis
	r.setClient(client)

	log.Info().Msg("Connecting to Redis successfully.")

	return r
}

// Close closes the redis.
func (r *Redis) Close() error {
	if err := r.client.Close(); err != nil {
		log.Err(err).Msg("Failed to close from Redis")
		return err
	}

	return nil
}

// Client returns the Redis client
func (r *Redis) Client() redis.UniversalClient {
	return r.client
}

// setClient adds a new client to the redis.
func (r *Redis) setClient(client redis.UniversalClient) {
	r.mu.Lock()
	r.client = client
	r.mu.Unlock()
}
