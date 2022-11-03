package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	permissionmodel "github.com/xdorro/golang-grpc-base-project/internal/module/permission/model"
	"github.com/xdorro/golang-grpc-base-project/pkg/redis"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
	"github.com/xdorro/golang-grpc-base-project/utils"
)

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func (s *Server) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := info.Server.(auth.ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
		} else {
			authFunc := s.authInterceptor(info.FullMethod)
			newCtx, err = authFunc(ctx)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// StreamServerInterceptor returns a new unary server interceptors that performs per-request auth.
func (s *Server) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
	) error {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := srv.(auth.ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(stream.Context(), info.FullMethod)
		} else {
			authFunc := s.authInterceptor(info.FullMethod)
			newCtx, err = authFunc(stream.Context())
		}
		if err != nil {
			return err
		}
		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}

func (s *Server) authInterceptor(fullMethod string) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		// get full method
		permissions := s.getListPermissions()
		if len(permissions) == 0 {
			return ctx, nil
		}

		per, ok := permissions[fullMethod]
		if !ok || per == nil {
			return ctx, nil
		}

		// check require auth
		if per.RequireAuth {
			var token string
			token, err := auth.AuthFromMD(ctx, utils.TokenType)
			if err != nil {
				log.Err(err).Msg("Error get token from header")
				return ctx, err
			}

			var claims *jwt.RegisteredClaims
			claims, err = utils.DecryptToken(token)
			if err != nil {
				return ctx, err
			}

			// check role
			var role string
			if len(claims.Audience) > 0 {
				role = claims.Audience[0]
			}

			allowed, _ := s.casbin.Client().Enforce(role, fullMethod)
			if !allowed {
				err = fmt.Errorf("Permission denied")
				return ctx, err
			}
		}

		return ctx, nil
	}
}

// getAllPermissions returns all permissions.
func (s *Server) getListPermissions() map[string]*permissionmodel.Permission {
	// get all permissions
	permissions := make(map[string]*permissionmodel.Permission)

	if val := redis.Get(s.redis.Client(), utils.RedisKeyListAuthPermissions); val != "" {
		_ = json.Unmarshal([]byte(val), &permissions)
		return permissions
	}

	// count all permissions with filter
	filter := bson.M{
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	// find all permissions with filter and option
	opt := options.
		Find().
		SetSort(bson.M{"created_at": -1})

	data, err := repo.Find[permissionmodel.Permission](s.permissionCollection, filter, opt)
	if err != nil {
		return permissions
	}

	for _, per := range data {
		permissions[per.Slug] = per
	}

	log.Info().
		Interface("permissions", permissions).
		Msg("Log get all permissions")

	go func() {
		_ = redis.SetObject(s.redis.Client(), utils.RedisKeyListAuthPermissions, permissions, 7*24*time.Hour)
	}()

	return permissions
}
