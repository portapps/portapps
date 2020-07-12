package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/portapps/portapps/v2/pkg/proc"
	reg "golang.org/x/sys/windows/registry"
)

// Key the registry key structure
type Key struct {
	Key     string
	Arch    string
	Default string
}

const (
	maxBackup = 19
)

// Add add a registry key
func (k *Key) Add(force bool) error {
	args := []string{"add", k.Key, fmt.Sprintf("/reg:%s", k.Arch)}
	if k.Default != "" {
		args = append(args, "/d", k.Default)
	}
	if force {
		args = append(args, "/f")
	}

	cmdResult, err := proc.Cmd(proc.CmdOptions{
		Command:    "reg",
		Args:       args,
		HideWindow: true,
	})
	if err != nil {
		return fmt.Errorf("cannot add registry key '%s': %v", k.Key, err)
	}

	if cmdResult.ExitCode != 0 {
		if len(cmdResult.Stderr) > 0 {
			return fmt.Errorf("%s, exit code %d", cmdResult.Stderr, cmdResult.ExitCode)
		}
		return fmt.Errorf("exit code %d", cmdResult.ExitCode)
	}

	return nil
}

// Delete removes a registry key
func (k *Key) Delete(force bool) error {
	args := []string{"delete", k.Key, fmt.Sprintf("/reg:%s", k.Arch)}
	if force {
		args = append(args, "/f")
	}

	cmdResult, err := proc.Cmd(proc.CmdOptions{
		Command:    "reg",
		Args:       args,
		HideWindow: true,
	})
	if err != nil {
		return fmt.Errorf("cannot remove registry key '%s': %v", k.Key, err)
	}

	if cmdResult.ExitCode != 0 {
		if len(cmdResult.Stderr) > 0 {
			return fmt.Errorf("%s, exit code %d", cmdResult.Stderr, cmdResult.ExitCode)
		}
		return fmt.Errorf("exit code %d", cmdResult.ExitCode)
	}

	return nil
}

// Exists checks if a registry key exists
func (k *Key) Exists() bool {
	args := []string{"query", k.Key, fmt.Sprintf("/reg:%s", k.Arch)}

	cmdResult, err := proc.Cmd(proc.CmdOptions{
		Command:    "reg",
		Args:       args,
		HideWindow: true,
	})

	return err == nil && cmdResult.ExitCode == 0
}

// Export exports a registry key
func (k *Key) Export(file string) error {
	if !k.Exists() {
		return nil
	}

	cmdResult, err := proc.Cmd(proc.CmdOptions{
		Command:    "reg",
		Args:       []string{"export", k.Key, file, "/y", fmt.Sprintf("/reg:%s", k.Arch)},
		HideWindow: true,
	})
	if err != nil {
		return fmt.Errorf("cannot export registry key '%s': %v", k.Key, err)
	}

	if cmdResult.ExitCode != 0 {
		if len(cmdResult.Stderr) > 0 {
			return fmt.Errorf("%s, exit code %d", cmdResult.Stderr, cmdResult.ExitCode)
		}
		return fmt.Errorf("exit code %d", cmdResult.ExitCode)
	}

	var regFiles []string
	err = filepath.Walk(filepath.Dir(file), func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".reg" {
			return nil
		}
		regFiles = append(regFiles, path)
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "Cannot retrieve files from reg directory")
	}

	sort.Strings(regFiles)
	if len(regFiles) <= maxBackup {
		return nil
	}

	for len(regFiles) > maxBackup {
		regFilePath := regFiles[0]
		if err := os.Remove(regFilePath); err != nil {
			return err
		}
		regFiles = append(regFiles[:0], regFiles[1:]...)
	}

	return nil
}

// Import imports a registry key
func (k *Key) Import(file string) error {
	// Save current reg key
	if err := k.Export(fmt.Sprintf("%s.%s", file, time.Now().Format("20060102150405"))); err != nil {
		return err
	}

	// Check if reg file exists
	if _, err := os.Stat(file); err != nil {
		return fmt.Errorf("reg file %s not found", file)
	}

	// Import
	cmdResult, err := proc.Cmd(proc.CmdOptions{
		Command:    "reg",
		Args:       []string{"import", file, fmt.Sprintf("/reg:%s", k.Arch)},
		HideWindow: true,
	})
	if err != nil {
		return fmt.Errorf("cannot import registry key '%s': %v", k.Key, err)
	}

	if cmdResult.ExitCode != 0 {
		if len(cmdResult.Stderr) > 0 {
			return fmt.Errorf("%s, exit code %d", cmdResult.Stderr, cmdResult.ExitCode)
		}
		return fmt.Errorf("exit code %d", cmdResult.ExitCode)
	}

	return nil
}

// Open opens a registry key
func (k *Key) Open() (reg.Key, error) {
	regSpl := strings.SplitN(k.Key, `\`, 2)

	var regKey reg.Key
	switch regSpl[0] {
	case "HKCR":
		regKey = reg.CLASSES_ROOT
	case "HKCU":
		regKey = reg.CURRENT_USER
	case "HKLM":
		regKey = reg.LOCAL_MACHINE
	case "HKU":
		regKey = reg.USERS
	case "HKCC":
		regKey = reg.CURRENT_CONFIG
	default:
		return reg.NONE, fmt.Errorf("unknown hive %s", regSpl[0])
	}

	return reg.OpenKey(regKey, regSpl[1], reg.ALL_ACCESS)
}
