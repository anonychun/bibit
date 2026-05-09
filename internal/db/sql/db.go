package sql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/uptrace/bun"
)

type IDB interface {
	DB(ctx context.Context) bun.IDB

	PgxPool(ctx context.Context) *pgxpool.Pool
}
