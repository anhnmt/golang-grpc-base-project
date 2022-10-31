package grpc

import (
	"github.com/google/wire"
)

// ProviderGrpcServerSet is gRPC Server providers.
var ProviderGrpcServerSet = wire.NewSet(
	NewGrpcServer,
)
