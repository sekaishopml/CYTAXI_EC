package retry

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type RetryConfig struct {
	MaxRetries  int
	InitialDelay time.Duration
	MaxDelay    time.Duration
	Multiplier   float64
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:   3,
		InitialDelay: 500 * time.Millisecond,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
	}
}

type RetryManager struct {
	config RetryConfig
	logger *slog.Logger
}

func NewRetryManager(config RetryConfig, logger *slog.Logger) *RetryManager {
	return &RetryManager{config: config, logger: logger}
}

func (rm *RetryManager) Execute(ctx context.Context, operation string, fn func(ctx context.Context) error) error {
	var lastErr error
	delay := rm.config.InitialDelay

	for attempt := 0; attempt <= rm.config.MaxRetries; attempt++ {
		if attempt > 0 {
			rm.logger.Warn("retrying", "operation", operation, "attempt", attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * rm.config.Multiplier)
			if delay > rm.config.MaxDelay {
				delay = rm.config.MaxDelay
			}
		}

		err := fn(ctx)
		if err == nil {
			if attempt > 0 {
				rm.logger.Info("retry succeeded", "operation", operation, "attempt", attempt)
			}
			return nil
		}
		lastErr = err
		rm.logger.Warn("operation failed", "operation", operation, "attempt", attempt, "error", err)
	}

	return fmt.Errorf("retry exhausted for %s after %d attempts: %w", operation, rm.config.MaxRetries, lastErr)
}

type DeadLetterQueue struct {
	messages map[string][]DeadLetterMessage
	logger   *slog.Logger
}

type DeadLetterMessage struct {
	ID        string
	EventType string
	Payload   any
	Error     string
	Retries   int
	CreatedAt time.Time
}

type DeadLetterRepository interface {
	Save(ctx context.Context, msg DeadLetterMessage) error
	FindByEventType(ctx context.Context, eventType string) ([]DeadLetterMessage, error)
	Replay(ctx context.Context, id string) error
}

func NewDeadLetterQueue(logger *slog.Logger) *DeadLetterQueue {
	return &DeadLetterQueue{
		messages: make(map[string][]DeadLetterMessage),
		logger:   logger,
	}
}

func (dlq *DeadLetterQueue) Enqueue(ctx context.Context, eventType string, payload any, err error) {
	msg := DeadLetterMessage{
		ID:        fmt.Sprintf("dlq_%d", time.Now().UnixNano()),
		EventType: eventType,
		Payload:   payload,
		Error:     err.Error(),
		CreatedAt: time.Now(),
	}
	dlq.messages[eventType] = append(dlq.messages[eventType], msg)
	dlq.logger.Error("dead letter queued", "event_type", eventType, "error", err)
}

func (dlq *DeadLetterQueue) Pending() []DeadLetterMessage {
	var all []DeadLetterMessage
	for _, msgs := range dlq.messages {
		all = append(all, msgs...)
	}
	return all
}
