package casbin

import (
	"github.com/google/wire"
)

// ProviderCasbinSet is Casbin providers.
var ProviderCasbinSet = wire.NewSet(
	NewCasbin,
)
