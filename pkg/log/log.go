// Package log source from https://github.com/rs/zerolog/tree/master/log
package log

import (
	"context"
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Output duplicates the global logger and sets w as its output.
func Output(w io.Writer) zerolog.Logger {
	return log.Logger.Output(w)
}

// With creates a child logger with the field added to its context.
func With() zerolog.Context {
	return log.Logger.With()
}

// Level creates a child logger with the minimum accepted level set to level.
func Level(level zerolog.Level) zerolog.Logger {
	return log.Logger.Level(level)
}

// Sample returns a logger with the s sampler.
func Sample(s zerolog.Sampler) zerolog.Logger {
	return log.Logger.Sample(s)
}

// Hook returns a logger with the h Hook.
func Hook(h zerolog.Hook) zerolog.Logger {
	return log.Logger.Hook(h)
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func Err(err error) *zerolog.Event {
	return log.Logger.Err(err)
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func Trace() *zerolog.Event {
	return log.Logger.Trace()
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return log.Logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return log.Logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *zerolog.Event {
	return log.Logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return log.Logger.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return log.Logger.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *zerolog.Event {
	return log.Logger.Panic()
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func WithLevel(level zerolog.Level) *zerolog.Event {
	return log.Logger.WithLevel(level)
}

// Log starts a new message with no level. Setting zerolog.GlobalLevel to
// zerolog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func Log() *zerolog.Event {
	return log.Logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	log.Logger.Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	log.Logger.Printf(format, v...)
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
