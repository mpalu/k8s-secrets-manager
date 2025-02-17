package logging

import (
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	log = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func GetLogger() *zerolog.Logger {
	return &log
}
