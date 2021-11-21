package utl

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

// FindElectronAppFolder returns the app electron folder
func FindElectronAppFolder(prefix string, source string) (string, error) {
	rootFiles, _ := os.ReadDir(source)
	for _, f := range rootFiles {
		if strings.HasPrefix(f.Name(), prefix) && f.IsDir() {
			return f.Name(), nil
		}
	}
	return "", errors.Errorf("Electron main path does not exist with prefix '%s' in %s", prefix, source)
}
