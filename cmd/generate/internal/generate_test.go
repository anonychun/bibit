package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFile(t *testing.T) {
	t.Run("renders a template into the target file", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "generated.go")

		err := generateFile(filePath, "package {{.PackageName}}\n", TemplateData{PackageName: "account"})

		require.NoError(t, err)

		fileContent, err := os.ReadFile(filePath)
		require.NoError(t, err)
		assert.Equal(t, "package account\n", string(fileContent))
	})

	t.Run("returns template parse errors", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "generated.go")

		err := generateFile(filePath, "package {{.PackageName", TemplateData{PackageName: "account"})

		require.Error(t, err)
	})

	t.Run("returns file creation errors", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "missing", "generated.go")

		err := generateFile(filePath, "package {{.PackageName}}\n", TemplateData{PackageName: "account"})

		require.Error(t, err)
	})

	t.Run("returns template execution errors", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "generated.go")

		err := generateFile(filePath, "{{call .PackageName}}", TemplateData{PackageName: "account"})

		require.Error(t, err)
	})
}

func TestExtractPackageName(t *testing.T) {
	t.Run("returns the last path segment", func(t *testing.T) {
		packageName := extractPackageName(" api/v1/app/profile ")

		assert.Equal(t, "profile", packageName)
	})
}

func TestGetModuleName(t *testing.T) {
	t.Run("returns a non-empty module name", func(t *testing.T) {
		moduleName := getModuleName()

		assert.NotEmpty(t, moduleName)
	})
}
