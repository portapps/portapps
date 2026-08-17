package netutil

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// DownloadFile downloads url to path, replacing the destination after a successful download.
func DownloadFile(path, url string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	dir := filepath.Dir(path)
	tmpPattern := filepath.Base(path) + ".*.tmp"
	out, err := os.CreateTemp(dir, tmpPattern)
	if err != nil {
		return err
	}
	tmpPath := out.Name()
	defer os.Remove(tmpPath)

	if _, err := io.Copy(out, resp.Body); err != nil {
		_ = out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}
