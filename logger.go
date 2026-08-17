package portapps

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/ilya1st/rotatewriter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"golang.org/x/sys/windows"
)

type logWriter interface {
	io.Writer
	CloseWriteFile() error
}

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

	logfolder := filepath.Join(app.RootPath, "log")
	if err := os.MkdirAll(logfolder, 0o755); err != nil {
		return err
	}
	logfile := filepath.Join(logfolder, fmt.Sprintf("%s.log", app.ID))
	rwriter, err := rotatewriter.NewRotateWriter(logfile, 5)
	if err != nil {
		return err
	}
	app.logfile = rwriter

	sighupChan := make(chan os.Signal, 1)
	signal.Notify(sighupChan, windows.SIGHUP)
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
