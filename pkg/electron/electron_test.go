package electron

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindAppFolderReturnsMatchingDirectory(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(root, "app-1.0.0"), 0o755))

	dir, err := FindAppFolder("app-", root)
	require.NoError(t, err)
	assert.Equal(t, "app-1.0.0", dir)
}

func TestFindAppFolderReturnsReadDirError(t *testing.T) {
	_, err := FindAppFolder("app-", filepath.Join(t.TempDir(), "missing"))
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestFindAppFolderReturnsErrorWithoutMatchingDirectory(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(root, "other-1.0.0"), 0o755))

	_, err := FindAppFolder("app-", root)
	require.Error(t, err)
	assert.EqualError(t, err, `electron main path does not exist with prefix "app-" in `+root)
}
