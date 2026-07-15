package ai

import (
	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/domain/entity"
)

type Context struct {
	SessionID      string
	ConversationID string
	Phone          string
	State          string
	Entries        map[string]string
	MessageHistory []MessageSummary
}

type MessageSummary struct {
	Role    string
	Content string
}

func BuildContext(session *entity.Session, convCtx *entity.ConversationContext) Context {
	ctx := Context{
		SessionID:      string(session.ID),
		ConversationID: string(session.ConversationID),
		Phone:          session.Phone,
		State:          string(session.CurrentState),
		Entries:        make(map[string]string),
	}

	if convCtx != nil {
		for k, entry := range convCtx.Entries {
			ctx.Entries[k] = entry.Value
		}
	}

	return ctx
}

type ContextEnricher interface {
	Enrich(ctx Context, session *entity.Session, convCtx *entity.ConversationContext) Context
}

type StaticContextEnricher struct {
	data map[string]string
}

func NewStaticContextEnricher(data map[string]string) *StaticContextEnricher {
	return &StaticContextEnricher{data: data}
}

func (e *StaticContextEnricher) Enrich(ctx Context, session *entity.Session, convCtx *entity.ConversationContext) Context {
	for k, v := range e.data {
		if _, exists := ctx.Entries[k]; !exists {
			ctx.Entries[k] = v
		}
	}
	return ctx
}
