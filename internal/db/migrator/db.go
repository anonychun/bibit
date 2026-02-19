package migrator

import (
	"context"
	"fmt"
	"net/url"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewDB)
}

type IDB interface {
	Migrate(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type DB struct {
	pgxPool  *pgxpool.Pool
	provider *goose.Provider
}

var _ IDB = (*DB)(nil)

func NewDB(i do.Injector) (*DB, error) {
	ctx := context.Background()
	cfg := do.MustInvoke[*config.Config](i)
	dsn := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.DB.Sql.User, cfg.DB.Sql.Password),
		Host:     fmt.Sprintf("%s:%d", cfg.DB.Sql.Host, cfg.DB.Sql.Port),
		Path:     cfg.DB.Sql.Name,
		RawQuery: "sslmode=disable",
	}

	pgxConfig, err := pgxpool.ParseConfig(dsn.String())
	if err != nil {
		return nil, err
	}
	pgxConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pgxPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}

	sqlDB := stdlib.OpenDBFromPool(pgxPool)
	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	provider, err := goose.NewProvider(
		"postgres",
		sqlDB,
		migrations.MigrationsFs,
		goose.WithVerbose(true),
	)
	if err != nil {
		return nil, err
	}

	return &DB{
		pgxPool:  pgxPool,
		provider: provider,
	}, nil
}

func (d *DB) Migrate(ctx context.Context) error {
	_, err := d.provider.Up(ctx)
	return err
}

func (d *DB) Rollback(ctx context.Context) error {
	_, err := d.provider.Down(ctx)
	return err
}

func (d *DB) Shutdown(ctx context.Context) error {
	d.pgxPool.Close()
	return nil
}
