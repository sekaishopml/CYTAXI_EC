package entity

import "time"

type MessageID string

type Message struct {
	ID             MessageID
	ConversationID ConversationID
	Content        string
	Role           MessageRole
	CreatedAt      time.Time
}

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleSystem    MessageRole = "system"
)
