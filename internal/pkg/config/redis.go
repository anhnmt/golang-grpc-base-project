package config

func RedisEnabled() bool {
	return Default().RedisEnabled
}

func RedisAddress() string {
	return Default().RedisAddress
}

func RedisPassword() string {
	return Default().RedisPassword
}

func RedisDB() int {
	return Default().RedisDB
}
