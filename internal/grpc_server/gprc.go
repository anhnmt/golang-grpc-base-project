package grpc_server

import (
	"context"
	"log/slog"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"google.golang.org/grpc"

	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/config"
)

type Server struct {
	mu         sync.Mutex
	logPayload bool
}

func New() *grpc.Server {
	logEvents := []logging.LoggableEvent{
		logging.StartCall,
		logging.FinishCall,
	}

	// log payload if enabled
	if config.LogPayload() {
		logEvents = append(logEvents,
			logging.PayloadReceived,
			logging.PayloadSent,
		)
	}

	opts := []logging.Option{
		logging.WithLogOnEvents(logEvents...),
	}

	logger := InterceptorLogger(slog.Default())
	streamInterceptors := []grpc.StreamServerInterceptor{
		logging.StreamServerInterceptor(logger, opts...),
		recovery.StreamServerInterceptor(),
		validator.StreamServerInterceptor(),
	}
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(logger, opts...),
		recovery.UnaryServerInterceptor(),
		validator.UnaryServerInterceptor(),
	}

	// register grpc service Server
	grpcServer := grpc.NewServer(
		grpc.ChainStreamInterceptor(streamInterceptors...),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	)

	return grpcServer
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
