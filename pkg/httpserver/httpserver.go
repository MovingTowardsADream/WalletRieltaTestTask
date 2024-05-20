package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout     = 10 * time.Second
	_defaultWriteTimeout    = 10 * time.Second
	_defaultAddr            = ":8080"
	_defaultShutdownTimeout = 5 * time.Second
)

type Server struct {
	log             *slog.Logger
	httpServer      *http.Server
	shutdownTimeout time.Duration
}

func New(log *slog.Logger, handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         _defaultAddr,
	}

	s := &Server{
		log:             log,
		httpServer:      httpServer,
		shutdownTimeout: _defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic("cannot run http server: " + err.Error())
	}
}

func (s *Server) Run() error {
	const op = "httpserver.Run"

	l, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("http server started", slog.String("addr", l.Addr().String()))

	if err := s.httpServer.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	const op = "httpserver.Shutdown"

	s.log.With(slog.String("op", op)).
		Info("stopping http server", slog.String("port", s.httpServer.Addr))

	return fmt.Errorf("Shutdown - s.httpServer.Shutdown: %w", s.httpServer.Shutdown(ctx))
}
