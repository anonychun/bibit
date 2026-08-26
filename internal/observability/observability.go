package observability

import (
	"log/slog"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/lib"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	do.Provide(bootstrap.Injector, NewObservability)
}

type IObservability interface {
	Logger() *slog.Logger
	Meter() metric.Meter
	Tracer() trace.Tracer
}

type Observability struct {
	logger *slog.Logger
	meter  metric.Meter
	tracer trace.Tracer
}

var _ IObservability = (*Observability)(nil)

func NewObservability(i do.Injector) (*Observability, error) {
	cfg := do.MustInvoke[*config.Config](i)
	serviceName := lib.GetModuleName()

	logger, err := newLogger()
	if err != nil {
		return nil, err
	}

	meter, err := newMeter(serviceName)
	if err != nil {
		return nil, err
	}

	tracer, err := newTracer(serviceName, cfg.OTLP.Endpoint)
	if err != nil {
		return nil, err
	}

	return &Observability{
		logger: logger,
		meter:  meter,
		tracer: tracer,
	}, nil
}

func (o *Observability) Logger() *slog.Logger {
	return o.logger
}

func (o *Observability) Meter() metric.Meter {
	return o.meter
}

func (o *Observability) Tracer() trace.Tracer {
	return o.tracer
}
