package user

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

func TestRepository_FindById(t *testing.T) {
	t.Run("returns the user selected by id", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()
		bunDB, sqlMock := newMockedBunDB(t)
		sqlDB := dbSql.NewMockIDB(t)
		repository := &Repository{sqlDB: sqlDB}

		sqlDB.EXPECT().DB(ctx).Return(bunDB).Once()
		sqlMock.ExpectQuery(fmt.Sprintf(`SELECT .* FROM "users" AS "user" WHERE \(id = '%s'\) LIMIT 1`, regexp.QuoteMeta(userID.String()))).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email_address", "password_digest"}).
				AddRow(userID.String(), "Ada Lovelace", "ada@example.com", "password-digest"))

		actualUser, err := repository.FindById(ctx, userID.String())

		require.NoError(t, err)
		require.NotNil(t, actualUser)
		assert.Equal(t, userID, actualUser.Id)
		assert.Equal(t, "Ada Lovelace", actualUser.Name)
		assert.Equal(t, "ada@example.com", actualUser.EmailAddress)
		assert.Equal(t, "password-digest", actualUser.PasswordDigest)
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("returns an error when the select fails", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New().String()
		expectedErr := errors.New("select user by id")
		bunDB, sqlMock := newMockedBunDB(t)
		sqlDB := dbSql.NewMockIDB(t)
		repository := &Repository{sqlDB: sqlDB}

		sqlDB.EXPECT().DB(ctx).Return(bunDB).Once()
		sqlMock.ExpectQuery(fmt.Sprintf(`SELECT .* FROM "users" AS "user" WHERE \(id = '%s'\) LIMIT 1`, regexp.QuoteMeta(userID))).
			WillReturnError(expectedErr)

		actualUser, err := repository.FindById(ctx, userID)

		require.ErrorIs(t, err, expectedErr)
		assert.Nil(t, actualUser)
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})
}

func TestRepository_FindByEmailAddress(t *testing.T) {
	t.Run("returns the user selected by email address", func(t *testing.T) {
		ctx := context.Background()
		emailAddress := "ada@example.com"
		userID := uuid.New()
		bunDB, sqlMock := newMockedBunDB(t)
		sqlDB := dbSql.NewMockIDB(t)
		repository := &Repository{sqlDB: sqlDB}

		sqlDB.EXPECT().DB(ctx).Return(bunDB).Once()
		sqlMock.ExpectQuery(fmt.Sprintf(`SELECT .* FROM "users" AS "user" WHERE \(email_address = '%s'\) LIMIT 1`, regexp.QuoteMeta(emailAddress))).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email_address", "password_digest"}).
				AddRow(userID.String(), "Ada Lovelace", emailAddress, "password-digest"))

		actualUser, err := repository.FindByEmailAddress(ctx, emailAddress)

		require.NoError(t, err)
		require.NotNil(t, actualUser)
		assert.Equal(t, userID, actualUser.Id)
		assert.Equal(t, "Ada Lovelace", actualUser.Name)
		assert.Equal(t, emailAddress, actualUser.EmailAddress)
		assert.Equal(t, "password-digest", actualUser.PasswordDigest)
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})
}

func TestRepository_Create(t *testing.T) {
	t.Run("inserts the user", func(t *testing.T) {
		ctx := context.Background()
		newUser := &entity.User{
			Name:           "Ada Lovelace",
			EmailAddress:   "ada@example.com",
			PasswordDigest: "password-digest",
		}
		createdAt := time.Now()
		updatedAt := createdAt.Add(time.Second)
		bunDB, sqlMock := newMockedBunDB(t)
		sqlDB := dbSql.NewMockIDB(t)
		repository := &Repository{sqlDB: sqlDB}

		sqlDB.EXPECT().DB(ctx).Return(bunDB).Once()
		sqlMock.ExpectQuery(fmt.Sprintf(
			`INSERT INTO "users" .* VALUES \(DEFAULT, DEFAULT, DEFAULT, '%s', '%s', '%s'\) RETURNING`,
			regexp.QuoteMeta(newUser.Name),
			regexp.QuoteMeta(newUser.EmailAddress),
			regexp.QuoteMeta(newUser.PasswordDigest),
		)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(uuid.New().String(), createdAt, updatedAt))

		err := repository.Create(ctx, newUser)

		require.NoError(t, err)
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("returns an error when the insert fails", func(t *testing.T) {
		ctx := context.Background()
		newUser := &entity.User{
			Name:           "Ada Lovelace",
			EmailAddress:   "ada@example.com",
			PasswordDigest: "password-digest",
		}
		expectedErr := errors.New("insert user")
		bunDB, sqlMock := newMockedBunDB(t)
		sqlDB := dbSql.NewMockIDB(t)
		repository := &Repository{sqlDB: sqlDB}

		sqlDB.EXPECT().DB(ctx).Return(bunDB).Once()
		sqlMock.ExpectQuery(fmt.Sprintf(
			`INSERT INTO "users" .* VALUES \(DEFAULT, DEFAULT, DEFAULT, '%s', '%s', '%s'\) RETURNING`,
			regexp.QuoteMeta(newUser.Name),
			regexp.QuoteMeta(newUser.EmailAddress),
			regexp.QuoteMeta(newUser.PasswordDigest),
		)).
			WillReturnError(expectedErr)

		err := repository.Create(ctx, newUser)

		require.ErrorIs(t, err, expectedErr)
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})
}

func TestRepository_ExistsByEmailAddress(t *testing.T) {
	t.Run("returns true when a matching user exists", func(t *testing.T) {
		ctx := context.Background()
		emailAddress := "ada@example.com"
		bunDB, sqlMock := newMockedBunDB(t)
		sqlDB := dbSql.NewMockIDB(t)
		repository := &Repository{sqlDB: sqlDB}

		sqlDB.EXPECT().DB(ctx).Return(bunDB).Once()
		sqlMock.ExpectQuery(fmt.Sprintf(`SELECT EXISTS \(SELECT .* FROM "users" AS "user" WHERE \(email_address = '%s'\)\)`, regexp.QuoteMeta(emailAddress))).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		exists, err := repository.ExistsByEmailAddress(ctx, emailAddress)

		require.NoError(t, err)
		assert.True(t, exists)
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("returns an error when the existence check fails", func(t *testing.T) {
		ctx := context.Background()
		emailAddress := "ada@example.com"
		expectedErr := errors.New("check user exists")
		bunDB, sqlMock := newMockedBunDB(t)
		sqlDB := dbSql.NewMockIDB(t)
		repository := &Repository{sqlDB: sqlDB}

		sqlDB.EXPECT().DB(ctx).Return(bunDB).Once()
		sqlMock.ExpectQuery(fmt.Sprintf(`SELECT EXISTS \(SELECT .* FROM "users" AS "user" WHERE \(email_address = '%s'\)\)`, regexp.QuoteMeta(emailAddress))).
			WillReturnError(expectedErr)

		exists, err := repository.ExistsByEmailAddress(ctx, emailAddress)

		require.ErrorIs(t, err, expectedErr)
		assert.False(t, exists)
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
