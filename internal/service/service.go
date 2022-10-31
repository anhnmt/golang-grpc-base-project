package service

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// Service struct.
type Service struct {
}

// NewService new service.
func NewService() *Service {
	s := &Service{}

	return s
}

// Close the Service.
func (s *Service) Close() error {
	group := new(errgroup.Group)

	// group.Go(func() error {
	// 	return nil
	// })

	return group.Wait()
}

// RegisterGrpcServerHandler adds a serviceHandler.
func (s *Service) RegisterGrpcServerHandler(grpcServer *grpc.Server) {

}

// RegisterGatewayServerHandler adds a serviceHandler.
func (s *Service) RegisterGatewayServerHandler(gatewayServer *runtime.ServeMux) {

}
