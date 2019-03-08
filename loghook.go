package portapps

import (
	"github.com/rs/zerolog"
)

// LogHook is a logging hook for zerolog
type LogHook struct {
	app *App
}

// Run logging hook
func (h LogHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level == zerolog.FatalLevel {
		h.app.ErrorBox(msg)
	}
}
