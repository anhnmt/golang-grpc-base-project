package service

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	authservice "github.com/xdorro/golang-grpc-base-project/internal/module/auth/service"
	userservice "github.com/xdorro/golang-grpc-base-project/internal/module/user/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/redis"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
	authv1 "github.com/xdorro/golang-grpc-base-project/proto/pb/auth/v1"
	userv1 "github.com/xdorro/golang-grpc-base-project/proto/pb/user/v1"
)

// Service struct.
type Service struct {
	repo  *repo.Repo
	redis *redis.Redis

	// services
	userService *userservice.Service
	authService *authservice.Service
}

// NewService new service.
func NewService(
	repo *repo.Repo,
	redis *redis.Redis,

	// services
	userService *userservice.Service,
	authService *authservice.Service,
) *Service {
	s := &Service{
		repo:        repo,
		redis:       redis,
		userService: userService,
		authService: authService,
	}

	return s
}

// Close the Service.
func (s *Service) Close() error {
	group := new(errgroup.Group)

	group.Go(func() error {
		return s.repo.Close()
	})

	group.Go(func() error {
		return s.redis.Close()
	})

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

// RegisterGrpcServerHandler adds a serviceHandler.
func (s *Service) RegisterGrpcServerHandler(grpcServer *grpc.Server) {
	userv1.RegisterUserServiceServer(grpcServer, s.userService)
	authv1.RegisterAuthServiceServer(grpcServer, s.authService)
}

// RegisterGatewayServerHandler adds a serviceHandler.
func (s *Service) RegisterGatewayServerHandler(gatewayServer *runtime.ServeMux) error {
	ctx := context.Background()

	if err := userv1.RegisterUserServiceHandlerServer(ctx, gatewayServer, s.userService); err != nil {
		return err
	}

	if err := authv1.RegisterAuthServiceHandlerServer(ctx, gatewayServer, s.authService); err != nil {
		return err
	}

	return nil
}
