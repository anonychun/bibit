package user_session

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	dbSql "github.com/anonychun/bibit/internal/db/sql"
	"github.com/anonychun/bibit/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func TestRepository_FindByToken(t *testing.T) {
	t.Run("returns the user session selected by token", func(t *testing.T) {
		ctx := context.Background()
		token := "session-token"
		userSessionID := uuid.New()
		userID := uuid.New()
		bunDB, sqlMock := newMockedBunDB(t)
		sqlDB := dbSql.NewMockIDB(t)
		repository := &Repository{sqlDB: sqlDB}

		sqlDB.EXPECT().DB(ctx).Return(bunDB).Once()
		sqlMock.ExpectQuery(fmt.Sprintf(`SELECT .* FROM "user_sessions" AS "user_session" WHERE \(token = '%s'\) LIMIT 1`, regexp.QuoteMeta(token))).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token", "ip_address", "user_agent"}).
				AddRow(userSessionID.String(), userID.String(), token, "127.0.0.1", "Go test"))

		actualSession, err := repository.FindByToken(ctx, token)

		require.NoError(t, err)
		require.NotNil(t, actualSession)
		assert.Equal(t, userSessionID, actualSession.Id)
		assert.Equal(t, userID, actualSession.UserId)
		assert.Equal(t, token, actualSession.Token)
		assert.Equal(t, "127.0.0.1", actualSession.IpAddress)
		assert.Equal(t, "Go test", actualSession.UserAgent)
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("returns an error when the select fails", func(t *testing.T) {
		ctx := context.Background()
		token := "session-token"
		expectedErr := errors.New("select user session")
		bunDB, sqlMock := newMockedBunDB(t)
		sqlDB := dbSql.NewMockIDB(t)
		repository := &Repository{sqlDB: sqlDB}

		sqlDB.EXPECT().DB(ctx).Return(bunDB).Once()
		sqlMock.ExpectQuery(fmt.Sprintf(`SELECT .* FROM "user_sessions" AS "user_session" WHERE \(token = '%s'\) LIMIT 1`, regexp.QuoteMeta(token))).
			WillReturnError(expectedErr)

		actualSession, err := repository.FindByToken(ctx, token)

		require.ErrorIs(t, err, expectedErr)
		assert.Nil(t, actualSession)
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})
}

func TestRepository_Create(t *testing.T) {
	t.Run("inserts the user session", func(t *testing.T) {
		ctx := context.Background()
		newSession := &entity.UserSession{
			UserId:    uuid.New(),
			Token:     "session-token",
			IpAddress: "127.0.0.1",
			UserAgent: "Go test",
		}
		createdAt := time.Now()
		updatedAt := createdAt.Add(time.Second)
		bunDB, sqlMock := newMockedBunDB(t)
		sqlDB := dbSql.NewMockIDB(t)
		repository := &Repository{sqlDB: sqlDB}

		sqlDB.EXPECT().DB(ctx).Return(bunDB).Once()
		sqlMock.ExpectQuery(fmt.Sprintf(
			`INSERT INTO "user_sessions" .* VALUES \(DEFAULT, DEFAULT, DEFAULT, '%s', '%s', '%s', '%s'\) RETURNING`,
			regexp.QuoteMeta(newSession.UserId.String()),
			regexp.QuoteMeta(newSession.Token),
			regexp.QuoteMeta(newSession.IpAddress),
			regexp.QuoteMeta(newSession.UserAgent),
		)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(uuid.New().String(), createdAt, updatedAt))

		err := repository.Create(ctx, newSession)

		require.NoError(t, err)
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("returns an error when the insert fails", func(t *testing.T) {
		ctx := context.Background()
		newSession := &entity.UserSession{
			UserId:    uuid.New(),
			Token:     "session-token",
			IpAddress: "127.0.0.1",
			UserAgent: "Go test",
		}
		expectedErr := errors.New("insert user session")
		bunDB, sqlMock := newMockedBunDB(t)
		sqlDB := dbSql.NewMockIDB(t)
		repository := &Repository{sqlDB: sqlDB}

		sqlDB.EXPECT().DB(ctx).Return(bunDB).Once()
		sqlMock.ExpectQuery(fmt.Sprintf(
			`INSERT INTO "user_sessions" .* VALUES \(DEFAULT, DEFAULT, DEFAULT, '%s', '%s', '%s', '%s'\) RETURNING`,
			regexp.QuoteMeta(newSession.UserId.String()),
			regexp.QuoteMeta(newSession.Token),
			regexp.QuoteMeta(newSession.IpAddress),
			regexp.QuoteMeta(newSession.UserAgent),
		)).
			WillReturnError(expectedErr)

		err := repository.Create(ctx, newSession)

		require.ErrorIs(t, err, expectedErr)
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})
}

func TestRepository_DeleteByToken(t *testing.T) {
	t.Run("deletes the user session by token", func(t *testing.T) {
		ctx := context.Background()
		token := "session-token"
		bunDB, sqlMock := newMockedBunDB(t)
		sqlDB := dbSql.NewMockIDB(t)
		repository := &Repository{sqlDB: sqlDB}

		sqlDB.EXPECT().DB(ctx).Return(bunDB).Once()
		sqlMock.ExpectExec(fmt.Sprintf(`DELETE FROM "user_sessions" AS "user_session" WHERE \(token = '%s'\)`, regexp.QuoteMeta(token))).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repository.DeleteByToken(ctx, token)

		require.NoError(t, err)
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("returns an error when the delete fails", func(t *testing.T) {
		ctx := context.Background()
		token := "session-token"
		expectedErr := errors.New("delete user session")
		bunDB, sqlMock := newMockedBunDB(t)
		sqlDB := dbSql.NewMockIDB(t)
		repository := &Repository{sqlDB: sqlDB}

		sqlDB.EXPECT().DB(ctx).Return(bunDB).Once()
		sqlMock.ExpectExec(fmt.Sprintf(`DELETE FROM "user_sessions" AS "user_session" WHERE \(token = '%s'\)`, regexp.QuoteMeta(token))).
			WillReturnError(expectedErr)

		err := repository.DeleteByToken(ctx, token)

		require.ErrorIs(t, err, expectedErr)
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})
}

func newMockedBunDB(t *testing.T) (*bun.DB, sqlmock.Sqlmock) {
	t.Helper()

	rawDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	bunDB := bun.NewDB(rawDB, pgdialect.New())
	t.Cleanup(func() {
		sqlMock.ExpectClose()
		_ = bunDB.Close()
	})

	return bunDB, sqlMock
}
