package migrator

import (
	"context"
	"fmt"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/migrations"
	"github.com/pressly/goose/v3"
	"github.com/samber/do/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Database.Sql.Host,
		cfg.Database.Sql.User,
		cfg.Database.Sql.Password,
		cfg.Database.Sql.Name,
		cfg.Database.Sql.Port,
	)

	gormConfig := &gorm.Config{
		QueryFields: true,
		Logger:      logger.Default.LogMode(logger.Info),
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
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
