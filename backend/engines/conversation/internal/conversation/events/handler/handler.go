package handler

import (
	"log/slog"
)

type EventHandler struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *EventHandler {
	return &EventHandler{logger: logger}
}
