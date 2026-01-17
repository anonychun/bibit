package scheduler

import (
	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewScheduler)
}

type Scheduler interface {
}

type SchedulerImpl struct {
}

func NewScheduler(i do.Injector) (Scheduler, error) {
	return &SchedulerImpl{}, nil
}
