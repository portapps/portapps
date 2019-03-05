package portapps

import (
	"os"
	"path/filepath"

	"github.com/google/logger"
	_ "github.com/josephspurrier/goversioninfo"
	"github.com/mitchellh/mapstructure"
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
func Init(appcfg interface{}) {
	var err error

	Papp.Path, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Log.Fatal("Current path:", err)
	}

	Papp.DataPath = AppPathJoin("data")

	logfile, err = os.OpenFile(PathJoin(Papp.Path, Papp.ID+".log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		Log.Fatal("Cannot open log file:", err)
	}

	// Startup
	Log = logger.Init(Papp.Name, false, false, logfile)
	Log.Info("--------")
	Log.Infof("Starting %s...", Papp.Name)
	Log.Infof("Current path: %s", Papp.Path)

	// Configuration
	Log.Info("Loading configuration...")
	if err = loadConfig(appcfg); err != nil {
		Log.Fatal("Cannot load configuration:", err)
	}
	if err := mapstructure.Decode(Papp.config.App, appcfg); err != nil {
		Log.Fatal("Cannot decode app configuration:", err)
	}
	b, _ := yaml.Marshal(Papp.config)
	Log.Infof("Configuration:\n%s", string(b))
}
