package entity

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBase_BeforeUpdate(t *testing.T) {
	t.Run("refreshes UpdatedAt to the current time", func(t *testing.T) {
		originalUpdatedAt := time.Now().Add(-time.Hour)
		base := &Base{UpdatedAt: originalUpdatedAt}

		startedAt := time.Now()
		err := base.BeforeUpdate(context.Background(), nil)
		finishedAt := time.Now()

		require.NoError(t, err)
		assert.True(t, base.UpdatedAt.After(originalUpdatedAt))
		assert.False(t, base.UpdatedAt.Before(startedAt))
		assert.False(t, base.UpdatedAt.After(finishedAt))
	})
}
