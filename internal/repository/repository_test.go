package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"unsafe"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/current"
	dbSql "github.com/anonychun/bibit/internal/db/sql"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func TestTransaction(t *testing.T) {
	t.Run("runs the callback with a transaction on the context and commits", func(t *testing.T) {
		ctx := context.Background()
		sqlMock := registerTransactionDB(t)
		called := false

		sqlMock.ExpectBegin()
		sqlMock.ExpectCommit()

		err := Transaction(ctx, func(ctx context.Context) error {
			called = true
			assert.NotNil(t, current.Tx(ctx))
			return nil
		})

		require.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("rolls back and returns the callback error", func(t *testing.T) {
		ctx := context.Background()
		expectedErr := errors.New("transaction callback")
		sqlMock := registerTransactionDB(t)

		sqlMock.ExpectBegin()
		sqlMock.ExpectRollback()

		err := Transaction(ctx, func(ctx context.Context) error {
			assert.NotNil(t, current.Tx(ctx))
			return expectedErr
		})

		require.ErrorIs(t, err, expectedErr)
	})
}

func registerTransactionDB(t *testing.T) sqlmock.Sqlmock {
	t.Helper()

	rawDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	bunDB := bun.NewDB(rawDB, pgdialect.New())
	postgresDB := &dbSql.PostgresDB{}
	bunDBField := reflect.ValueOf(postgresDB).Elem().FieldByName("bunDB")
	reflect.NewAt(bunDBField.Type(), unsafe.Pointer(bunDBField.UnsafeAddr())).Elem().Set(reflect.ValueOf(bunDB))

	do.OverrideValue[*dbSql.PostgresDB](bootstrap.Injector, postgresDB)
	t.Cleanup(func() {
		sqlMock.ExpectClose()
		require.NoError(t, bunDB.Close())
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})

	return sqlMock
}
