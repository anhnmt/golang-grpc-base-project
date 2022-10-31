package repo

import (
	"github.com/google/wire"
)

// ProviderRepoSet is Repo providers.
var ProviderRepoSet = wire.NewSet(
	NewRepo,
)
