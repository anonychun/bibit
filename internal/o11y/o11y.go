package o11y

import (
	"context"
	"log/slog"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/lib"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
)

func init() {
	do.Provide(bootstrap.Injector, NewO11y)
	do.MustInvoke[*O11y](bootstrap.Injector)
}

type IO11y interface {
	Logger() *slog.Logger
	Meter() metric.Meter
	Tracer() trace.Tracer
}

type O11y struct {
	meterProvider  *sdkmetric.MeterProvider
	tracerProvider *sdktrace.TracerProvider

	logger *slog.Logger
	meter  metric.Meter
	tracer trace.Tracer
}

var _ IO11y = (*O11y)(nil)

func NewO11y(i do.Injector) (*O11y, error) {
	cfg := do.MustInvoke[*config.Config](i)
	serviceName := lib.GetModuleName()

	logger, err := newLogger()
	if err != nil {
		return nil, err
	}

	meterProvider, err := newMeterProvider()
	if err != nil {
		return nil, err
	}
	meter := meterProvider.Meter(serviceName)

	tracerProvider, err := newTracerProvider(serviceName, cfg.OTLP.Endpoint)
	if err != nil {
		return nil, err
	}
	tracer := tracerProvider.Tracer(serviceName)

	return &O11y{
		meterProvider:  meterProvider,
		tracerProvider: tracerProvider,

		logger: logger,
		meter:  meter,
		tracer: tracer,
	}, nil
}

func (o *O11y) Logger() *slog.Logger {
	return o.logger
}

func (o *O11y) Meter() metric.Meter {
	return o.meter
}

func (o *O11y) Tracer() trace.Tracer {
	return o.tracer
}

func (o *O11y) Shutdown(ctx context.Context) error {
	slog.Info("shutting down o11y")

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return o.meterProvider.Shutdown(ctx) })
	g.Go(func() error { return o.tracerProvider.Shutdown(ctx) })

	return g.Wait()
}
