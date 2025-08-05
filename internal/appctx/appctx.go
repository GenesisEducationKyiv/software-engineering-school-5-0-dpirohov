package appctx

import (
	"context"
	"weatherApi/internal/common/constants"
)

func SetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, constants.TraceID, traceID)
}

func GetTraceID(ctx context.Context) string {
	traceID, _ := ctx.Value(constants.TraceID).(string) // nolint:errcheck // no check needed
	return traceID
}
