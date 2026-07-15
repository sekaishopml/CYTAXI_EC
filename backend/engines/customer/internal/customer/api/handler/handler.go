package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/application/port"
)

type Handler struct {
	service port.CustomerService
	logger  *slog.Logger
}

func New(service port.CustomerService, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"customer-engine"}`))
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customer_id")
	profile, err := h.service.GetProfile(r.Context(), customerID)
	if err != nil {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func (h *Handler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customer_id")
	prefs, err := h.service.GetPreferences(r.Context(), customerID)
	if err != nil {
		http.Error(w, "preferences not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, prefs)
}

func (h *Handler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customer_id")
	favorites, err := h.service.GetFavorites(r.Context(), customerID)
	if err != nil {
		http.Error(w, "favorites not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, favorites)
}

func (h *Handler) GetContext(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customer_id")
	ctx, err := h.service.GetContext(r.Context(), customerID)
	if err != nil {
		http.Error(w, "context not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, ctx)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
