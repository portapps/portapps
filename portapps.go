package portapps

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/josephspurrier/goversioninfo"
	"github.com/mitchellh/mapstructure"
	"github.com/portapps/portapps/pkg/dialog"
	"github.com/portapps/portapps/pkg/logging"
	"github.com/portapps/portapps/pkg/utl"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

// App represents an active app object
type App struct {
	ID   string
	Name string
	Args []string

	RootPath   string
	AppPath    string
	DataPath   string
	WorkingDir string
	Process    string

	logfile *os.File
	config  *Config
}

var (
	// Log represents an active zerolog object
	Log *zerolog.Logger
)

// New creates new app instance
func New(id string, name string) (app *App, err error) {
	return NewWithCfg(id, name, nil)
}

// NewWithCfg creates new app instance with an app config
func NewWithCfg(id string, name string, appcfg interface{}) (app *App, err error) {
	// Init
	app = &App{
		ID:   id,
		Name: name,
	}

	// Paths
	app.RootPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		app.FatalBox(err)
	}
	app.AppPath = utl.PathJoin(app.RootPath, "app")
	app.DataPath = utl.PathJoin(app.RootPath, "data")
	app.WorkingDir = app.AppPath

	// Logfile
	logfolder, err := utl.CreateFolder(utl.PathJoin(app.RootPath, "log"))
	if err != nil {
		app.FatalBox(err)
	}
	logpath := utl.PathJoin(logfolder, fmt.Sprintf("%s.log", app.ID))
	app.logfile, err = os.OpenFile(logpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		app.FatalBox(err)
	}

	// Configure logger
	if Log, err = logging.Configure("debug", logpath, LogHook{app}); err != nil {
		app.FatalBox(err)
	}

	// Startup
	Log.Info().Msg("--------")
	Log.Info().Msgf("Starting %s...", app.Name)
	Log.Info().Msgf("Root path: %s", app.RootPath)

	// Configuration
	var b []byte
	Log.Info().Msg("Loading main configuration...")
	if err = app.loadConfig(appcfg); err != nil {
		return nil, err
	}
	if appcfg != nil {
		Log.Info().Msg("Decoding app configuration...")
		if err = mapstructure.Decode(app.config.App, appcfg); err != nil {
			return nil, err
		}
	}
	b, _ = yaml.Marshal(app.config)
	Log.Info().Msgf("Configuration:\n%s", string(b))

	return app, nil
}

// Launch to execute the app with additional args
func (app *App) Launch(args []string) {
	Log.Info().Msgf("Process: %s", app.Process)
	Log.Info().Msgf("Args (config file): %s", strings.Join(app.config.Common.Args, " "))
	Log.Info().Msgf("Args (cmd line): %s", strings.Join(args, " "))
	Log.Info().Msgf("Args (hardcoded): %s", strings.Join(app.Args, " "))
	Log.Info().Msgf("Working dir: %s", app.WorkingDir)
	Log.Info().Msgf("Data path: %s", app.DataPath)

	if !utl.Exists(app.Process) {
		Log.Error().Msgf("Application not found in %s", app.Process)
		if _, err := dialog.MsgBox(
			fmt.Sprintf("%s portable", app.Name),
			fmt.Sprintf("%s application cannot be found in %s", app.Name, app.Process),
			dialog.MsgBoxBtnOk|dialog.MsgBoxIconError); err != nil {
			Log.Error().Err(err).Msgf("Cannot create dialog box")
		}
		return
	}

	Log.Info().Msgf("Launching %s...", app.Name)
	jArgs := append(append(app.config.Common.Args, args...), app.Args...)
	execute := exec.Command(app.Process, jArgs...)
	execute.Dir = app.WorkingDir

	execute.Stdout = app.logfile
	execute.Stderr = app.logfile

	Log.Info().Msgf("Exec %s %s", app.Process, strings.Join(jArgs, " "))
	if err := execute.Start(); err != nil {
		Log.Fatal().Err(err).Msg("Command failed")
	}
	if err := execute.Wait(); err != nil {
		Log.Error().Err(err).Msg("Command failed")
	}
}

// ErrorBox display an error message box
func (app *App) ErrorBox(msg interface{}) {
	_, _ = dialog.MsgBox(
		fmt.Sprintf("%s portable", app.Name),
		fmt.Sprintf("An error has occurred.\n\n%v", msg),
		dialog.MsgBoxBtnOk|dialog.MsgBoxIconError)
}

// FatalBox display an error message box and exit
func (app *App) FatalBox(msg interface{}) {
	app.ErrorBox(msg)
	os.Exit(1)
}
