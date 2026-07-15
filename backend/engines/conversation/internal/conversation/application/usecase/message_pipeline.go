package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/domain/entity"
	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/domain/event"
	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/domain/repository"
)

type MessagePipeline struct {
	sessionUC  *SessionUseCase
	convRepo   repository.ConversationRepository
	msgRepo    repository.MessageRepository
	ctxRepo    repository.ConversationContextRepository
	logger     *slog.Logger
	processors []MessageProcessor
}

type MessageProcessor interface {
	Process(ctx context.Context, msg *entity.Message, session *entity.Session) error
}

func NewMessagePipeline(
	sessionUC *SessionUseCase,
	convRepo repository.ConversationRepository,
	msgRepo repository.MessageRepository,
	ctxRepo repository.ConversationContextRepository,
	logger *slog.Logger,
) *MessagePipeline {
	return &MessagePipeline{
		sessionUC: sessionUC,
		convRepo:  convRepo,
		msgRepo:   msgRepo,
		ctxRepo:   ctxRepo,
		logger:    logger,
	}
}

func (p *MessagePipeline) RegisterProcessor(processor MessageProcessor) {
	p.processors = append(p.processors, processor)
}

func (p *MessagePipeline) ProcessIncoming(ctx context.Context, phone string, content string) error {
	session, err := p.sessionUC.GetActiveSession(ctx, phone)
	if err != nil {
		session, err = p.sessionUC.CreateSession(ctx, phone)
		if err != nil {
			return fmt.Errorf("create session: %w", err)
		}
	}

	if session == nil {
		session, err = p.sessionUC.CreateSession(ctx, phone)
		if err != nil {
			return fmt.Errorf("create session: %w", err)
		}
	}

	if err := entity.ValidateTransition(session.CurrentState, entity.StateProcessing); err != nil {
		p.logger.Warn("invalid state transition",
			"session", session.ID,
			"current_state", session.CurrentState,
		)
	}

	prevState := session.CurrentState
	session.CurrentState = entity.StateProcessing
	session.Touch()

	msg := &entity.Message{
		ID:             entity.MessageID(entity.NewID()),
		ConversationID: session.ConversationID,
		Content:        content,
		Role:           entity.MessageRoleUser,
	}

	if err := p.msgRepo.Save(ctx, msg); err != nil {
		return fmt.Errorf("save message: %w", err)
	}

	p.emitEvent(event.MessageReceivedEvent{
		MessageID:      string(msg.ID),
		ConversationID: string(msg.ConversationID),
		Content:        msg.Content,
		Role:           string(msg.Role),
		ReceivedAt:     time.Now(),
	})

	if prevState != session.CurrentState {
		p.emitEvent(event.StateChangedEvent{
			ConversationID: string(session.ConversationID),
			PreviousState:  string(prevState),
			NewState:       string(session.CurrentState),
			ChangedAt:      time.Now(),
		})
	}

	for _, proc := range p.processors {
		if err := proc.Process(ctx, msg, session); err != nil {
			p.logger.Error("processor failed", "error", err)
		}
	}

	session.CurrentState = entity.StateWaitingInput
	p.sessionUC.TouchSession(ctx, session.ID)

	return nil
}

func (p *MessagePipeline) emitEvent(evt any) {
	if p.sessionUC != nil {
		p.sessionUC.SetEventHandler(func(e any) {})
	}
}
