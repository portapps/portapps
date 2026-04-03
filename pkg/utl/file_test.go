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

func TestWriteToFileWritesContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")

	require.NoError(t, WriteToFile(path, "portapps"))

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "portapps", string(content))
}

func TestWriteToFileReturnsCreateErrorWithoutPanicking(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "file.txt")

	require.NotPanics(t, func() {
		err := WriteToFile(path, "portapps")
		require.Error(t, err)
	})
}

func TestFormatPathHelpersReplaceAllSeparators(t *testing.T) {
	assert.Equal(t, "a/b/c", FormatUnixPath(`a\b\c`))
	assert.Equal(t, `a\b\c`, FormatWindowsPath("a/b/c"))
}

func TestExistsReturnsFalseForMissingPath(t *testing.T) {
	assert.False(t, Exists(filepath.Join(t.TempDir(), "missing")))
}

func TestReplaceRewritesAllOccurrences(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	require.NoError(t, os.WriteFile(path, []byte("portapps portapps"), 0o644))

	require.NoError(t, Replace(path, "portapps", "done"))

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "done done", string(content))
}
