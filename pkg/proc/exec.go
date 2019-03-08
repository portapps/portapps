package proc

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
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
func Cmd(options CmdOptions) (CmdResult, error) {
	result := CmdResult{Options: options}

	command := exec.Command(options.Command, options.Args...)
	commandStdout := &bytes.Buffer{}
	command.Stdout = commandStdout
	commandStderr := &bytes.Buffer{}
	command.Stderr = commandStderr
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: options.HideWindow}

	if options.WorkingDir != "" {
		command.Dir = options.WorkingDir
	}

	if err := command.Start(); err != nil {
		return result, err
	}

	command.Wait()
	waitStatus := command.ProcessState.Sys().(syscall.WaitStatus)

	result.ExitCode = waitStatus.ExitCode
	result.Stdout = strings.TrimSpace(commandStdout.String())
	result.Stderr = strings.TrimSpace(commandStderr.String())

	return result, nil
}

// QuickCmd executes a cmd with args with default options
func QuickCmd(command string, args []string) error {
	cmdResult, err := Cmd(CmdOptions{
		Command:    command,
		Args:       args,
		HideWindow: true,
	})

	if err != nil {
		var errorOutput string
		if cmdResult.ExitCode != 0 {
			errorOutput = fmt.Sprintf(" (exit code %d)", cmdResult.ExitCode)
			if len(cmdResult.Stderr) > 0 {
				errorOutput += fmt.Sprintf("\n%s\n", cmdResult.Stderr)
			}
		}
		return fmt.Errorf("%s%s", err.Error(), errorOutput)
	}

	return nil
}
