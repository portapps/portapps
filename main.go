package portableapps

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/google/logger"
	"github.com/pkg/errors"
)

type papp struct {
	ID         string
	Name       string
	Path       string
	AppPath    string
	DataPath   string
	Process    string
	Args       []string
	WorkingDir string
}

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

type RegExportImport struct {
	Key  string
	Arch string
	File string
}

var (
	// Papp settings
	Papp papp

	// Log is the logger used by portapps
	Log *logger.Logger

	// Logfile is the log file used by logger
	Logfile *os.File
)

// Init must be used by every Portapp
func Init() {
	var err error

	Papp.Path, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Log.Fatal("Current path:", err)
	}

	Papp.DataPath = AppPathJoin("data")

	Logfile, err = os.OpenFile(PathJoin(Papp.Path, Papp.ID+".log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		Log.Fatal("Log file:", err)
	}

	Log = logger.Init(Papp.Name, false, false, Logfile)
	Log.Info("--------")
	Log.Infof("Starting %s...", Papp.Name)
	Log.Infof("Current path: %s", Papp.Path)
}

// FindElectronAppFolder retrieved the app electron folder
func FindElectronAppFolder(prefix string, source string) string {
	Log.Infof("Lookup app folder in: %s", source)
	rootFiles, _ := ioutil.ReadDir(source)
	for _, f := range rootFiles {
		if strings.HasPrefix(f.Name(), prefix) && f.IsDir() {
			Log.Infof("Electron app folder found: %s", f.Name())
			return f.Name()
		}
	}

	Log.Fatalf("Electron main path does not exist with prefix '%s' in %s", prefix, source)
	return ""
}

// OverrideEnv to override an env var
func OverrideEnv(key string, value string) {
	if err := os.Setenv(key, value); err != nil {
		Log.Fatalf("Cannot set %s env var: %v", key, err)
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

// CreateFolderCheck to create a folder and get its path and return error
func CreateFolderCheck(path string) (string, error) {
	if err := os.MkdirAll(path, 777); err != nil {
		return "", err
	}
	return path, nil
}

// CreateFolder to create a folder and get its path
func CreateFolder(path string) string {
	Log.Infof("Create folder %s...", path)
	if _, err := CreateFolderCheck(path); err != nil {
		Log.Fatalf("Cannot create folder: %v", err)
	}
	return path
}

// CreateFile creates / overwrites a file with content
func CreateFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err = file.Sync(); err != nil {
		return err
	}
	return nil
}

// PathJoin to join paths
func PathJoin(elem ...string) string {
	for i, e := range elem {
		if e != "" {
			return strings.Join(elem[i:], `\`)
		}
	}
	return ""
}

// AppPathJoin to join paths from Papp.Path
func AppPathJoin(elem ...string) string {
	return PathJoin(append([]string{Papp.Path}, elem...)...)
}

// FormatUnixPath to format a path for unix
func FormatUnixPath(path string) string {
	return strings.Replace(path, `\`, `/`, -1)
}

// FormatWindowsPath to format a path for windows
func FormatWindowsPath(path string) string {
	return strings.Replace(path, `/`, `\`, -1)
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
		return errors.New(fmt.Sprintf("%s%s", err.Error(), errorOutput))
	}

	return nil
}

// SetConsoleTitle sets windows console title
func SetConsoleTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer syscall.FreeLibrary(handle)

	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, err
	}

	rTitle, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return 0, err
	}

	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(rTitle)), 0, 0)
	return int(r), err
}
