package seeder

import (
	"context"
	"fmt"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func init() {
	do.Provide(bootstrap.Injector, NewDB)
}

type IDB interface {
	Seed(ctx context.Context) error
}

type DB struct {
	bunDB *bun.DB
}

var _ IDB = (*DB)(nil)

func NewDB(i do.Injector) (*DB, error) {
	ctx := context.Background()
	cfg := do.MustInvoke[*config.Config](i)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DB.Sql.User,
		cfg.DB.Sql.Password,
		cfg.DB.Sql.Host,
		cfg.DB.Sql.Port,
		cfg.DB.Sql.Name,
	)

	pgxConfig, err := pgxpool.ParseConfig(dsn)
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

	bunDB := bun.NewDB(sqlDB, pgdialect.New())
	bunDB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	return &DB{
		bunDB: bunDB,
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

	defaultAdminExists, err := d.bunDB.NewSelect().Model(defaultAdmin).Where("email_address = ?", defaultAdmin.EmailAddress).Exists(ctx)
	if err != nil {
		return err
	}

	if !defaultAdminExists {
		_, err = d.bunDB.NewInsert().Model(defaultAdmin).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
