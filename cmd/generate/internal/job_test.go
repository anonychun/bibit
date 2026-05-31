package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateJob(t *testing.T) {
	t.Run("creates a job file for the requested package", func(t *testing.T) {
		t.Chdir(t.TempDir())

		err := GenerateJob("email/send")

		require.NoError(t, err)

		fileContent, err := os.ReadFile(filepath.Join("internal", "job", "email", "send", "job.go"))
		require.NoError(t, err)

		content := string(fileContent)
		assert.Contains(t, content, "package send")
		assert.Contains(t, content, "type Args struct")
		assert.Contains(t, content, "return \"email/send\"")
		assert.Contains(t, content, "type Job struct")
		assert.Contains(t, content, "func NewJob")
		assert.Contains(t, content, "func (j *Job) Work")
	})

	t.Run("returns directory creation errors", func(t *testing.T) {
		t.Chdir(t.TempDir())
		require.NoError(t, os.MkdirAll("internal", os.ModePerm))
		require.NoError(t, os.WriteFile(filepath.Join("internal", "job"), []byte("job"), os.ModePerm))

		err := GenerateJob("email/send")

		require.Error(t, err)
	})

	t.Run("returns file generation errors", func(t *testing.T) {
		t.Chdir(t.TempDir())
		require.NoError(t, os.MkdirAll(filepath.Join("internal", "job", "email", "send", "job.go"), os.ModePerm))

		err := GenerateJob("email/send")

		require.Error(t, err)
	})
}
