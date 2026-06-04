package hello

import (
	"context"
	"fmt"

	"github.com/anonychun/bibit/internal/bootstrap"
	repositoryUser "github.com/anonychun/bibit/internal/repository/user"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewJob)
}

type Args struct {
	UserId uuid.UUID
}

func (Args) Kind() string {
	return "hello"
}

type Job struct {
	river.WorkerDefaults[Args]

	userRepository repositoryUser.IRepository
}

func NewJob(i do.Injector) (*Job, error) {
	return &Job{
		userRepository: do.MustInvoke[*repositoryUser.Repository](i),
	}, nil
}

func (j *Job) Work(ctx context.Context, job *river.Job[Args]) error {
	user, err := j.userRepository.FindById(ctx, job.Args.UserId)
	if err != nil {
		return err
	}

	fmt.Printf("Hello %s\n", user.Name)
	return nil
}
