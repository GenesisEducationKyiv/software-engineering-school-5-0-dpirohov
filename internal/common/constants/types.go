package constants

type contextKey string

const TraceID = contextKey("trace_id")

const HdrTraceID = "x-traceid"

const HdrRetries = "x-retries"
const HdrOriginalTopic = "x-original-topic"
