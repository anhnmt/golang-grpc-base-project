package roleservice

import (
	"github.com/google/wire"
)

// ProviderServiceSet is Service providers.
var ProviderServiceSet = wire.NewSet(
	NewService,
)
