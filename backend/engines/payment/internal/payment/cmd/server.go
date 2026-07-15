package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
)

type PaymentRecord struct {
	ID          string  `json:"id"`
	TripID      string  `json:"trip_id"`
	CustomerID  string  `json:"customer_id"`
	DriverID    string  `json:"driver_id"`
	Amount      float64 `json:"amount"`
	Method      string  `json:"method"`
	Status      string  `json:"status"`
	FareDetails FeeBreakdown `json:"fare_details"`
	CreatedAt   string  `json:"created_at"`
}

type FeeBreakdown struct {
	Base       float64 `json:"base"`
	Distance   float64 `json:"distance"`
	Time       float64 `json:"time"`
	Subtotal   float64 `json:"subtotal"`
	Tax        float64 `json:"tax"`
	Commission float64 `json:"commission_pct"`
	Total      float64 `json:"total"`
	Currency   string  `json:"currency"`
}

type Receipt struct {
	ID         string       `json:"id"`
	PaymentID  string       `json:"payment_id"`
	TripID     string       `json:"trip_id"`
	CustomerID string       `json:"customer_id"`
	DriverID   string       `json:"driver_id"`
	Amount     float64      `json:"amount"`
	Method     string       `json:"method"`
	FareDetails FeeBreakdown `json:"fare_details"`
	Date       string       `json:"date"`
	Status     string       `json:"status"`
}

type PaymentServer struct {
	service  port.PaymentService
	logger   *slog.Logger
	payments sync.Map
	receipts sync.Map
}

func NewPaymentServer(service port.PaymentService, logger *slog.Logger) *PaymentServer {
	return &PaymentServer{service: service, logger: logger}
}

func (s *PaymentServer) HandleCreatePayment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TripID     string  `json:"trip_id"`
		CustomerID string  `json:"customer_id"`
		DriverID   string  `json:"driver_id"`
		DistanceKM float64 `json:"distance_km"`
		DurationSec int   `json:"duration_sec"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// Calculate fare
	base := 1.00
	dist := req.DistanceKM * 0.50
	timeChg := float64(req.DurationSec) * 0.02
	subtotal := base + dist + timeChg
	tax := subtotal * 0.12
	total := subtotal + tax

	id := fmt.Sprintf("pay_%d", time.Now().UnixNano())

	record := PaymentRecord{
		ID: id, TripID: req.TripID, CustomerID: req.CustomerID, DriverID: req.DriverID,
		Amount: total, Method: "pending", Status: "pending",
		FareDetails: FeeBreakdown{
			Base: base, Distance: dist, Time: timeChg, Subtotal: subtotal,
			Tax: tax, Commission: 15, Total: total, Currency: "USD",
		},
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	s.payments.Store(id, record)

	s.logger.Info("payment created", "payment_id", id, "trip_id", req.TripID, "total", total)

	writeJSON(w, http.StatusOK, map[string]any{
		"status": "pending",
		"payment": record,
	})
}

func (s *PaymentServer) HandleConfirmPayment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PaymentID string `json:"payment_id"`
		Method    string `json:"method"` // card, cash, wallet, transfer
	}
	json.NewDecoder(r.Body).Decode(&req)

	v, ok := s.payments.Load(req.PaymentID)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "payment not found"})
		return
	}

	record := v.(PaymentRecord)
	record.Status = "paid"
	record.Method = req.Method
	s.payments.Store(req.PaymentID, record)

	receipt := Receipt{
		ID: fmt.Sprintf("rcp_%d", time.Now().UnixNano()),
		PaymentID: record.ID, TripID: record.TripID, CustomerID: record.CustomerID,
		DriverID: record.DriverID, Amount: record.Amount, Method: req.Method,
		FareDetails: record.FareDetails,
		Date: time.Now().UTC().Format(time.RFC3339), Status: "paid",
	}
	s.receipts.Store(receipt.ID, receipt)

	s.logger.Info("payment confirmed", "payment_id", req.PaymentID, "method", req.Method)

	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "paid",
		"payment": record,
		"receipt": receipt,
	})
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

func (s *PaymentServer) HandlePaymentHistory(w http.ResponseWriter, r *http.Request) {
	var payments []PaymentRecord
	s.payments.Range(func(_, v any) bool {
		payments = append(payments, v.(PaymentRecord))
		return true
	})
	writeJSON(w, http.StatusOK, map[string]any{
		"payments": payments,
		"total":    len(payments),
	})
}

func (s *PaymentServer) HandleRefund(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PaymentID string `json:"payment_id"`
		Reason    string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	v, ok := s.payments.Load(req.PaymentID)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
		return
	}

	record := v.(PaymentRecord)
	record.Status = "refunded"
	s.payments.Store(req.PaymentID, record)

	s.logger.Info("payment refunded", "payment_id", req.PaymentID, "reason", req.Reason)

	writeJSON(w, http.StatusOK, map[string]any{"status": "refunded", "payment_id": req.PaymentID})
}

func (s *PaymentServer) HandleDriverEarnings(w http.ResponseWriter, r *http.Request) {
	driverID := r.PathValue("driver_id")

	var total, commission float64
	var completed int
	s.payments.Range(func(_, v any) bool {
		r := v.(PaymentRecord)
		if r.DriverID == driverID && r.Status == "paid" {
			total += r.Amount
			commission += r.Amount * r.FareDetails.Commission / 100
			completed++
		}
		return true
	})

	writeJSON(w, http.StatusOK, map[string]any{
		"driver_id":      driverID,
		"trips_completed": completed,
		"total_earnings":  total,
		"platform_fee":    commission,
		"net_earnings":    total - commission,
	})
}

func (s *PaymentServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "payment-engine"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
