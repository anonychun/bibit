package observability

import (
	"log/slog"

	"github.com/anonychun/bibit/internal/bootstrap"
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
	Tracer() trace.Tracer
	Meter() metric.Meter
}

type Observability struct {
	logger *slog.Logger
	tracer trace.Tracer
	meter  metric.Meter
}

var _ IObservability = (*Observability)(nil)

func NewObservability(i do.Injector) (*Observability, error) {
	serviceName := lib.GetModuleName()
	logger, err := newLogger()
	if err != nil {
		return nil, err
	}

	tracer, err := newTracer(serviceName)
	if err != nil {
		return nil, err
	}

	meter, err := newMeter(serviceName)
	if err != nil {
		return nil, err
	}

	return &Observability{
		logger: logger,
		tracer: tracer,
		meter:  meter,
	}, nil
}

func (o *Observability) Logger() *slog.Logger {
	return o.logger
}

func (o *Observability) Tracer() trace.Tracer {
	return o.tracer
}

func (o *Observability) Meter() metric.Meter {
	return o.meter
}
