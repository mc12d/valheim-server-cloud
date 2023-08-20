package app

import (
	"os"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

func Zerolog() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.Stamp,
	}).With().Timestamp().Logger()
}

func CronLogger(z zerolog.Logger) cron.Logger {
	return cronZerolog{z}
}

type cronZerolog struct {
	logger zerolog.Logger
}

func (c cronZerolog) Info(msg string, keysAndValues ...interface{}) {
	c.logger.Info().Msgf(
		formatString(len(keysAndValues)+2),
		append([]interface{}{msg}, formatTimes(keysAndValues)...)...,
	)
}

func (c cronZerolog) Error(err error, msg string, keysAndValues ...interface{}) {
	c.logger.Error().Err(err).Msgf(
		formatString(len(keysAndValues)+2),
		append([]interface{}{msg}, formatTimes(keysAndValues)...)...,
	)
}

func formatString(numKeysAndValues int) string {
	var sb strings.Builder
	sb.WriteString("%s")
	if numKeysAndValues > 0 {
		sb.WriteString(", ")
	}
	for i := 0; i < numKeysAndValues/2; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("%v=%v")
	}
	return sb.String()
}

func formatTimes(keysAndValues []interface{}) []interface{} {
	var formattedArgs []interface{}
	for _, arg := range keysAndValues {
		if t, ok := arg.(time.Time); ok {
			arg = t.Format(time.RFC3339)
		}
		formattedArgs = append(formattedArgs, arg)
	}
	return formattedArgs
}
