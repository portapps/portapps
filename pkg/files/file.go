package files

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/sys/windows"
)

// SetFileAttributes set attributes to a file
func SetFileAttributes(path string, attrs uint32) error {
	pointer, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return windows.SetFileAttributes(pointer, attrs)
}

// CopyFile copies a file.
func CopyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	return err
}

// CopyFolder copies a folder.
func CopyFolder(source string, dest string) error {
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dest, relPath)
		if entry.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}
		return CopyFile(path, destPath)
	})
}

// Exists reports whether the named file or directory exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
	}
	return true
}

// ReplaceByPrefix replaces lines starting with a specific prefix.
func ReplaceByPrefix(filename, prefix, replace string) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	lines := strings.Split(string(input), "\n")
	for i := range lines {
		if strings.HasPrefix(lines[i], prefix) {
			lines[i] = replace
		}
	}
	return os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0o644)
}

// IsDirEmpty determines if directory is empty
func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if errors.Is(err, io.EOF) {
		return true, nil
	}
	return false, err
}

// StartMenuPath returns the user start menu path
func StartMenuPath() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, "Microsoft", "Windows", "Start Menu", "Programs")
}

// Cleanup removes leftover folders.
func Cleanup(folders ...string) {
	for _, folder := range folders {
		if err := os.RemoveAll(folder); err != nil {
			log.Error().Err(err).Msgf("Cannot cleanup %s", folder)
		}
	}
}
