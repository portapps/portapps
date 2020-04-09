package registry

import (
	"fmt"
	"os"
	"time"

	"github.com/portapps/portapps/pkg/proc"
)

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

// Delete removes a registry key
func Delete(key Key, force bool) error {
	args := []string{"delete", key.Key, fmt.Sprintf("/reg:%s", key.Arch)}
	if force {
		args = append(args, "/f")
	}

	cmdResult, err := proc.Cmd(proc.CmdOptions{
		Command:    "reg",
		Args:       args,
		HideWindow: true,
	})
	if err != nil {
		return fmt.Errorf("cannot remove registry key '%s': %v", key.Key, err)
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

// Export exports a registry key
func Export(key Key, file string) error {
	cmdResult, err := proc.Cmd(proc.CmdOptions{
		Command:    "reg",
		Args:       []string{"export", key.Key, file, "/y", fmt.Sprintf("/reg:%s", key.Arch)},
		HideWindow: true,
	})
	if err != nil {
		return fmt.Errorf("cannot export registry key '%s': %v", key.Key, err)
	}

	if cmdResult.ExitCode != 0 {
		if len(cmdResult.Stderr) > 0 {
			return fmt.Errorf("%s, exit code %d", cmdResult.Stderr, cmdResult.ExitCode)
		}
		return fmt.Errorf("exit code %d", cmdResult.ExitCode)
	}

	return nil
}

// Import imports a registry key
func Import(key Key, file string) error {
	// Save current reg key
	if err := Export(key, fmt.Sprintf("%s.%s", file, time.Now().Format("20060102150405"))); err != nil {
		return err
	}

	// Check if reg file exists
	if _, err := os.Stat(file); err != nil {
		return fmt.Errorf("reg file %s not found", file)
	}

	// Import
	cmdResult, err := proc.Cmd(proc.CmdOptions{
		Command:    "reg",
		Args:       []string{"import", file, fmt.Sprintf("/reg:%s", key.Arch)},
		HideWindow: true,
	})
	if err != nil {
		return fmt.Errorf("cannot import registry key '%s': %v", key.Key, err)
	}

	if cmdResult.ExitCode != 0 {
		if len(cmdResult.Stderr) > 0 {
			return fmt.Errorf("%s, exit code %d", cmdResult.Stderr, cmdResult.ExitCode)
		}
		return fmt.Errorf("exit code %d", cmdResult.ExitCode)
	}

	return nil
}
