package server

import (
	"github.com/VictoriaMetrics/metrics"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/spf13/viper"

	"github.com/kelein/trove-fiber/internal/handler"
	"github.com/kelein/trove-fiber/pkg/jwt"
	"github.com/kelein/trove-fiber/pkg/server/http"
	"github.com/kelein/trove-fiber/pkg/version"
)

// NewHTTPServer create a new HTTP server instance
func NewHTTPServer(conf *viper.Viper, jwt *jwt.JWT, userHandler *handler.UserHandler) *http.Server {
	server := http.NewServer(
		fiber.New(),
		http.WithHost(conf.GetString("http.host")),
		http.WithPort(conf.GetInt("http.port")),
	)

	setupRouter(server.App, userHandler)
	return server
}

func setupRouter(app *fiber.App, userHandler *handler.UserHandler) {
	app.Use(etag.New())
	app.Use(pprof.New())
	app.Use(logger.New())
	app.Use(requestid.New())

	app.Get("/", index)
	app.Get("/version", index)
	app.Get("/healthz", index)
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/metrics", func(ctx *fiber.Ctx) error {
		metrics.WritePrometheus(ctx.Response().BodyWriter(), true)
		return nil
	})

	// s.GET("/index", func(ctx *gin.Context) {
	// 	ctx.Redirect(gohttp.StatusFound, "/swagger/index.html")
	// })

	// // Non-strict permission routing group
	// noStrictAuthRouter := v1.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
	// noStrictAuthRouter.GET("/user", userHandler.GetProfile)

	// // Strict permission routing group
	// strictAuthRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger))
	// strictAuthRouter.PUT("/user", userHandler.UpdateProfile)

	v1 := app.Group("/v1")
	v1.Post("/login", userHandler.Login)
	v1.Post("/register", userHandler.Register)
	v1.Get("/user", userHandler.GetProfile)
	v1.Put("/user", userHandler.UpdateProfile)
}

func index(ctx *fiber.Ctx) error { return ctx.JSON(version.Runtime()) }
