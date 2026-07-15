package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sekaishopml/cytaxi/backend/gateways/whatsapp/internal/whatsapp/providers"
)

type Handler func(ctx context.Context, msg providers.Message) error

type Receiver struct {
	mu           sync.RWMutex
	secret       string
	verifyToken  string
	handlers     map[providers.MessageType][]Handler
	processedIDs map[string]time.Time
}

func NewReceiver(secret, verifyToken string) *Receiver {
	return &Receiver{
		secret:      secret,
		verifyToken: verifyToken,
		handlers:     make(map[providers.MessageType][]Handler),
		processedIDs: make(map[string]time.Time),
	}
}

func (r *Receiver) VerifySignature(payload []byte, signature string) bool {
	if r.secret == "" || signature == "" {
		return true
	}
	mac := hmac.New(sha256.New, []byte(r.secret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

func (r *Receiver) VerifyWebhookToken(mode, token, challenge string) (string, bool) {
	if mode == "subscribe" && token == r.verifyToken {
		return challenge, true
	}
	return "", false
}

func (r *Receiver) Handle(msgType providers.MessageType, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[msgType] = append(r.handlers[msgType], handler)
}

func (r *Receiver) ProcessWebhook(ctx context.Context, payload []byte) ([]providers.Message, error) {
	var evt providers.WebhookEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		return nil, fmt.Errorf("webhook parse: %w", err)
	}

	var messages []providers.Message

	for _, entry := range evt.Changes {
		for _, change := range entry.Changes {
			val := change.Value

			for _, msg := range val.Messages {
				if r.isProcessed(msg.ID) {
					continue
				}

				message := r.convertWebhookMessage(msg, val.Metadata)
				r.markProcessed(msg.ID)

				handlers := r.handlers[providers.MessageType(msg.Type)]
				for _, h := range handlers {
					if err := h(ctx, message); err != nil {
						return messages, err
					}
				}
				messages = append(messages, message)
			}

			for _, status := range val.Statuses {
				if !r.isProcessed(status.ID) {
					r.markProcessed(status.ID)
				}
			}
		}
	}

	return messages, nil
}

func (r *Receiver) convertWebhookMessage(msg providers.WebhookMessage, meta providers.WebhookMeta) providers.Message {
	message := providers.Message{
		ID:        msg.ID,
		From:      msg.From,
		Type:      providers.MessageType(msg.Type),
		Timestamp: time.Now(),
	}

	switch msg.Type {
	case "text":
		if msg.Text != nil {
			message.Text = &providers.TextBody{Body: msg.Text.Body}
		}
	case "location":
		if msg.Location != nil {
			message.Location = msg.Location
		}
	case "image":
		if msg.Image != nil {
			message.Image = msg.Image
		}
	case "interactive":
		if msg.Interactive != nil {
			message.Interactive = &providers.InteractiveBody{
				Type: msg.Interactive.Type,
				Body: &providers.TextBody{Body: ""},
			}
		}
	}

	return message
}

func (r *Receiver) isProcessed(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.processedIDs[id]
	return ok
}

func (r *Receiver) markProcessed(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.processedIDs[id] = time.Now()
}

func (r *Receiver) CleanupOlderThan(duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cutoff := time.Now().Add(-duration)
	for id, tm := range r.processedIDs {
		if tm.Before(cutoff) {
			delete(r.processedIDs, id)
		}
	}
}
