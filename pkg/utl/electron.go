package utl

import (
	"fmt"
	"os"
	"strings"
)

// FindElectronAppFolder returns the app electron folder
func FindElectronAppFolder(prefix string, source string) (string, error) {
	rootFiles, err := os.ReadDir(source)
	if err != nil {
		return "", err
	}
	for _, f := range rootFiles {
		if strings.HasPrefix(f.Name(), prefix) && f.IsDir() {
			return f.Name(), nil
		}
	}
	return "", fmt.Errorf("Electron main path does not exist with prefix '%s' in %s", prefix, source)
}
