package entity

import "time"

type ConversationID string

type Conversation struct {
	ID        ConversationID
	Phone     string
	Status    ConversationStatus
	State     ConversationState
	SessionID SessionID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ConversationStatus string

const (
	ConversationStatusActive    ConversationStatus = "active"
	ConversationStatusWaiting   ConversationStatus = "waiting"
	ConversationStatusClosed    ConversationStatus = "closed"
)
