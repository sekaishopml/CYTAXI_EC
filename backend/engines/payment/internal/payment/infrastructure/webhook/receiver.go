package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type Event struct {
	ID        string
	Type      string
	Provider  string
	Payload   map[string]any
	Signature string
	Timestamp time.Time
}

type EventType string

const (
	EventPaymentSucceeded EventType = "payment_intent.succeeded"
	EventPaymentFailed    EventType = "payment_intent.failed"
	EventCaptureCompleted EventType = "capture.completed"
	EventRefundProcessed  EventType = "refund.processed"
	EventRefundFailed     EventType = "refund.failed"
	EventDisputeCreated   EventType = "dispute.created"
	EventDisputeResolved  EventType = "dispute.resolved"
)

type Handler func(event Event) error

type Receiver struct {
	mu          sync.RWMutex
	secret      string
	handlers    map[EventType][]Handler
	processedIDs map[string]bool // idempotency cache
}

func NewReceiver(secret string) *Receiver {
	return &Receiver{
		secret:      secret,
		handlers:    make(map[EventType][]Handler),
		processedIDs: make(map[string]bool),
	}
}

func (r *Receiver) VerifySignature(payload []byte, signature string) bool {
	if r.secret == "" || signature == "" {
		return true // mock mode
	}
	mac := hmac.New(sha256.New, []byte(r.secret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

func (r *Receiver) Handle(eventType EventType, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[eventType] = append(r.handlers[eventType], handler)
}

func (r *Receiver) Process(event Event) error {
	r.mu.Lock()
	if r.processedIDs[event.ID] {
		r.mu.Unlock()
		return fmt.Errorf("webhook: event %s already processed (idempotent)", event.ID)
	}
	r.processedIDs[event.ID] = true
	r.mu.Unlock()

	handlers := r.handlers[EventType(event.Type)]
	if len(handlers) == 0 {
		return fmt.Errorf("webhook: no handler for event type %s", event.Type)
	}

	var errs []error
	for _, h := range handlers {
		if err := h(event); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("webhook: %d handlers failed", len(errs))
	}
	return nil
}

func (r *Receiver) CleanupOlderThan(duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cutoff := time.Now().Add(-duration)
	for id, tm := range r.processedIDs {
		if tm.ProcessedAt.Before(cutoff) {
			delete(r.processedIDs, id)
		}
	}
}

type ProcessedEvent struct {
	ProcessedAt time.Time
}
