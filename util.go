package portapps

import (
	"path/filepath"

	"github.com/portapps/portapps/v3/pkg/utl"
)

// ElectronAppPath returns the app electron path
func (app *App) ElectronAppPath() string {
	electronAppFolder, err := utl.FindElectronAppFolder("app-", app.AppPath)
	if err != nil {
		app.FatalBoxLog(err.Error())
	}
	return filepath.Join(app.AppPath, electronAppFolder)
}
