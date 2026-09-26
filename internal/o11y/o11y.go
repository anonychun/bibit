package o11y

import (
	"context"
	"log/slog"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/lib"
	"github.com/samber/do/v2"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
)

func init() {
	do.Provide(bootstrap.Injector, newO11y)
}

func Setup(i do.Injector) error {
	_, err := do.Invoke[*o11y](i)
	return err
}

type o11y struct {
	logger         *slog.Logger
	meterProvider  *sdkmetric.MeterProvider
	tracerProvider *sdktrace.TracerProvider
}

func newO11y(i do.Injector) (*o11y, error) {
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

	tracerProvider, err := newTracerProvider(serviceName, cfg.OTLP.Endpoint)
	if err != nil {
		return nil, err
	}

	return &o11y{
		logger:         logger,
		meterProvider:  meterProvider,
		tracerProvider: tracerProvider,
	}, nil
}

func (o *o11y) Shutdown(ctx context.Context) error {
	slog.Info("shutting down o11y")

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return o.meterProvider.Shutdown(ctx) })
	g.Go(func() error { return o.tracerProvider.Shutdown(ctx) })

	return g.Wait()
}
