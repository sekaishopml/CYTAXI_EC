package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/valueobject"
)

type Handler struct {
	service port.PricingService
	logger  *slog.Logger
}

func New(service port.PricingService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"pricing-engine"}`))
}

func (h *Handler) GetFare(w http.ResponseWriter, r *http.Request) {
	fareID := valueobject.FareID(r.PathValue("fare_id"))
	f, err := h.service.GetFare(r.Context(), query.GetFare{FareID: fareID})
	if err != nil {
		http.Error(w, "fare not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, f)
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	tripID := r.PathValue("trip_id")
	fares, err := h.service.GetFareHistory(r.Context(), query.GetFareHistory{TripID: tripID})
	if err != nil {
		http.Error(w, "history not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, fares)
}

func (h *Handler) GetPromotions(w http.ResponseWriter, r *http.Request) {
	promos, err := h.service.GetPromotions(r.Context(), query.GetPromotions{Active: true})
	if err != nil {
		http.Error(w, "promotions not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, promos)
}

// EstimateRequest body esperado desde el frontend
type EstimateRequest struct {
	TripID      string  `json:"trip_id"`
	DistanceKM  float64 `json:"distance_km"`
	DurationSec int     `json:"duration_sec"`
	WaitingSec  int     `json:"waiting_sec"`
	Region      string  `json:"region"`
	IsNight     bool    `json:"is_night"`
	DemandLevel int     `json:"demand_level"`
}

// HandleEstimate POST /pricing/estimate — Calcula tarifa estimada
func (h *Handler) HandleEstimate(w http.ResponseWriter, r *http.Request) {
	var req EstimateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json", "details": err.Error()})
		return
	}

	if req.TripID == "" {
		req.TripID = "preview_trip"
	}
	if req.Region == "" {
		req.Region = "ecuador"
	}
	if req.DistanceKM <= 0 {
		req.DistanceKM = 0.01
	}

	// Llamar al servicio
	f, err := h.service.CalculateFare(r.Context(), command.CalculateFare{
		TripID:      req.TripID,
		DistanceKM:  req.DistanceKM,
		DurationSec: req.DurationSec,
		WaitingSec:  req.WaitingSec,
		IsNight:     req.IsNight,
		DemandLevel: req.DemandLevel,
		Region:      req.Region,
	})
	if err != nil {
		h.logger.Error("estimate failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "could not estimate fare", "details": err.Error()})
		return
	}

	// Devolver respuesta esperada por el frontend
	c := f.Components
	writeJSON(w, http.StatusOK, map[string]any{
		"fare_id":     string(f.ID),
		"trip_id":     f.TripID,
		"currency":    f.Currency,
		"base":        c.BaseFare.Amount,
		"distance":    c.DistanceFare.Amount,
		"time":        c.TimeFare.Amount,
		"waiting":     c.WaitingFare.Amount,
		"night":       c.NightSurcharge.Amount,
		"demand":      c.DemandSurcharge.Amount,
		"subtotal":    c.Subtotal.Amount,
		"tax":         c.Tax.Amount,
		"discount":    c.Promotion.Amount + c.Coupon.Amount,
		"total":       c.Total.Amount,
		"driver_earn": c.DriverEarnings.Amount,
		"distance_km": req.DistanceKM,
		"duration_sec": req.DurationSec,
		"eta_minutes": req.DurationSec / 60,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
