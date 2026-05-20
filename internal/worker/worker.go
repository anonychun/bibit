package worker

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	clientRiver "github.com/anonychun/bibit/internal/client/river"
	"github.com/anonychun/bibit/internal/logger"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewWorker)
}

type IWorker interface {
	Start(ctx context.Context) error
}

type Worker struct {
	riverClient clientRiver.IClient
	logger      logger.ILogger
}

var _ IWorker = (*Worker)(nil)

func NewWorker(i do.Injector) (*Worker, error) {
	return &Worker{
		riverClient: do.MustInvoke[*clientRiver.Client](i),
		logger:      do.MustInvoke[*logger.Logger](i),
	}, nil
}

func (w *Worker) Start(ctx context.Context) error {
	w.logger.Log().Info("starting worker")
	err := w.riverClient.Client().Start(ctx)
	if err != nil {
		return err
	}

	<-w.riverClient.Client().Stopped()
	return nil
}

func (w *Worker) Shutdown(ctx context.Context) error {
	w.logger.Log().Info("shutting down worker")
	return w.riverClient.Client().Stop(ctx)
}
