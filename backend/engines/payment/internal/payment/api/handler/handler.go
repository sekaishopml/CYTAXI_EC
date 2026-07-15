package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
)

type Handler struct {
	service port.PaymentService
	logger  *slog.Logger
}

func New(service port.PaymentService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"payment-engine"}`))
}

func (h *Handler) GetPayment(w http.ResponseWriter, r *http.Request) {
	id := valueobject.PaymentID(r.PathValue("payment_id"))
	p, err := h.service.GetPayment(r.Context(), query.GetPayment{PaymentID: id})
	if err != nil {
		http.Error(w, "payment not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) GetWallet(w http.ResponseWriter, r *http.Request) {
	ownerID := r.PathValue("owner_id")
	w2, err := h.service.GetWallet(r.Context(), query.GetWallet{OwnerID: ownerID})
	if err != nil {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, w2)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
