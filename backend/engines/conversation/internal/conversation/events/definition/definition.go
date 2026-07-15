package definition

const (
	EventConversationStarted = "conversation.started"
	EventMessageReceived     = "message.received"
	EventConversationClosed  = "conversation.closed"
)

type ConversationStartedPayload struct {
	ConversationID string `json:"conversation_id"`
	Phone          string `json:"phone"`
}

type MessageReceivedPayload struct {
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
	Role           string `json:"role"`
}

type ConversationClosedPayload struct {
	ConversationID string `json:"conversation_id"`
}
