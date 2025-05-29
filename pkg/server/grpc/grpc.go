package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
)

// Server stands for a gRPC server implementation
type Server struct {
	*grpc.Server
	host string
	port int
}

// Option setup the server
type Option func(s *Server)

// NewServer creates a new gRPC server instance
func NewServer(opts ...Option) *Server {
	s := &Server{Server: grpc.NewServer()}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithHost setup host of the gRPC server
func WithHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

// WithPort setup port of the gRPC server
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

// Start make the gRPC server running on the given address
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("Failed to listen on", "addr", addr, "error", err)
	}
	if err = s.Server.Serve(lis); err != nil {
		slog.Error("Failed to serve", "addr", addr, "error", err)
	}
	return nil
}

// Stop makes the gRPC server gracefully stop
func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	s.Server.GracefulStop()
	return nil
}
