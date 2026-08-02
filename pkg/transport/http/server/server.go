package pkg_http_server

import (
	"context"
	pkg_logger "couriers/pkg/logger"
	pkg_http_middleware "couriers/pkg/transport/http/middleware"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type HTTPServer struct {
	mux        *http.ServeMux
	config     HTTPServerConfig
	log        *pkg_logger.Logger
	middleware []pkg_http_middleware.Middleware
}

func NewHTTPServer(
	config HTTPServerConfig,
	log *pkg_logger.Logger,
	middleware ...pkg_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterRouters(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prefix := fmt.Sprintf("/api/%s/%s", router.apiVersion, router.personRole)

		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, router))
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := pkg_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	channelErrors := make(chan error, 1)

	go func() {
		defer close(channelErrors)

		s.log.Warn(" http server started : ", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			channelErrors <- err
		}
	}()

	select {
	case <-ctx.Done():
		s.log.Debug("http server shutdown!")

		ctxWithTime, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctxWithTime); err != nil {
			server.Close()

			return fmt.Errorf("http server shutdown error: %w", err)
		}

		s.log.Debug("http server closed!")

	case err := <-channelErrors:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
	}

	return nil
}
