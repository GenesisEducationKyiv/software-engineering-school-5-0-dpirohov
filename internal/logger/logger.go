package logger

import (
	"context"
	"os"
	"time"
	"weatherApi/internal/common/constants"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init(serviceName string, level zerolog.Level) {
	zerolog.TimeFieldFormat = time.RFC3339
	Log = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger().
		Level(level)
}

func FromContext(ctx context.Context) *zerolog.Logger {
	if traceID, ok := ctx.Value(constants.TraceID).(string); ok {
		l := Log.With().Str("trace_id", traceID).Logger()
		return &l
	}
	return &Log
}


func InjectTraceID(ctx context.Context, msg amqp.Delivery) context.Context {
	if traceID, ok := msg.Headers[constants.HdrTraceID].(string); ok {
		return context.WithValue(ctx, constants.TraceID, traceID)
	}
	return ctx
}