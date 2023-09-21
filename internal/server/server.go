package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/config"
)

type Server struct {
	grpcServer *grpc.Server
}

func New(
	grpcServer *grpc.Server,
) *Server {
	s := &Server{
		grpcServer: grpcServer,
	}

	return s
}

func (s *Server) Start() error {
	g, _ := errgroup.WithContext(context.Background())

	if config.PprofEnabled() {
		g.TryGo(func() error {
			addr := fmt.Sprintf(":%d", config.PprofPort())
			slog.Info(fmt.Sprintf("starting pprof http://localhost%s", addr))

			return http.ListenAndServe(addr, nil)
		})
	}

	// Serve the http server on the http listener.
	g.TryGo(func() error {
		addr := fmt.Sprintf(":%d", config.AppPort())
		slog.Info(fmt.Sprintf("starting application http://localhost%s", addr))

		// create new http server
		srv := &http.Server{
			Addr: addr,
			// Use h2c, so we can serve HTTP/2 without TLS.
			Handler: h2c.NewHandler(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
						s.grpcServer.ServeHTTP(w, r)
						return
					}
				}),
				&http2.Server{},
			),
			// ReadHeaderTimeout: 10 * time.Second,
			// ReadTimeout:       1 * time.Minute,
			// WriteTimeout:      1 * time.Minute,
			// MaxHeaderBytes:    8 * 1024, // 8KiB
		}

		// run the server
		return srv.ListenAndServe()
	})

	return g.Wait()
}

func (s *Server) Close(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	g.TryGo(func() error {
		s.grpcServer.GracefulStop()
		return nil
	})

	return g.Wait()
}
