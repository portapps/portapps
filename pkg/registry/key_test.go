package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPruneRegistryBackupsOnlyRemovesOldBackupsForFile(t *testing.T) {
	dir := t.TempDir()
	regFile := filepath.Join(dir, "settings.reg")

	filesToKeep := []string{
		"settings.reg",
		"settings.reg.not-a-date",
		"notes.txt",
		"other.reg.20260101000000",
	}
	for _, name := range filesToKeep {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte("test"), 0o644))
	}

	for i := range maxBackup + 2 {
		name := fmt.Sprintf("settings.reg.202601010000%02d", i)
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte("backup"), 0o644))
	}

	require.NoError(t, pruneRegistryBackups(regFile))

	requireNoFileExists(t, filepath.Join(dir, "settings.reg.20260101000000"))
	requireNoFileExists(t, filepath.Join(dir, "settings.reg.20260101000001"))
	requireFileExists(t, filepath.Join(dir, "settings.reg.20260101000002"))
	requireFileExists(t, filepath.Join(dir, "settings.reg.20260101000020"))

	for _, name := range filesToKeep {
		requireFileExists(t, filepath.Join(dir, name))
	}
}

func requireFileExists(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	require.NoError(t, err)
}

func requireNoFileExists(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	assert.ErrorIs(t, err, os.ErrNotExist)
}
