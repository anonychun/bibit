package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateEntity(t *testing.T) {
	t.Run("creates an entity file", func(t *testing.T) {
		t.Chdir(t.TempDir())
		require.NoError(t, os.MkdirAll(filepath.Join("internal", "entity"), os.ModePerm))

		err := GenerateEntity("invoice")

		require.NoError(t, err)

		fileContent, err := os.ReadFile(filepath.Join("internal", "entity", "invoice.go"))
		require.NoError(t, err)
		assert.Equal(t, "package entity\n", string(fileContent))
	})

	t.Run("returns an error when the entity directory does not exist", func(t *testing.T) {
		t.Chdir(t.TempDir())

		err := GenerateEntity("invoice")

		require.Error(t, err)
	})
}
