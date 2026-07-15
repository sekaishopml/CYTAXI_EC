package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/valueobject"
)

type DriverRequest struct {
	ID         string    `json:"id"`
	TripID     string    `json:"trip_id"`
	Pickup     string    `json:"pickup"`
	Destination string   `json:"destination"`
	Fare       string    `json:"fare"`
	ETA        int       `json:"eta_seconds"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	Status     string    `json:"status"` // pending, accepted, rejected, expired
}

type DriverServer struct {
	service  port.DriverService
	logger   *slog.Logger
	requests sync.Map
	status   string // online, offline, busy
	mu       sync.RWMutex
}

func NewDriverServer(service port.DriverService, logger *slog.Logger) *DriverServer {
	ds := &DriverServer{
		service: service,
		logger:  logger,
		status:  "online",
	}

	ds.requests.Store("req_001", DriverRequest{
		ID:          "req_001",
		TripID:      "trip_demo",
		Pickup:      "Centro Comercial",
		Destination: "Aeropuerto",
		Fare:        "$12.50",
		ETA:         300,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(30 * time.Second),
		Status:      "pending",
	})

	return ds
}

func (s *DriverServer) HandleGetRequests(w http.ResponseWriter, r *http.Request) {
	var requests []DriverRequest
	s.requests.Range(func(_, v any) bool {
		req := v.(DriverRequest)
		if req.Status == "pending" && time.Now().Before(req.ExpiresAt) {
			requests = append(requests, req)
		}
		return true
	})

	if len(requests) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{
			"requests": []DriverRequest{},
			"status":   s.status,
			"message":  "No pending requests",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"requests": requests,
		"status":   s.status,
	})
}

func (s *DriverServer) HandleAccept(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RequestID string `json:"request_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	v, ok := s.requests.Load(req.RequestID)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "request not found"})
		return
	}

	dr := v.(DriverRequest)
	dr.Status = "accepted"
	s.requests.Store(req.RequestID, dr)

	s.mu.Lock()
	s.status = "busy"
	s.mu.Unlock()

	s.logger.Info("driver accepted trip", "request_id", req.RequestID)

	writeJSON(w, http.StatusOK, map[string]any{
		"status":     "accepted",
		"request_id": req.RequestID,
		"trip_id":    dr.TripID,
		"driver": map[string]any{
			"id":     "drv_1000",
			"name":   "Carlos M.",
			"vehicle": "Toyota Corolla",
			"plate":  "ABC-1234",
			"eta":    180,
		},
	})
}

func (s *DriverServer) HandleReject(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RequestID string `json:"request_id"`
		Reason    string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	v, ok := s.requests.Load(req.RequestID)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "request not found"})
		return
	}

	dr := v.(DriverRequest)
	dr.Status = "rejected"
	s.requests.Store(req.RequestID, dr)

	s.logger.Info("driver rejected trip", "request_id", req.RequestID, "reason", req.Reason)

	writeJSON(w, http.StatusOK, map[string]any{
		"status":     "rejected",
		"request_id": req.RequestID,
		"message":    "Trip request rejected",
	})
}

func (s *DriverServer) HandleStatus(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	st := s.status
	s.mu.RUnlock()

	writeJSON(w, http.StatusOK, map[string]any{
		"driver_id": "drv_1000",
		"status":    st,
		"name":      "Carlos M.",
		"rating":    4.8,
	})
}

func (s *DriverServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "driver-engine"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
