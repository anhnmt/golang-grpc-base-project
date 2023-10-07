package database

import (
	"github.com/google/wire"
)

// ProviderDatabaseSet is Database providers.
var ProviderDatabaseSet = wire.NewSet(
	New,
)
