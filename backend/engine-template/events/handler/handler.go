package handler

import "github.com/sekaishopml/cytaxi/backend/engine-template/domain/event"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(evt event.DomainEvent) error {
	return nil
}
