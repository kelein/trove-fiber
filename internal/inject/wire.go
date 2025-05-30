//go:build wireinject
// +build wireinject

package inject

import (
	"github.com/google/wire"
	"github.com/spf13/viper"

	"github.com/kelein/trove-fiber/internal/handler"
	"github.com/kelein/trove-fiber/internal/repository"
	"github.com/kelein/trove-fiber/internal/server"
	"github.com/kelein/trove-fiber/internal/service"
	"github.com/kelein/trove-fiber/pkg/app"
	"github.com/kelein/trove-fiber/pkg/jwt"
	"github.com/kelein/trove-fiber/pkg/server/http"
	"github.com/kelein/trove-fiber/pkg/sid"
	"github.com/kelein/trove-fiber/pkg/version"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
)

var handlerSet = wire.NewSet(
	handler.NewBaseHandler,
	handler.NewUserHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
)

func newApp(httpServer *http.Server) *app.App {
	return app.NewApp(
		app.WithServer(httpServer),
		app.WithName(version.AppName),
	)
}

func NewWire(*viper.Viper) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
