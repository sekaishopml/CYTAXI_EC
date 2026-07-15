package contracts

import (
	"context"
	"fmt"
	"time"
)

type EventEnvelope struct {
	ID            string            `json:"id"`
	Type          string            `json:"type"`
	Source        string            `json:"source"`
	Payload       any               `json:"payload"`
	Version       int               `json:"version"`
	CorrelationID string            `json:"correlation_id"`
	TraceID       string            `json:"trace_id"`
	Timestamp     time.Time         `json:"timestamp"`
	Headers       map[string]string `json:"headers,omitempty"`
}

func NewEnvelope(eventType, source string, payload any) EventEnvelope {
	return EventEnvelope{
		ID:        fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		Type:      eventType,
		Source:    source,
		Payload:   payload,
		Version:   1,
		Timestamp: time.Now(),
		Headers:   make(map[string]string),
	}
}

type CommandEnvelope struct {
	ID            string        `json:"id"`
	Command       string        `json:"command"`
	CorrelationID string        `json:"correlation_id"`
	TraceID       string        `json:"trace_id"`
	Timeout       time.Duration `json:"timeout"`
	Payload       any           `json:"payload"`
}

type IntegrationResult struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Ref     string `json:"ref,omitempty"`
}

type EventHandler func(ctx context.Context, envelope EventEnvelope) error
type CommandHandler func(ctx context.Context, cmd CommandEnvelope) (*IntegrationResult, error)

const (
	EventIntegrationStarted   = "integration.started"
	EventIntegrationCompleted = "integration.completed"
	EventEventPublished       = "integration.event_published"
	EventEventConsumed        = "integration.event_consumed"
	EventSagaStarted          = "integration.saga.started"
	EventSagaCompleted        = "integration.saga.completed"
	EventSagaCompensating     = "integration.saga.compensating"
	EventRetryExecuted        = "integration.retry.executed"
	EventDeadLetterCreated    = "integration.deadletter.created"
	EventCorrelationCreated   = "integration.correlation.created"
)
