package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/infrastructure/providers"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/infrastructure/webhook"
)

type PaymentRecord struct {
	ID            string    `json:"id"`
	TripID        string    `json:"trip_id"`
	CustomerID    string    `json:"customer_id"`
	DriverID      string    `json:"driver_id"`
	Amount        float64   `json:"amount"`
	Method        string    `json:"method"`
	Status        string    `json:"status"`
	ProviderRef   string    `json:"provider_ref,omitempty"`
	ReceiptNumber string    `json:"receipt_number,omitempty"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

type PaymentServer struct {
	registry    *providers.ProviderRegistry
	webhook     *webhook.Receiver
	logger      *slog.Logger
	payments    sync.Map
	receipts    sync.Map
	intents     sync.Map
}

func NewPaymentServer(registry *providers.ProviderRegistry, webhookReceiver *webhook.Receiver, logger *slog.Logger) *PaymentServer {
	s := &PaymentServer{
		registry: registry,
		webhook:  webhookReceiver,
		logger:   logger,
	}
	s.setupWebhooks()
	return s
}

func (s *PaymentServer) setupWebhooks() {
	s.webhook.Handle(webhook.EventPaymentSucceeded, func(evt webhook.Event) error {
		s.logger.Info("webhook: payment succeeded", "event_id", evt.ID)
		return nil
	})
	s.webhook.Handle(webhook.EventPaymentFailed, func(evt webhook.Event) error {
		s.logger.Error("webhook: payment failed", "event_id", evt.ID)
		return nil
	})
	s.webhook.Handle(webhook.EventRefundProcessed, func(evt webhook.Event) error {
		s.logger.Info("webhook: refund processed", "event_id", evt.ID)
		return nil
	})
}

func (s *PaymentServer) HandleCreateIntent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
		TripID   string  `json:"trip_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Currency == "" { req.Currency = "USD" }

	provider, _ := s.registry.Get("mock")
	intent, err := provider.CreatePaymentIntent(r.Context(), providers.PaymentIntentRequest{
		Amount: req.Amount, Currency: req.Currency,
		Description: fmt.Sprintf("Trip %s", req.TripID),
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	s.intents.Store(intent.ID, intent)

	writeJSON(w, http.StatusOK, map[string]any{"intent": intent})
}

func (s *PaymentServer) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IntentID      string  `json:"intent_id"`
		Amount        float64 `json:"amount"`
		PaymentMethod string  `json:"payment_method"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	provider, _ := s.registry.Get("mock")
	auth, err := provider.Authorize(r.Context(), providers.AuthorizationRequest{
		IntentID: req.IntentID, Amount: req.Amount, Currency: "USD", PaymentMethod: req.PaymentMethod,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"authorization": auth})
}

func (s *PaymentServer) HandleCapture(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AuthorizationID string  `json:"authorization_id"`
		Amount          float64 `json:"amount"`
		TripID          string  `json:"trip_id"`
		CustomerID      string  `json:"customer_id"`
		DriverID        string  `json:"driver_id"`
		PaymentMethod   string  `json:"payment_method"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	provider, _ := s.registry.Get("mock")
	cap, err := provider.Capture(r.Context(), providers.CaptureRequest{
		AuthorizationID: req.AuthorizationID, Amount: req.Amount,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	paymentID := fmt.Sprintf("pay_%d", time.Now().UnixNano())
	receiptNum := fmt.Sprintf("RCP-%d", time.Now().Unix())

	record := PaymentRecord{
		ID: paymentID, TripID: req.TripID, CustomerID: req.CustomerID,
		DriverID: req.DriverID, Amount: req.Amount, Method: req.PaymentMethod,
		Status: "paid", ProviderRef: cap.GatewayRef, ReceiptNumber: receiptNum,
		CreatedAt: time.Now().UTC().Format(time.RFC3339), UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	s.payments.Store(paymentID, record)

	receipt := map[string]any{
		"id": receiptNum, "payment_id": paymentID, "trip_id": req.TripID,
		"amount": req.Amount, "method": req.PaymentMethod,
		"status": "paid", "date": time.Now().UTC().Format(time.RFC3339),
	}
	s.receipts.Store(receiptNum, receipt)

	s.logger.Info("payment captured", "payment_id", paymentID, "amount", req.Amount, "provider", cap.GatewayRef)

	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "captured",
		"payment":   record,
		"receipt":   receipt,
		"capture":   cap,
	})
}

func (s *PaymentServer) HandleRefund(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PaymentID string  `json:"payment_id"`
		Amount    float64 `json:"amount"`
		Reason    string  `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	v, ok := s.payments.Load(req.PaymentID)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "payment not found"})
		return
	}
	record := v.(PaymentRecord)
	record.Status = "refunded"
	record.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	s.payments.Store(req.PaymentID, record)

	provider, _ := s.registry.Get("mock")
	ref, err := provider.Refund(r.Context(), providers.RefundRequest{
		PaymentID: req.PaymentID, Amount: req.Amount, Reason: req.Reason,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"status": "refunded", "refund": ref})
}

func (s *PaymentServer) HandleGetPayment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("payment_id")
	v, ok := s.payments.Load(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (s *PaymentServer) HandleGetReceipt(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("receipt_id")
	v, ok := s.receipts.Load(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (s *PaymentServer) HandleHistory(w http.ResponseWriter, r *http.Request) {
	var payments []PaymentRecord
	s.payments.Range(func(_, v any) bool {
		payments = append(payments, v.(PaymentRecord))
		return true
	})
	writeJSON(w, http.StatusOK, map[string]any{"payments": payments, "total": len(payments)})
}

func (s *PaymentServer) HandleSettlements(w http.ResponseWriter, r *http.Request) {
	driverID := r.PathValue("driver_id")
	var total, commission float64
	var completed int
	s.payments.Range(func(_, v any) bool {
		rec := v.(PaymentRecord)
		if rec.DriverID == driverID && rec.Status == "paid" {
			total += rec.Amount
			commission += rec.Amount * 0.15
			completed++
		}
		return true
	})
	writeJSON(w, http.StatusOK, map[string]any{
		"driver_id":       driverID,
		"trips_completed":  completed,
		"total_earnings":   total,
		"platform_fee":     commission,
		"net_earnings":     total - commission,
	})
}

func (s *PaymentServer) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	signature := r.Header.Get("X-Webhook-Signature")
	var evt webhook.Event
	json.NewDecoder(r.Body).Decode(&evt)
	evt.Signature = signature
	evt.Timestamp = time.Now()

	if err := s.webhook.Process(evt); err != nil {
		s.logger.Error("webhook process failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "webhook processing failed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"status": "received"})
}

func (s *PaymentServer) HandleProviders(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"providers": s.registry.List()})
}

func (s *PaymentServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	providers := s.registry.List()
	available := 0
	for _, name := range providers {
		p, _ := s.registry.Get(name)
		if p.IsAvailable(r.Context()) {
			available++
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":               "ok",
		"service":              "payment-engine",
		"providers_registered": len(providers),
		"providers_available":  available,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
