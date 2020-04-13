package portapps

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilya1st/rotatewriter"
	"github.com/portapps/portapps/v2/pkg/utl"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

// InitLogger configures logger
func (app *App) InitLogger() error {
	dialogHook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		if level == zerolog.FatalLevel {
			app.ErrorBox(msg)
		}
	})

	if app.config.Common.DisableLog {
		log.Logger = zerolog.New(zerolog.Nop()).With().Logger().Hook(dialogHook)
		return nil
	}

	var err error
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logfolder := utl.CreateFolder(utl.PathJoin(app.RootPath, "log"))
	logfile := utl.PathJoin(logfolder, fmt.Sprintf("%s.log", app.ID))
	app.logfile, err = os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rwriter, err := rotatewriter.NewRotateWriter(logfile, 5)
	if err != nil {
		return err
	}

	sighupChan := make(chan os.Signal, 1)
	signal.Notify(sighupChan, syscall.SIGHUP)
	go func() {
		for {
			_, ok := <-sighupChan
			if !ok {
				return
			}
			rwriter.Rotate(nil)
		}
	}()

	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        rwriter,
		TimeFormat: time.RFC1123,
		NoColor:    true,
	}).With().Caller().Timestamp().Logger().Hook(dialogHook)

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	return nil
}
