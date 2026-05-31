package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateUsecase(t *testing.T) {
	t.Run("creates usecase, handler, and dto files for the requested package", func(t *testing.T) {
		t.Chdir(t.TempDir())

		err := GenerateUsecase("api/v1/app/profile")

		require.NoError(t, err)

		baseDir := filepath.Join("internal", "usecase", "api", "v1", "app", "profile")
		usecaseContent, err := os.ReadFile(filepath.Join(baseDir, "usecase.go"))
		require.NoError(t, err)
		handlerContent, err := os.ReadFile(filepath.Join(baseDir, "handler.go"))
		require.NoError(t, err)
		dtoContent, err := os.ReadFile(filepath.Join(baseDir, "dto.go"))
		require.NoError(t, err)

		assert.Contains(t, string(usecaseContent), "package profile")
		assert.Contains(t, string(usecaseContent), "type IUsecase interface")
		assert.Contains(t, string(usecaseContent), "func NewUsecase")
		assert.Contains(t, string(handlerContent), "package profile")
		assert.Contains(t, string(handlerContent), "type IHandler interface")
		assert.Contains(t, string(handlerContent), "func NewHandler")
		assert.Equal(t, "package profile\n", string(dtoContent))
	})

	t.Run("returns directory creation errors", func(t *testing.T) {
		t.Chdir(t.TempDir())
		require.NoError(t, os.MkdirAll("internal", os.ModePerm))
		require.NoError(t, os.WriteFile(filepath.Join("internal", "usecase"), []byte("usecase"), os.ModePerm))

		err := GenerateUsecase("api/v1/app/profile")

		require.Error(t, err)
	})

	t.Run("returns usecase file generation errors", func(t *testing.T) {
		t.Chdir(t.TempDir())
		require.NoError(t, os.MkdirAll(filepath.Join("internal", "usecase", "api", "v1", "app", "profile", "usecase.go"), os.ModePerm))

		err := GenerateUsecase("api/v1/app/profile")

		require.Error(t, err)
	})

	t.Run("returns handler file generation errors", func(t *testing.T) {
		t.Chdir(t.TempDir())
		baseDir := filepath.Join("internal", "usecase", "api", "v1", "app", "profile")
		require.NoError(t, os.MkdirAll(filepath.Join(baseDir, "handler.go"), os.ModePerm))

		err := GenerateUsecase("api/v1/app/profile")

		require.Error(t, err)
	})

	t.Run("returns dto file generation errors", func(t *testing.T) {
		t.Chdir(t.TempDir())
		baseDir := filepath.Join("internal", "usecase", "api", "v1", "app", "profile")
		require.NoError(t, os.MkdirAll(filepath.Join(baseDir, "dto.go"), os.ModePerm))

		err := GenerateUsecase("api/v1/app/profile")

		require.Error(t, err)
	})
}
