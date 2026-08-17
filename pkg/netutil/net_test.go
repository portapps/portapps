package netutil

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadFileWritesResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("portapps"))
	}))
	defer server.Close()

	path := filepath.Join(t.TempDir(), "download.txt")
	require.NoError(t, DownloadFile(path, server.URL))

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "portapps", string(content))
}

func TestDownloadFileReplacesExistingFileAfterSuccessfulDownload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("new"))
	}))
	defer server.Close()

	path := filepath.Join(t.TempDir(), "download.txt")
	require.NoError(t, os.WriteFile(path, []byte("old"), 0o644))
	require.NoError(t, DownloadFile(path, server.URL))

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "new", string(content))
}

func TestDownloadFileReturnsHTTPStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer server.Close()

	path := filepath.Join(t.TempDir(), "download.txt")
	err := DownloadFile(path, server.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")

	_, statErr := os.Stat(path)
	assert.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestDownloadFileKeepsExistingFileOnCopyError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		_, _ = w.Write([]byte("partial"))
	}))
	defer server.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "download.txt")
	require.NoError(t, os.WriteFile(path, []byte("old"), 0o644))

	err := DownloadFile(path, server.URL)
	require.Error(t, err)

	content, readErr := os.ReadFile(path)
	require.NoError(t, readErr)
	assert.Equal(t, "old", string(content))

	entries, readDirErr := os.ReadDir(dir)
	require.NoError(t, readDirErr)
	assert.Len(t, entries, 1)
}
