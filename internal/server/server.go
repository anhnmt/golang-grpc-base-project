package server

import (
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Server struct.
type Server struct {
	mu sync.Mutex

	// config
	appName    string
	appPort    int
	pprofPort  int
	appDebug   bool
	logPayload bool

	grpcServer    *grpc.Server
	gatewayServer *runtime.ServeMux
}

// NewServer new server.
func NewServer(grpcServer *grpc.Server, gatewayServer *runtime.ServeMux) *Server {
	s := &Server{
		grpcServer:    grpcServer,
		gatewayServer: gatewayServer,
	}

	return s
}

func (s *Server) Close() error {
	s.grpcServer.GracefulStop()
	return nil
}
