package redis

import (
	"github.com/google/wire"
)

// ProviderRedisSet is redis providers.
var ProviderRedisSet = wire.NewSet(
	NewRedis,
)
