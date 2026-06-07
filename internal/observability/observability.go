package observability

import (
	"log/slog"
	"os"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/util"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	do.Provide(bootstrap.Injector, NewObservability)
}

type IObservability interface {
	Tracer() trace.Tracer
	Meter() metric.Meter
	Logger() *slog.Logger
}

type Observability struct {
	tracer trace.Tracer
	meter  metric.Meter
	logger *slog.Logger
}

var _ IObservability = (*Observability)(nil)

func NewObservability(i do.Injector) (*Observability, error) {
	moduleName := util.GetModuleName()

	tracer := otel.Tracer(moduleName)
	meter := otel.Meter(moduleName)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &Observability{
		tracer: tracer,
		meter:  meter,
		logger: logger,
	}, nil
}

func (o *Observability) Tracer() trace.Tracer {
	return o.tracer
}

func (o *Observability) Meter() metric.Meter {
	return o.meter
}

func (o *Observability) Logger() *slog.Logger {
	return o.logger
}
