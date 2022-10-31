package server

import (
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
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

	service       *service.Service
	grpcServer    *grpc.Server
	gatewayServer *runtime.ServeMux
}

// NewServer new server.
func NewServer(service *service.Service, grpcServer *grpc.Server, gatewayServer *runtime.ServeMux) *Server {
	s := &Server{
		service:       service,
		grpcServer:    grpcServer,
		gatewayServer: gatewayServer,
	}

	return s
}

func (s *Server) Close() error {
	s.grpcServer.GracefulStop()

	return s.service.Close()
}
