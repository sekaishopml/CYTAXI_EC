package correlation

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"
)

type Provider interface {
	GetCorrelationID(ctx context.Context) string
	GetTraceID(ctx context.Context) string
	NewTraceID() string
	WithCorrelation(ctx context.Context, correlationID string) context.Context
	WithTrace(ctx context.Context, traceID string) context.Context
}

type Manager struct{}

func NewManager() *Manager { return &Manager{} }

type ctxKey string

const correlationKey ctxKey = "correlation_id"
const traceKey ctxKey = "trace_id"

func (m *Manager) GetCorrelationID(ctx context.Context) string {
	id, _ := ctx.Value(correlationKey).(string)
	if id == "" {
		return newID("corr")
	}
	return id
}

func (m *Manager) GetTraceID(ctx context.Context) string {
	id, _ := ctx.Value(traceKey).(string)
	if id == "" {
		return newID("trace")
	}
	return id
}

func (m *Manager) NewTraceID() string {
	return newID("trace")
}

func (m *Manager) WithCorrelation(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationKey, correlationID)
}

func (m *Manager) WithTrace(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceKey, traceID)
}

func newID(prefix string) string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%s_%x", prefix, b)
}

type TraceProvider interface {
	StartSpan(ctx context.Context, name string) (context.Context, Span)
	GetTraceID(ctx context.Context) string
	GetSpanID(ctx context.Context) string
}

type Span struct {
	TraceID string
	SpanID  string
	Name    string
	StartAt time.Time
	EndAt   *time.Time
}

type TraceManager struct {
	manager *Manager
}

func NewTraceManager(manager *Manager) *TraceManager {
	return &TraceManager{manager: manager}
}

func (tm *TraceManager) StartSpan(ctx context.Context, name string) (context.Context, Span) {
	traceID := tm.manager.GetTraceID(ctx)
	span := Span{
		TraceID: traceID,
		SpanID:  newID("span"),
		Name:    name,
		StartAt: time.Now(),
	}
	return ctx, span
}

func (tm *TraceManager) GetTraceID(ctx context.Context) string {
	return tm.manager.GetTraceID(ctx)
}

func (tm *TraceManager) GetSpanID(ctx context.Context) string {
	return ""
}
