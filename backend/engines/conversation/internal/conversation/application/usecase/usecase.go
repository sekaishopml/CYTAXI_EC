package usecase

import (
	"context"
	"fmt"

	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/domain/entity"
	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/domain/repository"
)

type MessageUseCase struct {
	pipeline *MessagePipeline
	convRepo repository.ConversationRepository
}

func NewMessageUseCase(pipeline *MessagePipeline, convRepo repository.ConversationRepository) *MessageUseCase {
	return &MessageUseCase{
		pipeline: pipeline,
		convRepo: convRepo,
	}
}

func (uc *MessageUseCase) HandleIncomingMessage(ctx context.Context, phone string, content string) error {
	return uc.pipeline.ProcessIncoming(ctx, phone, content)
}

func (uc *MessageUseCase) GetConversation(ctx context.Context, id entity.ConversationID) (*entity.Conversation, error) {
	return uc.convRepo.FindByID(ctx, id)
}
