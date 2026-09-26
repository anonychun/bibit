package sql

import (
	"context"
	"log/slog"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/current"
	"github.com/anonychun/bibit/internal/db/internal"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

func init() {
	do.Provide(bootstrap.Injector, NewPostgresDB)
}

type PostgresDB struct {
	pgxPool *pgxpool.Pool
	bunDB   *bun.DB
}

var _ IDB = (*PostgresDB)(nil)

func NewPostgresDB(i do.Injector) (*PostgresDB, error) {
	ctx := context.Background()
	cfg := do.MustInvoke[*config.Config](i)

	pgxPool, bunDB, err := internal.OpenPostgres(ctx, cfg, cfg.DB.Sql.Name)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{
		bunDB:   bunDB,
		pgxPool: pgxPool,
	}, nil
}

func (pd *PostgresDB) DB(ctx context.Context) bun.IDB {
	tx := current.Tx(ctx)
	if tx != nil {
		return tx
	}

	return pd.bunDB
}

func (pd *PostgresDB) PgxPool(ctx context.Context) *pgxpool.Pool {
	return pd.pgxPool
}

func (pd *PostgresDB) Shutdown(ctx context.Context) error {
	slog.Info("shutting down postgres db")
	pd.pgxPool.Close()
	return nil
}
