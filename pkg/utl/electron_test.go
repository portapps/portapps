package utl

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindElectronAppFolderReturnsMatchingDirectory(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(root, "app-1.0.0"), 0o755))

	dir, err := FindElectronAppFolder("app-", root)
	require.NoError(t, err)
	assert.Equal(t, "app-1.0.0", dir)
}

func TestFindElectronAppFolderReturnsReadDirError(t *testing.T) {
	_, err := FindElectronAppFolder("app-", filepath.Join(t.TempDir(), "missing"))
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
}
