package telemetry

type Span interface {
	End()
	SetAttribute(key, value string)
	RecordError(err error)
}

type Tracer interface {
	StartSpan(name string) Span
}

type ContextKey string

const TraceIDKey ContextKey = "trace_id"
const SpanIDKey ContextKey = "span_id"
