package seeder

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/db/internal"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

func init() {
	do.Provide(bootstrap.Injector, NewDB)
}

type IDB interface {
	Seed(ctx context.Context) error
}

type DB struct {
	pgxPool *pgxpool.Pool
	bunDB   *bun.DB
}

var _ IDB = (*DB)(nil)

func NewDB(i do.Injector) (*DB, error) {
	ctx := context.Background()
	cfg := do.MustInvoke[*config.Config](i)
	pgxPool, bunDB, err := internal.OpenPostgres(ctx, cfg, cfg.DB.Sql.Name)
	if err != nil {
		return nil, err
	}

	return &DB{
		pgxPool: pgxPool,
		bunDB:   bunDB,
	}, nil
}

func (d *DB) Seed(ctx context.Context) error {
	return nil
}

func (d *DB) Shutdown(ctx context.Context) error {
	d.pgxPool.Close()
	return nil
}
