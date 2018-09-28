package portapps

import (
	"fmt"
	"os"
	"time"
)

type RegExportImport struct {
	Key  string
	Arch string
	File string
}

type RegKey struct {
	Key  string
	Arch string
}

// RegAdd add a registry key
func RegAdd(regKey RegKey, force bool) {
	args := []string{"add", regKey.Key, fmt.Sprintf("/reg:%s", regKey.Arch)}
	if force {
		args = append(args, "/f")
	}

	cmdResult, err := ExecCmd(CmdOptions{
		Command:    "reg",
		Args:       args,
		HideWindow: true,
	})
	if err != nil {
		Log.Fatalf("Cannot add registry key '%s': %v", regKey.Key, err)
	}

	if cmdResult.ExitCode != 0 {
		Log.Errorf(fmt.Sprintf("%d", cmdResult.ExitCode))
		if len(cmdResult.Stderr) > 0 {
			Log.Errorf(fmt.Sprintf("%s\n", cmdResult.Stderr))
		}
	}
}

// ExportRegKey export a registry key
func ExportRegKey(reg RegExportImport) {
	cmdResult, err := ExecCmd(CmdOptions{
		Command:    "reg",
		Args:       []string{"export", reg.Key, reg.File, "/y", fmt.Sprintf("/reg:%s", reg.Arch)},
		HideWindow: true,
	})
	if err != nil {
		Log.Fatalf("Cannot export registry key '%s': %v", reg.Key, err)
	}
	if cmdResult.ExitCode != 0 {
		Log.Errorf(fmt.Sprintf("%d", cmdResult.ExitCode))
		if len(cmdResult.Stderr) > 0 {
			Log.Errorf(fmt.Sprintf("%s\n", cmdResult.Stderr))
		}
	}
}

// ImportRegKey import a registry key
func ImportRegKey(reg RegExportImport) {
	// Save current reg key
	ExportRegKey(RegExportImport{
		Key:  reg.Key,
		Arch: reg.Arch,
		File: fmt.Sprintf("%s.%s", reg.File, time.Now().Format("20060102150405")),
	})

	// Check if reg file exists
	if _, err := os.Stat(reg.File); err != nil {
		return
	}

	// Import
	cmdResult, err := ExecCmd(CmdOptions{
		Command:    "reg",
		Args:       []string{"import", reg.File, fmt.Sprintf("/reg:%s", reg.Arch)},
		HideWindow: true,
	})
	if err != nil {
		Log.Fatalf("Cannot import registry file '%s': %v", reg.File, err)
	}
	if cmdResult.ExitCode != 0 {
		Log.Errorf(fmt.Sprintf("%d", cmdResult.ExitCode))
		if len(cmdResult.Stderr) > 0 {
			Log.Errorf(fmt.Sprintf("%s\n", cmdResult.Stderr))
		}
	}
}
