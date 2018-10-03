package portapps

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

// Launch to execute the app
func Launch(args []string) {
	Log.Infof("Process: %s", Papp.Process)
	Log.Infof("Args: %s %s", strings.Join(Papp.Args, " "), strings.Join(args, " "))
	Log.Infof("Working dir: %s", Papp.WorkingDir)
	Log.Infof("Data path: %s", Papp.DataPath)

	Log.Infof("Launch %s...", Papp.Name)
	execute := exec.Command(Papp.Process, append(Papp.Args, args...)...)
	execute.Dir = Papp.WorkingDir

	execute.Stdout = Logfile
	execute.Stderr = Logfile

	Log.Infof("Exec %s %s %s", Papp.Process, strings.Join(Papp.Args, " "), strings.Join(args, " "))
	if err := execute.Start(); err != nil {
		Log.Fatalf("Command failed: %v", err)
	}

	execute.Wait()
}

// ExecCmd to execute a system command
func ExecCmd(options CmdOptions) (CmdResult, error) {
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

	Log.Infof("Exec %s %s", options.Command, strings.Join(options.Args, " "))
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

// QuickExecCmd executes a cmd with args with default options
func QuickExecCmd(command string, args []string) error {
	cmdResult, err := ExecCmd(CmdOptions{
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
