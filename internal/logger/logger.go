package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func GetLogger() zerolog.Logger {
	// Wrap with SyncWriter to force immediate writes
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}

	return zerolog.New(
		zerolog.SyncWriter(consoleWriter), // This forces immediate flush
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()
}
