package portapps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/portapps/portapps/pkg/dialog"
	"github.com/portapps/portapps/pkg/logging"
	"github.com/portapps/portapps/pkg/utl"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

// App represents an active app object
type App struct {
	Info AppInfo

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

// AppInfo represents portapp.json file structure
type AppInfo struct {
	ID              string `json:"id"`
	GUID            string `json:"guid"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	Release         string `json:"release"`
	Date            string `json:"date"`
	Publisher       string `json:"publisher"`
	URL             string `json:"url"`
	PortappsVersion string `json:"portapps_version"`
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
		Info: AppInfo{},
		ID:   id,
		Name: name,
	}

	// Root path
	app.RootPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		app.FatalBox(err)
	}

	// Logfile
	logfolder := utl.CreateFolder(utl.PathJoin(app.RootPath, "log"))
	logpath := utl.PathJoin(logfolder, fmt.Sprintf("%s.log", app.ID))
	app.logfile, err = os.OpenFile(logpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		app.FatalBox(err)
	}

	// Configure logger
	if Log, err = logging.Configure("debug", logpath, LogHook{app}); err != nil {
		app.FatalBox(err)
	}

	// Load info
	infoFile := utl.PathJoin(app.RootPath, "portapp.json")
	infoRaw, err := ioutil.ReadFile(infoFile)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(infoRaw, &app.Info); err != nil {
		return nil, err
	}

	// Startup
	Log.Info().Msg("--------")
	Log.Info().Msgf("Starting %s %s-%s (portapps %s)...", app.Name, app.Info.Version, app.Info.Release, app.Info.PortappsVersion)
	Log.Info().Msgf("Release date: %s", app.Info.Date)
	Log.Info().Msgf("Publisher: %s (%s)", app.Info.Publisher, app.Info.URL)
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

	// Set paths
	app.AppPath = utl.PathJoin(app.RootPath, "app")
	if app.config.Common.AppPath != "" {
		app.AppPath = app.config.Common.AppPath
	}
	app.DataPath = utl.PathJoin(app.RootPath, "data")
	app.WorkingDir = app.AppPath

	// Load env vars from config
	if len(app.config.Common.Env) > 0 {
		Log.Info().Msg("Setting environment variables from config...")
		for key, value := range app.config.Common.Env {
			utl.OverrideEnv(key, app.extendPlaceholders(value))
		}
	}

	return app, nil
}

// Launch to execute the app with additional args
func (app *App) Launch(args []string) {
	Log.Info().Msgf("Process: %s", app.Process)
	Log.Info().Msgf("Args (config file): %s", strings.Join(app.config.Common.Args, " "))
	Log.Info().Msgf("Args (cmd line): %s", strings.Join(args, " "))
	Log.Info().Msgf("Args (hardcoded): %s", strings.Join(app.Args, " "))
	Log.Info().Msgf("Working dir: %s", app.WorkingDir)
	Log.Info().Msgf("App path: %s", app.AppPath)
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

func (app *App) extendPlaceholders(value string) string {
	placeholders := map[string]string{
		"@ROOT_PATH@":    app.RootPath,
		"@APP_PATH@":     app.AppPath,
		"@DATA_PATH@":    app.DataPath,
		"@DRIVE_LETTER@": app.RootPath[:1],
	}
	for placeholder, ext := range placeholders {
		value = strings.Replace(value, placeholder, ext, -1)
	}
	return value
}
