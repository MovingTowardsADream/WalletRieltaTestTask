package httpserver

import "time"

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.httpServer.Addr = port
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.httpServer.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.httpServer.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
