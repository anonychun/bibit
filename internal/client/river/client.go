package river

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	dbSql "github.com/anonychun/bibit/internal/db/sql"
	"github.com/anonychun/bibit/internal/logger"
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
	Workers() *river.Workers
}

type Client struct {
	riverClient *river.Client[pgx.Tx]
	workers     *river.Workers
}

var _ IClient = (*Client)(nil)

func NewClient(i do.Injector) (*Client, error) {
	sqlDB := do.MustInvoke[*dbSql.PostgresDB](i)
	l := do.MustInvoke[*logger.Logger](i)

	ctx := context.Background()
	workers := river.NewWorkers()

	riverClient, err := river.NewClient(riverpgxv5.New(sqlDB.PgxPool(ctx)), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: workers,
		Logger:  l.Log(),
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		riverClient: riverClient,
		workers:     workers,
	}, nil
}

func (c *Client) Client() *river.Client[pgx.Tx] {
	return c.riverClient
}

func (c *Client) Workers() *river.Workers {
	return c.workers
}
