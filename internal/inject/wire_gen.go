// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package inject

import (
	"github.com/google/wire"
	"github.com/kelein/trove-fiber/internal/handler"
	"github.com/kelein/trove-fiber/internal/repository"
	"github.com/kelein/trove-fiber/internal/server"
	"github.com/kelein/trove-fiber/internal/service"
	"github.com/kelein/trove-fiber/pkg/app"
	"github.com/kelein/trove-fiber/pkg/jwt"
	"github.com/kelein/trove-fiber/pkg/server/http"
	"github.com/kelein/trove-fiber/pkg/sid"
	"github.com/kelein/trove-fiber/pkg/version"
	"github.com/spf13/viper"
)

// Injectors from wire.go:

func NewWire(viperViper *viper.Viper) (*app.App, func(), error) {
	jwtJWT := jwt.NewJwt(viperViper)
	baseHandler := handler.NewBaseHandler()
	sidSid := sid.NewSid()
	db := repository.NewDB(viperViper)
	repositoryRepository := repository.NewRepository(db)
	transaction := repository.NewTransaction(repositoryRepository)
	serviceService := service.NewService(sidSid, jwtJWT, transaction)
	userRepository := repository.NewUserRepository(repositoryRepository)
	userService := service.NewUserService(serviceService, userRepository)
	userHandler := handler.NewUserHandler(baseHandler, userService)
	httpServer := server.NewHTTPServer(viperViper, jwtJWT, userHandler)
	appApp := newApp(httpServer)
	return appApp, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.NewRepository, repository.NewTransaction, repository.NewUserRepository)

var serviceSet = wire.NewSet(service.NewService, service.NewUserService)

var handlerSet = wire.NewSet(handler.NewBaseHandler, handler.NewUserHandler)

var serverSet = wire.NewSet(server.NewHTTPServer)

func newApp(httpServer *http.Server) *app.App {
	return app.NewApp(app.WithServer(httpServer), app.WithName(version.AppName))
}
