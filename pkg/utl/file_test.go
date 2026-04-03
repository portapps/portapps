package utl

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyFolderCopiesNestedFiles(t *testing.T) {
	src := filepath.Join(t.TempDir(), "src")
	dst := filepath.Join(t.TempDir(), "dst")

	require.NoError(t, os.MkdirAll(filepath.Join(src, "nested"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(src, "nested", "file.txt"), []byte("portapps"), 0o644))

	require.NoError(t, CopyFolder(src, dst))

	content, err := os.ReadFile(filepath.Join(dst, "nested", "file.txt"))
	require.NoError(t, err)
	assert.Equal(t, "portapps", string(content))
}

func TestCopyFolderReturnsSourceError(t *testing.T) {
	err := CopyFolder(filepath.Join(t.TempDir(), "missing"), filepath.Join(t.TempDir(), "dst"))
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
}
