package sql

import (
	"context"
	"database/sql"
	"fmt"
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
	bunDB *bun.DB
	sqlDB *sql.DB
}

var _ IDB = (*PostgresDB)(nil)

func NewPostgresDB(i do.Injector) (*PostgresDB, error) {
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

	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	pgxConfig.MaxConns = int32(maxOpenConns)

	maxLifeTime := 5 * time.Minute
	pgxConfig.MaxConnIdleTime = maxLifeTime
	pgxConfig.MaxConnLifetime = maxLifeTime

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

	return &PostgresDB{
		bunDB: bunDB,
		sqlDB: sqlDB,
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
