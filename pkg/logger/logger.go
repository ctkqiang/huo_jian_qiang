package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type LoggerInterface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct {
	logger *zerolog.Logger
}

func New(level string) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn", "warning":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &Logger{
		logger: &logger,
	}
}

func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.logMessage(zerolog.DebugLevel, message, args...)
}

func (l *Logger) logMessage(level zerolog.Level, message interface{}, args ...interface{}) {
	event := l.logger.WithLevel(level)

	if msg, ok := message.(string); ok {
		event.Msgf(msg, args...)
	} else {
		event.Interface("message", message).Send()
	}
}
