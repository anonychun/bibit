package manager

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
	CreateDatabase(ctx context.Context) error
	DropDatabase(ctx context.Context) error
}

type DB struct {
	pgxPool *pgxpool.Pool
	bunDB   *bun.DB
	config  *config.Config
}

var _ IDB = (*DB)(nil)

func NewDB(i do.Injector) (*DB, error) {
	ctx := context.Background()
	cfg := do.MustInvoke[*config.Config](i)
	pgxPool, bunDB, err := internal.OpenPostgres(ctx, cfg, "postgres")
	if err != nil {
		return nil, err
	}

	return &DB{
		pgxPool: pgxPool,
		bunDB:   bunDB,
		config:  cfg,
	}, nil
}

func (d *DB) CreateDatabase(ctx context.Context) error {
	var exists bool
	err := d.bunDB.NewRaw("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = ?)", d.config.DB.Sql.Name).Scan(ctx, &exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = d.bunDB.NewRaw("CREATE DATABASE ?", bun.Ident(d.config.DB.Sql.Name)).Exec(ctx)
	return err
}

func (d *DB) DropDatabase(ctx context.Context) error {
	var exists bool
	err := d.bunDB.NewRaw("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = ?)", d.config.DB.Sql.Name).Scan(ctx, &exists)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	_, err = d.bunDB.NewRaw("DROP DATABASE ?", bun.Ident(d.config.DB.Sql.Name)).Exec(ctx)
	return err
}

func (d *DB) Shutdown(ctx context.Context) error {
	d.pgxPool.Close()
	return nil
}
