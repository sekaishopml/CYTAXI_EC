package observability

import (
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/sekaishopml/cytaxi/backend/integration/contracts"
	"github.com/sekaishopml/cytaxi/backend/integration/correlation"
)

type IntegrationMetrics struct {
	EventsPublished atomic.Int64
	EventsConsumed  atomic.Int64
	EventsFailed    atomic.Int64
	SagasStarted    atomic.Int64
	SagasCompleted  atomic.Int64
	SagasFailed     atomic.Int64
	RetriesDone     atomic.Int64
	DeadLetters     atomic.Int64
}

type Observer struct {
	metrics   *IntegrationMetrics
	corr      *correlation.Manager
	logger    *slog.Logger
}

func NewObserver(corr *correlation.Manager, logger *slog.Logger) *Observer {
	return &Observer{
		metrics: &IntegrationMetrics{},
		corr:    corr,
		logger:  logger,
	}
}

func (o *Observer) Snapshot() IntegrationMetrics {
	return *o.metrics
}

func (o *Observer) OnPublish(envelope contracts.EventEnvelope) {
	o.metrics.EventsPublished.Add(1)
	o.logger.Info("event published",
		"event_type", envelope.Type,
		"event_id", envelope.ID,
		"source", envelope.Source,
	)
}

func (o *Observer) OnConsume(envelope contracts.EventEnvelope) {
	o.metrics.EventsConsumed.Add(1)
	o.logger.Info("event consumed",
		"event_type", envelope.Type,
		"event_id", envelope.ID,
		"correlation_id", envelope.CorrelationID,
	)
}

func (o *Observer) OnError(eventType string, err error) {
	o.metrics.EventsFailed.Add(1)
	o.logger.Error("event error", "event_type", eventType, "error", err)
}

func (o *Observer) OnSagaStarted(name string) {
	o.metrics.SagasStarted.Add(1)
}

func (o *Observer) OnSagaCompleted(name string) {
	o.metrics.SagasCompleted.Add(1)
}

func (o *Observer) OnSagaFailed(name string, err error) {
	o.metrics.SagasFailed.Add(1)
}

func (o *Observer) OnRetry(operation string) {
	o.metrics.RetriesDone.Add(1)
}

func (o *Observer) OnDeadLetter(eventType string) {
	o.metrics.DeadLetters.Add(1)
}

type HealthCheck struct {
	bus    contracts.Bus
	logger *slog.Logger
}

func NewHealthCheck(bus contracts.Bus, logger *slog.Logger) *HealthCheck {
	return &HealthCheck{bus: bus, logger: logger}
}

func (h *HealthCheck) Status() map[string]any {
	return map[string]any{
		"service": "integration-layer",
		"status":  "ok",
		"time":    time.Now().UTC(),
	}
}
