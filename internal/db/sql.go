package db

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
	do.Provide(bootstrap.Injector, NewSql)
}

type Sql interface {
	DB(ctx context.Context) *gorm.DB
	SqlDB() *sql.DB
}

type SqlImpl struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

func NewSql(i do.Injector) (Sql, error) {
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

	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return &SqlImpl{
		gormDB: gormDB,
		sqlDB:  sqlDB,
	}, nil
}

func (s *SqlImpl) DB(ctx context.Context) *gorm.DB {
	tx := current.Tx(ctx)
	if tx != nil {
		return tx
	}

	return s.gormDB.WithContext(ctx)
}

func (s *SqlImpl) SqlDB() *sql.DB {
	return s.sqlDB
}

func CreateSqlDatabase() error {
	cfg := do.MustInvoke[*config.Config](bootstrap.Injector)
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=disable",
		cfg.Database.Sql.Host,
		cfg.Database.Sql.User,
		cfg.Database.Sql.Password,
		cfg.Database.Sql.Port,
	)

	gormDB, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}

	var exists bool
	err = gormDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", cfg.Database.Sql.Name).Scan(&exists).Error
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	err = gormDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Database.Sql.Name)).Error
	if err != nil {
		return err
	}

	return nil
}

func DropSqlDatabase() error {
	cfg := do.MustInvoke[*config.Config](bootstrap.Injector)
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=disable",
		cfg.Database.Sql.Host,
		cfg.Database.Sql.User,
		cfg.Database.Sql.Password,
		cfg.Database.Sql.Port,
	)

	gormDB, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}

	var exists bool
	err = gormDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", cfg.Database.Sql.Name).Scan(&exists).Error
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	err = gormDB.Exec(fmt.Sprintf("DROP DATABASE %s", cfg.Database.Sql.Name)).Error
	if err != nil {
		return err
	}

	return nil
}
