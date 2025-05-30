package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/kelein/trove-fiber/docs"
	"github.com/kelein/trove-fiber/internal/inject"
	"github.com/kelein/trove-fiber/pkg/config"
	"github.com/kelein/trove-fiber/pkg/log"
	"github.com/kelein/trove-fiber/pkg/version"
)

var (
	v   = flag.Bool("v", false, "show the binary version")
	ver = flag.Bool("version", false, "show the binary version")
	cfg = flag.String("conf", "config/dev.yaml", "config file path")
)

func init() { initSwaggerInfo() }

// initSwaggerInfo setup swagger info
//
// * Basic Info
// @license.name Apache 2.0
// @contact.name trove-gin
// @contact.url https://github.com/kelein/trove-gin
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// * Authentication Info
// @securityDefinitions.apiKey Bearer
// @name Authorization
// @in header
func initSwaggerInfo() {
	docs.SwaggerInfo.BasePath = "v1"
	docs.SwaggerInfo.Host = "localhost:8000"
	docs.SwaggerInfo.Title = "Trove API Server"
	docs.SwaggerInfo.Version = version.AppVersion
	docs.SwaggerInfo.Description = version.String()
	docs.SwaggerInfo.InfoInstanceName = "swagger"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

func main() {
	flag.Parse()
	showVersion()

	conf := config.NewConfig(*cfg)
	log.SetupSlog(conf)

	app, cleanup, err := inject.NewWire(conf)
	defer cleanup()
	if err != nil {
		slog.Error("wire injection failed", "error", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf("http://%s:%d", conf.GetString("http.host"), conf.GetInt("http.port"))
	slog.Info("server start listen", "addr", addr)
	slog.Info("swagger docs", "addr", fmt.Sprintf("%s/swagger/index.html", addr))
	if err = app.Run(context.Background()); err != nil {
		slog.Error("server run failed", "error", err)
		os.Exit(1)
	}
}

func showVersion() {
	if *v || *ver {
		fmt.Println(version.String())
		os.Exit(0)
	}
}
