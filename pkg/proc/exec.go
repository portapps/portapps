package proc

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/pkg/errors"
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

	waitErr := command.Wait()
	if waitErr != nil {
		if _, ok := waitErr.(*exec.ExitError); !ok {
			return result, waitErr
		}
	}
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
	return cmdError(cmdResult, err)
}

func cmdError(cmdResult CmdResult, err error) error {
	var message string
	switch {
	case err != nil:
		message = err.Error()
	case cmdResult.ExitCode != 0:
		message = fmt.Sprintf("exit code %d", cmdResult.ExitCode)
	default:
		return nil
	}
	if err != nil && cmdResult.ExitCode != 0 {
		message += fmt.Sprintf(" (exit code %d)", cmdResult.ExitCode)
	}
	if cmdResult.Stderr != "" {
		message += fmt.Sprintf("\n%s\n", cmdResult.Stderr)
	}
	return errors.New(message)
}
