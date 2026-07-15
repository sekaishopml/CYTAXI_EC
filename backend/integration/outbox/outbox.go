package outbox

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/sekaishopml/cytaxi/backend/integration/contracts"
)

type OutboxMessage struct {
	ID          string
	EventType   string
	Payload     any
	CreatedAt   time.Time
	PublishedAt *time.Time
	Retries     int
	Status      string // pending, published, failed
}

type OutboxRepository interface {
	Save(ctx context.Context, msg OutboxMessage) error
	FindPending(ctx context.Context, limit int) ([]OutboxMessage, error)
	MarkPublished(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type OutboxPublisher struct {
	repo   OutboxRepository
	bus    contracts.Bus
	logger *slog.Logger
	mu     sync.Mutex
}

func NewOutboxPublisher(repo OutboxRepository, bus contracts.Bus, logger *slog.Logger) *OutboxPublisher {
	return &OutboxPublisher{repo: repo, bus: bus, logger: logger}
}

func (op *OutboxPublisher) Save(ctx context.Context, eventType, source string, payload any) error {
	envelope := contracts.NewEnvelope(eventType, source, payload)
	msg := OutboxMessage{
		ID:        envelope.ID,
		EventType: eventType,
		Payload:   envelope,
		CreatedAt: time.Now(),
		Status:    "pending",
	}
	return op.repo.Save(ctx, msg)
}

func (op *OutboxPublisher) PublishPending(ctx context.Context) (int, error) {
	messages, err := op.repo.FindPending(ctx, 50)
	if err != nil {
		return 0, fmt.Errorf("find pending: %w", err)
	}

	var published int
	for _, msg := range messages {
		env, ok := msg.Payload.(contracts.EventEnvelope)
		if !ok {
			op.repo.MarkFailed(ctx, msg.ID)
			continue
		}
		if err := op.bus.Publish(ctx, env); err != nil {
			op.repo.MarkFailed(ctx, msg.ID)
			op.logger.Error("outbox publish failed", "id", msg.ID, "error", err)
			continue
		}
		op.repo.MarkPublished(ctx, msg.ID)
		published++
	}
	return published, nil
}

type InboxMessage struct {
	ID          string
	EventID     string
	EventType   string
	Payload     any
	ProcessedAt *time.Time
	Status      string
}

type InboxRepository interface {
	Save(ctx context.Context, msg InboxMessage) error
	FindByEventID(ctx context.Context, eventID string) (*InboxMessage, error)
	MarkProcessed(ctx context.Context, id string) error
}

type InboxProcessor struct {
	repo     InboxRepository
	handlers map[string]contracts.EventHandler
	logger   *slog.Logger
}

func NewInboxProcessor(repo InboxRepository, logger *slog.Logger) *InboxProcessor {
	return &InboxProcessor{
		repo:     repo,
		handlers: make(map[string]contracts.EventHandler),
		logger:   logger,
	}
}

func (ip *InboxProcessor) RegisterHandler(eventType string, handler contracts.EventHandler) {
	ip.handlers[eventType] = handler
}

func (ip *InboxProcessor) Process(ctx context.Context, envelope contracts.EventEnvelope) error {
	existing, _ := ip.repo.FindByEventID(ctx, envelope.ID)
	if existing != nil && existing.Status == "processed" {
		return nil
	}

	msg := InboxMessage{
		ID:        fmt.Sprintf("inb_%d", time.Now().UnixNano()),
		EventID:   envelope.ID,
		EventType: envelope.Type,
		Payload:   envelope,
		Status:    "received",
	}
	ip.repo.Save(ctx, msg)

	handler, ok := ip.handlers[envelope.Type]
	if !ok {
		return fmt.Errorf("no handler for event %s", envelope.Type)
	}

	if err := handler(ctx, envelope); err != nil {
		ip.logger.Error("inbox handler failed", "event", envelope.Type, "error", err)
		return err
	}

	return ip.repo.MarkProcessed(ctx, msg.ID)
}
