package proc

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows"
)

// CmdOptions options of command
type CmdOptions struct {
	Command    string
	Args       []string
	WorkingDir string
	HideWindow bool
}

// CmdResult result of command
type CmdResult struct {
	Options  CmdOptions
	ExitCode uint32
	Stdout   string
	Stderr   string
}

// Cmd to execute a system command
func Cmd(options CmdOptions) (*CmdResult, error) {
	cmOut := &bytes.Buffer{}
	cmErr := &bytes.Buffer{}

	cm := exec.Command(options.Command, options.Args...)
	cm.Stdout = cmOut
	cm.Stderr = cmErr
	cm.SysProcAttr = &windows.SysProcAttr{
		HideWindow: options.HideWindow,
	}

	if options.WorkingDir != "" {
		cm.Dir = options.WorkingDir
	}

	if err := cm.Start(); err != nil {
		return nil, err
	}
	cm.Wait()

	return &CmdResult{
		Options:  options,
		ExitCode: cm.ProcessState.Sys().(windows.WaitStatus).ExitCode,
		Stdout:   strings.TrimSpace(cmOut.String()),
		Stderr:   strings.TrimSpace(cmErr.String()),
	}, nil
}

// QuickCmd executes a cmd with args with default options
func QuickCmd(command string, args []string) error {
	cmdResult, err := Cmd(CmdOptions{
		Command:    command,
		Args:       args,
		HideWindow: true,
	})
	if err != nil {
		return err
	}

	if cmdResult.ExitCode != 0 {
		stderr := fmt.Sprintf(" (exit code %d)", cmdResult.ExitCode)
		if len(cmdResult.Stderr) > 0 {
			stderr += fmt.Sprintf("\n%s\n", cmdResult.Stderr)
		}
		return errors.New(stderr)
	}

	return nil
}
