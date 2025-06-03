package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// DefaultBuckets prometheus buckets in seconds
var DefaultBuckets = []float64{0.5, 1.0, 5.0}

const (
	requestMetricName = "http_request_total"
	latencyMetricName = "http_request_seconds"
)

// Metrics in prometheus defination
type Metrics struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// NewProme returns a new prometheus middleware of Fiber
func NewProme(app *fiber.App, jobName, metricsPath string) *Metrics {
	m := Metrics{}
	m.reqs = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   jobName,
		Name:        requestMetricName,
		Help:        "How many http requests processed, with labels status code, method and path.",
		ConstLabels: prometheus.Labels{"job": jobName},
	},
		[]string{"method", "path", "code"},
	)

	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   jobName,
		Name:        latencyMetricName,
		Help:        "How long it took to process the request, with labels status code, method and path.",
		ConstLabels: prometheus.Labels{"job": jobName},
		Buckets:     DefaultBuckets,
	},
		[]string{"method", "path", "code"},
	)

	prometheus.MustRegister(m.reqs, m.latency)

	// * Register metrics router of the middleware
	app.Get(metricsPath, func(ctx *fiber.Ctx) error {
		handler := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
		handler(ctx.Context())
		return nil
	})

	return &m
}

// Run start prometheus middleware with context
func (m *Metrics) Run() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Path() == "/metrics" {
			return ctx.Next()
		}

		start := time.Now()
		lvs := []string{ctx.Method(), ctx.Path(), strconv.Itoa(ctx.Response().StatusCode())}

		err := ctx.Next()

		m.reqs.WithLabelValues(lvs...).Inc()
		m.latency.WithLabelValues(lvs...).Observe(float64(time.Since(start).Nanoseconds()) / 1e9)
		return err
	}
}
