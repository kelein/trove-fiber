package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Server stands for a http server implement
type Server struct {
	host string
	port int
	*fiber.App
}

// Option setup the server
type Option func(s *Server)

// NewServer creates a new HTTP server instance
func NewServer(app *fiber.App, opts ...Option) *Server {
	s := &Server{App: app}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithHost setup host of the gRPC server
func WithHost(host string) Option {
	return func(s *Server) { s.host = host }
}

// WithPort setup port of the gRPC server
func WithPort(port int) Option {
	return func(s *Server) { s.port = port }
}

// Start makes the HTTP server running on the given address
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	err := s.Listen(addr)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Failed to listen on", "addr", addr, "error", err)
	}
	return nil
}

// Stop makes the HTTP server gracefully stop
func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := s.Shutdown(); err != nil {
		slog.Error("Failed to shutdown server", "error", err)
	}
	return nil
}
