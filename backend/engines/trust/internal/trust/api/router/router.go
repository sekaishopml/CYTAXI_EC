package router

import (
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/api/handler"
)

func New(mux *http.ServeMux, h *handler.Handler) http.Handler {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /identity/{identity_id}", h.GetIdentity)
	mux.HandleFunc("GET /identity/{identity_id}/trust", h.GetTrustScore)
	return mux
}
