package entity

import (
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserSession_GenerateToken(t *testing.T) {
	t.Run("stores a valid ULID token", func(t *testing.T) {
		userSession := &UserSession{}

		userSession.GenerateToken()

		assert.NotEmpty(t, userSession.Token)

		_, err := ulid.ParseStrict(userSession.Token)
		require.NoError(t, err)
	})
}
