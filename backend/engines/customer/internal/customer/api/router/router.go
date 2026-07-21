package router

import (
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/api/handler"
)

func New(mux *http.ServeMux, h *handler.Handler) http.Handler {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /customers/{customer_id}/profile", h.GetProfile)
	mux.HandleFunc("GET /customers/{customer_id}/preferences", h.GetPreferences)
	mux.HandleFunc("GET /customers/{customer_id}/favorites", h.GetFavorites)
	mux.HandleFunc("GET /customers/{customer_id}/context", h.GetContext)
	return mux
}
