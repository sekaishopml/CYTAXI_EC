package event

import "time"

type EventType string

const (
	ConversationStarted EventType = "conversation.started"
	ConversationClosed  EventType = "conversation.closed"

	SessionCreated EventType = "session.created"
	SessionExpired EventType = "session.expired"
	SessionClosed  EventType = "session.closed"

	MessageReceived EventType = "message.received"
	MessageSent     EventType = "message.sent"
	StateChanged    EventType = "conversation.state_changed"
)

type ConversationStartedEvent struct {
	ConversationID string
	Phone          string
	StartedAt      time.Time
}

type ConversationClosedEvent struct {
	ConversationID string
	ClosedAt       time.Time
}

type SessionCreatedEvent struct {
	SessionID      string
	ConversationID string
	Phone          string
	CreatedAt      time.Time
}

type SessionExpiredEvent struct {
	SessionID      string
	ConversationID string
	ExpiredAt      time.Time
}

type SessionClosedEvent struct {
	SessionID      string
	ConversationID string
	ClosedAt       time.Time
}

type MessageReceivedEvent struct {
	MessageID      string
	ConversationID string
	Content        string
	Role           string
	ReceivedAt     time.Time
}

type MessageSentEvent struct {
	MessageID      string
	ConversationID string
	Content        string
	SentAt         time.Time
}

type StateChangedEvent struct {
	ConversationID string
	PreviousState  string
	NewState       string
	ChangedAt      time.Time
}
