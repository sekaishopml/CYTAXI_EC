package port

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/domain/entity"
)

type MessageInputPort interface {
	HandleIncomingMessage(ctx context.Context, phone string, content string) error
}

type ConversationOutputPort interface {
	SendMessage(ctx context.Context, conversationID entity.ConversationID, content string) error
}
