package sql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"runtime"
	"time"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/current"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func init() {
	do.Provide(bootstrap.Injector, NewPostgresDB)
}

type PostgresDB struct {
	pgxPool *pgxpool.Pool
	sqlDB   *sql.DB
	bunDB   *bun.DB
}

var _ IDB = (*PostgresDB)(nil)

func NewPostgresDB(i do.Injector) (*PostgresDB, error) {
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

	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	pgxConfig.MaxConns = int32(maxOpenConns)

	pgxConfig.MaxConnIdleTime = 5 * time.Minute
	pgxConfig.MaxConnLifetime = 30 * time.Minute

	pgxPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}

	sqlDB := stdlib.OpenDBFromPool(pgxPool)
	err = sqlDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	bunDB := bun.NewDB(sqlDB, pgdialect.New())
	bunDB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	return &PostgresDB{
		pgxPool: pgxPool,
		sqlDB:   sqlDB,
		bunDB:   bunDB,
	}, nil
}

func (pd *PostgresDB) DB(ctx context.Context) bun.IDB {
	tx := current.Tx(ctx)
	if tx != nil {
		return tx
	}

	return pd.bunDB
}

func (pd *PostgresDB) SqlDB(ctx context.Context) *sql.DB {
	return pd.sqlDB
}

func (pd *PostgresDB) Shutdown(ctx context.Context) error {
	pd.pgxPool.Close()
	return nil
}
