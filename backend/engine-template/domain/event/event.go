package event

type DomainEvent struct {
	ID      string
	Type    string
	Version int
	Payload any
}
