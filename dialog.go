package portapps

import (
	"fmt"
	"os"

	"github.com/portapps/portapps/v2/pkg/dialog"
)

// ErrorBox display an error message box
func (app *App) ErrorBox(msg interface{}) {
	_, _ = dialog.MsgBox(
		fmt.Sprintf("%s portable", app.Name),
		fmt.Sprintf("%v", msg),
		dialog.MsgBoxBtnOk|dialog.MsgBoxIconError)
}

// FatalBox display an error message box and exit
func (app *App) FatalBox(msg interface{}) {
	app.ErrorBox(msg)
	os.Exit(1)
}
