package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func GetLogger() zerolog.Logger {
	return zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()
}
