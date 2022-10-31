package grpc

import (
	"context"
	"time"

	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGrpcServer() *grpc.Server {
	logger := grpczerolog.InterceptorLogger(log.Logger)

	streamInterceptors := []grpc.StreamServerInterceptor{
		// tags.StreamServerInterceptor(tags.WithFieldExtractor(tags.CodeGenRequestFieldExtractor)),
		logging.StreamServerInterceptor(logger),
		recovery.StreamServerInterceptor(),
	}
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(logger),
		recovery.UnaryServerInterceptor(),
	}

	logPayload := viper.GetBool("log.payload")

	// log payload if enabled
	if logPayload {
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
		// grpc.Creds(tlsCredentials),
		grpc.ChainStreamInterceptor(streamInterceptors...),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	)

	reflection.Register(grpcServer)

	return grpcServer
}
