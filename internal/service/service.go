package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authservice "github.com/xdorro/golang-grpc-base-project/internal/module/auth/service"
	permissionmodel "github.com/xdorro/golang-grpc-base-project/internal/module/permission/model"
	permissionservice "github.com/xdorro/golang-grpc-base-project/internal/module/permission/service"
	roleservice "github.com/xdorro/golang-grpc-base-project/internal/module/role/service"
	userservice "github.com/xdorro/golang-grpc-base-project/internal/module/user/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/redis"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
	authv1 "github.com/xdorro/golang-grpc-base-project/proto/pb/auth/v1"
	permissionv1 "github.com/xdorro/golang-grpc-base-project/proto/pb/permission/v1"
	rolev1 "github.com/xdorro/golang-grpc-base-project/proto/pb/role/v1"
	userv1 "github.com/xdorro/golang-grpc-base-project/proto/pb/user/v1"
)

// Service struct.
type Service struct {
	mu    sync.Mutex
	repo  *repo.Repo
	redis *redis.Redis

	// services
	userService       *userservice.Service
	authService       *authservice.Service
	roleService       *roleservice.Service
	permissionService *permissionservice.Service
}

// NewService new service.
func NewService(
	repo *repo.Repo,
	redis *redis.Redis,

	// services
	userService *userservice.Service,
	authService *authservice.Service,
	roleService *roleservice.Service,
	permissionService *permissionservice.Service,
) *Service {
	s := &Service{
		repo:              repo,
		redis:             redis,
		userService:       userService,
		authService:       authService,
		roleService:       roleService,
		permissionService: permissionService,
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

	return group.Wait()
}

// RegisterGrpcServerHandler adds a serviceHandler.
func (s *Service) RegisterGrpcServerHandler(grpcServer *grpc.Server) {
	userv1.RegisterUserServiceServer(grpcServer, s.userService)
	authv1.RegisterAuthServiceServer(grpcServer, s.authService)
	rolev1.RegisterRoleServiceServer(grpcServer, s.roleService)
	permissionv1.RegisterPermissionServiceServer(grpcServer, s.permissionService)
}

// RegisterGatewayServerHandler adds a serviceHandler.
func (s *Service) RegisterGatewayServerHandler(gatewayServer *runtime.ServeMux) error {
	ctx := context.Background()

	appAddress := viper.GetString("app.address")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// we're going to run the different protocol servers in parallel, so
	// make an errgroup
	group := new(errgroup.Group)

	group.Go(func() error {
		s.mu.Lock()
		defer s.mu.Unlock()
		return userv1.RegisterUserServiceHandlerFromEndpoint(ctx, gatewayServer, appAddress, opts)
	})

	group.Go(func() error {
		s.mu.Lock()
		defer s.mu.Unlock()
		return authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, gatewayServer, appAddress, opts)
	})

	group.Go(func() error {
		s.mu.Lock()
		defer s.mu.Unlock()
		return rolev1.RegisterRoleServiceHandlerFromEndpoint(ctx, gatewayServer, appAddress, opts)
	})

	group.Go(func() error {
		s.mu.Lock()
		defer s.mu.Unlock()
		return permissionv1.RegisterPermissionServiceHandlerFromEndpoint(ctx, gatewayServer, appAddress, opts)
	})

	return group.Wait()
}

// SeederServiceInfo
func (s *Service) SeederServiceInfo(grpcServer *grpc.Server) {
	services := make([]string, 0)
	for name, val := range grpcServer.GetServiceInfo() {
		for _, info := range val.Methods {
			services = append(services, fmt.Sprintf("/%s/%s", name, info.Name))
		}
	}

	if len(services) == 0 {
		return
	}

	permissionCollection := s.repo.CollectionModel(&permissionmodel.Permission{})

	// find all permissions with filter
	filter := bson.M{
		"deleted_at": bson.M{
			"$exists": false,
		},
		"slug": bson.M{
			"$in": services,
		},
	}

	// find all permissions with filter and option
	opt := options.
		Find().
		SetSort(bson.M{"created_at": -1})

	bulk := make([]any, 0)
	permissions, _ := repo.Find[permissionmodel.Permission](permissionCollection, filter, opt)

	for _, slug := range services {
		if ok := s.hasSlugInPermissions(permissions, slug); !ok {
			name := slug[strings.LastIndex(slug, "/")+1:]
			per := &permissionmodel.Permission{
				Name: name,
				Slug: slug,
			}
			per.PreCreate()

			bulk = append(bulk, per)
		}
	}

	if len(bulk) > 0 {
		_, err := repo.InsertMany(permissionCollection, bulk)
		if err != nil {
			log.Err(err).Msg("Error create permission")
		}

		// _ = s.redis.Del(context.Background(), constants.ListAuthPermissionsKey)

		log.Info().
			Interface("data", bulk).
			Msg("Insert permissions")
	}

}

// hasSlugInPermissions
func (s *Service) hasSlugInPermissions(permissions []*permissionmodel.Permission, slug string) bool {
	for _, permission := range permissions {
		if strings.EqualFold(slug, permission.Slug) {
			return true
		}
	}

	return false
}
