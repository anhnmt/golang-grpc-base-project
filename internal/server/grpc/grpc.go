package grpc

import (
	"context"
	"time"

	metrics "github.com/grpc-ecosystem/go-grpc-middleware/providers/openmetrics/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/providers/opentracing/v2"
	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/tracing"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	permissionmodel "github.com/xdorro/golang-grpc-base-project/internal/module/permission/model"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/casbin"
	"github.com/xdorro/golang-grpc-base-project/pkg/redis"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
)

type Server struct {
	logPayload    bool
	seederService bool

	// options
	permissionCollection *mongo.Collection
	redis                *redis.Redis
	casbin               *casbin.Casbin
}

func NewGrpcServer(repo *repo.Repo, redis *redis.Redis, casbin *casbin.Casbin, service *service.Service) *grpc.Server {
	s := &Server{
		logPayload:           viper.GetBool("log.payload"),
		seederService:        viper.GetBool("seeder.service"),
		permissionCollection: repo.CollectionModel(&permissionmodel.Permission{}),
		redis:                redis,
		casbin:               casbin,
	}

	logger := grpczerolog.InterceptorLogger(log.Logger)
	optracing := opentracing.InterceptorTracer()
	opmetrics := metrics.NewServerMetrics()

	streamInterceptors := []grpc.StreamServerInterceptor{
		tracing.StreamServerInterceptor(optracing),
		metrics.StreamServerInterceptor(opmetrics),
		logging.StreamServerInterceptor(logger),
		recovery.StreamServerInterceptor(),
	}
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		tracing.UnaryServerInterceptor(optracing),
		metrics.UnaryServerInterceptor(opmetrics),
		logging.UnaryServerInterceptor(logger),
		recovery.UnaryServerInterceptor(),
	}

	// log payload if enabled
	if s.logPayload {
		payloadDecider := func(
			ctx context.Context, fullMethodName string, servingObject interface{},
		) logging.PayloadDecision {
			return logging.LogPayloadRequestAndResponse
		}

		streamInterceptors = append(streamInterceptors, logging.PayloadStreamServerInterceptor(logger, payloadDecider, time.RFC3339))
		unaryInterceptors = append(unaryInterceptors, logging.PayloadUnaryServerInterceptor(logger, payloadDecider, time.RFC3339))
	}

	// register grpc service Server
	grpcServer := grpc.NewServer(
		grpc.ChainStreamInterceptor(streamInterceptors...),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	)

	// register gRPC Server handler
	service.RegisterGrpcServerHandler(grpcServer)

	// seeder Service
	if s.seederService {
		// seeder service info
		service.SeederServiceInfo(grpcServer)
	}

	reflection.Register(grpcServer)

	return grpcServer
}
