package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetModuleName(t *testing.T) {
	t.Run("returns a non-empty module name", func(t *testing.T) {
		moduleName := GetModuleName()

		assert.NotEmpty(t, moduleName)
	})
}

func TestExtractPackageName(t *testing.T) {
	t.Run("returns the last path segment", func(t *testing.T) {
		packageName := ExtractPackageName(" api/v1/app/profile ")

		assert.Equal(t, "profile", packageName)
	})
}
