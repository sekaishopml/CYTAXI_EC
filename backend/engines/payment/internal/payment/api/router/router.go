package router

import (
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/api/handler"
)

func New(mux *http.ServeMux, h *handler.Handler) http.Handler {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /payments/{payment_id}", h.GetPayment)
	mux.HandleFunc("GET /wallets/{owner_id}", h.GetWallet)
	return mux
}
