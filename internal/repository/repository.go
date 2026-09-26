package repository

import (
	"context"
	"database/sql"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/current"
	dbSql "github.com/anonychun/bibit/internal/db/sql"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

func init() {
	do.Provide(bootstrap.Injector, NewRepository)
}

type IRepository interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type Repository struct {
	sqlDB dbSql.IDB
}

var _ IRepository = (*Repository)(nil)

func NewRepository(i do.Injector) (*Repository, error) {
	return &Repository{
		sqlDB: do.MustInvoke[*dbSql.PostgresDB](i),
	}, nil
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.sqlDB.DB(ctx).RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		ctx = current.SetTx(ctx, &tx)
		return fn(ctx)
	})
}
