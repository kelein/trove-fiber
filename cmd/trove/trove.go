package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

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

func init() {
	initTracerProvider()
	docs.InitSwaggerInfo()
}

func initTracerProvider() *sdktrace.TracerProvider {
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		slog.Error("initialize otel exporter failed", "error", err)
		os.Exit(1)
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(version.AppName),
				semconv.ServiceVersionKey.String(version.AppVersion),
			),
		),
	)
	otel.SetTracerProvider(provider)
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(prop)
	return provider
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
	slog.Info("server start listen on", "addr", addr)
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
