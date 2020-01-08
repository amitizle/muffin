package logger

import (
	"context"
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// ErrNoLoggerInContext is an error that's returned in case we try
	// to fetch logger from a context and it's not existing there
	ErrNoLoggerInContext = errors.New("logger cannot be found in context")
)

// ctxKeyLogger is a custom type that will be used as the key
// of the logger in context.Context
type ctxKeyLogger int

// LoggerCtxKey is the key that holds the unique logger ID in a context.
const LoggerCtxKey ctxKeyLogger = 0

// Init initializing zerolog with the given configuration
func Init(level string) error {
	log.Debug().Msgf("initializing log with level %s", level)
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.DisableSampling(true)
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	return nil
}

// StoreContext stores a zerolog `Logger` in a context and returns the
// new context
func StoreContext(ctx context.Context, lg zerolog.Logger) context.Context {
	return context.WithValue(ctx, LoggerCtxKey, lg)
}

// GetContext returns a zerolog.Logger from a context
func GetContext(ctx context.Context) (zerolog.Logger, error) {
	lg, ok := ctx.Value(LoggerCtxKey).(zerolog.Logger)
	if !ok {
		return zerolog.Logger{}, ErrNoLoggerInContext
	}
	return lg, nil
}
