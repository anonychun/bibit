package river

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	dbSql "github.com/anonychun/bibit/internal/db/sql"
	jobHello "github.com/anonychun/bibit/internal/job/hello"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewClient)
}

type IClient interface {
	Client() *river.Client[pgx.Tx]
}

type Client struct {
	riverClient *river.Client[pgx.Tx]
}

var _ IClient = (*Client)(nil)

func NewClient(i do.Injector) (*Client, error) {
	workers := river.NewWorkers()
	err := addJobs(workers,
		&jobHello.Job{},
	)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	sqlDB := do.MustInvoke[*dbSql.PostgresDB](i)
	riverClient, err := river.NewClient(riverpgxv5.New(sqlDB.PgxPool(ctx)), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: workers,
	})

	return &Client{
		riverClient: riverClient,
	}, nil
}

func (c *Client) Client() *river.Client[pgx.Tx] {
	return c.riverClient
}
