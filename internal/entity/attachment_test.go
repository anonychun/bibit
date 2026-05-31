package entity

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAttachmentFromFile(t *testing.T) {
	t.Run("builds attachment metadata from an open file", func(t *testing.T) {
		extension := ".png"
		filePath := filepath.Join(t.TempDir(), "avatar.png")
		fileContent := []byte("image bytes")
		require.NoError(t, os.WriteFile(filePath, fileContent, 0o600))

		file, err := os.Open(filePath)
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, file.Close())
		})

		attachment, err := NewAttachmentFromFile(file)

		require.NoError(t, err)
		require.NotNil(t, attachment)
		assert.Equal(t, "avatar.png", attachment.FileName)
		assert.Equal(t, int64(len(fileContent)), attachment.ByteSize)
		require.True(t, strings.HasSuffix(attachment.ObjectName, extension))

		objectToken := strings.TrimSuffix(attachment.ObjectName, extension)
		_, err = ulid.ParseStrict(objectToken)
		require.NoError(t, err)
	})
}

func TestNewAttachmentFromFileHeader(t *testing.T) {
	t.Run("builds attachment metadata from a multipart file header", func(t *testing.T) {
		extension := ".pdf"
		fileHeader := &multipart.FileHeader{
			Filename: "report.pdf",
			Size:     128,
		}

		attachment := NewAttachmentFromFileHeader(fileHeader)

		require.NotNil(t, attachment)
		assert.Equal(t, "report.pdf", attachment.FileName)
		assert.Equal(t, int64(128), attachment.ByteSize)
		require.True(t, strings.HasSuffix(attachment.ObjectName, extension))

		objectToken := strings.TrimSuffix(attachment.ObjectName, extension)
		_, err := ulid.ParseStrict(objectToken)
		require.NoError(t, err)
	})
}
