package eventbus

import (
	"context"
	"sync"

	"github.com/sekaishopml/cytaxi/backend/integration/contracts"
)

type Bus interface {
	Publish(ctx context.Context, envelope contracts.EventEnvelope) error
	Subscribe(eventType string, handler contracts.EventHandler)
	Unsubscribe(eventType string)
}

type BrokerProvider interface {
	Name() string
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(ctx context.Context, topic string, handler func(data []byte)) error
}

type MemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]contracts.EventHandler
}

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{handlers: make(map[string][]contracts.EventHandler)}
}

func (b *MemoryBus) Publish(ctx context.Context, envelope contracts.EventEnvelope) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, h := range b.handlers[envelope.Type] {
		go h(ctx, envelope)
	}
	return nil
}

func (b *MemoryBus) Subscribe(eventType string, handler contracts.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *MemoryBus) Unsubscribe(eventType string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.handlers, eventType)
}
