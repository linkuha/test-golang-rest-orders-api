package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
)

func InitLogger(logLvl string, customLog io.Writer) {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var l zerolog.Level
	switch strings.ToLower(logLvl) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	if customLog != nil {
		multi := zerolog.MultiLevelWriter(consoleWriter, customLog)
		log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
