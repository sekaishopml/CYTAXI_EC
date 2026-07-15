package events

type DomainEvent struct {
	ID      string
	Type    string
	Version int
	Payload any
}

type Handler func(event DomainEvent) error

type Bus interface {
	Publish(event DomainEvent) error
	Subscribe(eventType string, handler Handler)
}
