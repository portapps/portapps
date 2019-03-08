package registry

import (
	"fmt"
	"os"
	"time"

	"github.com/portapps/portapps/pkg/proc"
)

// ExportImport the registry export/import structure
type ExportImport struct {
	Key  string
	Arch string
	File string
}

// Key the registry key structure
type Key struct {
	Key  string
	Arch string
}

// Add add a registry key
func Add(key Key, force bool) error {
	args := []string{"add", key.Key, fmt.Sprintf("/reg:%s", key.Arch)}
	if force {
		args = append(args, "/f")
	}

	cmdResult, err := proc.Cmd(proc.CmdOptions{
		Command:    "reg",
		Args:       args,
		HideWindow: true,
	})
	if err != nil {
		return fmt.Errorf("cannot add registry key '%s': %v", key.Key, err)
	}

	if cmdResult.ExitCode != 0 {
		if len(cmdResult.Stderr) > 0 {
			return fmt.Errorf("%s, exit code %d", cmdResult.Stderr, cmdResult.ExitCode)
		} else {
			return fmt.Errorf("exit code %d", cmdResult.ExitCode)
		}
	}

	return nil
}

// ExportKey export a registry key
func ExportKey(reg ExportImport) error {
	cmdResult, err := proc.Cmd(proc.CmdOptions{
		Command:    "reg",
		Args:       []string{"export", reg.Key, reg.File, "/y", fmt.Sprintf("/reg:%s", reg.Arch)},
		HideWindow: true,
	})
	if err != nil {
		return fmt.Errorf("cannot export registry key '%s': %v", reg.Key, err)
	}

	if cmdResult.ExitCode != 0 {
		if len(cmdResult.Stderr) > 0 {
			return fmt.Errorf("%s, exit code %d", cmdResult.Stderr, cmdResult.ExitCode)
		}
		return fmt.Errorf("exit code %d", cmdResult.ExitCode)
	}

	return nil
}

// ImportKey import a registry key
func ImportKey(reg ExportImport) error {
	// Save current reg key
	if err := ExportKey(ExportImport{
		Key:  reg.Key,
		Arch: reg.Arch,
		File: fmt.Sprintf("%s.%s", reg.File, time.Now().Format("20060102150405")),
	}); err != nil {
		return err
	}

	// Check if reg file exists
	if _, err := os.Stat(reg.File); err != nil {
		return fmt.Errorf("reg file %s not found", reg.File)
	}

	// Import
	cmdResult, err := proc.Cmd(proc.CmdOptions{
		Command:    "reg",
		Args:       []string{"import", reg.File, fmt.Sprintf("/reg:%s", reg.Arch)},
		HideWindow: true,
	})
	if err != nil {
		return fmt.Errorf("cannot import registry key '%s': %v", reg.Key, err)
	}

	if cmdResult.ExitCode != 0 {
		if len(cmdResult.Stderr) > 0 {
			return fmt.Errorf("%s, exit code %d", cmdResult.Stderr, cmdResult.ExitCode)
		}
		return fmt.Errorf("exit code %d", cmdResult.ExitCode)
	}

	return nil
}
