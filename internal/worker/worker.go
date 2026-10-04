package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/anonychun/bibit/internal/bootstrap"
	brokerKafka "github.com/anonychun/bibit/internal/broker/kafka"
	brokerRabbitmq "github.com/anonychun/bibit/internal/broker/rabbitmq"
	clientRiver "github.com/anonychun/bibit/internal/client/river"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/consumer"
	consumerHello "github.com/anonychun/bibit/internal/consumer/hello"
	jobHello "github.com/anonychun/bibit/internal/job/hello"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
)

func init() {
	do.Provide(bootstrap.Injector, NewWorker)
}

type IWorker interface {
	Start(ctx context.Context) error
}

type IBroker interface {
	Start(ctx context.Context, routes []consumer.Route) error
}

type Worker struct {
	riverClient clientRiver.IClient
	broker      IBroker
	routes      []consumer.Route
}

var _ IWorker = (*Worker)(nil)

func NewWorker(i do.Injector) (*Worker, error) {
	riverClient := do.MustInvoke[*clientRiver.Client](i)

	err := addWorkers(riverClient.Workers(),
		do.MustInvoke[*jobHello.Job](i),
	)
	if err != nil {
		return nil, err
	}

	broker, err := newBroker(i)
	if err != nil {
		return nil, err
	}

	routes := []consumer.Route{}
	err = addRoutes(&routes,
		consumer.Route{
			Topic:    "hello",
			Consumer: do.MustInvoke[*consumerHello.Consumer](i),
		},
	)
	if err != nil {
		return nil, err
	}

	return &Worker{
		riverClient: riverClient,
		broker:      broker,
		routes:      routes,
	}, nil
}

func (w *Worker) Start(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		slog.Info("starting worker")
		err := w.riverClient.Client().Start(ctx)
		if err != nil {
			return err
		}

		<-w.riverClient.Client().Stopped()
		return nil
	})

	g.Go(func() error {
		return w.startConsumer(ctx)
	})

	return g.Wait()
}

func (w *Worker) Shutdown(ctx context.Context) error {
	slog.Info("shutting down worker")
	return w.riverClient.Client().Stop(ctx)
}

func (w *Worker) startConsumer(ctx context.Context) error {
	if w.broker == nil {
		slog.Info("consumer disabled, no broker backend configured")
		<-ctx.Done()
		return nil
	}

	slog.Info("starting consumer", slog.String("backend", config.Get().Broker.Backend))
	return w.broker.Start(ctx, w.routes)
}

func newBroker(i do.Injector) (IBroker, error) {
	backend := config.Get().Broker.Backend

	switch backend {
	case "":
		return nil, nil
	case "kafka":
		return do.MustInvoke[*brokerKafka.Broker](i), nil
	case "rabbitmq":
		return do.MustInvoke[*brokerRabbitmq.Broker](i), nil
	default:
		return nil, fmt.Errorf("unsupported broker backend %q, must be one of: kafka, rabbitmq", backend)
	}
}

func addWorkers[T river.JobArgs](workers *river.Workers, jobs ...river.Worker[T]) error {
	for _, job := range jobs {
		err := river.AddWorkerSafely(workers, job)
		if err != nil {
			return err
		}
	}

	return nil
}

func addRoutes(routes *[]consumer.Route, newRoutes ...consumer.Route) error {
	for _, route := range newRoutes {
		if route.Consumer == nil {
			return fmt.Errorf("topic %q has no consumer", route.Topic)
		}

		for _, existing := range *routes {
			if existing.Topic == route.Topic {
				return fmt.Errorf("duplicate topic %q", route.Topic)
			}
		}

		*routes = append(*routes, route)
	}

	return nil
}
