package logging

import (
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/ilya1st/rotatewriter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

// Configure configures logger
func Configure(level string, file string, hook zerolog.Hook) (*zerolog.Logger, error) {
	var err error
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logFile := path.Clean(file)
	if err := os.MkdirAll(path.Dir(logFile), os.ModePerm); err != nil {
		return nil, err
	}

	rwriter, err := rotatewriter.NewRotateWriter(logFile, 5)
	if err != nil {
		return nil, err
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
	}).With().Caller().Timestamp().Logger().Hook(hook)

	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	zerolog.SetGlobalLevel(logLevel)

	return &log.Logger, nil
}
