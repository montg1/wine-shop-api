package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Log is the global structured logger
var Log zerolog.Logger

// Init sets up the global logger.
// In debug/development mode, uses pretty console output.
// In release/production mode, uses JSON output.
func Init(ginMode string) {
	if ginMode == "release" {
		// Production: JSON logs for machine parsing
		Log = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	} else {
		// Development: human-friendly console output
		Log = zerolog.New(
			zerolog.ConsoleWriter{Out: os.Stderr},
		).With().Timestamp().Caller().Logger()
	}
}
