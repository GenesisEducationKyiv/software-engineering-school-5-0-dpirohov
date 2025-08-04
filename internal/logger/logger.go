package logger

import (
	"context"
	"io"
	"os"
	"weatherApi/internal/appctx"

	"github.com/rs/zerolog"
)

type Logger struct {
	base zerolog.Logger
}

func NewLogger(serviceName string, level zerolog.Level) *Logger {
	return NewWithWriter(serviceName, level, os.Stdout)
}

func NewWithWriter(serviceName string, level zerolog.Level, writer io.Writer) *Logger {
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z07:00"

	base := zerolog.New(writer).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger().
		Level(level)

	return &Logger{base: base}
}

// Test logger
func NewNoOpLogger() *Logger {
	return NewWithWriter("test-service", zerolog.InfoLevel, io.Discard)
}

func (l *Logger) Base() *zerolog.Logger {
	return &l.base
}

// Add trace_id from context to logger
func (l *Logger) FromContext(ctx context.Context) *zerolog.Logger {
	traceID := appctx.GetTraceID(ctx)
	loggerWithTraceID := l.base.With().Str("trace_id", traceID).Logger()
	return &loggerWithTraceID
}
