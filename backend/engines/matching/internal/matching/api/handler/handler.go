package handler

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/valueobject"
)

type Handler struct {
	service port.MatchingService
	logger  *slog.Logger
}

func New(service port.MatchingService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"matching-engine"}`))
}

func (h *Handler) GetMatching(w http.ResponseWriter, r *http.Request) {
	id := valueobject.MatchingID(r.PathValue("matching_id"))
	m, err := h.service.GetMatching(r.Context(), query.GetMatching{MatchingID: id})
	if err != nil {
		http.Error(w, "matching not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (h *Handler) GetCandidates(w http.ResponseWriter, r *http.Request) {
	id := valueobject.MatchingID(r.PathValue("matching_id"))
	cs, err := h.service.GetCandidates(r.Context(), query.GetCandidates{MatchingID: id})
	if err != nil {
		http.Error(w, "candidates not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, cs)
}

// StartMatchingRequest body esperado desde el frontend
type StartMatchingRequest struct {
	TripID    string  `json:"trip_id"`
	PickupLat float64 `json:"pickup_lat"`
	PickupLng float64 `json:"pickup_lng"`
	Strategy  string  `json:"strategy"`
}

// Candidate respuesta hacia el frontend
type Candidate struct {
	DriverID       string  `json:"driver_id"`
	Name           string  `json:"name"`
	DistanceMeters float64 `json:"distance_meters"`
	ETASeconds     int     `json:"eta_seconds"`
	Score          float64 `json:"score"`
	Vehicle        string  `json:"vehicle"`
	Plate          string  `json:"plate"`
	Rating         float64 `json:"rating"`
	Lat            float64 `json:"lat"`
	Lng            float64 `json:"lng"`
}

// HandleStartMatching POST /matching/start — Inicia matching y devuelve candidatos
func (h *Handler) HandleStartMatching(w http.ResponseWriter, r *http.Request) {
	var req StartMatchingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json", "details": err.Error()})
		return
	}
	if req.TripID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "trip_id required"})
		return
	}

	strategy := req.Strategy
	if strategy == "" {
		strategy = "balanced"
	}

	// Iniciar matching
	m, err := h.service.StartMatching(r.Context(), command.StartMatching{
		TripID:    valueobject.TripID(req.TripID),
		PickupLat: req.PickupLat,
		PickupLng: req.PickupLng,
		Strategy:  strategy,
	})
	if err != nil {
		h.logger.Error("start matching failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "could not start matching", "details": err.Error()})
		return
	}

	// Generar candidatos (mock hasta que el driver engine exponga un endpoint real)
	candidates := generateMockCandidates(req.PickupLat, req.PickupLng, 3)

	writeJSON(w, http.StatusOK, map[string]any{
		"matching_id": string(m.ID),
		"status":      string(m.Status),
		"candidates":  candidates,
		"strategy":    strategy,
	})
}

// HandleSelectDriver POST /matching/select — Asigna un conductor al viaje
func (h *Handler) HandleSelectDriver(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MatchingID string `json:"matching_id"`
		DriverID   string `json:"driver_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}
	if req.MatchingID == "" || req.DriverID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "matching_id and driver_id required"})
		return
	}

	// Aquí se llamaría al servicio de selección. En MVP devolvemos OK.
	writeJSON(w, http.StatusOK, map[string]any{
		"status":     "assigned",
		"driver_id":  req.DriverID,
		"message":    "Driver assigned successfully",
		"trip_id":    req.MatchingID,
	})
}

// generateMockCandidates crea N candidatos cercanos al pickup
func generateMockCandidates(pickupLat, pickupLng float64, n int) []Candidate {
	rng := rand.New(rand.NewSource(int64(pickupLat*1000) + int64(pickupLng*1000)))
	names := []string{"Carlos M.", "Ana P.", "Luis R.", "María G.", "José S.", "Diego T."}
	vehicles := []string{"Toyota Corolla", "Hyundai Accent", "Kia Rio", "Chevrolet Spark", "Nissan Versa"}
	plates := []string{"ABC-1234", "XYZ-5678", "GHI-9012", "DEF-3456", "JKL-7890"}

	out := make([]Candidate, 0, n)
	for i := 0; i < n; i++ {
		// Distancia aleatoria 200m - 3km
		dist := 200.0 + rng.Float64()*2800.0
		eta := int(dist / 8.33) // 30 km/h promedio
		score := 1.0 - (dist / 5000.0) + rng.Float64()*0.1
		if score > 1.0 {
			score = 1.0
		}
		if score < 0.0 {
			score = 0.0
		}
		// Posición cercana (offset aleatorio)
		lat := pickupLat + (rng.Float64()-0.5)*0.02
		lng := pickupLng + (rng.Float64()-0.5)*0.02
		// Rating
		rating := 4.5 + rng.Float64()*0.5

		out = append(out, Candidate{
			DriverID:       "drv_" + randStr(rng, 4),
			Name:           names[rng.Intn(len(names))],
			DistanceMeters: dist,
			ETASeconds:     eta,
			Score:          score,
			Vehicle:        vehicles[rng.Intn(len(vehicles))],
			Plate:          plates[rng.Intn(len(plates))],
			Rating:         rating,
			Lat:            lat,
			Lng:            lng,
		})
	}
	return out
}

func randStr(rng *rand.Rand, n int) string {
	const letters = "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
