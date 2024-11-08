package portapps

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/pkg/errors"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/utl"
	"github.com/portapps/portapps/v3/pkg/win"
	"gopkg.in/yaml.v3"
)

// App represents an active app object
type App struct {
	Info       AppInfo
	Prev       AppPrev
	WinVersion win.Version

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

// AppPrev represents previous settings file structure
type AppPrev struct {
	Info       AppInfo     `json:"info"`
	WinVersion win.Version `json:"win_version"`
	RootPath   string      `json:"root_path"`
	AppPath    string      `json:"app_path"`
	DataPath   string      `json:"data_path"`
}

// New creates new app instance
func New(id string, name string) (app *App, err error) {
	return NewWithCfg(id, name, nil)
}

// NewWithCfg creates new app instance with an app config
func NewWithCfg(id string, name string, appcfg interface{}) (app *App, err error) {
	// Init
	app = &App{
		Info: AppInfo{},
		Prev: AppPrev{},
		ID:   id,
		Name: name,
	}

	// WinVersion
	app.WinVersion, err = win.GetVersion()
	if err != nil {
		app.FatalBox(errors.Wrap(err, "Cannot get Windows version"))
	}

	// Root path
	ex, err := os.Executable()
	if err != nil {
		app.FatalBox(errors.Wrap(err, "Cannot get path name of the executable"))
	}
	app.RootPath, err = filepath.Abs(filepath.Dir(ex))
	if err != nil {
		app.FatalBox(errors.Wrap(err, "Cannot get root absolute path"))
	}

	// Load info
	infoFile := utl.PathJoin(app.RootPath, "portapp.json")
	infoRaw, err := os.ReadFile(infoFile)
	if err != nil {
		app.FatalBox(errors.Wrap(err, "Cannot load portapps.json"))
	}
	if err = json.Unmarshal(infoRaw, &app.Info); err != nil {
		app.FatalBox(errors.Wrap(err, "Cannot unmarshal portapps.json"))
	}

	// Load config
	if err = app.loadConfig(appcfg); err != nil {
		app.FatalBox(errors.Wrap(err, "Cannot load configuration"))
	}
	if appcfg != nil {
		if err = mapstructure.Decode(app.config.App, appcfg); err != nil {
			app.FatalBox(errors.Wrap(err, fmt.Sprintf("Cannot decode %s configuration", app.Name)))
		}
	}

	// Init logger
	if err = app.InitLogger(); err != nil {
		app.FatalBox(errors.Wrap(err, "Cannot configure logger"))
	}

	// Startup
	log.Info().Msg("--------")
	log.Info().Msgf("Operating System: Windows %d.%d.%d", app.WinVersion.Major, app.WinVersion.Minor, app.WinVersion.Build)
	log.Info().Msgf("Starting %s %s-%s (portapps %s)...", app.Name, app.Info.Version, app.Info.Release, app.Info.PortappsVersion)
	log.Info().Msgf("Release date: %s", app.Info.Date)
	log.Info().Msgf("Publisher: %s (%s)", app.Info.Publisher, app.Info.URL)
	log.Info().Msgf("Root path: %s", app.RootPath)

	// Display config
	b, _ := yaml.Marshal(app.config)
	log.Info().Msgf("Configuration:\n%s", string(b))

	// Set paths
	app.AppPath = utl.PathJoin(app.RootPath, "app")
	if app.config.Common.AppPath != "" {
		app.AppPath = app.config.Common.AppPath
	}
	app.DataPath = utl.PathJoin(app.RootPath, "data")
	app.WorkingDir = app.AppPath

	// Load previous
	prevFile := utl.PathJoin(app.RootPath, "portapp-prev.json")
	if utl.Exists(prevFile) {
		prevRaw, err := os.ReadFile(prevFile)
		if err != nil {
			app.FatalBox(errors.Wrap(err, "Cannot load portapp-prev"))
		}
		if err = json.Unmarshal(prevRaw, &app.Prev); err != nil {
			log.Error().Err(err).Msgf("Cannot unmarshal portapp-prev")
			_ = os.Remove(prevFile)
		}
	}

	// Load env vars from config
	if len(app.config.Common.Env) > 0 {
		log.Info().Msg("Setting environment variables from config...")
		for key, value := range app.config.Common.Env {
			os.Setenv(key, app.extendPlaceholders(value))
		}
	}

	return app, nil
}

// Config returns app configuration
func (app *App) Config() *Config {
	return app.config
}

// Launch to execute the app with additional args
func (app *App) Launch(args []string) {
	log.Info().Msgf("Process: %s", app.Process)
	log.Info().Msgf("Args (config file): %s", strings.Join(app.config.Common.Args, " "))
	log.Info().Msgf("Args (cmd line): %s", strings.Join(args, " "))
	log.Info().Msgf("Args (hardcoded): %s", strings.Join(app.Args, " "))
	log.Info().Msgf("Working dir: %s", app.WorkingDir)
	log.Info().Msgf("App path: %s", app.AppPath)
	log.Info().Msgf("Data path: %s", app.DataPath)
	log.Info().Msgf("Previous path: %s", app.Prev.RootPath)

	if !utl.Exists(app.Process) {
		log.Fatal().Msgf("Application not found in %s", app.Process)
	}

	log.Info().Msgf("Launching %s", app.Name)
	jArgs := append(append(app.config.Common.Args, args...), app.Args...)
	execute := exec.Command(app.Process, jArgs...)
	execute.Dir = app.WorkingDir

	if !app.config.Common.DisableLog {
		execute.Stdout = app.logfile
		execute.Stderr = app.logfile
	}

	log.Info().Msgf("Exec %s %s", app.Process, strings.Join(jArgs, " "))
	if err := execute.Run(); err != nil {
		log.Fatal().Err(err).Msg("Command failed")
	}
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

// Close closes the app
func (app *App) Close() {
	log.Info().Msgf("Closing %s", app.Name)

	// Update previous
	prevFile := utl.PathJoin(app.RootPath, "portapp-prev.json")
	jsonPrev, err := json.MarshalIndent(AppPrev{
		Info:       app.Info,
		WinVersion: app.WinVersion,
		RootPath:   app.RootPath,
		AppPath:    app.AppPath,
		DataPath:   app.DataPath,
	}, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Cannot marshal portapp-prev")
	}
	err = os.WriteFile(prevFile, jsonPrev, 0o644)
	if err != nil {
		log.Error().Err(err).Msg("Cannot write portapp-prev")
	}
}
