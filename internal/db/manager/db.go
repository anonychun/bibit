package manager

import (
	"context"
	"fmt"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/samber/do/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	do.Provide(bootstrap.Injector, NewDB)
}

type IDB interface {
	CreateDatabase(ctx context.Context) error
	DropDatabase(ctx context.Context) error
}

type DB struct {
	gormDB *gorm.DB
	config *config.Config
}

var _ IDB = (*DB)(nil)

func NewDB(i do.Injector) (*DB, error) {
	cfg := do.MustInvoke[*config.Config](i)
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=disable",
		cfg.Database.Sql.Host,
		cfg.Database.Sql.User,
		cfg.Database.Sql.Password,
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

	return &DB{
		gormDB: gormDB,
		config: cfg,
	}, nil
}

func (d *DB) CreateDatabase(ctx context.Context) error {
	var exists bool
	err := d.gormDB.WithContext(ctx).Raw("SELECT 1 FROM pg_database WHERE datname = ?", d.config.Database.Sql.Name).Scan(&exists).Error
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return d.gormDB.WithContext(ctx).Exec(fmt.Sprintf("CREATE DATABASE %s", d.config.Database.Sql.Name)).Error
}

func (d *DB) DropDatabase(ctx context.Context) error {
	var exists bool
	err := d.gormDB.WithContext(ctx).Raw("SELECT 1 FROM pg_database WHERE datname = ?", d.config.Database.Sql.Name).Scan(&exists).Error
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	return d.gormDB.WithContext(ctx).Exec(fmt.Sprintf("DROP DATABASE %s", d.config.Database.Sql.Name)).Error
}
