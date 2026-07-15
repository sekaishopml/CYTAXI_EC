package whatsapp

import "context"

type EventType string

const (
	EventConnected       EventType = "whatsapp.connected"
	EventDisconnected    EventType = "whatsapp.disconnected"
	EventQRReceived      EventType = "whatsapp.qr_received"
	EventMessageReceived EventType = "whatsapp.message_received"
	EventMessageSent     EventType = "whatsapp.message_sent"
	EventError           EventType = "whatsapp.error"
)

type Event struct {
	Type    EventType
	Session string
	Data    any
}

type EventHandler func(event Event)

type EventBus struct {
	handlers []EventHandler
}

func NewEventBus() *EventBus {
	return &EventBus{}
}

func (eb *EventBus) Subscribe(handler EventHandler) {
	eb.handlers = append(eb.handlers, handler)
}

func (eb *EventBus) Publish(ctx context.Context, event Event) {
	for _, h := range eb.handlers {
		h(event)
	}
}
