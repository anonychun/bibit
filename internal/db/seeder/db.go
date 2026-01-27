package seeder

import (
	"context"
	"fmt"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/entity"
	"github.com/samber/do/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	do.Provide(bootstrap.Injector, NewDB)
}

type IDB interface {
	Seed(ctx context.Context) error
}

type DB struct {
	gormDB *gorm.DB
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

	return &DB{
		gormDB: gormDB,
	}, nil
}

func (d *DB) Seed(ctx context.Context) error {
	defaultAdmin := &entity.Admin{
		Name:         "Achmad Chun Chun",
		EmailAddress: "anonychun@gmail.com",
	}

	defaultAdminPassword := "didbnyaada"
	err := defaultAdmin.HashPassword(defaultAdminPassword)
	if err != nil {
		return err
	}

	err = d.gormDB.WithContext(ctx).First(defaultAdmin, "email_address = ?", defaultAdmin.EmailAddress).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	err = d.gormDB.WithContext(ctx).Save(defaultAdmin).Error
	if err != nil {
		return err
	}

	return nil
}
