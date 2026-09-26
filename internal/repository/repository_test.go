package repository

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/anonychun/bibit/internal/bootstrap"
	dbSql "github.com/anonychun/bibit/internal/db/sql"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestRepository_Transaction(t *testing.T) {
	t.Run("keeps the writes when fn succeeds", func(t *testing.T) {
		ctx := t.Context()
		sqlDB := do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)
		repository := &Repository{sqlDB: sqlDB}
		table := newScratchTable(t, ctx, sqlDB)

		err := repository.Transaction(ctx, func(ctx context.Context) error {
			return insertRow(ctx, sqlDB, table)
		})

		require.NoError(t, err)
		assert.Equal(t, 1, countRows(t, ctx, sqlDB, table))
	})

	t.Run("discards the writes and returns the error when fn fails", func(t *testing.T) {
		ctx := t.Context()
		sqlDB := do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)
		repository := &Repository{sqlDB: sqlDB}
		table := newScratchTable(t, ctx, sqlDB)
		expectedErr := errors.New("second write failed")

		err := repository.Transaction(ctx, func(ctx context.Context) error {
			require.NoError(t, insertRow(ctx, sqlDB, table))
			return expectedErr
		})

		require.ErrorIs(t, err, expectedErr)
		assert.Equal(t, 0, countRows(t, ctx, sqlDB, table))
	})
}

func newScratchTable(t *testing.T, ctx context.Context, sqlDB dbSql.IDB) bun.Ident {
	t.Helper()

	table := bun.Ident("transaction_test_" + strings.ReplaceAll(uuid.NewString(), "-", ""))
	_, err := sqlDB.DB(ctx).NewRaw("CREATE TABLE ? (id integer)", table).Exec(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		ctx := context.WithoutCancel(ctx)
		_, err := sqlDB.DB(ctx).NewRaw("DROP TABLE IF EXISTS ?", table).Exec(ctx)
		require.NoError(t, err)
	})

	return table
}

func insertRow(ctx context.Context, sqlDB dbSql.IDB, table bun.Ident) error {
	_, err := sqlDB.DB(ctx).NewRaw("INSERT INTO ? (id) VALUES (1)", table).Exec(ctx)
	return err
}

func countRows(t *testing.T, ctx context.Context, sqlDB dbSql.IDB, table bun.Ident) int {
	t.Helper()

	var count int
	err := sqlDB.DB(ctx).NewRaw("SELECT count(*) FROM ?", table).Scan(ctx, &count)
	require.NoError(t, err)

	return count
}
