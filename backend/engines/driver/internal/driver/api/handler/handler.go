package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/valueobject"
)

type Handler struct {
	service port.DriverService
	logger  *slog.Logger
}

func New(service port.DriverService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"driver-engine"}`))
}

func (h *Handler) GetDriver(w http.ResponseWriter, r *http.Request) {
	driverID := valueobject.DriverID(r.PathValue("driver_id"))
	d, err := h.service.GetDriver(r.Context(), driverID)
	if err != nil {
		http.Error(w, "driver not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (h *Handler) GetVehicles(w http.ResponseWriter, r *http.Request) {
	driverID := valueobject.DriverID(r.PathValue("driver_id"))
	vehicles, err := h.service.GetVehicles(r.Context(), driverID)
	if err != nil {
		http.Error(w, "vehicles not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, vehicles)
}

func (h *Handler) GetLicenses(w http.ResponseWriter, r *http.Request) {
	driverID := valueobject.DriverID(r.PathValue("driver_id"))
	licenses, err := h.service.GetLicenses(r.Context(), driverID)
	if err != nil {
		http.Error(w, "licenses not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, licenses)
}

func (h *Handler) GetAvailability(w http.ResponseWriter, r *http.Request) {
	driverID := valueobject.DriverID(r.PathValue("driver_id"))
	avail, err := h.service.GetAvailability(r.Context(), driverID)
	if err != nil {
		http.Error(w, "availability not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, avail)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
