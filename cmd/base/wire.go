//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"

	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/internal/server/gateway"
	"github.com/xdorro/golang-grpc-base-project/internal/server/grpc"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

func initServer() *server.Server {
	wire.Build(
		service.ProviderServiceSet,
		grpc.ProviderGrpcServerSet,
		gateway.ProviderGatewayServerSet,
		server.ProviderServerSet,
	)

	return &server.Server{}
}
