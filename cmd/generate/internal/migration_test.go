package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateMigration(t *testing.T) {
	t.Run("creates a sql migration by default", func(t *testing.T) {
		t.Chdir(t.TempDir())
		require.NoError(t, os.MkdirAll("migrations", os.ModePerm))

		err := GenerateMigration("create_users", "")

		require.NoError(t, err)

		migrationFiles, err := filepath.Glob(filepath.Join("migrations", "*_create_users.sql"))
		require.NoError(t, err)
		require.Len(t, migrationFiles, 1)

		fileContent, err := os.ReadFile(migrationFiles[0])
		require.NoError(t, err)
		assert.Contains(t, string(fileContent), "-- +goose Up")
		assert.Contains(t, string(fileContent), "-- +goose Down")
	})

	t.Run("creates the requested migration type", func(t *testing.T) {
		t.Chdir(t.TempDir())
		require.NoError(t, os.MkdirAll("migrations", os.ModePerm))

		err := GenerateMigration("backfill_users", "go")

		require.NoError(t, err)

		migrationFiles, err := filepath.Glob(filepath.Join("migrations", "*_backfill_users.go"))
		require.NoError(t, err)
		require.Len(t, migrationFiles, 1)
	})
}
