package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Set is a setter for any.
func Set(r redis.UniversalClient, key string, value any, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := r.Set(ctx, key, value, expiration).Err(); err != nil {
		log.Err(err).Msg("Failed to set keys")
		return err
	}

	return nil
}

// SetProto is a setter for the proto.
func SetProto(r redis.UniversalClient, key string, value proto.Message, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	bytes, err := protojson.Marshal(value)
	if err != nil {
		return err
	}

	if err = r.Set(ctx, key, bytes, expiration).Err(); err != nil {
		log.Err(err).Msg("Failed to set keys")
		return err
	}

	return nil
}

// SetObject is a setter for the object.
func SetObject(r redis.UniversalClient, key string, value any, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err = r.Set(ctx, key, bytes, expiration).Err(); err != nil {
		log.Err(err).Msg("Failed to set keys")
		return err
	}

	return nil
}

// Exists checks if the keys exists.
func Exists(r redis.UniversalClient, keys ...string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	count := r.Exists(ctx, keys...).Val()
	return count > 0
}

// Del deletes the keys.
func Del(r redis.UniversalClient, keys ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := r.Del(ctx, keys...).Err(); err != nil {
		log.Err(err).Msg("Failed to delete keys")
		return err
	}

	return nil
}

func Get(r redis.UniversalClient, key string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	val, err := r.Get(ctx, key).Result()
	if err != nil {
		log.Err(err).Msg("Failed to get key")
	}

	return val
}
