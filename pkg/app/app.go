package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelein/trove-fiber/pkg/server"
)

// Option for config option of the App
type Option func(a *App)

// App stands for application
type App struct {
	name    string
	servers []server.Server
}

// NewApp initializes a new App instance with options
func NewApp(opts ...Option) *App {
	a := &App{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// WithServer sets the servers for the App
func WithServer(servers ...server.Server) Option {
	return func(a *App) { a.servers = servers }
}

// WithName sets the name of the App
func WithName(name string) Option {
	return func(a *App) { a.name = name }
}

// Run bootstraps the application and all servers
func (a *App) Run(ctx context.Context) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	for _, serv := range a.servers {
		go func(serv server.Server) {
			if err := serv.Start(ctx); err != nil {
				slog.Error("Server start failed", "error", err)
			}
		}(serv)
	}

	select {
	case <-quit:
		slog.Info("Received SIGTERM, exiting gracefully...")
	case <-ctx.Done():
		slog.Info("Server context canceled")
	}

	for _, serv := range a.servers {
		if err := serv.Stop(ctx); err != nil {
			slog.Error("Server stop failed", "error", err)
		}
	}
	return nil
}
