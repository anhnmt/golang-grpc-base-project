package grpc_server

import (
	"github.com/google/wire"
)

// ProviderGrpcServerSet is gRPC Server providers.
var ProviderGrpcServerSet = wire.NewSet(
	New,
)
