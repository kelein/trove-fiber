package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Slogger stands for a logger middleware of Fiber
func Slogger() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()
		err := ctx.Next()
		slog.Info("HTTP Request",
			slog.String("method", ctx.Method()),
			slog.String("path", ctx.Path()),
			slog.Int("code", ctx.Response().StatusCode()),
			slog.String("client", ctx.IP()),
			slog.String("agent", ctx.Get(fiber.HeaderUserAgent)),
			slog.Duration("latency", time.Since(start)),
		)
		return err
	}
}
