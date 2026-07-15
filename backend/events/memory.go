package events

import (
	"fmt"
	"sync"
)

type memoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewMemoryBus() Bus {
	return &memoryBus{
		handlers: make(map[string][]Handler),
	}
}

func (b *memoryBus) Publish(event DomainEvent) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	handlers, ok := b.handlers[event.Type]
	if !ok {
		return nil
	}

	var errs []error
	for _, h := range handlers {
		if err := h(event); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("events: %v", errs)
	}
	return nil
}

func (b *memoryBus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}
