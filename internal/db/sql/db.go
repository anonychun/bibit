package sql

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type IDB interface {
	DB(ctx context.Context) *gorm.DB
	SqlDB(ctx context.Context) *sql.DB
}
