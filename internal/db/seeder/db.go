package seeder

import (
	"context"
	"fmt"
	"net/url"

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
	pgxPool *pgxpool.Pool
	bunDB   *bun.DB
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

	bunDB := bun.NewDB(sqlDB, pgdialect.New())
	bunDB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	return &DB{
		pgxPool: pgxPool,
		bunDB:   bunDB,
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

func (d *DB) Shutdown(ctx context.Context) error {
	d.pgxPool.Close()
	return nil
}
