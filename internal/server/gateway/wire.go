package gateway

import (
	"github.com/google/wire"
)

// ProviderGatewayServerSet is Gateway Server providers.
var ProviderGatewayServerSet = wire.NewSet(
	NewGatewayServer,
)
