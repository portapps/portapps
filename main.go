package portapps

import (
	_ "github.com/josephspurrier/goversioninfo"

	"os"
	"path/filepath"

	"github.com/google/logger"
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
