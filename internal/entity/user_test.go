package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_HashAndComparePassword(t *testing.T) {
	t.Run("stores a bcrypt digest and accepts the original password", func(t *testing.T) {
		user := &User{}
		password := "correct horse battery staple"

		err := user.HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, user.PasswordDigest)
		assert.NotEqual(t, password, user.PasswordDigest)
		assert.NoError(t, user.ComparePassword(password))
	})

	t.Run("rejects a different password", func(t *testing.T) {
		user := &User{}
		password := "correct horse battery staple"

		require.NoError(t, user.HashPassword(password))

		assert.Error(t, user.ComparePassword("this is not the password"))
	})
}
