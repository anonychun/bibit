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

func Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	sqlDB, err := do.Invoke[*dbSql.PostgresDB](bootstrap.Injector)
	if err != nil {
		return err
	}

	return sqlDB.DB(ctx).RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		ctx = current.SetTx(ctx, &tx)
		return fn(ctx)
	})
}
