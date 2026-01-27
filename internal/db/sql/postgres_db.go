package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/current"
	"github.com/samber/do/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	do.Provide(bootstrap.Injector, NewPostgresDB)
}

type PostgresDB struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

var _ IDB = (*PostgresDB)(nil)

func NewPostgresDB(i do.Injector) (*PostgresDB, error) {
	cfg := do.MustInvoke[*config.Config](i)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DB.Sql.Host,
		cfg.DB.Sql.User,
		cfg.DB.Sql.Password,
		cfg.DB.Sql.Name,
		cfg.DB.Sql.Port,
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

	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresDB{
		gormDB: gormDB,
		sqlDB:  sqlDB,
	}, nil
}

func (pd *PostgresDB) DB(ctx context.Context) *gorm.DB {
	tx := current.Tx(ctx)
	if tx != nil {
		return tx
	}

	return pd.gormDB.WithContext(ctx)
}

func (pd *PostgresDB) SqlDB(ctx context.Context) *sql.DB {
	return pd.sqlDB
}
