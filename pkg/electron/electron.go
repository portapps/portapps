package electron

import (
	"fmt"
	"os"
	"strings"
)

// FindAppFolder returns the first directory in source matching prefix.
func FindAppFolder(prefix, source string) (string, error) {
	entries, err := os.ReadDir(source)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() && strings.HasPrefix(name, prefix) {
			return name, nil
		}
	}
	return "", fmt.Errorf("electron main path does not exist with prefix %q in %s", prefix, source)
}
