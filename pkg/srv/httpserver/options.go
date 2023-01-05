package httpserver

import (
	"net"
	"time"
)

// Option -.
type Option func(*Server)

// Port -.
func Port(port string) Option {
	return func(s *Server) {
		s.httpServer.Addr = net.JoinHostPort("", port)
	}
}

// ReadTimeout -.
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.httpServer.ReadTimeout = timeout
	}
}

// WriteTimeout -.
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.httpServer.WriteTimeout = timeout
	}
}
