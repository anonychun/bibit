package hello

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
)

type Args struct {
	Name string
}

func (Args) Kind() string {
	return "hello"
}

type Job struct {
	river.WorkerDefaults[Args]
}

func (j *Job) Work(ctx context.Context, job *river.Job[Args]) error {
	fmt.Printf("Hello %s\n", job.Args.Name)
	return nil
}
