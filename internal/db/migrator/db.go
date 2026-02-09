package migrator

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/migrations"
	"github.com/pressly/goose/v3"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun/driver/pgdriver"
)

func init() {
	do.Provide(bootstrap.Injector, NewDB)
}

type IDB interface {
	Migrate(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type DB struct {
	provider *goose.Provider
}

var _ IDB = (*DB)(nil)

func NewDB(i do.Injector) (*DB, error) {
	cfg := do.MustInvoke[*config.Config](i)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DB.Sql.User,
		cfg.DB.Sql.Password,
		cfg.DB.Sql.Host,
		cfg.DB.Sql.Port,
		cfg.DB.Sql.Name,
	)

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	err := sqlDB.Ping()
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
