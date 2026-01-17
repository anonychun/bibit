package db

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/migrations"
	"github.com/pressly/goose/v3"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewMigrator)
}

type Migrator interface {
	Migrate(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type MigratorImpl struct {
	sql      Sql
	provider *goose.Provider
}

func NewMigrator(i do.Injector) (Migrator, error) {
	sql := do.MustInvoke[Sql](i)
	provider, err := goose.NewProvider(
		"postgres",
		sql.SqlDB(),
		migrations.MigrationsFs,
		goose.WithVerbose(true),
	)
	if err != nil {
		return nil, err
	}

	return &MigratorImpl{
		sql:      sql,
		provider: provider,
	}, nil
}

func (m *MigratorImpl) Migrate(ctx context.Context) error {
	_, err := m.provider.Up(ctx)
	return err
}

func (m *MigratorImpl) Rollback(ctx context.Context) error {
	_, err := m.provider.Down(ctx)
	return err
}
