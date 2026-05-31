package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRepository(t *testing.T) {
	t.Run("creates a repository file for the requested package", func(t *testing.T) {
		t.Chdir(t.TempDir())

		err := GenerateRepository("billing/account")

		require.NoError(t, err)

		fileContent, err := os.ReadFile(filepath.Join("internal", "repository", "billing", "account", "repository.go"))
		require.NoError(t, err)

		content := string(fileContent)
		assert.Contains(t, content, "package account")
		assert.Contains(t, content, "type IRepository interface")
		assert.Contains(t, content, "type Repository struct")
		assert.Contains(t, content, "func NewRepository")
	})

	t.Run("returns directory creation errors", func(t *testing.T) {
		t.Chdir(t.TempDir())
		require.NoError(t, os.MkdirAll("internal", os.ModePerm))
		require.NoError(t, os.WriteFile(filepath.Join("internal", "repository"), []byte("repository"), os.ModePerm))

		err := GenerateRepository("billing/account")

		require.Error(t, err)
	})

	t.Run("returns file generation errors", func(t *testing.T) {
		t.Chdir(t.TempDir())
		require.NoError(t, os.MkdirAll(filepath.Join("internal", "repository", "billing", "account", "repository.go"), os.ModePerm))

		err := GenerateRepository("billing/account")

		require.Error(t, err)
	})
}
