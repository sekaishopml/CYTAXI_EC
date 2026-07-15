package repository

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/domain/entity"
)

type ConversationRepository interface {
	FindByID(ctx context.Context, id entity.ConversationID) (*entity.Conversation, error)
	FindByPhone(ctx context.Context, phone string) (*entity.Conversation, error)
	Save(ctx context.Context, conversation *entity.Conversation) error
}

type MessageRepository interface {
	FindByConversation(ctx context.Context, conversationID entity.ConversationID) ([]entity.Message, error)
	Save(ctx context.Context, message *entity.Message) error
}

type SessionRepository interface {
	FindByID(ctx context.Context, id entity.SessionID) (*entity.Session, error)
	FindByPhone(ctx context.Context, phone string) (*entity.Session, error)
	FindByConversation(ctx context.Context, convID entity.ConversationID) (*entity.Session, error)
	FindExpired(ctx context.Context) ([]entity.Session, error)
	Save(ctx context.Context, session *entity.Session) error
	Delete(ctx context.Context, id entity.SessionID) error
}

type ConversationContextRepository interface {
	FindByConversation(ctx context.Context, convID entity.ConversationID) (*entity.ConversationContext, error)
	Save(ctx context.Context, ctxData *entity.ConversationContext) error
	Delete(ctx context.Context, convID entity.ConversationID) error
}
