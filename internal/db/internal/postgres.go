package internal

import (
	"context"
	"fmt"
	"net/url"
	"runtime"
	"time"

	"github.com/anonychun/bibit/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func OpenPostgres(ctx context.Context, cfg *config.Config, dbName string) (*pgxpool.Pool, *bun.DB, error) {
	dsn := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.DB.Sql.User, cfg.DB.Sql.Password),
		Host:     fmt.Sprintf("%s:%d", cfg.DB.Sql.Host, cfg.DB.Sql.Port),
		Path:     dbName,
		RawQuery: "sslmode=disable",
	}

	pgxConfig, err := pgxpool.ParseConfig(dsn.String())
	if err != nil {
		return nil, nil, err
	}
	pgxConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	pgxConfig.MaxConns = int32(maxOpenConns)

	pgxConfig.MaxConnIdleTime = 5 * time.Minute
	pgxConfig.MaxConnLifetime = 30 * time.Minute

	pgxPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, nil, err
	}

	err = pgxPool.Ping(ctx)
	if err != nil {
		pgxPool.Close()
		return nil, nil, err
	}

	sqlDB := stdlib.OpenDBFromPool(pgxPool)
	bunDB := bun.NewDB(sqlDB, pgdialect.New())
	bunDB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	return pgxPool, bunDB, nil
}
