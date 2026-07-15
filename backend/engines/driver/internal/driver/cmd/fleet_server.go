package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/infrastructure/fleet"
)

type FleetServer struct {
	manager *fleet.Manager
	logger  *slog.Logger
}

func NewFleetServer(manager *fleet.Manager, logger *slog.Logger) *FleetServer {
	m := manager
	m.RegisterVehicle("", "ABC-1234", "Toyota", "Corolla", 2023, "standard")
	m.RegisterVehicle("", "XYZ-5678", "Hyundai", "Tucson", 2022, "xl")
	return &FleetServer{manager: m, logger: logger}
}

func (s *FleetServer) HandleFleet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req struct {
			Name        string `json:"name"`
			OwnerID     string `json:"owner_id"`
			Description string `json:"description"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		result := s.manager.CreateFleet(req.Name, req.OwnerID, req.Description)
		writeJSON(w, http.StatusCreated, result)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (s *FleetServer) HandleVehicles(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req struct {
			FleetID string `json:"fleet_id"`
			Plate   string `json:"plate"`
			Brand   string `json:"brand"`
			Model   string `json:"model"`
			Year    int    `json:"year"`
			Type    string `json:"type"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		result := s.manager.RegisterVehicle(req.FleetID, req.Plate, req.Brand, req.Model, req.Year, req.Type)
		s.logger.Info("vehicle registered", "plate", req.Plate, "fleet", req.FleetID)
		writeJSON(w, http.StatusCreated, result)
		return
	}

	vehicles := s.manager.GetVehicles("")
	writeJSON(w, http.StatusOK, map[string]any{"vehicles": vehicles, "count": len(vehicles)})
}

func (s *FleetServer) HandleAssignments(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VehicleID string `json:"vehicle_id"`
		DriverID  string `json:"driver_id"`
		Action    string `json:"action"` // assign, release
		Notes     string `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Action == "release" {
		s.manager.ReleaseDriver(req.VehicleID)
		s.logger.Info("driver released", "vehicle", req.VehicleID)
		writeJSON(w, http.StatusOK, map[string]any{"status": "released"})
		return
	}

	result := s.manager.AssignDriver(req.VehicleID, req.DriverID, req.Notes)
	s.logger.Info("driver assigned", "vehicle", req.VehicleID, "driver", req.DriverID)
	writeJSON(w, http.StatusCreated, result)
}

func (s *FleetServer) HandleMaintenance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VehicleID     string  `json:"vehicle_id"`
		Type          string  `json:"type"`
		Description   string  `json:"description"`
		Action        string  `json:"action"` // schedule, complete
		MaintenanceID string  `json:"maintenance_id,omitempty"`
		Cost          float64 `json:"cost,omitempty"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Action == "complete" {
		s.manager.CompleteMaintenance(req.VehicleID, req.MaintenanceID, req.Cost)
		writeJSON(w, http.StatusOK, map[string]any{"status": "completed"})
		return
	}

	result := s.manager.ScheduleMaintenance(req.VehicleID, req.Type, req.Description)
	s.logger.Info("maintenance scheduled", "vehicle", req.VehicleID, "type", req.Type)
	writeJSON(w, http.StatusCreated, result)
}

func (s *FleetServer) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	dash := s.manager.GetDashboard()
	writeJSON(w, http.StatusOK, dash)
}

func (s *FleetServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "fleet-management"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
