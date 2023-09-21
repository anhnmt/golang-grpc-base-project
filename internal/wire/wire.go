//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package wire

import (
	"github.com/google/wire"

	"github.com/anhnmt/golang-grpc-base-project/internal/grpc_server"
	"github.com/anhnmt/golang-grpc-base-project/internal/server"
)

func InitServer() (*server.Server, error) {
	wire.Build(
		grpc_server.ProviderGrpcServerSet,
		server.ProviderServerSet,
	)

	return &server.Server{}, nil
}
