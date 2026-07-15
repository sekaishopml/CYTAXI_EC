package router

import (
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/api/handler"
)

func New(mux *http.ServeMux, handlers *handler.Handler) http.Handler {
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("POST /messages/incoming", handlers.IncomingMessage)
	mux.HandleFunc("/", handlers.NotFound)
	return mux
}

func NewFromInput(logger *slog.Logger, mux *http.ServeMux) http.Handler {
	return nil
}
