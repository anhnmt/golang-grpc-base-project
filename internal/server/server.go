package server

import (
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/casbin"
)

// Server struct.
type Server struct {
	mu sync.Mutex

	// config
	appName      string
	appAddress   string
	pprofAddress string
	appDebug     bool
	corsEnabled  bool

	service       *service.Service
	casbin        *casbin.Casbin
	grpcServer    *grpc.Server
	gatewayServer *runtime.ServeMux
}

// NewServer new server.
func NewServer(
	service *service.Service, casbin *casbin.Casbin, grpcServer *grpc.Server, gatewayServer *runtime.ServeMux,
) *Server {
	s := &Server{
		appName:       viper.GetString("app.name"),
		appDebug:      viper.GetBool("app.debug"),
		appAddress:    viper.GetString("app.address"),
		pprofAddress:  viper.GetString("pprof.address"),
		corsEnabled:   viper.GetBool("cors.enabled"),
		service:       service,
		casbin:        casbin,
		grpcServer:    grpcServer,
		gatewayServer: gatewayServer,
	}

	return s
}

// Run runs the server.
func (s *Server) Run() error {
	// we're going to run the different protocol servers in parallel, so
	// make an errgroup
	group := new(errgroup.Group)

	// we need a webserver to get the pprof webserver
	if s.appDebug {
		group.Go(func() error {
			log.Info().Msgf("Starting pprof http://%s", s.pprofAddress)

			return http.ListenAndServe(s.pprofAddress, nil)
		})
	}

	// Serve the http server on the http listener.
	group.Go(func() error {
		log.Info().Msgf("Starting application http://%s", s.appAddress)

		// create new http server
		srv := &http.Server{
			Addr: s.appAddress,
			// Use h2c, so we can serve HTTP/2 without TLS.
			Handler: h2c.NewHandler(
				s.grpcHandlerFunc(),
				&http2.Server{},
			),
			ReadHeaderTimeout: time.Second,
			ReadTimeout:       1 * time.Minute,
			WriteTimeout:      1 * time.Minute,
			MaxHeaderBytes:    8 * 1024, // 8KiB
		}

		// run the server
		return srv.ListenAndServe()
	})

	return group.Wait()
}

func (s *Server) Close() error {
	s.grpcServer.GracefulStop()

	return s.service.Close()
}

func (s *Server) grpcHandlerFunc() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			s.grpcServer.ServeHTTP(w, r)
			return
		}

		var gwMux http.Handler = s.gatewayServer

		// add CORS if enabled
		if s.corsEnabled {
			gwMux = newCORS().Handler(gwMux)
		}

		gwMux.ServeHTTP(w, r)
	})
}
