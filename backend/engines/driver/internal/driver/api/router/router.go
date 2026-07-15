package router

import (
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/api/handler"
)

func New(mux *http.ServeMux, h *handler.Handler) http.Handler {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /drivers/{driver_id}", h.GetDriver)
	mux.HandleFunc("GET /drivers/{driver_id}/vehicles", h.GetVehicles)
	mux.HandleFunc("GET /drivers/{driver_id}/licenses", h.GetLicenses)
	mux.HandleFunc("GET /drivers/{driver_id}/availability", h.GetAvailability)
	return mux
}
