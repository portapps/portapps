package portableapps

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/logger"
)

type app struct {
	ID           string
	Name         string
	Path         string
	RootDataPath string
	MainPath     string
	DataPath     string
	Process      string
	Args         []string
	WorkingDir   string
}

var (
	// App main struct
	App app

	// Log is the logger used by portapps
	Log *logger.Logger

	// Logfile is the log file used by logger
	Logfile *os.File
)

// Init must be used by every Portapp
func Init() {
	var err error

	App.Path, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Log.Fatal("Current path:", err)
	}

	App.MainPath = App.Path
	App.RootDataPath = RootPathJoin("data")
	App.DataPath = App.RootDataPath

	Logfile, err = os.OpenFile(PathJoin(App.Path, App.ID+".log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		Log.Fatal("Log file:", err)
	}

	Log = logger.Init(App.Name, false, false, Logfile)
	Log.Info("--------")
	Log.Infof("Starting %s...", App.Name)
	Log.Infof("Current path: %s", App.Path)
}

// FindElectronMainFolder retrieved the app electron folder based on its prefix
func FindElectronMainFolder(prefix string) string {
	var electronMainPath string
	Log.Infof("Lookup app folder in: %s", App.Path)
	rootFiles, _ := ioutil.ReadDir(App.Path)
	for _, f := range rootFiles {
		if strings.HasPrefix(f.Name(), prefix) && f.IsDir() {
			Log.Infof("Main folder found: %s", f.Name())
			electronMainPath = PathJoin(App.Path, f.Name())
			break
		}
	}
	if _, err := os.Stat(electronMainPath); err != nil {
		Log.Fatalf("Electron main path does not exist with prefix '%s'", prefix)
	}

	Log.Infof("Electron main path found: %s", electronMainPath)
	return electronMainPath
}

// OverrideEnv to override an env var
func OverrideEnv(key string, value string) {
	if err := os.Setenv(key, value); err != nil {
		Log.Fatalf("Cannot set %s env var: %v", key, err)
	}
}

// Launch to execute the app
func Launch() {
	Log.Infof("Process: %s", App.Process)
	Log.Infof("Args: %s", strings.Join(App.Args, " "))
	Log.Infof("Working dir: %s", App.WorkingDir)
	Log.Infof("Data path: %s", App.DataPath)

	Log.Infof("Launch %s...", App.Name)
	execApp := exec.Command(App.Process, App.Args...)
	execApp.Dir = App.WorkingDir

	defer Logfile.Close()
	execApp.Stdout = Logfile
	execApp.Stderr = Logfile

	if err := execApp.Start(); err != nil {
		Log.Fatalf("Cmd Start: %v", err)
	}

	execApp.Wait()
}

// CreateFolder to create a folder and get its path
func CreateFolder(path string) string {
	Log.Infof("Create folder %s...", path)
	if err := os.MkdirAll(path, 777); err != nil {
		Log.Fatalf("Cannot create folder: %v", err)
	}
	return path
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

// RootPathJoin to join paths from App.Path
func RootPathJoin(elem ...string) string {
	return PathJoin(append([]string{App.Path}, elem...)...)
}
