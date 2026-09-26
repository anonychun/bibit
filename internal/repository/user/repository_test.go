package user

import (
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
	t.Run("inserts the user and fills in the generated columns", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
		user := &entity.User{Name: "Ada Lovelace", EmailAddress: uuid.NewString() + "@example.com", PasswordDigest: "password-digest"}

		err := repository.Create(ctx, user)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, user.Id)
		assert.False(t, user.CreatedAt.IsZero())
		assert.False(t, user.UpdatedAt.IsZero())
	})

	t.Run("rejects a second user with the same email address", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
		emailAddress := uuid.NewString() + "@example.com"
		require.NoError(t, repository.Create(ctx, &entity.User{Name: "Ada", EmailAddress: emailAddress, PasswordDigest: "x"}))

		err := repository.Create(ctx, &entity.User{Name: "Other Ada", EmailAddress: emailAddress, PasswordDigest: "y"})

		require.Error(t, err)
	})
}

func TestRepository_FindById(t *testing.T) {
	t.Run("returns the user with that id", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
		user := &entity.User{Name: "Ada Lovelace", EmailAddress: uuid.NewString() + "@example.com", PasswordDigest: "password-digest"}
		require.NoError(t, repository.Create(ctx, user))

		actualUser, err := repository.FindById(ctx, user.Id)

		require.NoError(t, err)
		assert.Equal(t, user.Id, actualUser.Id)
		assert.Equal(t, "Ada Lovelace", actualUser.Name)
		assert.Equal(t, user.EmailAddress, actualUser.EmailAddress)
		assert.Equal(t, "password-digest", actualUser.PasswordDigest)
	})

	t.Run("returns sql.ErrNoRows when no user has that id", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}

		actualUser, err := repository.FindById(ctx, uuid.New())

		require.ErrorIs(t, err, sql.ErrNoRows)
		assert.Nil(t, actualUser)
	})
}

func TestRepository_FindByEmailAddress(t *testing.T) {
	t.Run("returns the user with that email address", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
		user := &entity.User{Name: "Ada Lovelace", EmailAddress: uuid.NewString() + "@example.com", PasswordDigest: "password-digest"}
		require.NoError(t, repository.Create(ctx, user))

		actualUser, err := repository.FindByEmailAddress(ctx, user.EmailAddress)

		require.NoError(t, err)
		assert.Equal(t, user.Id, actualUser.Id)
	})

	t.Run("returns sql.ErrNoRows when no user has that email address", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}

		actualUser, err := repository.FindByEmailAddress(ctx, uuid.NewString()+"@example.com")

		require.ErrorIs(t, err, sql.ErrNoRows)
		assert.Nil(t, actualUser)
	})
}

func TestRepository_ExistsByEmailAddress(t *testing.T) {
	t.Run("reports whether a user has that email address", func(t *testing.T) {
		ctx := t.Context()
		repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
		emailAddress := uuid.NewString() + "@example.com"
		require.NoError(t, repository.Create(ctx, &entity.User{Name: "Ada", EmailAddress: emailAddress, PasswordDigest: "x"}))

		exists, err := repository.ExistsByEmailAddress(ctx, emailAddress)
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = repository.ExistsByEmailAddress(ctx, uuid.NewString()+"@example.com")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}
