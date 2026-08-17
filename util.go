package portapps

import (
	"path/filepath"

	"github.com/portapps/portapps/v3/pkg/electron"
)

// ElectronAppPath returns the app electron path
func (app *App) ElectronAppPath() string {
	electronAppFolder, err := electron.FindAppFolder("app-", app.AppPath)
	if err != nil {
		app.FatalBoxLog(err.Error())
	}
	return filepath.Join(app.AppPath, electronAppFolder)
}
