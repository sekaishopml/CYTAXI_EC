package telemetry

import (
	"context"
	"crypto/rand"
	"fmt"
)

type ctxTelemetry string

const telemetryKey ctxTelemetry = "telemetry"

type contextTracer struct{}

func NewTracer() Tracer {
	return &contextTracer{}
}

func (t *contextTracer) StartSpan(name string) Span {
	return &contextSpan{
		name:      name,
		traceID:   newID(),
		spanID:    newID(),
		startTime: now(),
	}
}

type contextSpan struct {
	name      string
	traceID   string
	spanID    string
	startTime int64
	attributes map[string]string
}

func (s *contextSpan) End() {}

func (s *contextSpan) SetAttribute(key, value string) {
	if s.attributes == nil {
		s.attributes = make(map[string]string)
	}
	s.attributes[key] = value
}

func (s *contextSpan) RecordError(err error) {
	s.SetAttribute("error", err.Error())
}

func WithSpan(ctx context.Context, name string) (context.Context, Span) {
	t := NewTracer()
	span := t.StartSpan(name)
	return context.WithValue(ctx, telemetryKey, span), span
}

func SpanFromContext(ctx context.Context) Span {
	span, ok := ctx.Value(telemetryKey).(Span)
	if !ok {
		return nil
	}
	return span
}

func newID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func now() int64 {
	return 0
}
