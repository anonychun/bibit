package sql

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type IDB interface {
	DB(ctx context.Context) bun.IDB
	SqlDB(ctx context.Context) *sql.DB
}
