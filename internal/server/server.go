package server

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/anhnmt/golang-grpc-base-project/ent"
	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/config"
)

type Server struct {
	redis      redis.UniversalClient
	database   *ent.Client
	grpcServer *grpc.Server
}

func New(
	redis redis.UniversalClient,
	database *ent.Client,
	grpcServer *grpc.Server,
) *Server {
	s := &Server{
		redis:      redis,
		database:   database,
		grpcServer: grpcServer,
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	if config.PprofEnabled() {
		g.TryGo(func() error {
			addr := fmt.Sprintf(":%d", config.PprofPort())
			log.Info().Msg(fmt.Sprintf("starting pprof http://localhost%s", addr))

			return http.ListenAndServe(addr, nil)
		})
	}

	// Serve the http server on the http listener.
	g.TryGo(func() error {
		addr := fmt.Sprintf(":%d", config.AppPort())
		log.Info().Msg(fmt.Sprintf("starting application http://localhost%s", addr))

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

	g.TryGo(func() error {
		return s.database.Close()
	})

	g.TryGo(func() error {
		return s.redis.Close()
	})

	return g.Wait()
}
