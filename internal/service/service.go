package service

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
)

// Service struct.
type Service struct {
	repo *repo.Repo
}

// NewService new service.
func NewService(repo *repo.Repo) *Service {
	s := &Service{
		repo: repo,
	}

	return s
}

// Close the Service.
func (s *Service) Close() error {
	group := new(errgroup.Group)

	group.Go(func() error {
		return s.repo.Close()
	})

	return group.Wait()
}

// RegisterGrpcServerHandler adds a serviceHandler.
func (s *Service) RegisterGrpcServerHandler(grpcServer *grpc.Server) {

}

// RegisterGatewayServerHandler adds a serviceHandler.
func (s *Service) RegisterGatewayServerHandler(gatewayServer *runtime.ServeMux) {

}
