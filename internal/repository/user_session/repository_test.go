package user_session

import (
	"context"
	"database/sql"
	"testing"

	"github.com/anonychun/bibit/internal/bootstrap"
	dbSql "github.com/anonychun/bibit/internal/db/sql"
	"github.com/anonychun/bibit/internal/entity"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_Create(t *testing.T) {
	t.Run("inserts the session and fills in the generated columns", func(t *testing.T) {
		ctx := t.Context()
		sqlDB := do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)
		repository := &Repository{sqlDB: sqlDB}
		userSession := newUserSession(t, ctx, sqlDB)

		err := repository.Create(ctx, userSession)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, userSession.Id)
	})

	t.Run("rejects a session for a user that does not exist", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
		userSession := &entity.UserSession{UserId: uuid.New(), IpAddress: "127.0.0.1", UserAgent: "Go test"}
		userSession.GenerateToken()

		err := repository.Create(ctx, userSession)

		require.Error(t, err)
	})
}

func TestRepository_FindByToken(t *testing.T) {
	t.Run("returns the session with that token", func(t *testing.T) {
		ctx := t.Context()
		sqlDB := do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)
		repository := &Repository{sqlDB: sqlDB}
		userSession := newUserSession(t, ctx, sqlDB)
		require.NoError(t, repository.Create(ctx, userSession))

		actualSession, err := repository.FindByToken(ctx, userSession.Token)

		require.NoError(t, err)
		assert.Equal(t, userSession.Id, actualSession.Id)
		assert.Equal(t, userSession.UserId, actualSession.UserId)
		assert.Equal(t, "127.0.0.1", actualSession.IpAddress)
		assert.Equal(t, "Go test", actualSession.UserAgent)
	})

	t.Run("returns sql.ErrNoRows when no session has that token", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}

		actualSession, err := repository.FindByToken(ctx, "missing-token")

		require.ErrorIs(t, err, sql.ErrNoRows)
		assert.Nil(t, actualSession)
	})
}

func TestRepository_DeleteByToken(t *testing.T) {
	t.Run("deletes only the session with that token", func(t *testing.T) {
		ctx := t.Context()
		sqlDB := do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)
		repository := &Repository{sqlDB: sqlDB}
		deleted := newUserSession(t, ctx, sqlDB)
		kept := &entity.UserSession{UserId: deleted.UserId, IpAddress: "127.0.0.1", UserAgent: "Go test"}
		kept.GenerateToken()
		require.NoError(t, repository.Create(ctx, deleted))
		require.NoError(t, repository.Create(ctx, kept))

		err := repository.DeleteByToken(ctx, deleted.Token)

		require.NoError(t, err)
		_, err = repository.FindByToken(ctx, deleted.Token)
		require.ErrorIs(t, err, sql.ErrNoRows)
		_, err = repository.FindByToken(ctx, kept.Token)
		require.NoError(t, err)
	})

	t.Run("succeeds when no session has that token", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}

		err := repository.DeleteByToken(ctx, "missing-token")

		require.NoError(t, err)
	})
}

func newUserSession(t *testing.T, ctx context.Context, sqlDB dbSql.IDB) *entity.UserSession {
	t.Helper()

	user := &entity.User{Name: "Ada Lovelace", EmailAddress: uuid.NewString() + "@example.com", PasswordDigest: "password-digest"}
	_, err := sqlDB.DB(ctx).NewInsert().Model(user).Exec(ctx)
	require.NoError(t, err)

	userSession := &entity.UserSession{UserId: user.Id, IpAddress: "127.0.0.1", UserAgent: "Go test"}
	userSession.GenerateToken()

	return userSession
}
