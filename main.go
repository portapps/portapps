package portapps

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/google/logger"
	_ "github.com/josephspurrier/goversioninfo"
	"github.com/mitchellh/mapstructure"
	"github.com/portapps/portapps/pkg/dialog"
	"gopkg.in/yaml.v2"
)

type papp struct {
	ID         string
	Name       string
	Path       string
	AppPath    string
	DataPath   string
	WorkingDir string
	Process    string
	Args       []string

	config *Config
}

var (
	// Papp settings
	Papp papp

	// Log is the logger used by portapps
	Log     *logger.Logger
	logfile *os.File
)

// Init must be used by every Portapp
func Init() {
	InitWithCfg(nil)
}

// Init with app configuration
func InitWithCfg(appcfg interface{}) {
	var err error

	Papp.Path, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		LogFatal("Current path:", err)
	}

	Papp.DataPath = AppPathJoin("data")

	logfile, err = os.OpenFile(PathJoin(Papp.Path, Papp.ID+".log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		LogFatal("Cannot open log file:", err)
	}

	// Startup
	Log = logger.Init(Papp.Name, false, false, logfile)
	Log.Info("--------")
	Log.Infof("Starting %s...", Papp.Name)
	Log.Infof("Current path: %s", Papp.Path)

	// Configuration
	Log.Info("Loading configuration...")
	if err = loadConfig(appcfg); err != nil {
		LogFatal("Cannot load configuration:", err)
	}
	if appcfg != reflect.Zero(reflect.TypeOf(appcfg)).Interface() {
		if err := mapstructure.Decode(Papp.config.App, appcfg); err != nil {
			LogFatal("Cannot decode app configuration:", err)
		}
	}
	b, _ := yaml.Marshal(Papp.config)
	Log.Infof("Configuration:\n%s", string(b))
}

// Launch to execute the app
func Launch(args []string) {
	Log.Infof("Process: %s", Papp.Process)
	Log.Infof("Args (config file): %s", strings.Join(Papp.config.Common.Args, " "))
	Log.Infof("Args (cmd line): %s", strings.Join(args, " "))
	Log.Infof("Args (hardcoded): %s", strings.Join(Papp.Args, " "))
	Log.Infof("Working dir: %s", Papp.WorkingDir)
	Log.Infof("Data path: %s", Papp.DataPath)

	if !Exists(Papp.Process) {
		Log.Errorf("Application not found in %s", Papp.Process)
		if _, err := dialog.MsgBox(
			fmt.Sprintf("%s portable", Papp.Name),
			fmt.Sprintf("%s application cannot be found in %s", Papp.Name, Papp.Process),
			dialog.MsgBoxBtnOk|dialog.MsgBoxIconError); err != nil {
			Log.Error("Cannot create dialog box", err)
		}
		return
	}

	Log.Infof("Launch %s...", Papp.Name)
	jArgs := append(append(Papp.config.Common.Args, args...), Papp.Args...)
	execute := exec.Command(Papp.Process, jArgs...)
	execute.Dir = Papp.WorkingDir

	execute.Stdout = logfile
	execute.Stderr = logfile

	Log.Infof("Exec %s %s", Papp.Process, strings.Join(jArgs, " "))
	if err := execute.Start(); err != nil {
		LogFatal("Command failed: %v", err)
	}

	execute.Wait()
}

func LogFatal(v ...interface{}) {
	Log.Error(v...)
	msgBoxFatal()
}

func LogFatalf(format string, v ...interface{}) {
	Log.Errorf(format, v...)
	msgBoxFatal()
}

func msgBoxFatal() {
	if _, err := dialog.MsgBox(
		fmt.Sprintf("%s portable", Papp.Name),
		fmt.Sprintf("An error has occurred, see your %s.log file for more info", Papp.ID),
		dialog.MsgBoxBtnOk|dialog.MsgBoxIconError); err != nil {
		Log.Error("Cannot create dialog box", err)
	}
	os.Exit(1)
}
