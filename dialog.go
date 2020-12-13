package portapps

import (
	"fmt"
	"os"

	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/win"
)

// ErrorBox display an error message box
func (app *App) ErrorBox(msg interface{}) {
	_, _ = win.MsgBox(
		fmt.Sprintf("%s portable", app.Name),
		fmt.Sprintf("%v", msg),
		win.MsgBoxBtnOk|win.MsgBoxIconError)
}

// ErrorBoxLog display an error message box nad log
func (app *App) ErrorBoxLog(msg interface{}) {
	log.Error().Msgf("%s", msg)
	app.ErrorBox(msg)
}

// FatalBox display an error message box and exit
func (app *App) FatalBox(msg interface{}) {
	app.ErrorBox(msg)
	os.Exit(1)
}

// FatalBoxLog display an error message box, log and exit
func (app *App) FatalBoxLog(msg interface{}) {
	log.Error().Msgf("%s", msg)
	app.FatalBox(msg)
}
