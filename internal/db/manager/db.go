package manager

import (
	"context"
	"fmt"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
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
	CreateDatabase(ctx context.Context) error
	DropDatabase(ctx context.Context) error
}

type DB struct {
	bunDB  *bun.DB
	config *config.Config
}

var _ IDB = (*DB)(nil)

func NewDB(i do.Injector) (*DB, error) {
	ctx := context.Background()
	cfg := do.MustInvoke[*config.Config](i)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=disable",
		cfg.DB.Sql.User,
		cfg.DB.Sql.Password,
		cfg.DB.Sql.Host,
		cfg.DB.Sql.Port,
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
		bunDB:  bunDB,
		config: cfg,
	}, nil
}

func (d *DB) CreateDatabase(ctx context.Context) error {
	var exists bool
	err := d.bunDB.NewRaw("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = ?)", d.config.DB.Sql.Name).Scan(ctx, &exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = d.bunDB.NewRaw("CREATE DATABASE ?", bun.Ident(d.config.DB.Sql.Name)).Exec(ctx)
	return err
}

func (d *DB) DropDatabase(ctx context.Context) error {
	var exists bool
	err := d.bunDB.NewRaw("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = ?)", d.config.DB.Sql.Name).Scan(ctx, &exists)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	_, err = d.bunDB.NewRaw("DROP DATABASE ?", bun.Ident(d.config.DB.Sql.Name)).Exec(ctx)
	return err
}
