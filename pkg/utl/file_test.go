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

func TestCopyFolderOverwritesExistingFiles(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")

	require.NoError(t, os.MkdirAll(filepath.Join(src, "nested"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(dst, "nested"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(src, "nested", "file.txt"), []byte("new"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dst, "nested", "file.txt"), []byte("old"), 0o644))

	require.NoError(t, CopyFolder(src, dst))

	content, err := os.ReadFile(filepath.Join(dst, "nested", "file.txt"))
	require.NoError(t, err)
	assert.Equal(t, "new", string(content))
}

func TestCopyFolderReturnsSourceError(t *testing.T) {
	err := CopyFolder(filepath.Join(t.TempDir(), "missing"), filepath.Join(t.TempDir(), "dst"))
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestCopyFileOverwritesDestination(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src.txt")
	dst := filepath.Join(root, "dst.txt")
	require.NoError(t, os.WriteFile(src, []byte("new"), 0o644))
	require.NoError(t, os.WriteFile(dst, []byte("old"), 0o644))

	require.NoError(t, CopyFile(src, dst))

	content, err := os.ReadFile(dst)
	require.NoError(t, err)
	assert.Equal(t, "new", string(content))
}

func TestExistsReturnsFalseForMissingPath(t *testing.T) {
	assert.False(t, Exists(filepath.Join(t.TempDir(), "missing")))
}

func TestReplaceByPrefixReplacesMatchingLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	require.NoError(t, os.WriteFile(path, []byte("foo=old\nbar=keep\nfoo2=keep"), 0o644))

	require.NoError(t, ReplaceByPrefix(path, "foo=", "foo=new"))

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "foo=new\nbar=keep\nfoo2=keep", string(content))
}

func TestIsDirEmptyReturnsTrueForEmptyDirectory(t *testing.T) {
	empty, err := IsDirEmpty(t.TempDir())

	require.NoError(t, err)
	assert.True(t, empty)
}

func TestIsDirEmptyReturnsFalseForNonEmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "file.txt"), []byte("portapps"), 0o644))

	empty, err := IsDirEmpty(dir)

	require.NoError(t, err)
	assert.False(t, empty)
}

func TestIsDirEmptyReturnsOpenError(t *testing.T) {
	empty, err := IsDirEmpty(filepath.Join(t.TempDir(), "missing"))

	require.Error(t, err)
	assert.False(t, empty)
}
