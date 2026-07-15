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

type SessionUseCase struct {
	sessionRepo repository.SessionRepository
	convRepo    repository.ConversationRepository
	ctxRepo     repository.ConversationContextRepository
	logger      *slog.Logger
	eventHandler func(event any)
}

func NewSessionUseCase(
	sessionRepo repository.SessionRepository,
	convRepo repository.ConversationRepository,
	ctxRepo repository.ConversationContextRepository,
	logger *slog.Logger,
) *SessionUseCase {
	return &SessionUseCase{
		sessionRepo: sessionRepo,
		convRepo:    convRepo,
		ctxRepo:     ctxRepo,
		logger:      logger,
	}
}

func (uc *SessionUseCase) SetEventHandler(handler func(event any)) {
	uc.eventHandler = handler
}

func (uc *SessionUseCase) CreateSession(ctx context.Context, phone string) (*entity.Session, error) {
	session, err := uc.sessionRepo.FindByPhone(ctx, phone)
	if err == nil && session != nil && session.Status == entity.SessionStatusActive {
		session.Touch()
		uc.sessionRepo.Save(ctx, session)
		return session, nil
	}

	conv, err := uc.convRepo.FindByPhone(ctx, phone)
	if err != nil {
		conv = &entity.Conversation{
			ID:     entity.ConversationID(entity.NewID()),
			Phone:  phone,
			Status: entity.ConversationStatusActive,
		}
		if err := uc.convRepo.Save(ctx, conv); err != nil {
			return nil, fmt.Errorf("create conversation: %w", err)
		}
	}

	session = entity.NewSession(phone)
	session.ConversationID = conv.ID

	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	conversationCtx := entity.NewConversationContext(conv.ID)
	if err := uc.ctxRepo.Save(ctx, conversationCtx); err != nil {
		return nil, fmt.Errorf("create context: %w", err)
	}

	uc.emitEvent(event.SessionCreatedEvent{
		SessionID:      string(session.ID),
		ConversationID: string(session.ConversationID),
		Phone:          session.Phone,
		CreatedAt:      session.CreatedAt,
	})

	return session, nil
}

func (uc *SessionUseCase) ExpireSession(ctx context.Context, sessionID entity.SessionID) error {
	session, err := uc.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	session.MarkExpired()

	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return fmt.Errorf("save expired session: %w", err)
	}

	conv, _ := uc.convRepo.FindByID(ctx, session.ConversationID)
	if conv != nil {
		conv.Status = entity.ConversationStatusClosed
		uc.convRepo.Save(ctx, conv)
	}

	uc.emitEvent(event.SessionExpiredEvent{
		SessionID:      string(session.ID),
		ConversationID: string(session.ConversationID),
		ExpiredAt:      time.Now(),
	})

	return nil
}

func (uc *SessionUseCase) ExpireIdleSessions(ctx context.Context) (int, error) {
	sessions, err := uc.sessionRepo.FindExpired(ctx)
	if err != nil {
		return 0, fmt.Errorf("find expired sessions: %w", err)
	}

	var expired int
	for _, s := range sessions {
		if s.IsExpired() && s.Status == entity.SessionStatusActive {
			if err := uc.ExpireSession(ctx, s.ID); err != nil {
				uc.logger.Error("expire session failed", "session_id", s.ID, "error", err)
				continue
			}
			expired++
		}
	}

	return expired, nil
}

func (uc *SessionUseCase) GetActiveSession(ctx context.Context, phone string) (*entity.Session, error) {
	session, err := uc.sessionRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, nil
	}

	if session.IsExpired() {
		session.MarkExpired()
		uc.sessionRepo.Save(ctx, session)
		return nil, nil
	}

	session.Touch()
	uc.sessionRepo.Save(ctx, session)
	return session, nil
}

func (uc *SessionUseCase) TouchSession(ctx context.Context, sessionID entity.SessionID) error {
	session, err := uc.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return err
	}
	session.Touch()
	return uc.sessionRepo.Save(ctx, session)
}

func (uc *SessionUseCase) emitEvent(evt any) {
	if uc.eventHandler != nil {
		uc.eventHandler(evt)
	}
}
