package sql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/anonychun/bibit/internal/current"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func TestPostgresDB_DB(t *testing.T) {
	t.Run("returns the bun database by default", func(t *testing.T) {
		rawDB, sqlMock, err := sqlmock.New()
		require.NoError(t, err)

		bunDB := bun.NewDB(rawDB, pgdialect.New())
		t.Cleanup(func() {
			sqlMock.ExpectClose()
			require.NoError(t, bunDB.Close())
		})

		postgresDB := &PostgresDB{bunDB: bunDB}

		db := postgresDB.DB(context.Background())

		actualDB, ok := db.(*bun.DB)
		require.True(t, ok)
		assert.Same(t, bunDB, actualDB)
	})

	t.Run("returns the transaction stored on the context", func(t *testing.T) {
		tx := &bun.Tx{}
		ctx := current.SetTx(context.Background(), tx)
		postgresDB := &PostgresDB{}

		db := postgresDB.DB(ctx)

		actualTx, ok := db.(*bun.Tx)
		require.True(t, ok)
		assert.Same(t, tx, actualTx)
	})
}

func TestPostgresDB_PgxPool(t *testing.T) {
	t.Run("returns the configured pgx pool", func(t *testing.T) {
		pgxPool := &pgxpool.Pool{}
		postgresDB := &PostgresDB{pgxPool: pgxPool}

		actualPool := postgresDB.PgxPool(context.Background())

		assert.Same(t, pgxPool, actualPool)
	})
}
