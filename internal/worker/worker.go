package worker

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	clientRiver "github.com/anonychun/bibit/internal/client/river"
	jobHello "github.com/anonychun/bibit/internal/job/hello"
	"github.com/anonychun/bibit/internal/observability"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewWorker)
}

type IWorker interface {
	Start(ctx context.Context) error
}

type Worker struct {
	riverClient   clientRiver.IClient
	observability observability.IObservability
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

	return &Worker{
		riverClient:   riverClient,
		observability: do.MustInvoke[*observability.Observability](i),
	}, nil
}

func (w *Worker) Start(ctx context.Context) error {
	w.observability.Logger().Info("starting worker")
	err := w.riverClient.Client().Start(ctx)
	if err != nil {
		return err
	}

	<-w.riverClient.Client().Stopped()
	return nil
}

func (w *Worker) Shutdown(ctx context.Context) error {
	w.observability.Logger().Info("shutting down worker")
	return w.riverClient.Client().Stop(ctx)
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
