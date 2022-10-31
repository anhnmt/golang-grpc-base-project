package server

import (
	"github.com/google/wire"
)

// ProviderServerSet is Server providers.
var ProviderServerSet = wire.NewSet(
	NewServer,
)
