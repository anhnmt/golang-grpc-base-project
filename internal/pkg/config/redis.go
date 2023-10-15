package config

import (
	"github.com/spf13/viper"
)

func RedisEnabled() bool {
	return viper.GetBool("REDIS_ENABLED")
}

func RedisAddress() string {
	return viper.GetString("REDIS_ADDRESS")
}

func RedisPassword() string {
	return viper.GetString("REDIS_PASSWORD")
}

func RedisDB() int {
	return viper.GetInt("REDIS_DB")
}
