package redis

import (
	"github.com/google/wire"
)

// ProviderRedisSet is Redis providers.
var ProviderRedisSet = wire.NewSet(
	New,
)
